package main

import (
	"os"

	"github.com/cosmos/cosmos-sdk/server"
	serverCmd "github.com/cosmos/cosmos-sdk/server/cmd"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/AssetMantle/node/application"
	"github.com/AssetMantle/node/server/mantleNode/commands"
)

func main() {
	config := sdkTypes.GetConfig()
	config.SetBech32PrefixForAccount("mantle", "mantlepub")
	config.SetBech32PrefixForValidator("mantlevaloper", "mantlevaloper")
	config.SetBech32PrefixForConsensusNode("mantlevalcons", "mantlevalconspub")
	config.Seal()

	rootCmd, _ := commands.NewRootCommand()

	if err := serverCmd.Execute(rootCmd, application.DefaultNodeHome); err != nil {
		switch e := err.(type) {
		case server.ErrorCode:
			os.Exit(e.Code)

		default:
			os.Exit(1)
		}
	}
}
