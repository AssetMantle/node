// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package application

import (
	"github.com/AssetMantle/modules/modules/assets"
	"github.com/AssetMantle/modules/modules/classifications"
	"github.com/AssetMantle/modules/modules/identities"
	"github.com/AssetMantle/modules/modules/maintainers"
	"github.com/AssetMantle/modules/modules/metas"
	"github.com/AssetMantle/modules/modules/orders"
	"github.com/AssetMantle/modules/modules/splits"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsClient "github.com/cosmos/cosmos-sdk/x/params/client"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeClient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
)

var moduleAccountPermissions = map[string][]string{
	auth.FeeCollectorName:     nil,
	distribution.ModuleName:   nil,
	mint.ModuleName:           {supply.Minter},
	staking.BondedPoolName:    {supply.Burner, supply.Staking},
	staking.NotBondedPoolName: {supply.Burner, supply.Staking},
	gov.ModuleName:            {supply.Burner},
	splits.Prototype().Name(): nil,
}
var tokenReceiveAllowedModules = map[string]bool{
	distribution.ModuleName: true,
}
var ModuleBasics = module.NewBasicManager(
	genutil.AppModuleBasic{},
	auth.AppModuleBasic{},
	bank.AppModuleBasic{},
	staking.AppModuleBasic{},
	mint.AppModuleBasic{},
	distribution.AppModuleBasic{},
	gov.NewAppModuleBasic(paramsClient.ProposalHandler, distribution.ProposalHandler, upgradeClient.ProposalHandler),
	params.AppModuleBasic{},
	crisis.AppModuleBasic{},
	slashing.AppModuleBasic{},
	supply.AppModuleBasic{},
	upgrade.AppModuleBasic{},
	evidence.AppModuleBasic{},

	assets.Prototype(),
	classifications.Prototype(),
	identities.Prototype(),
	maintainers.Prototype(),
	metas.Prototype(),
	orders.Prototype(),
	splits.Prototype(),
)
