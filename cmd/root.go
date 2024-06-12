package cmd

import (
	"database/sql"
	"fmt"
	"os"
	"os/user"

	"github.com/dhth/hours/internal/ui"
	"github.com/spf13/cobra"
)

var (
	dbPath           string
	db               *sql.DB
	reportAgg        bool
	reportOrLogPlain bool
)

func die(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

var rootCmd = &cobra.Command{
	Use:   "hours",
	Short: "Track time on your tasks via a simple TUI.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if dbPath == "" {
			die("dbpath cannot be empty")
		}

		dbPathFull := expandTilde(dbPath)

		var err error
		db, err = setupDB(dbPathFull)
		if err != nil {
			die("Couldn't set up \"hours\"' local database. This is a fatal error; let @dhth know about this.\n%s\n", err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ui.RenderUI(db)
	},
}

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Output a report based on tasks/log entries",
	Long: `Output a report based on tasks/log entries.

Reports show time spent on tasks in the last n days. These can also be
aggregated (using -a) to consolidate all task entries and show the
cumulative time spent on each task per day.

Accepts an argument, which can be one of the following:

  today:     for today's report
  yest:      for yesterday's report
  3d:        for a report on the last 3 days (default)
  week:      for a report on the last 7 days
  date:      for a report on a specific date (eg. "2024/06/08")
  range:     for a report on a date range (eg. "2024/06/08...2024/06/12")

Note: If a task log continues past midnight in your local timezone, it
will be reported on the day it ends.
    `,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var period string
		if len(args) == 0 {
			period = "3d"
		} else {
			period = args[0]
		}

		ui.RenderReport(db, os.Stdout, reportOrLogPlain, period, reportAgg)
	},
}

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Output task log entries",
	Long: `Output task log entries

Accepts an argument, which can be one of the following:

  today:     for log entries from today
  yest:      for log entries from yesterday
  3d:        for log entries from the last 3 days (default)
  week:      for log entries from the last 7 days
  date:      for log entries from that date (eg. "2024/06/08")
  range:     for log entries from that date range (eg. "2024/06/08...2024/06/12")
  all:       for all recent log entries (in reverse chronological order)

Note: If a task log continues past midnight in your local timezone, it'll
appear in the log on the day it ends.
    `,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var period string
		if len(args) == 0 {
			period = "3d"
		} else {
			period = args[0]
		}

		ui.RenderTaskLog(db, os.Stdout, reportOrLogPlain, period)
	},
}

var activeCmd = &cobra.Command{
	Use:   "active",
	Short: "Show the task being actively tracked by \"hours\"",
	Run: func(cmd *cobra.Command, args []string) {
		ui.ShowActiveTask(db, os.Stdout)
	},
}

func init() {
	currentUser, err := user.Current()

	if err != nil {
		die("Error getting your home directory, This is a fatal error; use --dbpath to specify database path manually\n%s\n", err)
	}

	defaultDBPath := fmt.Sprintf("%s/hours.v%s.db", currentUser.HomeDir, "1")
	rootCmd.PersistentFlags().StringVarP(&dbPath, "dbpath", "d", defaultDBPath, "location of hours' database file")

	reportCmd.Flags().BoolVarP(&reportAgg, "agg", "a", false, "whether to aggregate data by task in report")
	reportCmd.Flags().BoolVarP(&reportOrLogPlain, "plain", "p", false, "whether to output report without any formatting")

	logCmd.Flags().BoolVarP(&reportOrLogPlain, "plain", "p", false, "whether to output log without any formatting")

	rootCmd.AddCommand(reportCmd)
	rootCmd.AddCommand(logCmd)
	rootCmd.AddCommand(activeCmd)

	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		die("Something went wrong: %s\n", err)
	}
}
