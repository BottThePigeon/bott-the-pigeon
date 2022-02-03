package init

import (
	"bott-the-pigeon/bot-utils/handlers"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Functions called by the initialisation process of a Discord bot.

// Function caller for bot initialisation steps
func InitBot(botTokenKey string) (*discordgo.Session) {

	// Instantiate Bot
	bot, err := discordgo.New("Bot " + os.Getenv(botTokenKey))
	if err != nil {
		log.Fatal("Could not initialise bot: ", err)
	}

	addHandlers(bot)
	openBot(bot)
	addListeners()

	return bot
}

// Declarative event handler attachment
func addHandlers(bot *discordgo.Session) {

	// Handle messages sent in group
	bot.AddHandler(handlers.OnMessage)
	bot.Identify.Intents = discordgo.IntentsGuildMessages
}

// Opens the provided session, handling errors
func openBot(s *discordgo.Session) {

	err := s.Open()
	if err != nil {
		log.Fatal("Could not open bot session: ", err)
		return
	}

	// Return success message
	fmt.Println("Bot is running.");
}

// Make external event listener channel for bot
func addListeners() {

	//Notify for when killed. TODO: + `os.Kill` maybe? But that throws untrappable warning
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}