package init

import (
	handlers "bott-the-pigeon/bot-utils/handlers"

	"log"

	"github.com/bwmarrin/discordgo"
)

// Functions called by the initialisation process of a Discord bot.

// Function caller for bot initialisation steps
func InitBot(botTokenKey string) *discordgo.Session {

	// Instantiate Bot
	bot, err := discordgo.New("Bot " + botTokenKey)
	if err != nil {
		log.Fatal("Could not initialise bot: ", err)
	}

	// Handlers can be added after opening a session,
	// so we can run these concurrently with no worries.
	go addHandlers(bot)
	go openBot(bot)

	return bot
}

// Declarative event handler attachment
func addHandlers(bot *discordgo.Session) {

	// Handle messages sent in group
	bot.AddHandler(handlers.OnMessage)
	bot.Identify.Intents = discordgo.IntentsGuildMessages
}

// Opens the provided session, handling errors
func openBot(bot *discordgo.Session) {

	err := bot.Open()
	if err != nil {
		log.Fatal("Could not open bot session: ", err)
		return
	}

	// Return success message
	log.Println("Bot is running.")
}
