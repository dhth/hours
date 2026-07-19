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

const (
	logTaskCharsBudget     = 20
	logCommentCharsBudget  = 40
	logDurationCharsBudget = 39
	logTimeCharsBudget     = 6
	interactiveLogDayLimit = 1
	logLimit               = 10000
)

var errCouldntGenerateLogs = errors.New("couldn't generate logs")

func RenderTaskLog(db *sql.DB,
	style Style,
	writer io.Writer,
	plain bool,
	dateRange types.DateRange,
	period string,
	taskStatus types.TaskStatus,
	interactive bool,
	noTruncate bool,
) error {
	if interactive && dateRange.NumDays > interactiveLogDayLimit {
		return fmt.Errorf("%w (limited to %d day); use non-interactive mode to see logs for a larger time period", errInteractiveModeNotApplicable, interactiveLogDayLimit)
	}

	log, err := getTaskLog(db, style, dateRange.Start, dateRange.End, taskStatus, logLimit, plain, noTruncate)
	if err != nil {
		return fmt.Errorf("%w: %s", errCouldntGenerateLogs, err.Error())
	}

	if interactive {
		p := tea.NewProgram(initialRecordsModel(
			reportLogs,
			db,
			style,
			types.RealTimeProvider{},
			dateRange,
			period,
			taskStatus,
			plain,
			noTruncate,
			log,
		))
		_, err := p.Run()
		if err != nil {
			return err
		}
	} else {
		fmt.Fprint(writer, log)
	}
	return nil
}

func getTaskLog(db *sql.DB,
	style Style,
	start,
	end time.Time,
	taskStatus types.TaskStatus,
	limit int,
	plain bool,
	noTruncate bool) (string,
	error,
) {
	entries, err := pers.FetchTLEntriesBetweenTS(db, start, end, taskStatus, limit)
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
		data[0] = []string{"", "", "", ""}
		if !noTruncate {
			data[0] = []string{
				utils.RightPadTrim("", logTaskCharsBudget, false),
				utils.RightPadTrim("", logCommentCharsBudget, false),
				utils.RightPadTrim("", logDurationCharsBudget, false),
				utils.RightPadTrim("", logTimeCharsBudget, false),
			}
		}
	}

	var timeSpentStr string

	rs := style.getReportStyles(plain)
	styleCache := make(map[string]lipgloss.Style)

	for i, entry := range entries {
		timeSpentStr = types.HumanizeDuration(entry.SecsSpent)

		taskSummary := entry.TaskSummary
		comment := entry.GetComment()
		duration := fmt.Sprintf("%s  ...  %s", entry.BeginTS.Format(timeFormat), entry.EndTS.Format(timeFormat))

		if !noTruncate {
			taskSummary = utils.RightPadTrim(taskSummary, logTaskCharsBudget, false)
			comment = utils.RightPadTrimWithMoreLinesIndicator(comment, logCommentCharsBudget)
			timeSpentStr = utils.RightPadTrim(timeSpentStr, logTimeCharsBudget, false)
		}

		if plain {
			data[i] = []string{taskSummary, comment, duration, timeSpentStr}
		} else {
			rowStyle, ok := styleCache[entry.TaskSummary]
			if !ok {
				rowStyle = style.getDynamicStyle(entry.TaskSummary)
				styleCache[entry.TaskSummary] = rowStyle
			}
			data[i] = []string{
				rowStyle.Render(taskSummary),
				rowStyle.Render(comment),
				rowStyle.Render(duration),
				rowStyle.Render(timeSpentStr),
			}
		}
	}

	headerValues := []string{"Task", "Comment", "Duration", "TimeSpent"}
	headers := make([]string, len(headerValues))
	for i, h := range headerValues {
		headers[i] = rs.headerStyle.Render(h)
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
		}),
		tablewriter.WithRenderer(renderer.NewBlueprint(tw.Rendition{Symbols: rs.symbols(tw.StyleASCII)})),
		tablewriter.WithHeader(headers),
	)

	if err := table.Bulk(data); err != nil {
		return "", fmt.Errorf("%w: %s", errCouldntAddDataToTable, err.Error())
	}

	if err := table.Render(); err != nil {
		return "", fmt.Errorf("%w: %s", errCouldntRenderTable, err.Error())
	}

	return b.String(), nil
}
