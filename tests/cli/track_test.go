package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func skipIntegration(t *testing.T) {
	t.Helper()
	if os.Getenv("INTEGRATION") != "1" {
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

	t.Run("Help", func(t *testing.T) {
		// GIVEN
		// WHEN
		c := exec.Command(binPath, "-h")
		b, err := c.CombinedOutput()

		// THEN
		assert.NoError(t, err, "output:\n%s", b)
	})

	t.Run("Generate Data", func(t *testing.T) {
		// GIVEN
		// WHEN
		fileName := fmt.Sprintf("%s.db", uuid.New().String())
		dbPath := filepath.Join(tempDir, fileName)
		c := exec.Command(binPath, "gen", "-y", "-d", dbPath)
		b, err := c.CombinedOutput()

		// THEN
		assert.NoError(t, err, "output:\n%s", b)
	})

	t.Run("List tasks", func(t *testing.T) {
		// GIVEN
		fileName := fmt.Sprintf("%s.db", uuid.New().String())
		dbPath := filepath.Join(tempDir, fileName)
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

	t.Run("Start Tracking", func(t *testing.T) {
		// GIVEN
		fileName := fmt.Sprintf("%s.db", uuid.New().String())
		dbPath := filepath.Join(tempDir, fileName)
		c := exec.Command(binPath, "gen", "-y", "-d", dbPath)
		b, err := c.CombinedOutput()
		require.NoError(t, err, "couldn't generate dummy data, output:\n\n%s", b)

		// WHEN
		lc := exec.Command(binPath, "track", "start", "1", "-c", "comment goes here", "-d", dbPath)
		o, err := lc.CombinedOutput()

		// THEN
		require.NoError(t, err, "output:\n%s", o)
		var js interface{}
		err = json.Unmarshal(o, &js)
		assert.NoError(t, err, "output is not valid json; output:\n%s", o)
	})

	t.Run("Start Tracking fails if task ID argument is invalid", func(t *testing.T) {
		// GIVEN
		fileName := fmt.Sprintf("%s.db", uuid.New().String())
		dbPath := filepath.Join(tempDir, fileName)
		c := exec.Command(binPath, "gen", "-y", "-d", dbPath)
		b, err := c.CombinedOutput()
		require.NoError(t, err, "couldn't generate dummy data, output:\n\n%s", b)

		// WHEN
		lc := exec.Command(binPath, "track", "start", "blah", "-c", "comment goes here", "-d", dbPath)
		o, err := lc.CombinedOutput()

		// THEN
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			exitCode := exitError.ExitCode()
			require.Equal(t, 1, exitCode, "exit code is not correct: got %d, expected: 1; output:\n%s", exitCode, o)
			assert.Contains(t, string(o), "couldn't parse the argument for task ID as an integer:")
		} else {
			t.Fatalf("couldn't get error code")
		}
	})

	t.Run("Start tracking fails if task ID is non existent", func(t *testing.T) {
		// GIVEN
		fileName := fmt.Sprintf("%s.db", uuid.New().String())
		dbPath := filepath.Join(tempDir, fileName)
		c := exec.Command(binPath, "gen", "-y", "-d", dbPath)
		b, err := c.CombinedOutput()
		require.NoError(t, err, "couldn't generate dummy data, output:\n\n%s", b)

		// WHEN
		lc := exec.Command(binPath, "track", "start", "100000", "-c", "comment goes here", "-d", dbPath)
		o, err := lc.CombinedOutput()

		// THEN
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			exitCode := exitError.ExitCode()
			require.Equal(t, 1, exitCode, "exit code is not correct: got %d, expected: 1; output:\n%s", exitCode, o)
			assert.Contains(t, string(o), "task does not exist")
		} else {
			t.Fatalf("couldn't get error code")
		}
	})
}
