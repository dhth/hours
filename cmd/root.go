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
	reportNumDays    int
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
	Run: func(cmd *cobra.Command, args []string) {
		out := os.Stdout

		if reportNumDays <= 0 || reportNumDays > 7 {
			die("--num-days/-n needs to be between [1-7] (both inclusive)")
		}

		if reportAgg {
			ui.RenderNDaysReportAgg(db, out, reportNumDays, reportOrLogPlain)
		} else {
			ui.RenderNDaysReport(db, out, reportNumDays, reportOrLogPlain)
		}
	},
}

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Output task log entries",
	Long: `Output task log entries

Accepts an argument, which can be one of the following:

  all:   all recent log entries (in reverse chronological order)
  today: for log entries from today
  yest:  for log entries from yesterday
  date:  for log entries from that day (eg. "2024/06/08")
  range: for log entries from that range (eg. "2024/06/08...2024/06/12")
    `,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			ui.RenderTaskLog(db, os.Stdout, reportOrLogPlain, "all")
		} else {
			if args[0] == "" {
				die("Time period shouldn't be empty\n")
			}
			ui.RenderTaskLog(db, os.Stdout, reportOrLogPlain, args[0])
		}
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
	reportCmd.Flags().IntVarP(&reportNumDays, "num-days", "n", 3, "number of days to gather data for")

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
