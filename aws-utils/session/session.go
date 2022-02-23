package session

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// The AWS SDK Session pointer is stored, and can be accessed later.
var (
	awssess *session.Session
)

// Return the stored AWS Session or create one if it doesn't exist.
func GetAWSSession() (*session.Session, error) {
	if awssess != nil {
		return awssess, nil
	} else {
		sess, err := session.NewSessionWithOptions(session.Options{
			Config:            aws.Config{Region: aws.String("eu-west-2")},
			SharedConfigState: session.SharedConfigEnable,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to initialise session with AWS: %v", err)
		}
		awssess = sess
		log.Println("New AWS session created.")
		return awssess, nil
	}
}