package theme

const themeNameMonokaiClassic = "monokai-classic"

func paletteMonokai() builtInThemePalette {
	return builtInThemePalette{
		primary:     "#66d9ef",
		secondary:   "#ae81ff",
		tertiary:    "#a6e22e",
		quaternary: "#fd971f",
		foreground:  "#272822",
		text:        "#fdfff1",
		subtext:     "#c0c1b5",
		muted:       "#c0c1b5",
		help:        "#e6db74",
		info:        "#a6e22e",
		error:       "#f92672",
		warn:        "#fd971f",
		tasks: []string{
			"#75e6da",
			"#9effff",
			"#c4e88a",
			"#d4bfff",
			"#e6db74",
			"#ff6a9e",
			"#ffb86c",
		},
	}
}
