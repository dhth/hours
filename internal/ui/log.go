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

const (
	logNumDaysUpperBound   = 7
	logTimeCharsBudget     = 6
	interactiveLogDayLimit = 1
)

var (
	errInteractiveModeNotApplicable = errors.New("interactive mode is not applicable")
	errCouldntGenerateLogs          = errors.New("couldn't generate logs")
)

func RenderTaskLog(db *sql.DB, writer io.Writer, plain bool, period string, interactive bool) error {
	if period == "" {
		return nil
	}

	ts, err := types.GetTimePeriod(period, time.Now(), false)
	if err != nil {
		return err
	}

	if interactive && ts.NumDays > interactiveLogDayLimit {
		return fmt.Errorf("%w (limited to %d day); use non-interactive mode to see logs for a larger time period", errInteractiveModeNotApplicable, interactiveLogDayLimit)
	}

	log, err := getTaskLog(db, ts.Start, ts.End, 100, plain)
	if err != nil {
		return fmt.Errorf("%w: %s", errCouldntGenerateLogs, err.Error())
	}

	if interactive {
		p := tea.NewProgram(initialRecordsModel(reportLogs, db, ts.Start, ts.End, plain, period, ts.NumDays, log))
		_, err := p.Run()
		if err != nil {
			return err
		}
	} else {
		fmt.Fprint(writer, log)
	}
	return nil
}

func getTaskLog(db *sql.DB, start, end time.Time, limit int, plain bool) (string, error) {
	entries, err := pers.FetchTLEntriesBetweenTS(db, start, end, limit)
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

	rs := getReportStyles(plain)
	styleCache := make(map[string]lipgloss.Style)

	for i, entry := range entries {
		timeSpentStr = types.HumanizeDuration(entry.SecsSpent)

		if plain {
			data[i] = []string{
				utils.RightPadTrim(entry.TaskSummary, 20, false),
				utils.RightPadTrim(entry.Comment, 40, false),
				fmt.Sprintf("%s  ...  %s", entry.BeginTS.Format(timeFormat), entry.EndTS.Format(timeFormat)),
				utils.RightPadTrim(timeSpentStr, logTimeCharsBudget, false),
			}
		} else {
			rowStyle, ok := styleCache[entry.TaskSummary]
			if !ok {
				rowStyle = getDynamicStyle(entry.TaskSummary)
				styleCache[entry.TaskSummary] = rowStyle
			}
			data[i] = []string{
				rowStyle.Render(utils.RightPadTrim(entry.TaskSummary, 20, false)),
				rowStyle.Render(utils.RightPadTrim(entry.Comment, 40, false)),
				rowStyle.Render(fmt.Sprintf("%s  ...  %s", entry.BeginTS.Format(timeFormat), entry.EndTS.Format(timeFormat))),
				rowStyle.Render(utils.RightPadTrim(timeSpentStr, logTimeCharsBudget, false)),
			}
		}
	}

	b := bytes.Buffer{}
	table := tablewriter.NewWriter(&b)

	headerValues := []string{"Task", "Comment", "Duration", "TimeSpent"}
	headers := make([]string, len(headerValues))
	for i, h := range headerValues {
		headers[i] = rs.headerStyle.Render(h)
	}
	table.SetHeader(headers)

	table.SetRowSeparator(rs.borderStyle.Render("-"))
	table.SetColumnSeparator(rs.borderStyle.Render("|"))
	table.SetCenterSeparator(rs.borderStyle.Render("+"))
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(false)
	table.AppendBulk(data)

	table.Render()

	return b.String(), nil
}
