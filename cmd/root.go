package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var (
	cluster string
)

// Configure the root command
var rootCmd = &cobra.Command{
	Use:   "ecs",
	Short: "Manage ECS",
	Long:  `A lightweight tool for working with ECS`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute validates input the Cobra CLI
func Execute(version string) {
	rootCmd.Version = version
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

// Log errors if exist and exit
func check(err error) {
	if err != nil {
		log.Fatalf("ERROR\t%s", err.Error())
	}
}
