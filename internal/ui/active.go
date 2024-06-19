package ui

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	ActiveTaskPlaceholder = "{{task}}"
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

	activeStr := strings.Replace(template, ActiveTaskPlaceholder, activeTaskDetails.taskSummary, 1)
	fmt.Fprint(writer, activeStr)
}
