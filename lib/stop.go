package ecs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
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
			logWarning("Successfully stopped " + *task.TaskArn)
		}
	}
}
