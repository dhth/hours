package ui

import (
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
