package session

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// The AWS SDK Session pointer is stored, and can be accessed later.
var (
	awssess *session.Session
)

// Return the stored AWS Session or create one if not.
// Therefore, initialisation code should only run once.
func GetAWSSession() *session.Session {

	if awssess != nil {
		return awssess
	} else {
		sess, err := session.NewSessionWithOptions(session.Options{
			Config:            aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))},
			SharedConfigState: session.SharedConfigEnable,
		})
		if err != nil {
			log.Fatal("Could not initialise session with AWS: ", err)
		}

		awssess = sess
		log.Println("New AWS session created.")
		return awssess
	}
}
