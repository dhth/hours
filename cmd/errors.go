package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dhth/hours/internal/ui"
)

func handleErrors(err error) {
	if errors.Is(err, errCouldntGenerateData) {
		fmt.Printf("\n%s\n", msgReportIssue)
		return
	}

	if errors.Is(err, ErrThemeDoesntExist) {
		fmt.Printf(`
Run "hours themes list" to list themes or create a new one using "hours themes add".
`)
		return
	}

	if errors.Is(err, ui.ErrThemeFileHasInvalidSchema) {
		defaultTheme := ui.DefaultTheme()
		defaultThemeBytes, err := json.MarshalIndent(defaultTheme, "", "  ")
		if err != nil {
			return
		}

		fmt.Printf(`
A valid theme file looks like this:

%s
`, defaultThemeBytes)
	}
}
