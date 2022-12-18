package api

import (
	e "bott-the-pigeon/app/error"

	"github.com/bwmarrin/discordgo"

	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

var _bot *discordgo.Session

type Message struct {
	General bool
	Message string
}

func broadcastMsg(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	var msg Message
	err := decoder.Decode(&msg)
	if err != nil {
		err = fmt.Errorf("failed to send help message: %v", err)
		e.ThrowBotError(_bot, os.Getenv("BOSS_CHANNEL_ID"), err)
		return
	} else {
		var channel string
		if msg.General {
			channel = os.Getenv("BOSS_GENERAL_CHANNEL_ID")
		} else {
			channel = os.Getenv("BOSS_BIRDHOUSE_CHANNEL_ID")
		}
		_, err := _bot.ChannelMessageSend(channel,
			msg.Message)
		if err != nil {
			err = fmt.Errorf("failed to send help message: %v", err)
			e.ThrowBotError(_bot, os.Getenv("BOSS_CHANNEL_ID"), err)
			return
		}
	}
}

func Handler(bot *discordgo.Session) {
	_bot = bot
	http.HandleFunc("/broadcast", broadcastMsg)
	log.Fatal(http.ListenAndServe(":8192", nil))
}
