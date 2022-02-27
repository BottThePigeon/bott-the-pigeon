package handlers

import (
	c "bott-the-pigeon/app/common"
	e "bott-the-pigeon/app/errors"
	handlers "bott-the-pigeon/app/handlers/on-message-handlers"

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
	case mention && strings.Contains(msg.Content, "support"):
		handlers.OnHelp(bot, msg)

	case strings.Contains(msg.Content, ">pigeon"):
		handlers.OnImage(bot, msg)

	case strings.Contains(msg.Content, ">todo"):
		handlers.OnTodo(bot, msg)

	case mention:
		handlers.OnTag(bot, msg)
	}
}
