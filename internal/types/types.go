package types

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/dhth/hours/internal/utils"
	"github.com/dustin/go-humanize"
)

const emptyCommentIndicator = "∅"

var ErrIncorrectTaskStatusProvided = errors.New("incorrect task status provided")

type TaskLogEntry struct {
	ID          int
	TaskID      int
	TaskSummary string
	BeginTS     time.Time
	EndTS       time.Time
	SecsSpent   int
	Comment     *string
	ListTitle   string
	ListDesc    string
}

type ActiveTaskLogEntry struct {
	ID          int
	TaskID      int
	TaskSummary string
	BeginTS     time.Time
	Comment     *string
}

type ActiveTaskDetails struct {
	TaskID            int
	TaskSummary       string
	CurrentLogBeginTS time.Time
	CurrentLogComment *string
}

type TaskReportEntry struct {
	TaskID      int
	TaskSummary string
	NumEntries  int
	SecsSpent   int
}

type TimeProvider interface {
	Now() time.Time
}

type RealTimeProvider struct{}

func (RealTimeProvider) Now() time.Time {
	return time.Now()
}

type TestTimeProvider struct {
	FixedTime time.Time
}

func (t TestTimeProvider) Now() time.Time {
	return t.FixedTime
}

func (tl *TaskLogEntry) UpdateListTitle() {
	tl.ListTitle = utils.TrimWithMoreLinesIndicator(tl.GetComment(), 60)
}

func (tl *TaskLogEntry) UpdateListDesc(timeProvider TimeProvider) {
	timeSpentStr := HumanizeDuration(tl.SecsSpent)

	var timeStr string
	var durationMsg string

	now := timeProvider.Now()
	endTSRelative := getTSRelative(tl.EndTS, now)

	switch endTSRelative {
	case tsFromToday:
		durationMsg = fmt.Sprintf("%s  ...  %s", tl.BeginTS.Format(timeOnlyFormat), tl.EndTS.Format(timeOnlyFormat))
	case tsFromYesterday:
		durationMsg = "Yesterday"
	case tsFromThisWeek:
		durationMsg = tl.EndTS.Format(dayFormat)
	default:
		durationMsg = humanize.RelTime(tl.EndTS, now, "ago", "from now")
	}

	timeStr = fmt.Sprintf("%s (%s)",
		utils.RightPadTrim(durationMsg, 40, true),
		timeSpentStr)

	tl.ListDesc = fmt.Sprintf("%s %s", utils.RightPadTrim(tl.TaskSummary, 60, true), timeStr)
}

func (tl *TaskLogEntry) GetComment() string {
	if tl.Comment == nil {
		return emptyCommentIndicator
	}

	return *tl.Comment
}

func (tl TaskLogEntry) Title() string {
	return tl.ListTitle
}

func (tl TaskLogEntry) Description() string {
	return tl.ListDesc
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

type TaskStatus uint8

const (
	TSValueActive   = "active"
	TSValueInactive = "inactive"
	TSValueAny      = "any"
)

const (
	TaskStatusActive TaskStatus = iota
	TaskStatusInactive
	TaskStatusAny
)

func ParseTaskStatus(value string) (TaskStatus, error) {
	switch value {
	case TSValueActive:
		return TaskStatusActive, nil
	case TSValueInactive:
		return TaskStatusInactive, nil
	case TSValueAny:
		return TaskStatusAny, nil
	default:
		return TaskStatusAny, ErrIncorrectTaskStatusProvided
	}
}

var ValidTaskStatusValues = []string{TSValueActive, TSValueInactive, TSValueAny}

type DateRange struct {
	Start   time.Time
	End     time.Time
	NumDays int
}
