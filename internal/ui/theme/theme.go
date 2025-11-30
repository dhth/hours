package theme

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

const (
	defaultThemeName  = "default"
	customThemePrefix = "custom:"
)

var (
	errThemeFileIsInvalidJSON     = errors.New("theme file is not valid JSON")
	ErrThemeFileHasInvalidSchema  = errors.New("theme file's schema is incorrect")
	ErrThemeColorsAreInvalid      = errors.New("invalid colors provided")
	errCouldntReadCustomThemeFile = errors.New("couldn't read custom theme file")
	errCouldntLoadCustomTheme     = errors.New("couldn't load custom theme")
	errEmptyThemeNameProvided     = errors.New("empty theme name provided")
	ErrCustomThemeDoesntExist     = errors.New("custom theme doesn't exist")
	ErrBuiltInThemeDoesntExist    = errors.New("built-in theme doesn't exist")
)

var hexCodeRegex = regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`)

type Theme struct {
	ActiveTask              string   `json:"activeTask,omitempty"`
	ActiveTaskBeginTime     string   `json:"activeTaskBeginTime,omitempty"`
	ActiveTasks             string   `json:"activeTasks,omitempty"`
	FormContext             string   `json:"formContext,omitempty"`
	FormFieldName           string   `json:"formFieldName,omitempty"`
	FormHelp                string   `json:"formHelp,omitempty"`
	HelpMsg                 string   `json:"helpMsg,omitempty"`
	HelpPrimary             string   `json:"helpPrimary,omitempty"`
	HelpSecondary           string   `json:"helpSecondary,omitempty"`
	InactiveTasks           string   `json:"inactiveTasks,omitempty"`
	InitialHelpMsg          string   `json:"initialHelpMsg,omitempty"`
	ListItemDesc            string   `json:"listItemDesc,omitempty"`
	ListItemTitle           string   `json:"listItemTitle,omitempty"`
	RecordsBorder           string   `json:"recordsBorder,omitempty"`
	RecordsDateRange        string   `json:"recordsDateRange,omitempty"`
	RecordsFooter           string   `json:"recordsFooter,omitempty"`
	RecordsHeader           string   `json:"recordsHeader,omitempty"`
	RecordsHelp             string   `json:"recordsHelp,omitempty"`
	TaskEntry               string   `json:"taskEntry,omitempty"`
	TaskLogDetailsViewTitle string   `json:"taskLogDetails,omitempty"`
	TaskLogEntry            string   `json:"taskLogEntry,omitempty"`
	TaskLogFormError        string   `json:"taskLogFormError,omitempty"`
	TaskLogFormInfo         string   `json:"taskLogFormInfo,omitempty"`
	TaskLogFormWarn         string   `json:"taskLogFormWarn,omitempty"`
	TaskLogList             string   `json:"taskLogList,omitempty"`
	Tasks                   []string `json:"tasks,omitempty"`
	TitleForeground         string   `json:"titleForeground,omitempty"`
	ToolName                string   `json:"toolName,omitempty"`
	Tracking                string   `json:"tracking,omitempty"`
}

func Get(themeName string, themesDir string) (Theme, error) {
	var zero Theme
	themeName = strings.TrimSpace(themeName)

	if len(themeName) == 0 {
		return zero, errEmptyThemeNameProvided
	}

	if themeName == defaultThemeName {
		return Default(), nil
	}

	if customThemeName, ok := strings.CutPrefix(themeName, customThemePrefix); ok {
		if len(customThemeName) == 0 {
			return zero, errEmptyThemeNameProvided
		}

		themeFilePath := path.Join(themesDir, fmt.Sprintf("%s.json", customThemeName))
		themeBytes, err := os.ReadFile(themeFilePath)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				return zero, fmt.Errorf("%w: %q", ErrCustomThemeDoesntExist, customThemeName)
			}
			return zero, fmt.Errorf("%w %q: %s", errCouldntReadCustomThemeFile, themeFilePath, err.Error())
		}

		theme, err := loadCustom(themeBytes)
		if err != nil {
			return zero, fmt.Errorf("%w from file %q: %w", errCouldntLoadCustomTheme, themeFilePath, err)
		}

		return theme, nil
	}

	builtInTheme, err := getBuiltIn(themeName)
	if err != nil {
		return zero, err
	}

	return builtInTheme, nil
}

func Default() Theme {
	return getBuiltInTheme(paletteGruvboxDark())
}

func BuiltIn() []string {
	return []string{
		themeNameCatppuccinMocha,
		themeNameGruvboxDark,
		themeNameMonokaiClassic,
		themeNameTokyonight,
	}
}

func loadCustom(themeJSON []byte) (Theme, error) {
	theme := Default()
	err := json.Unmarshal(themeJSON, &theme)
	var syntaxError *json.SyntaxError

	if err != nil {
		if errors.As(err, &syntaxError) {
			return theme, fmt.Errorf("%w: %w", errThemeFileIsInvalidJSON, err)
		}
		return theme, fmt.Errorf("%w: %s", ErrThemeFileHasInvalidSchema, err.Error())
	}

	invalidColors := getInvalidColors(theme)
	if len(invalidColors) > 0 {
		return theme, fmt.Errorf("%w: %q", ErrThemeColorsAreInvalid, invalidColors)
	}

	return theme, err
}

func getBuiltIn(theme string) (Theme, error) {
	var palette builtInThemePalette
	switch theme {
	case themeNameCatppuccinMocha:
		palette = paletteCatppuccinMocha()
	case themeNameGruvboxDark:
		palette = paletteGruvboxDark()
	case themeNameMonokaiClassic:
		palette = paletteMonokai()
	case themeNameTokyonight:
		palette = paletteTokyonight()
	default:
		return Theme{}, fmt.Errorf("%w: %q", ErrBuiltInThemeDoesntExist, theme)
	}

	return getBuiltInTheme(palette), nil
}

func getInvalidColors(theme Theme) []string {
	var invalidColors []string

	if !isValidColor(theme.ActiveTask) {
		invalidColors = append(invalidColors, "activeTask")
	}
	if !isValidColor(theme.ActiveTaskBeginTime) {
		invalidColors = append(invalidColors, "activeTaskBeginTime")
	}
	if !isValidColor(theme.ActiveTasks) {
		invalidColors = append(invalidColors, "activeTasks")
	}
	if !isValidColor(theme.FormContext) {
		invalidColors = append(invalidColors, "formContext")
	}
	if !isValidColor(theme.FormFieldName) {
		invalidColors = append(invalidColors, "formFieldName")
	}
	if !isValidColor(theme.FormHelp) {
		invalidColors = append(invalidColors, "formHelp")
	}
	if !isValidColor(theme.HelpMsg) {
		invalidColors = append(invalidColors, "helpMsg")
	}
	if !isValidColor(theme.HelpPrimary) {
		invalidColors = append(invalidColors, "helpPrimary")
	}
	if !isValidColor(theme.HelpSecondary) {
		invalidColors = append(invalidColors, "helpSecondary")
	}
	if !isValidColor(theme.InactiveTasks) {
		invalidColors = append(invalidColors, "inactiveTasks")
	}
	if !isValidColor(theme.InitialHelpMsg) {
		invalidColors = append(invalidColors, "initialHelpMsg")
	}
	if !isValidColor(theme.ListItemDesc) {
		invalidColors = append(invalidColors, "listItemDesc")
	}
	if !isValidColor(theme.ListItemTitle) {
		invalidColors = append(invalidColors, "listItemTitle")
	}
	if !isValidColor(theme.RecordsBorder) {
		invalidColors = append(invalidColors, "recordsBorder")
	}
	if !isValidColor(theme.RecordsDateRange) {
		invalidColors = append(invalidColors, "recordsDateRange")
	}
	if !isValidColor(theme.RecordsFooter) {
		invalidColors = append(invalidColors, "recordsFooter")
	}
	if !isValidColor(theme.RecordsHeader) {
		invalidColors = append(invalidColors, "recordsHeader")
	}
	if !isValidColor(theme.RecordsHelp) {
		invalidColors = append(invalidColors, "recordsHelp")
	}
	if !isValidColor(theme.TaskLogDetailsViewTitle) {
		invalidColors = append(invalidColors, "taskLogDetails")
	}
	if !isValidColor(theme.TaskEntry) {
		invalidColors = append(invalidColors, "taskEntry")
	}
	if !isValidColor(theme.TaskLogEntry) {
		invalidColors = append(invalidColors, "taskLogEntry")
	}
	if !isValidColor(theme.TaskLogList) {
		invalidColors = append(invalidColors, "taskLogList")
	}
	if !isValidColor(theme.TaskLogFormInfo) {
		invalidColors = append(invalidColors, "taskLogFormInfo")
	}
	if !isValidColor(theme.TaskLogFormWarn) {
		invalidColors = append(invalidColors, "taskLogFormWarn")
	}
	if !isValidColor(theme.TaskLogFormError) {
		invalidColors = append(invalidColors, "taskLogFormError")
	}
	for i, color := range theme.Tasks {
		if !isValidColor(color) {
			invalidColors = append(invalidColors, fmt.Sprintf("tasks[%d]", i+1))
		}
	}
	if !isValidColor(theme.TitleForeground) {
		invalidColors = append(invalidColors, "titleForeground")
	}
	if !isValidColor(theme.ToolName) {
		invalidColors = append(invalidColors, "toolName")
	}
	if !isValidColor(theme.Tracking) {
		invalidColors = append(invalidColors, "tracking")
	}

	return invalidColors
}

func isValidColor(s string) bool {
	if len(s) == 0 {
		return false
	}

	if strings.HasPrefix(s, "#") {
		return hexCodeRegex.MatchString(s)
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		return false
	}

	if i < 0 || i > 255 {
		return false
	}

	return true
}
