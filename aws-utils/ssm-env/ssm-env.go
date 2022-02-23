package ssmenv

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

// Gets environment vars from the provided AWS SSM parameter store path (With an valid session).
func GetEnv(awssess *session.Session, ssmPath string) (map[string]string, error) {
	ssmEnv, err := getEnvFromSSM(awssess, ssmPath)
	if err != nil {
		return nil, err
	}
	env := getEnvMap(ssmEnv, ssmPath)
	return env, nil
}

// Retrieves the environment variables from AWS.
func getEnvFromSSM(awssess *session.Session, ssmPath string) (*ssm.GetParametersByPathOutput, error) {
	ssmsvc := ssm.New(awssess, aws.NewConfig())
	withDecryption := true
	ssmparams, err := ssmsvc.GetParametersByPath(&ssm.GetParametersByPathInput{
		Path:           &ssmPath,
		WithDecryption: &withDecryption,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to obtain application credentials from AWS. error: %v", err)
	}
	return ssmparams, nil
}

// Returns a k,v map of the provided SSM Parameters.
func getEnvMap(ssmparams *ssm.GetParametersByPathOutput, ssmPath string) map[string]string {
	env := make(map[string]string)
	for _, p := range ssmparams.Parameters {
		k := strings.ReplaceAll(*p.Name, ssmPath, "")
		env[k] = *p.Value
	}
	return env
}
