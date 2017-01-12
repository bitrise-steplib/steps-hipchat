package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/bitrise-io/go-utils/log"
)

// -----------------------
// --- Constants
// -----------------------

const (
	baseURL = "https://api.hipchat.com/v1"
)

// -----------------------
// --- Functions
// -----------------------

func logFail(format string, v ...interface{}) {
	errorMsg := fmt.Sprintf(format, v...)
	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", errorMsg)
	os.Exit(1)
}

func logWarn(format string, v ...interface{}) {
	errorMsg := fmt.Sprintf(format, v...)
	fmt.Printf("\x1b[33;1m%s\x1b[0m\n", errorMsg)
}

func logInfo(format string, v ...interface{}) {
	fmt.Println()
	errorMsg := fmt.Sprintf(format, v...)
	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", errorMsg)
}

func logDetails(format string, v ...interface{}) {
	errorMsg := fmt.Sprintf(format, v...)
	fmt.Printf("  %s\n", errorMsg)
}

func logDone(format string, v ...interface{}) {
	errorMsg := fmt.Sprintf(format, v...)
	fmt.Printf("  \x1b[32;1m%s\x1b[0m\n", errorMsg)
}

func validateRequiredInput(key string) string {
	value := os.Getenv(key)
	if value == "" {
		logFail("missing required input: %s", key)
	}
	return value
}

func validateInput(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		logWarn("%s input not provided, use default!", key)
		return defaultValue
	}
	return value
}

func printConfig(roomID, fromName, message, color, fromNameOnError, messageOnError, colorOnError, messageFormat string) {
	log.Infof("Configs:")
	log.Printf("token: ***")
	log.Printf("romm_id: %s", roomID)
	log.Printf("from_name: %s", fromName)
	log.Printf("color: %s", color)
	log.Printf("message: %s", message)
	log.Printf("from_name_on_error: %s", fromNameOnError)
	log.Printf("message_on_error: %s", messageOnError)
	log.Printf("color_on_error: %s", colorOnError)
	log.Printf("message_format: %s", messageFormat)
}

// -----------------------
// --- Main
// -----------------------

func main() {
	//
	// Validate options
	token := validateRequiredInput("auth_token")
	roomID := validateRequiredInput("room_id")
	fromName := validateRequiredInput("from_name")
	message := validateRequiredInput("message")

	//optional inputs
	messageColor := validateInput("color", "yellow")

	errorFromName := validateInput("from_name_on_error", fromName)
	errorMessage := validateInput("message_on_error", message)
	errorMessageColor := validateInput("color_on_error", messageColor)

	messageFormat := validateInput("message_format", "text")

	isBuildFailedMode := (os.Getenv("STEPLIB_BUILD_STATUS") != "0")
	if isBuildFailedMode {
		fromName = errorFromName
		message = errorMessage
		messageColor = errorMessageColor
	}

	fmt.Println()

	printConfig(roomID, fromName, message, messageColor, errorFromName, errorMessage, errorMessageColor, messageFormat)

	//
	// Create request
	logInfo("Performing request")

	values := url.Values{
		"room_id":        {roomID},
		"from":           {fromName},
		"message":        {message},
		"color":          {messageColor},
		"message_format": {messageFormat},
	}
	valuesReader := *strings.NewReader(values.Encode())

	url := baseURL + "/rooms/message?auth_token=" + token

	request, err := http.NewRequest("POST", url, &valuesReader)
	if err != nil {
		logFail("Failed to create request, error: %s", err)
	}
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	response, requestErr := client.Do(request)

	defer response.Body.Close()
	contents, readErr := ioutil.ReadAll(response.Body)

	//
	// Process response

	// Error
	if requestErr != nil {
		if readErr != nil {
			logWarn("Failed to read response body, error: %#v", readErr)
		} else {
			logInfo("Response:")
			logDetails("status code: %d", response.StatusCode)
			logDetails("body: %s", string(contents))
		}
		logFail("Performing request failed, error: %#v", requestErr)
	}

	if response.StatusCode < 200 || response.StatusCode > 300 {
		if readErr != nil {
			logWarn("Failed to read response body, error: %#v", readErr)
		} else {
			logInfo("Response:")
			logDetails("status code: %d", response.StatusCode)
			logDetails("body: %s", string(contents))
		}
		logFail("Performing request failed, status code: %d", response.StatusCode)
	}

	// Success
	logDone("Request succed")

	logInfo("Response:")
	logDetails("status code: %d", response.StatusCode)
	logDetails("body: %s", contents)

	if readErr != nil {
		logFail("Failed to read response body, error: %#v", readErr)
	}
}
