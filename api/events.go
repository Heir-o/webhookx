package api

import (
	"context"
	"encoding/json"
	"github.com/webhookx-io/webhookx/db"
	"github.com/webhookx-io/webhookx/db/entities"
	"github.com/webhookx-io/webhookx/db/query"
	"github.com/webhookx-io/webhookx/model"
	"github.com/webhookx-io/webhookx/pkg/queue"
	"github.com/webhookx-io/webhookx/utils"
	"net/http"
)

func (api *API) PageEvent(w http.ResponseWriter, r *http.Request) {
	var q query.EventQuery
	api.bindQuery(r, &q.Query)

	list, total, err := api.DB.Events.Page(r.Context(), &q)
	api.assert(err)

	api.json(200, w, NewPagination(total, list))
}

func (api *API) GetEvent(w http.ResponseWriter, r *http.Request) {
	id := api.param(r, "id")
	event, err := api.DB.Events.Get(r.Context(), id)
	api.assert(err)

	if event == nil {
		api.json(404, w, ErrorResponse{Message: MsgNotFound})
		return
	}

	api.json(200, w, event)
}

func (api *API) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var event entities.Event
	event.ID = utils.UUID()

	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		api.error(400, w, err)
		return
	}

	if err := event.Validate(); err != nil {
		api.error(400, w, err)
		return
	}

	err := DispatchEvent(api, r.Context(), &event)
	api.assert(err)

	api.json(201, w, event)
}

func DispatchEvent(api *API, ctx context.Context, event *entities.Event) error {
	// TODO: db.begin()
	err := api.DB.Events.Insert(ctx, event)
	if err != nil {
		return err
	}

	endpoints, err := listSubscribedEndpoints(ctx, api.DB, event.EventType)
	if err != nil {
		return err
	}

	tasks := make([]*queue.TaskMessage, 0, len(endpoints))
	for _, endpoint := range endpoints {
		task := &queue.TaskMessage{
			ID: utils.UUID(),
			Data: &model.MessageData{
				EventID:     event.ID,
				EndpointId:  endpoint.ID,
				Attempt:     1,
				AttemptLeft: len(endpoint.Retry.Config.Attempts) - 1,
				Delay:       endpoint.Retry.Config.Attempts[0],
			},
		}
		tasks = append(tasks, task)
	}
	// TODO: db.commit()

	for _, task := range tasks {
		err := api.queue.Add(task, utils.DurationS(task.Data.(*model.MessageData).Delay))
		if err != nil {
			return err
		}
	}

	// TODO: record tasks has been successfully queued

	return nil
}

func listSubscribedEndpoints(ctx context.Context, db *db.DB, eventType string) (list []*entities.Endpoint, err error) {
	var q query.EndpointQuery
	endpoints, err := db.Endpoints.List(ctx, &q)
	if err != nil {
		return nil, err
	}

	for _, endpoint := range endpoints {
		if !endpoint.Enabled {
			continue
		}
		for _, event := range endpoint.Events {
			if eventType == event {
				list = append(list, endpoint)
			}
		}
	}

	return
}
