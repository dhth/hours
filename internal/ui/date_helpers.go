package ui

import (
	"errors"
	"strings"
	"time"
)

var (
	timePeriodNotValidErr = errors.New("time period is not valid")
)

type timePeriod struct {
	start time.Time
	end   time.Time
}

func parseDateDuration(period string) (timePeriod, int, error) {
	var tp timePeriod
	var numDaysBothInclusive int

	elements := strings.Split(period, "...")
	if len(elements) != 2 {
		return tp, numDaysBothInclusive, timePeriodNotValidErr
	}

	start, err := time.ParseInLocation(string(dateFormat), elements[0], time.Local)
	if err != nil {
		return tp, numDaysBothInclusive, timePeriodNotValidErr
	}

	end, err := time.ParseInLocation(string(dateFormat), elements[1], time.Local)
	if err != nil {
		return tp, numDaysBothInclusive, timePeriodNotValidErr
	}

	if end.Sub(start) <= 0 {
		return tp, numDaysBothInclusive, timePeriodNotValidErr
	}

	tp.start = start
	tp.end = end
	numDaysBothInclusive = int(end.Sub(start).Hours()/24) + 1

	return tp, numDaysBothInclusive, nil
}
