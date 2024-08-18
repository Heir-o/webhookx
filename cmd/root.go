package cmd

import (
	"github.com/spf13/cobra"
	"github.com/webhookx-io/webhookx/config"
)

var (
	configurationFile string
	cfg               *config.Config

	cmd = &cobra.Command{
		Use:          "webhookx",
		Short:        "",
		Long:         ``,
		SilenceUsage: true,
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	cmd.PersistentFlags().StringVarP(&configurationFile, "config", "", "", "The configuration filename")

	cmd.AddCommand(newVersionCmd())
	cmd.AddCommand(newMigrationsCmd())
	cmd.AddCommand(newStartCmd())
}

func initConfig() {
	var err error
	if configurationFile != "" {
		cfg, err = config.InitWithFile(configurationFile)
	} else {
		cfg, err = config.Init()
	}
	cobra.CheckErr(err)

	err = cfg.Validate()
	cobra.CheckErr(err)
}

func Execute() {
	cmd.Execute()
}
