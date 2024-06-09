package ui

import (
	"database/sql"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

func InitialModel(db *sql.DB) model {
	var stackItems []list.Item
	var tasklogListItems []list.Item

	trackingInputs := make([]textinput.Model, 3)
	trackingInputs[entryBeginTS] = textinput.New()
	trackingInputs[entryBeginTS].Placeholder = "09:30"
	trackingInputs[entryBeginTS].Focus()
	trackingInputs[entryBeginTS].CharLimit = len(string(timeFormat))
	trackingInputs[entryBeginTS].Width = 30

	trackingInputs[entryEndTS] = textinput.New()
	trackingInputs[entryEndTS].Placeholder = "12:30pm"
	trackingInputs[entryEndTS].Focus()
	trackingInputs[entryEndTS].CharLimit = len(string(timeFormat))
	trackingInputs[entryEndTS].Width = 30

	trackingInputs[entryComment] = textinput.New()
	trackingInputs[entryComment].Placeholder = "Your comment goes here"
	trackingInputs[entryComment].Focus()
	trackingInputs[entryComment].CharLimit = 255
	trackingInputs[entryComment].Width = 60

	taskInputs := make([]textinput.Model, 3)
	taskInputs[summaryField] = textinput.New()
	taskInputs[summaryField].Placeholder = "task summary goes here"
	taskInputs[summaryField].Focus()
	taskInputs[summaryField].CharLimit = 100
	taskInputs[entryBeginTS].Width = 60

	m := model{
		db:                db,
		taskList:          list.New(stackItems, newItemDelegate(lipgloss.Color(taskListColor)), listWidth, 0),
		taskMap:           make(map[int]*task),
		taskIndexMap:      make(map[int]int),
		taskLogList:       list.New(tasklogListItems, newItemDelegate(lipgloss.Color(taskLogListColor)), listWidth, 0),
		showHelpIndicator: true,
		trackingInputs:    trackingInputs,
		taskInputs:        taskInputs,
	}
	m.taskList.Title = "Tasks"
	m.taskList.SetStatusBarItemName("task", "tasks")
	m.taskList.DisableQuitKeybindings()
	m.taskList.SetShowHelp(false)
	m.taskList.Styles.Title = m.taskList.Styles.Title.Foreground(lipgloss.Color(defaultBackgroundColor)).Background(lipgloss.Color(taskListColor)).Bold(true)

	m.taskLogList.Title = "Task Log"
	m.taskLogList.SetStatusBarItemName("entry", "entries")
	m.taskLogList.SetFilteringEnabled(false)
	m.taskLogList.DisableQuitKeybindings()
	m.taskLogList.SetShowHelp(false)
	m.taskLogList.Styles.Title = m.taskLogList.Styles.Title.Foreground(lipgloss.Color(defaultBackgroundColor)).Background(lipgloss.Color(taskLogListColor)).Bold(true)

	return m
}
