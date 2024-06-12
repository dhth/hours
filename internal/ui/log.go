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
	logLimit = 100
)

func RenderTaskLog(db *sql.DB, writer io.Writer, plain bool, period string) {
	if period == "" {
		return
	}

	switch period {
	case "all":
		taskLogEntries, err := fetchTLEntriesFromDB(db, true, 20)
		if err != nil {
			fmt.Fprintf(writer, "Something went wrong generating the log: %s\n", err)
			os.Exit(1)
		}
		renderTaskLog(writer, plain, taskLogEntries)

	case "today":
		today := time.Now()

		start := time.Date(today.Year(),
			today.Month(),
			today.Day(),
			0,
			0,
			0,
			0,
			today.Location(),
		)
		taskLogEntries, err := fetchTLEntriesBetweenTSFromDB(db, start, start.AddDate(0, 0, 1), logLimit)
		if err != nil {
			fmt.Fprintf(writer, "Something went wrong generating the log: %s\n", err)
			os.Exit(1)
		}
		renderTaskLog(writer, plain, taskLogEntries)

	case "yest":
		yest := time.Now().AddDate(0, 0, -1)

		start := time.Date(yest.Year(),
			yest.Month(),
			yest.Day(),
			0,
			0,
			0,
			0,
			yest.Location(),
		)
		taskLogEntries, err := fetchTLEntriesBetweenTSFromDB(db, start, start.AddDate(0, 0, 1), logLimit)
		if err != nil {
			fmt.Fprintf(writer, "Something went wrong generating the log: %s\n", err)
			os.Exit(1)
		}
		renderTaskLog(writer, plain, taskLogEntries)

	case "3d":
		threeDaysAgo := time.Now().AddDate(0, 0, -2)

		start := time.Date(threeDaysAgo.Year(),
			threeDaysAgo.Month(),
			threeDaysAgo.Day(),
			0,
			0,
			0,
			0,
			threeDaysAgo.Location(),
		)
		taskLogEntries, err := fetchTLEntriesBetweenTSFromDB(db, start, time.Now(), logLimit)
		if err != nil {
			fmt.Fprintf(writer, "Something went wrong generating the log: %s\n", err)
			os.Exit(1)
		}
		renderTaskLog(writer, plain, taskLogEntries)

	case "week":
		aWeekBack := time.Now().AddDate(0, 0, -6)

		start := time.Date(aWeekBack.Year(),
			aWeekBack.Month(),
			aWeekBack.Day(),
			0,
			0,
			0,
			0,
			aWeekBack.Location(),
		)
		taskLogEntries, err := fetchTLEntriesBetweenTSFromDB(db, start, time.Now(), logLimit)
		if err != nil {
			fmt.Fprintf(writer, "Something went wrong generating the log: %s\n", err)
			os.Exit(1)
		}
		renderTaskLog(writer, plain, taskLogEntries)

	default:
		var start, end time.Time
		var err error

		if strings.Contains(period, "...") {
			var ts timePeriod
			ts, err = parseDateDuration(period)
			if err != nil {
				fmt.Fprintf(writer, "%s\n", err)
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

		taskLogEntries, err := fetchTLEntriesBetweenTSFromDB(db, start, end, logLimit)
		if err != nil {
			fmt.Fprintf(writer, "Something went wrong generating the log: %s\n", err)
			os.Exit(1)
		}
		renderTaskLog(writer, plain, taskLogEntries)
	}
}

func renderTaskLog(writer io.Writer, plain bool, entries []taskLogEntry) {

	if len(entries) == 0 {
		return
	}

	data := make([][]string, len(entries))
	var timeSpentStr string

	rs := getReportStyles(plain)
	styleCache := make(map[string]lipgloss.Style)

	for i, entry := range entries {
		timeSpentStr = humanizeDuration(entry.secsSpent)

		if plain {
			data[i] = []string{
				Trim(entry.taskSummary, 50),
				Trim(entry.comment, 80),
				fmt.Sprintf("%s  ...  %s", entry.beginTs.Format(timeFormat), entry.beginTs.Format(timeFormat)),
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
				reportStyle.Render(Trim(entry.comment, 80)),
				reportStyle.Render(fmt.Sprintf("%s  ...  %s", entry.beginTs.Format(timeFormat), entry.endTs.Format(timeFormat))),
				reportStyle.Render(timeSpentStr),
			}
		}
	}
	table := tablewriter.NewWriter(writer)

	headerValues := []string{"Task", "Comment", "Duration", "TimeSpent"}
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
