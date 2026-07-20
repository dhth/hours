package ui

import (
	"fmt"
	"time"

	"github.com/dhth/hours/internal/domain"
	"github.com/dhth/hours/internal/types"
	"github.com/dhth/hours/internal/utils"
	"github.com/dustin/go-humanize"
)

const (
	emptyCommentIndicator = "∅"
	dayFormat             = "Monday"
)

type taskLogListItem struct {
	domain.TaskLogEntry
	listTitle string
	listDesc  string
}

func (tl *taskLogListItem) updateListTitle() {
	tl.listTitle = utils.TrimWithMoreLinesIndicator(taskLogComment(tl.Comment), 60)
}

func (tl *taskLogListItem) updateListDesc(timeProvider types.TimeProvider) {
	timeSpentStr := types.HumanizeDuration(tl.SecsSpent)

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

	timeStr := fmt.Sprintf("%s (%s)",
		utils.RightPadTrim(durationMsg, 40, true),
		timeSpentStr)

	tl.listDesc = fmt.Sprintf("%s %s", utils.RightPadTrim(tl.TaskSummary, 60, true), timeStr)
}

func (tl taskLogListItem) Title() string {
	return tl.listTitle
}

func (tl taskLogListItem) Description() string {
	return tl.listDesc
}

func (tl taskLogListItem) FilterValue() string {
	return fmt.Sprintf("%d", tl.ID)
}

func taskLogComment(comment *string) string {
	if comment == nil {
		return emptyCommentIndicator
	}

	return *comment
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
