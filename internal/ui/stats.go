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

var errCouldntGenerateStats = errors.New("couldn't generate stats")

const (
	statsLogEntriesLimit   = 10000
	statsNumDaysUpperBound = 3650
	statsTimeCharsBudget   = 6
	periodAll              = "all"
)

func RenderStats(
	db *sql.DB,
	style Style,
	writer io.Writer,
	plain bool,
	period string,
	activeFilter types.TaskActiveStatusFilter,
	interactive bool,
) error {
	if period == "" {
		return nil
	}

	var stats string
	var err error

	if interactive && period == periodAll {
		return fmt.Errorf("%w when period=all", errInteractiveModeNotApplicable)
	}

	fetcher := newStatsFetcher(activeFilter)

	if period == periodAll {
		stats, err = getStats(db, style, fetcher, plain)
		if err != nil {
			return fmt.Errorf("%w: %s", errCouldntGenerateStats, err.Error())
		}
		fmt.Fprint(writer, stats)
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

	stats, err = getStats(db, style, fetcher.ts(ts.Start, ts.End), plain)
	if err != nil {
		return fmt.Errorf("%w: %s", errCouldntGenerateStats, err.Error())
	}

	if interactive {
		p := tea.NewProgram(initialRecordsModel(
			reportStats,
			db,
			style,
			ts.Start,
			ts.End,
			activeFilter,
			plain,
			period,
			ts.NumDays,
			stats,
		))
		_, err := p.Run()
		if err != nil {
			return err
		}
	} else {
		fmt.Fprint(writer, stats)
	}
	return nil
}

type statsFetcher struct {
	all          bool
	start, end   time.Time
	activeFilter types.TaskActiveStatusFilter
}

func newStatsFetcher(activeFilter types.TaskActiveStatusFilter) statsFetcher {
	return statsFetcher{
		all:          true,
		activeFilter: activeFilter,
	}
}

func newStatsFetcherFromPeriod(
	period string,
	activeFilter types.TaskActiveStatusFilter,
	start, end time.Time,
) statsFetcher {
	f := newStatsFetcher(activeFilter)
	if period == periodAll {
		return f
	}
	return f.ts(start, end)
}

func (f statsFetcher) ts(start, end time.Time) statsFetcher {
	return statsFetcher{
		all:          false,
		start:        start,
		end:          end,
		activeFilter: f.activeFilter,
	}
}

func (f statsFetcher) fetch(db *sql.DB) ([]types.TaskReportEntry, error) {
	if f.all {
		return pers.FetchStats(db, statsLogEntriesLimit, f.activeFilter)
	}
	return pers.FetchStatsBetweenTS(db, f.start, f.end, statsLogEntriesLimit, f.activeFilter)
}

func getStats(db *sql.DB, style Style, fetcher statsFetcher, plain bool) (string, error) {
	var entries []types.TaskReportEntry
	var err error

	if entries, err = fetcher.fetch(db); err != nil {
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
			"",
			utils.RightPadTrim("", statsTimeCharsBudget, false),
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
				fmt.Sprintf("%d", entry.NumEntries),
				utils.RightPadTrim(timeSpentStr, statsTimeCharsBudget, false),
			}
		} else {
			rowStyle, ok := styleCache[entry.TaskSummary]
			if !ok {
				rowStyle = style.getDynamicStyle(entry.TaskSummary)
				styleCache[entry.TaskSummary] = rowStyle
			}
			data[i] = []string{
				rowStyle.Render(utils.RightPadTrim(entry.TaskSummary, 20, false)),
				rowStyle.Render(fmt.Sprintf("%d", entry.NumEntries)),
				rowStyle.Render(utils.RightPadTrim(timeSpentStr, statsTimeCharsBudget, false)),
			}
		}
	}
	b := bytes.Buffer{}
	table := tablewriter.NewWriter(&b)

	headerValues := []string{"Task", "#LogEntries", "TimeSpent"}
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
