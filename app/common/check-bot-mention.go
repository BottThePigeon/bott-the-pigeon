package common

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// Returns true if the provided list of mentions contains the provided bot.
func CheckForBotMention(bot *discordgo.Session, mentions []*discordgo.User) (bool, error) {
	user, err := bot.User("@me")
	if err != nil {
		return false, fmt.Errorf("failed to get bot session user: %v", err)
	}
	for _, m := range mentions {
		if m.ID == user.ID {
			return true, nil
		}
	}
	return false, nil
}
