package main

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

const (
	BASE_URL           = "https://api.hipchat.com/v1"
	ResponseStatusSent = "sent"
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
		RoomId: os.Getenv("key"),
	}

	fromName := os.Getenv("HIPCHAT_FROMNAME")
	if isBuildFailedMode {
		errorFromName := os.Getenv("HIPCHAT_ERROR_FROMNAME")
		if errorFromName == "" {
			log.Infoln("Build failed, but no HIPCHAT_ERROR_FROMNAME defined, use default")
		} else {
			fromName = errorFromName
		}
	}
	req.FromName = fromName

	message := os.Getenv("HIPCHAT_MESSAGE")
	if isBuildFailedMode {
		errorMessage := os.Getenv("HIPCHAT_ERROR_MESSAGE")
		if errorMessage == "" {
			log.Infoln("Build failed, but no HIPCHAT_ERROR_MESSAGE defined, use default")
		} else {
			message = errorMessage
		}
	}
	req.Message = message

	req.MessageColor = os.Getenv("HIPCHAT_MESSAGE_COLOR")

	return req
}

func main() {
	isBuildFailedMode := (os.Getenv("STEPLIB_BUILD_STATUS") != "0")
	req := buildMessageRequest(isBuildFailedMode)

	token := os.Getenv("HIPCHAT_TOKEN")
	c := NewClient(token)
	if err := c.PostMessage(req); err != nil {
		log.Printf("Expected no error, but got %q", err)
	}
}
