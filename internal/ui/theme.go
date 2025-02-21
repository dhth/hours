package ui

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	errThemeFileIsInvalidJSON    = errors.New("theme file is not valid JSON")
	ErrThemeFileHasInvalidSchema = errors.New("theme file's schema is incorrect")
	ErrThemeColorsAreInvalid     = errors.New("invalid colors provided")
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
	TaskLogDetailsViewTitle string   `json:"taskLogDetails,omitempty"`
	TaskEntry               string   `json:"taskEntry,omitempty"`
	TaskLogEntry            string   `json:"taskLogEntry,omitempty"`
	TaskLogList             string   `json:"taskLogList,omitempty"`
	Tasks                   []string `json:"tasks,omitempty"`
	TitleForeground         string   `json:"titleForeground,omitempty"`
	ToolName                string   `json:"toolName,omitempty"`
	Tracking                string   `json:"tracking,omitempty"`
}

func DefaultTheme() Theme {
	return Theme{
		ActiveTask:          "#8ec07c",
		ActiveTaskBeginTime: "#d3869b",
		ActiveTasks:         "#fe8019",
		FormContext:         "#fabd2f",
		FormFieldName:       "#8ec07c",
		FormHelp:            "#928374",
		HelpMsg:             "#83a598",
		HelpPrimary:         "#83a598",
		HelpSecondary:       "#bdae93",
		InactiveTasks:       "#928374",
		InitialHelpMsg:      "#a58390",
		ListItemDesc:        "#777777",
		ListItemTitle:       "#dddddd",
		RecordsBorder:       "#665c54",
		RecordsDateRange:    "#fabd2f",
		RecordsFooter:       "#ef8f62",
		RecordsHeader:       "#d85d5d",
		RecordsHelp:         "#928374",
		Tasks: []string{
			"#d3869b",
			"#b5e48c",
			"#90e0ef",
			"#ca7df9",
			"#ada7ff",
			"#bbd0ff",
			"#48cae4",
			"#8187dc",
			"#ffb4a2",
			"#b8bb26",
			"#ffc6ff",
			"#4895ef",
			"#83a598",
			"#fabd2f",
		},
		TaskEntry:               "#8ec07c",
		TaskLogDetailsViewTitle: "#d3869b",
		TaskLogEntry:            "#fabd2f",
		TaskLogList:             "#b8bb26",
		TitleForeground:         "#282828",
		ToolName:                "#fe8019",
		Tracking:                "#fabd2f",
	}
}

func LoadTheme(themeJSON []byte) (Theme, error) {
	theme := DefaultTheme()
	err := json.Unmarshal(themeJSON, &theme)
	var syntaxError *json.SyntaxError

	if err != nil {
		if errors.As(err, &syntaxError) {
			return theme, fmt.Errorf("%w: %w", errThemeFileIsInvalidJSON, err)
		}
		return theme, ErrThemeFileHasInvalidSchema
	}

	invalidColors := getInvalidColors(theme)
	if len(invalidColors) > 0 {
		return theme, fmt.Errorf("%w: %q", ErrThemeColorsAreInvalid, invalidColors)
	}

	return theme, err
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
		invalidColors = append(invalidColors, "ListItemDesc")
	}
	if !isValidColor(theme.ListItemTitle) {
		invalidColors = append(invalidColors, "ListItemTitle")
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
	if !isValidColor(theme.TaskEntry) {
		invalidColors = append(invalidColors, "taskEntry")
	}
	if !isValidColor(theme.TaskLogDetailsViewTitle) {
		invalidColors = append(invalidColors, "taskLogDetails")
	}
	if !isValidColor(theme.TaskLogEntry) {
		invalidColors = append(invalidColors, "taskLogEntry")
	}
	if !isValidColor(theme.TaskLogList) {
		invalidColors = append(invalidColors, "taskLogList")
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

	for i, color := range theme.Tasks {
		if !isValidColor(color) {
			invalidColors = append(invalidColors, fmt.Sprintf("tasks[%d]", i+1))
		}
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
