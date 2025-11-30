package theme

const themeNameTokyonight = "tokyonight"

func paletteTokyonight() builtInThemePalette {
	return builtInThemePalette{
		primary:     "#f7768e",
		secondary:   "#7aa2f7",
		tertiary:    "#9ece6a",
		quaternary: "#bb9af7",
		foreground:  "#1a1b26",
		text:        "#7aa2f7",
		subtext:     "#c0caf5",
		muted:       "#7dcfff",
		help:        "#e0af68",
		info:        "#9ece6a",
		error:       "#f7768e",
		warn:        "#e0af68",
		tasks: []string{
			"#2ac3de",
			"#73daca",
			"#7dcfff",
			"#b4f9f8",
			"#c0caf5",
			"#ff9e64",
		},
	}
}
