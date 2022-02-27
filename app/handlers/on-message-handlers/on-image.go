package onmessagehandlers

import (
	e "bott-the-pigeon/app/errors"
	awssess "bott-the-pigeon/lib/aws/session"
	"fmt"
	"io"

	"math/rand"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bwmarrin/discordgo"
)

// Sends a random image of a pigeon from the provided bot.
func OnImage(bot *discordgo.Session, msg *discordgo.MessageCreate) error {
	img, err := getRandomObjectFromS3("btp-pigeon-pics")
	if err != nil {
		e.ThrowBotError(bot, msg.ChannelID, err)
		return err
	}
	defer img.Close()
	_, err = bot.ChannelFileSend(msg.ChannelID, "pigeon.png", img)
	if err != nil {
		err = fmt.Errorf("failed to send message with file attachment: %v", err)
		e.ThrowBotError(bot, msg.ChannelID, err)
		return err
	}
	return nil
}

// Returns a random object from a specified S3 bucket.
func getRandomObjectFromS3(bucket string) (io.ReadCloser, error) {
	awssess, err := awssess.GetAWSSession()
	if err != nil {
		return nil, fmt.Errorf("failed to get AWS session: %v", err)
	}
	s3svc := s3.New(awssess)
	s3ObjKey, err := getRandomS3Key(s3svc, bucket)
	if err != nil {
		return nil, err
	}
	s3ObjReader, err := getS3ObjectIOStream(s3svc, bucket, *s3ObjKey)
	if err != nil {
		return nil, err
	}
	return s3ObjReader, nil
}

// Returns the key of a random object from a specified S3 bucket.
func getRandomS3Key(s3svc *s3.S3, bucketLoc string) (*string, error) {
	objList, err := s3svc.ListObjects(&s3.ListObjectsInput{
		Bucket: &bucketLoc,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list objects in bucket %v: %v", bucketLoc, err)
	}
	randomIndex := rand.Intn(len(objList.Contents))
	randomKey := *objList.Contents[randomIndex].Key
	return &randomKey, nil
}

// Returns an object by provided bucket/key, as an io.ReadCloser.
func getS3ObjectIOStream(s3svc *s3.S3, bucketLoc string, objKey string) (io.ReadCloser, error) {
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