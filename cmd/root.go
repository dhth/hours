package cmd

import (
	"flag"
	"fmt"
	"os"
	"os/user"

	"github.com/dhth/hours/internal/ui"
)

func die(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func Execute() {
	currentUser, err := user.Current()

	if err != nil {
		die("Error getting your home directory, This is a fatal error; use -db-path to specify database path manually.\n%s\n", err)
	}

	defaultDBPath := fmt.Sprintf("%s/hours.v%s.db", currentUser.HomeDir, DB_VERSION)
	dbPath := flag.String("db-path", defaultDBPath, "location where hours should create its DB file")

	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, `Track time on your tasks via a simple TUI.

Usage:
  hours [flags] [command]

Flags:
`)
		flag.CommandLine.SetOutput(os.Stdout)
		flag.PrintDefaults()
		fmt.Fprintf(os.Stdout, `
Commands:
  7d
        outputs a report of log entries from the last 7 days
  3d
        outputs a report of log entries from the last 3 days
  24h
        outputs a report of log entries from the last 24h
  tasks
        outputs a report of time spent on tasks
  log
        outputs the last few saved task log entries
  active
        shows the task currently being tracked
`)
	}
	flag.Parse()

	if *dbPath == "" {
		die("db-path cannot be empty")
	}

	dbPathFull := expandTilde(*dbPath)

	db, err := setupDB(dbPathFull)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't set up hours' local database. This is a fatal error; let @dhth know about this.\n%s\n", err)
		os.Exit(1)
	}

	args := os.Args[1:]
	out := os.Stdout

	if len(args) > 0 {
		if args[0] == "tasks" {
			ui.RenderTaskReport(db, out)
		} else if args[0] == "24h" {
			ui.Render24hReport(db, out)
		} else if args[0] == "7d" {
			ui.Render7DReport(db, out)
		} else if args[0] == "3d" {
			ui.Render3DReport(db, out)
		} else if args[0] == "log" {
			ui.RenderTaskLogReport(db, out)
		} else if args[0] == "active" {
			ui.ShowActiveTask(db, out)
		}
	} else {
		ui.RenderUI(db)
	}
}
