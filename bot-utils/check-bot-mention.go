package botutils

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

// Simple loop encapsulation to check if user has been mentioned
func CheckForBotMention(bot *discordgo.Session, mentions []*discordgo.User) bool {
	user, err := bot.User("@me")

	if err != nil {
		log.Fatal("Could not get session current user: ", err)
	}

	for _, m := range mentions {
		if m.ID == user.ID {
			return true
		}
	}

	return false
}
