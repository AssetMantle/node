package main

import (
	"os"

	"github.com/cosmos/cosmos-sdk/server"
	serverCmd "github.com/cosmos/cosmos-sdk/server/cmd"

	"github.com/AssetMantle/node/application"
	"github.com/AssetMantle/node/server/commands"
)

func main() {
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
