package onmessagehandlers

import (
	e "bott-the-pigeon/app/errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// Sends a bot usage help message from the provided bot.
func OnHelp(bot *discordgo.Session, msg *discordgo.MessageCreate) error {
	_, err := bot.ChannelMessageSend(msg.ChannelID,
		"Hello! My command prefix is a `>`. Get it? Because it looks like a beak.\n\n"+
			"**Commands**:\n"+
			"`pigeon`: Send a random picture of a pigeon.\n"+
			"`todo {Some feature}`: Submit a suggestion to the project todo list.\n\n"+
			"_That's all for now, because I'm a dumb bird._")
	if err != nil {
		err = fmt.Errorf("failed to send help message: %v", err)
		e.ThrowBotError(bot, msg.ChannelID, err)
		return err
	}
	return nil
}
// TODO: This is very hard-coded, and should be called from an API in some way.