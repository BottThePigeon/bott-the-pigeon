package s3utils

import (
	aws "bott-the-pigeon/lib/aws/session"

	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// An S3 object with the associated key.
type S3_ObjectWithKey struct {
	Key string
	Body io.ReadCloser
}

// The CloudWatch client pointer is stored, and can be accessed later.
var s3svc *s3.S3

// Returns the stored CloudWatch client or creates one if it doesn't exist,
// using the provided AWS session.
func getClient(awssess *session.Session) *s3.S3 {
	if s3svc != nil {
		return s3svc
	} else {
		s3svc := s3.New(awssess)
		return s3svc
	}
}

// Returns an object by provided bucket/key, as an io.ReadCloser.
func GetS3ObjectIOStream(bucketLoc string, objKey string) (io.ReadCloser, error) {
	awssess, err := aws.GetSession()
	if err != nil {
		return nil, fmt.Errorf("failed to get AWS session: %v", err)
	}
	s3svc := getClient(awssess)
	s3in := &s3.GetObjectInput {
		Bucket: &bucketLoc,
		Key:	&objKey,
	}
	s3out, err := s3svc.GetObject(s3in)
	if err != nil {
		return nil, fmt.Errorf("failed to get object %v from bucket %v: %v", objKey, bucketLoc, err)
	}
	return s3out.Body, nil
}