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

// supervisorCmd represents the service command
var supervisorCmd = &cobra.Command{
	Use:   "supervisor",
	Short: "Manage and monitor background supervisors(group of services)",
	Long: `The supervisor command acts as the central process manager for your application.
It provides subcommands to start, stop, restart, and monitor background services.

Before any subcommand is executed, the supervisor automatically loads the required 
configuration file to determine the state and rules for the processes. 
(Default config path: /etc/myapp/config.yaml)`,

	// Run: func(cmd *cobra.Command, args []string) {
	// },
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
	rootCmd.AddCommand(supervisorCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serviceCmd.PersistentFlags().String("foo", "", "A help for foo")
	supervisorCmd.PersistentFlags().StringVarP(&cfgFilePath, "config", "c", defaultConfigFilePath, "config file")
	supervisorCmd.PersistentFlags().StringVar(&supervisorName, "supervisor-name", "", "name of the service")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serviceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
