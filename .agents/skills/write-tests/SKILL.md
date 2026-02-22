---
name: write-tests
description: Use this skill when creating or updating tests in this repository.
---

# Write Tests

## Scope

- Apply these patterns when adding or editing tests in this repository.
- Use existing test style and keep test changes focused on behavior.

## Commands

- Run tests with `just test`.
- Update snapshots with `just update-snapshots`.

## Testing Patterns

- Unit tests use `testify/assert` alongside source files
- UI components use snapshot testing (`go-snaps`); snapshots in `internal/ui/__snapshots__/`
- `internal/types` has a `TimeProvider` interface for mocking time in tests
- CLI integration tests use Go tests in `tests/cli` with snapshots in `tests/cli/__snapshots__/`

### CLI integration tests

- Keep integration tests in package `tests/cli` (file-per-command is preferred)
- `TestMain` builds the `hours` binary once per package run and is shared by all tests in the package
- `NewFixture(t, testBinaryPath)` gives each test a `t.TempDir()` workspace (auto-cleaned)
- Use `cmd.UseDB()` to append a per-test DB path under the fixture temp dir
- `RunCmd` executes with timeout and deterministic env (`HOME` set to fixture temp dir, `PATH` propagated, explicit overrides via `SetEnv`)

### Given-When-Then Structure

All tests in this codebase should follow the **Given-When-Then** (GWT) pattern
to improve readability and maintainability. This structure makes tests
self-documenting and helps both humans and AI agents quickly understand:
- What state is being set up (Given)
- What action is being performed (When)
- What outcome is expected (Then)

Format:

```go
func TestSomething(t *testing.T) {
    // GIVEN
    // ... setup code, test data, initial state

    // WHEN
    // ... the action being tested

    // THEN
    // ... assertions and expected outcomes
}
```

### Snapshot tests

- Use snapshot tests with `go-snaps` wherever it makes sense, especially for:
  - Complex data structures (serialization/deserialization)
  - Generated output (TUI frames, command line output, etc.)
  - Regression testing where exact output matters
