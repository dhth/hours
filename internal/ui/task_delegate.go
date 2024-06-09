package ui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

func newItemDelegate(color lipgloss.Color) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.Styles.SelectedTitle = d.Styles.
		SelectedTitle.
		Foreground(color).
		BorderLeftForeground(color)

	d.Styles.SelectedDesc = d.Styles.
		SelectedTitle

	return d
}
