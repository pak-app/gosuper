/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/pak-app/gosuper/internal/client"
	"github.com/spf13/cobra"
)

// stopServiceCmd represents the stopService command
var stopServiceCmd = &cobra.Command{
	Use:   "stop",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := client.New("tmp/gosuper.sock")

		if err != nil {
			log.Println(err)
			return
		}

		var name string

		if appConfig.Supervisor.Name != "" {
			name = appConfig.Supervisor.Name
		} else if supervisorName != "" {
			name = supervisorName
		} else {
			log.Println("Group name doesn't set for services")
			return
		}

		err = c.ServiceStopRequest(name)

		if err != nil {
			log.Println(err)
		}
	},
}

func init() {
	serviceCmd.AddCommand(stopServiceCmd)

	stopServiceCmd.Flags().String("service-name", "", "Name of the service to stop")
}
