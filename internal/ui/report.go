package ui

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	"github.com/olekukonko/tablewriter"
)

func RenderNDaysReport(db *sql.DB, writer io.Writer, numDays int, plain bool) {
	now := time.Now()

	nDaysBack := now.AddDate(0, 0, -1*(numDays-1))

	start := time.Date(nDaysBack.Year(),
		nDaysBack.Month(),
		nDaysBack.Day(),
		0,
		0,
		0,
		0,
		nDaysBack.Location(),
	)

	day := start
	var nextDay time.Time

	var maxEntryForADay int
	reportData := make(map[int][]taskLogEntry)

	for i := 0; i < numDays; i++ {
		nextDay = day.AddDate(0, 0, 1)
		taskLogEntries, err := fetchTLEntriesBetweenTSFromDB(db, day, nextDay, 100)
		if err != nil {
			fmt.Fprintf(writer, "Something went wrong generating the report:\n%s", err)
			os.Exit(1)
		}

		day = nextDay
		reportData[i] = taskLogEntries
		if len(taskLogEntries) > maxEntryForADay {
			maxEntryForADay = len(taskLogEntries)
		}
	}
	data := make([][]string, maxEntryForADay)
	totalSecsPerDay := make(map[int]int)

	for j := 0; j < numDays; j++ {
		totalSecsPerDay[j] = 0
	}

	rs := getReportStyles(plain)

	styleCache := make(map[string]lipgloss.Style)
	for rowIndex := 0; rowIndex < maxEntryForADay; rowIndex++ {
		row := make([]string, numDays)
		for colIndex := 0; colIndex < numDays; colIndex++ {
			if rowIndex >= len(reportData[colIndex]) {
				row[colIndex] = ""
				continue
			}

			tr := reportData[colIndex][rowIndex]
			timeSpentStr := humanizeDuration(tr.secsSpent)

			if plain {
				row[colIndex] = fmt.Sprintf("%s %s",
					RightPadTrim(tr.taskSummary, 8, false),
					timeSpentStr,
				)
			} else {
				reportStyle, ok := styleCache[tr.taskSummary]

				if !ok {
					reportStyle = getDynamicStyle(tr.taskSummary)
					styleCache[tr.taskSummary] = reportStyle
				}

				row[colIndex] = fmt.Sprintf("%s %s",
					reportStyle.Render(RightPadTrim(tr.taskSummary, 8, false)),
					reportStyle.Render(timeSpentStr),
				)
			}
			totalSecsPerDay[colIndex] += tr.secsSpent
		}
		data[rowIndex] = row
	}

	totalTimePerDay := make([]string, numDays)

	for i, ts := range totalSecsPerDay {
		if ts != 0 {
			totalTimePerDay[i] = rs.footerStyle.Render(humanizeDuration(ts))
		} else {
			totalTimePerDay[i] = " "
		}
	}
	table := tablewriter.NewWriter(writer)

	headersValues := make([]string, numDays)

	day = start
	counter := 0

	if numDays > 2 {
		for counter < numDays-2 {
			headersValues[counter] = day.Weekday().String()
			day = day.AddDate(0, 0, 1)
			counter++
		}
	}
	if numDays > 1 {
		headersValues[counter] = "Yesterday"
		counter++
	}
	headersValues[counter] = "Today"

	headers := make([]string, numDays)
	for i := 0; i < numDays; i++ {
		headers[i] = rs.headerStyle.Render(headersValues[i])
	}

	table.SetHeader(headers)

	table.SetRowSeparator(rs.borderStyle.Render("-"))
	table.SetColumnSeparator(rs.borderStyle.Render("|"))
	table.SetCenterSeparator(rs.borderStyle.Render("+"))
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(false)
	table.AppendBulk(data)
	table.SetFooter(totalTimePerDay)

	table.Render()
}

