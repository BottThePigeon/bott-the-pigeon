package onmessagehandlers

import (
	e "bott-the-pigeon/app/error"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// Sends a bot usage help message from the provided bot.
func OnHelp(bot *discordgo.Session, msg *discordgo.MessageCreate) {
	_, err := bot.ChannelMessageSend(msg.ChannelID,
		"Hello! My command prefix is a `>`. Get it? Because it looks like a beak.\n\n"+
			"**Commands**:\n"+
			"`support`: Sends bot usage help - like you're seeing right now!\n"+
			"`pigeon`: Sends a random picture of a pigeon.\n"+
			"`todo {Some feature}`: Submit a suggestion to the project todo list.\n"+
			"`mc`: Starts the Minecraft server if it's not already running. It takes a few minutes though. If you want to run the Tekkit server, use `>mc --tekkit`.\n"+
			"`mc-status`: Checks the current status of the Minecraft server. If you want to know the status of the Tekkit server, use `>mc-status --tekkit`.\n\n"+
			"_That's all for now folks, because I'm a dumb bird._")
	if err != nil {
		err = fmt.Errorf("failed to send help message: %v", err)
		e.ThrowBotError(bot, msg.ChannelID, err)
		return
	}
}

// TODO: This is very hard-coded, and should be called from an API in some way.
