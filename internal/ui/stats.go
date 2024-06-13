package ui

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/olekukonko/tablewriter"
)

const (
	statsLogEntriesLimit   = 10000
	statsNumDaysUpperBound = 3650
)

func RenderStats(db *sql.DB, writer io.Writer, plain bool, period string) {
	if period == "" {
		return
	}

	var start, end time.Time
	var entries []taskReportEntry
	var err error

	switch period {
	case "all":
		entries, err = fetchStatsFromDB(db, statsLogEntriesLimit)

	case "today":
		today := time.Now()

		start = time.Date(today.Year(),
			today.Month(),
			today.Day(),
			0,
			0,
			0,
			0,
			today.Location(),
		)
		end = start.AddDate(0, 0, 1)
		entries, err = fetchStatsBetweenTSFromDB(db, start, end, statsLogEntriesLimit)

	case "yest":
		yest := time.Now().AddDate(0, 0, -1)

		start = time.Date(yest.Year(),
			yest.Month(),
			yest.Day(),
			0,
			0,
			0,
			0,
			yest.Location(),
		)
		end = start.AddDate(0, 0, 1)
		entries, err = fetchStatsBetweenTSFromDB(db, start, end, statsLogEntriesLimit)

	case "3d":
		threeDaysAgo := time.Now().AddDate(0, 0, -2)

		start = time.Date(threeDaysAgo.Year(),
			threeDaysAgo.Month(),
			threeDaysAgo.Day(),
			0,
			0,
			0,
			0,
			threeDaysAgo.Location(),
		)
		end = time.Now()
		entries, err = fetchStatsBetweenTSFromDB(db, start, end, statsLogEntriesLimit)

	case "week":
		now := time.Now()
		weekday := now.Weekday()
		offset := (7 + weekday - time.Monday) % 7
		startOfWeek := now.AddDate(0, 0, -int(offset))
		start = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, startOfWeek.Location())

		end = time.Now()
		entries, err = fetchStatsBetweenTSFromDB(db, start, end, statsLogEntriesLimit)

	case "month":
		now := time.Now()

		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		end = time.Now()
		entries, err = fetchStatsBetweenTSFromDB(db, start, end, statsLogEntriesLimit)

	default:
		if strings.Contains(period, "...") {
			var ts timePeriod
			var nd int
			ts, nd, err = parseDateDuration(period)
			if err != nil {
				fmt.Fprintf(writer, "%s\n", err)
				os.Exit(1)
			}
			if nd > statsNumDaysUpperBound {
				fmt.Fprintf(writer, "Time period is too large; maximum number of days allowed in range (both inclusive): %d\n", statsNumDaysUpperBound)
				os.Exit(1)
			}
			start = ts.start
			end = ts.end.AddDate(0, 0, 1)
		} else {
			start, err = time.ParseInLocation(string(dateFormat), period, time.Local)
			if err != nil {
				fmt.Fprintf(writer, "Couldn't parse date: %s\n", err)
				os.Exit(1)
			}
			end = start.AddDate(0, 0, 1)
		}

		entries, err = fetchStatsBetweenTSFromDB(db, start, end, statsLogEntriesLimit)
	}

	if err != nil {
		fmt.Fprintf(writer, "Something went wrong generating the log: %s\n", err)
		os.Exit(1)
	}

	if len(entries) == 0 {
		return
	}
	renderStats(writer, plain, entries)
}

func renderStats(writer io.Writer, plain bool, entries []taskReportEntry) {

	data := make([][]string, len(entries))
	var timeSpentStr string

	rs := getReportStyles(plain)
	styleCache := make(map[string]lipgloss.Style)

	for i, entry := range entries {
		timeSpentStr = humanizeDuration(entry.secsSpent)

		if plain {
			data[i] = []string{
				Trim(entry.taskSummary, 50),
				fmt.Sprintf("%d", entry.numEntries),
				timeSpentStr,
			}
		} else {
			reportStyle, ok := styleCache[entry.taskSummary]
			if !ok {
				reportStyle = getDynamicStyle(entry.taskSummary)
				styleCache[entry.taskSummary] = reportStyle
			}
			data[i] = []string{
				reportStyle.Render(Trim(entry.taskSummary, 50)),
				reportStyle.Render(fmt.Sprintf("%d", entry.numEntries)),
				reportStyle.Render(timeSpentStr),
			}
		}
	}
	table := tablewriter.NewWriter(writer)

	headerValues := []string{"Task", "#LogEntries", "TimeSpent"}
	headers := make([]string, len(headerValues))
	for i, h := range headerValues {
		headers[i] = rs.headerStyle.Render(h)
	}
	table.SetHeader(headers)

	table.SetRowSeparator(rs.borderStyle.Render("-"))
	table.SetColumnSeparator(rs.borderStyle.Render("|"))
	table.SetCenterSeparator(rs.borderStyle.Render("+"))
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(false)
	table.AppendBulk(data)

	table.Render()
}
