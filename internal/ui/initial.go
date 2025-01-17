package ui

import (
	"database/sql"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/dhth/hours/internal/types"
)

func InitialModel(db *sql.DB) Model {
	var activeTaskItems []list.Item
	var inactiveTaskItems []list.Item
	var tasklogListItems []list.Item

	tLInputs := make([]textinput.Model, 3)
	tLInputs[entryBeginTS] = textinput.New()
	tLInputs[entryBeginTS].Placeholder = "09:30"
	tLInputs[entryBeginTS].CharLimit = len(timeFormat)
	tLInputs[entryBeginTS].Width = 30

	tLInputs[entryEndTS] = textinput.New()
	tLInputs[entryEndTS].Placeholder = "12:30pm"
	tLInputs[entryEndTS].CharLimit = len(timeFormat)
	tLInputs[entryEndTS].Width = 30

	tLInputs[entryComment] = textinput.New()
	tLInputs[entryComment].Placeholder = "Your comment goes here"
	tLInputs[entryComment].CharLimit = 255
	tLInputs[entryComment].Width = 100

	tLDescriptionInput := textarea.New()
	tLDescriptionInput.Placeholder = `Task Log Description goes here.

This can be used to record additional details about your work on this task.`
	tLDescriptionInput.CharLimit = 1024
	tLDescriptionInput.SetWidth(100)
	tLDescriptionInput.SetHeight(8)
	tLDescriptionInput.ShowLineNumbers = false
	tLDescriptionInput.Prompt = "  ┃ "

	taskInputs := make([]textinput.Model, 3)
	taskInputs[summaryField] = textinput.New()
	taskInputs[summaryField].Placeholder = "task summary goes here"
	taskInputs[summaryField].Focus()
	taskInputs[summaryField].CharLimit = 100
	taskInputs[entryBeginTS].Width = 60

	m := Model{
		db:                db,
		activeTasksList:   list.New(activeTaskItems, newItemDelegate(lipgloss.Color(activeTaskListColor)), listWidth, 0),
		inactiveTasksList: list.New(inactiveTaskItems, newItemDelegate(lipgloss.Color(inactiveTaskListColor)), listWidth, 0),
		taskMap:           make(map[int]*types.Task),
		taskIndexMap:      make(map[int]int),
		taskLogList:       list.New(tasklogListItems, newItemDelegate(lipgloss.Color(taskLogListColor)), listWidth, 0),
		showHelpIndicator: true,
		tLInputs:          tLInputs,
		tLDescInput:       tLDescriptionInput,
		taskInputs:        taskInputs,
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
