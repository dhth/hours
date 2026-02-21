package cli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

const runCmdTimeout = 10 * time.Second

type Fixture struct {
	tempDir string
	binPath string
}

type HoursCmd struct {
	args  []string
	useDB bool
	env   map[string]string
}

func NewCmd(args []string) HoursCmd {
	return HoursCmd{
		args: args,
		env:  make(map[string]string),
	}
}

func (c *HoursCmd) AddArgs(args ...string) {
	c.args = append(c.args, args...)
}

func (c *HoursCmd) SetEnv(key, value string) {
	c.env[key] = value
}

func (c *HoursCmd) UseDB() {
	c.useDB = true
}

func NewFixture(t *testing.T, binPath string) Fixture {
	t.Helper()

	tempDir := t.TempDir()

	return Fixture{
		tempDir: tempDir,
		binPath: binPath,
	}
}

func (f Fixture) RunCmd(cmd HoursCmd) (string, error) {
	argsToUse := cmd.args
	if cmd.useDB {
		dbPath := filepath.Join(f.tempDir, "hours.db")
		argsToUse = append(argsToUse, "--dbpath", dbPath)
	}
	ctx, cancel := context.WithTimeout(context.Background(), runCmdTimeout)
	defer cancel()

	cmdToRun := exec.CommandContext(ctx, f.binPath, argsToUse...)

	cmdToRun.Env = []string{
		fmt.Sprintf("HOME=%s", f.tempDir),
		fmt.Sprintf("PATH=%s", os.Getenv("PATH")),
	}
	for key, value := range cmd.env {
		cmdToRun.Env = append(cmdToRun.Env, fmt.Sprintf("%s=%s", key, value))
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	cmdToRun.Stdout = &stdoutBuf
	cmdToRun.Stderr = &stderrBuf

	err := cmdToRun.Run()
	exitCode := 0
	success := true

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return "", fmt.Errorf("command timed out after %s", runCmdTimeout)
		}

		if exitError, ok := errors.AsType[*exec.ExitError](err); ok {
			success = false
			exitCode = exitError.ExitCode()
		} else {
			return "", fmt.Errorf("couldn't run command: %s", err.Error())
		}
	}

	output := fmt.Sprintf(`success: %t
exit_code: %d
----- stdout -----
%s
----- stderr -----
%s
`, success, exitCode, stdoutBuf.String(), stderrBuf.String())

	return output, nil
}
