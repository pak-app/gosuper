/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/pak-app/gosuper/internal/client"
	"github.com/spf13/cobra"
)

// statusServiceCmd represents the statusService command
var statusServiceCmd = &cobra.Command{
	Use:   "status",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := client.New("tmp/gosuper.sock")

		if err != nil {
			log.Println("Error occured: ", err)
		}

		var supName string

		if appConfig.Supervisor.Name != "" {
			supName = appConfig.Supervisor.Name
		} else if supervisorName != "" {
			supName = supervisorName
		} else {
			log.Println("supervisor name doesn't available")
			return
		}

		supervisorStatus, err := c.ServiceStatusRequest(supName)

		if err != nil {
			log.Println("Error during requesting status of supervisor/supervisors:", err)
		}

		log.Println("Status:\n", supervisorStatus)
	},
}

func init() {
	serviceCmd.AddCommand(statusServiceCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusServiceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusServiceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
