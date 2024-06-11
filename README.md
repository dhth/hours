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

I wanted to keep track of the time I spend on side projects and other
non-day-job activities. I also wanted to be able to generate plaintext reports
of the time tracked, so I could get a general sense of how much of my time was
being spent on what. I couldn't find a tool that precisely fit these needs, so I
decided to build one myself.

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

This subcommand accepts a `-p` flag, which can be anything in the range [1-7]
(both inclusive) to see reports for the last "n" days (including today).

```
hours report -h

Output reports based on tasks/log entries.

Usage:
  hours report [flags]

Flags:
  -a, --agg            whether to aggregate data by task in report
  -n, --num-days int   number of days to gather data for (default 3)
  -p, --plain          whether to output report without any formatting
```

```bash
hours report
# or
hours report -n=7
```

### Logs

As the name suggests, logs are just that: list of task entries you've saved
using `hours`. This subcommand accepts an argument, which can be one of the following:

- `all`:     all recent log entries (in reverse chronological order)
- `today`:   for log entries from today
- `yest`:    for log entries from yesterday
- `<date>`:  for log entries from that day
- `<range>`: for log entries from in that range

```bash
hours log today
# or
hours log 2024/06/08
# or
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
