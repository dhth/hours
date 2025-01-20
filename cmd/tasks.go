package cmd

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"

	pers "github.com/dhth/hours/internal/persistence"
)

const (
	tasksLimit = 500
)

func renderTasks(db *sql.DB, writer io.Writer, limit uint) error {
	limitToUse := limit
	if limit > tasksLimit {
		limitToUse = tasksLimit
	}

	tasks, err := pers.FetchTasks(db, true, limitToUse)
	if err != nil {
		return err
	}

	if len(tasks) == 0 {
		return nil
	}

	result, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("%w: %s", errCouldntUnmarshalToJSON, err.Error())
	}

	fmt.Fprintln(writer, string(result))

	return nil
}
