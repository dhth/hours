package types

import (
	"errors"
	"fmt"
	"math"
	"time"
)

var ErrIncorrectTaskStatusProvided = errors.New("incorrect task status provided")

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
