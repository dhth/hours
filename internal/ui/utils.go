package ui

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

type reportStyles struct {
	headerStyle lipgloss.Style
	footerStyle lipgloss.Style
	borderStyle lipgloss.Style
}

func RightPadTrim(s string, length int, dots bool) string {
	if len(s) >= length {
		if dots && length > 3 {
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

	if duration.Seconds() < 60 {
		return fmt.Sprintf("%ds", int(duration.Seconds()))
	}

	if duration.Minutes() < 60 {
		return fmt.Sprintf("%dm", int(duration.Minutes()))
	}

	modMins := int(math.Mod(duration.Minutes(), 60))

	if modMins == 0 {
		return fmt.Sprintf("%dh", int(duration.Hours()))
	}

	return fmt.Sprintf("%dh %dm", int(duration.Hours()), modMins)
}

func getReportStyles(plain bool) reportStyles {
	if plain {
		return reportStyles{
			emptyStyle,
			emptyStyle,
			emptyStyle,
		}
	}
	return reportStyles{
		recordsHeaderStyle,
		recordsFooterStyle,
		recordsBorderStyle,
	}
}
