package ui

import (
	"encoding/json"
)

type Theme struct {
	DefaultBackground   string   `json:"defaultBackground,omitempty"`
	ActiveTaskList      string   `json:"activeTaskList,omitempty"`
	InactiveTaskList    string   `json:"inactiveTaskList,omitempty"`
	TaskEntry           string   `json:"taskEntry,omitempty"`
	TaskLogEntry        string   `json:"taskLogEntry,omitempty"`
	TaskLogList         string   `json:"taskLogList,omitempty"`
	Tracking            string   `json:"tracking,omitempty"`
	ActiveTask          string   `json:"activeTask,omitempty"`
	ActiveTaskBeginTime string   `json:"activeTaskBeginTime,omitempty"`
	FormFieldName       string   `json:"formFieldName,omitempty"`
	FormHelp            string   `json:"formHelp,omitempty"`
	FormContext         string   `json:"formContext,omitempty"`
	ToolName            string   `json:"toolName,omitempty"`
	RecordsHeader       string   `json:"recordsHeader,omitempty"`
	RecordsFooter       string   `json:"recordsFooter,omitempty"`
	RecordsBorder       string   `json:"recordsBorder,omitempty"`
	InitialHelpMsg      string   `json:"initialHelpMsg,omitempty"`
	RecordsDateRange    string   `json:"recordsDateRange,omitempty"`
	RecordsHelp         string   `json:"recordsHelp,omitempty"`
	TLDetailsViewTitle  string   `json:"tLDetailsViewTitle,omitempty"`
	HelpMsg             string   `json:"helpMsg,omitempty"`
	HelpViewTitle       string   `json:"helpViewTitle,omitempty"`
	HelpHeader          string   `json:"helpHeader,omitempty"`
	HelpSection         string   `json:"helpSection,omitempty"`
	FallbackTask        string   `json:"fallbackTask,omitempty"`
	Warning             string   `json:"warning,omitempty"`
	Tasks               []string `json:"tasks,omitempty"`
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

func LoadTheme(themeJSON []byte) (Theme, error) {
	theme := DefaultTheme()
	err := json.Unmarshal(themeJSON, &theme)
	return theme, err
}
