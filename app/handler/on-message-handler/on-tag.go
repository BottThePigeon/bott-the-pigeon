package onmessagehandlers

import (
	e "bott-the-pigeon/app/error"
	"log"

	"github.com/bwmarrin/discordgo"
)

// Sends a simple message from the provided bot.
func OnTag(bot *discordgo.Session, msg *discordgo.MessageCreate) {
	_, err := bot.ChannelMessageSend(msg.ChannelID, "Watashi wa『DUMB BIRD』desu.")
	if err != nil {
		log.Println("failed to send simple message: ", err)
		e.ThrowBotError(bot, msg.ChannelID, err)
		return
	}
}
