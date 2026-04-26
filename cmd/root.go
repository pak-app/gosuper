/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/pak-app/gosuper/internal/config"
	"github.com/spf13/cobra"
)

var AppConfig *config.Config
var cfgFilePath string
var serviceName string
const defaultConfigFilePath string = "gosuper.yaml"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gosuper",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cfgFilePath == "" {
			cfgFilePath = defaultConfigFilePath
		}

		c, err := config.LoadConfig(cfgFilePath)
		if err != nil {
			log.Println("config file is not available")
			return fmt.Errorf("failed to load config: %w", err)
		}

		AppConfig = c
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVarP(&cfgFilePath, "config", "c", defaultConfigFilePath, "config file (default is gosuper.yaml)")
	rootCmd.PersistentFlags().StringVar(&serviceName, "name", "", "name of the service")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
