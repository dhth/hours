package ui

import (
	"hash/fnv"

	"github.com/charmbracelet/lipgloss"
)

type Style struct {
	theme                Theme
	helpMsg              lipgloss.Style
	initialHelpMsg       lipgloss.Style
	viewPort             lipgloss.Style
	list                 lipgloss.Style
	toolName             lipgloss.Style
	taskEntryHeading     lipgloss.Style
	taskLogEntryHeading  lipgloss.Style
	formContext          lipgloss.Style
	formFieldName        lipgloss.Style
	formHelp             lipgloss.Style
	tracking             lipgloss.Style
	activeTaskSummaryMsg lipgloss.Style
	activeTaskBeginTime  lipgloss.Style
	recordsHeader        lipgloss.Style
	recordsFooter        lipgloss.Style
	recordsBorder        lipgloss.Style
	recordsDateRange     lipgloss.Style
	recordsHelp          lipgloss.Style
	tLDetailsViewTitle   lipgloss.Style
	helpTitle            lipgloss.Style
	helpHeader           lipgloss.Style
	helpSection          lipgloss.Style
	empty                lipgloss.Style
	Warning              lipgloss.Style
}

func NewStyle(theme Theme) Style {
	base := lipgloss.NewStyle().
		PaddingLeft(1).
		PaddingRight(1).
		Foreground(lipgloss.Color(theme.DefaultBackground))

	baseList := lipgloss.NewStyle().PaddingTop(1).PaddingRight(2).PaddingBottom(1)

	baseHeading := lipgloss.NewStyle().
		Bold(true).
		PaddingLeft(1).
		PaddingRight(1).
		Foreground(lipgloss.Color(theme.DefaultBackground))

	helpMsg := lipgloss.NewStyle().
		PaddingLeft(1).
		Bold(true).
		Foreground(lipgloss.Color(theme.HelpMsg))

	tracking := lipgloss.NewStyle().
		PaddingLeft(2).
		Bold(true).
		Foreground(lipgloss.Color(theme.Tracking))

	helpTitle := base.
		Bold(true).
		Background(lipgloss.Color(theme.HelpViewTitle)).
		Align(lipgloss.Left)

	return Style{
		theme: theme,

		helpMsg: helpMsg,

		initialHelpMsg: helpMsg.
			Foreground(lipgloss.Color(theme.InitialHelpMsg)),

		viewPort: lipgloss.NewStyle().
			PaddingTop(1).
			PaddingLeft(2).
			PaddingRight(2).
			PaddingBottom(1),

		list: baseList,

		toolName: base.
			Align(lipgloss.Center).
			Bold(true).
			Background(lipgloss.Color(theme.ToolName)),

		taskEntryHeading: baseHeading.
			Background(lipgloss.Color(theme.TaskEntry)),

		taskLogEntryHeading: baseHeading.
			Background(lipgloss.Color(theme.TaskLogEntry)),

		formContext: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.FormContext)),

		formFieldName: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.FormFieldName)),

		formHelp: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.FormHelp)),

		tracking: tracking,

		activeTaskSummaryMsg: tracking.
			PaddingLeft(1).
			Foreground(lipgloss.Color(theme.ActiveTask)),

		activeTaskBeginTime: lipgloss.NewStyle().
			PaddingLeft(1).
			Foreground(lipgloss.Color(theme.ActiveTaskBeginTime)),

		recordsHeader: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.RecordsHeader)),

		recordsFooter: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.RecordsFooter)),

		recordsBorder: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.RecordsBorder)),

		recordsDateRange: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.RecordsDateRange)),

		recordsHelp: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.RecordsHelp)),

		tLDetailsViewTitle: helpTitle.
			Background(lipgloss.Color(theme.TLDetailsViewTitle)),

		helpTitle: helpTitle,

		helpHeader: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(theme.HelpHeader)),

		helpSection: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.HelpSection)),

		empty: lipgloss.NewStyle(),

		Warning: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Warning)),
	}
}

func (s *Style) getDynamicStyle(str string) lipgloss.Style {
	h := fnv.New32()
	_, err := h.Write([]byte(str))
	if err != nil {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(s.theme.FallbackTask))
	}

	hash := h.Sum32()

	color := s.theme.Tasks[hash%uint32(len(s.theme.Tasks))]
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(color))
}

type reportStyles struct {
	headerStyle lipgloss.Style
	footerStyle lipgloss.Style
	borderStyle lipgloss.Style
}

func (s *Style) getReportStyles(plain bool) reportStyles {
	if plain {
		return reportStyles{
			s.empty,
			s.empty,
			s.empty,
		}
	}
	return reportStyles{
		s.recordsHeader,
		s.recordsFooter,
		s.recordsBorder,
	}
}
