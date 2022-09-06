// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package initialize

import (
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/genutil/client/cli"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"
)

func Command(
	serverContext *server.Context,
	codec *codec.Codec,
	moduleBasicManager module.BasicManager,
	defaultNodeHome string,
) *cobra.Command {
	return cli.InitCmd(
		serverContext,
		codec,
		moduleBasicManager,
		defaultNodeHome,
	)
}
