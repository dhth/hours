package cmd

import (
	"bufio"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"io/fs"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	c "github.com/dhth/hours/internal/common"
	pers "github.com/dhth/hours/internal/persistence"
	"github.com/dhth/hours/internal/ui"
	"github.com/spf13/cobra"
)

const (
	defaultDBName     = "hours.db"
	configDirName     = "hours"
	themeDirName      = "themes"
	numDaysThreshold  = 30
	numTasksThreshold = 20

	themeEnvVar      = "HOURS_THEME"
	defaultThemeName = "default"
	warningColor     = "#fb4934"
)

//go:embed static/sample-theme.txt
var sampleThemeConfig string

var (
	errCouldntGetHomeDir         = errors.New("couldn't get home directory")
	errCouldntGetConfigDir       = errors.New("couldn't get config directory")
	errDBFileExtIncorrect        = errors.New("db file needs to end with .db")
	errCouldntCreateDBDirectory  = errors.New("couldn't create directory for database")
	errCouldntCreateDB           = errors.New("couldn't create database")
	errCouldntInitializeDB       = errors.New("couldn't initialize database")
	errCouldntOpenDB             = errors.New("couldn't open database")
	errCouldntReadTheme          = errors.New("couldn't read theme file")
	errCouldntLoadTheme          = errors.New("couldn't load theme file")
	errCouldntGenerateData       = errors.New("couldn't generate dummy data")
	errNumDaysExceedsThreshold   = errors.New("number of days exceeds threshold")
	errNumTasksExceedsThreshold  = errors.New("number of tasks exceeds threshold")
	errCouldntReadInput          = errors.New("couldn't read input")
	errIncorrectCodeEntered      = errors.New("incorrect code entered")
	errCouldntListThemes         = errors.New("couldn't list themes in config directory")
	errCouldntCheckIfThemeExists = errors.New("couldn't check if theme already exists")
	errThemeAlreadyExists        = errors.New("theme already exists")
	ErrThemeDoesntExist          = errors.New("theme doesn't exist")

	msgReportIssue = fmt.Sprintf("This isn't supposed to happen; let %s know about this error via \n%s.", c.Author, c.RepoIssuesURL)
)

func Execute() error {
	rootCmd, err := NewRootCommand()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		if errors.Is(err, errCouldntGetHomeDir) {
			fmt.Printf("\n%s\n", msgReportIssue)
		}
		return err
	}

	err = rootCmd.Execute()
	if err != nil {
		handleErrors(err)
	}
	return err
}

func setupDB(dbPathFull string) (*sql.DB, error) {
	var db *sql.DB
	var err error

	_, err = os.Stat(dbPathFull)
	if errors.Is(err, fs.ErrNotExist) {

		dir := filepath.Dir(dbPathFull)
		err = os.MkdirAll(dir, 0o755)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", errCouldntCreateDBDirectory, err.Error())
		}

		db, err = pers.GetDB(dbPathFull)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", errCouldntCreateDB, err.Error())
		}

		err = pers.InitDB(db)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", errCouldntInitializeDB, err.Error())
		}
		err = pers.UpgradeDB(db, 1)
		if err != nil {
			return nil, err
		}
	} else {
		db, err = pers.GetDB(dbPathFull)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", errCouldntOpenDB, err.Error())
		}
		err = pers.UpgradeDBIfNeeded(db)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}

func getTheme(themesDir string, themeName string) (ui.Style, error) {
	var zero ui.Style
	var theme ui.Theme
	if themeName == defaultThemeName {
		return ui.NewStyle(ui.DefaultTheme()), nil
	}

	themeFilePath := path.Join(themesDir, fmt.Sprintf("%s.json", themeName))
	themeBytes, err := os.ReadFile(themeFilePath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return zero, fmt.Errorf("%w: %s", ErrThemeDoesntExist, themeName)
		}
		return zero, fmt.Errorf("%w: %s", errCouldntReadTheme, err.Error())
	}
	if theme, err = ui.LoadTheme(themeBytes); err != nil {
		return zero, fmt.Errorf("%w: %w", errCouldntLoadTheme, err)
	}

	return ui.NewStyle(theme), nil
}

