package ui

import (
	"database/sql"
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	"github.com/olekukonko/tablewriter"
)

func RenderTaskLog(db *sql.DB, writer io.Writer, plain bool) {
	taskLogEntries, err := fetchTLEntriesFromDB(db, true, 20)
	if err != nil {
		fmt.Fprintf(writer, "Something went wrong generating the log:\n%s", err)
		os.Exit(1)
	}

	if len(taskLogEntries) == 0 {
		return
	}

	data := make([][]string, len(taskLogEntries))
	var timeSpentStr string

	rs := getReportStyles(plain)
	styleCache := make(map[string]lipgloss.Style)

	for i, entry := range taskLogEntries {
		timeSpentStr = humanizeDuration(entry.secsSpent)

		if plain {
			data[i] = []string{
				Trim(entry.taskSummary, 50),
				Trim(entry.comment, 80),
				fmt.Sprintf("%s (%s)", entry.beginTS.Format(timeFormat), humanize.Time(entry.beginTS)),
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
				reportStyle.Render(fmt.Sprintf("%s (%s)", entry.beginTS.Format(timeFormat), humanize.Time(entry.beginTS))),
				reportStyle.Render(timeSpentStr),
			}
		}
	}
	table := tablewriter.NewWriter(writer)

	headerValues := []string{"Task", "Comment", "Begin", "TimeSpent"}
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
