package cmd

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"math/rand"
	"os"
	"os/user"
	"strings"

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
	genNumDays         uint8
	genNumTasks        uint8
)

func die(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func setupDB() {

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
}

var rootCmd = &cobra.Command{
	Use:   "hours",
	Short: "\"hours\" is a no-frills time tracking toolkit for the command line",
	Long: `"hours" is a no-frills time tracking toolkit for the command line.

You can use "hours" to track time on your tasks, or view logs, reports, and
summary statistics for your tracked time.
`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if cmd.CalledAs() == "gen" {
			return
		}
		setupDB()
	},
	Run: func(cmd *cobra.Command, args []string) {
		ui.RenderUI(db)
	},
}

var generateCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate dummy log entries",
	Long: `Generate dummy log entries.
This is intended for new users of 'hours' so they can get a sense of its
capabilities without actually tracking any time. It's recommended to always use
this with a --dbpath/-d flag that points to a throwaway database.
`,
	Run: func(cmd *cobra.Command, args []string) {
		if genNumDays > 30 {
			die("Maximum value for number of days is 30")
		}
		if genNumTasks > 20 {
			die("Maximum value for number of days is 20")
		}

		dbPathFull := expandTilde(dbPath)

		_, statErr := os.Stat(dbPathFull)
		if statErr == nil {
			die(`A file already exists at %s. Either delete it, or use a different path.

Tip: 'gen' should always be used on a throwaway database file.`, dbPathFull)
		}

		fmt.Print(ui.WarningStyle.Render(`
WARNING: You shouldn't run 'gen' on hours' actively used database as it'll
create dummy entries in it. You can run it out on a throwaway database by
passing a path for it via --dbpath/-d (use it for all further invocations of
'hours' as well).
`))
		fmt.Print(`
The 'gen' subcommand is intended for new users of 'hours' so they can get a
sense of its capabilities without actually tracking any time.

---

`)
		confirm := getConfirmation()
		if !confirm {
			fmt.Printf("\nIncorrect code; exiting\n")
			os.Exit(1)
		}

		setupDB()
		genErr := ui.GenerateData(db, genNumDays, genNumTasks)
		if genErr != nil {
			die(`Something went wrong generating dummy data.
let %s know about this via %s.

Error: %s`, author, repoIssuesUrl, genErr)
		}
		fmt.Printf(`
Successfully generated dummy data in the database file: %s

If this is not the default database file path, use --dbpath/-d with 'hours' when
you want to access the dummy data.

Go ahead and try the following!

hours --dbpath=%s
hours --dbpath=%s report week -i
hours --dbpath=%s log today -i
hours --dbpath=%s stats today -i
`, dbPath, dbPath, dbPath, dbPath, dbPath)
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

	generateCmd.Flags().Uint8Var(&genNumDays, "num-days", 30, "number of days to generate fake data for")
	generateCmd.Flags().Uint8Var(&genNumTasks, "num-tasks", 10, "number of tasks to generate fake data for")

	reportCmd.Flags().BoolVarP(&reportAgg, "agg", "a", false, "whether to aggregate data by task for each day in report")
	reportCmd.Flags().BoolVarP(&recordsInteractive, "interactive", "i", false, "whether to view report interactively")
	reportCmd.Flags().BoolVarP(&recordsOutputPlain, "plain", "p", false, "whether to output report without any formatting")

	logCmd.Flags().BoolVarP(&recordsOutputPlain, "plain", "p", false, "whether to output logs without any formatting")
	logCmd.Flags().BoolVarP(&recordsInteractive, "interactive", "i", false, "whether to view logs interactively")

	statsCmd.Flags().BoolVarP(&recordsOutputPlain, "plain", "p", false, "whether to output stats without any formatting")
	statsCmd.Flags().BoolVarP(&recordsInteractive, "interactive", "i", false, "whether to view stats interactively")

	activeCmd.Flags().StringVarP(&activeTemplate, "template", "t", ui.ActiveTaskPlaceholder,
		fmt.Sprintf("string template to use for outputting active task; use \"%s\" as placeholder for the task", ui.ActiveTaskPlaceholder))

	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(reportCmd)
	rootCmd.AddCommand(logCmd)
	rootCmd.AddCommand(statsCmd)
	rootCmd.AddCommand(activeCmd)

	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		die("Something went wrong: %s", err)
	}
}

func getRandomChars(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz"

	var code string
	for i := 0; i < length; i++ {
		code += string(charset[rand.Intn(len(charset))])
	}
	return code
}

func getConfirmation() bool {

	code := getRandomChars(2)
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Type %s to proceed: ", code)

	response, err := reader.ReadString('\n')
	if err != nil {
		die("Something went wrong reading input: %s", err)
	}
	response = strings.TrimSpace(response)

	return response == code
}
