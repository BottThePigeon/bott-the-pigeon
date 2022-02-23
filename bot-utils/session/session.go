package session

import (
	handlers "bott-the-pigeon/bot-utils/handlers"
	"fmt"

	"log"

	"github.com/bwmarrin/discordgo"
)

// The Bot Session pointer is stored, and can be accessed later.
var (
	bot *discordgo.Session
)

// Returns the stored Bot Session or creates one if it doesn't exist, using the provided token.
func GetBotSession(token string) (*discordgo.Session, error) {
	if bot != nil {
		return bot, nil
	} else {
		bot, err := discordgo.New("Bot " + token)
		if err != nil {
			return nil, err
		}
		go addMessageHandler(bot)
		go openBot(bot)
		return bot, nil
	}
}

// Adds OnMessage handler.
func addMessageHandler(bot *discordgo.Session) {
	bot.AddHandler(handlers.OnMessage)
	bot.Identify.Intents = discordgo.IntentsGuildMessages
}

// Opens the provided bot session.
func openBot(bot *discordgo.Session) error {
	err := bot.Open()
	if err != nil {
		return fmt.Errorf("failed to open bot session: %v", err)
	}
	log.Println("Bot is running.")
	return nil
}
