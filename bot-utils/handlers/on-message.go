package handlers

import (
	c "bott-the-pigeon/bot-utils/common"
	e "bott-the-pigeon/bot-utils/errors"
	handlers "bott-the-pigeon/bot-utils/handlers/on-message-handlers"

	"strings"

	"github.com/bwmarrin/discordgo"
)

// TODO: Create some sort of handler architecture/flow diagram. The basic idea should be:
// - The top level of the parent handler ("E.g., OnMessage") should be a VERY simple switch/case,
//   calling the specific handler.
// - Only the top level of a specific handler (i.e., those called by switch/case) can make a system call.
// - Only the top level of a handler should be able to use the bot (i.e., for errors, responses, etc.)
// - Therefore, only the top level of a handler should effectively be "impure".

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
