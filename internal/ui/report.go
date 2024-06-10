package ui

import (
	"database/sql"
	"fmt"
	"io"
	"os"

	"github.com/olekukonko/tablewriter"
)

func RenderTaskLogReport(db *sql.DB, writer io.Writer) {
	taskLogEntries, err := fetchTLEntriesFromDB(db, 20)
	if err != nil {
		fmt.Fprintf(writer, "Something went wrong generating the report:\n%s", err)
		os.Exit(1)
	}

	if len(taskLogEntries) == 0 {
		return
	}

	data := make([][]string, len(taskLogEntries))
	var secsSpent int
	var timeSpentStr string

	for i, entry := range taskLogEntries {
		secsSpent = int(entry.endTS.Sub(entry.beginTS).Seconds())
		timeSpentStr = humanizeDuration(secsSpent)
		data[i] = []string{
			fmt.Sprintf("%d", entry.id),
			entry.taskSummary,
			entry.beginTS.Format(timeFormat),
			timeSpentStr,
		}
	}

	table := tablewriter.NewWriter(writer)
	table.SetHeader([]string{"ID", "Task", "Begin TS", "Time Spent"})

	table.AppendBulk(data)
	table.Render()
}
