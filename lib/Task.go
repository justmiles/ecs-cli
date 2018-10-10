package ecs

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/ecs"
)

// Task represents a single, runnable task
type Task struct {
	Cluster           string
	Name              string
	Image             string
	ExecutionRoleArn  string
	Detach            bool
	Public            bool
	Fargate           bool
	Count             int64
	Memory            int64
	MemoryReservation int64
	CPUReservation    int64
	Publish           []string
	Environment       []string
	SecurityGroups    []string
	Subnets           []string
	Volumes           []string
	Command           []string
	TaskDefinition    ecs.TaskDefinition
	Tasks             []*ecs.Task
}

var (
	sess = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
)

// Stop a task
func (t *Task) Stop() {
	var svc = ecs.New(sess)
	logInfo("Stopping tasks")
	for _, task := range t.Tasks {
		_, err := svc.StopTask(&ecs.StopTaskInput{
			Cluster: task.ClusterArn,
			Reason:  aws.String("recieved a ^C"),
			Task:    task.TaskArn,
		})

		if err != nil {
			logError(err)
		} else {
			logInfo("Successfully stopped " + *task.TaskArn)
		}
	}
}

// Run a task
func (t *Task) Run() error {
	var launchType string
	var publicIP string
	var svc = ecs.New(sess)
	t.createLogGroup()
	logInfo("Creating task definition")
	v, m := buildMountPoint(t.Volumes)
	taskDefInput := ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions: []*ecs.ContainerDefinition{
			&ecs.ContainerDefinition{
				Name:    aws.String(t.Name),
				Image:   aws.String(t.Image),
				Command: aws.StringSlice(t.Command),
				LogConfiguration: &ecs.LogConfiguration{
					LogDriver: aws.String("awslogs"),
					Options: aws.StringMap(map[string]string{
						"awslogs-group":         "/" + t.Cluster + "/ecs/ephemeral-task-from-ecs-cli",
						"awslogs-region":        "us-east-1",
						"awslogs-stream-prefix": t.Name,
					}),
				},
				Environment:  buildEnvironmentKeyValuePair(t.Environment),
				PortMappings: buildPortMapping(t.Publish),
				MountPoints:  m,
			},
		},
		Volumes: v,
		Family:  aws.String("ephemeral-task-from-ecs-cli"),
	}

	if t.Memory > 0 {
		taskDefInput.ContainerDefinitions[0].Memory = aws.Int64(t.Memory)
	}

	if t.MemoryReservation > 0 {
		taskDefInput.ContainerDefinitions[0].MemoryReservation = aws.Int64(t.MemoryReservation)
	}

	if t.Fargate {
		taskDefInput.RequiresCompatibilities = aws.StringSlice([]string{"FARGATE"})
		taskDefInput.NetworkMode = aws.String("awsvpc")
		taskDefInput.ExecutionRoleArn = aws.String(t.ExecutionRoleArn)
		taskDefInput.Cpu = aws.String(fmt.Sprintf("%d", t.CPUReservation))
		taskDefInput.Memory = aws.String(fmt.Sprintf("%d", t.MemoryReservation))
	}

	// Register a new task definition
	taskDef, err := svc.RegisterTaskDefinition(&taskDefInput)
	if err != nil {
		return err
	}
	t.TaskDefinition = *taskDef.TaskDefinition

	// fmt.Println(taskDef) //debug
	logInfo("Running task definition: " + *t.TaskDefinition.TaskDefinitionArn)
	// Build the task parametes
	runTaskInput := &ecs.RunTaskInput{
		Cluster:        aws.String(t.Cluster),
		Count:          aws.Int64(t.Count),
		StartedBy:      aws.String("ecs cli"),
		TaskDefinition: taskDef.TaskDefinition.Family,
	}

	// Configure for Fargate
	if t.Fargate {

		if t.Public {
			publicIP = "ENABLED"
		} else {
			publicIP = "DISABLED"
		}
		launchType = "FARGATE"
		runTaskInput.NetworkConfiguration = &ecs.NetworkConfiguration{
			AwsvpcConfiguration: &ecs.AwsVpcConfiguration{
				AssignPublicIp: aws.String(publicIP),
				SecurityGroups: aws.StringSlice(t.SecurityGroups),
				Subnets:        aws.StringSlice(t.Subnets),
			},
		}

		// Default to EC2 launch tye
	} else {
		launchType = "EC2"
	}

	runTaskInput.LaunchType = aws.String(launchType)

	// Run the task
	runTaskResponse, err := svc.RunTask(runTaskInput)
	if err != nil {
		return err
	}
	t.Tasks = runTaskResponse.Tasks
	return nil
}

