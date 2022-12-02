// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package application

import (
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/genutil/client/cli"

	"github.com/spf13/cobra"
)

func Command(basicManager module.BasicManager, defaultNodeHome string) *cobra.Command {
	return cli.InitCmd(basicManager, defaultNodeHome)
}
