package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/dhth/hours/internal/ui/theme"
)

func handleError(err error) {
	if errors.Is(err, errCouldntGenerateData) {
		fmt.Fprintf(os.Stderr, "\n%s\n", msgReportIssue)
		return
	}

	if errors.Is(err, theme.ErrBuiltInThemeDoesntExist) {
		fmt.Fprintf(os.Stderr, `
If you intended to use a custom theme, prefix it with "custom:". Run "hours themes list" to list all themes. 
`)
		return
	}

	if errors.Is(err, theme.ErrCustomThemeDoesntExist) {
		fmt.Fprintf(os.Stderr, `
Run "hours themes list" to list custom themes.
`)
		return
	}

	if errors.Is(err, theme.ErrThemeFileHasInvalidSchema) {
		defaultTheme := theme.Default()
		defaultThemeBytes, err := json.MarshalIndent(defaultTheme, "", "  ")
		if err != nil {
			return
		}

		fmt.Fprintf(os.Stderr, `
A valid theme file looks like this:

%s
`, defaultThemeBytes)
		return
	}

	if errors.Is(err, theme.ErrThemeColorsAreInvalid) {
		fmt.Fprintf(os.Stderr, `
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