// Stream logs to stdout
func (t *Task) Stream() {
	logInfo("Streaming from Cloudwatch Logs")
	var svc = cloudwatchlogs.New(sess)
	nextToken := ""
	for _, task := range t.Tasks {
		for {
			logEventsInput := cloudwatchlogs.GetLogEventsInput{
				StartFromHead: aws.Bool(true),
				LogGroupName:  aws.String(*t.TaskDefinition.ContainerDefinitions[0].LogConfiguration.Options["awslogs-group"]),
				LogStreamName: aws.String(t.Name + "/" + t.Name + "/" + strings.Split(*task.TaskArn, "/")[1]),
			}

			if nextToken != "" {
				logEventsInput.NextToken = aws.String(nextToken)
			}

			logEvents, err := svc.GetLogEvents(&logEventsInput)
			if err != nil {
				if awsErr, ok := err.(awserr.Error); ok {
					// Get error details
					if awsErr.Code() == "ResourceNotFoundException" {
						time.Sleep(time.Second / 5)
						continue
					}
				} else {
					logFatalError(err)
				}
			}

			for _, log := range logEvents.Events {
				logCloudWatchEvent(log)
			}

			if *logEvents.NextForwardToken != "" {
				nextToken = *logEvents.NextForwardToken
			}

			time.Sleep(time.Second / 5)
		}
	}
}

// Check the container is still running
func (t *Task) Check() {
	var svc = ecs.New(sess)
	var tasks []string
	var cluster *string
	var stoppedCount int
	var exitCode int64 = 1
	var reportedPorts = false
	var ip *string

	for _, task := range t.Tasks {
		tasks = append(tasks, *task.TaskArn)
		cluster = task.ClusterArn
	}

	for {
		res, err := svc.DescribeTasks(&ecs.DescribeTasksInput{
			Cluster: cluster,
			Tasks:   aws.StringSlice(tasks),
		})
		logError(err)

		for _, ecsTask := range res.Tasks {

			if ip == nil && ecsTask.ContainerInstanceArn != nil {
				res, err := svc.DescribeContainerInstances(&ecs.DescribeContainerInstancesInput{
					Cluster:            &t.Cluster,
					ContainerInstances: aws.StringSlice([]string{*ecsTask.ContainerInstanceArn}),
				})
				logError(err)
				// getEc2Ip
				ip = getEc2InstanceIp(*res.ContainerInstances[0].Ec2InstanceId)
				logInfo(fmt.Sprintf("Container is starting on EC2 instance %v (%v).", *res.ContainerInstances[0].Ec2InstanceId, *ip))
			}

			if !reportedPorts {
				for _, container := range ecsTask.Containers {

					if container.NetworkBindings != nil {
						for _, networkBind := range container.NetworkBindings {
							//  get container instance ip from container.ContainerInstanceArn
							logInfo(fmt.Sprintf("Container is available here\n\thttp://%v:%v\n\tTCP %v %v", *ip, *networkBind.HostPort, *ip, *networkBind.HostPort))
							reportedPorts = true
						}
					}
				}
			}

			if *ecsTask.LastStatus == "STOPPED" {
				for _, container := range ecsTask.Containers {
					if container.ExitCode != nil {
						exitCode = *container.ExitCode
					}

					logInfo(fmt.Sprintf("Task %v has stopped (exit code %v):\n\t%v", *ecsTask.TaskArn, exitCode, *ecsTask.StoppedReason))
					if container.Reason != nil {
						logInfo(fmt.Sprintf("\t%v", *container.Reason))
					}
					stoppedCount++
				}
			}
		}
		if stoppedCount == len(res.Tasks) {
			logInfo("All containers have exited")
			time.Sleep(time.Second * 5) // give the logs another chance to come in
			os.Exit(int(exitCode))
		}
		if t.Detach {
			return
		}
		time.Sleep(time.Second * 5)

	}

}

func (t *Task) createLogGroup() {
	var svc = cloudwatchlogs.New(sess)
	var logGroupName = aws.String("/" + t.Cluster + "/ecs/ephemeral-task-from-ecs-cli")

	output, err := svc.DescribeLogGroups(&cloudwatchlogs.DescribeLogGroupsInput{
		LogGroupNamePrefix: logGroupName,
	})
	if err != nil {
		logError(err)
		return
	}
	if len(output.LogGroups) == 0 {
		logInfo(fmt.Sprintf("Creating Log Group %s\n", *logGroupName))
		svc.CreateLogGroup(&cloudwatchlogs.CreateLogGroupInput{
			LogGroupName: logGroupName,
		})
	}
}
