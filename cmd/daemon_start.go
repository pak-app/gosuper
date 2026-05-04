/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os/exec"
	"syscall"
)

// daemonStartCmd represents the daemonStart command
var daemonStartCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting Gosuper Daemon...")

		exePath := "./gosuper"
		// Create a command that runs the hidden "serve" subcommand
		bgCmd := exec.Command(exePath, "daemon", "serve")

		// Detach the process from the terminal (Unix-like systems)
		bgCmd.SysProcAttr = &syscall.SysProcAttr{
			Setsid: true, // Creates a new session, detaching from current terminal
		}

		// Start the command but do NOT wait for it to finish
		if err := bgCmd.Start(); err != nil {
			err = fmt.Errorf("failed to start daemon: %w", err)
			return
		}

		fmt.Printf("Daemon started successfully in background (PID: %d)\n", bgCmd.Process.Pid)
	},
}

func init() {
	daemonCmd.AddCommand(daemonStartCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// daemonStartCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// daemonStartCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
