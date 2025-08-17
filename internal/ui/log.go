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
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
)

const (
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
) error {
	if interactive && dateRange.NumDays > interactiveLogDayLimit {
		return fmt.Errorf("%w (limited to %d day); use non-interactive mode to see logs for a larger time period", errInteractiveModeNotApplicable, interactiveLogDayLimit)
	}

	log, err := getTaskLog(db, style, dateRange.Start, dateRange.End, taskStatus, logLimit, plain)
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
	plain bool) (string,
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
		data[0] = []string{
			utils.RightPadTrim("", 20, false),
			utils.RightPadTrim("", 40, false),
			utils.RightPadTrim("", 39, false),
			utils.RightPadTrim("", logTimeCharsBudget, false),
		}
	}

	var timeSpentStr string

	rs := style.getReportStyles(plain)
	styleCache := make(map[string]lipgloss.Style)

	for i, entry := range entries {
		timeSpentStr = types.HumanizeDuration(entry.SecsSpent)

		if plain {
			data[i] = []string{
				utils.RightPadTrim(entry.TaskSummary, 20, false),
				utils.RightPadTrimWithMoreLinesIndicator(entry.GetComment(), 40),
				fmt.Sprintf("%s  ...  %s", entry.BeginTS.Format(timeFormat), entry.EndTS.Format(timeFormat)),
				utils.RightPadTrim(timeSpentStr, logTimeCharsBudget, false),
			}
		} else {
			rowStyle, ok := styleCache[entry.TaskSummary]
			if !ok {
				rowStyle = style.getDynamicStyle(entry.TaskSummary)
				styleCache[entry.TaskSummary] = rowStyle
			}
			data[i] = []string{
				rowStyle.Render(utils.RightPadTrim(entry.TaskSummary, 20, false)),
				rowStyle.Render(utils.RightPadTrimWithMoreLinesIndicator(entry.GetComment(), 40)),
				rowStyle.Render(fmt.Sprintf("%s  ...  %s", entry.BeginTS.Format(timeFormat), entry.EndTS.Format(timeFormat))),
				rowStyle.Render(utils.RightPadTrim(timeSpentStr, logTimeCharsBudget, false)),
			}
		}
	}

	headerValues := []string{"Task", "Comment", "Duration", "TimeSpent"}
	headers := make([]string, len(headerValues))
	for i, h := range headerValues {
		headers[i] = rs.headerStyle.Render(h)
	}

	b := bytes.Buffer{}
	table := tablewriter.NewTable(&b,
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
		tablewriter.WithRenderer(renderer.NewBlueprint(tw.Rendition{Symbols: tw.NewSymbols(tw.StyleASCII)})),
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
