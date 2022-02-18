package boterrors

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func ThrowBotError(bot *discordgo.Session, msg *discordgo.MessageCreate) {

	_, err := bot.ChannelMessageSendEmbed(msg.ChannelID, &discordgo.MessageEmbed{
		Title:       "Uh-oh. Something went wrong.",
		Description: "As you know, I'm a pigeon, so things like this happen. Please don't kill me.",
		Color:       0xDD4400,
	})

	if err != nil {
		log.Println("Bot could not send a message - something went seriously wrong: ", err)
	}
}