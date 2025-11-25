package types

import (
	"errors"
	"strings"
	"time"
)

var (
	errBeginTimeIsEmpty       = errors.New("begin time is empty")
	errEndTimeIsEmpty         = errors.New("end time is empty")
	errBeginTimeIsInvalid     = errors.New("begin time is invalid")
	errEndTimeIsInvalid       = errors.New("end time is invalid")
	errEndTimeBeforeBeginTime = errors.New("end time is before begin time")
	ErrDurationNotLongEnough  = errors.New("end time needs to be at least a minute after begin time")
)

func ParseTaskLogTimes(beginStr, endStr string) (time.Time, time.Time, error) {
	var zero time.Time
	if strings.TrimSpace(beginStr) == "" {
		return zero, zero, errBeginTimeIsEmpty
	}

	if strings.TrimSpace(endStr) == "" {
		return zero, zero, errEndTimeIsEmpty
	}

	beginTS, err := time.ParseInLocation(timeFormat, beginStr, time.Local)
	if err != nil {
		return zero, zero, errBeginTimeIsInvalid
	}

	endTS, err := time.ParseInLocation(timeFormat, endStr, time.Local)
	if err != nil {
		return zero, zero, errEndTimeIsInvalid
	}

	durationErr := IsTaskLogDurationValid(beginTS, endTS)
	if durationErr != nil {
		return zero, zero, durationErr
	}

	return beginTS, endTS, nil
}

func IsTaskLogDurationValid(begin, end time.Time) error {
	if end.Before(begin) {
		return errEndTimeBeforeBeginTime
	}

	if end.Sub(begin) < time.Minute {
		return ErrDurationNotLongEnough
	}
	return nil
}
