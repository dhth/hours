package ui

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	timePeriodHoursUpperBound = 168
)

var (
	timePeriodNotValidErr = errors.New("time period is not valid")
	timePeriodTooLargeErr = fmt.Errorf("time period is too large; maximum allowed time period is %d hours", timePeriodHoursUpperBound)
)

type timePeriod struct {
	start time.Time
	end   time.Time
}

func parseDateDuration(period string) (timePeriod, error) {
	var tp timePeriod

	elements := strings.Split(period, "...")
	if len(elements) != 2 {
		return tp, timePeriodNotValidErr
	}

	start, err := time.ParseInLocation(string(dateFormat), elements[0], time.Local)
	if err != nil {
		return tp, timePeriodNotValidErr
	}

	end, err := time.ParseInLocation(string(dateFormat), elements[1], time.Local)
	if err != nil {
		return tp, timePeriodNotValidErr
	}

	if end.Sub(start) <= 0 {
		return tp, timePeriodNotValidErr
	}

	if end.Sub(start).Hours() >= timePeriodHoursUpperBound {
		return tp, timePeriodTooLargeErr
	}

	tp.start = start
	tp.end = end

	return tp, nil
}
