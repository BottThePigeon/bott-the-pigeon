package onmessagehandlers

import (
	awssess "bott-the-pigeon/aws-utils/session"

	"log"
	"math/rand"
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bwmarrin/discordgo"
)

// TODO: This will be refactored - most of which will probably be in its own Lambda

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

	img := &discordgo.MessageEmbed{
		Image: &discordgo.MessageEmbedImage{
			URL: *imgUrl,
		},
	}

	return img
}

// Wrapper (providing parameters etc. to make its called functions more agnostic)
// to pull a random object from S3. (In our case, an image).
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
func getRandomS3Key(s3svc *s3.S3, bucketLoc string) string {

	objList, err := s3svc.ListObjects(&s3.ListObjectsInput{
		Bucket: &bucketLoc,
	})

	if err != nil {
		log.Println("Could not retrieve object list from bucket: ", err)
	}

	randomIndex := rand.Intn(len(objList.Contents))
	randomKey := *objList.Contents[randomIndex].Key

	return randomKey
}

// Returns a presigned URL for access to an S3 object based on the bucket and key.
func getS3Object(s3svc *s3.S3, bucketLoc string, objKey string) (*string, error) {
	s3req, _ := s3svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: &bucketLoc,
		Key:    &objKey,
	})

	if s3req.Error != nil {
		return nil, s3req.Error
	}

	url, err := s3req.Presign(time.Hour * 24)

	if err != nil {
		return nil, err
	}

	return &url, nil
}
