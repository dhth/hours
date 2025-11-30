package theme

type builtInThemePalette struct {
	primary    string
	secondary  string
	tertiary   string
	quaternary string
	foreground string
	text       string
	subtext    string
	muted      string
	help       string
	info       string
	error      string
	warn       string
	tasks []string
}

func getBuiltInTheme(palette builtInThemePalette) Theme {
	taskColors := []string{
		palette.tertiary,
		palette.quaternary,
	}
	taskColors = append(taskColors, palette.tasks...)

	return Theme{
		ActiveTask:              palette.primary,
		ActiveTaskBeginTime:     palette.tertiary,
		ActiveTasks:             palette.primary,
		FormContext:             palette.help,
		FormFieldName:           palette.text,
		FormHelp:                palette.subtext,
		HelpMsg:                 palette.help,
		HelpPrimary:             palette.help,
		HelpSecondary:           palette.subtext,
		InactiveTasks:           palette.quaternary,
		InitialHelpMsg:          palette.help,
		ListItemDesc:            palette.subtext,
		ListItemTitle:           palette.text,
		RecordsBorder:           palette.muted,
		RecordsDateRange:        palette.help,
		RecordsFooter:           palette.secondary,
		RecordsHeader:           palette.primary,
		RecordsHelp:             palette.subtext,
		TaskEntry:               palette.primary,
		TaskLogDetailsViewTitle: palette.tertiary,
		TaskLogEntry:            palette.secondary,
		TaskLogFormError:        palette.error,
		TaskLogFormInfo:         palette.info,
		TaskLogFormWarn:         palette.warn,
		TaskLogList:             palette.tertiary,
		Tasks:                   taskColors,
		TitleForeground:         palette.foreground,
		ToolName:                palette.primary,
		Tracking:                palette.secondary,
	}
}
