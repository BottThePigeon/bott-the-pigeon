package onmessagehandlers

import (
	e "bott-the-pigeon/bot-utils/errors"
	"errors"
	"time"

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

// TODO: Create generic error handling wrapper for functions.

// Add new key/values as required - no point storing values we don't need.
type Response_Post_GHProjectCard struct {
	ID int `json:"id"`
}

// Bot sends provided To-Do to GH project
func OnTodo(bot *discordgo.Session, msg *discordgo.MessageCreate) {
	req := initGHProjectCardPostRequest(bot, msg)
	doGHProjectCardPostRequest(bot, msg, req)
}

// Initialise a GitHub Project Card creation HTTP POST request
func initGHProjectCardPostRequest(bot *discordgo.Session, msg *discordgo.MessageCreate) (*http.Request) {
	todoText := trimTodoMessage(msg.Content)
	reqBody, err := genGHProjectCardPostBody(todoText)

	if err != nil {
		e.ThrowBotError(bot, msg)
		return nil
	}

	reqBuffer := bytes.NewBuffer(reqBody)
	req, err := genGHProjectCardPostRequest(reqBuffer)

	if err != nil {
		e.ThrowBotError(bot, msg)
		return nil
	}

	addGHProjectCardPostRequestHeaders(req)
	return req
}

// Send provided request and manage (read, parse, etc.) response for GitHub Project Card creation POST
func doGHProjectCardPostRequest(bot *discordgo.Session, msg *discordgo.MessageCreate, req *http.Request) {
	httpResp, err := doRequest(req, 201)
	if err != nil {
		e.ThrowBotError(bot, msg)
		return
	}

	defer httpResp.Body.Close()

	readResp, err := readResponse(httpResp)
	if err != nil {
		e.ThrowBotError(bot, msg)
		return
	}

	objResp, err := parseGHProjectCardPostResponse(readResp)
	if err != nil {
		e.ThrowBotError(bot, msg)
		return
	}

	ghProjectCardLink := genGHProjectCardLink(objResp.ID)
	bot.ChannelMessageSendEmbed(msg.ChannelID, &discordgo.MessageEmbed{
		Title:       "To-Do added successfully.",
		Description: "Cheers lad, your suggestion was added to the To-Do list [here](" + ghProjectCardLink + ").",
		Color:       0x44DD00,
	})
}

// Trims todo message to remove command and trailing spaces
func trimTodoMessage(msg string) string {
	str := strings.TrimSpace(strings.Replace(msg, ">todo", "", 1))
	return str
}

// Generate a body for POST request to create a new GitHub project card
func genGHProjectCardPostBody(cardText string) ([]byte, error) {
	body, err := json.Marshal(map[string]string{
		"note": cardText,
	}) 

	if err != nil {
		log.Println("Could not create POST request JSON body: ", err)
		return nil, err
	}

	return body, nil
}

// Generate a HTTP POST request to generate a new GitHub project card using a body buffer
func genGHProjectCardPostRequest(buffer *bytes.Buffer) (*http.Request, error) {
	req, err := http.NewRequest("POST", "https://api.github.com/projects/columns/"+
	os.Getenv("GITHUB_SUGGESTIONS_COLUMN_ID")+
	"/cards",
	buffer)

	if err != nil {
		log.Println("Could not create POST request: ", err)
		return nil, err
	}

	return req, nil
}

// Add the necessary headers to the provided GitHub Project Card POST request
func addGHProjectCardPostRequestHeaders(req *http.Request) *http.Request {

	// Basic Auth - Personal Access Token for Bott-The-Pigeon Bot GitHub account
	req.Header.Add("Authorization", "Bearer " + os.Getenv("GITHUB_PROJECTS_ACCESS_TOKEN"))
	return req
}

// Generic HTTP caller function using provided request (TODO: Maybe extract into http-utils?)
func doRequest(req *http.Request, mustStatus int) (*http.Response, error) {
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)

	//Throw an error if the request fails, or the response status code doesn't match expected
	if err != nil {

		log.Println("HTTP request failed: ", err)
		return nil, err

	} else if resp.StatusCode != mustStatus {

		err := errors.New("Got bad status code, expected " +
		strconv.Itoa(mustStatus) + ", got " + 
		strconv.Itoa(resp.StatusCode) + ": " + 
		resp.Status)
		return nil, err

	}

	return resp, nil
}

// Generic HTTP response reader returning as []byte (TODO: Maybe extract into http-utils?)
func readResponse(resp *http.Response) ([]byte, error) {
	respBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println("Could not read response body: ", err)
		return nil, err
	}

	return respBody, nil
}

// Maps and unmarshals provided response into Response_Post_GHProjectCard type
func parseGHProjectCardPostResponse(resp []byte) (*Response_Post_GHProjectCard, error) {
	var respJson Response_Post_GHProjectCard
	err := json.Unmarshal(resp, &respJson)

	if err != nil {
		log.Println("Could not parse response body: ", err)
		return nil, err
	}

	return &respJson, nil
}

// Generate the link for a GitHub project card, given its ID
func genGHProjectCardLink(cardID int) string {
	link :=
		"https://github.com/orgs/" +
			os.Getenv("GITHUB_REPO_ACCOUNT") +
			"/projects/" +
			os.Getenv("GITHUB_PROJECT_ID") +
			"#card-" +
			strconv.Itoa(cardID)

	return link
}