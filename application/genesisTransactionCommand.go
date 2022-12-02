// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package application

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/types/module"
)

func GenesisTransactionCommand(basicManager module.BasicManager, txEncodingConfig client.TxEncodingConfig, genesisBalancesIterator types.GenesisBalancesIterator, defaultNodeHome string) *cobra.Command {
	return cli.GenTxCmd(basicManager, txEncodingConfig, genesisBalancesIterator,
		defaultNodeHome,
	)
}
