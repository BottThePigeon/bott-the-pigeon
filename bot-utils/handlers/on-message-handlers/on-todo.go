package onmessagehandlers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

//TODO: Refactor

type Response_Post_GHProjectCard struct {
	ID int `json:"id"`
}

func OnTodo(bot *discordgo.Session, msg *discordgo.MessageCreate) {

	// Remove the leading command and trim spaces either side of the output
	crText := strings.TrimSpace(strings.Replace(msg.Content, ">todo", "", 1))

	// "note" key is used to define the text content for a new note
	postBody, _ := json.Marshal(map[string]string{
		"note": crText,
	})

	reqBody := bytes.NewBuffer(postBody)
	req, _ := http.NewRequest("POST", "https://api.github.com/projects/columns/"+
		os.Getenv("GITHUB_SUGGESTIONS_COLUMN_ID")+
		"/cards",
		reqBody)

	req.Header.Add("Authorization", "Bearer " + os.Getenv("GITHUB_PROJECTS_ACCESS_TOKEN"))
	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil || resp.StatusCode != 201 {

		log.Println("HTTP request failed.")

		_, err := bot.ChannelMessageSendEmbed(msg.ChannelID, &discordgo.MessageEmbed{
			Title:       "Uh-oh. Something went wrong.",
			Description: "As you know, I'm a pigeon, so things like this happen. Please don't kill me.",
			Color:       0xDD4400,
		})

		if err != nil {
			log.Println("Bot could not send a message - something went seriously wrong: ", err)
		}

		return
	}

	defer resp.Body.Close()

	byteResp, err := ioutil.ReadAll(resp.Body)

	if err != nil {

		log.Println("Could not read response body: ", err)

		_, err := bot.ChannelMessageSendEmbed(msg.ChannelID, &discordgo.MessageEmbed{
			Title:       "Uh-oh. Something went wrong.",
			Description: "As you know, I'm a pigeon, so things like this happen. Please don't kill me.",
			Color:       0xDD4400,
		})

		if err != nil {
			log.Println("Bot could not send a message - something went seriously wrong: ", err)
		}

		return
	}

	var jsonResp Response_Post_GHProjectCard
	json.Unmarshal(byteResp, &jsonResp)

	ghProjectCardLink :=
		"https://github.com/orgs/" +
			os.Getenv("GITHUB_REPO_ACCOUNT") +
			"/projects/" +
			os.Getenv("GITHUB_PROJECT_ID") +
			"#card-" +
			strconv.Itoa(jsonResp.ID)

	bot.ChannelMessageSendEmbed(msg.ChannelID, &discordgo.MessageEmbed{
		Title:       "To-Do added successfully.",
		Description: "Cheers lad, your suggestion was added to the To-Do list [here](" + ghProjectCardLink + ").",
		Color:       0x44DD00,
	})
}
