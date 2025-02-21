package ui

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
)

var (
	errThemeFileIsInvalidJSON    = errors.New("theme file is not valid JSON")
	ErrThemeFileHasInvalidSchema = errors.New("theme file's schema is incorrect")
	errThemeDataIsInvalid        = errors.New("theme data is invalid")
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

	invalidColors := getInvalidHEXColors(theme)
	if len(invalidColors) > 0 {
		return theme, fmt.Errorf("%w: invalid HEX colors: %v", errThemeDataIsInvalid, invalidColors)
	}

	return theme, err
}

func getInvalidHEXColors(theme Theme) []string {
	var invalidColors []string

	if !hexCodeRegex.MatchString(theme.ActiveTask) {
		invalidColors = append(invalidColors, "activeTask")
	}
	if !hexCodeRegex.MatchString(theme.ActiveTaskBeginTime) {
		invalidColors = append(invalidColors, "activeTaskBeginTime")
	}
	if !hexCodeRegex.MatchString(theme.ActiveTasks) {
		invalidColors = append(invalidColors, "activeTasks")
	}
	if !hexCodeRegex.MatchString(theme.FormContext) {
		invalidColors = append(invalidColors, "formContext")
	}
	if !hexCodeRegex.MatchString(theme.FormFieldName) {
		invalidColors = append(invalidColors, "formFieldName")
	}
	if !hexCodeRegex.MatchString(theme.FormHelp) {
		invalidColors = append(invalidColors, "formHelp")
	}
	if !hexCodeRegex.MatchString(theme.HelpMsg) {
		invalidColors = append(invalidColors, "helpMsg")
	}
	if !hexCodeRegex.MatchString(theme.HelpPrimary) {
		invalidColors = append(invalidColors, "helpPrimary")
	}
	if !hexCodeRegex.MatchString(theme.HelpSecondary) {
		invalidColors = append(invalidColors, "helpSecondary")
	}
	if !hexCodeRegex.MatchString(theme.InactiveTasks) {
		invalidColors = append(invalidColors, "inactiveTasks")
	}
	if !hexCodeRegex.MatchString(theme.InitialHelpMsg) {
		invalidColors = append(invalidColors, "initialHelpMsg")
	}
	if !hexCodeRegex.MatchString(theme.ListItemDesc) {
		invalidColors = append(invalidColors, "ListItemDesc")
	}
	if !hexCodeRegex.MatchString(theme.ListItemTitle) {
		invalidColors = append(invalidColors, "ListItemTitle")
	}
	if !hexCodeRegex.MatchString(theme.RecordsBorder) {
		invalidColors = append(invalidColors, "recordsBorder")
	}
	if !hexCodeRegex.MatchString(theme.RecordsDateRange) {
		invalidColors = append(invalidColors, "recordsDateRange")
	}
	if !hexCodeRegex.MatchString(theme.RecordsFooter) {
		invalidColors = append(invalidColors, "recordsFooter")
	}
	if !hexCodeRegex.MatchString(theme.RecordsHeader) {
		invalidColors = append(invalidColors, "recordsHeader")
	}
	if !hexCodeRegex.MatchString(theme.RecordsHelp) {
		invalidColors = append(invalidColors, "recordsHelp")
	}
	if !hexCodeRegex.MatchString(theme.TaskEntry) {
		invalidColors = append(invalidColors, "taskEntry")
	}
	if !hexCodeRegex.MatchString(theme.TaskLogDetailsViewTitle) {
		invalidColors = append(invalidColors, "taskLogDetails")
	}
	if !hexCodeRegex.MatchString(theme.TaskLogEntry) {
		invalidColors = append(invalidColors, "taskLogEntry")
	}
	if !hexCodeRegex.MatchString(theme.TaskLogList) {
		invalidColors = append(invalidColors, "taskLogList")
	}
	if !hexCodeRegex.MatchString(theme.TitleForeground) {
		invalidColors = append(invalidColors, "titleForeground")
	}
	if !hexCodeRegex.MatchString(theme.ToolName) {
		invalidColors = append(invalidColors, "toolName")
	}
	if !hexCodeRegex.MatchString(theme.Tracking) {
		invalidColors = append(invalidColors, "tracking")
	}

	for i, color := range theme.Tasks {
		if !hexCodeRegex.MatchString(color) {
			invalidColors = append(invalidColors, fmt.Sprintf("tasks[%d]", i+1))
		}
	}

	return invalidColors
}
