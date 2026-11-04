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
)

func ShowActiveTask(db *sql.DB, writer io.Writer, template string) error {
	activeTaskDetails, err := pers.FetchActiveTaskDetails(db)
	if err != nil {
		return err
	}

	if activeTaskDetails.TaskID == -1 {
		return nil
	}

	timeSpentStr := types.HumanizeDuration(int(time.Since(activeTaskDetails.CurrentLogBeginTS).Seconds()))

	activeStr := strings.Replace(template, ActiveTaskPlaceholder, activeTaskDetails.TaskSummary, 1)
	activeStr = strings.Replace(activeStr, ActiveTaskTimePlaceholder, timeSpentStr, 1)
	fmt.Fprint(writer, activeStr)
	return nil
}
