package cli

import (
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stretchr/testify/require"
)

func TestShowConfig(t *testing.T) {
	fx := NewFixture(t, testBinaryPath)

	commonArgs := []string{
		"themes",
		"show-config",
	}

	t.Run("help flag works", func(t *testing.T) {
		cmd := NewCmd(commonArgs)
		cmd.AddArgs("--help")

		result, err := fx.RunCmd(cmd)

		require.NoError(t, err)
		snaps.MatchStandaloneSnapshot(t, result)
	})

	t.Run("works for built-in theme", func(t *testing.T) {
		cmd := NewCmd(commonArgs)
		cmd.AddArgs("--theme", "monokai-classic")

		result, err := fx.RunCmd(cmd)

		require.NoError(t, err)
		snaps.MatchStandaloneSnapshot(t, result)
	})

	t.Run("fails for incorrect builtin theme", func(t *testing.T) {
		cmd := NewCmd(commonArgs)
		cmd.AddArgs("--theme", "unknown")

		result, err := fx.RunCmd(cmd)

		require.NoError(t, err)
		snaps.MatchStandaloneSnapshot(t, result)
	})

	t.Run("fails for incorrect builtin theme provided via env var", func(t *testing.T) {
		cmd := NewCmd(commonArgs)
		cmd.SetEnv("HOURS_THEME", "unknown")

		result, err := fx.RunCmd(cmd)

		require.NoError(t, err)
		snaps.MatchStandaloneSnapshot(t, result)
	})
}
