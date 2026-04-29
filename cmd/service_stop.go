/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/pak-app/gosuper/internal/client"
	"github.com/spf13/cobra"
	"log"
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

		if AppConfig.Supervisor.Name != "" {
			name = AppConfig.Supervisor.Name
		} else if serviceName != "" {
			name = serviceName
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
