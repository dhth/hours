package ui

import (
	"database/sql"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dhth/hours/internal/types"
)

type trackingChange uint

const (
	trackingStarted trackingChange = iota
	trackingFinished
)

type stateView uint

const (
	taskListView               stateView = iota // Main list of active tasks
	taskLogView                                 // View showing task log entries
	taskLogDetailsView                          // Detailed view of a specific task log entry
	inactiveTaskListView                        // List of inactive tasks
	editActiveTLView                            // Form to edit currently active task log (ie, begin TS)
	finishActiveTLView                          // Form to finish active task log
	manualTasklogEntryView                      // Form to manually create a new task log entry
	editSavedTLView                             // Form to edit an existing task log
	taskInputView                               // Form to create or edit task details
	moveTaskLogView                             // View to select target task for moving log entry
	helpView                                    // Help documentation view
	insufficientDimensionsView                  // Error view when terminal is too small
)

type taskMgmtContext uint

const (
	taskCreateCxt taskMgmtContext = iota
	taskUpdateCxt
)

type taskInputField uint

const (
	summaryField taskInputField = iota
)

type tLTrackingFormField uint

const (
	entryBeginTS tLTrackingFormField = iota
	entryEndTS
	entryComment
)

type tasklogSaveType uint

type recordsKind uint

const (
	reportRecords recordsKind = iota
	reportAggRecords
	reportLogs
	reportStats
)

const (
	tasklogInsert tasklogSaveType = iota
	tasklogUpdate
)

const (
	timeFormat           = "2006/01/02 15:04"
	timeOnlyFormat       = "15:04"
	dateFormat           = "2006/01/02"
	userMsgDefaultFrames = 3
)

type userMsgKind uint

const (
	userMsgInfo userMsgKind = iota
	userMsgErr
)

type userMsg struct {
	value      string
	kind       userMsgKind
	framesLeft uint
}

type logFramesConfig struct {
	log       bool
	framesDir string
}

type Model struct {
	activeView                     stateView
	lastView                       stateView
	lastViewBeforeInsufficientDims stateView
	db                             *sql.DB
	style                          Style
	timeProvider                   types.TimeProvider
	activeTasksList                list.Model
	inactiveTasksList              list.Model
	taskMap                        map[int]*types.Task
	taskIndexMap                   map[int]int
	activeTLBeginTS                time.Time
	activeTLEndTS                  time.Time
	activeTLComment                *string
	tasksFetched                   bool
	taskLogList                    list.Model
	tLInputs                       []textinput.Model
	trackingFocussedField          tLTrackingFormField
	tLCommentInput                 textarea.Model
	taskInputs                     []textinput.Model
	taskMgmtContext                taskMgmtContext
	taskInputFocussedField         taskInputField
	helpVP                         viewport.Model
	helpVPReady                    bool
	tLDetailsVP                    viewport.Model
	tLDetailsVPReady               bool
	lastTrackingChange             trackingChange
	changesLocked                  bool
	activeTaskID                   int
	tasklogSaveType                tasklogSaveType
	message                        userMsg
	showHelpIndicator              bool
	terminalWidth                  int
	terminalHeight                 int
	trackingActive                 bool
	debug                          bool
	frameCounter                   uint
	logFramesCfg                   logFramesConfig
	targetTasksList                list.Model
	moveTLID                       int
	moveOldTaskID                  int
	moveSecsSpent                  int
}

func (m *Model) blurTLTrackingInputs() {
	for i := range m.tLInputs {
		m.tLInputs[i].Blur()
	}
	m.tLCommentInput.Blur()
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		hideHelp(time.Minute*1),
		fetchTasks(m.db, true),
		fetchTLS(m.db, nil),
		fetchTasks(m.db, false),
	)
}

type recordsModel struct {
	db           *sql.DB
	style        Style
	timeProvider types.TimeProvider
	kind         recordsKind
	dateRange    types.DateRange
	period       string
	plain        bool
	taskStatus   types.TaskStatus
	report       string
	quitting     bool
	busy         bool
	err          error
}

func (recordsModel) Init() tea.Cmd {
	return nil
}

func infoMsg(msg string) userMsg {
	return userMsg{
		value:      msg,
		kind:       userMsgInfo,
		framesLeft: userMsgDefaultFrames,
	}
}

func errMsg(msg string) userMsg {
	return userMsg{
		value:      msg,
		kind:       userMsgErr,
		framesLeft: userMsgDefaultFrames,
	}
}

func errMsgQuick(msg string) userMsg {
	return userMsg{
		value:      msg,
		kind:       userMsgErr,
		framesLeft: 2,
	}
}
