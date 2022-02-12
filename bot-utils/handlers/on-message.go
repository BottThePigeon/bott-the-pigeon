package handlers

import (
	botutils "bott-the-pigeon/bot-utils"
	handlers "bott-the-pigeon/bot-utils/handlers/on-message-handlers"

	"strings"

	"github.com/bwmarrin/discordgo"
)

//  Handles Discord MessageCreate event and handles several relevant subconditions
func OnMessage(bot *discordgo.Session, msg *discordgo.MessageCreate) {

	switch {
	case botutils.CheckForBotMention(bot, msg.Mentions) &&
		strings.Contains(msg.Content, "support"):
		handlers.OnHelp(bot, msg)

	case strings.Contains(msg.Content, ">pigeon"):
		handlers.OnPigeon(bot, msg)

	case botutils.CheckForBotMention(bot, msg.Mentions):
		handlers.OnTag(bot, msg)
	}
}
