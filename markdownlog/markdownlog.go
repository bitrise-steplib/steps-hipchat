package markdownlog

import (
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
)

func ClearLogFile() error {
	pth := os.Getenv("BITRISE_STEP_FORMATTED_OUTPUT_FILE_PATH")
	if pth != "" {
		err := os.Remove(pth)
		if err != nil {
			return err
		}

		log.Info("Log file cleared")
	} else {
		log.Error("No BITRISE_STEP_FORMATTED_OUTPUT_FILE_PATH defined")
	}

	return nil
}

func ErrorMessageToOutput(msg string) error {
	pth := os.Getenv("BITRISE_STEP_FORMATTED_OUTPUT_FILE_PATH")
	if pth != "" {
		f, err := os.OpenFile(pth, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		defer func() error {
			err := f.Close()
			if err != nil {
				return err
			}

			return nil
		}()

		f.Write([]byte(msg))
	} else {
		log.Errorln("No BITRISE_STEP_FORMATTED_OUTPUT_FILE_PATH defined")
	}

	lines := strings.Split(msg, "\n")
	for _, line := range lines {
		log.Infoln(line)
	}

	return nil
}

func ErrorSectionToOutput(section string) error {
	msg := "\n" + section + "\n"

	return ErrorMessageToOutput(msg)
}

func ErrorSectionStartToOutput(section string) error {
	msg := section + "\n"

	return ErrorMessageToOutput(msg)
}

func MessageToOutput(msg string) error {
	pth := os.Getenv("BITRISE_STEP_FORMATTED_OUTPUT_FILE_PATH")
	if pth != "" {
		f, err := os.OpenFile(pth, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		defer func() error {
			err := f.Close()
			if err != nil {
				return err
			}

			return nil
		}()

		f.Write([]byte(msg))
	} else {
		log.Error("No BITRISE_STEP_FORMATTED_OUTPUT_FILE_PATH defined")
	}

	lines := strings.Split(msg, "\n")
	for _, line := range lines {
		log.Infoln(line)
	}

	return nil
}

func SectionToOutput(section string) error {
	msg := "\n" + section + "\n"

	return MessageToOutput(msg)
}

func SectionStartToOutput(section string) error {
	msg := section + "\n"

	return MessageToOutput(msg)
}
