/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/pak-app/gosuper/internal/client"
	"github.com/spf13/cobra"
)

// startServiceCmd represents the startService command
var startServiceCmd = &cobra.Command{
	Use:   "start",
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

		err = c.ServiceStartRequest(AppConfig)

		if err != nil {
			log.Println(err)
		}
	},
}

func init() {
	serviceCmd.AddCommand(startServiceCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startServiceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startServiceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
