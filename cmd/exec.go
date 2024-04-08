package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/AlecAivazis/survey/v2"
	ecs "github.com/justmiles/ecs-cli/lib"
	"github.com/spf13/cobra"
)

var (
	execInput ecs.ExecInput
)

func init() {
	log.SetFlags(0)

	rootCmd.AddCommand(ExecCmd)
	ExecCmd.PersistentFlags().StringVarP(&execInput.Cluster, "cluster", "c", "", "ECS cluster")
	ExecCmd.PersistentFlags().StringVarP(&execInput.Service, "service", "s", "", "ECS service")
	ExecCmd.PersistentFlags().StringVarP(&execInput.Task, "task", "t", "", "ECS task")
	ExecCmd.PersistentFlags().StringVar(&execInput.Container, "container", "", "ECS container")
	ExecCmd.PersistentFlags().StringVar(&execInput.Command, "cmd", "", "ECS container")
	ExecCmd.PersistentFlags().BoolVarP(&execInput.Interactive, "interactive", "i", true, "open interative session")
}

var ExecCmd = &cobra.Command{
	Use:   "exec",
	Short: "Start and interactive prompt to select and esc-exec into a running container.",
	Run: func(cmd *cobra.Command, args []string) {
		promptCluster()
		promptService()
		promptTask()
		promptContainer()
		promptCommand()

		err := ecs.ExecuteCommand(&execInput)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func promptCluster() {
	if execInput.Cluster == "" {
		clusters, err := ecs.GetClusters()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		clusterPrompt := &survey.Select{
			Message: "Select a cluster:",
			Options: clusters,
		}

		err = survey.AskOne(clusterPrompt, &execInput.Cluster)
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			os.Exit(1)
		}
	}
}

func promptService() {
	if execInput.Service == "" && execInput.Task == "" {

		services, err := ecs.GetServices(execInput.Cluster)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		prompt := &survey.Select{
			Message: "Select a service:",
			Options: services,
		}

		err = survey.AskOne(prompt, &execInput.Service)
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			os.Exit(1)
		}
	}
}

func promptTask() {
	if execInput.Task == "" {
		tasks, err := ecs.GetRunningTasks(execInput.Cluster, execInput.Service)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if len(tasks) == 1 {
			execInput.Task = tasks[0]
		} else {
			taskPrompt := &survey.Select{
				Message: "Select a task",
				Options: tasks,
			}

			err = survey.AskOne(taskPrompt, &execInput.Task)
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				os.Exit(1)
			}
		}
	}
}

func promptContainer() {
	if execInput.Container == "" {
		containers, err := ecs.GetContainers(execInput.Cluster, execInput.Task)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if len(containers) == 1 {
			execInput.Container = containers[0]
		} else {
			prompt := &survey.Select{
				Message: "Select a container",
				Options: containers,
			}

			err = survey.AskOne(prompt, &execInput.Container)
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				os.Exit(1)
			}
		}
	}
}

func promptCommand() {
	if execInput.Command == "" {
		prompt := &survey.Input{
			Message: "Command",
			Default: "bash",
		}
		err := survey.AskOne(prompt, &execInput.Command)

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			os.Exit(1)
		}
	}
}
