package main

import (
	awsenv "bott-the-pigeon/aws-utils/aws-env"
	aws "bott-the-pigeon/aws-utils/session"
	bot "bott-the-pigeon/bot-utils/session"

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
		awssess := aws.GetAWSSession()
		awsenv.InitEnv(awssess)
	}()

	var botTokenKey string
	go func() {
		defer wg.Done()
		config := *flagHandler()

		//Works with the presence of --prod as a flag
		botTokenKey = getBotTokenKey(*config.prod)
	}()

	wg.Wait()

	// Return a bot instance. Everything happens in here
	bot := bot.GetBotSession(os.Getenv(botTokenKey))
	defer bot.Close()

	addCloseListener()
}

// Miscellaneous (and non-confidential) environment variable initialisation
// (That doesn't need AWS) goes here
func initEnv() {
	os.Setenv("GITHUB_REPO_ACCOUNT", "BottThePigeon")	  //The org that the repo belongs to
	os.Setenv("GITHUB_PROJECT_ID", "1")                   // The ID of the repo-associated project
	os.Setenv("GITHUB_SUGGESTIONS_COLUMN_ID", "17803319") // Where GH suggestions go
	os.Setenv("AWS_REGION", "eu-west-2")                  // AWS SDK Session Region
	os.Setenv("AWS_SSM_PARAMETER_PATH", "/btp/")          // SSM location of project variables
	// The EC2 instance shouldn't have permission to parameters outside this path
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
