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
	dbPath string
	db     *sql.DB
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
	Short: "Output reports based on tasks/log entries.",
	Long: `Output reports based on tasks/log entries.

Available reports:
  tasks     outputs a report of time spent on tasks
  log       outputs a report of the last few saved task log entries
  24h       outputs a report of log entries from the last 24h
  3d        outputs a report of log entries from the last 3 days (from beginning of day) (default)
  7d        outputs a report of log entries from the last 7 days (from beginning of day)
`,
	ValidArgs: []string{"tasks", "log", "24h", "3d", "7d"},
	Args:      cobra.MatchAll(cobra.MaximumNArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		out := os.Stdout

		if len(args) == 0 {
			ui.Render3DReport(db, out)
			return
		}
		switch args[0] {
		case "tasks":
			ui.RenderTaskReport(db, out)
		case "log":
			ui.RenderTaskLogReport(db, out)
		case "24h":
			ui.Render24hReport(db, out)
		case "3d":
			ui.Render3DReport(db, out)
		case "7d":
			ui.Render7DReport(db, out)
		}
	},
}

var activeCmd = &cobra.Command{
	Use:   "active",
	Short: "Shows task being actively tracked by \"hours\".",
	Run: func(cmd *cobra.Command, args []string) {
		ui.ShowActiveTask(db, os.Stdout)
	},
}

func init() {
	currentUser, err := user.Current()

	if err != nil {
		die("Error getting your home directory, This is a fatal error; use --dbpath to specify database path manually.\n%s\n", err)
	}

	defaultDBPath := fmt.Sprintf("%s/hours.v%s.db", currentUser.HomeDir, "1")
	rootCmd.PersistentFlags().StringVarP(&dbPath, "dbpath", "d", defaultDBPath, "location where hours should create its DB file")

	rootCmd.AddCommand(reportCmd)
	rootCmd.AddCommand(activeCmd)

	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		die("Something went wrong: %s\n", err)
	}
}
