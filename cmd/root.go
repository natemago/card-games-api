/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"toggl.com/services/card-games-api/app"
	"toggl.com/services/card-games-api/config"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "card-games-api",
	Short: "Card Games REST API",

	RunE: func(cmd *cobra.Command, args []string) error {
		return app.RunApp(&Config)
	},
}

var Config config.Config = config.Config{}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVar(&Config.DBConfig.URL, "db-url", "", "URL to PostgreSQL database.")
	rootCmd.Flags().StringVar(&Config.APIConfig.Host, "bind-host", "", "Bind to hostname.")
	rootCmd.Flags().IntVar(&Config.APIConfig.Port, "bind-port", 8080, "Listen on port.")
}
