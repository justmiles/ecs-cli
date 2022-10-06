package cmd

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"

	ecs "github.com/justmiles/ecs-cli/lib"
	"github.com/spf13/cobra"
)

var (
	task         ecs.Task
	wg           sync.WaitGroup
	validMemCPU  map[int][]int
	noDeregister bool
)

func init() {
	log.SetFlags(0)

	rootCmd.AddCommand(runCmd)
	runCmd.PersistentFlags().StringVarP(&task.Cluster, "cluster", "", "", "ECS cluster")
	runCmd.PersistentFlags().StringVarP(&task.Name, "name", "n", "ephemeral-task-from-ecs-cli", "Assign a name to the task")
	runCmd.PersistentFlags().StringVar(&task.Family, "family", "", "Family for ECS task")
	runCmd.PersistentFlags().StringVar(&task.ExecutionRoleArn, "execution-role", "", "Execution role ARN (required for Fargate)")
	runCmd.PersistentFlags().StringVar(&task.TaskRoleArn, "role", "", "Task role ARN")
	runCmd.PersistentFlags().StringVar(&task.CLIRoleArn, "cli-role", "", "An IAM role ARN to assume before creating/executing a task")
	runCmd.PersistentFlags().BoolVarP(&task.Detach, "detach", "d", false, "Run the task in the background")
	runCmd.PersistentFlags().Int64VarP(&task.Count, "count", "c", 1, "Spawn n tasks")
	runCmd.PersistentFlags().Int64VarP(&task.Memory, "memory", "m", 0, "Memory limit")
	runCmd.PersistentFlags().Int64Var(&task.CPUReservation, "cpu-reservation", 256, "CPU reservation")
	runCmd.PersistentFlags().Int64Var(&task.MemoryReservation, "memory-reservation", 2048, "Memory reservation")
	runCmd.PersistentFlags().StringArrayVarP(&task.Environment, "env", "e", nil, "Set environment variables")
	runCmd.PersistentFlags().StringArrayVarP(&task.Publish, "publish", "p", nil, "Publish a container's port(s) to the host")
	// TODO: attach a specific security group
	runCmd.PersistentFlags().StringArrayVar(&task.SecurityGroups, "security-groups", nil, "attach security groups to task")
	runCmd.PersistentFlags().StringArrayVar(&task.SubnetFilters, "subnet-filter", nil, "'Key=Value' filters for your subnet, eg tag:Name=private")
	runCmd.PersistentFlags().StringArrayVarP(&task.Volumes, "volume", "v", nil, "Map volume to ECS Container Instance")
	runCmd.PersistentFlags().StringArrayVarP(&task.EfsVolumes, "efs-volume", "", nil, "Map EFS volume to ECS Container Instance (ex. fs-23kj2f:/efs/dir:/container/mnt/dir)")
	// TODO: support assigning public ip address
	runCmd.PersistentFlags().BoolVar(&task.Public, "public", false, "assign public ip")
	runCmd.PersistentFlags().BoolVar(&task.Fargate, "fargate", false, "Launch in Fargate")
	runCmd.PersistentFlags().BoolVar(&noDeregister, "no-deregister", false, "do not deregister the task definition")
	runCmd.PersistentFlags().BoolVar(&task.Debug, "debug", false, "Verbose logging")
	runCmd.Flags().SetInterspersed(false)

	// Init CPU/Memory configs
	validMemCPU = make(map[int][]int)
	validMemCPU[256] = []int{512, 1024, 2048}
	validMemCPU[512] = []int{1024, 2048, 3072, 4096}
	validMemCPU[1024] = []int{2048, 3072, 4096, 5120, 6144, 7168, 8192}
	validMemCPU[2048] = []int{4096, 5120, 6144, 7168, 8192, 9216, 10240, 11264, 12288, 13312, 14336, 15360, 16384}
	validMemCPU[4096] = []int{8192, 9216, 10240, 11264, 12288, 13312, 14336, 15360, 16384, 17408, 18432, 19456, 20480, 21504, 22528, 23552, 24575, 25600, 26624, 27648, 28678, 29696, 30720}
	validMemCPU[8192] = []int{16384, 20480, 24576, 28672, 32768, 36864, 40960, 45056, 49152, 53248, 57344, 61440}
	validMemCPU[16384] = []int{32768, 40960, 49152, 57344, 65536, 73728, 81920, 90112, 98304, 106496, 114688, 122880}
}

// process the list command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a command in a new task",
	Run: func(cmd *cobra.Command, args []string) {

		task.Deregister = !noDeregister
		if len(args) < 1 {
			log.Fatal("Please pass an image to run")
		}

		task.Image = args[0]

		if len(args) > 1 {
			task.Command = args[1:len(args)]
		}

		// fargate validation
		if task.Fargate {
			if len(task.SubnetFilters) == 0 {
				log.Fatal("Fargate requires at least one subnet")
			}
			if !isValidMemCPU(int(task.CPUReservation), int(task.MemoryReservation)) {
				log.Fatal("CPU/Memory unsupported. See supported values here: https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task-cpu-memory-error.html")
			}
		}

		// efs-volume validation
		for _, volume := range task.EfsVolumes {
			av := strings.Split(volume, ":")
			if len(av) != 3 {
				log.Fatal("Incorrect usage (--efs-volume)")
			}
		}

		// Run the task
		err := task.Run()
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

func isValidMemCPU(cpu, memory int) bool {
	for _, allowedCpuValue := range validMemCPU[cpu] {
		if allowedCpuValue == memory {
			return true
		}
	}

	return false
}
