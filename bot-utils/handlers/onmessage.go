package handlers

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

//  Handles Discord MessageCreate event and handles multiple relevant
func OnMessage(bot *discordgo.Session, msg *discordgo.MessageCreate) {

	// Respond with random response if the bot is mentioned
	if checkForBotMention(bot, msg.Mentions) {
		onTag(bot, msg)
	}

	if strings.Contains(msg.Content, "!pigeon") {
		onPigeon(bot, msg)
	}
}

// Bot response to being tagged in a group.
func onTag(bot *discordgo.Session, msg *discordgo.MessageCreate) {
	bot.ChannelMessageSend(msg.ChannelID, "G'day, am Scott de racing pigeon.");
}

// Bot response to "!pigeon" - sending an image.
func onPigeon(bot *discordgo.Session, msg *discordgo.MessageCreate) {
	
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