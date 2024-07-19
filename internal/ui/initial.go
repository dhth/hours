package ui

import (
	"database/sql"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

func InitialModel(db *sql.DB) model {
	var activeTaskItems []list.Item
	var inactiveTaskItems []list.Item
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
	trackingInputs[entryComment].Width = 80

	taskInputs := make([]textinput.Model, 3)
	taskInputs[summaryField] = textinput.New()
	taskInputs[summaryField].Placeholder = "task summary goes here"
	taskInputs[summaryField].Focus()
	taskInputs[summaryField].CharLimit = 100
	taskInputs[entryBeginTS].Width = 60

	m := model{
		db:                 db,
		activeTasksList:    list.New(activeTaskItems, newItemDelegate(lipgloss.Color(activeTaskListColor)), listWidth, 0),
		inactiveTasksList:  list.New(inactiveTaskItems, newItemDelegate(lipgloss.Color(inactiveTaskListColor)), listWidth, 0),
		activeTaskMap:      make(map[int]*task),
		activeTaskIndexMap: make(map[int]int),
		taskLogList:        list.New(tasklogListItems, newItemDelegate(lipgloss.Color(taskLogListColor)), listWidth, 0),
		showHelpIndicator:  true,
		trackingInputs:     trackingInputs,
		taskInputs:         taskInputs,
	}
	m.activeTasksList.Title = "Tasks"
	m.activeTasksList.SetStatusBarItemName("task", "tasks")
	m.activeTasksList.DisableQuitKeybindings()
	m.activeTasksList.SetShowHelp(false)
	m.activeTasksList.Styles.Title = m.activeTasksList.Styles.Title.Foreground(lipgloss.Color(defaultBackgroundColor)).Background(lipgloss.Color(activeTaskListColor)).Bold(true)

	m.taskLogList.Title = "Task Logs (last 50)"
	m.taskLogList.SetStatusBarItemName("entry", "entries")
	m.taskLogList.SetFilteringEnabled(false)
	m.taskLogList.DisableQuitKeybindings()
	m.taskLogList.SetShowHelp(false)
	m.taskLogList.Styles.Title = m.taskLogList.Styles.Title.Foreground(lipgloss.Color(defaultBackgroundColor)).Background(lipgloss.Color(taskLogListColor)).Bold(true)

	m.inactiveTasksList.Title = "Inactive Tasks"
	m.inactiveTasksList.SetStatusBarItemName("task", "tasks")
	m.inactiveTasksList.DisableQuitKeybindings()
	m.inactiveTasksList.SetShowHelp(false)
	m.inactiveTasksList.Styles.Title = m.inactiveTasksList.Styles.Title.Foreground(lipgloss.Color(defaultBackgroundColor)).Background(lipgloss.Color(inactiveTaskListColor)).Bold(true)

	return m
}

func initialRecordsModel(typ recordsType, db *sql.DB, start, end time.Time, plain bool, period string, numDays int, initialData string) recordsModel {

	return recordsModel{
		typ:     typ,
		db:      db,
		start:   start,
		end:     end,
		period:  period,
		numDays: numDays,
		plain:   plain,
		report:  initialData,
	}
}
