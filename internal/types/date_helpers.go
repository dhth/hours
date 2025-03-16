package types

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	dateRangeDaysUpperBound = 7
	TimePeriodWeek          = "week"
	timeFormat              = "2006/01/02 15:04"
	timeOnlyFormat          = "15:04"
	dayFormat               = "Monday"
	friendlyTimeFormat      = "Mon, 15:04"
	dateFormat              = "2006/01/02"
)

var (
	errDateRangeIncorrect         = errors.New("date range is incorrect")
	errStartDateIncorrect         = errors.New("start date is incorrect")
	errEndDateIncorrect           = errors.New("end date is incorrect")
	errEndDateIsNotAfterStartDate = errors.New("end date is not after start date")
	errTimePeriodNotValid         = errors.New("time period is not valid")
	errTimePeriodTooLarge         = errors.New("time period is too large")
)

func parseDateDuration(dateRangeStr string) (DateRange, error) {
	var dr DateRange

	elements := strings.Split(dateRangeStr, "...")
	if len(elements) != 2 {
		return dr, fmt.Errorf("%w: date range needs to be of the format: %s...%s", errDateRangeIncorrect, dateFormat, dateFormat)
	}

	start, err := time.ParseInLocation(string(dateFormat), elements[0], time.Local)
	if err != nil {
		return dr, fmt.Errorf("%w: %s", errStartDateIncorrect, err.Error())
	}

	end, err := time.ParseInLocation(string(dateFormat), elements[1], time.Local)
	if err != nil {
		return dr, fmt.Errorf("%w: %s", errEndDateIncorrect, err.Error())
	}

	if end.Sub(start) <= 0 {
		return dr, fmt.Errorf("%w", errEndDateIsNotAfterStartDate)
	}

	dr.Start = start
	dr.End = end
	dr.NumDays = int(end.Sub(start).Hours()/24) + 1

	return dr, nil
}

func GetDateRange(period string, now time.Time, fullWeek bool) (DateRange, error) {
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

	case TimePeriodWeek:
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
			var ts DateRange
			ts, err = parseDateDuration(period)
			if err != nil {
				return ts, fmt.Errorf("%w: %s", errTimePeriodNotValid, err.Error())
			}

			if ts.NumDays > dateRangeDaysUpperBound {
				return ts, fmt.Errorf("%w: maximum number of days allowed (both inclusive): %d", errTimePeriodTooLarge, dateRangeDaysUpperBound)
			}

			start = ts.Start
			end = ts.End.AddDate(0, 0, 1)
			numDays = ts.NumDays
		} else {
			start, err = time.ParseInLocation(string(dateFormat), period, time.Local)
			if err != nil {
				return DateRange{}, fmt.Errorf("%w: %s", errTimePeriodNotValid, err.Error())
			}
			end = start.AddDate(0, 0, 1)
			numDays = 1
		}
	}

	return DateRange{
		Start:   start,
		End:     end,
		NumDays: numDays,
	}, nil
}

func GetShiftedTime(ts time.Time, direction TimeShiftDirection, duration TimeShiftDuration) time.Time {
	var d time.Duration

	switch duration {
	case ShiftMinute:
		d = time.Minute
	case ShiftFiveMinutes:
		d = time.Minute * 5
	case ShiftHour:
		d = time.Hour
	case ShiftDay:
		d = time.Hour * 24
	}

	if direction == ShiftBackward {
		d = -1 * d
	}
	return ts.Add(d)
}

type tsRelative uint8

const (
	tsFromFuture tsRelative = iota
	tsFromToday
	tsFromYesterday
	tsFromThisWeek
	tsFromBeforeThisWeek
)

func getTSRelative(ts time.Time, reference time.Time) tsRelative {
	if ts.Sub(reference) > 0 {
		return tsFromFuture
	}

	startOfReferenceDay := time.Date(reference.Year(), reference.Month(), reference.Day(), 0, 0, 0, 0, reference.Location())

	if ts.Sub(startOfReferenceDay) > 0 {
		return tsFromToday
	}

	startOfYest := startOfReferenceDay.AddDate(0, 0, -1)

	if ts.Sub(startOfYest) > 0 {
		return tsFromYesterday
	}

	weekday := reference.Weekday()
	offset := (7 + weekday - time.Monday) % 7
	startOfWeek := startOfReferenceDay.AddDate(0, 0, -int(offset))

	if ts.Sub(startOfWeek) > 0 {
		return tsFromThisWeek
	}
	return tsFromBeforeThisWeek
}
