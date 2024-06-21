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
	statsLogEntriesLimit   = 10000
	statsNumDaysUpperBound = 3650
	statsTimeCharsBudget   = 6
)

func RenderStats(db *sql.DB, writer io.Writer, plain bool, period string, interactive bool) {
	if period == "" {
		return
	}

	var stats string
	var err error

	if interactive && period == "all" {
		fmt.Print("Interactive mode cannot be used when period='all'\n")
		os.Exit(1)
	}

	if period == "all" {
		// TODO: find a better way for this, passing start, end for "all" doesn't make sense
		stats, err = renderStats(db, period, time.Now(), time.Now(), plain)
		if err != nil {
			fmt.Fprintf(writer, "Something went wrong generating the log: %s\n", err)
			os.Exit(1)
		}
		fmt.Fprint(writer, stats)
		return
	}

	var fullWeek bool
	if interactive {
		fullWeek = true
	}
	ts, tsErr := getTimePeriod(period, time.Now(), fullWeek)

	if tsErr != nil {
		fmt.Printf("error: %s\n", tsErr)
		os.Exit(1)
	}
	stats, err = renderStats(db, period, ts.start, ts.end, plain)

	if err != nil {
		fmt.Fprintf(writer, "Something went wrong generating the log: %s\n", err)
		os.Exit(1)
	}

	if interactive {
		p := tea.NewProgram(initialRecordsModel(reportStats, db, ts.start, ts.end, plain, period, ts.numDays, stats))
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there has been an error: %v", err)
			os.Exit(1)
		}
	} else {
		fmt.Fprint(writer, stats)
	}
}

func renderStats(db *sql.DB, period string, start, end time.Time, plain bool) (string, error) {
	var entries []taskReportEntry
	var err error

	if period == "all" {
		entries, err = fetchStatsFromDB(db, statsLogEntriesLimit)
	} else {
		entries, err = fetchStatsBetweenTSFromDB(db, start, end, statsLogEntriesLimit)
	}

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
			"",
			RightPadTrim("", statsTimeCharsBudget, false),
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
				fmt.Sprintf("%d", entry.numEntries),
				RightPadTrim("", statsTimeCharsBudget, false),
			}
		} else {
			rowStyle, ok := styleCache[entry.taskSummary]
			if !ok {
				rowStyle = getDynamicStyle(entry.taskSummary)
				styleCache[entry.taskSummary] = rowStyle
			}
			data[i] = []string{
				rowStyle.Render(RightPadTrim(entry.taskSummary, 20, false)),
				rowStyle.Render(fmt.Sprintf("%d", entry.numEntries)),
				rowStyle.Render(RightPadTrim(timeSpentStr, statsTimeCharsBudget, false)),
			}
		}
	}
	b := bytes.Buffer{}
	table := tablewriter.NewWriter(&b)

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

	return b.String(), nil
}
