package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.com/justmiles/ecs-cli/lib"
)

var (
	task ecs.Task
)

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.PersistentFlags().StringVarP(&task.Cluster, "cluster", "", "", "ECS cluster")
	runCmd.PersistentFlags().StringVarP(&task.Name, "name", "n", "ecs-cli-app", "Assign a name to the task")
	runCmd.PersistentFlags().BoolVarP(&task.Detach, "detach", "d", false, "[TODO] Run the task in the background")
	runCmd.PersistentFlags().Int64VarP(&task.Count, "count", "c", 1, "Spawn n tasks")
	runCmd.PersistentFlags().Int64VarP(&task.Memory, "memory", "m", 0, "Memory limit")
	runCmd.PersistentFlags().Int64Var(&task.MemoryReservation, "memory-reservation", 1024, "Memory reservation")
	runCmd.PersistentFlags().StringArrayVarP(&task.Environment, "env", "e", nil, "Set environment variables")
	runCmd.PersistentFlags().StringArrayVarP(&task.Publish, "publish", "p", nil, "Publish a container's port(s) to the host")
	runCmd.PersistentFlags().StringArrayVar(&task.SecurityGroups, "security-groups", nil, "[TODO] Attach security groups to task")
	runCmd.PersistentFlags().StringArrayVar(&task.Subnets, "subnet", nil, "[TODO] Subnet(s) where task should run")
	// mark subnets required if using fargate
	runCmd.PersistentFlags().StringArrayVarP(&task.Volumes, "volume", "v", nil, "[TODO] Map volume to ECS Container Instance")
	runCmd.PersistentFlags().BoolVar(&task.Public, "public", false, "[TODO] Assign public IP")
	runCmd.PersistentFlags().BoolVar(&task.Fargate, "fargate", false, "[TODO] Launch in Fargate")
}

// process the list command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a command in a new task",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Please pass an image to run")
			return
		}

		task.Image = args[0]

		if len(args) > 1 {
			task.Command = args[1:len(args)]
		}
		// Run the task
		err := task.Run()
		if err != nil {
			fmt.Println(err.Error())
		}
	},
}
