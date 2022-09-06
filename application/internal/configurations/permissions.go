// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package configurations

import (
	"github.com/AssetMantle/modules/modules/splits"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
)

var ModuleAccountPermissions = map[string][]string{
	auth.FeeCollectorName:     nil,
	distribution.ModuleName:   nil,
	mint.ModuleName:           {supply.Minter},
	staking.BondedPoolName:    {supply.Burner, supply.Staking},
	staking.NotBondedPoolName: {supply.Burner, supply.Staking},
	gov.ModuleName:            {supply.Burner},
	splits.Prototype().Name(): nil,
}
var TokenReceiveAllowedModules = map[string]bool{
	distribution.ModuleName: true,
}
