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
	var errorMsg string
	if m.errorMessage != "" {
		errorMsg = "error: " + Trim(m.errorMessage, 120)
	}
	var activeMsg string
	if m.tasksFetched && m.trackingActive {
		var taskSummaryMsg string
		task, ok := m.activeTaskMap[m.activeTaskId]
		if ok {
			taskSummaryMsg = fmt.Sprintf("(%s)", Trim(task.summary, 50))
		}
		activeMsg = fmt.Sprintf("%s%s",
			trackingStyle.Render("tracking:"),
			activeTaskSummaryMsgStyle.Render(taskSummaryMsg),
		)
	}

	switch m.activeView {
	case activeTaskListView:
		content = taskListStyle.Render(m.activeTasksList.View())
	case taskLogView:
		content = taskListStyle.Render(m.taskLogList.View())
	case inactiveTaskListView:
		content = taskListStyle.Render(m.inactiveTasksList.View())
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
		content = fmt.Sprintf(
			`
    %s

    %s


    %s
`,
			formFieldNameStyle.Render(RightPadTrim("Comment:", 16)),
			m.trackingInputs[entryComment].View(),
			formContextStyle.Render("Press enter to submit"),
		)
		for i := 0; i < m.terminalHeight-20+10; i++ {
			content += "\n"
		}
	case manualTasklogEntryView:
		var formHeadingText string
		switch m.tasklogSaveType {
		case tasklogInsert:
			formHeadingText = "Adding a manual entry. Enter the following details:"
		case tasklogUpdate:
			formHeadingText = "Updating task log entry. Enter the following details:"
		}

		content = fmt.Sprintf(
			`
    %s

    %s

    %s

    %s

    %s

    %s

    %s


    %s
`,
			formContextStyle.Render(formHeadingText),
			formFieldNameStyle.Render("Begin Time  (format: 2006/01/02 15:04)"),
			m.trackingInputs[entryBeginTS].View(),
			formFieldNameStyle.Render("End Time  (format: 2006/01/02 15:04)"),
			m.trackingInputs[entryEndTS].View(),
			formFieldNameStyle.Render(RightPadTrim("Comment:", 16)),
			m.trackingInputs[entryComment].View(),
			formContextStyle.Render("Press enter to submit"),
		)
		for i := 0; i < m.terminalHeight-20; i++ {
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
		helpMsg = " " + helpMsgStyle.Render("Press ? for help")
	}

	footerStr := fmt.Sprintf("%s%s%s%s",
		toolNameStyle.Render("hours"),
		helpMsg,
		activeMsg,
		errorMsg,
	)
	footer = footerStyle.Render(footerStr)

	return lipgloss.JoinVertical(lipgloss.Left,
		content,
		statusBar,
		footer,
	)
}
