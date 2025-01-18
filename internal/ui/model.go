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

type trackingStatus uint

const (
	trackingInactive trackingStatus = iota
	trackingActive
)

type trackingChange uint

const (
	trackingStarted trackingChange = iota
	trackingFinished
)

type stateView uint

const (
	taskListView stateView = iota
	taskLogView
	taskLogDetailsView
	inactiveTaskListView
	editActiveTLView
	finishActiveTLView
	manualTasklogEntryView
	taskInputView
	helpView
	insufficientDimensionsView
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

type recordsType uint

const (
	reportRecords recordsType = iota
	reportAggRecords
	reportLogs
	reportStats
)

const (
	tasklogInsert tasklogSaveType = iota
	tasklogUpdate
)

const (
	timeFormat         = "2006/01/02 15:04"
	timeOnlyFormat     = "15:04"
	dayFormat          = "Monday"
	friendlyTimeFormat = "Mon, 15:04"
	dateFormat         = "2006/01/02"
)

type Model struct {
	activeView                     stateView
	lastView                       stateView
	lastViewBeforeInsufficientDims stateView
	db                             *sql.DB
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
	message                        string
	showHelpIndicator              bool
	terminalWidth                  int
	terminalHeight                 int
	trackingActive                 bool
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
		fetchTLS(m.db),
		fetchTasks(m.db, false),
	)
}

type recordsModel struct {
	db       *sql.DB
	typ      recordsType
	start    time.Time
	end      time.Time
	period   string
	numDays  int
	plain    bool
	report   string
	quitting bool
	busy     bool
	err      error
}

func (recordsModel) Init() tea.Cmd {
	return nil
}
