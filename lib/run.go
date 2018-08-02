package ecs

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

var (
	sess = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
)

// Run a task
func (t *Task) Run() error {
	var launchType string
	var publicIP string
	var svc = ecs.New(sess)

	taskDefInput := ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions: []*ecs.ContainerDefinition{
			&ecs.ContainerDefinition{
				Name:    aws.String(t.Name),
				Image:   aws.String(t.Image),
				Command: aws.StringSlice(t.Command),
				LogConfiguration: &ecs.LogConfiguration{
					LogDriver: aws.String("awslogs"),
					Options: aws.StringMap(map[string]string{
						"awslogs-group":         "/" + t.Cluster + "/ecs/epemeral-task-from-ecs-cli",
						"awslogs-region":        "us-east-1",
						"awslogs-stream-prefix": t.Name,
					}),
				},
				Environment:  buildEnvironmentKeyValuePair(t.Environment),
				PortMappings: buildPortMapping(t.Publish),
			},
		},
		Family: aws.String("epemeral-task-from-ecs-cli"),
	}

	if t.Memory > 0 {
		taskDefInput.ContainerDefinitions[0].Memory = aws.Int64(t.Memory)
	}

	if t.MemoryReservation > 0 {
		taskDefInput.ContainerDefinitions[0].MemoryReservation = aws.Int64(t.MemoryReservation)
	}

	// Register a new task definition
	taskDef, err := svc.RegisterTaskDefinition(&taskDefInput)

	if err != nil {
		return err
	}

	fmt.Println(taskDef) //debug

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

	fmt.Println(runTaskInput)
	// Run the task
	runTaskResponse, err := svc.RunTask(runTaskInput)
	if err != nil {
		return err
	}
	fmt.Println(runTaskResponse)

	return nil
}

// TODO: create log group 	/ecs/qa/epemeral-task-from-ecs-cli
