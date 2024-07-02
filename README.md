# hours

`hours` is a no-frills time tracking toolkit for the command line.

It's designed for users who want basic time tracking for their tasks/projects
right in the terminal. With a simple and minimalistic UI, almost everything in
`hours` can be achieved with one or two keypresses. It can also generate
plaintext reports, summary statistics, and logs based on time tracked.

![Usage](https://tools.dhruvs.space/images/hours/hours.gif)

[Link to Video][2]

🤔 Motivation
---

For a while, I've been wanting to keep track of time I spend on side projects
and other non-day-job activities. I also wanted to be able to generate plain
text reports to get an overview of time allocation. All of this needed to be
done via a fast command line tool that prioritised ease of use over
fancy-but-ultimately-not-so-useful features. I couldn't find a tool that
precisely fit these needs, so I decided to build one for myself.

💾 Install
---

**homebrew**:

```sh
brew install dhth/tap/hours
```

**go**:

```sh
go install github.com/dhth/hours@latest
```

Or get the binaries directly from a
[release](https://github.com/dhth/hours/releases).

⚡️ Usage
---

> Newbie tip: If you want to see how `hours` works without having to track time,
> you can have it generate dummy data for you. See [here](#generate-dummy-data)
> for more details.

### TUI

Open the TUI by simply running `hours`. The TUI lets you do the following:

- create/update tasks
- start/stop tracking time on a task
- add manual task log entries
- deactivate/activate a task
- view historical task log entries

![Usage](https://tools.dhruvs.space/images/hours/tui-1.png)

![Usage](https://tools.dhruvs.space/images/hours/tui-2.png)

![Usage](https://tools.dhruvs.space/images/hours/tui-3.png)

Besides a TUI, `hours` also offers reports, statistics, and logs based on the
time tracking you do. These can be viewed using the subcommands `report`,
`stats`, and `log` respectively.

### Reports

```bash
hours report [flags] [arg]
```

Output a report based on task log entries.

Reports show time spent on tasks per day in the time period you specify. These
can also be aggregated (using `-a`) to consolidate all task entries and show the
cumulative time spent on each task per day.

Accepts an argument, which can be one of the following:

    today:     for today's report
    yest:      for yesterday's report
    3d:        for a report on the last 3 days (default)
    week:      for a report on the current week
    date:      for a report for a specific date (eg. "2024/06/08")
    range:     for a report for a date range (eg. "2024/06/08...2024/06/12")

*Note: If a task log continues past midnight in your local timezone, it will be
reported on the day it ends.*

![Usage](https://tools.dhruvs.space/images/hours/report-1.png)

Reports can also be viewed via an interactive interface using the
`--interactive`/`-i` flag.

![Usage](https://tools.dhruvs.space/images/hours/report-interactive-1.gif)

### Log

```bash
hours log [flags] [arg]
```

Output task log entries.

Accepts an argument, which can be one of the following:

    today:     for log entries from today (default)
    yest:      for log entries from yesterday
    3d:        for log entries from the last 3 days
    week:      for log entries from the current week
    date:      for log entries from a specific date (eg. "2024/06/08")
    range:     for log entries from a specific date range (eg. "2024/06/08...2024/06/12")

*Note: If a task log continues past midnight in your local timezone, it'll
appear in the log for the day it ends.*

![Usage](https://tools.dhruvs.space/images/hours/log-1.png)

Logs can also be viewed via an interactive interface using the
`--interactive`/`-i` flag.

![Usage](https://tools.dhruvs.space/images/hours/log-interactive-1.gif)


### Statistics

```bash
hours stats [flag] [arg]
```

Output statistics for tracked time.

Accepts an argument, which can be one of the following:

    today:     show stats for today
    yest:      show stats for yesterday
    3d:        show stats for the last 3 days (default)
    week:      show stats for the current week
    date:      show stats for a specific date (eg. "2024/06/08")
    range:     show stats for a specific date range (eg. "2024/06/08...2024/06/12")
    all:       show stats for all log entries

*Note: If a task log continues past midnight in your local timezone, it'll
be considered in the stats for the day it ends.*

![Usage](https://tools.dhruvs.space/images/hours/stats-1.png)

Stats can also be viewed via an interactive interface using the
`--interactive`/`-i` flag.

![Usage](https://tools.dhruvs.space/images/hours/stats-interactive-1.gif)

### Active Task

`hours` can show you the task being actively tracked using the `active`
subcommand. This subcommand supports the following placeholders using the
`--template`/`-t` flag:

    {{task}}:  for the task summary
    {{time}}:  for the time spent so far on the active log entry

Tip: This can be used to display the active task in tmux's (or similar terminal
multiplexers) status line using:

```
set -g status-right "#(hours active -t ' {{task}} ({{time}}) ')".
```

### Generate Dummy Data

You can have `hours` generate dummy data for you, so you can play around with
it, and see if its approach of showing reports/logs/stats works for you. You can
do so using the `gen` subcommand.

```bash
hours gen --dbpath=/var/tmp/throwaway.db
```


📋 TUI Reference Manual
---

```text
"hours" has 6 views:
  - Tasks List View                       Shows active tasks
  - Task Management View                  Shows a form to create/update tasks
  - Task Log List View                    Shows your task log entries
  - Inactive Tasks List View              Shows inactive tasks
  - Task Log Entry View                   Shows a form to save/update a task log entry
  - Help View

Keyboard Shortcuts

General

  1                                       Switch to Tasks List View
  2                                       Switch to Task Log List View
  3                                       Switch to Inactive Task Log List View
  <tab>                                   Go to next view/form entry
  <shift+tab>                             Go to previous view/form entry
    ?                                     Show help view

General List Controls

  h/<Up>                                  Move cursor up
  k/<Down>                                Move cursor down
  h<Left>                                 Go to previous page
  l<Right>                                Go to next page
  <ctrl+r>                                Refresh list

Task List View

  a                                       Add a task
  u                                       Update task details
  s                                       Toggle recording time on the currently selected task,
                                              will open up a form to record a task log entry on
                                              the second "s" keypress
  <ctrl+s>                                Edit the currently active task log/Add a new manual task log entry
  <ctrl+t>                                Go to currently tracked item
  <ctrl+d>                                Deactivate task

Task Log List View

  <ctrl+d>                                Delete task log entry
  <ctrl+r>                                Refresh list

Inactive Task List View

  <ctrl+d>                                Activate task

Task Log Entry View

  enter                                   Save task log entry
  k                                       Move timestamp backwards by one minute
  j                                       Move timestamp forwards by one minute
  K                                       Move timestamp backwards by five minutes
  J                                       Move timestamp forwards by five minutes
```

Acknowledgements
---

`hours` is built using [bubbletea][1], and is released using [goreleaser][2],
both of which are amazing tools.

[1]: https://github.com/charmbracelet/bubbletea
[2]: https://github.com/goreleaser/goreleaser
