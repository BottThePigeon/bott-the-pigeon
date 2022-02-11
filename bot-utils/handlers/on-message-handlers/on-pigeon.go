package onmessagehandlers

import (
	awssess "bott-the-pigeon/aws-utils/session"

	"log"
	"math/rand"
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bwmarrin/discordgo"
)

// TODO: Refactors:
// Modularise this a bit. It's likely we'll need GET requests to S3 elsewhere.
// Code should be a bit more agnostic, and therefore reusable (i.e., passing in parameters a bit more).
// Embedding < Actually sending the Base64 image (or something that isn't a link).
// Note: If we do continue to embed (not ideal), maybe don't have an expiry time.
// All of this could be solved potentially by moving S3 get requests into its own module, as an API - maybe use Lambda?

// Bot response to "!pigeon" - sending an image.
func OnPigeon(bot *discordgo.Session, msg *discordgo.MessageCreate) {
	bot.ChannelMessageSendEmbed(msg.ChannelID, getBotImageResponse(bot, msg))
}

// Generate a MessageEmbed struct using image URL.
func getBotImageResponse(bot *discordgo.Session, msg *discordgo.MessageCreate) (*discordgo.MessageEmbed) {
	imgUrl, err := getRandomImageFromS3()

	if err != nil {
		log.Println("Could not retrieve image from S3 bucket: ", err)
	}

	img := &discordgo.MessageEmbed {
		Image: &discordgo.MessageEmbedImage {
			URL: *imgUrl,
		},
	}

	return img
}

// Wrapper (providing parameters etc. to make it's called functions more agnostic) to pull a random object from S3.
// (In our case, an image).
func getRandomImageFromS3() (*string, error) {
	s3svc := s3.New(awssess.GetAWSSession())
	bucketLoc := "btp-pigeon-pics"
	s3ObjKey := getRandomS3Key(s3svc, bucketLoc)

	s3img, err := getS3Object(s3svc, bucketLoc, s3ObjKey)

	if err != nil {
		return nil, err
	}

	return s3img, nil
}

// Returns a random item from a specified bucket.
func getRandomS3Key(s3svc *s3.S3, bucketLoc string) (string) {

	objList, err := s3svc.ListObjects( &s3.ListObjectsInput{
		Bucket: &bucketLoc,
	})

	if err != nil {
		log.Println("Could not retrieve object list from bucket: ", err)
	}

	randomIndex := rand.Intn(len(objList.Contents))
	randomKey := *objList.Contents[randomIndex].Key

	return randomKey
}

// Returns a presigned URL for access to an S3 object based upon the bucket and key provided.
func getS3Object(s3svc *s3.S3, bucketLoc string, objKey string) (*string, error) {
	s3req, _ := s3svc.GetObjectRequest( &s3.GetObjectInput {
		Bucket: &bucketLoc,
		Key: &objKey,
	})

	if s3req.Error != nil {
		return nil, s3req.Error
	}

	url, err := s3req.Presign(time.Hour*24)

	if err != nil {
		return nil, err
	}

	return &url, nil
}