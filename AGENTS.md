# AGENTS.md

This file provides guidance to AI agents when working with code in this repository.

## Project Overview

`hours` is a no-frills CLI time tracking toolkit written in Go. It provides a TUI (Terminal User Interface) built with BubbleTea/Lipgloss and uses SQLite for persistence. Users track time on tasks, then generate plaintext reports, stats, and logs.

## Common Commands

All commands available via `justfile` aliases:

```bash
just all                          # Format + lint + test (all-in-one)
just fmt                          # Format: gofumpt -l -w .
just check                        # Lint: golangci-lint run
just build                        # Build: go build .
just test                         # Run tests (no cache): go test -count=1 ./...
just run                          # Run: go run .
just update-snapshots             # To run tests while updating snapshots
```

**Important**: always use the `just` recipes to invoke Go commands.

## Architecture

```
cmd/           CLI commands (Cobra). root.go is the main entry point.
internal/
  persistence/ SQLite database layer: schema init, migrations, all SQL queries
  types/       Core domain types (Task, TaskLogEntry) and date/duration helpers
  ui/          BubbleTea TUI: model/view/update split across files
    theme/     Theme system with customizable colors
  utils/       String utilities (trim, padding)
  common/      Shared constants
tests/         Integration tests (test.sh runs 28+ CLI tests)
```

**Data flow**: CLI commands (cmd/) → persistence layer (queries.go) → SQLite DB. The TUI (ui/) uses BubbleTea's Elm architecture (model → update → view) with async commands in cmds.go.

**Database**: Two main tables: `task` (projects) and `task_log` (time entries). A trigger prevents multiple simultaneously active task logs.

## Tests

Use the local skill `write-tests` (.agents/skills/write-tests/SKILL.md) to get more context.

## Key Conventions

- Linting: golangci-lint with revive rules, gofumpt formatting (not gofmt)
- Error handling: custom error types with `fmt.Errorf` wrapping
- DB migrations tracked in `db_versions` table (`internal/persistence/migrations.go`)
