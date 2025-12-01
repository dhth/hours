package theme

const (
	themeNameDracula = "dracula"
)

func paletteDracula() builtInThemePalette {
	return builtInThemePalette{
		primary:    "#ff6e6e",
		secondary:  "#50fa7b",
		tertiary:   "#bd93f9",
		quaternary: "#8be9fd",
		foreground: "#282a36",
		text:       "#ffffff",
		subtext:    "#f8f8f2",
		muted:      "#bd93f9",
		help:       "#f1fa8c",
		info:       "#69ff94",
		error:      "#ff5555",
		warn:       "#ffffa5",
		tasks: []string{
			"#6ecbff",
			"#7be0ad",
			"#a4ffff",
			"#c9a0dc",
			"#e8a0e8",
			"#ff79c6",
			"#ff92df",
			"#ffb86c",
		},
	}
}
