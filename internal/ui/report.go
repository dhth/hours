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
	reportTimeCharsBudget = 6
)

func RenderReport(db *sql.DB, writer io.Writer, plain bool, period string, agg bool, interactive bool) {
	if period == "" {
		return
	}

	var fullWeek bool
	if interactive {
		fullWeek = true
	}
	ts, err := getTimePeriod(period, time.Now(), fullWeek)

	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}

	var report string
	var analyticsType recordsType

	if agg {
		analyticsType = reportAggRecords
		report, err = getReportAgg(db, ts.start, ts.numDays, plain)
	} else {
		analyticsType = reportRecords
		report, err = getReport(db, ts.start, ts.numDays, plain)
	}
	if err != nil {
		fmt.Printf("Something went wrong generating the report: %s\n", err)
	}

	if interactive {
		p := tea.NewProgram(initialRecordsModel(analyticsType, db, ts.start, ts.end, plain, period, ts.numDays, report))
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there has been an error: %v", err)
			os.Exit(1)
		}
	} else {
		fmt.Fprint(writer, report)
	}
}

func getReport(db *sql.DB, start time.Time, numDays int, plain bool) (string, error) {
	day := start
	var nextDay time.Time

	var maxEntryForADay int
	reportData := make(map[int][]taskLogEntry)

	noEntriesFound := true
	for i := 0; i < numDays; i++ {
		nextDay = day.AddDate(0, 0, 1)
		taskLogEntries, err := fetchTLEntriesBetweenTSFromDB(db, day, nextDay, 100)
		if err != nil {
			return "", err
		}
		if noEntriesFound && len(taskLogEntries) > 0 {
			noEntriesFound = false
		}

		day = nextDay
		reportData[i] = taskLogEntries

		if len(taskLogEntries) > maxEntryForADay {
			maxEntryForADay = len(taskLogEntries)
		}
	}

	if noEntriesFound {
		maxEntryForADay = 1
	}

	data := make([][]string, maxEntryForADay)
	totalSecsPerDay := make(map[int]int)

	for j := 0; j < numDays; j++ {
		totalSecsPerDay[j] = 0
	}

	rs := getReportStyles(plain)

	var summaryBudget int
	switch numDays {
	case 7:
		summaryBudget = 8
	case 6:
		summaryBudget = 10
	case 5:
		summaryBudget = 14
	default:
		summaryBudget = 16
	}

	styleCache := make(map[string]lipgloss.Style)
	for rowIndex := 0; rowIndex < maxEntryForADay; rowIndex++ {
		row := make([]string, numDays)
		for colIndex := 0; colIndex < numDays; colIndex++ {
			if rowIndex >= len(reportData[colIndex]) {
				row[colIndex] = fmt.Sprintf("%s  %s",
					RightPadTrim("", summaryBudget, false),
					RightPadTrim("", reportTimeCharsBudget, false),
				)
				continue
			}

			tr := reportData[colIndex][rowIndex]
			timeSpentStr := humanizeDuration(tr.secsSpent)

			if plain {
				row[colIndex] = fmt.Sprintf("%s  %s",
					RightPadTrim(tr.taskSummary, summaryBudget, false),
					RightPadTrim(timeSpentStr, reportTimeCharsBudget, false),
				)
			} else {
				reportStyle, ok := styleCache[tr.taskSummary]

				if !ok {
					reportStyle = getDynamicStyle(tr.taskSummary)
					styleCache[tr.taskSummary] = reportStyle
				}

				row[colIndex] = fmt.Sprintf("%s  %s",
					reportStyle.Render(RightPadTrim(tr.taskSummary, summaryBudget, false)),
					reportStyle.Render(RightPadTrim(timeSpentStr, reportTimeCharsBudget, false)),
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

	b := bytes.Buffer{}
	table := tablewriter.NewWriter(&b)

	headersValues := make([]string, numDays)

	day = start
	counter := 0

	for counter < numDays {
		headersValues[counter] = day.Format(dateFormat)
		day = day.AddDate(0, 0, 1)
		counter++
	}

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

	return b.String(), nil
}

func getReportAgg(db *sql.DB, start time.Time, numDays int, plain bool) (string, error) {

	day := start
	var nextDay time.Time

	var maxEntryForADay int
	reportData := make(map[int][]taskReportEntry)

	noEntriesFound := true
	for i := 0; i < numDays; i++ {
		nextDay = day.AddDate(0, 0, 1)
		taskLogEntries, err := fetchReportBetweenTSFromDB(db, day, nextDay, 100)
		if err != nil {
			return "", err
		}
		if noEntriesFound && len(taskLogEntries) > 0 {
			noEntriesFound = false
		}

		day = nextDay
		reportData[i] = taskLogEntries
		if len(taskLogEntries) > maxEntryForADay {
			maxEntryForADay = len(taskLogEntries)
		}
	}

	if noEntriesFound {
		maxEntryForADay = 1
	}

	data := make([][]string, maxEntryForADay)
	totalSecsPerDay := make(map[int]int)

	for j := 0; j < numDays; j++ {
		totalSecsPerDay[j] = 0
	}

	rs := getReportStyles(plain)

	var summaryBudget int
	switch numDays {
	case 7:
		summaryBudget = 8
	case 6:
		summaryBudget = 10
	case 5:
		summaryBudget = 14
	default:
		summaryBudget = 16
	}

	styleCache := make(map[string]lipgloss.Style)
	for rowIndex := 0; rowIndex < maxEntryForADay; rowIndex++ {
		row := make([]string, numDays)
		for colIndex := 0; colIndex < numDays; colIndex++ {
			if rowIndex >= len(reportData[colIndex]) {
				row[colIndex] = fmt.Sprintf("%s  %s",
					RightPadTrim("", summaryBudget, false),
					RightPadTrim("", reportTimeCharsBudget, false),
				)
				continue
			}

			tr := reportData[colIndex][rowIndex]
			timeSpentStr := humanizeDuration(tr.secsSpent)

			if plain {
				row[colIndex] = fmt.Sprintf("%s  %s",
					RightPadTrim(tr.taskSummary, summaryBudget, false),
					RightPadTrim(timeSpentStr, reportTimeCharsBudget, false),
				)
			} else {
				reportStyle, ok := styleCache[tr.taskSummary]
				if !ok {
					reportStyle = getDynamicStyle(tr.taskSummary)
					styleCache[tr.taskSummary] = reportStyle
				}

				row[colIndex] = fmt.Sprintf("%s  %s",
					reportStyle.Render(RightPadTrim(tr.taskSummary, summaryBudget, false)),
					reportStyle.Render(RightPadTrim(timeSpentStr, reportTimeCharsBudget, false)),
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

	b := bytes.Buffer{}
	table := tablewriter.NewWriter(&b)

	headersValues := make([]string, numDays)

	day = start
	counter := 0

	for counter < numDays {
		headersValues[counter] = day.Format(dateFormat)
		day = day.AddDate(0, 0, 1)
		counter++
	}

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

	return b.String(), nil
}
