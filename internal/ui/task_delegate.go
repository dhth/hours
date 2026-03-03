package ui

import (
	"image/color"

	"charm.land/bubbles/v2/list"
)

func newItemDelegate(titleColor, descColor, selectedColor color.Color) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.Styles.NormalTitle = d.Styles.
		NormalTitle.
		Foreground(titleColor)

	d.Styles.NormalDesc = d.Styles.
		NormalDesc.
		Foreground(descColor)

	d.Styles.SelectedTitle = d.Styles.
		SelectedTitle.
		Foreground(selectedColor).
		BorderLeftForeground(selectedColor)

	d.Styles.SelectedDesc = d.Styles.
		SelectedTitle

	return d
}
