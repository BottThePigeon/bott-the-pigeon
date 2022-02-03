package handlers

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

// Contains handler function(s) for on Discord MessageCreate event

func OnMessage(bot *discordgo.Session, msg *discordgo.MessageCreate) {

	// Respond with random response if the bot is mentioned
	if checkForBotMention(bot, msg.Mentions) {
		onTagBot(bot, msg)
	}
}

// Bot response to being tagged in a group.
func onTagBot(bot *discordgo.Session, msg *discordgo.MessageCreate) {
	bot.ChannelMessageSend(msg.ChannelID, "G'day, am Scott de racing pigeon.");
}

// Simple loop encapsulation to check if user has been mentioned
func checkForBotMention(bot *discordgo.Session, mentions []*discordgo.User) (bool) {
	user, err := bot.User("@me")

	if err != nil {
		log.Fatal("Could not get session current user: ", err)
	}

	for i := 0; i < len(mentions); i++ {
		if mentions[i].ID == user.ID {
			return true
		}
	}
	
	return false;
}