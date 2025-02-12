package ui

import (
	"encoding/json"
	"fmt"
	"os"
)

type Theme struct {
	DefaultBackground   string   `json:"defaultBackground"`
	ActiveTaskList      string   `json:"activeTaskList"`
	InactiveTaskList    string   `json:"inactiveTaskList"`
	TaskEntry           string   `json:"taskEntry"`
	TaskLogEntry        string   `json:"taskLogEntry"`
	TaskLogList         string   `json:"taskLogList"`
	Tracking            string   `json:"tracking"`
	ActiveTask          string   `json:"activeTask"`
	ActiveTaskBeginTime string   `json:"activeTaskBeginTime"`
	FormFieldName       string   `json:"formFieldName"`
	FormHelp            string   `json:"formHelp"`
	FormContext         string   `json:"formContext"`
	ToolName            string   `json:"toolName"`
	RecordsHeader       string   `json:"recordsHeader"`
	RecordsFooter       string   `json:"recordsFooter"`
	RecordsBorder       string   `json:"recordsBorder"`
	InitialHelpMsg      string   `json:"initialHelpMsg"`
	RecordsDateRange    string   `json:"recordsDateRange"`
	RecordsHelp         string   `json:"recordsHelp"`
	TLDetailsViewTitle  string   `json:"tLDetailsViewTitle"`
	HelpMsg             string   `json:"helpMsg"`
	HelpViewTitle       string   `json:"helpViewTitle"`
	HelpHeader          string   `json:"helpHeader"`
	HelpSection         string   `json:"helpSection"`
	FallbackTask        string   `json:"fallbackTask"`
	Warning             string   `json:"warning"`
	Tasks               []string `json:"tasks"`
}

func DefaultTheme() Theme {
	return Theme{
		DefaultBackground:   "#282828",
		ActiveTaskList:      "#fe8019",
		InactiveTaskList:    "#928374",
		TaskEntry:           "#8ec07c",
		TaskLogEntry:        "#fabd2f",
		TaskLogList:         "#b8bb26",
		Tracking:            "#fabd2f",
		ActiveTask:          "#8ec07c",
		ActiveTaskBeginTime: "#d3869b",
		FormFieldName:       "#8ec07c",
		FormHelp:            "#928374",
		FormContext:         "#fabd2f",
		ToolName:            "#fe8019",
		RecordsHeader:       "#d85d5d",
		RecordsFooter:       "#ef8f62",
		RecordsBorder:       "#665c54",
		InitialHelpMsg:      "#a58390",
		RecordsDateRange:    "#fabd2f",
		RecordsHelp:         "#928374",
		TLDetailsViewTitle:  "#d3869b",
		HelpMsg:             "#83a598",
		HelpViewTitle:       "#83a598",
		HelpHeader:          "#83a598",
		HelpSection:         "#bdae93",
		FallbackTask:        "#ada7ff",
		Warning:             "#fb4934",
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
	}
}

func LoadTheme(path string) (Theme, error) {
	var (
		err   error
		theme = DefaultTheme()
	)
	themeFile, err := os.ReadFile(path)
	if err != nil {
		return theme, fmt.Errorf("failed to read theme file %q: %w", path, err)
	}
	if err = json.Unmarshal(themeFile, &theme); err != nil {
		return theme, fmt.Errorf("failed to parse theme file %q: %w", path, err)
	}
	return theme, err
}
