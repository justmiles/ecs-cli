package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.com/justmiles/ecs-cli/lib"
)

var (
	name              string
	detach            bool
	public            bool
	fargate           bool
	count             int64
	memory            int64
	memoryReservation int64
	publish           []string
	environment       []string
	securityGroups    []string
	subnets           []string
	volume            []string
)

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.PersistentFlags().StringVarP(&name, "name", "n", "ecs-cli-app", "[TODO] Assign a name to the task")
	runCmd.PersistentFlags().BoolVarP(&detach, "detach", "d", false, "[TODO] Run the task in the background")
	runCmd.PersistentFlags().Int64VarP(&count, "count", "c", 1, "[TODO] Spawn n tasks")
	runCmd.PersistentFlags().Int64VarP(&memory, "memory", "m", 0, "[TODO] Memory limit")
	runCmd.PersistentFlags().Int64Var(&memoryReservation, "memory-reservation", 0, "[TODO] Memory reservation")
	runCmd.PersistentFlags().StringArrayVarP(&environment, "environment", "e", nil, "[TODO] Set environment variables")
	runCmd.PersistentFlags().StringArrayVarP(&publish, "publish", "p", nil, "[TODO] Publish a container's port(s) to the host")
	runCmd.PersistentFlags().StringArrayVar(&securityGroups, "security-groups", nil, "[TODO] Attach security groups to task")
	runCmd.PersistentFlags().StringArrayVar(&subnets, "subnets", nil, "[TODO] Subnet(s) where task should run")
	// mark subnets required
	runCmd.PersistentFlags().StringArrayVarP(&volume, "volume", "v", nil, "[TODO] Map volume to ECS Container Instance")
	runCmd.PersistentFlags().BoolVar(&public, "public", false, "[TODO] Assign public IP")
	runCmd.PersistentFlags().BoolVar(&fargate, "fargate", false, "[TODO] Launch in Fargate")
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
		var command []string
		if len(args) > 1 {
			command = args[1:len(args)]
		}
		err := ecs.Run(cluster, name, args[0], detach, public, fargate, count, memory, memoryReservation, publish, environment, securityGroups, subnets, volume, command)
		if err != nil {
			fmt.Println(err.Error())
		}
	},
}
