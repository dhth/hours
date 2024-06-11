package ui

import (
	"database/sql"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type trackingStatus uint

const (
	trackingInactive trackingStatus = iota
	trackingActive
)

type dBChange uint

const (
	insertChange dBChange = iota
	updateChange
)

type stateView uint

const (
	activeTaskListView stateView = iota
	taskLogView
	inactiveTaskListView
	askForCommentView
	manualTasklogEntryView
	taskInputView
	helpView
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

type trackingFocussedField uint

const (
	entryBeginTS trackingFocussedField = iota
	entryEndTS
	entryComment
)

type tasklogSaveType uint

const (
	tasklogInsert tasklogSaveType = iota
	tasklogUpdate
)

const (
	timeFormat         = "2006/01/02 15:04"
	friendlyTimeFormat = "Mon, 15:04"
	dateFormat         = "2006/01/02"
)

type model struct {
	activeView             stateView
	lastView               stateView
	db                     *sql.DB
	activeTasksList        list.Model
	inactiveTasksList      list.Model
	activeTaskMap          map[int]*task
	activeTaskIndexMap     map[int]int
	activeTLBeginTS        time.Time
	activeTLEndTS          time.Time
	tasksFetched           bool
	taskLogList            list.Model
	trackingInputs         []textinput.Model
	trackingFocussedField  trackingFocussedField
	taskInputs             []textinput.Model
	taskMgmtContext        taskMgmtContext
	taskInputFocussedField taskInputField
	helpVP                 viewport.Model
	helpVPReady            bool
	lastChange             dBChange
	changesLocked          bool
	activeTaskId           int
	tasklogSaveType        tasklogSaveType
	message                string
	messages               []string
	showHelpIndicator      bool
	terminalHeight         int
	trackingActive         bool
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		hideHelp(time.Minute*1),
		fetchTasks(m.db, true),
		fetchTaskLogEntries(m.db),
		fetchTasks(m.db, false),
	)
}
