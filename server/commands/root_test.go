package commands_test

import (
	"testing"

	serverCmd "github.com/cosmos/cosmos-sdk/server/cmd"

	"github.com/stretchr/testify/require"

	"github.com/AssetMantle/node/application"
	"github.com/AssetMantle/node/server/commands"
)

func TestRootCmdConfig(t *testing.T) {

	rootCmd, _ := commands.NewRootCommand()
	rootCmd.SetArgs([]string{
		"config",          // Test the config command
		"keyring-backend", // key
		"test",            // value
	})

	require.NoError(t, serverCmd.Execute(rootCmd, application.DefaultNodeHome))
}