func NewRootCommand() (*cobra.Command, error) {
	var (
		userHomeDir         string
		userConfigDir       string
		themesDir           string
		dbPath              string
		dbPathFull          string
		db                  *sql.DB
		themeName           string
		style               ui.Style
		reportAgg           bool
		recordsInteractive  bool
		recordsOutputPlain  bool
		activeTemplate      string
		genNumDays          uint8
		genNumTasks         uint8
		genSkipConfirmation bool
	)

	preRun := func(cmd *cobra.Command, _ []string) error {
		dbPathFull = expandTilde(dbPath, userHomeDir)
		if filepath.Ext(dbPathFull) != ".db" {
			return errDBFileExtIncorrect
		}

		var err error
		db, err = setupDB(dbPathFull)
		switch {
		case errors.Is(err, errCouldntCreateDB):
			fmt.Fprintf(os.Stderr, `Couldn't create hours' local database.
%s

`, msgReportIssue)
		case errors.Is(err, errCouldntInitializeDB):
			fmt.Fprintf(os.Stderr, `Couldn't initialise hours' local database.
%s

`, msgReportIssue)
			// cleanup
			cleanupErr := os.Remove(dbPathFull)
			if cleanupErr != nil {
				fmt.Fprintf(os.Stderr, `Failed to remove hours' database file as well (at %s). Remove it manually.
Clean up error: %s

`, dbPathFull, cleanupErr.Error())
			}
		case errors.Is(err, errCouldntOpenDB):
			fmt.Fprintf(os.Stderr, `Couldn't open hours' local database.
%s

`, msgReportIssue)
		case errors.Is(err, pers.ErrCouldntFetchDBVersion):
			fmt.Fprintf(os.Stderr, `Couldn't get hours' latest database version.
%s

`, msgReportIssue)
		case errors.Is(err, pers.ErrDBDowngraded):
			fmt.Fprintf(os.Stderr, `Looks like you downgraded hours. You should either delete hours' database file (you
will lose data by doing that), or upgrade hours to the latest version.

`)
		case errors.Is(err, pers.ErrDBMigrationFailed):
			fmt.Fprintf(os.Stderr, `Something went wrong migrating hours' database.

You can try running hours by passing it a custom database file path (using
--db-path; this will create a new database) to see if that fixes things. If that
works, you can either delete the previous database, or keep using this new
database (both are not ideal).

%s
Sorry for breaking the upgrade step!

---

`, msgReportIssue)
		}

		if err != nil {
			return err
		}

		if !cmd.Flags().Changed("theme") {
			themeFromEnv := strings.TrimSpace(os.Getenv(themeEnvVar))
			if themeFromEnv != "" {
				themeName = themeFromEnv
			}
		}

		if style, err = getTheme(themesDir, themeName); err != nil {
			return err
		}

		return nil
	}

	rootCmd := &cobra.Command{
		Use:   "hours",
		Short: "\"hours\" is a no-frills time tracking toolkit for the command line",
		Long: `"hours" is a no-frills time tracking toolkit for the command line.

You can use "hours" to track time on your tasks, or view logs, reports, and
summary statistics for your tracked time.
`,
		SilenceUsage: true,
		PreRunE:      preRun,
		RunE: func(_ *cobra.Command, _ []string) error {
			return ui.RenderUI(db, style)
		},
	}

	generateCmd := &cobra.Command{
		Use:   "gen",
		Short: "Generate dummy log entries (helpful for beginners)",
		Long: `Generate dummy log entries.
This is intended for new users of 'hours' so they can get a sense of its
capabilities without actually tracking any time. It's recommended to always use
this with a --dbpath/-d flag that points to a throwaway database.
`,
		PreRunE: preRun,
		RunE: func(_ *cobra.Command, _ []string) error {
			if genNumDays > numDaysThreshold {
				return fmt.Errorf("%w (%d)", errNumDaysExceedsThreshold, numDaysThreshold)
			}
			if genNumTasks > numTasksThreshold {
				return fmt.Errorf("%w (%d)", errNumTasksExceedsThreshold, numTasksThreshold)
			}

			if !genSkipConfirmation {
				fmt.Print(lipgloss.NewStyle().Foreground(lipgloss.Color(warningColor)).Render(`
WARNING: You shouldn't run 'gen' on hours' actively used database as it'll
create dummy entries in it. You can run it on a throwaway database by passing a
path for it via --dbpath/-d (use it for all further invocations of 'hours' as
well).
`))
				fmt.Print(`
The 'gen' subcommand is intended for new users of 'hours' so they can get a
sense of its capabilities without actually tracking any time.

---

`)
				confirm, err := getConfirmation()
				if err != nil {
					return err
				}
				if !confirm {
					return fmt.Errorf("%w", errIncorrectCodeEntered)
				}
			}

			genErr := ui.GenerateData(db, genNumDays, genNumTasks)
			if genErr != nil {
				return fmt.Errorf("%w: %s", errCouldntGenerateData, genErr.Error())
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
			return nil
		},
	}

	reportCmd := &cobra.Command{
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
		Args:    cobra.MaximumNArgs(1),
		PreRunE: preRun,
		RunE: func(_ *cobra.Command, args []string) error {
			var period string
			if len(args) == 0 {
				period = "3d"
			} else {
				period = args[0]
			}

			return ui.RenderReport(db, style, os.Stdout, recordsOutputPlain, period, reportAgg, recordsInteractive)
		},
	}

	logCmd := &cobra.Command{
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
		Args:    cobra.MaximumNArgs(1),
		PreRunE: preRun,
		RunE: func(_ *cobra.Command, args []string) error {
			var period string
			if len(args) == 0 {
				period = "today"
			} else {
				period = args[0]
			}

			return ui.RenderTaskLog(db, style, os.Stdout, recordsOutputPlain, period, recordsInteractive)
		},
	}

	statsCmd := &cobra.Command{
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
		Args:    cobra.MaximumNArgs(1),
		PreRunE: preRun,
		RunE: func(_ *cobra.Command, args []string) error {
			var period string
			if len(args) == 0 {
				period = "3d"
			} else {
				period = args[0]
			}

			return ui.RenderStats(db, style, os.Stdout, recordsOutputPlain, period, recordsInteractive)
		},
	}

	themesCmd := &cobra.Command{
		Use:   "themes",
		Short: "Generate or view hours' themes",
	}

	addThemeCmd := &cobra.Command{
		Use:   "add <THEME_NAME>",
		Short: "Add a UI theme for hours",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			themeToCreate := args[0]
			themePath := filepath.Join(themesDir, fmt.Sprintf("%s.json", themeToCreate))
			fileInfo, err := os.Stat(themePath)
			if err == nil && !fileInfo.IsDir() {
				return fmt.Errorf("%w (at %q)", errThemeAlreadyExists, themePath)
			} else if !errors.Is(err, fs.ErrNotExist) {
				return errCouldntCheckIfThemeExists
			}

			themeFilePath, err := addTheme(themeToCreate, themesDir)
			if err != nil {
				return err
			}

			fmt.Printf(`Theme file with default values created at: %s

You can edit it as per your liking.
`, themeFilePath)

			return nil
		},
	}

	listThemesCmd := &cobra.Command{
		Use:   "list",
		Short: "List themes set up for hours",
		RunE: func(_ *cobra.Command, _ []string) error {
			var themes []string
			walkErr := filepath.Walk(themesDir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					if errors.Is(err, fs.ErrNotExist) {
						return nil
					}

					return err
				}
				ext := filepath.Ext(path)
				if !info.IsDir() && ext == ".json" {
					base := filepath.Base(path)
					themes = append(themes, strings.TrimSuffix(base, ext))
				}
				return nil
			})

			if walkErr != nil {
				return fmt.Errorf("%w: %s", errCouldntListThemes, walkErr.Error())
			}

			if len(themes) == 0 {
				fmt.Println("no themes configured (run \"hours theme add\" to add one)")
				return nil
			}

			fmt.Printf("%s\n", strings.Join(themes, "\n"))
			return nil
		},
	}

	sampleThemeCmd := &cobra.Command{
		Use:   "show-sample",
		Short: "Show a sample theme config",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Printf(`A sample theme config looks like the following.
You can choose to provide only the attributes you want to change.

%s
`,
				sampleThemeConfig)
		},
	}

	activeCmd := &cobra.Command{
		Use:   "active",
		Short: "Show the task being actively tracked by \"hours\"",
		Long: `Show the task being actively tracked by "hours".

You can pass in a template using the --template/-t flag, which supports the
following placeholders:

  {{task}}:  for the task summary
  {{time}}:  for the time spent so far on the active log entry

eg. hours active -t ' {{task}} ({{time}}) '
`,
		PreRunE: preRun,
		RunE: func(_ *cobra.Command, _ []string) error {
			return ui.ShowActiveTask(db, os.Stdout, activeTemplate)
		},
	}

	var err error
	userHomeDir, err = os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errCouldntGetHomeDir, err.Error())
	}

	userConfigDir, err = os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errCouldntGetConfigDir, err.Error())
	}

	themesDir = filepath.Join(userConfigDir, configDirName, themeDirName)

	defaultDBPath := filepath.Join(userHomeDir, defaultDBName)
	rootCmd.Flags().StringVarP(&dbPath, "dbpath", "d", defaultDBPath, "location of hours' database file")
	rootCmd.Flags().StringVarP(&themeName, "theme", "t", defaultThemeName, "UI theme to use")

	generateCmd.Flags().Uint8Var(&genNumDays, "num-days", 30, "number of days to generate fake data for")
	generateCmd.Flags().Uint8Var(&genNumTasks, "num-tasks", 10, "number of tasks to generate fake data for")
	generateCmd.Flags().BoolVarP(&genSkipConfirmation, "yes", "y", false, "to skip confirmation")
	generateCmd.Flags().StringVarP(&dbPath, "dbpath", "d", defaultDBPath, "location of hours' database file")

	reportCmd.Flags().BoolVarP(&reportAgg, "agg", "a", false, "whether to aggregate data by task for each day in report")
	reportCmd.Flags().BoolVarP(&recordsInteractive, "interactive", "i", false, "whether to view report interactively")
	reportCmd.Flags().BoolVarP(&recordsOutputPlain, "plain", "p", false, "whether to output report without any formatting")
	reportCmd.Flags().StringVarP(&dbPath, "dbpath", "d", defaultDBPath, "location of hours' database file")
	reportCmd.Flags().StringVarP(&themeName, "theme", "t", defaultThemeName,
		fmt.Sprintf("UI theme to use; themes live in %s", themesDir))

	logCmd.Flags().BoolVarP(&recordsOutputPlain, "plain", "p", false, "whether to output logs without any formatting")
	logCmd.Flags().BoolVarP(&recordsInteractive, "interactive", "i", false, "whether to view logs interactively")
	logCmd.Flags().StringVarP(&dbPath, "dbpath", "d", defaultDBPath, "location of hours' database file")
	logCmd.Flags().StringVarP(&themeName, "theme", "t", defaultThemeName,
		fmt.Sprintf("UI theme to use; themes live in %s", themesDir))

	statsCmd.Flags().BoolVarP(&recordsOutputPlain, "plain", "p", false, "whether to output stats without any formatting")
	statsCmd.Flags().BoolVarP(&recordsInteractive, "interactive", "i", false, "whether to view stats interactively")
	statsCmd.Flags().StringVarP(&dbPath, "dbpath", "d", defaultDBPath, "location of hours' database file")
	statsCmd.Flags().StringVarP(&themeName, "theme", "t", defaultThemeName,
		fmt.Sprintf("UI theme to use; themes live in %s", themesDir))

	activeCmd.Flags().StringVarP(&activeTemplate, "template", "t", ui.ActiveTaskPlaceholder, "string template to use for outputting active task")
	activeCmd.Flags().StringVarP(&dbPath, "dbpath", "d", defaultDBPath, "location of hours' database file")

	themesCmd.AddCommand(addThemeCmd)
	themesCmd.AddCommand(listThemesCmd)
	themesCmd.AddCommand(sampleThemeCmd)

	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(reportCmd)
	rootCmd.AddCommand(logCmd)
	rootCmd.AddCommand(statsCmd)
	rootCmd.AddCommand(activeCmd)
	rootCmd.AddCommand(themesCmd)

	rootCmd.CompletionOptions.DisableDefaultCmd = true

	return rootCmd, nil
}

func getRandomChars(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz"

	var code string
	for i := 0; i < length; i++ {
		code += string(charset[rand.Intn(len(charset))])
	}
	return code
}

func getConfirmation() (bool, error) {
	code := getRandomChars(2)
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Type %s to proceed: ", code)

	response, err := reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("%w: %s", errCouldntReadInput, err.Error())
	}
	response = strings.TrimSpace(response)

	return response == code, nil
}
