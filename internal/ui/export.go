package ui

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/dhth/hours/internal/types"
)

type exportLogEntry struct {
	TaskID    int    `json:"taskId"`
	Task      string `json:"task"`
	BeginTS   string `json:"beginTs"`
	EndTS     string `json:"endTs"`
	SecsSpent int    `json:"secsSpent"`
	TimeSpent string `json:"timeSpent"`
	Comment   string `json:"comment"`
}

type exportReportDay struct {
	Date           string           `json:"date"`
	TotalSecsSpent int              `json:"totalSecsSpent"`
	TotalTimeSpent string           `json:"totalTimeSpent"`
	Entries        []exportLogEntry `json:"entries"`
}

type exportReportJSON struct {
	Days []exportReportDay `json:"days"`
}

type exportAggEntry struct {
	TaskID     int    `json:"taskId"`
	Task       string `json:"task"`
	NumEntries int    `json:"numEntries"`
	SecsSpent  int    `json:"secsSpent"`
	TimeSpent  string `json:"timeSpent"`
}

type exportReportAggDay struct {
	Date           string           `json:"date"`
	TotalSecsSpent int              `json:"totalSecsSpent"`
	TotalTimeSpent string           `json:"totalTimeSpent"`
	Entries        []exportAggEntry `json:"entries"`
}

type exportReportAggJSON struct {
	Days []exportReportAggDay `json:"days"`
}

type exportLogJSON struct {
	Entries []exportLogEntry `json:"entries"`
}

type exportStatsEntry struct {
	TaskID     int    `json:"taskId"`
	Task       string `json:"task"`
	NumEntries int    `json:"numEntries"`
	SecsSpent  int    `json:"secsSpent"`
	TimeSpent  string `json:"timeSpent"`
}

type exportStatsJSON struct {
	Entries []exportStatsEntry `json:"entries"`
}

func exportComment(comment *string) string {
	if comment == nil {
		return ""
	}

	return *comment
}

func toExportLogEntry(entry types.TaskLogEntry) exportLogEntry {
	return exportLogEntry{
		TaskID:    entry.TaskID,
		Task:      entry.TaskSummary,
		BeginTS:   entry.BeginTS.Format(time.RFC3339),
		EndTS:     entry.EndTS.Format(time.RFC3339),
		SecsSpent: entry.SecsSpent,
		TimeSpent: types.HumanizeDuration(entry.SecsSpent),
		Comment:   exportComment(entry.Comment),
	}
}

func toExportAggEntry(entry types.TaskReportEntry) exportAggEntry {
	return exportAggEntry{
		TaskID:     entry.TaskID,
		Task:       entry.TaskSummary,
		NumEntries: entry.NumEntries,
		SecsSpent:  entry.SecsSpent,
		TimeSpent:  types.HumanizeDuration(entry.SecsSpent),
	}
}

func toExportStatsEntry(entry types.TaskReportEntry) exportStatsEntry {
	return exportStatsEntry{
		TaskID:     entry.TaskID,
		Task:       entry.TaskSummary,
		NumEntries: entry.NumEntries,
		SecsSpent:  entry.SecsSpent,
		TimeSpent:  types.HumanizeDuration(entry.SecsSpent),
	}
}

func writeJSON(writer io.Writer, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(writer, string(data))
	return err
}

func writeReportJSON(
	writer io.Writer,
	start time.Time,
	numDays int,
	reportData map[int][]types.TaskLogEntry,
) error {
	days := make([]exportReportDay, 0, numDays)
	day := start
	for i := range numDays {
		entries := reportData[i]
		totalSecs := 0
		exportEntries := make([]exportLogEntry, 0, len(entries))
		for _, entry := range entries {
			totalSecs += entry.SecsSpent
			exportEntries = append(exportEntries, toExportLogEntry(entry))
		}

		days = append(days, exportReportDay{
			Date:           day.Format(dateFormat),
			TotalSecsSpent: totalSecs,
			TotalTimeSpent: types.HumanizeDuration(totalSecs),
			Entries:        exportEntries,
		})
		day = day.AddDate(0, 0, 1)
	}

	return writeJSON(writer, exportReportJSON{Days: days})
}

func writeReportAggJSON(
	writer io.Writer,
	start time.Time,
	numDays int,
	reportData map[int][]types.TaskReportEntry,
) error {
	days := make([]exportReportAggDay, 0, numDays)
	day := start
	for i := range numDays {
		entries := reportData[i]
		totalSecs := 0
		exportEntries := make([]exportAggEntry, 0, len(entries))
		for _, entry := range entries {
			totalSecs += entry.SecsSpent
			exportEntries = append(exportEntries, toExportAggEntry(entry))
		}

		days = append(days, exportReportAggDay{
			Date:           day.Format(dateFormat),
			TotalSecsSpent: totalSecs,
			TotalTimeSpent: types.HumanizeDuration(totalSecs),
			Entries:        exportEntries,
		})
		day = day.AddDate(0, 0, 1)
	}

	return writeJSON(writer, exportReportAggJSON{Days: days})
}

