package ui

import (
	"database/sql"
	"fmt"
	"io"
	"os"

	"github.com/dustin/go-humanize"
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
			fmt.Sprintf("%d", i+1),
			Trim(entry.taskSummary, 50),
			Trim(entry.comment, 50),
			entry.beginTS.Format(timeFormat),
			timeSpentStr,
		}
	}

	table := tablewriter.NewWriter(writer)
	table.SetHeader([]string{"#", "Task", "Log Comment", "Begin TS", "Time Spent"})

	table.SetAutoWrapText(false)
	table.AppendBulk(data)
	table.Render()
}

func RenderTaskReport(db *sql.DB, writer io.Writer) {
	tasks, err := fetchTasksFromDB(db, true, 30)
	if err != nil {
		fmt.Fprintf(writer, "Something went wrong generating the report:\n%s", err)
		os.Exit(1)
	}

	if len(tasks) == 0 {
		return
	}

	data := make([][]string, len(tasks))
	var timeSpentStr string

	for i, entry := range tasks {
		timeSpentStr = humanizeDuration(entry.secsSpent)
		data[i] = []string{
			fmt.Sprintf("%d", i+1),
			Trim(entry.summary, 50),
			timeSpentStr,
			humanize.Time(entry.updatedAt),
		}
	}

	table := tablewriter.NewWriter(writer)
	table.SetHeader([]string{"#", "Task", "Time Spent", "Last Updated"})

	table.SetAutoWrapText(false)
	table.AppendBulk(data)
	table.Render()
}
