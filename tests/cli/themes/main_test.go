package themes

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

var testBinaryPath string

func TestMain(m *testing.M) {
	tempDir, err := os.MkdirTemp("", "hours-cli-tests-*")
	if err != nil {
		panic(fmt.Sprintf("couldn't create temporary directory: %s", err.Error()))
	}

	binPath := filepath.Join(tempDir, "hours")
	buildCmd := exec.Command("go", "build", "-o", binPath, "../../..")
	buildOutput, err := buildCmd.CombinedOutput()
	if err != nil {
		_ = os.RemoveAll(tempDir)
		panic(fmt.Sprintf("couldn't build binary: %s\noutput:\n%s", err.Error(), buildOutput))
	}

	testBinaryPath = binPath
	code := m.Run()

	if err := os.RemoveAll(tempDir); err != nil {
		panic(fmt.Sprintf("couldn't clean up temporary directory (%s): %s", tempDir, err.Error()))
	}

	os.Exit(code)
}
