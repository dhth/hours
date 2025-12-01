package theme

const (
	themeNameGithubDark = "github-dark"
)

func paletteGithubDark() builtInThemePalette {
	return builtInThemePalette{
		primary:    "#f78166",
		secondary:  "#56d364",
		tertiary:   "#db61a2",
		quaternary: "#6ca4f8",
		foreground: "#101216",
		text:       "#ffffff",
		subtext:    "#8b949e",
		muted:      "#8b949e",
		help:       "#e3b341",
		info:       "#56d364",
		error:      "#db61a2",
		warn:       "#f78166",
		tasks: []string{
			"#79c0ff",
			"#7ee787",
			"#89dceb",
			"#a5d6ff",
			"#d2a8ff",
			"#ff9bce",
			"#ffa657",
			"#ffd33d",
		},
	}
}
