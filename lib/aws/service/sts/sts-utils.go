package stsutils

import (
	aws "bott-the-pigeon/lib/aws/session"

	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

// The Lambda client pointer is stored, and can be accessed later.
var stssvc *sts.STS

// Returns the stored Lambda client or creates one if it doesn't exist,
// using the provided AWS session.
func getClient(awssess *session.Session) *sts.STS {
	if stssvc != nil {
		return stssvc
	} else {
		stssvc := sts.New(awssess)
		return stssvc
	}
}

// Thin wrapper for the STS Assume Role function, using a managed
// STS Service.
func GetCrossAccountCredentials(input *sts.AssumeRoleInput) (*sts.AssumeRoleOutput, error) {
	awssess, err := aws.GetSession()
	if err != nil {
		return nil, fmt.Errorf("failed to get AWS session: %v", err)
	}
	stssvc := getClient(awssess)
	return stssvc.AssumeRole(input)
}
