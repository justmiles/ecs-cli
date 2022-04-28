package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	// "os"
	// "os/signal"
	// "strings"
	// "sync"
	// ecs "github.com/justmiles/ecs-cli/lib"
	"github.com/spf13/cobra"
)

func init() {
	log.SetFlags(0)
	rootCmd.AddCommand(runTaskDefCmd)
	runTaskDefCmd.PersistentFlags().StringVarP(&task.Cluster, "cluster", "", "default", "ECS cluster")
	runTaskDefCmd.PersistentFlags().StringVarP(&task.Family, "family", "", "", "The family name of the task definition. The latest ACTIVE revision is used.")
	runTaskDefCmd.PersistentFlags().StringVarP(&task.ImageVersion, "image-version", "", "", "Optionally pass in an image-version to override the current image version.")
	runTaskDefCmd.PersistentFlags().StringVarP(&task.Name, "name", "n", "ephemeral-task-from-ecs-cli", "Assign a name to the task")
	// TODO: attach a specific security group
	runTaskDefCmd.PersistentFlags().StringArrayVar(&task.SecurityGroups, "security-groups", nil, "attach security groups to task")
	runTaskDefCmd.PersistentFlags().StringArrayVar(&task.SubnetFilters, "subnet-filter", nil, "'Key=Value' filters for your subnet, eg tag:Name=private")
	// TODO: support assigning public ip address
	runTaskDefCmd.PersistentFlags().BoolVar(&task.Public, "public", false, "assign public ip")
	runTaskDefCmd.PersistentFlags().BoolVar(&task.Wait, "wait", false, "wait for container to finish")
	runTaskDefCmd.PersistentFlags().BoolVarP(&task.Detach, "detach", "d", false, "Run the task in the background")
	runTaskDefCmd.PersistentFlags().BoolVar(&task.Deregister, "deregister", false, "deregister the task definition after completion")

	runTaskDefCmd.PersistentFlags().Int64VarP(&task.Count, "count", "c", 1, "Spawn n tasks")
	runTaskDefCmd.Flags().SetInterspersed(false)
}

// process the list command
var runTaskDefCmd = &cobra.Command{
	Use:   "run-task-def",
	Short: "Run a task from an existing task definition",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello world")
		fmt.Println(task)

		task.Fargate = true

		if len(task.SubnetFilters) == 0 {
			log.Fatal("Fargate requires at least one subnet")
		}

		// Run the task
		err := task.RunTaskDef()
		check(err)
		if task.Detach {
			task.Check()
		} else {
			defer task.Stop()
			wg.Add(2)
			go task.Stream()
			go task.Check()
			if err != nil {
				log.Fatal(err.Error())
			}
			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt)
			go func() {
				for sig := range c {
					log.Printf("I got a %T\n", sig)
					task.Stop()
					os.Exit(0)
				}
			}()
			wg.Wait()
		}
	},
}
