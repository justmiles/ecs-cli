package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	cluster string
)

// Configure the root command
var rootCmd = &cobra.Command{
	Use:     "ecs",
	Short:   "Manage ECS",
	Long:    `A lightweight tool for working with ECS`,
	Version: "0.0.5",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute validates input the Cobra CLI
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Log errors if exist and exit
func check(err error) {
	if err != nil {
		fmt.Printf("ERROR\t%s", err.Error())
		os.Exit(1)
	}
}
