package onmessagehandlers

import "github.com/bwmarrin/discordgo"

// Bot response to being tagged in a group
func OnTag(bot *discordgo.Session, msg *discordgo.MessageCreate) {
	bot.ChannelMessageSend(msg.ChannelID, "G'day, am Scott de racing pigeon. Tag me with \"support\" for help.")
}
