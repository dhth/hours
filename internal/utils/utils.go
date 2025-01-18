package utils

import "strings"

func RightPadTrim(s string, length int, dots bool) string {
	if len(s) > length {
		if dots && length > 3 {
			return s[:length-3] + "..."
		}
		return s[:length]
	}
	return s + strings.Repeat(" ", length-len(s))
}

func Trim(s string, length int) string {
	if len(s) > length {
		if length > 3 {
			return s[:length-3] + "..."
		}
		return s[:length]
	}
	return s
}

func TrimWithMoreLinesIndicator(s string, length int) string {
	lines := strings.SplitN(s, "\n", 2)

	if len(lines) > 1 {
		if length <= 5 {
			return Trim(lines[0], length)
		}
		return Trim(lines[0], length-2) + " ~"
	}

	return Trim(lines[0], length)
}

func RightPadTrimWithMoreLinesIndicator(s string, length int) string {
	lines := strings.SplitN(s, "\n", 2)

	if len(lines) > 1 {
		if length <= 5 {
			return RightPadTrim(lines[0], length, true)
		}

		return RightPadTrim(Trim(lines[0], length-2)+" ~", length, false)
	}

	return RightPadTrim(lines[0], length, true)
}
