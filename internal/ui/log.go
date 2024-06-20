package ui

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/olekukonko/tablewriter"
)

const (
	logNumDaysUpperBound = 7
	logTimeCharsBudget   = 6
)

func RenderTaskLog(db *sql.DB, writer io.Writer, plain bool, period string, interactive bool) {
	if period == "" {
		return
	}

	ts, err := getTimePeriod(period, time.Now(), false)

	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}

	if interactive && ts.numDays > 1 {
		fmt.Print("Interactive mode for logs is limited to a day; use non-interactive mode to see logs for a larger time period\n")
		os.Exit(1)
	}

	log, err := renderTaskLog(db, ts.start, ts.end, 100, plain)

	if interactive {
		p := tea.NewProgram(initialRecordsModel(reportLogs, db, ts.start, ts.end, plain, period, ts.numDays, log))
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there has been an error: %v", err)
			os.Exit(1)
		}
	} else {
		fmt.Fprint(writer, log)
	}
}

func renderTaskLog(db *sql.DB, start, end time.Time, limit int, plain bool) (string, error) {

	entries, err := fetchTLEntriesBetweenTSFromDB(db, start, end, limit)

	if err != nil {
		return "", err
	}

	var numEntriesInTable int

	if len(entries) == 0 {
		numEntriesInTable = 1
	} else {
		numEntriesInTable = len(entries)
	}

	data := make([][]string, numEntriesInTable)

	if len(entries) == 0 {
		data[0] = []string{
			RightPadTrim("", 20, false),
			RightPadTrim("", 40, false),
			RightPadTrim("", 39, false),
			RightPadTrim("", logTimeCharsBudget, false),
		}
	}

	var timeSpentStr string

	rs := getReportStyles(plain)
	styleCache := make(map[string]lipgloss.Style)

	for i, entry := range entries {
		timeSpentStr = humanizeDuration(entry.secsSpent)

		if plain {
			data[i] = []string{
				RightPadTrim(entry.taskSummary, 20, false),
				RightPadTrim(entry.comment, 40, false),
				fmt.Sprintf("%s  ...  %s", entry.beginTs.Format(timeFormat), entry.beginTs.Format(timeFormat)),
				RightPadTrim(timeSpentStr, logTimeCharsBudget, false),
			}
		} else {
			reportStyle, ok := styleCache[entry.taskSummary]
			if !ok {
				reportStyle = getDynamicStyle(entry.taskSummary)
				styleCache[entry.taskSummary] = reportStyle
			}
			data[i] = []string{
				reportStyle.Render(RightPadTrim(entry.taskSummary, 20, false)),
				reportStyle.Render(RightPadTrim(entry.comment, 40, false)),
				reportStyle.Render(fmt.Sprintf("%s  ...  %s", entry.beginTs.Format(timeFormat), entry.endTs.Format(timeFormat))),
				reportStyle.Render(RightPadTrim(timeSpentStr, logTimeCharsBudget, false)),
			}
		}
	}

	b := bytes.Buffer{}
	table := tablewriter.NewWriter(&b)

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

	return b.String(), nil
}
