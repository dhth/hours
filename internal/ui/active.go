package ui

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

const (
	ActiveTaskPlaceholder     = "{{task}}"
	ActiveTaskTimePlaceholder = "{{time}}"
	activeSecsThreshold       = 60
	activeSecsThresholdStr    = "<1m"
)

func ShowActiveTask(db *sql.DB, writer io.Writer, template string) {
	activeTaskDetails, err := fetchActiveTaskFromDB(db)

	if err != nil {
		fmt.Fprintf(os.Stdout, "Something went wrong:\n%s", err)
		os.Exit(1)
	}

	if activeTaskDetails.taskId == -1 {
		return
	}

	now := time.Now()
	timeSpent := now.Sub(activeTaskDetails.lastLogEntryBeginTs).Seconds()
	var timeSpentStr string
	if timeSpent <= activeSecsThreshold {
		timeSpentStr = activeSecsThresholdStr
	} else {
		timeSpentStr = humanizeDuration(int(timeSpent))
	}

	activeStr := strings.Replace(template, ActiveTaskPlaceholder, activeTaskDetails.taskSummary, 1)
	activeStr = strings.Replace(activeStr, ActiveTaskTimePlaceholder, timeSpentStr, 1)
	fmt.Fprint(writer, activeStr)
}
