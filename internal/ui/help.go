package ui

import "fmt"

func getHelpText(style *Style) string {
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
		style.helpHeader.Render("\"hours\" Reference Manual"),
		style.helpSection.Render(`
"hours" has 6 views:
  - Tasks List View                       Shows active tasks
  - Task Management View                  Shows a form to create/update tasks
  - Task Logs List View                   Shows your task logs
  - Task Log Details View                 Shows details for a task log
  - Inactive Tasks List View              Shows inactive tasks
  - Task Log Entry View                   Shows a form to save/update a task log entry
  - Help View (this one)
`),
		style.helpHeader.Render("Keyboard Shortcuts"),
		style.helpHeader.Render("General"),
		style.helpSection.Render(`
  1                                       Switch to Tasks List View
  2                                       Switch to Task Logs List View
  3                                       Switch to Inactive Tasks List View
  <tab>                                   Go to next view/form entry
  <shift+tab>                             Go to previous view/form entry
  q/<ctrl+c>                              Go back
  ?                                       Show help view
`),
		style.helpHeader.Render("General List Controls"),
		style.helpSection.Render(`
  k/<Up>                                  Move cursor up
  j/<Down>                                Move cursor down
  h<Left>                                 Go to previous page
  l<Right>                                Go to next page
  <ctrl+r>                                Refresh list
`),
		style.helpHeader.Render("Task List View"),
		style.helpSection.Render(`
  a                                       Add a task
  u                                       Update task details
  s                                       Start/stop recording time on a task; stopping
                                              will open up the "Task Log Entry View"
  S                                       Quick switch recording; will save a task log
                                              entry for the currently active task, and
                                              start recording time for another
  <ctrl+s>                                Edit the currently active task log/Add a new
                                              manual task log entry
  <ctrl+x>                                Discard currently active recording
  <ctrl+t>                                Go to currently tracked item
  <ctrl+d>                                Deactivate task
`),
		style.helpHeader.Render("Task Logs List View"),
		style.helpSection.Render(`
  ~ at the end of a task log comment indicates that it has more lines that are not
  visible in the list view

  d                                       Show task log details
  <ctrl+s>/u                              Update task log entry
  <ctrl+d>                                Delete task log entry
`),
		style.helpHeader.Render("Task Log Details View"),
		style.helpSection.Render(`
  h                                       Go to previous entry
  l                                       Go to next entry
`),
		style.helpHeader.Render("Inactive Task List View"),
		style.helpSection.Render(`
  <ctrl+d>                                Activate task
`),
		style.helpHeader.Render("Task Log Entry View"),
		style.helpSection.Render(`
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
