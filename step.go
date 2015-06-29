package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"./markdownlog"
)

const (
	BASE_URL             = "https://api.hipchat.com/v1"
	RESPONSE_STATUS_SENT = "sent"
)

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
	pth := os.Getenv("BITRISE_STEP_FORMATTED_OUTPUT_FILE_PATH")
	markdownlog.Setup(pth)
	err := markdownlog.ClearLogFile()
	if err != nil {
		fmt.Errorf("Failed to clear log file", err)
	}

	// required inputs
	token := os.Getenv("HIPCHAT_TOKEN")
	if token == "" {
		errorMessageToOutput("$HIPCHAT_TOKEN is not provided!")
		os.Exit(1)
	}
	roomId := os.Getenv("HIPCHAT_ROOMID")
	if roomId == "" {
		errorMessageToOutput("$HIPCHAT_ROOMID is not provided!")
		os.Exit(1)
	}
	fromName := os.Getenv("HIPCHAT_FROMNAME")
	if fromName == "" {
		errorMessageToOutput("$HIPCHAT_FROMNAME is not provided!")
		os.Exit(1)
	}
	message := os.Getenv("HIPCHAT_MESSAGE")
	if message == "" {
		errorMessageToOutput("$HIPCHAT_MESSAGE is not provided!")
		os.Exit(1)
	}
	//optional inputs
	messageColor := os.Getenv("HIPCHAT_MESSAGE_COLOR")
	if messageColor == "" {
		markdownlog.SectionToOutput("$HIPCHAT_MESSAGE_COLOR is not provided, use default!")
		messageColor = "yellow"
	}
	errorFromName := os.Getenv("HIPCHAT_ERROR_FROMNAME")
	if errorFromName == "" {
		markdownlog.SectionToOutput("$HIPCHAT_ERROR_FROMNAME is not provided!")
	}
	errorMessage := os.Getenv("HIPCHAT_ERROR_MESSAGE")
	if errorMessage == "" {
		markdownlog.SectionToOutput("$HIPCHAT_ERROR_MESSAGE is not provided!")
	}

	isBuildFailedMode := (os.Getenv("STEPLIB_BUILD_STATUS") != "0")
	if isBuildFailedMode {
		if errorFromName == "" {
			fmt.Errorf("Build failed, but no HIPCHAT_ERROR_FROMNAME defined, use default")
		} else {
			fromName = errorFromName
		}
		if errorMessage == "" {
			fmt.Errorf("Build failed, but no HIPCHAT_ERROR_MESSAGE defined, use default")
		} else {
			message = errorMessage
		}
	}

	// request payload
	values := url.Values{
		"room_id": {roomId},
		"from":    {fromName},
		"message": {message},
		"color":   {messageColor},
	}
	valuesReader := *strings.NewReader(values.Encode())

	// request
	url := BASE_URL + "/rooms/message?auth_token=" + token

	request, err := http.NewRequest("POST", url, &valuesReader)
	if err != nil {
		fmt.Println("Failed to create requestuest:", err)
		os.Exit(1)
	}

	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// perform request
	client := &http.Client{}
	response, err := client.Do(request)
	if response.StatusCode == 200 {
		successMessageToOutput(fromName, roomId, message)
	} else {
		var data map[string]interface{}
		bodyBytes, _ := ioutil.ReadAll(response.Body)
		err := json.Unmarshal(bodyBytes, &data)
		if err == nil {
			fmt.Println("Response:", data)
		}

		errorMsg := fmt.Sprintf("Status code: %s Body: %s", response.StatusCode, response.Body)
		errorMessageToOutput(errorMsg)

		os.Exit(1)
	}
}
