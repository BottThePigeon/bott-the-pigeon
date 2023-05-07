package lambdautils

import (
	aws "bott-the-pigeon/lib/aws/session"

	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

// The Lambda client pointer is stored, and can be accessed later.
var lambdasvc *lambda.Lambda

// Returns the stored Lambda client or creates one if it doesn't exist,
// using the provided AWS session.
func getClient(awssess *session.Session) *lambda.Lambda {
	if lambdasvc != nil {
		return lambdasvc
	} else {
		lambdasvc := lambda.New(awssess)
		return lambdasvc
	}
}

// Thin wrapper for the Lambda Invoke function, using a managed
// Lambda service, in us-east-1.
func InvokeLambda(input *lambda.InvokeInput) (*lambda.InvokeOutput, error) {
	awssess, err := aws.GetRegionalSession("us-east-1")
	if err != nil {
		return nil, fmt.Errorf("failed to get AWS session: %v", err)
	}
	lambdasvc := getClient(awssess)
	return lambdasvc.Invoke(input)
}
