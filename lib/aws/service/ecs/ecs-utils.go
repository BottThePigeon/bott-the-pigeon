package ecsutils

import (
	aws "bott-the-pigeon/lib/aws/session"

	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

// The SSM client pointer is stored, and can be accessed later.
var ecssvc *ecs.ECS

// Returns the stored SSM client or creates one if it doesn't exist,
// using the provided AWS session.
func getClient(awssess *session.Session) *ecs.ECS {
	if ecssvc != nil {
		return ecssvc
	} else {
		ecssvc := ecs.New(awssess)
		return ecssvc
	}
}

// Thin wrapper for the ECS RunTask function, using a managed
// ECS service.
func RunTask(input *ecs.RunTaskInput) (*ecs.RunTaskOutput, error) {
	awssess, err := aws.GetSession()
	if err != nil {
		return nil, fmt.Errorf("failed to get AWS session: %v", err)
	}
	ecssvc := getClient(awssess)
	return ecssvc.RunTask(input)
}