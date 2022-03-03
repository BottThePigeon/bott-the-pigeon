package main

import (
	bot "bott-the-pigeon/app/session"
	ssm "bott-the-pigeon/lib/aws/service/ssm-env"
	aws "bott-the-pigeon/lib/aws/session"
	"fmt"
	"log"

	"flag"
	"os"
	"os/signal"
	"syscall"
)

// Note: main should give a high-level overview of the E2E flow of the application.

func main() {

	// This is the only place where logs can (should) be fatal, and terminate the app.
	// Flag management; whether to use test or prod configs.
	config := *flagHandler()
	botTokenKey := getBotTokenKey(*config.prod)

	// Set any non-confidential configs (stored in repo, in getConfigs()).
	err := setEnvs(getConfigs())
	if err != nil {
		log.Fatal(err)
	}

	// Initialise a session with AWS.
	_, err = aws.GetSession()
	if err != nil {
		log.Fatal(err)
	}

	// Get credentials from SSM not to be stored in repo.
	ssmEnv, err := ssm.Getenv(os.Getenv("AWS_SSM_PARAMETER_PATH"))
	if err != nil {
		log.Fatal(err)
	}

	// Add new environment variables based on those returned from SSM.
	err = setEnvs(ssmEnv)
	if err != nil {
		log.Fatal(err)
	}

	// Create a Discord bot session. This is where most of the magic happens.
	bot, err := bot.GetSession(os.Getenv(botTokenKey))
	if err != nil {
		log.Fatal(err)
	}

	// Closure functions
	defer bot.Close()
	addCloseListener()
}

// Returns a k,v map of base configs. NON-SENSITIVE CONFIGS GO HERE.
func getConfigs() map[string]string {
	env := make(map[string]string)
	env["GITHUB_REPO_ACCOUNT"] = "BottThePigeon"
	env["GITHUB_PROJECT_ID"] = "1"
	env["GITHUB_SUGGESTIONS_COLUMN_ID"] = "17943099"
	env["AWS_SSM_PARAMETER_PATH"] = "/btp/"
	env["AWS_CW_ERROR_LOG_GROUP"] = "bot-error"
	return env
}

// Initialises environment based upon provided k,v map.
func setEnvs(env map[string]string) error {
	errs := []error{}
	for k, v := range env {
		err := os.Setenv(k, v)
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("failed env variable initialisation. error(s): %v", errs)
	}
	return nil
}

// Flag configurations for the application.
type flagConfig struct {
	prod *bool
}

// Parses the flag configurations for the application.
func flagHandler() *flagConfig {
	flags := &flagConfig{
		prod: flag.Bool("prod",
			false,
			"Should the production bot application be used?"),
	}
	flag.Parse()
	return flags
}

// Returns the token environment variable key based on isProd.
func getBotTokenKey(isProd bool) string {
	botTokenKey := "BOT_TOKEN_TEST"
	if isProd {
		botTokenKey = "BOT_TOKEN"
	}
	return botTokenKey
}

// Waits for a termination/kill etc. signal (Holding the application open).
func addCloseListener() {

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sigChan
}
