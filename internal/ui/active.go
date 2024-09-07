package ui

import (
	"database/sql"
	"fmt"
	"io"
	"strings"
	"time"

	pers "github.com/dhth/hours/internal/persistence"
	"github.com/dhth/hours/internal/types"
)

const (
	ActiveTaskPlaceholder     = "{{task}}"
	ActiveTaskTimePlaceholder = "{{time}}"
	activeSecsThreshold       = 60
	activeSecsThresholdStr    = "<1m"
)

func ShowActiveTask(db *sql.DB, writer io.Writer, template string) error {
	activeTaskDetails, err := pers.FetchActiveTask(db)
	if err != nil {
		return err
	}

	if activeTaskDetails.TaskID == -1 {
		return nil
	}

	now := time.Now()
	timeSpent := now.Sub(activeTaskDetails.LastLogEntryBeginTS).Seconds()
	var timeSpentStr string
	if timeSpent <= activeSecsThreshold {
		timeSpentStr = activeSecsThresholdStr
	} else {
		timeSpentStr = types.HumanizeDuration(int(timeSpent))
	}

	activeStr := strings.Replace(template, ActiveTaskPlaceholder, activeTaskDetails.TaskSummary, 1)
	activeStr = strings.Replace(activeStr, ActiveTaskTimePlaceholder, timeSpentStr, 1)
	fmt.Fprint(writer, activeStr)
	return nil
}
