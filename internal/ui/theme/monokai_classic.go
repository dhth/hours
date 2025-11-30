package theme

const themeNameMonokaiClassic = "monokai-classic"

func paletteMonokaiClassic() builtInThemePalette {
	return builtInThemePalette{
		primary:    "#66d9ef",
		secondary:  "#ae81ff",
		tertiary:   "#a6e22e",
		quaternary: "#fd971f",
		foreground: "#272822",
		text:       "#fdfff1",
		subtext:    "#c0c1b5",
		muted:      "#57584f",
		help:       "#e6db74",
		info:       "#a6e22e",
		error:      "#f92672",
		warn:       "#fd971f",
		tasks: []string{
			"#78dce8",
			"#a9dc76",
			"#ab9df2",
			"#c4e88a",
			"#d4bfff",
			"#ff6a9e",
			"#ffb86c",
			"#ffd866",
		},
	}
}
