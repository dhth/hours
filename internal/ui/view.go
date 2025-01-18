package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/dhth/hours/internal/utils"
)

const (
	taskLogEntryViewHeading = "Task Log Entry"
	minHeightNeeded         = 32
	minWidthNeeded          = 110
)

var listWidth = 140

func (m Model) View() string {
	var content string
	var footer string

	var statusBar string
	if m.message != "" {
		statusBar = utils.Trim(m.message, 120)
	}

	var activeMsg string
	if m.tasksFetched && m.trackingActive {
		var taskSummaryMsg, taskStartedSinceMsg string
		task, ok := m.taskMap[m.activeTaskID]
		if ok {
			taskSummaryMsg = utils.Trim(task.Summary, 50)
			if m.activeView != finishActiveTLView {
				taskStartedSinceMsg = fmt.Sprintf("(since %s)", m.activeTLBeginTS.Format(timeOnlyFormat))
			}
		}
		activeMsg = fmt.Sprintf("%s%s%s",
			trackingStyle.Render("tracking:"),
			activeTaskSummaryMsgStyle.Render(taskSummaryMsg),
			activeTaskBeginTimeStyle.Render(taskStartedSinceMsg),
		)
	}

	formHelp := "Use tab/shift-tab to move between sections; esc to go back."
	formBeginTimeHelp := "Begin Time* (format: 2006/01/02 15:04)"
	formEndTimeHelp := "End Time* (format: 2006/01/02 15:04)"
	formTimeShiftHelp := "(j/k/J/K/h/l moves time)"

	var formCommentContext string
	if m.tLCommentInput.Length() == 0 {
		formCommentContext = "optional"
	} else {
		formCommentContext = fmt.Sprintf("%d/%d", m.tLCommentInput.Length(), tlCommentLengthLimit)
	}
	formCommentHelp := fmt.Sprintf("Comment (%s)", formCommentContext)
	var formSubmitHelp string

	switch m.activeView {
	case taskInputView, editActiveTLView, finishActiveTLView, manualTasklogEntryView:
		if m.trackingFocussedField == entryComment {
			formSubmitHelp = "Press <ctrl+s> to submit"
		} else {
			formSubmitHelp = "Press <ctrl+s>/<enter> to submit"
		}
	}

	switch m.activeView {
	case taskListView:
		content = listStyle.Render(m.activeTasksList.View())
	case taskLogView:
		content = listStyle.Render(m.taskLogList.View())
	case taskLogDetailsView:
		if !m.helpVPReady {
			content = "\n  Initializing..."
		} else {
			content = viewPortStyle.Render(fmt.Sprintf("%s\n\n%s",
				tLDetailsViewTitleStyle.Render("Task Log Details"), m.tLDetailsVP.View()))
		}
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
			formHelpStyle.Render(formSubmitHelp),
		)
		for i := 0; i < m.terminalHeight-20+10; i++ {
			content += "\n"
		}
	case finishActiveTLView:
		formHeadingText := "Saving log entry. Enter the following details."

		content = fmt.Sprintf(
			`
  %s

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
			taskLogEntryHeadingStyle.Render(taskLogEntryViewHeading),
			formContextStyle.Render(formHeadingText),
			formHelpStyle.Render(formHelp),
			formFieldNameStyle.Render(formBeginTimeHelp),
			m.tLInputs[entryBeginTS].View(),
			formHelpStyle.Render(formTimeShiftHelp),
			formFieldNameStyle.Render(formEndTimeHelp),
			m.tLInputs[entryEndTS].View(),
			formHelpStyle.Render(formTimeShiftHelp),
			formFieldNameStyle.Render(formCommentHelp),
			m.tLCommentInput.View(),
			formHelpStyle.Render(formSubmitHelp),
		)
		for i := 0; i < m.terminalHeight-32; i++ {
			content += "\n"
		}
	case editActiveTLView:
		formHeadingText := "Updating log entry. Enter the following details."

		content = fmt.Sprintf(
			`
  %s

  %s

  %s

  %s    %s

  %s

%s

  %s
`,
			taskLogEntryHeadingStyle.Render(taskLogEntryViewHeading),
			formContextStyle.Render(formHeadingText),
			formFieldNameStyle.Render(formBeginTimeHelp),
			m.tLInputs[entryBeginTS].View(),
			formHelpStyle.Render(formTimeShiftHelp),
			formFieldNameStyle.Render(formCommentHelp),
			m.tLCommentInput.View(),
			formHelpStyle.Render(formSubmitHelp),
		)
		for i := 0; i < m.terminalHeight-26; i++ {
			content += "\n"
		}
	case manualTasklogEntryView:
		var formHeadingText string
		switch m.tasklogSaveType {
		case tasklogInsert:
			formHeadingText = "Adding a manual log entry. Enter the following details."
		case tasklogUpdate:
			formHeadingText = "Updating log entry. Enter the following details."
		}

		content = fmt.Sprintf(
			`
  %s

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
			taskLogEntryHeadingStyle.Render(taskLogEntryViewHeading),
			formContextStyle.Render(formHeadingText),
			formHelpStyle.Render(formHelp),
			formFieldNameStyle.Render(formBeginTimeHelp),
			m.tLInputs[entryBeginTS].View(),
			formHelpStyle.Render(formTimeShiftHelp),
			formFieldNameStyle.Render(formEndTimeHelp),
			m.tLInputs[entryEndTS].View(),
			formHelpStyle.Render(formTimeShiftHelp),
			formFieldNameStyle.Render(formCommentHelp),
			m.tLCommentInput.View(),
			formHelpStyle.Render(formSubmitHelp),
		)
		for i := 0; i < m.terminalHeight-32; i++ {
			content += "\n"
		}
	case helpView:
		if !m.helpVPReady {
			content = "\n  Initializing..."
		} else {
			content = viewPortStyle.Render(fmt.Sprintf("%s  %s\n\n%s\n",
				helpTitleStyle.Render("Help"),
				helpSectionStyle.Render("(scroll with j/k/↓/↑)"),
				m.helpVP.View()))
		}
	case insufficientDimensionsView:
		return fmt.Sprintf(`
    Terminal size too small:
      Width = %d Height = %d

    Minimum dimensions needed:
      Width = %d Height = %d

    Press (q/<ctrl+c>/<esc> to exit)
`, m.terminalWidth, m.terminalHeight, minWidthNeeded, minHeightNeeded)
	}

	var helpMsg string
	if m.showHelpIndicator {
		// first time directions
		if m.activeView == taskListView && len(m.activeTasksList.Items()) <= 1 {
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

	footer = fmt.Sprintf("%s%s%s",
		toolNameStyle.Render("hours"),
		helpMsg,
		activeMsg,
	)

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
