package session

import (
	handlers "bott-the-pigeon/bot-utils/handlers"

	"log"

	"github.com/bwmarrin/discordgo"
)

// The Bot Session pointer is stored, and can be accessed later.
// This probably won't be needed very often and should usually
// just be passed between functions, since the Bot is the
// primary continuous flow of the application - but sometimes
// it may be inconvenient (or unclean) to pass the session.
var (
	bot *discordgo.Session
)

// Functions called by the initialisation process of a Discord bot.

// Function caller for bot initialisation steps
func GetBotSession(botTokenKey string) *discordgo.Session {

	// Return the stored Bot Session or create one if not.
	// Therefore, initialisation code should only run once.
	if bot != nil {
		return bot
	} else {
		// Instantiate Bot
		var err error
		bot, err = discordgo.New("Bot " + botTokenKey)
		if err != nil {
			log.Fatal("Could not initialise bot: ", err)
		}

		// Handlers can be added after opening a session,
		// so we can run these concurrently with no worries.
		go addHandlers(bot)
		go openBot(bot)

		return bot
	}
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
