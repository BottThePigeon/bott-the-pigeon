package init

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// Initialise an AWS Session for use throughout the application
func InitAws() (*session.Session) {

	awssess, err := session.NewSessionWithOptions(session.Options {
		Config:				aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))},
		SharedConfigState: 	session.SharedConfigEnable,
	})
	if err != nil {
		log.Fatal("Could not initialise session with AWS: ", err)
	}

	return awssess;
}