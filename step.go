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
	token := os.Getenv("auth_token")
	if token == "" {
		errorMessageToOutput("$auth_token is not provided!")
		os.Exit(1)
	}
	roomId := os.Getenv("room_id")
	if roomId == "" {
		errorMessageToOutput("$room_id is not provided!")
		os.Exit(1)
	}
	fromName := os.Getenv("from_name")
	if fromName == "" {
		errorMessageToOutput("$from_name is not provided!")
		os.Exit(1)
	}
	message := os.Getenv("message")
	if message == "" {
		errorMessageToOutput("$message is not provided!")
		os.Exit(1)
	}
	//optional inputs
	messageFormat := os.Getenv("message_format")
	if messageFormat == "" {
		markdownlog.SectionToOutput("$message_format is not provided, use default - html!")
		messageFormat = "html"
	}
	messageColor := os.Getenv("color")
	if messageColor == "" {
		markdownlog.SectionToOutput("$color is not provided, use default!")
		messageColor = "yellow"
	}
	errorFromName := os.Getenv("from_name_on_error")
	if errorFromName == "" {
		markdownlog.SectionToOutput("$from_name_on_error is not provided!")
	}
	errorMessage := os.Getenv("message_on_error")
	if errorMessage == "" {
		markdownlog.SectionToOutput("$message_on_error is not provided!")
	}
	errorMessageFormat := os.Getenv("message_on_error_format")
	if errorMessageFormat == "" {
		markdownlog.SectionToOutput("$message_on_error_format is not provided, use default - html!")
		errorMessageFormat = "html"
	}
	errorMessageColor := os.Getenv("color_on_error")
	if errorMessageColor == "" {
		markdownlog.SectionToOutput("$color_on_error is not provided, use default!")
	}

	isBuildFailedMode := (os.Getenv("STEPLIB_BUILD_STATUS") != "0")
	if isBuildFailedMode {
		if errorFromName == "" {
			fmt.Errorf("Build failed, but no from_name_on_error defined, use default")
		} else {
			fromName = errorFromName
		}
		if errorMessage == "" {
			fmt.Errorf("Build failed, but no message_on_error defined, use default")
		} else {
			message = errorMessage
		}
		if errorMessageFormat == "" {
			fmt.Errorf("Build failed, but no message_format_on_error defined, use default")
		} else {
			messageFormat = errorMessageFormat
		}
		if errorMessageColor == "" {
			fmt.Errorf("Build failed, but no color_on_error defined, use default")
		} else {
			messageColor = errorMessageColor
		}
	}

	// request payload
	values := url.Values{
		"room_id": {roomId},
		"from":    {fromName},
		"message": {message},
		"color":   {messageColor},
		"message_format": {messageFormat}
	}
	valuesReader := *strings.NewReader(values.Encode())

	// request
	url := BASE_URL + "/rooms/message?auth_token=" + token

	request, err := http.NewRequest("POST", url, &valuesReader)
	if err != nil {
		fmt.Println("Failed to create request:", err)
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
