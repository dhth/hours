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

const (
	tlCommentLengthLimit = 3000
	textInputWidth       = 80
)

func InitialModel(db *sql.DB, style Style) Model {
	var activeTaskItems []list.Item
	var inactiveTaskItems []list.Item
	var tasklogListItems []list.Item

	tLInputs := make([]textinput.Model, 2)
	tLInputs[entryBeginTS] = textinput.New()
	tLInputs[entryBeginTS].Placeholder = "09:30"
	tLInputs[entryBeginTS].CharLimit = len(timeFormat)
	tLInputs[entryBeginTS].Width = 30

	tLInputs[entryEndTS] = textinput.New()
	tLInputs[entryEndTS].Placeholder = "12:30pm"
	tLInputs[entryEndTS].CharLimit = len(timeFormat)
	tLInputs[entryEndTS].Width = 30

	tLCommentInput := textarea.New()
	tLCommentInput.Placeholder = `Task log comment goes here.

This can be used to record details about your work on this task.`
	tLCommentInput.CharLimit = tlCommentLengthLimit
	tLCommentInput.SetWidth(textInputWidth)
	tLCommentInput.SetHeight(10)
	tLCommentInput.ShowLineNumbers = false
	tLCommentInput.Prompt = "  ┃ "

	taskInputs := make([]textinput.Model, 1)
	taskInputs[summaryField] = textinput.New()
	taskInputs[summaryField].Placeholder = "task summary goes here"
	taskInputs[summaryField].Focus()
	taskInputs[summaryField].CharLimit = 100
	taskInputs[entryBeginTS].Width = textInputWidth

	m := Model{
		db:    db,
		style: style,
		activeTasksList: list.New(activeTaskItems,
			newItemDelegate(style.listItemTitleColor,
				style.listItemDescColor,
				lipgloss.Color(style.theme.ActiveTasks),
			), listWidth, 0),
		inactiveTasksList: list.New(inactiveTaskItems,
			newItemDelegate(style.listItemTitleColor,
				style.listItemDescColor,
				lipgloss.Color(style.theme.InactiveTasks),
			), listWidth, 0),
		taskMap:      make(map[int]*types.Task),
		taskIndexMap: make(map[int]int),
		taskLogList: list.New(tasklogListItems,
			newItemDelegate(style.listItemTitleColor,
				style.listItemDescColor,
				lipgloss.Color(style.theme.TaskLogList),
			), listWidth, 0),
		showHelpIndicator: true,
		tLInputs:          tLInputs,
		tLCommentInput:    tLCommentInput,
		taskInputs:        taskInputs,
	}
	m.activeTasksList.Title = "Tasks"
	m.activeTasksList.SetStatusBarItemName("task", "tasks")
	m.activeTasksList.DisableQuitKeybindings()
	m.activeTasksList.SetShowHelp(false)
	m.activeTasksList.Styles.Title = m.activeTasksList.Styles.Title.
		Foreground(lipgloss.Color(style.theme.TitleForeground)).
		Background(lipgloss.Color(style.theme.ActiveTasks)).
		Bold(true)
	m.activeTasksList.KeyMap.PrevPage.SetKeys("left", "h", "pgup")
	m.activeTasksList.KeyMap.NextPage.SetKeys("right", "l", "pgdown")

	m.taskLogList.Title = "Task Logs (last 50)"
	m.taskLogList.SetStatusBarItemName("entry", "entries")
	m.taskLogList.SetFilteringEnabled(false)
	m.taskLogList.DisableQuitKeybindings()
	m.taskLogList.SetShowHelp(false)
	m.taskLogList.Styles.Title = m.taskLogList.Styles.Title.
		Foreground(lipgloss.Color(style.theme.TitleForeground)).
		Background(lipgloss.Color(style.theme.TaskLogList)).
		Bold(true)
	m.taskLogList.KeyMap.PrevPage.SetKeys("left", "h", "pgup")
	m.taskLogList.KeyMap.NextPage.SetKeys("right", "l", "pgdown")

	m.inactiveTasksList.Title = "Inactive Tasks"
	m.inactiveTasksList.SetStatusBarItemName("task", "tasks")
	m.inactiveTasksList.DisableQuitKeybindings()
	m.inactiveTasksList.SetShowHelp(false)
	m.inactiveTasksList.Styles.Title = m.inactiveTasksList.Styles.Title.
		Foreground(lipgloss.Color(style.theme.TitleForeground)).
		Background(lipgloss.Color(style.theme.InactiveTasks)).
		Bold(true)
	m.inactiveTasksList.KeyMap.PrevPage.SetKeys("left", "h", "pgup")
	m.inactiveTasksList.KeyMap.NextPage.SetKeys("right", "l", "pgdown")

	return m
}

func initialRecordsModel(typ recordsType, db *sql.DB, style Style, start, end time.Time, plain bool, period string, numDays int, initialData string) recordsModel {
	return recordsModel{
		typ:     typ,
		db:      db,
		style:   style,
		start:   start,
		end:     end,
		period:  period,
		numDays: numDays,
		plain:   plain,
		report:  initialData,
	}
}
