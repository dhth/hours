package ui

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

var errFailedToConfigureDebugging = errors.New("failed to configure debugging")

func RenderUI(db *sql.DB, style Style) error {
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			return fmt.Errorf("%w: %s", errFailedToConfigureDebugging, err.Error())
		}
		defer f.Close()
	}

	debug := os.Getenv("HOURS_DEBUG") == "1"

	p := tea.NewProgram(InitialModel(db, style, debug), tea.WithAltScreen())
	_, err := p.Run()

	return err
}
