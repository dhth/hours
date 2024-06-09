package ui

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	defaultBackgroundColor = "#282828"
	taskListColor          = "#fe8019"
	taskLogListColor       = "#b8bb26"
	trackingColor          = "#fabd2f"
	activeTaskolor         = "#8ec07c"
	formFieldNameColor     = "#8ec07c"
	formContextColor       = "#fabd2f"
	toolNameColor          = "#fe8019"
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

	baseListStyle = lipgloss.NewStyle().PaddingTop(1).PaddingRight(2).PaddingLeft(1).PaddingBottom(1)
	viewPortStyle = lipgloss.NewStyle().
			PaddingTop(1).
			PaddingRight(2).
			PaddingLeft(1).
			PaddingBottom(1)

	stackListStyle = baseListStyle

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

	helpTitleStyle = baseStyle.
			Bold(true).
			Background(lipgloss.Color(helpViewTitleColor)).
			Align(lipgloss.Left)

	helpHeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(helpHeaderColor))

	helpSectionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(helpSectionColor))
)
