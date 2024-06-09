package ui

import (
	"fmt"
	"math"
	"strings"
	"time"
)

func RightPadTrim(s string, length int) string {
	if len(s) >= length {
		if length > 3 {
			return s[:length-3] + "..."
		}
		return s[:length]
	}
	return s + strings.Repeat(" ", length-len(s))
}

func Trim(s string, length int) string {
	if len(s) >= length {
		if length > 3 {
			return s[:length-3] + "..."
		}
		return s[:length]
	}
	return s
}

func humanizeDuration(durationInSecs int) string {
	duration := time.Duration(durationInSecs) * time.Second

	dSecs := duration.Seconds()
	dMins := duration.Minutes()
	dHours := duration.Hours()

	if dSecs < 60 {
		if dSecs == 1 {
			return "1 second"
		}
		return fmt.Sprintf("%d seconds", int(duration.Seconds()))
	}

	if duration.Minutes() < 60 {
		if dMins == 1 {
			return "1 minute"
		}
		return fmt.Sprintf("%d minutes", int(duration.Minutes()))
	}

	modMins := int(math.Mod(duration.Minutes(), 60))

	hourStr := "hours"
	modMinsStr := "minutes"

	if dHours == 1 {
		hourStr = "hour"
	}
	if modMins == 1 {
		modMinsStr = "minute"
	}

	if modMins == 0 {
		return fmt.Sprintf("%d %s", int(duration.Hours()), hourStr)
	}

	return fmt.Sprintf("%d %s %d %s", int(duration.Hours()), hourStr, modMins, modMinsStr)
}
