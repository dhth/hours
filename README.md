# hours

[![Build Workflow Status](https://img.shields.io/github/actions/workflow/status/dhth/hours/main.yml?style=flat-square)](https://github.com/dhth/hours/actions/workflows/main.yml)
[![Vulncheck Workflow Status](https://img.shields.io/github/actions/workflow/status/dhth/hours/vulncheck.yml?style=flat-square&label=vulncheck)](https://github.com/dhth/hours/actions/workflows/vulncheck.yml)
[![Latest Release](https://img.shields.io/github/release/dhth/hours.svg?style=flat-square)](https://github.com/dhth/hours/releases/latest)
[![Commits Since Latest Release](https://img.shields.io/github/commits-since/dhth/hours/latest?style=flat-square)](https://github.com/dhth/hours/releases)

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
- edit task logs
- view task log details
- deactivate/activate a task
- view historical task log entries

![Usage](https://github.com/user-attachments/assets/16e34df0-fab3-42d9-a183-c8a07af06cca)

![Usage](https://github.com/user-attachments/assets/1213b61b-498a-4840-9ba3-f17097801b9d)

![Usage](https://github.com/user-attachments/assets/d804cc05-53d0-4740-ac53-6f8bf636be6c)

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

    today      for today's report
    yest       for yesterday's report
    3d         for a report on the last 3 days (default)
    week       for a report on the current week
    date       for a report for a specific date (eg. "2024/06/08")
    range      for a report for a date range (eg. "2024/06/08...2024/06/12", "2024/06/08...today", "2024/06/08..."; shouldn't be greater than 7 days)

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

    today      for log entries from today (default)
    yest       for log entries from yesterday
    3d         for log entries from the last 3 days
    week       for log entries from the current week
    date       for log entries from a specific date (eg. "2024/06/08")
    range      for log entries for a date range (eg. "2024/06/08...2024/06/12", "2024/06/08...today", "2024/06/08...")

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

    today      show stats for today
    yest       show stats for yesterday
    3d         show stats for the last 3 days (default)
    week       show stats for the current week
    date       show stats for a specific date (eg. "2024/06/08")
    range      show stats for a date range (eg. "2024/06/08...2024/06/12", "2024/06/08...today", "2024/06/08...")
    all        show stats for all log entries

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

🎨 Custom Themes
---

`hours` supports custom themes for its user interface (for the TUI and the
output of the `logs`, `report`, and `stats` commands. New themes can be added
using `hours themes add`, which will create a JSON file in `hours`' config
directory. You can then tweak this file as per your liking.

A sample theme config looks like the following. Colors codes can be provided in
ANSI 16, ANSI 256, or HEX formats. You can choose to provide only the attributes
you want to change.

```text
{
  "activeTask": "#8ec07c",                   # color for the active task in the footer
  "activeTaskBeginTime": "#d3869b",          # color for the active task begin time in the footer
  "activeTasks": "#fe8019",                  # primary color for the active task list view
  "formContext": "#fabd2f",                  # color for the context message in all forms
  "formFieldName": "#8ec07c",                # color for field names in all forms
  "formHelp": "#928374",                     # color for the help text in all forms
  "helpMsg": "#83a598",                      # color for help messages in the footer
  "helpPrimary": "#83a598",                  # primary color for the help view
  "helpSecondary": "#bdae93",                # secondary color for the help view
  "inactiveTasks": "#928374",                # primary color for the inactive task list view
  "initialHelpMsg": "#a58390",               # color of the initial help message in the footer
  "listItemDesc": "#777777",                 # color to be used for the title of list items (when they're not selected)
  "listItemTitle": "#dddddd",                # color to be used for the title of list items (when they're not selected)
  "recordsBorder": "#665c54",                # color for the table border in the output of logs, reports, stats
  "recordsDateRange": "#fabd2f",             # color for the data range picker in the output of logs, reports, stats
  "recordsFooter": "#ef8f62",                # color for the footer row in the output of logs, reports, stats
  "recordsHeader": "#d85d5d",                # color for the header row in the output of logs, reports, stats
  "recordsHelp": "#928374",                  # color for the help message in the output of logs, reports, stats
  "taskLogDetails": "#d3869b",               # primary color for the task log details view
  "taskEntry": "#8ec07c",                    # primary color for the task entry view
  "taskLogEntry": "#fabd2f",                 # primary color for the task log entry view
  "taskLogList": "#b8bb26",                  # primary color for the task log list view
  "taskLogFormInfo": "#d3869b",              # color to use for contextual information in the task log form
  "taskLogFormWarn": "#fe8019",              # color to use for contextual warnings in the task log form
  "taskLogFormError": "#fb4934",             # color to use for contextual errors in the task log form
  "tasks": [                                 # colors to be used for tasks in the output of logs, report, stats
    "#d3869b",
    "#b5e48c",
    "#90e0ef",
    "#ca7df9",
    "#ada7ff",
    "#bbd0ff",
    "#48cae4",
    "#8187dc",
    "#ffb4a2",
    "#b8bb26",
    "#ffc6ff",
    "#4895ef",
    "#83a598",
    "#fabd2f"
  ],
  "titleForeground": "#282828",              # foreground color to use for the title of all views
  "toolName": "#fe8019",                     # color for the tool name in the footer
  "tracking": "#fabd2f"                      # color for the tracking message in the footer
}
```

You can view configured themes using `hours themes list`.

Running hours with the `--theme <THEME_NAME>` flag will load up that theme.
Alternatively, you can set `$HOURS_THEME` to the theme name so you don't have to
pass the flag every time.

Here's a sampling of custom themes in action.

| Theme          | Preview                                                                                            |
|----------------|----------------------------------------------------------------------------------------------------|
| Solarized Dark | ![solarized-dark](https://github.com/user-attachments/assets/f68c0863-c45f-41d9-be2a-395f768b43ea) |
| Monokai        | ![monokai](https://github.com/user-attachments/assets/42e1ed59-b9be-42c3-953c-553bd94ff8e2)        |
| Nord           | ![nord](https://github.com/user-attachments/assets/407d54f3-e48a-4c08-8688-f19058e4c373)           |
| Dracula        | ![dracula](https://github.com/user-attachments/assets/854273e9-be0c-4457-bb19-86a9e1a04434)        |
| Gruvbox        | ![gruvbox](https://github.com/user-attachments/assets/b15982fb-0597-4457-940f-0e90b0d2cc06)        |
| Catppuccin     | ![catppuccin](https://github.com/user-attachments/assets/2dfdd9ec-7a87-4d18-819f-f5135b77fb23)     |
| Tokyonight     | ![tokyonight](https://github.com/user-attachments/assets/21ebe806-3159-4c5d-abbc-5405ef75087b)     |

📋 TUI Reference Manual
---

`hours` has 6 views:

  - Tasks List View                       Shows active tasks
  - Task Management View                  Shows a form to create/update tasks
  - Task Logs List View                   Shows your task logs
  - Task Log Details View                 Shows details for a task log
  - Inactive Tasks List View              Shows inactive tasks
  - Task Log Entry View                   Shows a form to save/update a task log entry
  - Help View

### Keyboard Shortcuts

#### General

| Shortcut      | Action                             |
|---------------|------------------------------------|
| `1`           | Switch to Tasks List View          |
| `2`           | Switch to Task Logs List View      |
| `3`           | Switch to Inactive Tasks List View |
| `<tab>`       | Go to next view/form entry         |
| `<shift+tab>` | Go to previous view/form entry     |
| `q`/`<esc>`   | Go back or quit                    |
| `<ctrl+c>`    | Quit immediately                   |
| `?`           | Show help view                     |

#### General List Controls

| Shortcut      | Action              |
|---------------|---------------------|
| `k`/`<Up>`    | Move cursor up      |
| `j`/`<Down>`  | Move cursor down    |
| `h`/`<Left>`  | Go to previous page |
| `l`/`<Right>` | Go to next page     |
| `<ctrl+r>`    | Refresh list        |

#### Task List View

| Shortcut   | Action                                                                                                                 |
|------------|------------------------------------------------------------------------------------------------------------------------|
| `a`        | Add a task                                                                                                             |
| `u`        | Update task details                                                                                                    |
| `s`        | Start/stop recording time on a task; stopping will open up the "Task Log Entry View"                                   |
| `S`        | Quick switch recording; will save a task log entry for the currently active task, and start recording time for another |
| `f`        | Finish the currently active task log without comment                                                                   |
| `<ctrl+s>` | Edit the currently active task log/Add a new manual task log entry                                                     |
| `<ctrl+x>` | Discard currently active recording                                                                                     |
| `<ctrl+t>` | Go to currently tracked item                                                                                           |
| `<ctrl+d>` | Deactivate task                                                                                                        |

#### Task Logs List View

*Note: `~` at the end of a task log comment indicates that it has more lines that are not visible in the list view*

| Shortcut       | Action                |
|----------------|-----------------------|
| `d`            | Show task log details |
| `<ctrl+s>`/`u` | Update task log entry |
| `<ctrl+d>`     | Delete task log entry |

#### Task Log Details View

| Shortcut | Action               |
|----------|----------------------|
| `h`      | Go to previous entry |
| `l`      | Go to next entry     |

#### Inactive Task List View

| Shortcut   | Action        |
|------------|---------------|
| `<ctrl+d>` | Activate task |

#### Task Log Entry View

| Shortcut           | Action                                   |
|--------------------|------------------------------------------|
| `enter`/`<ctrl+s>` | Save entered details for the task log    |
| `k`                | Move timestamp backwards by one minute   |
| `j`                | Move timestamp forwards by one minute    |
| `K`                | Move timestamp backwards by five minutes |
| `J`                | Move timestamp forwards by five minutes  |
| `h`                | Move timestamp backwards by a day        |
| `l`                | Move timestamp forwards by a day         |

Acknowledgements
---

`hours` is built using [bubbletea][1], and is released using [goreleaser][2],
both of which are amazing tools.

[1]: https://github.com/charmbracelet/bubbletea
[2]: https://github.com/goreleaser/goreleaser
