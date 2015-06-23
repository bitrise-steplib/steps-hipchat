package outpututil

import (
	"io"
	"os"

	log "github.com/Sirupsen/logrus"
)

func IsPathExists(pth string) (bool, error) {
	if pth == "" {
		return false, errors.New("No path provided")
	}
	_, err := os.Stat(pth)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func stringToOutput(msg string) error {
	var w io.Writer

	pth := os.Getenv("BITRISE_STEP_FORMATTED_OUTPUT_FILE_PATH")

	if pth != "" {
		f, err := os.OpenFile(pth, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			t.Fatalf("Failed to open file: %v", err)
		}
		defer func() {
			err := file.Close()
			if err != nil {
				log.Fatalln("Failed to close file:", err)
			}
		}()

		w = io.MultiWriter(f, os.Stdout)
	} else {
		w = io.MultiWriter(f, os.Stdout)
	}

	n, err := w.Write(msg)
	if err != nil {
		log.Fatalln("Failed to write message:", err)
	}

	return nil
}
