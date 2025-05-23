package ui

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	pers "github.com/dhth/hours/internal/persistence"
)

const (
	nonEmptyCommentChance = 0.8
	longCommentChance     = 0.3
	sampleLongCommentBody = `

This is a sample task log comment. The comment can be used to record
additional information for a task log.

You can include:
- Detailed steps taken during the task
- Observations and notes
- Any issues encountered and how they were resolved
- Future actions or follow-ups required
- References to related tasks or documents

Use this section to ensure all relevant details are captured for each task,
providing a comprehensive log that can be referred to later.`
)

var (
	tasks = []string{
		".net",
		"assembly",
		"c",
		"c#",
		"c++",
		"clojure",
		"dart",
		"elixir",
		"erlang",
		"f#",
		"go",
		"haskell",
		"java",
		"javascript",
		"julia",
		"kotlin",
		"lisp",
		"lua",
		"ocaml",
		"objective-c",
		"php",
		"perl",
		"prolog",
		"python",
		"r",
		"roc",
		"ruby",
		"rust",
		"sql",
		"scala",
		"swift",
		"typescript",
		"zig",
	}
	verbs = []string{
		"write",
		"fix",
		"deploy",
		"review",
		"test",
		"refactor",
		"design",
		"implement",
		"document",
		"update",
		"create",
		"analyze",
		"optimize",
		"integrate",
		"configure",
		"build",
		"debug",
		"monitor",
		"automate",
		"maintain",
	}
	nouns = []string{
		"documentation",
		"tests",
		"code",
		"review",
		"feature",
		"bug",
		"module",
		"api",
		"interface",
		"function",
		"pipeline",
		"database",
		"service",
		"deployment",
		"configuration",
		"component",
		"report",
		"script",
		"workflow",
		"log",
	}
)

func GenerateData(db *sql.DB, numDays, numTasks uint8) error {
	for i := range numTasks {
		summary := tasks[rand.Intn(len(tasks))]
		_, err := pers.InsertTask(db, summary)
		if err != nil {
			return err
		}
		numLogs := int(numDays/2) + rand.Intn(int(numDays/2))
		for range numLogs {
			beginTs := randomTimestamp(int(numDays))
			numMinutes := 30 + rand.Intn(60)
			endTs := beginTs.Add(time.Minute * time.Duration(numMinutes))
			var comment *string
			commentStr := fmt.Sprintf("%s %s", verbs[rand.Intn(len(verbs))], nouns[rand.Intn(len(nouns))])
			if rand.Float64() < nonEmptyCommentChance {
				if rand.Float64() < longCommentChance {
					commentStr += sampleLongCommentBody
				}
				comment = &commentStr
			}

			_, err = pers.InsertManualTL(db, int(i+1), beginTs, endTs, comment)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func randomTimestamp(numDays int) time.Time {
	now := time.Now().Local()

	maxSeconds := numDays * 24 * 60 * 60
	randomSeconds := rand.Intn(maxSeconds)
	randomTime := now.Add(-time.Duration(randomSeconds) * time.Second)
	return randomTime
}
