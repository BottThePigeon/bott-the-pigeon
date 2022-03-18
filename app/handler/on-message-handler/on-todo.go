package onmessagehandlers

import (
	e "bott-the-pigeon/app/error"
	httputil "bott-the-pigeon/lib/http"
	client "bott-the-pigeon/lib/http/util"
	"fmt"
	"strings"

	"encoding/json"
	"os"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

// Response fields that should be captured for POST GitHub Create Project Card.
type Response_Post_GHProjectCard struct {
	ID int `json:"id"`
}

// Creates a new Todo on the repo's associated GitHub project.
func OnTodo(bot *discordgo.Session, msg *discordgo.MessageCreate) error {
	todo := strings.TrimSpace(strings.Replace(msg.Content, ">todo", "", 1))
	cardID, err := createGHTodo(os.Getenv("GITHUB_PROJECTS_ACCESS_TOKEN"), os.Getenv("GITHUB_SUGGESTIONS_COLUMN_ID"), todo)
	if err != nil {
		e.ThrowBotError(bot, msg.ChannelID, err)
		return err
	}
	cardLink := genGHProjectCardLink(os.Getenv("GITHUB_REPO_ACCOUNT"), os.Getenv("GITHUB_PROJECT_ID"), strconv.Itoa(*cardID))
	res := genGHTodoSuccessMessage(cardLink)
	_, err = bot.ChannelMessageSendEmbed(msg.ChannelID, res)
	if err != nil {
		e.ThrowBotError(bot, msg.ChannelID, err)
		return err
	}
	return nil
}

// Sends a todo to the repo's GitHub project using the provided text, and returns the card ID.
func createGHTodo(ghAccessToken string, ghProjectColumnID string, todoText string) (*int, error) {
	reqBody := map[string]string{
		"note": todoText,
	}
	reqParams := &httputil.HTTP_Request{
		Method: "POST",
		URL: 	fmt.Sprintf("https://api.github.com/projects/columns/%v/cards", ghProjectColumnID),
		Body:	reqBody,
	}

	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + ghAccessToken
	successCode := 201
	r, err := client.CreateDoHTTPRequest(*reqParams, headers, successCode)
	if err != nil {
		return nil, err
	}
	res, err := parseGHTodoPostResponseBody(r)
	if err != nil {
		return nil, err
	}
	return &res.ID, nil
}

// Generates the link for a GitHub project card, given its ID.
func genGHProjectCardLink(account string, projectID string, cardID string) string {
	link := fmt.Sprintf("https://github.com/orgs/%v/projects/%v#card-%v", account, projectID, cardID)
	return link
}

// Returns a Discord message for successful GitHub Todo creation.
func genGHTodoSuccessMessage(link string) *discordgo.MessageEmbed {
	msg := &discordgo.MessageEmbed{
		Title:       "To-Do added successfully.",
		Description: fmt.Sprintf("Cheers lad, your suggestion was added to the To-Do list [here](%v).", link),
		Color:       0x44DD00,
	}
	return msg
}
// This may get reused in future for other success messages, but it can be left for now.


// Maps and unmarshals provided response into Response_Post_GHProjectCard type.
func parseGHTodoPostResponseBody(resp []byte) (*Response_Post_GHProjectCard, error) {
	var respJson Response_Post_GHProjectCard
	err := json.Unmarshal(resp, &respJson)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response body. error: %v", err)
	}
	return &respJson, nil
}