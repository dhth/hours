package theme

const (
	themeNameGruvboxDark = "gruvbox-dark"
)

func paletteGruvboxDark() builtInThemePalette {
	return builtInThemePalette{
		primary:    "#fe8019",
		secondary:  "#83a598",
		tertiary:   "#b8bb26",
		quaternary: "#d3869b",
		foreground: "#282828",
		text:       "#ebdbb2",
		subtext:    "#a89984",
		muted:      "#928374",
		help:       "#fabd2f",
		info:       "#8ec07c",
		error:      "#fb4934",
		warn:       "#d79921",
		tasks: []string{
			"#7ec8e3",
			"#98d4bb",
			"#a3d9a5",
			"#c9b1ff",
			"#d5c4a1",
			"#e09f9f",
			"#f4a261",
			"#f5d67b",
		},
	}
}
