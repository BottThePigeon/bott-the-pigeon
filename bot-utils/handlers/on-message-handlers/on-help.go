package onmessagehandlers

import "github.com/bwmarrin/discordgo"

// Provides a list of commands and their use to help the user.
func OnHelp(bot *discordgo.Session, msg *discordgo.MessageCreate) {

	// TODO: This could get massive, so suggest using a script
	// to detect commands in the code and add them to a hosted file we can use.

	bot.ChannelMessageSend(msg.ChannelID,
		"Hello! My command prefix is a `>`. Get it? Because it looks like a beak.\n\n"+
			"**Commands**:\n"+
			"`pigeon`: Send a random picture of a pigeon.\n"+
			"`todo {Some feature}`: Submit a suggestion to the project todo list.\n\n"+
			"_That's all for now, because I'm a dumb bird._")
}
