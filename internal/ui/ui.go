package ui

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dhth/hours/internal/types"
)

var (
	errFailedToConfigureDebugging = errors.New("failed to configure debugging")
	errCouldnCreateFramesDir      = errors.New("couldn't create frames directory")
)

func RenderUI(db *sql.DB, style Style, timeProvider types.TimeProvider) error {
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			return fmt.Errorf("%w: %s", errFailedToConfigureDebugging, err.Error())
		}
		defer f.Close()
	}

	debug := os.Getenv("HOURS_DEBUG") == "1"
	logFrames := os.Getenv("HOURS_LOG_FRAMES") == "1"
	logFramesCfg := logFramesConfig{
		log: logFrames,
	}
	if logFrames {
		framesDir := filepath.Join(".frames", fmt.Sprintf("%d", timeProvider.Now().Unix()))
		err := os.MkdirAll(framesDir, 0o755)
		if err != nil {
			return fmt.Errorf("%w: %s", errCouldnCreateFramesDir, err.Error())
		}
		logFramesCfg.framesDir = framesDir
	}

	p := tea.NewProgram(
		InitialModel(
			db,
			style,
			timeProvider,
			debug,
			logFramesCfg,
		),
		tea.WithAltScreen(),
	)
	_, err := p.Run()

	return err
}
