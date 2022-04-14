package main

import (
	"os"

	"github.com/cosmos/cosmos-sdk/server"
	serverCmd "github.com/cosmos/cosmos-sdk/server/cmd"

	"github.com/AssetMantle/node/application"
	command "github.com/AssetMantle/node/cmd/gaiad/cmd"
)

func main() {
	rootCmd, _ := command.NewRootCommand()

	if err := serverCmd.Execute(rootCmd, application.DefaultNodeHome); err != nil {
		switch e := err.(type) {
		case server.ErrorCode:
			os.Exit(e.Code)

		default:
			os.Exit(1)
		}
	}
}
