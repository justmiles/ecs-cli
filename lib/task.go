package ecs

import "github.com/aws/aws-sdk-go/service/ecs"

// Task represents a single, runnable task
type Task struct {
	Cluster           string
	Name              string
	Image             string
	Detach            bool
	Public            bool
	Fargate           bool
	Count             int64
	Memory            int64
	MemoryReservation int64
	Publish           []string
	Environment       []string
	SecurityGroups    []string
	Subnets           []string
	Volumes           []string
	Command           []string
	TaskDefinition    ecs.TaskDefinition
	Tasks             []*ecs.Task
}
