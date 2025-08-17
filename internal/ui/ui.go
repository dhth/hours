package ui

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	errFailedToConfigureDebugging = errors.New("failed to configure debugging")
	errCouldnCreateFramesDir      = errors.New("couldn't create frames directory")
)

func RenderUI(db *sql.DB, style Style) error {
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			return fmt.Errorf("%w: %s", errFailedToConfigureDebugging, err.Error())
		}
		defer f.Close()
	}

	debug := os.Getenv("HOURS_DEBUG") == "1"
	logFrames := os.Getenv("HOURS_LOG_FRAMES") == "1"
	if logFrames {
		err := os.MkdirAll("frames", 0o755)
		if err != nil {
			return fmt.Errorf("%w: %s", errCouldnCreateFramesDir, err.Error())
		}
	}

	p := tea.NewProgram(InitialModel(db, style, debug, logFrames), tea.WithAltScreen())
	_, err := p.Run()

	return err
}
