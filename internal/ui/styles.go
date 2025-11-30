package ui

import (
	"hash/fnv"

	"github.com/charmbracelet/lipgloss"
	"github.com/dhth/hours/internal/ui/theme"
	"github.com/olekukonko/tablewriter/tw"
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
	theme                theme.Theme
	titleForegroundColor lipgloss.Color
	tlFormOkStyle        lipgloss.Style
	tlFormWarnStyle      lipgloss.Style
	tlFormErrStyle       lipgloss.Style
	toolName             lipgloss.Style
	tracking             lipgloss.Style
	viewPort             lipgloss.Style
}

func NewStyle(theme theme.Theme) Style {
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
	plain       bool
}

func (rs reportStyles) symbols(borderStyle tw.BorderStyle) tw.Symbols {
	base := tw.NewSymbols(borderStyle)
	if rs.plain {
		return base
	}

	return styledTableSymbols{base: base, style: rs.borderStyle}
}

type styledTableSymbols struct {
	base  tw.Symbols
	style lipgloss.Style
}

func (s styledTableSymbols) Name() string        { return s.base.Name() }
func (s styledTableSymbols) Center() string      { return s.style.Render(s.base.Center()) }
func (s styledTableSymbols) Row() string         { return s.style.Render(s.base.Row()) }
func (s styledTableSymbols) Column() string      { return s.style.Render(s.base.Column()) }
func (s styledTableSymbols) TopLeft() string     { return s.style.Render(s.base.TopLeft()) }
func (s styledTableSymbols) TopMid() string      { return s.style.Render(s.base.TopMid()) }
func (s styledTableSymbols) TopRight() string    { return s.style.Render(s.base.TopRight()) }
func (s styledTableSymbols) MidLeft() string     { return s.style.Render(s.base.MidLeft()) }
func (s styledTableSymbols) MidRight() string    { return s.style.Render(s.base.MidRight()) }
func (s styledTableSymbols) BottomLeft() string  { return s.style.Render(s.base.BottomLeft()) }
func (s styledTableSymbols) BottomMid() string   { return s.style.Render(s.base.BottomMid()) }
func (s styledTableSymbols) BottomRight() string { return s.style.Render(s.base.BottomRight()) }
func (s styledTableSymbols) HeaderLeft() string  { return s.style.Render(s.base.HeaderLeft()) }
func (s styledTableSymbols) HeaderMid() string   { return s.style.Render(s.base.HeaderMid()) }
func (s styledTableSymbols) HeaderRight() string { return s.style.Render(s.base.HeaderRight()) }

func (s *Style) getReportStyles(plain bool) reportStyles {
	if plain {
		return reportStyles{
			headerStyle: s.empty,
			footerStyle: s.empty,
			borderStyle: s.empty,
			plain:       true,
		}
	}

	return reportStyles{
		headerStyle: s.recordsHeader,
		footerStyle: s.recordsFooter,
		borderStyle: s.recordsBorder,
		plain:       false,
	}
}
