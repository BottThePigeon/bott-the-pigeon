package errors

import (
	logger "bott-the-pigeon/lib/aws/service/cw-logger"
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

// Sends the provided error message using the provided bot.
func ThrowBotError(bot *discordgo.Session, channel string, e error) {
	log.Println(e)
	uuid, err := logger.Log(os.Getenv("AWS_CW_ERROR_LOG_GROUP"), e.Error())
	var footer string
	if err != nil {
		footer = "(I failed to store the error log too! God I'm a pencil.)"
		log.Println(err)
	} else {
		footer = fmt.Sprintf("BEEP BOOP: %v", *uuid)
	}
	_, err = bot.ChannelMessageSendEmbed(channel, &discordgo.MessageEmbed{
		Title:       "Uh-oh. Something went wrong.",
		Description: "As you know, I'm a pigeon, so things like this happen. Please don't kill me.",
		Color:       0xDD4400,
		Footer: &discordgo.MessageEmbedFooter{
			Text: footer,
		},
	})
	if err != nil {
		log.Println("failed to send error message: ", err)
		return
	}
}
