package onmessagehandlers

import (
	awssess "bott-the-pigeon/aws-utils/session"
	e "bott-the-pigeon/bot-utils/errors"
	"io"

	"log"
	"math/rand"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bwmarrin/discordgo"
)

// Bot response to "!pigeon" - sending an image.
func OnPigeon(bot *discordgo.Session, msg *discordgo.MessageCreate) {
	img, err := getRandomImageFromS3()

	if err != nil {
		log.Println("Could not retrieve image from S3 bucket: ", err)
		e.ThrowBotError(bot, msg)
	}

	bot.ChannelFileSend(msg.ChannelID, "pigeon.png", img)
}

// Wrapper (providing parameters etc. to make its called functions more agnostic)
// to pull a random object from S3. (In our case, an image).
func getRandomImageFromS3() (io.Reader, error) {
	s3svc := s3.New(awssess.GetAWSSession())
	bucketLoc := "btp-pigeon-pics"
	s3ObjKey, err := getRandomS3Key(s3svc, bucketLoc)

	if err != nil {
		return nil, err
	}

	s3ObjReader, err := getS3ObjectReader(s3svc, bucketLoc, *s3ObjKey)

	if err != nil {
		return nil, err
	}

	return s3ObjReader, nil
}

// Returns a random item from a specified bucket.
func getRandomS3Key(s3svc *s3.S3, bucketLoc string) (*string, error) {

	objList, err := s3svc.ListObjects(&s3.ListObjectsInput{
		Bucket: &bucketLoc,
	})

	if err != nil {
		return nil, err
	}

	randomIndex := rand.Intn(len(objList.Contents))
	randomKey := *objList.Contents[randomIndex].Key

	return &randomKey, nil
}

func getS3ObjectReader(s3svc *s3.S3, bucketLoc string, objKey string) (io.ReadCloser, error) {

	s3in := &s3.GetObjectInput {
		Bucket: &bucketLoc,
		Key:	&objKey,
	}

	s3out, err := s3svc.GetObject(s3in)

	if err != nil {
		return nil, err
	}

	return s3out.Body, nil
}
