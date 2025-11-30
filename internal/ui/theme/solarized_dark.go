package theme

const (
	themeNameSolarizedDark = "solarized-dark"
)

func themeSolarizedDark() Theme {
	return Theme{
		ActiveTask:              "#859900",
		ActiveTaskBeginTime:     "#dc322f",
		ActiveTasks:             "#268bd2",
		FormContext:             "#b58900",
		FormFieldName:           "#859900",
		FormHelp:                "#586e75",
		HelpMsg:                 "#268bd2",
		HelpPrimary:             "#268bd2",
		HelpSecondary:           "#93a1a1",
		InactiveTasks:           "#586e75",
		InitialHelpMsg:          "#d33682",
		ListItemDesc:            "#586e75",
		ListItemTitle:           "#93a1a1",
		RecordsBorder:           "#586e75",
		RecordsDateRange:        "#93a1a1",
		RecordsFooter:           "#cb4b16",
		RecordsHeader:           "#859900",
		RecordsHelp:             "#586e75",
		TaskEntry:               "#859900",
		TaskLogDetailsViewTitle: "#d33682",
		TaskLogEntry:            "#b58900",
		TaskLogFormError:        "#dc322f",
		TaskLogFormInfo:         "#6c71c4",
		TaskLogFormWarn:         "#cb4b16",
		TaskLogList:             "#268bd2",
		Tasks: []string{
			"#dc322f",
			"#859900",
			"#268bd2",
			"#d33682",
			"#cb4b16",
			"#b58900",
			"#2aa198",
			"#6c71c4",
			"#93a1a1",
			"#dc322f",
			"#859900",
			"#268bd2",
			"#d33682",
			"#2aa198",
		},
		TitleForeground: "#002b36",
		ToolName:        "#268bd2",
		Tracking:        "#b58900",
	}
}
