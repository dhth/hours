package cmd

import (
	"path/filepath"
	"strings"
)

func expandTilde(path string, homeDir string) string {
	if strings.HasPrefix(path, "~/") {
		return filepath.Join(homeDir, path[2:])
	}
	return path
}
