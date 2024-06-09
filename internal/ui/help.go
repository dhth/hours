package ui

import "fmt"

var (
	helpText = fmt.Sprintf(`
  %s
%s
  %s

  %s
%s
  %s
%s
  %s
%s
  %s
%s
  %s
%s
`,
		helpHeaderStyle.Render("\"hours\" Reference Manual"),
		helpSectionStyle.Render(`
  (scroll line by line with j/k/arrow keys or by half a page with <c-d>/<c-u>)

  "hours" has a simple to use TUI, indended for those who want to track time on the tasks they
  care about with minimal keypresses.

  "hours" has 4 panes:
    - Tasks List View                      Shows your tasks
    - Task Management View                 Allows you to create/update tasks
    - Task Log List View                   Shows your task log entries
    - Help View (this one)
`),
		helpHeaderStyle.Render("Keyboard Shortcuts"),
		helpHeaderStyle.Render("General"),
		helpSectionStyle.Render(`
    1                                       Switch to Tasks List View
    2                                       Switch to Task Log List View
    <tab>                                   Go to next view/form entry
    <shift+tab>                             Go to previous view/form entry
      ?                                     Show help view
`),
		helpHeaderStyle.Render("General List Controls"),
		helpSectionStyle.Render(`
    h/<Up>                                  Move cursor up
    k/<Down>                                Move cursor down
    h<Left>                                 Go to previous page
    l<Right>                                Go to next page
    /                                       Start filtering
`),
		helpHeaderStyle.Render("Task List View"),
		helpSectionStyle.Render(`
    s                                       Toggle recording time on the currently selected task,
                                                will open up a form to record a comment on the
                                                second "s" keypress
    <ctrl+s>                                Add a manual task log entry
    <ctrl+t>                                Go to currently tracked item
`),
		helpHeaderStyle.Render("Task Log List View"),
		helpSectionStyle.Render(`
    <ctrl+d>                                Delete task log entry
    <ctrl+r>                                Refresh list
`),
		helpHeaderStyle.Render("Task Log Entry View"),
		helpSectionStyle.Render(`
    enter                                   Save task log entry
`),
	)
)
