package ui

import (
	"hash/fnv"

	"github.com/charmbracelet/lipgloss"
)

const (
	defaultBackgroundColor = "#282828"
	activeTaskListColor    = "#fe8019"
	inactiveTaskListColor  = "#928374"
	taskLogListColor       = "#b8bb26"
	trackingColor          = "#fabd2f"
	activeTaskolor         = "#8ec07c"
	formFieldNameColor     = "#8ec07c"
	formContextColor       = "#fabd2f"
	toolNameColor          = "#fe8019"
	reportHeaderColor      = "#d85d5d"
	reportFooterColor      = "#ef8f62"
	reportBorderColor      = "#665c54"
	initialHelpMsgColor    = "#a58390"
	helpMsgColor           = "#83a598"
	helpViewTitleColor     = "#83a598"
	helpHeaderColor        = "#83a598"
	helpSectionColor       = "#fabd2f"
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

	baseListStyle = lipgloss.NewStyle().PaddingTop(1).PaddingRight(2).PaddingLeft(1).PaddingBottom(1)
	viewPortStyle = lipgloss.NewStyle().
			PaddingTop(1).
			PaddingRight(2).
			PaddingLeft(1).
			PaddingBottom(1)

	listStyle = baseListStyle

	toolNameStyle = baseStyle.
			Align(lipgloss.Center).
			Bold(true).
			Background(lipgloss.Color(toolNameColor))

	formContextStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(formContextColor))

	formFieldNameStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(formFieldNameColor))

	trackingStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			Bold(true).
			Foreground(lipgloss.Color(trackingColor))

	activeTaskSummaryMsgStyle = trackingStyle.
					PaddingLeft(1).
					Foreground(lipgloss.Color(activeTaskolor))

	reportHeaderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(reportHeaderColor))

	reportFooterStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(reportFooterColor))

	reportBorderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(reportBorderColor))

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

		color := taskColors[int(hash)%len(taskColors)]
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(color))
	}

	emptyStyle = lipgloss.NewStyle()
)
