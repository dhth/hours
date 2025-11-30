package theme

const (
	themeNameGruvboxDark = "gruvbox-dark"
)

func themeGruvboxDark() Theme {
	return Theme{
		ActiveTask:              "#8ec07c",
		ActiveTaskBeginTime:     "#d3869b",
		ActiveTasks:             "#fe8019",
		FormContext:             "#fabd2f",
		FormFieldName:           "#8ec07c",
		FormHelp:                "#928374",
		HelpMsg:                 "#83a598",
		HelpPrimary:             "#83a598",
		HelpSecondary:           "#bdae93",
		InactiveTasks:           "#928374",
		InitialHelpMsg:          "#a58390",
		ListItemDesc:            "#777777",
		ListItemTitle:           "#dddddd",
		RecordsBorder:           "#665c54",
		RecordsDateRange:        "#fabd2f",
		RecordsFooter:           "#ef8f62",
		RecordsHeader:           "#d85d5d",
		RecordsHelp:             "#928374",
		TaskEntry:               "#8ec07c",
		TaskLogDetailsViewTitle: "#d3869b",
		TaskLogEntry:            "#fabd2f",
		TaskLogFormError:        "#fb4934",
		TaskLogFormInfo:         "#d3869b",
		TaskLogFormWarn:         "#fe8019",
		TaskLogList:             "#b8bb26",
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
		TitleForeground: "#282828",
		ToolName:        "#fe8019",
		Tracking:        "#fabd2f",
	}
}
