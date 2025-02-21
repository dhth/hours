package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/dhth/hours/internal/ui"
)

const themeNameRegexPattern = `^[a-zA-Z0-9-]{1,20}$`

var themeNameRegExp = regexp.MustCompile(themeNameRegexPattern)

var (
	errThemeNameInvalid        = fmt.Errorf("theme name is invalid; valid regex: %s", themeNameRegexPattern)
	errMarshallingDefaultTheme = errors.New("couldn't marshall default theme to bytes")
	errCouldntCreateThemesDir  = errors.New("couldn't create themes directory")
	errCouldntCreateThemeFile  = errors.New("couldn't create theme file")
	errCouldntWriteToThemeFile = errors.New("couldn't write to theme file")
)

func addTheme(themeName string, themesDir string) (string, error) {
	var zero string
	if !themeNameRegExp.MatchString(themeName) {
		return zero, errThemeNameInvalid
	}

	defaultTheme := ui.DefaultTheme()
	themeBytes, err := json.MarshalIndent(defaultTheme, "", "  ")
	if err != nil {
		return zero, fmt.Errorf("%w: %s", errMarshallingDefaultTheme, err.Error())
	}

	err = os.MkdirAll(themesDir, 0o755)
	if err != nil {
		return zero, fmt.Errorf("%w: %s", errCouldntCreateThemesDir, err.Error())
	}

	themePath := filepath.Join(themesDir, fmt.Sprintf("%s.json", themeName))

	file, err := os.Create(themePath)
	if err != nil {
		return zero, fmt.Errorf("%w: %s", errCouldntCreateThemeFile, err.Error())
	}
	defer file.Close()

	_, err = file.Write(themeBytes)
	if err != nil {
		return zero, fmt.Errorf("%w: %s", errCouldntWriteToThemeFile, err.Error())
	}

	return themePath, nil
}
