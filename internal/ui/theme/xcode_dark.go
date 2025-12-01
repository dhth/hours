package theme

const (
	themeNameXcodeDark = "xcode-dark"
)

func paletteXcodeDark() builtInThemePalette {
	return builtInThemePalette{
		primary:    "#ff7ab2",
		secondary:  "#4eb0cc",
		tertiary:   "#ff8170",
		quaternary: "#b281eb",
		foreground: "#292a30",
		text:       "#dfdfe0",
		subtext:    "#7f8c98",
		muted:      "#7f8c98",
		help:       "#d9c97c",
		info:       "#78c2b3",
		error:      "#ff7ab2",
		warn:       "#ffa14f",
		tasks: []string{
			"#6bdfff",
			"#83d9a2",
			"#a8c8ff",
			"#acf2e4",
			"#d0a8ff",
			"#ff9cac",
			"#ffc1a6",
			"#ffcc66",
		},
	}
}
