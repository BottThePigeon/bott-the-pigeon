package ssmenv

import (
	aws "bott-the-pigeon/lib/aws/session"

	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

// The SSM client pointer is stored, and can be accessed later.
var ssmsvc *ssm.SSM

// Returns the stored SSM client or creates one if it doesn't exist,
// using the provided AWS session.
func getClient(awssess *session.Session) (*ssm.SSM) {
	if ssmsvc != nil {
		return ssmsvc
	} else {
		ssmsvc := ssm.New(awssess)
		return ssmsvc
	}
}

// Gets environment vars from the provided AWS SSM parameter store path (With an valid session).
func Getenv(ssmPath string) (map[string]string, error) {
	awssess, err := aws.GetSession()
	if err != nil {
		return nil, err
	}
	ssmEnv, err := getEnvFromSSM(awssess, ssmPath)
	if err != nil {
		return nil, err
	}
	env := getEnvMap(ssmEnv, ssmPath)
	return env, nil
}

// Retrieves the environment variables from AWS.
func getEnvFromSSM(awssess *session.Session, ssmPath string) (*ssm.GetParametersByPathOutput, error) {
	ssmsvc := getClient(awssess)
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
