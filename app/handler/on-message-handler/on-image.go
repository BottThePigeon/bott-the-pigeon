package onmessagehandlers

import (
	e "bott-the-pigeon/app/error"
	s3utils "bott-the-pigeon/lib/aws/service/s3"
	awssess "bott-the-pigeon/lib/aws/session"
	"fmt"

	"math/rand"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bwmarrin/discordgo"
)

// Sends a random image of a pigeon from the provided bot.
func OnImage(bot *discordgo.Session, msg *discordgo.MessageCreate) {
	s3Obj, err := getRandomObjectFromS3("btp-pigeon-pics")
	if err != nil {
		e.ThrowBotError(bot, msg.ChannelID, err)
		return
	}
	img := s3Obj.Body
	defer img.Close()
	_, err = bot.ChannelFileSend(msg.ChannelID, s3Obj.Key, img)
	if err != nil {
		err = fmt.Errorf("failed to send message with file attachment: %v", err)
		e.ThrowBotError(bot, msg.ChannelID, err)
		return
	}
}

// Returns a random object from a specified S3 bucket.
func getRandomObjectFromS3(bucket string) (*s3utils.S3_ObjectWithKey, error) {
	awssess, err := awssess.GetSession()
	if err != nil {
		return nil, fmt.Errorf("failed to get AWS session: %v", err)
	}
	s3svc := s3.New(awssess)
	s3ObjKey, err := getRandomS3Key(s3svc, bucket)
	if err != nil {
		return nil, err
	}
	s3ObjReader, err := s3utils.GetS3ObjectIOStream(bucket, *s3ObjKey)
	if err != nil {
		return nil, err
	}
	s3Obj := &s3utils.S3_ObjectWithKey{
		Key: *s3ObjKey,
		Body: s3ObjReader,
	}
	return s3Obj, nil
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