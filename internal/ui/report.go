package ui

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	pers "github.com/dhth/hours/internal/persistence"
	"github.com/dhth/hours/internal/types"
	"github.com/dhth/hours/internal/utils"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
)

var errCouldntGenerateReport = errors.New("couldn't generate report")

const (
	reportTimeCharsBudget = 6
)

func RenderReport(db *sql.DB,
	style Style,
	writer io.Writer,
	plain bool,
	dateRange types.DateRange,
	period string,
	taskStatus types.TaskStatus,
	agg bool,
	interactive bool,
	format types.OutputFormat,
) error {
	if interactive && format != types.OutputFormatPlain {
		return fmt.Errorf("%w with output format %q", errInteractiveModeNotApplicable, format)
	}

	if format != types.OutputFormatPlain {
		return renderReportExport(db, writer, dateRange, taskStatus, agg, format)
	}

	var report string
	var analyticsType recordsKind
	var err error

	if agg {
		analyticsType = reportAggRecords
		report, err = getReportAgg(db, style, dateRange.Start, dateRange.NumDays, taskStatus, plain)
	} else {
		analyticsType = reportRecords
		report, err = getReport(db, style, dateRange.Start, dateRange.NumDays, taskStatus, plain)
	}
	if err != nil {
		return fmt.Errorf("%w: %s", errCouldntGenerateReport, err.Error())
	}

	if interactive {
		p := tea.NewProgram(initialRecordsModel(
			analyticsType,
			db,
			style,
			types.RealTimeProvider{},
			dateRange,
			period,
			taskStatus,
			plain,
			report,
		))
		_, err := p.Run()
		if err != nil {
			return err
		}
	} else {
		fmt.Fprint(writer, report)
	}
	return nil
}

func renderReportExport(
	db *sql.DB,
	writer io.Writer,
	dateRange types.DateRange,
	taskStatus types.TaskStatus,
	agg bool,
	format types.OutputFormat,
) error {
	if agg {
		reportData, err := collectReportAggData(db, dateRange.Start, dateRange.NumDays, taskStatus)
		if err != nil {
			return fmt.Errorf("%w: %s", errCouldntGenerateReport, err.Error())
		}

		switch format {
		case types.OutputFormatJSON:
			return writeReportAggJSON(writer, dateRange.Start, dateRange.NumDays, reportData)
		case types.OutputFormatCSV:
			return writeReportAggCSV(writer, dateRange.Start, dateRange.NumDays, reportData)
		default:
			return types.ErrIncorrectOutputFormatProvided
		}
	}

	reportData, err := collectReportData(db, dateRange.Start, dateRange.NumDays, taskStatus)
	if err != nil {
		return fmt.Errorf("%w: %s", errCouldntGenerateReport, err.Error())
	}

	switch format {
	case types.OutputFormatJSON:
		return writeReportJSON(writer, dateRange.Start, dateRange.NumDays, reportData)
	case types.OutputFormatCSV:
		return writeReportCSV(writer, dateRange.Start, dateRange.NumDays, reportData)
	default:
		return types.ErrIncorrectOutputFormatProvided
	}
}

func collectReportData(
	db *sql.DB,
	start time.Time,
	numDays int,
	taskStatus types.TaskStatus,
) (map[int][]types.TaskLogEntry, error) {
	day := start
	reportData := make(map[int][]types.TaskLogEntry, numDays)

	for i := range numDays {
		nextDay := day.AddDate(0, 0, 1)
		taskLogEntries, err := pers.FetchTLEntriesBetweenTS(db, day, nextDay, taskStatus, 100)
		if err != nil {
			return nil, err
		}

		reportData[i] = taskLogEntries
		day = nextDay
	}

	return reportData, nil
}

func collectReportAggData(
	db *sql.DB,
	start time.Time,
	numDays int,
	taskStatus types.TaskStatus,
) (map[int][]types.TaskReportEntry, error) {
	day := start
	reportData := make(map[int][]types.TaskReportEntry, numDays)

	for i := range numDays {
		nextDay := day.AddDate(0, 0, 1)
		taskReportEntries, err := pers.FetchReportBetweenTS(db, day, nextDay, taskStatus, 100)
		if err != nil {
			return nil, err
		}

		reportData[i] = taskReportEntries
		day = nextDay
	}

	return reportData, nil
}

