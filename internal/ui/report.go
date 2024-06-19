package ui

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/olekukonko/tablewriter"
)

const (
	reportNumDaysUpperBound = 7
	timeCharsBudget         = 6
)

func RenderReport(db *sql.DB, writer io.Writer, plain bool, period string, agg bool, interactive bool) {
	if period == "" {
		return
	}

	var start time.Time
	var numDays int

	switch period {
	case "today":
		numDays = 1
		now := time.Now()
		nDaysBack := now.AddDate(0, 0, -1*(numDays-1))

		start = time.Date(nDaysBack.Year(), nDaysBack.Month(), nDaysBack.Day(), 0, 0, 0, 0, nDaysBack.Location())

	case "yest":
		numDays = 1
		now := time.Now().AddDate(0, 0, -1)
		nDaysBack := now.AddDate(0, 0, -1*(numDays-1))

		start = time.Date(nDaysBack.Year(), nDaysBack.Month(), nDaysBack.Day(), 0, 0, 0, 0, nDaysBack.Location())

	case "3d":
		numDays = 3
		now := time.Now()
		nDaysBack := now.AddDate(0, 0, -1*(numDays-1))

		start = time.Date(nDaysBack.Year(), nDaysBack.Month(), nDaysBack.Day(), 0, 0, 0, 0, nDaysBack.Location())

	case "week":
		now := time.Now()
		weekday := now.Weekday()
		offset := (7 + weekday - time.Monday) % 7
		startOfWeek := now.AddDate(0, 0, -int(offset))
		start = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, startOfWeek.Location())
		if interactive {
			numDays = 7
		} else {
			numDays = int(offset) + 1
		}

	default:
		var end time.Time
		var err error

		if strings.Contains(period, "...") {
			var ts timePeriod
			var nd int
			ts, nd, err = parseDateDuration(period)
			if err != nil {
				fmt.Fprintf(writer, "%s\n", err)
				os.Exit(1)
			}
			if nd > reportNumDaysUpperBound {
				fmt.Fprintf(writer, "Time period is too large; maximum number of days allowed in range (both inclusive): %d\n", reportNumDaysUpperBound)
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
		numDays = int(end.Sub(start).Hours() / 24)
	}

	var report string
	var err error

	if agg {
		report, err = getReportAgg(db, start, numDays, plain)
	} else {
		report, err = getReport(db, start, numDays, plain)
	}
	if err != nil {
		fmt.Printf("Something went wrong generating the report: %s\n", err)
	}

	if interactive {
		p := tea.NewProgram(initialReportModel(db, start, plain, period, numDays, agg, report))
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
					RightPadTrim("", timeCharsBudget, false),
				)
				continue
			}

			tr := reportData[colIndex][rowIndex]
			timeSpentStr := humanizeDuration(tr.secsSpent)

			if plain {
				row[colIndex] = fmt.Sprintf("%s  %s",
					RightPadTrim(tr.taskSummary, summaryBudget, false),
					RightPadTrim(timeSpentStr, timeCharsBudget, false),
				)
			} else {
				reportStyle, ok := styleCache[tr.taskSummary]

				if !ok {
					reportStyle = getDynamicStyle(tr.taskSummary)
					styleCache[tr.taskSummary] = reportStyle
				}

				row[colIndex] = fmt.Sprintf("%s  %s",
					reportStyle.Render(RightPadTrim(tr.taskSummary, summaryBudget, false)),
					reportStyle.Render(RightPadTrim(timeSpentStr, timeCharsBudget, false)),
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
					RightPadTrim("", timeCharsBudget, false),
				)
				continue
			}

			tr := reportData[colIndex][rowIndex]
			timeSpentStr := humanizeDuration(tr.secsSpent)

			if plain {
				row[colIndex] = fmt.Sprintf("%s  %s",
					RightPadTrim(tr.taskSummary, summaryBudget, false),
					RightPadTrim(timeSpentStr, timeCharsBudget, false),
				)
			} else {
				reportStyle, ok := styleCache[tr.taskSummary]
				if !ok {
					reportStyle = getDynamicStyle(tr.taskSummary)
					styleCache[tr.taskSummary] = reportStyle
				}

				row[colIndex] = fmt.Sprintf("%s  %s",
					reportStyle.Render(RightPadTrim(tr.taskSummary, summaryBudget, false)),
					reportStyle.Render(RightPadTrim(timeSpentStr, timeCharsBudget, false)),
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
