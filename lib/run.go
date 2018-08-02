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

// Run a container
func Run(cluster, name, image string, detach, public, fargate bool, count, memory, memoryReservation int64, publish, environment, securityGroups, subnets, volume, command []string) error {
	var launchType string
	var publicIP string
	var svc = ecs.New(sess)

	taskDefInput := ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions: []*ecs.ContainerDefinition{
			&ecs.ContainerDefinition{
				Name:    aws.String(name),
				Image:   aws.String(image),
				Command: aws.StringSlice(command),
				LogConfiguration: &ecs.LogConfiguration{
					LogDriver: aws.String("awslogs"),
					Options: aws.StringMap(map[string]string{
						"awslogs-group":         "/" + cluster + "/ecs/epemeral-task-from-ecs-cli",
						"awslogs-region":        "us-east-1",
						"awslogs-stream-prefix": name,
					}),
				},
				Environment:  buildEnvironmentKeyValuePair(environment),
				PortMappings: buildPortMapping(publish),
			},
		},
		Family: aws.String("epemeral-task-from-ecs-cli"),
	}

	if memory > 0 {
		taskDefInput.ContainerDefinitions[0].Memory = aws.Int64(memory)
	}

	if memoryReservation > 0 {
		taskDefInput.ContainerDefinitions[0].MemoryReservation = aws.Int64(memoryReservation)
	}

	// Register a new task definition
	taskDef, err := svc.RegisterTaskDefinition(&taskDefInput)

	if err != nil {
		return err
	}

	fmt.Println(taskDef)

	// Build the task parametes
	runTaskInput := &ecs.RunTaskInput{
		Cluster:        aws.String(cluster),
		Count:          aws.Int64(count),
		StartedBy:      aws.String("ecs cli"),
		TaskDefinition: taskDef.TaskDefinition.Family,
	}

	// Configure for Fargate
	if fargate {

		if public {
			publicIP = "ENABLED"
		} else {
			publicIP = "DISABLED"
		}
		launchType = "FARGATE"
		runTaskInput.NetworkConfiguration = &ecs.NetworkConfiguration{
			AwsvpcConfiguration: &ecs.AwsVpcConfiguration{
				AssignPublicIp: aws.String(publicIP),
				SecurityGroups: aws.StringSlice(securityGroups),
				Subnets:        aws.StringSlice(subnets),
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
	fmt.Println(runTaskResponse)

	return nil
}

// TODO: create log group 	/ecs/qa/epemeral-task-from-ecs-cli
