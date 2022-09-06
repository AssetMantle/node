// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package initialize

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/types/module"
)

func GenesisTransactionCommand(
	serverContext *server.Context,
	codec *codec.Codec,
	moduleBasicManager module.BasicManager,
	stakingMessageBuildingHelpers cli.StakingMsgBuildingHelpers,
	genesisBalancesIterator types.GenesisAccountsIterator,
	defaultNodeHome string,
	defaultClientHome string,
) *cobra.Command {
	return cli.GenTxCmd(
		serverContext,
		codec,
		moduleBasicManager,
		stakingMessageBuildingHelpers,
		genesisBalancesIterator,
		defaultNodeHome,
		defaultClientHome,
	)
}
