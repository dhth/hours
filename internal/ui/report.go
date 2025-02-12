package ui

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	pers "github.com/dhth/hours/internal/persistence"
	"github.com/dhth/hours/internal/types"
	"github.com/dhth/hours/internal/utils"
	"github.com/olekukonko/tablewriter"
)

var errCouldntGenerateReport = errors.New("couldn't generate report")

const (
	reportTimeCharsBudget = 6
)

func RenderReport(db *sql.DB, style *Style, writer io.Writer, plain bool, period string, agg bool, interactive bool) error {
	if period == "" {
		return nil
	}

	var fullWeek bool
	if interactive {
		fullWeek = true
	}
	ts, err := types.GetTimePeriod(period, time.Now(), fullWeek)
	if err != nil {
		return err
	}

	var report string
	var analyticsType recordsType

	if agg {
		analyticsType = reportAggRecords
		report, err = getReportAgg(db, style, ts.Start, ts.NumDays, plain)
	} else {
		analyticsType = reportRecords
		report, err = getReport(db, style, ts.Start, ts.NumDays, plain)
	}
	if err != nil {
		return fmt.Errorf("%w: %s", errCouldntGenerateReport, err.Error())
	}

	if interactive {
		p := tea.NewProgram(initialRecordsModel(analyticsType, db, style, ts.Start, ts.End, plain, period, ts.NumDays, report))
		_, err := p.Run()
		if err != nil {
			return err
		}
	} else {
		fmt.Fprint(writer, report)
	}
	return nil
}

func getReport(db *sql.DB, style *Style, start time.Time, numDays int, plain bool) (string, error) {
	day := start
	var nextDay time.Time

	var maxEntryForADay int
	reportData := make(map[int][]types.TaskLogEntry)

	noEntriesFound := true
	for i := 0; i < numDays; i++ {
		nextDay = day.AddDate(0, 0, 1)
		taskLogEntries, err := pers.FetchTLEntriesBetweenTS(db, day, nextDay, 100)
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

	rs := style.getReportStyles(plain)

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
					utils.RightPadTrim("", summaryBudget, false),
					utils.RightPadTrim("", reportTimeCharsBudget, false),
				)
				continue
			}

			tr := reportData[colIndex][rowIndex]
			timeSpentStr := types.HumanizeDuration(tr.SecsSpent)

			if plain {
				row[colIndex] = fmt.Sprintf("%s  %s",
					utils.RightPadTrim(tr.TaskSummary, summaryBudget, false),
					utils.RightPadTrim(timeSpentStr, reportTimeCharsBudget, false),
				)
			} else {
				rowStyle, ok := styleCache[tr.TaskSummary]

				if !ok {
					rowStyle = style.getDynamicStyle(tr.TaskSummary)
					styleCache[tr.TaskSummary] = rowStyle
				}

				row[colIndex] = fmt.Sprintf("%s  %s",
					rowStyle.Render(utils.RightPadTrim(tr.TaskSummary, summaryBudget, false)),
					rowStyle.Render(utils.RightPadTrim(timeSpentStr, reportTimeCharsBudget, false)),
				)
			}
			totalSecsPerDay[colIndex] += tr.SecsSpent
		}
		data[rowIndex] = row
	}

	totalTimePerDay := make([]string, numDays)

	for i, ts := range totalSecsPerDay {
		if ts != 0 {
			totalTimePerDay[i] = rs.footerStyle.Render(types.HumanizeDuration(ts))
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

func getReportAgg(db *sql.DB, style *Style, start time.Time, numDays int, plain bool) (string, error) {
	day := start
	var nextDay time.Time

	var maxEntryForADay int
	reportData := make(map[int][]types.TaskReportEntry)

	noEntriesFound := true
	for i := 0; i < numDays; i++ {
		nextDay = day.AddDate(0, 0, 1)
		taskLogEntries, err := pers.FetchReportBetweenTS(db, day, nextDay, 100)
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

	rs := style.getReportStyles(plain)

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
					utils.RightPadTrim("", summaryBudget, false),
					utils.RightPadTrim("", reportTimeCharsBudget, false),
				)
				continue
			}

			tr := reportData[colIndex][rowIndex]
			timeSpentStr := types.HumanizeDuration(tr.SecsSpent)

			if plain {
				row[colIndex] = fmt.Sprintf("%s  %s",
					utils.RightPadTrim(tr.TaskSummary, summaryBudget, false),
					utils.RightPadTrim(timeSpentStr, reportTimeCharsBudget, false),
				)
			} else {
				rowStyle, ok := styleCache[tr.TaskSummary]
				if !ok {
					rowStyle = style.getDynamicStyle(tr.TaskSummary)
					styleCache[tr.TaskSummary] = rowStyle
				}

				row[colIndex] = fmt.Sprintf("%s  %s",
					rowStyle.Render(utils.RightPadTrim(tr.TaskSummary, summaryBudget, false)),
					rowStyle.Render(utils.RightPadTrim(timeSpentStr, reportTimeCharsBudget, false)),
				)
			}
			totalSecsPerDay[colIndex] += tr.SecsSpent
		}
		data[rowIndex] = row
	}
	totalTimePerDay := make([]string, numDays)
	for i, ts := range totalSecsPerDay {
		if ts != 0 {
			totalTimePerDay[i] = rs.footerStyle.Render(types.HumanizeDuration(ts))
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
