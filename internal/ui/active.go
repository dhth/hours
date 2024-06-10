package ui

import (
	"database/sql"
	"fmt"
	"io"
	"os"
)

func ShowActiveTask(db *sql.DB, writer io.Writer) {
	activeTaskDetails, err := fetchActiveTaskFromDB(db)

	if err != nil {
		fmt.Fprintf(os.Stdout, "Something went wrong:\n%s", err)
		os.Exit(1)
	}

	if activeTaskDetails.taskId == -1 {
		return
	}

	fmt.Fprintf(writer, "%s", activeTaskDetails.taskSummary)
}
