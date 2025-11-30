package theme

const (
	themeNameCatppuccinMocha = "catppuccin-mocha"
)

func paletteCatppuccinMocha() builtInThemePalette {
	return builtInThemePalette{
		primary:    "#f37799",
		secondary:  "#74a8fc",
		tertiary:   "#a6e3a1",
		quaternary: "#f2aede",
		foreground: "#1e1e2e",
		text:       "#74a8fc",
		subtext:    "#a6adc8",
		muted:      "#a6adc8",
		help:       "#ebd391",
		info:       "#89d88b",
		error:      "#f37799",
		warn:       "#ebd391",
		tasks: []string{
			"#6bd7ca",
			"#89b4fa",
			"#94e2d5",
			"#cba6f7",
			"#eba0ac",
			"#f38ba8",
			"#f9e2af",
			"#fab387",
		},
	}
}
