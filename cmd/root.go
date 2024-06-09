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
		die("Error getting your home directory, This is a fatal error; let @dhth know about this.\n%s\n", err)
	}

	defaultDBPath := fmt.Sprintf("%s/hours.v%s.db", currentUser.HomeDir, DB_VERSION)
	dbPath := flag.String("db-path", defaultDBPath, "location where hours should create its DB file")

	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, "Track time on your tasks.\n\nFlags:\n")
		flag.CommandLine.SetOutput(os.Stdout)
		flag.PrintDefaults()
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

	ui.RenderUI(db)

}
