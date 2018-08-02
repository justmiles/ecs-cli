package ecs

import (
	"testing"
	// "github.com/aws/aws-sdk-go/service/ecs/ecsiface"
)

func TestRun(t *testing.T) {
	err := Run("qa", "test", "hello-world", false, false, false, 1, 100, 100, nil, nil, nil, []string{"subnet-1fa54557"}, nil, nil)
	if err != nil {
		t.Errorf("got: %v", err.Error())
	}
}

// type mockedReceiveMsgs struct {
// 	sqsiface.SQSAPI
// 	Resp sqs.ReceiveMessageOutput
// }
