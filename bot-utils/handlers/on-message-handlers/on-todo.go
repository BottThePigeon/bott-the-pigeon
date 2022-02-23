package onmessagehandlers

import (
	e "bott-the-pigeon/bot-utils/errors"
	"errors"
	"fmt"
	"strings"
	"time"

	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
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
	cardID, err := botCreateGitHubTodo(os.Getenv("GITHUB_PROJECTS_ACCESS_TOKEN"), os.Getenv("GITHUB_SUGGESTIONS_COLUMN_ID"), todo)
	if err != nil {
		e.ThrowBotError(bot, msg.ChannelID, err)
		return err
	}
	cardLink := genGHProjectCardLink(os.Getenv(""), os.Getenv("GITHUB_PROJECT_ID"), strconv.Itoa(*cardID))
	res := getGitHubTodoSuccessMessage(cardLink)
	_, err = bot.ChannelMessageSendEmbed(msg.ChannelID, res)
	if err != nil {
		e.ThrowBotError(bot, msg.ChannelID, err)
		return err
	}
	return nil
}

// Sends a todo to the repo's GitHub project using the provided text.
func botCreateGitHubTodo(ghAccessToken string, ghProjectColumnID string, todoText string) (*int, error) {
	req, err := genGHProjectCardPostRequest(todoText, ghProjectColumnID, ghAccessToken)
	if err != nil {
		return nil, err
	}

	res, err := doGHProjectCardPostRequest(req)
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
func getGitHubTodoSuccessMessage(link string) *discordgo.MessageEmbed {
	msg := &discordgo.MessageEmbed{
		Title:       "To-Do added successfully.",
		Description: fmt.Sprintf("Cheers lad, your suggestion was added to the To-Do list [here](%v).", link),
		Color:       0x44DD00,
	}
	return msg
}

// Generates a GitHub Project Card creation HTTP POST request
func genGHProjectCardPostRequest(todoText string, ghProjectColumnID string, ghAccessToken string) (*http.Request, error) {
	reqBody, err := genGHProjectCardPostBody(todoText)
	if err != nil {
		return nil, err
	}

	reqBuffer := bytes.NewBuffer(reqBody)
	req, err := getGHProjectCardPostRequest(reqBuffer, ghProjectColumnID)
	if err != nil {
		return nil, err
	}

	addGHProjectCardPostRequestHeaders(req, ghAccessToken)
	return req, nil
}

// Sends the provided request and reads the response for GitHub Project Card creation POST
func doGHProjectCardPostRequest(req *http.Request) (*Response_Post_GHProjectCard, error) {
	httpResp, err := doRequest(req, 201)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	readResp, err := readResponse(httpResp)
	if err != nil {
		return nil, err
	}

	objResp, err := parseGHProjectCardPostResponse(readResp)
	if err != nil {
		return nil, err
	}

	return objResp, nil
}

// Generates a body for POST request to create a new GitHub project card
func genGHProjectCardPostBody(cardText string) ([]byte, error) {
	body, err := json.Marshal(map[string]string{
		"note": cardText,
	}) 
	if err != nil {
		return nil, fmt.Errorf("failed to create POST request JSON body: %v", err)
	}
	return body, nil
}

// Generates a HTTP POST request to create a new GitHub project card using a body buffer
func getGHProjectCardPostRequest(buffer *bytes.Buffer, columnID string) (*http.Request, error) {
	endpoint := fmt.Sprintf("https://api.github.com/projects/columns/%v/cards", columnID)
	req, err := http.NewRequest("POST", endpoint, buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to create post request: %v", err)
	}
	return req, nil
}

// Adds the necessary headers to the provided GitHub Project Card POST request
func addGHProjectCardPostRequestHeaders(req *http.Request, ghAccessToken string) *http.Request {
	req.Header.Add("Authorization", "Bearer " + os.Getenv("GITHUB_PROJECTS_ACCESS_TOKEN"))
	return req
}

// Generic HTTP caller function using provided request (TODO: Maybe extract into http-utils?)
func doRequest(req *http.Request, mustStatus int) (*http.Response, error) {
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)

	//Throw an error if the request fails, or the response status code doesn't match mustStatus
	if err != nil {
		return nil, fmt.Errorf("failed http request: %v", err)
	} else if resp.StatusCode != mustStatus {
		err := errors.New("got bad status code, expected " +
		strconv.Itoa(mustStatus) + ", got " + 
		resp.Status)
		return nil, fmt.Errorf("failed http request: %v", err)
	}
	return resp, nil
}

// Generic HTTP response reader returning as []byte (TODO: Maybe extract into http-utils?)
func readResponse(resp *http.Response) ([]byte, error) {
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}
	return respBody, nil
}

// Maps and unmarshals provided response into Response_Post_GHProjectCard type
func parseGHProjectCardPostResponse(resp []byte) (*Response_Post_GHProjectCard, error) {
	var respJson Response_Post_GHProjectCard
	err := json.Unmarshal(resp, &respJson)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response body. error: %v", err)
	}
	return &respJson, nil
}