package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type reportStyles struct {
	headerStyle lipgloss.Style
	footerStyle lipgloss.Style
	borderStyle lipgloss.Style
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

func isCommentValid(comment string) bool {
	return strings.TrimSpace(comment) != ""
}
