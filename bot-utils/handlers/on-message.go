package handlers

import (
	botutils "bott-the-pigeon/bot-utils"
	handlers "bott-the-pigeon/bot-utils/handlers/on-message-handlers"

	"strings"

	"github.com/bwmarrin/discordgo"
)

//  Handles Discord MessageCreate event and handles several relevant subconditions
func OnMessage(bot *discordgo.Session, msg *discordgo.MessageCreate) {

	// Respond with random response if the bot is mentioned
	if botutils.CheckForBotMention(bot, msg.Mentions) {
		handlers.OnTag(bot, msg)
	}

	if strings.Contains(msg.Content, "!pigeon") {
		handlers.OnPigeon(bot, msg)
	}
}