func writeReportCSV(
	writer io.Writer,
	start time.Time,
	numDays int,
	reportData map[int][]types.TaskLogEntry,
) error {
	w := csv.NewWriter(writer)
	defer w.Flush()

	if err := w.Write([]string{"date", "taskId", "task", "beginTs", "endTs", "secsSpent", "timeSpent", "comment"}); err != nil {
		return err
	}

	day := start
	for i := range numDays {
		date := day.Format(dateFormat)
		for _, entry := range reportData[i] {
			exportEntry := toExportLogEntry(entry)
			if err := w.Write([]string{
				date,
				strconv.Itoa(exportEntry.TaskID),
				exportEntry.Task,
				exportEntry.BeginTS,
				exportEntry.EndTS,
				strconv.Itoa(exportEntry.SecsSpent),
				exportEntry.TimeSpent,
				exportEntry.Comment,
			}); err != nil {
				return err
			}
		}
		day = day.AddDate(0, 0, 1)
	}

	return w.Error()
}

func writeReportAggCSVRows(
	w *csv.Writer,
	start time.Time,
	numDays int,
	reportData map[int][]types.TaskReportEntry,
) error {
	day := start
	for i := range numDays {
		date := day.Format(dateFormat)
		for _, entry := range reportData[i] {
			exportEntry := toExportAggEntry(entry)
			if err := w.Write([]string{
				date,
				strconv.Itoa(exportEntry.TaskID),
				exportEntry.Task,
				strconv.Itoa(exportEntry.NumEntries),
				strconv.Itoa(exportEntry.SecsSpent),
				exportEntry.TimeSpent,
			}); err != nil {
				return err
			}
		}
		day = day.AddDate(0, 0, 1)
	}

	return w.Error()
}

func writeReportAggCSV(
	writer io.Writer,
	start time.Time,
	numDays int,
	reportData map[int][]types.TaskReportEntry,
) error {
	w := csv.NewWriter(writer)
	defer w.Flush()

	if err := w.Write([]string{"date", "taskId", "task", "numEntries", "secsSpent", "timeSpent"}); err != nil {
		return err
	}

	return writeReportAggCSVRows(w, start, numDays, reportData)
}

func writeLogJSON(writer io.Writer, entries []types.TaskLogEntry) error {
	exportEntries := make([]exportLogEntry, 0, len(entries))
	for _, entry := range entries {
		exportEntries = append(exportEntries, toExportLogEntry(entry))
	}

	return writeJSON(writer, exportLogJSON{Entries: exportEntries})
}

func writeLogCSV(writer io.Writer, entries []types.TaskLogEntry) error {
	w := csv.NewWriter(writer)
	defer w.Flush()

	if err := w.Write([]string{"taskId", "task", "beginTs", "endTs", "secsSpent", "timeSpent", "comment"}); err != nil {
		return err
	}

	for _, entry := range entries {
		exportEntry := toExportLogEntry(entry)
		if err := w.Write([]string{
			strconv.Itoa(exportEntry.TaskID),
			exportEntry.Task,
			exportEntry.BeginTS,
			exportEntry.EndTS,
			strconv.Itoa(exportEntry.SecsSpent),
			exportEntry.TimeSpent,
			exportEntry.Comment,
		}); err != nil {
			return err
		}
	}

	return w.Error()
}

func writeStatsJSON(writer io.Writer, entries []types.TaskReportEntry) error {
	exportEntries := make([]exportStatsEntry, 0, len(entries))
	for _, entry := range entries {
		exportEntries = append(exportEntries, toExportStatsEntry(entry))
	}

	return writeJSON(writer, exportStatsJSON{Entries: exportEntries})
}

func writeStatsCSV(writer io.Writer, entries []types.TaskReportEntry) error {
	w := csv.NewWriter(writer)
	defer w.Flush()

	if err := w.Write([]string{"taskId", "task", "numEntries", "secsSpent", "timeSpent"}); err != nil {
		return err
	}

	for _, entry := range entries {
		exportEntry := toExportStatsEntry(entry)
		if err := w.Write([]string{
			strconv.Itoa(exportEntry.TaskID),
			exportEntry.Task,
			strconv.Itoa(exportEntry.NumEntries),
			strconv.Itoa(exportEntry.SecsSpent),
			exportEntry.TimeSpent,
		}); err != nil {
			return err
		}
	}

	return w.Error()
}
