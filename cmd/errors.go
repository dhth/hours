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
		return
	}

	if errors.Is(err, ui.ErrThemeColorsAreInvalid) {
		fmt.Printf(`
Colors codes can only be provided in ANSI 16, ANSI 256, or HEX formats.

For example:

"activeTask": "9"           # red in ANSI 16
"activeTask": "201"         # hot pink in ANSI 256
"activeTask": "#0000FF"     # blue in HEX (true color)

Fun fact: There are 16,777,216 true color choices. Go nuts.
`)
		return
	}
}
