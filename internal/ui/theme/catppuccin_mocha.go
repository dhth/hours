package theme

const (
	themeNameCatppuccinMocha = "catppuccin-mocha"
)

func themeCatppuccinMocha() Theme {
	return Theme{
		ActiveTask:              "#a6e3a1",
		ActiveTaskBeginTime:     "#f38ba8",
		ActiveTasks:             "#89b4fa",
		FormContext:             "#cba6f7",
		FormFieldName:           "#a6e3a1",
		FormHelp:                "#6c7086",
		HelpMsg:                 "#89b4fa",
		HelpPrimary:             "#89b4fa",
		HelpSecondary:           "#cdd6f4",
		InactiveTasks:           "#6c7086",
		InitialHelpMsg:          "#f38ba8",
		ListItemDesc:            "#6c7086",
		ListItemTitle:           "#cdd6f4",
		RecordsBorder:           "#6c7086",
		RecordsDateRange:        "#cdd6f4",
		RecordsFooter:           "#f38ba8",
		RecordsHeader:           "#a6e3a1",
		RecordsHelp:             "#6c7086",
		TaskEntry:               "#a6e3a1",
		TaskLogDetailsViewTitle: "#f38ba8",
		TaskLogEntry:            "#cba6f7",
		TaskLogFormError:        "#f38ba8",
		TaskLogFormInfo:         "#cba6f7",
		TaskLogFormWarn:         "#fab387",
		TaskLogList:             "#89b4fa",
		Tasks: []string{
			"#74c7ec",
			"#89b4fa",
			"#89dceb",
			"#94e2d5",
			"#a6e3a1",
			"#b4befe",
			"#cba6f7",
			"#cdd6f4",
			"#eba0ac",
			"#f2cdcd",
			"#f38ba8",
			"#f5e0dc",
			"#f9e2af",
			"#fab387",
		},
		TitleForeground: "#1e1e2e",
		ToolName:        "#89b4fa",
		Tracking:        "#cba6f7",
	}
}
