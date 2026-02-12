package ui

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"io"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	pers "github.com/dhth/hours/internal/persistence"
	"github.com/dhth/hours/internal/types"
	"github.com/dhth/hours/internal/utils"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
)

var errCouldntGenerateStats = errors.New("couldn't generate stats")

const (
	statsLogEntriesLimit = 10000
	statsTimeCharsBudget = 6
)

func RenderStats(db *sql.DB,
	style Style,
	writer io.Writer,
	plain bool,
	dateRange *types.DateRange,
	period string,
	taskStatus types.TaskStatus,
	interactive bool,
) error {
	var stats string
	var err error

	if interactive && dateRange == nil {
		return fmt.Errorf("%w when period=all", errInteractiveModeNotApplicable)
	}

	if dateRange == nil {
		stats, err = getStats(db, style, dateRange, taskStatus, plain)
		if err != nil {
			return fmt.Errorf("%w: %s", errCouldntGenerateStats, err.Error())
		}

		fmt.Fprint(writer, stats)
		return nil
	}

	stats, err = getStats(db, style, dateRange, taskStatus, plain)
	if err != nil {
		return fmt.Errorf("%w: %s", errCouldntGenerateStats, err.Error())
	}

	if interactive {
		p := tea.NewProgram(initialRecordsModel(
			reportStats,
			db,
			style,
			types.RealTimeProvider{},
			*dateRange,
			period,
			taskStatus,
			plain,
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

func getStats(db *sql.DB,
	style Style,
	dateRange *types.DateRange,
	taskStatus types.TaskStatus,
	plain bool) (string,
	error,
) {
	var entries []types.TaskReportEntry
	var err error

	if dateRange == nil {
		entries, err = pers.FetchStats(db, taskStatus, statsLogEntriesLimit)
	} else {
		entries, err = pers.FetchStatsBetweenTS(db, dateRange.Start, dateRange.End, taskStatus, statsLogEntriesLimit)
	}

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
			"",
			utils.RightPadTrim("", statsTimeCharsBudget, false),
		}
	}

	var timeSpentStr string

	rs := style.getReportStyles(plain)
	styleCache := make(map[string]lipgloss.Style)

	var totalSecs int
	var totalNumEntries int
	for i, entry := range entries {
		timeSpentStr = types.HumanizeDuration(entry.SecsSpent)
		totalSecs += entry.SecsSpent
		totalNumEntries += entry.NumEntries

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

	headerValues := []string{"Task", "#LogEntries", "TimeSpent"}
	headers := make([]string, len(headerValues))
	for i, h := range headerValues {
		headers[i] = rs.headerStyle.Render(h)
	}

	var footer []string
	if len(entries) > 0 {
		totalTimeStr := types.HumanizeDuration(totalSecs)
		if plain {
			footer = []string{
				utils.RightPadTrim("Total", 20, false),
				fmt.Sprintf("%d", totalNumEntries),
				utils.RightPadTrim(totalTimeStr, statsTimeCharsBudget, false),
			}
		} else {
			footer = []string{
				rs.footerStyle.Render(utils.RightPadTrim("Total", 20, false)),
				rs.footerStyle.Render(fmt.Sprintf("%d", totalNumEntries)),
				rs.footerStyle.Render(utils.RightPadTrim(totalTimeStr, statsTimeCharsBudget, false)),
			}
		}
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
		tablewriter.WithFooter(footer),
	)

	if err := table.Bulk(data); err != nil {
		return "", fmt.Errorf("%w: %s", errCouldntAddDataToTable, err.Error())
	}

	if err := table.Render(); err != nil {
		return "", fmt.Errorf("%w: %s", errCouldntRenderTable, err.Error())
	}

	return b.String(), nil
}