func getReport(db *sql.DB, style Style, start time.Time, numDays int, taskStatus types.TaskStatus, plain bool) (string, error) {
	reportData, err := collectReportData(db, start, numDays, taskStatus)
	if err != nil {
		return "", err
	}

	var maxEntryForADay int
	noEntriesFound := true
	for i := range numDays {
		if noEntriesFound && len(reportData[i]) > 0 {
			noEntriesFound = false
		}

		if len(reportData[i]) > maxEntryForADay {
			maxEntryForADay = len(reportData[i])
		}
	}

	if noEntriesFound {
		maxEntryForADay = 1
	}

	data := make([][]string, maxEntryForADay)
	totalSecsPerDay := make(map[int]int)

	for j := range numDays {
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
	for rowIndex := range maxEntryForADay {
		row := make([]string, numDays)
		for colIndex := range numDays {
			if rowIndex >= len(reportData[colIndex]) {
				row[colIndex] = fmt.Sprintf(
					"%s  %s",
					utils.RightPadTrim("", summaryBudget, false),
					utils.RightPadTrim("", reportTimeCharsBudget, false),
				)
				continue
			}

			tr := reportData[colIndex][rowIndex]
			timeSpentStr := types.HumanizeDuration(tr.SecsSpent)

			if plain {
				row[colIndex] = fmt.Sprintf(
					"%s  %s",
					utils.RightPadTrim(tr.TaskSummary, summaryBudget, false),
					utils.RightPadTrim(timeSpentStr, reportTimeCharsBudget, false),
				)
			} else {
				rowStyle, ok := styleCache[tr.TaskSummary]

				if !ok {
					rowStyle = style.getDynamicStyle(tr.TaskSummary)
					styleCache[tr.TaskSummary] = rowStyle
				}

				row[colIndex] = fmt.Sprintf(
					"%s  %s",
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

	headersValues := make([]string, numDays)

	day := start
	counter := 0

	for counter < numDays {
		headersValues[counter] = day.Format(dateFormat)
		day = day.AddDate(0, 0, 1)
		counter++
	}

	headers := make([]string, numDays)
	for i := range numDays {
		headers[i] = rs.headerStyle.Render(headersValues[i])
	}

	b := bytes.Buffer{}
	table := tablewriter.NewTable(
		&b,
		tablewriter.WithConfig(tablewriter.Config{
			Header: tw.CellConfig{
				Formatting: tw.CellFormatting{
					Alignment:  tw.AlignCenter,
					AutoWrap:   tw.WrapNone,
					AutoFormat: tw.Off,
				},
			},
			Row: tw.CellConfig{
				Formatting: tw.CellFormatting{
					Alignment: tw.AlignLeft,
					AutoWrap:  tw.WrapNone,
				},
			},
			Footer: tw.CellConfig{
				Formatting: tw.CellFormatting{
					Alignment:  tw.AlignCenter,
					AutoWrap:   tw.WrapNone,
					AutoFormat: tw.Off,
				},
			},
		}),
		tablewriter.WithRenderer(renderer.NewBlueprint(tw.Rendition{Symbols: rs.symbols(tw.StyleASCII)})),
		tablewriter.WithHeader(headers),
		tablewriter.WithFooter(totalTimePerDay),
	)

	if err := table.Bulk(data); err != nil {
		return "", fmt.Errorf("%w: %s", errCouldntAddDataToTable, err.Error())
	}

	if err := table.Render(); err != nil {
		return "", fmt.Errorf("%w: %s", errCouldntRenderTable, err.Error())
	}

	return b.String(), nil
}

func getReportAgg(db *sql.DB,
	style Style,
	start time.Time,
	numDays int,
	taskStatus types.TaskStatus,
	plain bool) (string,
	error,
) {
	reportData, err := collectReportAggData(db, start, numDays, taskStatus)
	if err != nil {
		return "", err
	}

	var maxEntryForADay int
	noEntriesFound := true
	for i := range numDays {
		if noEntriesFound && len(reportData[i]) > 0 {
			noEntriesFound = false
		}

		if len(reportData[i]) > maxEntryForADay {
			maxEntryForADay = len(reportData[i])
		}
	}

	if noEntriesFound {
		maxEntryForADay = 1
	}

	data := make([][]string, maxEntryForADay)
	totalSecsPerDay := make(map[int]int)

	for j := range numDays {
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
	for rowIndex := range maxEntryForADay {
		row := make([]string, numDays)
		for colIndex := range numDays {
			if rowIndex >= len(reportData[colIndex]) {
				row[colIndex] = fmt.Sprintf(
					"%s  %s",
					utils.RightPadTrim("", summaryBudget, false),
					utils.RightPadTrim("", reportTimeCharsBudget, false),
				)
				continue
			}

			tr := reportData[colIndex][rowIndex]
			timeSpentStr := types.HumanizeDuration(tr.SecsSpent)

			if plain {
				row[colIndex] = fmt.Sprintf(
					"%s  %s",
					utils.RightPadTrim(tr.TaskSummary, summaryBudget, false),
					utils.RightPadTrim(timeSpentStr, reportTimeCharsBudget, false),
				)
			} else {
				rowStyle, ok := styleCache[tr.TaskSummary]
				if !ok {
					rowStyle = style.getDynamicStyle(tr.TaskSummary)
					styleCache[tr.TaskSummary] = rowStyle
				}

				row[colIndex] = fmt.Sprintf(
					"%s  %s",
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

	headersValues := make([]string, numDays)

	day := start
	counter := 0

	for counter < numDays {
		headersValues[counter] = day.Format(dateFormat)
		day = day.AddDate(0, 0, 1)
		counter++
	}

	headers := make([]string, numDays)
	for i := range numDays {
		headers[i] = rs.headerStyle.Render(headersValues[i])
	}

	b := bytes.Buffer{}
	table := tablewriter.NewTable(
		&b,
		tablewriter.WithConfig(tablewriter.Config{
			Header: tw.CellConfig{
				Formatting: tw.CellFormatting{
					Alignment:  tw.AlignCenter,
					AutoWrap:   tw.WrapNone,
					AutoFormat: tw.Off,
				},
			},
			Row: tw.CellConfig{
				Formatting: tw.CellFormatting{
					Alignment: tw.AlignLeft,
					AutoWrap:  tw.WrapNone,
				},
			},
			Footer: tw.CellConfig{
				Formatting: tw.CellFormatting{
					Alignment:  tw.AlignCenter,
					AutoWrap:   tw.WrapNone,
					AutoFormat: tw.Off,
				},
			},
		}),
		tablewriter.WithRenderer(renderer.NewBlueprint(tw.Rendition{Symbols: rs.symbols(tw.StyleASCII)})),
		tablewriter.WithHeader(headers),
		tablewriter.WithFooter(totalTimePerDay),
	)

	if err := table.Bulk(data); err != nil {
		return "", fmt.Errorf("%w: %s", errCouldntAddDataToTable, err.Error())
	}

	if err := table.Render(); err != nil {
		return "", fmt.Errorf("%w: %s", errCouldntRenderTable, err.Error())
	}

	return b.String(), nil
}
