package ecs

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
)

func TestRun(t *testing.T) {
	// err := Run("qa", "test", "hello-world", false, false, false, 1, 0, 1024, nil, nil, nil, nil, nil, nil)
	// if err != nil {
	// 	t.Errorf("got: %v", err.Error())
	// }
	x := []string{
		"x=test",
		"y=thisisthevalueforY",
		"SOMEENVVAR",
	}

	// task := Task{
	// 	Name: "test",
	// }
	// task.Stream()
	buildEnvironmentKeyValuePair(x)
	return
}

type mockedReceiveMsgs struct {
	ecsiface.ECSAPI
	Resp ecs.RunTaskOutput
}

func (m mockedReceiveMsgs) ReceiveMessage(in *ecs.RunTaskInput) (*ecs.RunTaskOutput, error) {
	// Only need to return mocked response output
	res := `{
    Failures: [],
    Tasks: [{
        Attachments: [],
        ClusterArn: "arn:aws:ecs:us-east-1:965579072529:cluster/qa",
        ContainerInstanceArn: "arn:aws:ecs:us-east-1:965579072529:container-instance/90da1657-b5bd-461e-84f7-62775f670a09",
        Containers: [{
            ContainerArn: "arn:aws:ecs:us-east-1:965579072529:container/01f8b585-7688-48f4-8a68-31e7b0c594a5",
            LastStatus: "PENDING",
            Name: "test",
            NetworkInterfaces: [],
            TaskArn: "arn:aws:ecs:us-east-1:965579072529:task/8ce3fc54-1630-4703-9bfd-3b4d55dc6219"
          }],
        Cpu: "0",
        CreatedAt: 2018-08-02 14:39:10 +0000 UTC,
        DesiredStatus: "RUNNING",
        Group: "family:epemeral-task-from-ecs-cli",
        LastStatus: "PENDING",
        LaunchType: "EC2",
        Memory: "100",
        Overrides: {
          ContainerOverrides: [{
              Name: "test"
            }]
        },
        StartedBy: "ecs cli",
        TaskArn: "arn:aws:ecs:us-east-1:965579072529:task/8ce3fc54-1630-4703-9bfd-3b4d55dc6219",
        TaskDefinitionArn: "arn:aws:ecs:us-east-1:965579072529:task-definition/epemeral-task-from-ecs-cli:22",
        Version: 1
      }]
  }`

	json.Unmarshal([]byte(res), &m.Resp)

	return &m.Resp, nil
}
