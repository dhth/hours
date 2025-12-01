package theme

const (
	themeNameNightOwl = "night-owl"
)

func paletteNightOwl() builtInThemePalette {
	return builtInThemePalette{
		primary:    "#22da6e",
		secondary:  "#82aaff",
		tertiary:   "#c792ea",
		quaternary: "#ef5350",
		foreground: "#011627",
		text:       "#ffffff",
		subtext:    "#d6deeb",
		muted:      "#d6deeb",
		help:       "#ffeb95",
		info:       "#22da6e",
		error:      "#ef5350",
		warn:       "#c792ea",
		tasks: []string{
			"#7fdbca",
			"#80cbc4",
			"#a3c4f3",
			"#b8e994",
			"#dbb2ff",
			"#ecc48d",
			"#f78c6c",
			"#ff9eb5",
		},
	}
}
