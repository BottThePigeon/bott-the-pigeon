package session

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// The AWS SDK Session pointer is stored globally.
var (
	AWSSess *session.Session
)
// TODO: Consider refactoring this. It's not the worst thing to have global scope,
// and it seems much more practical than passing the session through layers of functions,
// but it's still not ideal.

// Return the stored AWS Session or create one if not. Therefore, initialisation code should only run once.
func GetAWSSession() (*session.Session) {

	// We return the stored AWSSess if it exists, so we're not creating multiple sessions.
	if AWSSess != nil {
		return AWSSess
	} else {
		AWSSess, err := session.NewSessionWithOptions(session.Options {
			Config:				aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))},
			SharedConfigState: 	session.SharedConfigEnable,
		})
		if err != nil {
			log.Fatal("Could not initialise session with AWS: ", err)
		}
	return AWSSess;
	}
}