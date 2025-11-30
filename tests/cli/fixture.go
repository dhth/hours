package cli

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

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

func NewFixture() (Fixture, error) {
	var zero Fixture
	tempDir, err := os.MkdirTemp("", "")
	if err != nil {
		return zero, fmt.Errorf("couldn't create temporary directory: %s", err.Error())
	}

	binPath := filepath.Join(tempDir, "hours")
	buildArgs := []string{"build", "-o", binPath, "../../.."}

	c := exec.Command("go", buildArgs...)
	buildOutput, err := c.CombinedOutput()
	if err != nil {
		cleanupErr := os.RemoveAll(tempDir)
		if cleanupErr != nil {
			fmt.Fprintf(os.Stderr, "couldn't clean up temporary directory (%s): %s", tempDir, cleanupErr.Error())
		}

		return zero, fmt.Errorf(`couldn't build binary: %s
output:
%s`, err.Error(), buildOutput)
	}

	return Fixture{
		tempDir: tempDir,
		binPath: binPath,
	}, nil
}

func (f Fixture) Cleanup() error {
	err := os.RemoveAll(f.tempDir)
	if err != nil {
		return fmt.Errorf("couldn't clean up temporary directory (%s): %s", f.tempDir, err.Error())
	}

	return nil
}

func (f Fixture) RunCmd(cmd HoursCmd) (string, error) {
	argsToUse := cmd.args
	if cmd.useDB {
		dbPath := filepath.Join(f.tempDir, "hours.db")
		argsToUse = append(argsToUse, "--dbpath", dbPath)
	}
	cmdToRun := exec.Command(f.binPath, argsToUse...)

	cmdToRun.Env = os.Environ()
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
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
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
