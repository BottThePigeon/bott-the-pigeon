package ecsutils

import (
	awsutils "bott-the-pigeon/lib/aws/session"

	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

// The ECS client pointer is stored, and can be accessed later.
var ecssvc *ecs.ECS

// Returns the stored ECS client or creates one if it doesn't exist,
// using the provided AWS session.
func getClient(awssess *session.Session) *ecs.ECS {
	if ecssvc != nil {
		return ecssvc
	} else {
		lambdasvc := ecs.New(awssess)
		return lambdasvc
	}
}

// Returns an ECS client with the provided AWS Session and Configs.
func getClientWithCfg(awssess *session.Session, cfg *aws.Config) *ecs.ECS {
	if ecssvc != nil {
		return ecssvc
	} else {
		lambdasvc := ecs.New(awssess, cfg)
		return lambdasvc
	}
}

// Thin wrapper for the ECS Describe Services function, using a managed
// ECS service.
func DescribeServices(input *ecs.DescribeServicesInput) (*ecs.DescribeServicesOutput, error) {
	awssess, err := awsutils.GetSession()
	if err != nil {
		return nil, fmt.Errorf("failed to get AWS session: %v", err)
	}
	ecssvc := getClient(awssess)
	return ecssvc.DescribeServices(input)
}

func DescribeServicesWithCfg(input *ecs.DescribeServicesInput, cfg *aws.Config) (*ecs.DescribeServicesOutput, error) {
	awssess, err := awsutils.GetSession()
	if err != nil {
		return nil, fmt.Errorf("failed to get AWS session: %v", err)
	}
	ecssvc := getClientWithCfg(awssess, cfg)
	return ecssvc.DescribeServices(input)
}
