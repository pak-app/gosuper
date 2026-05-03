/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/pak-app/gosuper/internal/config"
	"github.com/spf13/cobra"
)

var appConfig *config.Config
var cfgFilePath string
var supervisorName string

const defaultConfigFilePath string = "gosuper.yaml"

// serviceCmd represents the service command
var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("service called")
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cfgFilePath == "" {
			cfgFilePath = defaultConfigFilePath
		}

		c, err := config.LoadConfig(cfgFilePath)
		if err != nil {
			log.Println("config file is not available")
			return fmt.Errorf("failed to load config: %w", err)
		}

		appConfig = c
		return nil
	},
}

func init() {
	rootCmd.AddCommand(serviceCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serviceCmd.PersistentFlags().String("foo", "", "A help for foo")
	serviceCmd.PersistentFlags().StringVarP(&cfgFilePath, "config", "c", defaultConfigFilePath, "config file (default is gosuper.yaml)")
	serviceCmd.PersistentFlags().StringVar(&supervisorName, "supervisor-name", "", "name of the service")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serviceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
