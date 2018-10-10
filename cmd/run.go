package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/spf13/cobra"
	"gitlab.com/justmiles/ecs-cli/lib"
)

var (
	task ecs.Task
	wg   sync.WaitGroup
)

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.PersistentFlags().StringVarP(&task.Cluster, "cluster", "", "", "ECS cluster")
	runCmd.PersistentFlags().StringVarP(&task.Name, "name", "n", "ecs-cli-app", "Assign a name to the task")
	runCmd.PersistentFlags().StringVar(&task.ExecutionRoleArn, "execution-role", "", "Execution role ARN (required for Fargate)")
	runCmd.PersistentFlags().BoolVarP(&task.Detach, "detach", "d", false, "[TODO] Run the task in the background")
	runCmd.PersistentFlags().Int64VarP(&task.Count, "count", "c", 1, "Spawn n tasks")
	runCmd.PersistentFlags().Int64VarP(&task.Memory, "memory", "m", 0, "Memory limit")
	runCmd.PersistentFlags().Int64Var(&task.CPUReservation, "cpu-reservation", 1024, "CPU reservation")
	runCmd.PersistentFlags().Int64Var(&task.MemoryReservation, "memory-reservation", 2048, "Memory reservation")
	runCmd.PersistentFlags().StringArrayVarP(&task.Environment, "env", "e", nil, "Set environment variables")
	runCmd.PersistentFlags().StringArrayVarP(&task.Publish, "publish", "p", nil, "Publish a container's port(s) to the host")
	runCmd.PersistentFlags().StringArrayVar(&task.SecurityGroups, "security-groups", nil, "[TODO] Attach security groups to task")
	runCmd.PersistentFlags().StringArrayVar(&task.Subnets, "subnet", nil, "Subnet(s) where task should run")
	runCmd.PersistentFlags().StringArrayVarP(&task.Volumes, "volume", "v", nil, "Map volume to ECS Container Instance")
	runCmd.PersistentFlags().BoolVar(&task.Public, "public", false, "[TODO] Assign public IP")
	runCmd.PersistentFlags().BoolVar(&task.Fargate, "fargate", false, "Launch in Fargate")
	runCmd.Flags().SetInterspersed(false)
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

		// fargate validation
		if task.Fargate {
			if len(task.Subnets) == 0 {
				fmt.Println("Fargate requires at least one subnet (--subnet)")
				os.Exit(1)
			}
			if task.ExecutionRoleArn == "" {
				fmt.Println("Fargate requires an executino role (--execution-role)")
				os.Exit(1)
			}
		}
		// Run the task
		err := task.Run()
		defer task.Stop()
		check(err)

		wg.Add(2)
		go task.Stream()
		go task.Check()

		if err != nil {
			fmt.Println(err.Error())
		}
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		go func() {
			for sig := range c {
				fmt.Printf("I got a %T\n", sig)
				task.Stop()
				os.Exit(0)
			}
		}()

		wg.Wait()
	},
}
