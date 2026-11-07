#!/usr/bin/env bash
# Drives hours through a set of scenarios with termctrl and saves a PNG and a
# plain-text dump of each screen. scenarios.tsv lists each scenario's file name
# and description, in order.
#
# usage: capture.sh <hours-binary> <out-dir>
#
# Needs termctrl on PATH and the FiraCode Nerd Font Mono font installed. If a
# step fails, the screen at that point is saved as FAILED-<step>.{png,txt} and
# the script exits non-zero.

# Backticks in scenario descriptions are Markdown, not command substitution.
# shellcheck disable=SC2016

set -euo pipefail

if [ $# -ne 2 ]; then
	echo "usage: $0 <hours-binary> <out-dir>" >&2
	exit 2
fi

HOURS="$(cd "$(dirname "$1")" && pwd)/$(basename "$1")"
OUT=$2
FONT='FiraCode Nerd Font Mono'
SESSION="hours-screenshots-$$"
TMP=$(mktemp -d)
DB="$TMP/hours.db"
STEP=setup

# Same seed, same generated tasks and log entries on every run.
export HOURS_GEN_SEED=${HOURS_GEN_SEED:-42}

command -v termctrl >/dev/null || {
	echo "error: termctrl not found on PATH" >&2
	exit 1
}

# termctrl silently falls back to another font, so check for it here.
font_installed() {
	if command -v fc-match >/dev/null; then
		fc-match -f '%{family}\n' "$FONT" | tr ',' '\n' | grep -Fqx "$FONT"
	else
		ls ~/Library/Fonts/FiraCodeNerdFontMono-* /Library/Fonts/FiraCodeNerdFontMono-* >/dev/null 2>&1
	fi
}
font_installed || {
	echo "error: font \"$FONT\" is not installed" >&2
	exit 1
}

mkdir -p "$OUT"
echo "termctrl: $(termctrl --version)"

save() { # save <name> [termctrl save args...]
	local name=$1
	shift
	termctrl save --font-family "$FONT" --format png --format txt --out "$OUT/$name" "$@" >/dev/null
}

finish() {
	local status=$?
	if [ $status -ne 0 ]; then
		echo "error: step \"$STEP\" failed" >&2
		save "FAILED-$STEP" "$SESSION" 2>/dev/null || true
	fi
	termctrl stop "$SESSION" >/dev/null 2>&1 || true
	rm -rf "$TMP"
	exit $status
}
trap finish EXIT

step() { # step <name> <description>
	STEP=$1
	printf '%s\t%s\n' "$1" "$2" >>"$OUT/scenarios.tsv"
}
shot() { save "$STEP" "$SESSION"; }
send() { termctrl send "$SESSION" "$@" >/dev/null; }
wait_for() { termctrl wait "$SESSION" "$1" >/dev/null; }
tui() { termctrl start "$SESSION" --color always --cols "$1" --rows "$2" -- "$HOURS" -d "$DB" "${@:3}" >/dev/null; }
cli() { # one-off capture of a non-interactive command
	save "$STEP" --color always --cols 140 --rows 30 -- "$HOURS" -d "$DB" "$@"
}

"$HOURS" -d "$DB" gen -y --num-days 7 --num-tasks 6 >/dev/null

tui 120 36

step 01-tasks-list "Tasks list with generated data"
wait_for "6 tasks"
shot

step 02-add-task-form "Add a task: form filled in"
send text:a
wait_for "Add task"
send "text:write screenshot tests"
shot

step 03-task-added "Add a task: new task in the list"
send enter
wait_for "7 tasks"
shot

step 04-timer-started "Start a timer"
send text:s
wait_for "tracking:"
shot

step 05-active-while-tracking '`hours active` while a timer runs'
cli active

# Stopping the timer opens the entry form with the comment focused.
step 06-stop-timer-entry-form "Stop the timer: entry form, begin time moved back 30m"
send text:s
wait_for "Task Log Entry"
send "text:Drafted termctrl scenarios for hours"
send tab text:KKKKKK # focus begin time, move it back 30 minutes
wait_for "You're recording 30m"
shot

step 07-task-logs-list "Task logs list after saving the entry"
send ctrl-s
wait_for "7 tasks"
send text:2
wait_for "Task Logs (last"
shot

step 08-task-log-details "Task log details"
send text:d
wait_for "Task Log Details"
shot

step 09-help "Help view"
send esc
wait_for "Task Logs (last" # without this, esc followed by ? arrives as alt+?
send text:?
wait_for "Keyboard Shortcuts"
shot

termctrl stop "$SESSION" >/dev/null

step 10-log-today '`hours log today`'
cli log today

step 11-report-3d '`hours report 3d`'
cli report 3d

step 12-stats-week '`hours stats week`'
cli stats week

step 13-report-interactive '`hours report week -i`'
tui 140 24 report week -i
wait_for "go backwards"
shot

step 14-report-interactive-prev-week '`hours report week -i`, previous week (h)'
send text:h
sleep 0.5 # the previous week's dates depend on today, so there's no fixed text to wait for
shot

STEP=finished
echo "saved screenshots to $OUT"
