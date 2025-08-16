package ui

import (
	"hash/fnv"

	"github.com/charmbracelet/lipgloss"
)

const (
	fallbackTaskColor = "#ada7ff"
)

type Style struct {
	activeTaskBeginTime  lipgloss.Style
	activeTaskSummaryMsg lipgloss.Style
	empty                lipgloss.Style
	formContext          lipgloss.Style
	formFieldName        lipgloss.Style
	formHelp             lipgloss.Style
	helpMsg              lipgloss.Style
	helpPrimary          lipgloss.Style
	helpSecondary        lipgloss.Style
	helpTitle            lipgloss.Style
	initialHelpMsg       lipgloss.Style
	list                 lipgloss.Style
	listItemDescColor    lipgloss.Color
	listItemTitleColor   lipgloss.Color
	recordsBorder        lipgloss.Style
	recordsDateRange     lipgloss.Style
	recordsFooter        lipgloss.Style
	recordsHeader        lipgloss.Style
	recordsHelp          lipgloss.Style
	taskEntryHeading     lipgloss.Style
	taskLogDetails       lipgloss.Style
	taskLogEntryHeading  lipgloss.Style
	theme                Theme
	titleForegroundColor lipgloss.Color
	tlFormOkStyle        lipgloss.Style
	tlFormWarnStyle      lipgloss.Style
	tlFormErrStyle       lipgloss.Style
	toolName             lipgloss.Style
	tracking             lipgloss.Style
	viewPort             lipgloss.Style
}

func NewStyle(theme Theme) Style {
	base := lipgloss.NewStyle().
		PaddingLeft(1).
		PaddingRight(1).
		Foreground(lipgloss.Color(theme.TitleForeground))

	baseList := lipgloss.NewStyle().PaddingTop(1).PaddingRight(2).PaddingBottom(1)

	baseHeading := lipgloss.NewStyle().
		Bold(true).
		PaddingLeft(1).
		PaddingRight(1).
		Foreground(lipgloss.Color(theme.TitleForeground))

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
		Background(lipgloss.Color(theme.HelpPrimary)).
		Align(lipgloss.Left)

	return Style{
		activeTaskBeginTime:  lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color(theme.ActiveTaskBeginTime)),
		activeTaskSummaryMsg: tracking.PaddingLeft(1).Foreground(lipgloss.Color(theme.ActiveTask)),
		empty:                lipgloss.NewStyle(),
		formContext:          lipgloss.NewStyle().Foreground(lipgloss.Color(theme.FormContext)),
		formFieldName:        lipgloss.NewStyle().Foreground(lipgloss.Color(theme.FormFieldName)),
		formHelp:             lipgloss.NewStyle().Foreground(lipgloss.Color(theme.FormHelp)),
		helpMsg:              helpMsg,
		helpPrimary:          lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(theme.HelpPrimary)),
		helpSecondary:        lipgloss.NewStyle().Foreground(lipgloss.Color(theme.HelpSecondary)),
		helpTitle:            helpTitle,
		initialHelpMsg:       helpMsg.Foreground(lipgloss.Color(theme.InitialHelpMsg)),
		list:                 baseList,
		listItemDescColor:    lipgloss.Color(theme.ListItemDesc),
		listItemTitleColor:   lipgloss.Color(theme.ListItemTitle),
		recordsBorder:        lipgloss.NewStyle().Foreground(lipgloss.Color(theme.RecordsBorder)),
		recordsDateRange:     lipgloss.NewStyle().Foreground(lipgloss.Color(theme.RecordsDateRange)),
		recordsFooter:        lipgloss.NewStyle().Foreground(lipgloss.Color(theme.RecordsFooter)),
		recordsHeader:        lipgloss.NewStyle().Foreground(lipgloss.Color(theme.RecordsHeader)),
		recordsHelp:          lipgloss.NewStyle().Foreground(lipgloss.Color(theme.RecordsHelp)),
		taskEntryHeading:     baseHeading.Background(lipgloss.Color(theme.TaskEntry)),
		taskLogDetails:       helpTitle.Background(lipgloss.Color(theme.TaskLogDetailsViewTitle)),
		taskLogEntryHeading:  baseHeading.Background(lipgloss.Color(theme.TaskLogEntry)),
		theme:                theme,
		titleForegroundColor: lipgloss.Color(theme.TitleForeground),
		tlFormOkStyle:        lipgloss.NewStyle().Foreground(lipgloss.Color(theme.TaskLogFormInfo)),
		tlFormWarnStyle:      lipgloss.NewStyle().Foreground(lipgloss.Color(theme.TaskLogFormWarn)),
		tlFormErrStyle:       lipgloss.NewStyle().Foreground(lipgloss.Color(theme.TaskLogFormError)),
		toolName:             base.Align(lipgloss.Center).Bold(true).Background(lipgloss.Color(theme.ToolName)),
		tracking:             tracking,
		viewPort:             lipgloss.NewStyle().PaddingTop(1).PaddingLeft(2).PaddingRight(2).PaddingBottom(1),
	}
}

func (s *Style) getDynamicStyle(str string) lipgloss.Style {
	h := fnv.New32()
	_, err := h.Write([]byte(str))
	if err != nil {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(fallbackTaskColor))
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
