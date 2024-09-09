package cmd

import (
	"path/filepath"
	"strings"
)

func expandTilde(path string, homeDir string) string {
	pathWithoutTilde, found := strings.CutPrefix(path, "~/")
	if !found {
		return path
	}
	return filepath.Join(homeDir, pathWithoutTilde)
}
