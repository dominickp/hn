package logger

import (
	"io"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

var Logger = log.New(os.Stdout, "", log.LstdFlags)

func Init(logfilePath string) {
	if logfilePath != "" {
		// Only capture logs if the log file path is available (e.g. when executing from this repo -> ./logs)
		f, err := tea.LogToFile(logfilePath, "simple")
		if err != nil {
			Logger = log.New(io.Discard, "", 0) // switch to a dummy logger
		} else {
			Logger = log.New(f, "", log.LstdFlags) // switch to a file logger
		}
	}

}
