package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var integration = os.Getenv("INTEGRATION")

func skipIntegration(t *testing.T) {
	t.Helper()
	if integration != "1" {
		t.Skip("Skipping integration tests")
	}
}

func TestCLI(t *testing.T) {
	skipIntegration(t)

	tempDir, err := os.MkdirTemp("", "")
	require.NoErrorf(t, err, "error creating temporary directory: %s", err)

	binPath := filepath.Join(tempDir, "hours")
	buildArgs := []string{"build", "-o", binPath, "../.."}

	c := exec.Command("go", buildArgs...)
	err = c.Run()
	require.NoErrorf(t, err, "error building binary: %s", err)

	defer func() {
		err := os.RemoveAll(tempDir)
		if err != nil {
			fmt.Printf("couldn't clean up temporary directory (%s): %s", binPath, err)
		}
	}()

	t.Run("TestHelp", func(t *testing.T) {
		// GIVEN
		// WHEN
		c := exec.Command(binPath, "-h")
		b, err := c.CombinedOutput()

		// THEN
		assert.NoError(t, err, "output:\n%s", b)
	})

	t.Run("TestGen", func(t *testing.T) {
		// GIVEN
		// WHEN
		dbPath := filepath.Join(tempDir, "db.db")
		c := exec.Command(binPath, "gen", "-y", "-d", dbPath)
		b, err := c.CombinedOutput()

		// THEN
		assert.NoError(t, err, "output:\n%s", b)
	})

	t.Run("TestListTasks", func(t *testing.T) {
		// GIVEN
		dbPath := filepath.Join(tempDir, "db.db")
		c := exec.Command(binPath, "gen", "-y", "-d", dbPath)
		b, err := c.CombinedOutput()
		require.NoError(t, err, "couldn't generate dummy data, output:\n\n%s", b)

		// WHEN
		lc := exec.Command(binPath, "tasks", "-l", "3", "-d", dbPath)
		o, err := lc.CombinedOutput()

		// THEN
		require.NoError(t, err, "output:\n%s", o)
		var js interface{}
		err = json.Unmarshal(o, &js)
		assert.NoError(t, err, "output is not valid json; output:\n%s", o)
	})
}
