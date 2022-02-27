package boterrors

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

// TODO: Add CloudWatch logging and generate a Correlation GUID shared between it and the Discord message

// Sends the provided error message using the provided bot.
func ThrowBotError(bot *discordgo.Session, channel string, e error) {
	log.Println(e.Error())
	_, err := bot.ChannelMessageSendEmbed(channel, &discordgo.MessageEmbed{
		Title:       "Uh-oh. Something went wrong.",
		Description: "As you know, I'm a pigeon, so things like this happen. Please don't kill me.",
		Color:       0xDD4400,
		Footer: &discordgo.MessageEmbedFooter{
			Text: e.Error(),
		},
	})

	//If we're throwing this error, welp, everything's gone bloody topsy-turvy.
	if err != nil {
		log.Println("failed to send error message - something went seriously wrong: ", err)
	}
}