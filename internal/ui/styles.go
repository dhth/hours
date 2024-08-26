package ui

import (
	"hash/fnv"

	"github.com/charmbracelet/lipgloss"
)

const (
	defaultBackgroundColor   = "#282828"
	activeTaskListColor      = "#fe8019"
	inactiveTaskListColor    = "#928374"
	taskLogEntryColor        = "#fabd2f"
	taskLogListColor         = "#b8bb26"
	trackingColor            = "#fabd2f"
	activeTaskColor          = "#8ec07c"
	activeTaskBeginTimeColor = "#d3869b"
	formFieldNameColor       = "#8ec07c"
	formHelpColor            = "#928374"
	formContextColor         = "#fabd2f"
	toolNameColor            = "#fe8019"
	recordsHeaderColor       = "#d85d5d"
	recordsFooterColor       = "#ef8f62"
	recordsBorderColor       = "#665c54"
	initialHelpMsgColor      = "#a58390"
	recordsDateRangeColor    = "#fabd2f"
	recordsHelpColor         = "#928374"
	helpMsgColor             = "#83a598"
	helpViewTitleColor       = "#83a598"
	helpHeaderColor          = "#83a598"
	helpSectionColor         = "#bdae93"
	warningColor             = "#fb4934"
)

var (
	baseStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1).
			Foreground(lipgloss.Color(defaultBackgroundColor))

	helpMsgStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			Bold(true).
			Foreground(lipgloss.Color(helpMsgColor))

	initialHelpMsgStyle = helpMsgStyle.
				Foreground(lipgloss.Color(initialHelpMsgColor))

	baseListStyle = lipgloss.NewStyle().PaddingTop(1).PaddingRight(2).PaddingBottom(1)

	baseHeadingStyle = lipgloss.NewStyle().
				Bold(true).
				PaddingLeft(1).
				PaddingRight(1).
				Foreground(lipgloss.Color(defaultBackgroundColor))

	viewPortStyle = lipgloss.NewStyle().
			PaddingTop(1).
			PaddingRight(2).
			PaddingBottom(1)

	listStyle = baseListStyle

	toolNameStyle = baseStyle.
			Align(lipgloss.Center).
			Bold(true).
			Background(lipgloss.Color(toolNameColor))

	taskLogEntryHeadingStyle = baseHeadingStyle.
					Background(lipgloss.Color(taskLogEntryColor))

	formContextStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(formContextColor))

	formFieldNameStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(formFieldNameColor))

	formHelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(formHelpColor))

	trackingStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			Bold(true).
			Foreground(lipgloss.Color(trackingColor))

	activeTaskSummaryMsgStyle = trackingStyle.
					PaddingLeft(1).
					Foreground(lipgloss.Color(activeTaskColor))

	activeTaskBeginTimeStyle = lipgloss.NewStyle().
					PaddingLeft(1).
					Foreground(lipgloss.Color(activeTaskBeginTimeColor))

	recordsHeaderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(recordsHeaderColor))

	recordsFooterStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(recordsFooterColor))

	recordsBorderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(recordsBorderColor))

	recordsDateRangeStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(recordsDateRangeColor))

	recordsHelpStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(recordsHelpColor))

	helpTitleStyle = baseStyle.
			Bold(true).
			Background(lipgloss.Color(helpViewTitleColor)).
			Align(lipgloss.Left)

	helpHeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(helpHeaderColor))

	helpSectionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(helpSectionColor))

	taskColors = []string{
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
	}

	getDynamicStyle = func(str string) lipgloss.Style {
		h := fnv.New32()
		h.Write([]byte(str))
		hash := h.Sum32()

		color := taskColors[hash%uint32(len(taskColors))]
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(color))
	}

	emptyStyle = lipgloss.NewStyle()

	WarningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(warningColor))
)
