package themes

import (
	"testing"

	"github.com/dhth/hours/tests/cli"
	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stretchr/testify/require"
)

func TestShowConfig(t *testing.T) {
	fx, err := cli.NewFixture()
	require.NoErrorf(t, err, "error setting up fixture: %s", err)

	defer func() {
		err := fx.Cleanup()
		require.NoErrorf(t, err, "error cleaning up fixture: %s", err)
	}()

	commonArgs := []string{
		"themes",
		"show-config",
	}

	//-------------//
	//  SUCCESSES  //
	//-------------//

	t.Run("help flag works", func(t *testing.T) {
		// GIVEN
		cmd := cli.NewCmd(commonArgs)
		cmd.AddArgs("--help")

		// WHEN
		result, err := fx.RunCmd(cmd)

		// THEN
		require.NoError(t, err)
		snaps.MatchStandaloneSnapshot(t, result)
	})

	t.Run("works for built-in theme", func(t *testing.T) {
		// GIVEN
		cmd := cli.NewCmd(commonArgs)
		cmd.AddArgs("--theme", "monokai-classic")

		// WHEN
		result, err := fx.RunCmd(cmd)

		// THEN
		require.NoError(t, err)
		snaps.MatchStandaloneSnapshot(t, result)
	})

	//------------//
	//  FAILURES  //
	//------------//

	t.Run("fails for incorrect builtin theme", func(t *testing.T) {
		// GIVEN
		cmd := cli.NewCmd(commonArgs)
		cmd.AddArgs("--theme", "unknown")

		// WHEN
		result, err := fx.RunCmd(cmd)

		// THEN
		require.NoError(t, err)
		snaps.MatchStandaloneSnapshot(t, result)
	})

	t.Run("fails for incorrect builtin theme provided via env var", func(t *testing.T) {
		// GIVEN
		cmd := cli.NewCmd(commonArgs)
		cmd.SetEnv("HOURS_THEME", "unknown")

		// WHEN
		result, err := fx.RunCmd(cmd)

		// THEN
		require.NoError(t, err)
		snaps.MatchStandaloneSnapshot(t, result)
	})
}
