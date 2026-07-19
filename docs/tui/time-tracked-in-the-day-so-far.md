# Time tracked in the day so far

## Purpose and scope

The TUI presents the amount of time tracked during the user's current local
calendar day as a lightweight footer summary. It combines completed and active
tracking time and appears across TUI views once the total reaches one minute.

This metric is a day-progress signal, not a reporting primitive. It does not
provide a saved-versus-active breakdown and is not intended to replace the app's
reporting flows.

## Domain constraints

`Today` means the user's local calendar day, beginning at local midnight. It is
not a rolling 24-hour period. Calendar-day boundaries must respect the current
location and daylight-saving transitions.

Only elapsed time that overlaps the interval from local midnight to the current
time belongs to the total. This applies equally to completed and active task
logs. In particular, logs that cross midnight contribute only their portion
within the current day.

## Architectural direction

The daily total is derived TUI state. Its inputs are:

- the recent completed task logs already cached by the TUI;
- active tracking state;
- the current time.

The TUI does not fetch a separate database aggregate for this metric. The
summary and the rest of the interface therefore share the same in-memory view of
task-log data.

This choice keeps the feature proportional to its role as a convenience summary.
It avoids a second source of truth and avoids coordinating aggregate refreshes
with every task-log mutation.

Rendering consumes the derived value without calculating it or reading the
clock. This keeps rendering deterministic and free of side effects.

## Time-driven updates

Model changes keep the total synchronized with task-log mutations, but they are
not sufficient while the user is idle. Active tracking continues to accrue, and
the local day can change at midnight without any user input.

The TUI therefore maintains a lightweight periodic update signal for its
lifetime. The signal performs no database I/O; it only gives the model an
opportunity to derive the total using the latest time. It remains active even
when nothing is currently being tracked so that midnight transitions are
reflected automatically.

A coarse cadence is sufficient because the footer displays minute-level
progress. Keeping one lifetime update mechanism is intentionally preferred over
switching between active-tracking updates and a separate midnight wake-up.

## Completeness tradeoff

The total is based on a bounded window of recent task logs rather than all
matching rows in the database. This means it can undercount on unusually
high-volume days and can remain stale when the database changes outside the
running TUI.

That limitation is accepted in exchange for:

- simpler state management;
- consistency between the footer and the visible TUI state;
- deterministic rendering;
- no periodic database work;
- no mutation-specific aggregate refresh orchestration.

Reporting flows remain the authoritative path when completeness is required.

## Architectural invariants

- Interpret today as the current local calendar day.
- Count only elapsed time that overlaps the current day.
- Include both completed and active tracking time.
- Derive the total from existing TUI state rather than maintaining it
    incrementally.
- Keep rendering free of clock reads and state changes.
- Keep a time-driven update active for the lifetime of the TUI.
- Do not perform database I/O solely to refresh this footer metric.
- Treat the value as a convenience summary, not an authoritative aggregate.
