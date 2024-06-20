package ui

import (
	"fmt"
	"strings"
	"time"
)

const (
	timePeriodDaysUpperBound = 7
)

var (
	timePeriodNotValidErr = fmt.Errorf("time period is not valid; accepted format: %s or %s...%s", dateFormat, dateFormat, dateFormat)
	timePeriodTooLargeErr = fmt.Errorf("time period is too large; maximum number of days allowed (both inclusive): %d", timePeriodDaysUpperBound)
)

type timePeriod struct {
	start   time.Time
	end     time.Time
	numDays int
}

func parseDateDuration(dateRange string) (timePeriod, error) {
	var tp timePeriod

	elements := strings.Split(dateRange, "...")
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

	tp.start = start
	tp.end = end
	tp.numDays = int(end.Sub(start).Hours()/24) + 1

	return tp, nil
}

func getTimePeriod(period string, now time.Time, fullWeek bool) (timePeriod, error) {
	var start, end time.Time
	var numDays int

	switch period {

	case "today":
		start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		end = start.AddDate(0, 0, 1)
		numDays = 1

	case "yest":
		aDayBefore := now.AddDate(0, 0, -1)

		start = time.Date(aDayBefore.Year(), aDayBefore.Month(), aDayBefore.Day(), 0, 0, 0, 0, aDayBefore.Location())
		end = start.AddDate(0, 0, 1)
		numDays = 1

	case "3d":
		threeDaysBefore := now.AddDate(0, 0, -2)

		start = time.Date(threeDaysBefore.Year(), threeDaysBefore.Month(), threeDaysBefore.Day(), 0, 0, 0, 0, threeDaysBefore.Location())
		end = start.AddDate(0, 0, 3)
		numDays = 3

	case "week":
		weekday := now.Weekday()
		offset := (7 + weekday - time.Monday) % 7
		startOfWeek := now.AddDate(0, 0, -int(offset))
		start = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, startOfWeek.Location())
		if fullWeek {
			numDays = 7
		} else {
			numDays = int(offset) + 1
		}
		end = start.AddDate(0, 0, numDays)

	default:
		var err error

		if strings.Contains(period, "...") {
			var ts timePeriod
			ts, err = parseDateDuration(period)
			if err != nil {
				return ts, err
			}
			if ts.numDays > timePeriodDaysUpperBound {
				return ts, timePeriodTooLargeErr
			}

			start = ts.start
			end = ts.end.AddDate(0, 0, 1)
			numDays = ts.numDays
		} else {
			start, err = time.ParseInLocation(string(dateFormat), period, time.Local)
			if err != nil {
				return timePeriod{}, timePeriodNotValidErr
			}
			end = start.AddDate(0, 0, 1)
			numDays = 1
		}
	}

	return timePeriod{
		start:   start,
		end:     end,
		numDays: numDays,
	}, nil
}
