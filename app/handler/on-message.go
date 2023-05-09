package handlers

import (
	c "bott-the-pigeon/app/common"
	e "bott-the-pigeon/app/error"
	handlers "bott-the-pigeon/app/handler/on-message-handler"

	"strings"

	"github.com/bwmarrin/discordgo"
)

// Handles Discord MessageCreate event and handles several relevant subconditions.
func OnMessage(bot *discordgo.Session, msg *discordgo.MessageCreate) {

	mention, err := c.CheckForBotMention(bot, msg.Mentions)
	if err != nil {
		e.ThrowBotError(bot, msg.ChannelID, err)
	}

	switch {
	case strings.ToLower(msg.Content) == ">support":
		handlers.OnHelp(bot, msg)

	case strings.ToLower(msg.Content) == ">pigeon":
		handlers.OnImage(bot, msg)

	// HasPrefix because we want to take whatever comes after as the input
	case strings.HasPrefix(strings.ToLower(msg.Content), ">todo"):
		handlers.OnTodo(bot, msg)

	case strings.ToLower(msg.Content) == ">mc-status":
		handlers.OnMinecraftStatus(bot, msg)

	case strings.ToLower(msg.Content) == ">mc":
		handlers.OnMinecraft(bot, msg)

	case mention:
		handlers.OnTag(bot, msg)
	}
}
