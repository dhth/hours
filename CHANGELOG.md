# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.5.0] - Feb 22, 2025

### Added

- Support for custom themes and subcommands to manage them

### Changed

- Go, dependency upgrades

## [v0.4.1] - Feb 03, 2025

### Changed

- Minor wording changes in error messages

## [v0.4.0] - Jan 19, 2025

### Added

- Time tracking can now be switched between tasks with a single keypress
- The active task log can now be edited before it's finished
- Task logs can now be edited after saving
- Adds a view for viewing task log details

### Changed

- Allow for longer task log comments
- Task log comments can now be empty

## [v0.3.0] - Jun 29, 2024

### Added

- Timestamps in the "Task Log Entry" view can be moved forwards/backwards using
  j/k/J/K
- The TUI now shows the start time of an active recording
- An active task log recording can now be cancelled

### Changed

- Timestamps in "Task Log" view show up differently based on the end timestamp
- "active" subcommand supports a time placeholder, eg. hours active -t 'working
  on {{task}} for {{time}}'

## [v0.2.0] - Jun 21, 2024

### Added

- Adds the ability to view reports/logs/stats interactively (using the
  --interactive/-i flag)
- Adds the "gen" subcommand to allow new users of "hours" to generate dummy data

[unreleased]: https://github.com/dhth/hours/compare/v0.5.0...HEAD
[v0.5.0]: https://github.com/dhth/hours/compare/v0.4.1...v0.5.0
[v0.4.1]: https://github.com/dhth/hours/compare/v0.4.0...v0.4.1
[v0.4.0]: https://github.com/dhth/hours/compare/v0.3.0...v0.4.0
[v0.3.0]: https://github.com/dhth/hours/compare/v0.2.0...v0.3.0
[v0.2.0]: https://github.com/dhth/hours/compare/v0.1.0...v0.2.0
