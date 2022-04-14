package cmd_test

import (
	"testing"

	serverCmd "github.com/cosmos/cosmos-sdk/server/cmd"

	"github.com/stretchr/testify/require"

	"github.com/AssetMantle/node/application"
	"github.com/AssetMantle/node/cmd/gaiad/cmd"
)

func TestRootCmdConfig(t *testing.T) {

	rootCmd, _ := cmd.NewRootCommand()
	rootCmd.SetArgs([]string{
		"config",          // Test the config cmd
		"keyring-backend", // key
		"test",            // value
	})

	require.NoError(t, serverCmd.Execute(rootCmd, application.DefaultNodeHome))
}
