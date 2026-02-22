# Implementation Plan: Move Log Entry to Another Task

## Overview
This feature allows users to move a task log entry from one task (parent) to another task. The implementation follows the existing UI patterns in the codebase (Bubble Tea TUI framework).

## Feature Requirements (from user discussion)
1. **Keyboard shortcut**: `m` (no conflicts found in codebase)
2. **Task list**: Only show active tasks as potential targets
3. **After moving**: Return to task log view and refresh the list
4. **Time recalculation**: Automatically update secs_spent on both old and new parent tasks

## Architecture

### Technology Stack
- **Language**: Go
- **TUI Framework**: Bubble Tea (charmbracelet/bubbletea)
- **UI Components**: Bubbles (list, textinput, textarea)
- **Styling**: Lipgloss
- **Database**: SQLite (modernc.org/sqlite)
- **CLI**: Cobra

### UI Pattern
The app uses **view state switching** instead of modal dialogs. Each "screen" is a different `stateView` value, and the UI renders different content based on the active view.

## Files to Modify (9 files)

### 1. internal/persistence/queries.go
**Purpose**: Add database operation to move log entry between tasks

**Add new function**:
```go
func MoveTaskLog(db *sql.DB, tlID int, oldTaskID int, newTaskID int, secsSpent int) error
```

**Logic**:
- Use `runInTx()` for atomic transaction
- Update task_log.task_id to newTaskID
- Decrease old task's secs_spent by secsSpent
- Increase new task's secs_spent by secsSpent
- Update timestamps on both tasks

**Reference**: Similar pattern to `DeleteTL()` (line 703) which already handles secs_spent updates

---

### 2. internal/ui/msgs.go
**Purpose**: Add message type for move operation result

**Add new type** (around line 95):
```go
type taskLogMovedMsg struct {
	tlID      int
	oldTaskID int
	newTaskID int
	err       error
}
```

---

### 3. internal/ui/cmds.go
**Purpose**: Add command to execute the database operation

**Add new function** (after line 168):
```go
func moveTaskLog(db *sql.DB, tlID int, oldTaskID int, newTaskID int, secsSpent int) tea.Cmd {
	return func() tea.Msg {
		err := pers.MoveTaskLog(db, tlID, oldTaskID, newTaskID, secsSpent)
		return taskLogMovedMsg{tlID, oldTaskID, newTaskID, err}
	}
}
```

---

### 4. internal/ui/model.go
**Purpose**: Add new view state and tracking fields

**Add new state constant** (after line 35):
```go
moveTaskLogView           stateView = iota // View to select target task for moving log entry
```

**Add new fields to Model struct** (after line 137):
```go	targetTasksList    list.Model    // List of active tasks for selecting move target
	moveTLID           int           // ID of task log entry being moved
	moveOldTaskID      int           // ID of original parent task
	moveSecsSpent      int           // Seconds spent on the log entry being moved
```

---

### 5. internal/ui/initial.go
**Purpose**: Initialize the target task list component

**In the initialization function** (where activeTasksList is initialized around line 60):
Add similar initialization for `targetTasksList`:
```go
m.targetTasksList = list.New([]list.Item{}, newItemDelegate(style.theme), 0, 0)
m.targetTasksList.Title = "Select Target Task"
m.targetTasksList.SetStatusBarItemName("task", "tasks")
m.targetTasksList.DisableQuitKeybindings()
m.targetTasksList.SetShowHelp(false)
m.targetTasksList.Styles.Title = m.targetTasksList.Styles.Title.
	Foreground(lipgloss.Color(style.theme.TitleForeground)).
	Background(lipgloss.Color(style.theme.ActiveTasks)).
	Bold(true)
m.targetTasksList.KeyMap.PrevPage.SetKeys("left", "h", "pgup")
m.targetTasksList.KeyMap.NextPage.SetKeys("right", "l", "pgdown")
```

---

### 6. internal/ui/update.go
**Purpose**: Add keyboard handling and view state management

**A. Add "m" key handler in taskLogView section** (around line 227, inside taskLogView case):
```go
case "m":
	m.handleRequestToMoveTaskLog()
```

**B. Add moveTaskLogView case in view switch** (around the view state handling, similar to other views):
```go
case moveTaskLogView:
	// Delegate to targetTasksList's Update method
	var cmd tea.Cmd
	m.targetTasksList, cmd = m.targetTasksList.Update(msg)
	cmds = append(cmds, cmd)
```

**C. Add taskLogMovedMsg handler** (in the message type switch, around line 315):
```go
case taskLogMovedMsg:
	if msg.err != nil {
		m.message = errMsg(fmt.Sprintf("Error moving task log: %s", msg.err))
	} else {
		cmds = append(cmds, fetchTLS(m.db, nil))
		cmds = append(cmds, fetchTasks(m.db, true))
	}
	m.activeView = taskLogView
```

---

### 7. internal/ui/handle.go
**Purpose**: Add handlers for move request and target selection

