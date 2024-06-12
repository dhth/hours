# hours

✨ Overview
---

`hours` is a no-frills command-line app for tracking time on tasks. It's
designed for users who want basic time tracking for their tasks/projects right
in the terminal. With a simple and minimalistic UI, almost everything in `hours`
can be achieved with one or two keypresses. It can also generate plaintext
reports and logs based on time tracked.

🤔 Motivation
---

For a while, I've been wanting to keep track of time I spend on side projects
and other non-day-job activities. I also wanted to be able to generate plain
text reports to get an overview of time allocation. All of this needed to be
done via a fast command line tool that prioritised ease of use over unnecessary
features. After I couldn't find a tool that precisely fit these needs, I decided
to build one myself.

💾 Install
---

**go**:

```sh
go install github.com/dhth/hours@latest
```

⚡️ Usage
---

Open the TUI by simply running `hours`. The TUI lets you do the following:

- create/update tasks
- start/stop tracking time on a task
- add manual task log entries
- deactivate/activate a task
- view historical task log entries

Besides a TUI, `hours` also offers reports and logs based on the time tracking
you do. These can be viewed using the subcommands `report` and `log`
respectively.

### Reports

Reports show time spent on tasks in the last `n` days. These can also be
aggregated (using `-a`) to consolidate all task entries and show the cumulative
time spent on each task per day.

```
hours report -h

Output a report based on tasks/log entries.

Reports show time spent on tasks in the last n days. These can also be
aggregated (using -a) to consolidate all task entries and show the
cumulative time spent on each task per day.

Accepts an argument, which can be one of the following:

  today:     for today's report
  yest:      for yesterday's report
  3d:        for a report on the last 3 days (default)
  week:      for a report on the last 7 days
  date:      for a report on a specific date (eg. "2024/06/08")
  range:     for a report on a date range (eg. "2024/06/08...2024/06/12")

Note: If a task log continues past midnight in your local timezone, it
will be reported on the day it ends.

Flags:
  -a, --agg     whether to aggregate data by task in report
  -p, --plain   whether to output report without any formatting
```

```bash
# see report from last 3 days
hours report

# see aggregated time spent on tasks
hours report -a

# see report for the 7 days
hours report week

# see report for a specific date range
hours report 2024/06/08...2024/06/12
```

### Logs

As the name suggests, logs are just that: list of task entries you've saved
using `hours`.

```
hours log -h

Output task log entries

Accepts an argument, which can be one of the following:

  today:     for log entries from today
  yest:      for log entries from yesterday
  3d:        for log entries from the last 3 days (default)
  week:      for log entries from the last 7 days
  date:      for log entries from that date (eg. "2024/06/08")
  range:     for log entries from that date range (eg. "2024/06/08...2024/06/12")
  all:       for all recent log entries (in reverse chronological order)

Note: If a task log continues past midnight in your local timezone, it'll
appear in the log on the day it ends.

Usage:
  hours log [flags]

Flags:
  -p, --plain   whether to output log without any formatting
```

```bash
# see log entries from today
hours log today

# see log entries from a specific day
hours log 2024/06/08

# see log entries from a specific date range
hours log 2024/06/08...2024/06/12
```

📋 TUI Reference Manual
---

```
"hours" has 5 panes:
  - Tasks List View                      Shows your tasks
  - Task Management View                 Allows you to create/update tasks
  - Task Log List View                   Shows your task log entries
  - Inactive Tasks List View             Shows inactive tasks
  - Help View (this one)

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
  /                                       Start filtering

Task List View

  s                                       Toggle recording time on the currently selected task,
                                              will open up a form to record a comment on the
                                              second "s" keypress
  <ctrl+s>                                Add a manual task log entry
  <ctrl+t>                                Go to currently tracked item
  <ctrl+d>                                Deactivate task

Task Log List View

  <ctrl+d>                                Delete task log entry
  <ctrl+r>                                Refresh list

Inactive Task List View

  <ctrl+d>                                Activate task

Task Log Entry View

  enter                                   Save task log entry
```

Acknowledgements
---

`hours` is built using the TUI framework [bubbletea][1].

[1]: https://github.com/charmbracelet/bubbletea
