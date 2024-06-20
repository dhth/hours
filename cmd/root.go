package cmd

import (
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/user"

	"github.com/dhth/hours/internal/ui"
	"github.com/spf13/cobra"
)

const (
	author        = "@dhth"
	repoIssuesUrl = "https://github.com/dhth/hours/issues"
)

var (
	dbPath             string
	db                 *sql.DB
	reportAgg          bool
	recordsInteractive bool
	recordsOutputPlain bool
	activeTemplate     string
)

func die(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

var rootCmd = &cobra.Command{
	Use:   "hours",
	Short: "\"hours\" is a no-frills time tracking toolkit for the command line",
	Long: `"hours" is a no-frills time tracking toolkit for the command line.

You can use "hours" to track time on your tasks, or view logs, reports, and
summary statistics for your tracked time.
`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if dbPath == "" {
			die("dbpath cannot be empty")
		}

		dbPathFull := expandTilde(dbPath)

		var err error

		_, err = os.Stat(dbPathFull)
		if errors.Is(err, fs.ErrNotExist) {
			db, err = getDB(dbPathFull)

			if err != nil {
				die(`Couldn't create hours' local database. This is a fatal error;
let %s know about this via %s.

Error: %s`,
					author,
					repoIssuesUrl,
					err)
			}

			err = initDB(db)
			if err != nil {
				die(`Couldn't create hours' local database. This is a fatal error;
let %s know about this via %s.

Error: %s`,
					author,
					repoIssuesUrl,
					err)
			}
			upgradeDB(db, 1)
		} else {
			db, err = getDB(dbPathFull)
			if err != nil {
				die(`Couldn't open hours' local database. This is a fatal error;
let %s know about this via %s.

Error: %s`,
					author,
					repoIssuesUrl,
					err)
			}
			upgradeDBIfNeeded(db)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ui.RenderUI(db)
	},
}

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Output a report based on task log entries",
	Long: `Output a report based on task log entries.

Reports show time spent on tasks per day in the time period you specify. These
can also be aggregated (using -a) to consolidate all task entries and show the
cumulative time spent on each task per day.

Accepts an argument, which can be one of the following:

  today:     for today's report
  yest:      for yesterday's report
  3d:        for a report on the last 3 days (default)
  week:      for a report on the current week
  date:      for a report for a specific date (eg. "2024/06/08")
  range:     for a report for a date range (eg. "2024/06/08...2024/06/12")

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

		ui.RenderReport(db, os.Stdout, recordsOutputPlain, period, reportAgg, recordsInteractive)
	},
}

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Output task log entries",
	Long: `Output task log entries.

Accepts an argument, which can be one of the following:

  today:     for log entries from today (default)
  yest:      for log entries from yesterday
  3d:        for log entries from the last 3 days
  week:      for log entries from the current week
  date:      for log entries from a specific date (eg. "2024/06/08")
  range:     for log entries from a specific date range (eg. "2024/06/08...2024/06/12")

Note: If a task log continues past midnight in your local timezone, it'll
appear in the log for the day it ends.
    `,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var period string
		if len(args) == 0 {
			period = "today"
		} else {
			period = args[0]
		}

		ui.RenderTaskLog(db, os.Stdout, recordsOutputPlain, period, recordsInteractive)
	},
}

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Output statistics for tracked time",
	Long: `Output statistics for tracked time.

Accepts an argument, which can be one of the following:

  today:     show stats for today
  yest:      show stats for yesterday
  3d:        show stats for the last 3 days (default)
  week:      show stats for the current week
  month:     show stats for the current month
  date:      show stats for a specific date (eg. "2024/06/08")
  range:     show stats for a specific date range (eg. "2024/06/08...2024/06/12")
  all:       show stats for all log entries

Note: If a task log continues past midnight in your local timezone, it'll
be considered in the stats for the day it ends.
    `,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var period string
		if len(args) == 0 {
			period = "3d"
		} else {
			period = args[0]
		}

		ui.RenderStats(db, os.Stdout, recordsOutputPlain, period, recordsInteractive)
	},
}

var activeCmd = &cobra.Command{
	Use:   "active",
	Short: "Show the task being actively tracked by \"hours\"",
	Run: func(cmd *cobra.Command, args []string) {
		ui.ShowActiveTask(db, os.Stdout, activeTemplate)
	},
}

func init() {
	currentUser, err := user.Current()

	if err != nil {
		die(`Couldn't get your home directory. This is a fatal error;
use --dbpath to specify database path manually
let %s know about this via %s.

Error: %s`, author, repoIssuesUrl, err)
	}

	defaultDBPath := fmt.Sprintf("%s/hours.db", currentUser.HomeDir)
	rootCmd.PersistentFlags().StringVarP(&dbPath, "dbpath", "d", defaultDBPath, "location of hours' database file")

	reportCmd.Flags().BoolVarP(&reportAgg, "agg", "a", false, "whether to aggregate data by task for each day in report")
	reportCmd.Flags().BoolVarP(&recordsInteractive, "interactive", "i", false, "whether to view report interactively")
	reportCmd.Flags().BoolVarP(&recordsOutputPlain, "plain", "p", false, "whether to output report without any formatting")

	logCmd.Flags().BoolVarP(&recordsOutputPlain, "plain", "p", false, "whether to output logs without any formatting")
	logCmd.Flags().BoolVarP(&recordsInteractive, "interactive", "i", false, "whether to view logs interactively")

	statsCmd.Flags().BoolVarP(&recordsOutputPlain, "plain", "p", false, "whether to output stats without any formatting")
	statsCmd.Flags().BoolVarP(&recordsInteractive, "interactive", "i", false, "whether to view stats interactively")

	activeCmd.Flags().StringVarP(&activeTemplate, "template", "t", ui.ActiveTaskPlaceholder,
		fmt.Sprintf("string template to use for outputting active task; use \"%s\" as placeholder for the task", ui.ActiveTaskPlaceholder))

	rootCmd.AddCommand(reportCmd)
	rootCmd.AddCommand(logCmd)
	rootCmd.AddCommand(statsCmd)
	rootCmd.AddCommand(activeCmd)

	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		die("Something went wrong: %s\n", err)
	}
}