**A. Add move request handler** (new function):
```go
func (m *Model) handleRequestToMoveTaskLog() tea.Cmd {
	if m.taskLogList.IsFiltered() {
		m.message = errMsg(removeFilterMsg)
		return nil
	}

	entry, ok := m.taskLogList.SelectedItem().(*types.TaskLogEntry)
	if !ok {
		m.message = errMsg(msgCouldntSelectATask)
		return nil
	}

	// Store the log entry details
	m.moveTLID = entry.ID
	m.moveOldTaskID = entry.TaskID
	m.moveSecsSpent = entry.SecsSpent

	// Initialize target list with active tasks, excluding current parent
	targetItems := []list.Item{}
	for i := range m.activeTasksList.Items() {
		task, ok := m.activeTasksList.Items()[i].(*types.Task)
		if !ok {
			continue
		}
		// Exclude the current parent task
		if task.ID != entry.TaskID {
			targetItems = append(targetItems, task)
		}
	}
	m.targetTasksList.SetItems(targetItems)

	m.activeView = moveTaskLogView
	return nil
}
```

**B. Add target selection handler** (new function):
```go
func (m *Model) handleTargetTaskSelection() tea.Cmd {
	if m.targetTasksList.IsFiltered() {
		m.message = errMsg(removeFilterMsg)
		return nil
	}

	task, ok := m.targetTasksList.SelectedItem().(*types.Task)
	if !ok {
		m.message = errMsg(msgCouldntSelectATask)
		return nil
	}

	return moveTaskLog(m.db, m.moveTLID, m.moveOldTaskID, task.ID, m.moveSecsSpent)
}
```

**C. Update handleEnterKey function** to handle moveTaskLogView:
Around line where other views handle Enter key (around line 500), add:
```go
case moveTaskLogView:
	cmd := m.handleTargetTaskSelection()
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
```

**D. Update handleEscapeInForms function** to handle moveTaskLogView:
Around line 150, add:
```go
case moveTaskLogView:
	m.activeView = taskLogView
	m.targetTasksList.ResetFilter()
```

---

### 8. internal/ui/view.go
**Purpose**: Add rendering for the move task log view

**Add new case in view switch** (around line 102):
```go
case moveTaskLogView:
	content = m.style.list.Render(m.targetTasksList.View())
```

---

### 9. internal/ui/help.go
**Purpose**: Add help text for the new "m" shortcut

**Add entry under "Task Logs List View" section** (around line 77-79):
```
  m                                       Move task log entry to another task
```

Update the section to:
```
  d                                       Show task log details
  <ctrl+s>/u                              Update task log entry
  <ctrl+d>                                Delete task log entry
  m                                       Move task log entry to another task
```

---

## Implementation Flow

```
User in taskLogView
    ↓
Press "m" key
    ↓
handleRequestToMoveTaskLog() called
    ↓
Validate no filter active on taskLogList
    ↓
Get selected task log entry
    ↓
Store: moveTLID, moveOldTaskID, moveSecsSpent
    ↓
Initialize targetTasksList (active tasks minus current parent)
    ↓
Switch to moveTaskLogView
    ↓
User sees filtered list of target tasks
    ↓
User can filter/search or select task + press Enter
    ↓
handleTargetTaskSelection() called
    ↓
Validate no filter active on targetTasksList
    ↓
Get selected target task
    ↓
Execute moveTaskLog() command
    ↓
MoveTaskLog() in persistence layer:
    - UPDATE task_log SET task_id = newTaskID WHERE id = tlID
    - UPDATE task SET secs_spent = secs_spent - ? WHERE id = oldTaskID
    - UPDATE task SET secs_spent = secs_spent + ? WHERE id = newTaskID
    ↓
taskLogMovedMsg received
    ↓
If success: fetchTLS(), fetchTasks() to refresh lists
    ↓
Return to taskLogView
```

## Testing Considerations

1. **Database transaction**: Ensure all updates happen atomically
2. **Time calculations**: Verify secs_spent is correctly subtracted/added
3. **Edge cases**: 
   - What if target task is deleted during the move?
   - What if the log entry is deleted during the move?
4. **UI state**: Ensure filter state is properly reset when entering/exiting move view
5. **Help text**: Verify "m" appears in help documentation

## Database Schema Context

**task table** (from init.go line 21-28):
- id, summary, secs_spent, active, created_at, updated_at

**task_log table** (from init.go line 30-39):
- id, task_id (FK to task), begin_ts, end_ts, secs_spent, comment, active

The secs_spent field on task is a denormalized cache of the sum of all task_log entries for that task. When moving an entry, we need to:
1. Update the task_log.task_id reference
2. Subtract secs_spent from old task
3. Add secs_spent to new task

## Existing Patterns to Follow

1. **Transaction handling**: Use `runInTx()` from persistence/queries.go (line 754)
2. **Message handling**: Follow pattern of tLDeletedMsg (line 80-83 in msgs.go)
3. **View switching**: Use the stateView pattern in model.go
4. **List filtering**: Check `list.IsFiltered()` before operations (see handle.go line 393-397)
5. **Error messages**: Use `errMsg()` function from model.go (line 183)

## Notes for Future Implementation

- The `m` key is confirmed to have no conflicts in the codebase
- Only active tasks are shown as targets (not inactive ones)
- The current parent task is excluded from the target list
- This implementation reuses the existing list component and filtering mechanism
- No new external dependencies needed
