package ui

import "fmt"

func getHelpText(style Style) string {
	return fmt.Sprintf(`%s
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
%s
%s
%s`,
		style.helpPrimary.Render("\"hours\" Reference Manual"),
		style.helpSecondary.Render(`
"hours" has 6 views:
  - Tasks List View                       Shows active tasks
  - Task Management View                  Shows a form to create/update tasks
  - Task Logs List View                   Shows your task logs
  - Task Log Details View                 Shows details for a task log
  - Inactive Tasks List View              Shows inactive tasks
  - Task Log Entry View                   Shows a form to save/update a task log entry
  - Help View (this one)
`),
		style.helpPrimary.Render("Keyboard Shortcuts"),
		style.helpPrimary.Render("General"),
		style.helpSecondary.Render(`
  1                                       Switch to Tasks List View
  2                                       Switch to Task Logs List View
  3                                       Switch to Inactive Tasks List View
  <tab>                                   Go to next view/form entry
  <shift+tab>                             Go to previous view/form entry
  q/<esc>                                 Go back or quit
  <ctrl+c>                                Quit immediately
  ?                                       Show help view
`),
		style.helpPrimary.Render("General List Controls"),
		style.helpSecondary.Render(`
  k/<Up>                                  Move cursor up
  j/<Down>                                Move cursor down
  h<Left>                                 Go to previous page
  l<Right>                                Go to next page
  <ctrl+r>                                Refresh list
`),
		style.helpPrimary.Render("Task List View"),
		style.helpSecondary.Render(`
  a                                       Add a task
  u                                       Update task details
  c                                       Copy task summary to clipboard
  s                                       Start/stop recording time on a task; stopping
                                              will open up the "Task Log Entry View"
  S                                       Quick switch recording; will save a task log
                                              entry for the currently active task, and
                                              start recording time for another
  f                                       Finish the currently active task log without
                                              comment
  <ctrl+s>                                Edit the currently active task log/Add a new
                                              manual task log entry
  <ctrl+x>                                Discard currently active recording
  <ctrl+t>                                Go to currently tracked item
  <ctrl+d>                                Deactivate task
`),
		style.helpPrimary.Render("Task Logs List View"),
		style.helpSecondary.Render(`
  ~ at the end of a task log comment indicates that it has more lines that are not
  visible in the list view

  d                                       Show task log details
  <ctrl+s>/u                              Update task log entry
  <ctrl+d>                                Delete task log entry
  m                                       Move task log entry to another task
`),
		style.helpPrimary.Render("Task Log Details View"),
		style.helpSecondary.Render(`
  h                                       Go to previous entry
  l                                       Go to next entry
`),
		style.helpPrimary.Render("Inactive Task List View"),
		style.helpSecondary.Render(`
  c                                       Copy task summary to clipboard
  <ctrl+d>                                Activate task
`),
		style.helpPrimary.Render("Task Log Entry View"),
		style.helpSecondary.Render(`
  enter/<ctrl+s>                          Save entered details for the task log
  k                                       Move timestamp backwards by one minute
  j                                       Move timestamp forwards by one minute
  K                                       Move timestamp backwards by five minutes
  J                                       Move timestamp forwards by five minutes
  h                                       Move timestamp backwards by a day
  l                                       Move timestamp forwards by a day
`),
	)
}