func RenderNDaysReportAgg(db *sql.DB, writer io.Writer, numDays int, plain bool) {
	now := time.Now()

	nDaysBack := now.AddDate(0, 0, -1*(numDays-1))

	start := time.Date(nDaysBack.Year(),
		nDaysBack.Month(),
		nDaysBack.Day(),
		0,
		0,
		0,
		0,
		nDaysBack.Location(),
	)

	day := start
	var nextDay time.Time

	var maxEntryForADay int
	reportData := make(map[int][]taskReportEntry)

	for i := 0; i < numDays; i++ {
		nextDay = day.AddDate(0, 0, 1)
		taskLogEntries, err := fetchReportBetweenTSFromDB(db, day, nextDay, 100)
		if err != nil {
			fmt.Fprintf(writer, "Something went wrong generating the report:\n%s", err)
			os.Exit(1)
		}

		day = nextDay
		reportData[i] = taskLogEntries
		if len(taskLogEntries) > maxEntryForADay {
			maxEntryForADay = len(taskLogEntries)
		}
	}
	data := make([][]string, maxEntryForADay)
	totalSecsPerDay := make(map[int]int)

	for j := 0; j < numDays; j++ {
		totalSecsPerDay[j] = 0
	}

	rs := getReportStyles(plain)

	styleCache := make(map[string]lipgloss.Style)
	for rowIndex := 0; rowIndex < maxEntryForADay; rowIndex++ {
		row := make([]string, numDays)
		for colIndex := 0; colIndex < numDays; colIndex++ {
			if rowIndex >= len(reportData[colIndex]) {
				row[colIndex] = ""
				continue
			}

			tr := reportData[colIndex][rowIndex]
			timeSpentStr := humanizeDuration(tr.secsSpent)

			if plain {
				row[colIndex] = fmt.Sprintf("%s %s",
					RightPadTrim(tr.taskSummary, 8, false),
					timeSpentStr,
				)
			} else {
				reportStyle, ok := styleCache[tr.taskSummary]
				if !ok {
					reportStyle = getDynamicStyle(tr.taskSummary)
					styleCache[tr.taskSummary] = reportStyle
				}

				row[colIndex] = fmt.Sprintf("%s %s",
					reportStyle.Render(RightPadTrim(tr.taskSummary, 8, false)),
					reportStyle.Render(timeSpentStr),
				)
			}
			totalSecsPerDay[colIndex] += tr.secsSpent
		}
		data[rowIndex] = row
	}
	totalTimePerDay := make([]string, numDays)
	for i, ts := range totalSecsPerDay {
		if ts != 0 {
			totalTimePerDay[i] = rs.footerStyle.Render(humanizeDuration(ts))
		} else {
			totalTimePerDay[i] = " "
		}
	}
	table := tablewriter.NewWriter(writer)

	headersValues := make([]string, numDays)

	day = start
	counter := 0

	if numDays > 2 {
		for counter < numDays-2 {
			headersValues[counter] = day.Weekday().String()
			day = day.AddDate(0, 0, 1)
			counter++
		}
	}

	if numDays > 2 {
		headersValues[counter] = "Yesterday"
		counter++
	}
	headersValues[counter] = "Today"

	headers := make([]string, numDays)
	for i := 0; i < numDays; i++ {
		headers[i] = rs.headerStyle.Render(headersValues[i])
	}

	table.SetHeader(headers)

	table.SetRowSeparator(rs.borderStyle.Render("-"))
	table.SetColumnSeparator(rs.borderStyle.Render("|"))
	table.SetCenterSeparator(rs.borderStyle.Render("+"))
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(false)
	table.AppendBulk(data)
	table.SetFooter(totalTimePerDay)

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

	styleCache := make(map[string]lipgloss.Style)

	for i, entry := range tasks {
		reportStyle, ok := styleCache[entry.summary]
		if !ok {
			reportStyle = getDynamicStyle(entry.summary)
			styleCache[entry.summary] = reportStyle
		}

		timeSpentStr = humanizeDuration(entry.secsSpent)
		data[i] = []string{
			reportStyle.Render(fmt.Sprintf("%d", i+1)),
			reportStyle.Render(Trim(entry.summary, 50)),
			reportStyle.Render(timeSpentStr),
			reportStyle.Render(humanize.Time(entry.updatedAt)),
		}
	}

	table := tablewriter.NewWriter(writer)

	headerValues := []string{"#", "Task", "TimeSpent", "LastUpdated"}
	headers := make([]string, len(headerValues))
	for i, h := range headerValues {
		headers[i] = reportHeaderStyle.Render(h)
	}
	table.SetHeader(headers)

	table.SetRowSeparator(reportBorderStyle.Render("-"))
	table.SetColumnSeparator(reportBorderStyle.Render("|"))
	table.SetCenterSeparator(reportBorderStyle.Render("+"))
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(false)
	table.AppendBulk(data)

	table.Render()
}
