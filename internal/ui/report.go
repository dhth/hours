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

	flip := 0
	var reportStyle lipgloss.Style
	for i, entry := range taskLogEntries {
		if flip == 0 {
			reportStyle = reportStyle1
			flip = 1
		} else {
			reportStyle = reportStyle2
			flip = 0
		}
		secsSpent = int(entry.endTS.Sub(entry.beginTS).Seconds())
		timeSpentStr = humanizeDuration(secsSpent)
		data[i] = []string{
			reportStyle.Render(fmt.Sprintf("%d", i+1)),
			reportStyle.Render(Trim(entry.taskSummary, 50)),
			reportStyle.Render(Trim(entry.comment, 50)),
			reportStyle.Render(entry.beginTS.Format(timeFormat)),
			reportStyle.Render(timeSpentStr),
		}
	}
	table := tablewriter.NewWriter(writer)

	headerValues := []string{"#", "Task", "Comment", "Begin", "TimeSpent"}
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

var (
	daysMap = map[int]string{
		0: "Mon",
		1: "Tues",
		2: "Wed",
		3: "Thurs",
		4: "Fri",
		5: "Sat",
		6: "Sun",
	}
)

func RenderWeeklyReport(db *sql.DB, writer io.Writer) {
	numDays := 7
	now := time.Now().Local()

	// say time.Now is Sunday 16:00
	// start:     mon
	// day2:      tues
	// day3:      wed
	// day4:      thurs
	// day5:      fri
	// yesterday: sat
	// today:     sun
	sevenDaysFromNow := now.AddDate(0, 0, -7)

	start := time.Date(sevenDaysFromNow.Year(),
		sevenDaysFromNow.Month(),
		sevenDaysFromNow.Day(),
		0,
		0,
		0,
		0,
		sevenDaysFromNow.Location(),
	)

	day1 := start.AddDate(0, 0, 1)
	day2 := start.AddDate(0, 0, 2)
	day3 := start.AddDate(0, 0, 3)
	day4 := start.AddDate(0, 0, 4)
	day5 := start.AddDate(0, 0, 5)
	yesterday := start.AddDate(0, 0, 6)
	today := start.AddDate(0, 0, 7)

	taskLogEntries, err := fetchTLEntriesFromDBAfterTS(db, day1, 100)
	if err != nil {
		fmt.Fprintf(writer, "Something went wrong generating the report:\n%s", err)
		os.Exit(1)
	}

	if len(taskLogEntries) == 0 {
		return
	}

	dayEntries := make(map[int][]*taskLogEntry)

	for _, entry := range taskLogEntries {
		switch {
		case entry.beginTS.Before(day2):
			dayEntries[0] = append(dayEntries[0], &entry)
		case entry.beginTS.Before(day3):
			dayEntries[1] = append(dayEntries[1], &entry)
		case entry.beginTS.Before(day4):
			dayEntries[2] = append(dayEntries[2], &entry)
		case entry.beginTS.Before(day5):
			dayEntries[3] = append(dayEntries[3], &entry)
		case entry.beginTS.Before(yesterday):
			dayEntries[4] = append(dayEntries[4], &entry)
		case entry.beginTS.Before(today):
			dayEntries[5] = append(dayEntries[5], &entry)
		default:
			dayEntries[6] = append(dayEntries[6], &entry)
		}
	}
	var maxEntryForADay int
	for _, entries := range dayEntries {
		if len(entries) > maxEntryForADay {
			maxEntryForADay = len(entries)
		}
	}

	data := make([][]string, len(taskLogEntries))

	totalSecsPerDay := make(map[int]int)

	for j := 0; j < numDays; j++ {
		totalSecsPerDay[j] = 0
	}

	for rowIndex := 0; rowIndex < maxEntryForADay; rowIndex++ {
		row := make([]string, numDays)
		for colIndex := 0; colIndex < numDays; colIndex++ {
			if rowIndex >= len(dayEntries[colIndex]) {
				row[colIndex] = ""
				continue
			}

			tl := dayEntries[colIndex][rowIndex]
			secsSpent := int(tl.endTS.Sub(tl.beginTS).Seconds())
			timeSpent := humanizeDuration(secsSpent)
			row[colIndex] = fmt.Sprintf("%s %s",
				reportStyle1.Render(RightPadTrim(tl.taskSummary, 8, false)),
				reportStyle3.Render(fmt.Sprintf("(%s)", timeSpent)),
			)
			totalSecsPerDay[colIndex] += secsSpent
		}
		data[rowIndex] = row
	}

	totalTimePerDay := make([]string, numDays)
	for i, ts := range totalSecsPerDay {
		if ts != 0 {
			totalTimePerDay[i] = reportFooterStyle.Render(fmt.Sprintf("%s", humanizeDuration(ts)))
		} else {
			totalTimePerDay[i] = " "
		}
	}

	table := tablewriter.NewWriter(writer)

	table.SetHeader([]string{
		reportHeaderStyle.Render(day1.Weekday().String()),
		reportHeaderStyle.Render(day2.Weekday().String()),
		reportHeaderStyle.Render(day3.Weekday().String()),
		reportHeaderStyle.Render(day4.Weekday().String()),
		reportHeaderStyle.Render(day5.Weekday().String()),
		reportHeaderStyle.Render("Yesterday"),
		reportHeaderStyle.Render("Today"),
	})

	table.SetRowSeparator(reportBorderStyle.Render("-"))
	table.SetColumnSeparator(reportBorderStyle.Render("|"))
	table.SetCenterSeparator(reportBorderStyle.Render("+"))
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

	var reportStyle lipgloss.Style

	flip := 0
	for i, entry := range tasks {
		if flip == 0 {
			reportStyle = reportStyle1
			flip = 1
		} else {
			reportStyle = reportStyle2
			flip = 0
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
