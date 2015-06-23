package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"./markdownlog"
)

const (
	BASE_URL           = "https://api.hipchat.com/v1"
	ResponseStatusSent = "sent"
)

var (
	token string

	roomId   string
	fromName string
	message  string

	messageFormat string
	notify        string
	messageColor  string

	errorFromName string
	errorMessage  string
)

type MessageRequest struct {
	// - required
	RoomId   string `json:"room_id"` // ID or name of the room.
	FromName string `json:"from"`    // Name the message will appear be sent from. (<15 chars)
	Message  string `json:"message"` // The message body. 10,000 characters max.
	// - optional
	MessageFormat string `json:"message_format"` // html or text (default: html)
	Notify        bool   `json:"notify"`         // Whether or not this message should trigger a notification for people in the room
	MessageColor  string `json:"color"`          // One of "yellow", "red", "green", "purple", "gray", or "random". (default: yellow)
}

type HipchatError struct {
	Code    int
	Type    string
	Message string
}

func (e HipchatError) Error() string {
	return e.Message
}

type ErrorResponse struct {
	Error HipchatError
}

type Client struct {
	AuthToken string
	BaseURL   string
}

func NewClient(authToken string) Client {
	return Client{AuthToken: authToken, BaseURL: BASE_URL}
}

func urlValuesFromMessageRequest(req MessageRequest) (url.Values, error) {
	payload := url.Values{
		"room_id": {req.RoomId},
		"from":    {req.FromName},
		"message": {req.Message},
	}
	if req.Notify == true {
		payload.Add("notify", "1")
	}
	if len(req.MessageColor) > 0 {
		payload.Add("color", req.MessageColor)
	}
	if len(req.MessageFormat) > 0 {
		payload.Add("message_format", req.MessageFormat)
	}
	return payload, nil
}

func (c *Client) PostMessage(req MessageRequest) error {
	if len(c.BaseURL) == 0 {
		c.BaseURL = BASE_URL
	}
	uri := fmt.Sprintf("%s/rooms/message?auth_token=%s", c.BaseURL, url.QueryEscape(c.AuthToken))

	payload, err := urlValuesFromMessageRequest(req)
	if err != nil {
		return err
	}

	resp, err := http.PostForm(uri, payload)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	msgResp := &struct{ Status string }{}
	if err := json.Unmarshal(body, msgResp); err != nil {
		return err
	}
	if msgResp.Status != ResponseStatusSent {
		return getError(body)
	}

	return nil
}

func getError(body []byte) error {
	var errResp ErrorResponse
	if err := json.Unmarshal(body, &errResp); err != nil {
		return err
	}
	return errResp.Error
}

func buildMessageRequest(isBuildFailedMode bool) MessageRequest {
	req := MessageRequest{
		RoomId: roomId,
	}

	if isBuildFailedMode {
		if errorFromName == "" {
			fmt.Println("Build failed, but no HIPCHAT_ERROR_FROMNAME defined, use default")
		} else {
			fromName = errorFromName
		}
	}
	req.FromName = fromName

	if isBuildFailedMode {
		if errorMessage == "" {
			fmt.Println("Build failed, but no HIPCHAT_ERROR_MESSAGE defined, use default")
		} else {
			message = errorMessage
		}
	}
	req.Message = message

	req.MessageColor = messageColor

	return req
}

func errorMessageToOutput(msg string) error {
	message := "Message send failed!\n"
	message = message + "Error message:\n"
	message = message + msg

	return markdownlog.ErrorSectionToOutput(message)
}

func successMessageToOutput(from, roomId, msg string) error {
	message := "Message successfully sent!\n"
	message = message + "From:\n"
	message = message + from + "\n"
	message = message + "To Romm:\n"
	message = message + roomId + "\n"
	message = message + "Message:\n"
	message = message + msg

	return markdownlog.SectionToOutput(message)
}

func main() {
	// init / cleanup the formatted output
	err := markdownlog.ClearLogFile()
	if err != nil {
		fmt.Errorf("Failed to clear log file", err)
	}

	// input validation
	// required
	token = os.Getenv("HIPCHAT_TOKEN")
	if token == "" {
		errorMessageToOutput("$HIPCHAT_TOKEN is not provided!")
		os.Exit(1)
	}
	roomId = os.Getenv("HIPCHAT_ROOMID")
	if roomId == "" {
		errorMessageToOutput("$HIPCHAT_ROOMID is not provided!")
		os.Exit(1)
	}
	fromName = os.Getenv("HIPCHAT_FROMNAME")
	if fromName == "" {
		errorMessageToOutput("$HIPCHAT_FROMNAME is not provided!")
		os.Exit(1)
	}
	message = os.Getenv("HIPCHAT_MESSAGE")
	if message == "" {
		errorMessageToOutput("$HIPCHAT_MESSAGE is not provided!")
		os.Exit(1)
	}
	//optional
	messageColor = os.Getenv("HIPCHAT_MESSAGE_COLOR")
	if messageColor == "" {
		markdownlog.SectionToOutput("$HIPCHAT_MESSAGE_COLOR is not provided!")
	}
	errorFromName = os.Getenv("HIPCHAT_ERROR_FROMNAME")
	if errorFromName == "" {
		markdownlog.SectionToOutput("$HIPCHAT_ERROR_FROMNAME is not provided!")
	}
	errorMessage = os.Getenv("HIPCHAT_ERROR_MESSAGE")
	if errorMessage == "" {
		markdownlog.SectionToOutput("$HIPCHAT_ERROR_MESSAGE is not provided!")
	}

	// perform step
	isBuildFailedMode := (os.Getenv("STEPLIB_BUILD_STATUS") != "0")
	req := buildMessageRequest(isBuildFailedMode)

	token := os.Getenv("HIPCHAT_TOKEN")
	c := NewClient(token)
	if err := c.PostMessage(req); err != nil {
		errorMessageToOutput(err.Error())
		os.Exit(1)
	}

	successMessageToOutput(req.FromName, req.RoomId, req.Message)
}
