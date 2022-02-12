package main

import (
	awsenv "bott-the-pigeon/aws-utils/aws-env"
	aws "bott-the-pigeon/aws-utils/session"
	bot "bott-the-pigeon/bot-utils/init"

	"flag"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// Main should give a high-level overview of the E2E flow of the application.

func main() {

	// Run concurrently because it's just initialisation stuff,
	// all of which is ultimately for the bot.
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		// All need to be run sequentially
		initEnv()
		aws.GetAWSSession()
		awsenv.InitEnv()
	}()

	var botTokenKey string
	go func() {
		defer wg.Done()
		config := *flagHandler()

		//Works with the presence of --prod as a flag
		botTokenKey = getBotTokenKey(*config.prod)
	}()

	wg.Wait()

	// Return a bot instance. Everything happens in here.
	bot := bot.InitBot(os.Getenv(botTokenKey))
	defer bot.Close()

	addCloseListener()
}

// Miscellaneous (and non-confidential) environment variable initialisation
// (That doesn't need AWS) goes here
func initEnv() {
	os.Setenv("AWS_REGION", "eu-west-2")     // AWS SDK Session Region
	os.Setenv("SSM_PARAMETER_PATH", "/btp/") // SSM location of project variables.
	// The EC2 instance shouldn't have permission to parameters outside this path.
}

// Flag configurations for the application
type flagConfig struct {
	prod *bool
}

//Parse and return the flag configurations for the application
func flagHandler() *flagConfig {

	flags := &flagConfig{
		prod: flag.Bool("prod",
			false,
			"Should the production bot application be used?"),
	}

	flag.Parse()
	return flags
}

// Determine what the key for the bot token needed is, based on if running in prod
func getBotTokenKey(isProd bool) string {

	botTokenKey := "BOT_TOKEN_TEST"
	if isProd {
		botTokenKey = "BOT_TOKEN"
	}

	return botTokenKey
}

// Waits for a termination/kill etc. signal (Holding the application open.)
func addCloseListener() {

	//There's also os.Kill
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sigChan
}
