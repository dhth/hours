package types

import (
	"fmt"
	"math"
	"time"

	"github.com/dhth/hours/internal/utils"
	"github.com/dustin/go-humanize"
)

type Task struct {
	ID             int
	Summary        string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	TrackingActive bool
	SecsSpent      int
	Active         bool
	TaskTitle      string
	TaskDesc       string
}

type TaskLogEntry struct {
	ID          int
	TaskID      int
	TaskSummary string
	BeginTS     time.Time
	EndTS       time.Time
	SecsSpent   int
	Comment     *string
	Desc        *string
	TLTitle     string
	TLDesc      string
}

type ActiveTaskDetails struct {
	TaskID            int
	TaskSummary       string
	CurrentLogBeginTS time.Time
	CurrentLogComment *string
	CurrentLogDesc    *string
}

type TaskReportEntry struct {
	TaskID      int
	TaskSummary string
	NumEntries  int
	SecsSpent   int
}

func (t *Task) UpdateTitle() {
	var trackingIndicator string
	if t.TrackingActive {
		trackingIndicator = "⏲ "
	}

	t.TaskTitle = trackingIndicator + t.Summary
}

func (t *Task) UpdateDesc() {
	var timeSpent string

	if t.SecsSpent != 0 {
		timeSpent = "worked on for " + HumanizeDuration(t.SecsSpent)
	} else {
		timeSpent = "no time spent"
	}
	lastUpdated := fmt.Sprintf("last updated: %s", humanize.Time(t.UpdatedAt))

	t.TaskDesc = fmt.Sprintf("%s %s", utils.RightPadTrim(lastUpdated, 60, true), timeSpent)
}

func (tl *TaskLogEntry) UpdateTitle() {
	tl.TLTitle = utils.Trim(tl.GetComment(), 60)
}

func (tl *TaskLogEntry) UpdateDesc() {
	timeSpentStr := HumanizeDuration(tl.SecsSpent)

	var timeStr string
	var durationMsg string

	endTSRelative := getTSRelative(tl.EndTS, time.Now())

	switch endTSRelative {
	case tsFromToday:
		durationMsg = fmt.Sprintf("%s  ...  %s", tl.BeginTS.Format(timeOnlyFormat), tl.EndTS.Format(timeOnlyFormat))
	case tsFromYesterday:
		durationMsg = "Yesterday"
	case tsFromThisWeek:
		durationMsg = tl.EndTS.Format(dayFormat)
	default:
		durationMsg = humanize.Time(tl.EndTS)
	}

	timeStr = fmt.Sprintf("%s (%s)",
		utils.RightPadTrim(durationMsg, 40, true),
		timeSpentStr)

	tl.TLDesc = fmt.Sprintf("%s %s", utils.RightPadTrim(tl.TaskSummary, 60, true), timeStr)
}

func (tl *TaskLogEntry) GetComment() string {
	if tl.Comment == nil {
		return "∅"
	}

	return *tl.Comment
}

func (tl *TaskLogEntry) GetDescription() string {
	if tl.Desc == nil {
		return "∅"
	}

	return *tl.Desc
}

func (tl *TaskLogEntry) GetDetails() string {
	timeSpentStr := HumanizeDuration(tl.SecsSpent)

	return fmt.Sprintf(`
Comment: %s

%s..%s (%s)

%s
`, tl.GetComment(), tl.BeginTS.Format(timeFormat), tl.EndTS.Format(timeFormat), timeSpentStr, tl.GetDescription())
}

func (t Task) Title() string {
	return t.TaskTitle
}

func (t Task) Description() string {
	return t.TaskDesc
}

func (t Task) FilterValue() string {
	return t.Summary
}

func (tl TaskLogEntry) Title() string {
	if tl.Desc != nil {
		return fmt.Sprintf("📃 %s", tl.TLTitle)
	}
	return tl.TLTitle
}

func (tl TaskLogEntry) Description() string {
	return tl.TLDesc
}

func (tl TaskLogEntry) FilterValue() string {
	return fmt.Sprintf("%d", tl.ID)
}

func HumanizeDuration(durationInSecs int) string {
	duration := time.Duration(durationInSecs) * time.Second

	if duration.Seconds() < 60 {
		return fmt.Sprintf("%ds", int(duration.Seconds()))
	}

	if duration.Minutes() < 60 {
		return fmt.Sprintf("%dm", int(duration.Minutes()))
	}

	modMins := int(math.Mod(duration.Minutes(), 60))

	if modMins == 0 {
		return fmt.Sprintf("%dh", int(duration.Hours()))
	}

	return fmt.Sprintf("%dh %dm", int(duration.Hours()), modMins)
}

type TimeShiftDirection uint8

const (
	ShiftForward TimeShiftDirection = iota
	ShiftBackward
)

type TimeShiftDuration uint8

const (
	ShiftMinute TimeShiftDuration = iota
	ShiftFiveMinutes
	ShiftHour
	ShiftDay
)
