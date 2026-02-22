package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/charmbracelet/lipgloss"
	"github.com/dhth/hours/internal/types"
	"github.com/dhth/hours/internal/utils"
)

const (
	taskLogEntryViewHeading = "Task Log Entry"
	minHeightNeeded         = 32
	minWidthNeeded          = 80
	tlWarningThresholdSecs  = 8 * 60 * 60
)

var listWidth = 140

type tlFormValidity uint

const (
	tlSubmitOk tlFormValidity = iota
	tlSubmitWarn
	tlSubmitErr
)

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func (m Model) View() string {
	var content string
	var footer string

	var statusBar string
	if m.message.framesLeft > 0 && m.message.value != "" {
		statusBar = m.message.value
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
			m.style.tracking.Render("tracking:"),
			m.style.activeTaskSummaryMsg.Render(taskSummaryMsg),
			m.style.activeTaskBeginTime.Render(taskStartedSinceMsg),
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

	var submissionCtx string
	var submissionValidity tlFormValidity
	var durationCtx string
	if m.activeView == finishActiveTLView || m.activeView == manualTasklogEntryView || m.activeView == editSavedTLView {
		durationCtx, submissionValidity = getDurationValidityContext(m.tLInputs[entryBeginTS].Value(), m.tLInputs[entryEndTS].Value())

		switch submissionValidity {
		case tlSubmitOk:
			submissionCtx = m.style.tlFormOkStyle.Render(durationCtx)
		case tlSubmitWarn:
			submissionCtx = m.style.tlFormWarnStyle.Render(durationCtx)
		case tlSubmitErr:
			submissionCtx = m.style.tlFormErrStyle.Render(durationCtx)
		}
	}

	var formSubmitHelp string
	switch m.activeView {
	case taskInputView:
		formSubmitHelp = "Press <ctrl+s>/<enter> to submit"
	case editActiveTLView, finishActiveTLView, manualTasklogEntryView, editSavedTLView:
		if submissionValidity != tlSubmitErr {
			if m.trackingFocussedField == entryComment {
				formSubmitHelp = m.style.formHelp.Render("Press <ctrl+s> to submit")
			} else {
				formSubmitHelp = m.style.formHelp.Render("Press <ctrl+s>/<enter> to submit")
			}
		}
	}

	switch m.activeView {
	case taskListView:
		content = m.style.list.Render(m.activeTasksList.View())
	case taskLogView:
		content = m.style.list.Render(m.taskLogList.View())
	case taskLogDetailsView:
		if !m.helpVPReady {
			content = "\n  Initializing..."
		} else {
			content = m.style.viewPort.Render(fmt.Sprintf("%s\n\n%s",
				m.style.taskLogDetails.Render("Task Log Details"), m.tLDetailsVP.View()))
		}
	case inactiveTaskListView:
		content = m.style.list.Render(m.inactiveTasksList.View())
	case taskInputView:
		var formTitle string
		switch m.taskMgmtContext {
		case taskCreateCxt:
			formTitle = "Add task"
		case taskUpdateCxt:
			formTitle = "Update task"
		}
		content = fmt.Sprintf(
			`
  %s

  %s

  %s
`,
			m.style.taskEntryHeading.Render(formTitle),
			m.taskInputs[summaryField].View(),
			m.style.formHelp.Render(formSubmitHelp),
		)
		for range m.terminalHeight - 9 {
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

  %s
`,
			m.style.taskLogEntryHeading.Render(taskLogEntryViewHeading),
			m.style.formContext.Render(formHeadingText),
			m.style.formHelp.Render(formHelp),
			m.style.formFieldName.Render(formBeginTimeHelp),
			m.tLInputs[entryBeginTS].View(),
			m.style.formHelp.Render(formTimeShiftHelp),
			m.style.formFieldName.Render(formEndTimeHelp),
			m.tLInputs[entryEndTS].View(),
			m.style.formHelp.Render(formTimeShiftHelp),
			m.style.formFieldName.Render(formCommentHelp),
			m.tLCommentInput.View(),
			submissionCtx,
			formSubmitHelp,
		)
		for range m.terminalHeight - 34 {
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
			m.style.taskLogEntryHeading.Render(taskLogEntryViewHeading),
			m.style.formContext.Render(formHeadingText),
			m.style.formFieldName.Render(formBeginTimeHelp),
			m.tLInputs[entryBeginTS].View(),
			m.style.formHelp.Render(formTimeShiftHelp),
			m.style.formFieldName.Render(formCommentHelp),
			m.tLCommentInput.View(),
			m.style.formHelp.Render(formSubmitHelp),
		)
		for range m.terminalHeight - 26 {
			content += "\n"
		}
	case manualTasklogEntryView, editSavedTLView:
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

  %s
`,
			m.style.taskLogEntryHeading.Render(taskLogEntryViewHeading),
			m.style.formContext.Render(formHeadingText),
			m.style.formHelp.Render(formHelp),
			m.style.formFieldName.Render(formBeginTimeHelp),
			m.tLInputs[entryBeginTS].View(),
			m.style.formHelp.Render(formTimeShiftHelp),
			m.style.formFieldName.Render(formEndTimeHelp),
			m.tLInputs[entryEndTS].View(),
			m.style.formHelp.Render(formTimeShiftHelp),
			m.style.formFieldName.Render(formCommentHelp),
			m.tLCommentInput.View(),
			submissionCtx,
			formSubmitHelp,
		)
		for range m.terminalHeight - 34 {
			content += "\n"
		}
	case moveTaskLogView:
		helpText := "Press <enter> to move task log, <esc>/<q> to cancel"
		content = m.style.list.Render(m.targetTasksList.View()) + "\n\n" + m.style.formHelp.Render(helpText)
	case helpView:
		if !m.helpVPReady {
			content = "\n  Initializing..."
		} else {
			content = m.style.viewPort.Render(fmt.Sprintf("%s  %s\n\n%s\n",
				m.style.helpTitle.Render("Help"),
				m.style.helpSecondary.Render("(scroll with j/k/↓/↑)"),
				m.helpVP.View()))
		}
	case insufficientDimensionsView:
		return fmt.Sprintf(`
  Terminal size too small:
    Width = %d Height = %d

  Minimum dimensions needed:
    Width = %d Height = %d

  Press q/<ctrl+c>/<esc>
    to exit
`, m.terminalWidth, m.terminalHeight, minWidthNeeded, minHeightNeeded)
	}

	var helpMsg string
	if m.showHelpIndicator {
		// first time directions
		if m.activeView == taskListView && len(m.activeTasksList.Items()) <= 1 {
			if len(m.activeTasksList.Items()) == 0 {
				helpMsg += " " + m.style.initialHelpMsg.Render("Press a to add a task")
			} else if len(m.taskLogList.Items()) == 0 {
				if m.trackingActive {
					helpMsg += " " + m.style.initialHelpMsg.Render("Press s to stop tracking time")
				} else {
					helpMsg += " " + m.style.initialHelpMsg.Render("Press s to start tracking time")
				}
			}
		}

		helpMsg += " " + m.style.helpMsg.Render("Press ? for help")
	}

	footer = fmt.Sprintf("%s%s%s",
		m.style.toolName.Render("hours"),
		helpMsg,
		activeMsg,
	)

	if m.debug {
		footer = fmt.Sprintf("%s [term: %dx%d] [msg frames left: %d] [frames rendered: %d]",
			footer,
			m.terminalWidth,
			m.terminalHeight,
			m.message.framesLeft,
			m.frameCounter,
		)
	}

	result := lipgloss.JoinVertical(lipgloss.Left,
		content,
		statusBar,
		footer,
	)

	if m.logFramesCfg.log {
		logFrame(result, m.frameCounter, m.logFramesCfg.framesDir)
	}

	return result
}

func (m recordsModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("Something went wrong: %s\n", m.err)
	}
	var help string

	var dateRangeStr string
	var dateRange string
	if m.dateRange.NumDays > 1 {
		dateRangeStr = fmt.Sprintf(`
 range:             %s...%s
 `,
			m.dateRange.Start.Format(dateFormat), m.dateRange.End.AddDate(0, 0, -1).Format(dateFormat))
	} else {
		dateRangeStr = fmt.Sprintf(`
 date:              %s
`,
			m.dateRange.Start.Format(dateFormat))
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
		help = m.style.recordsHelp.Render(helpStr)
		dateRange = m.style.recordsDateRange.Render(dateRangeStr)
	}

	return fmt.Sprintf("%s%s%s", m.report, dateRange, help)
}

func getDurationValidityContext(beginStr, endStr string) (string, tlFormValidity) {
	beginTS, endTS, err := types.ParseTaskLogTimes(beginStr, endStr)
	if err != nil {
		return fmt.Sprintf("Error: %s", err.Error()), tlSubmitErr
	}

	dur := endTS.Sub(beginTS)
	totalSeconds := int(dur.Seconds())

	humanized := types.HumanizeDuration(totalSeconds)
	msg := fmt.Sprintf("You're recording %s", humanized)
	if totalSeconds > tlWarningThresholdSecs {
		return msg, tlSubmitWarn
	}

	return msg, tlSubmitOk
}

func logFrame(content string, frameIndex uint, framesDir string) {
	cleanContent := stripANSI(content)

	filename := fmt.Sprintf("%d.txt", frameIndex)
	filepath := filepath.Join(framesDir, filename)

	_ = os.WriteFile(filepath, []byte(cleanContent), 0o644)
}

func stripANSI(s string) string {
	return ansiRegex.ReplaceAllString(s, "")
}
