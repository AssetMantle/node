// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package application

import (
	"github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/spf13/cobra"
)

func CollectGenesisTransactionsCommand(genesisBalancesIterator types.GenesisBalancesIterator, defaultNodeHome string) *cobra.Command {
	return cli.CollectGenTxsCmd(genesisBalancesIterator, defaultNodeHome)
}
