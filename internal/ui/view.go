package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	listWidth = 140
)

func (m model) View() string {
	var content string
	var footer string

	var statusBar string
	if m.message != "" {
		statusBar = Trim(m.message, 120)
	}

	var activeMsg string
	if m.tasksFetched && m.trackingActive {
		var taskSummaryMsg, taskStartedSinceMsg string
		task, ok := m.activeTaskMap[m.activeTaskId]
		if ok {
			taskSummaryMsg = Trim(task.summary, 50)
			if m.activeView != askForCommentView {
				taskStartedSinceMsg = fmt.Sprintf("(since %s)", m.activeTLBeginTS.Format(timeOnlyFormat))
			}
		}
		activeMsg = fmt.Sprintf("%s%s%s",
			trackingStyle.Render("tracking:"),
			activeTaskSummaryMsgStyle.Render(taskSummaryMsg),
			activeTaskBeginTimeStyle.Render(taskStartedSinceMsg),
		)
	}

	switch m.activeView {
	case activeTaskListView:
		content = listStyle.Render(m.activeTasksList.View())
	case taskLogView:
		content = listStyle.Render(m.taskLogList.View())
	case inactiveTaskListView:
		content = listStyle.Render(m.inactiveTasksList.View())
	case taskInputView:
		var formTitle string
		switch m.taskMgmtContext {
		case taskCreateCxt:
			formTitle = "Add a task"
		case taskUpdateCxt:
			formTitle = "Update task"
		}
		content = fmt.Sprintf(
			`
    %s

    %s


    %s
`,
			formFieldNameStyle.Render(formTitle),
			m.taskInputs[summaryField].View(),
			formContextStyle.Render("Press enter to submit"),
		)
		for i := 0; i < m.terminalHeight-20+10; i++ {
			content += "\n"
		}
	case askForCommentView:
		formHeadingText := "Saving task log entry. Enter the following details."

		content = fmt.Sprintf(
			`
    %s

    %s

    %s

    %s    %s

    %s

    %s    %s

    %s

    %s


    %s
`,
			formContextStyle.Render(formHeadingText),
			formHelpStyle.Render("Use tab/shift-tab to move between sections; esc to go back."),
			formFieldNameStyle.Render("Begin Time (format: 2006/01/02 15:04)"),
			m.trackingInputs[entryBeginTS].View(),
			formHelpStyle.Render("(k/j/K/J moves time, when correct)"),
			formFieldNameStyle.Render("End Time (format: 2006/01/02 15:04)"),
			m.trackingInputs[entryEndTS].View(),
			formHelpStyle.Render("(k/j/K/J moves time, when correct)"),
			formFieldNameStyle.Render(RightPadTrim("Comment:", 16, true)),
			m.trackingInputs[entryComment].View(),
			formContextStyle.Render("Press enter to submit"),
		)
		for i := 0; i < m.terminalHeight-22; i++ {
			content += "\n"
		}
	case editStartTsView:
		formHeadingText := "Update task log entry"

		content = fmt.Sprintf(
			`
    %s

    %s

    %s    %s


    %s
`,
			formContextStyle.Render(formHeadingText),
			formFieldNameStyle.Render("Begin Time (format: 2006/01/02 15:04)"),
			m.trackingInputs[entryBeginTS].View(),
			formHelpStyle.Render("(k/j/K/J moves time, when correct)"),
			formContextStyle.Render("Press enter to submit"),
		)
		for i := 0; i < m.terminalHeight-12; i++ {
			content += "\n"
		}
	case manualTasklogEntryView:
		var formHeadingText string
		switch m.tasklogSaveType {
		case tasklogInsert:
			formHeadingText = "Adding a manual entry. Enter the following details."
		case tasklogUpdate:
			formHeadingText = "Updating task log entry. Enter the following details."
		}

		content = fmt.Sprintf(
			`
    %s

    %s

    %s

    %s    %s

    %s

    %s    %s

    %s

    %s


    %s
`,
			formContextStyle.Render(formHeadingText),
			formHelpStyle.Render("Use tab/shift-tab to move between sections; esc to go back."),
			formFieldNameStyle.Render("Begin Time (format: 2006/01/02 15:04)"),
			m.trackingInputs[entryBeginTS].View(),
			formHelpStyle.Render("(k/j/K/J moves time, when correct)"),
			formFieldNameStyle.Render("End Time (format: 2006/01/02 15:04)"),
			m.trackingInputs[entryEndTS].View(),
			formHelpStyle.Render("(k/j/K/J moves time, when correct)"),
			formFieldNameStyle.Render(RightPadTrim("Comment:", 16, true)),
			m.trackingInputs[entryComment].View(),
			formContextStyle.Render("Press enter to submit"),
		)
		for i := 0; i < m.terminalHeight-22; i++ {
			content += "\n"
		}
	case helpView:
		if !m.helpVPReady {
			content = "\n  Initializing..."
		} else {
			content = viewPortStyle.Render(fmt.Sprintf("  %s\n\n%s\n", helpTitleStyle.Render("Help"), m.helpVP.View()))
		}
	}

	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#282828")).
		Background(lipgloss.Color("#7c6f64"))

	var helpMsg string
	if m.showHelpIndicator {
		// first time directions
		if m.activeView == activeTaskListView && len(m.activeTasksList.Items()) <= 1 {
			if len(m.activeTasksList.Items()) == 0 {
				helpMsg += " " + initialHelpMsgStyle.Render("Press a to add a task")
			} else if len(m.taskLogList.Items()) == 0 {
				if m.trackingActive {
					helpMsg += " " + initialHelpMsgStyle.Render("Press s to stop tracking time")
				} else {
					helpMsg += " " + initialHelpMsgStyle.Render("Press s to start tracking time")
				}
			}
		}

		helpMsg += " " + helpMsgStyle.Render("Press ? for help")
	}

	footerStr := fmt.Sprintf("%s%s%s",
		toolNameStyle.Render("hours"),
		helpMsg,
		activeMsg,
	)
	footer = footerStyle.Render(footerStr)

	return lipgloss.JoinVertical(lipgloss.Left,
		content,
		statusBar,
		footer,
	)
}

func (m recordsModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("Something went wrong: %s\n", m.err)
	}
	var help string

	var dateRangeStr string
	var dateRange string
	if m.numDays > 1 {
		dateRangeStr = fmt.Sprintf(`
 range:             %s...%s
 `,
			m.start.Format(dateFormat), m.end.AddDate(0, 0, -1).Format(dateFormat))
	} else {
		dateRangeStr = fmt.Sprintf(`
 date:              %s
`,
			m.start.Format(dateFormat))
	}

	helpStr := `
 go backwards:      h or <-
 go forwards:       l or ->
 go to today:       ctrl+t

 press ctrl+c/q to quit
`

	if m.plain {
		help = helpStr
		dateRange = dateRangeStr
	} else {
		help = recordsHelpStyle.Render(helpStr)
		dateRange = recordsDateRangeStyle.Render(dateRangeStr)
	}

	return fmt.Sprintf("%s%s%s", m.report, dateRange, help)
}
