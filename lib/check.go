package ecs

import (
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
)

// Check the container is still running
func (t *Task) Check() {
	var svc = ecs.New(sess)
	var tasks []string
	var cluster *string
	var stoppedCount int

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
			if *ecsTask.LastStatus == "STOPPED" {
				logInfo("Task " + *ecsTask.TaskArn + " has stopped: " + *ecsTask.StoppedReason)
				stoppedCount++
			}
			// fmt.Println(*ecsTask)
		}
		if stoppedCount == len(res.Tasks) {
			logInfo("All containers have exited")
			os.Exit(0)
		}

		time.Sleep(time.Second * 5)

	}

}
