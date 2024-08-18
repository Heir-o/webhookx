package config

import "errors"

type ProxyResponse struct {
	Code        uint   `yaml:"code" default:"200"`
	ContentType string `yaml:"contentType" default:"application/json"`
	Body        string `yaml:"body" default:"{\"message\": \"OK\"}"`
}

type Queue struct {
	Redis RedisConfig `yaml:"redis"`
}

type ProxyConfig struct {
	Listen             string        `yaml:"listen" default:"127.0.0.1:8081"`
	MaxRequestBodySize int64         `yaml:"max_request_body_size" default:"1048576"`
	Response           ProxyResponse `yaml:"response"`
	Queue              Queue         `yaml:"queue"`
}

func (cfg ProxyConfig) Validate() error {
	if cfg.MaxRequestBodySize < 0 {
		return errors.New("max_request_body_size cannot be negative value")
	}
	return nil
}
