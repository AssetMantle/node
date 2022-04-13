package cmd_test

import (
	"testing"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"

	"github.com/stretchr/testify/require"

	application "github.com/AssetMantle/node/application"
	"github.com/AssetMantle/node/cmd/gaiad/cmd"
)

func TestRootCmdConfig(t *testing.T) {

	rootCmd, _ := cmd.NewRootCmd()
	rootCmd.SetArgs([]string{
		"config",          // Test the config cmd
		"keyring-backend", // key
		"test",            // value
	})

	require.NoError(t, svrcmd.Execute(rootCmd, application.DefaultNodeHome))
}
