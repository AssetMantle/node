// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"

	"github.com/cosmos/cosmos-sdk/server"
	serverCmd "github.com/cosmos/cosmos-sdk/server/cmd"

	"github.com/AssetMantle/node/application"
	"github.com/AssetMantle/node/application/commands"
)

func main() {
	rootCmd, _ := commands.RootCommand()

	if err := serverCmd.Execute(rootCmd, "", application.Prototype.GetDefaultNodeHome()); err != nil {
		switch e := err.(type) {
		case server.ErrorCode:
			os.Exit(e.Code)

		default:
			os.Exit(1)
		}
	}
}
