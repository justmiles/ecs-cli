package ecs

import (
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

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
						time.Sleep(time.Second * 10)
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
