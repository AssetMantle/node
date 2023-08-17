// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package constants

import (
	"github.com/AssetMantle/modules/x/assets"
	"github.com/AssetMantle/modules/x/classifications"
	"github.com/AssetMantle/modules/x/identities"
	"github.com/AssetMantle/modules/x/maintainers"
	"github.com/AssetMantle/modules/x/orders"
	"github.com/AssetMantle/modules/x/splits"
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmClient "github.com/CosmWasm/wasmd/x/wasm/client"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	authzModule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/capability"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	distributionClient "github.com/cosmos/cosmos-sdk/x/distribution/client"
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	feeGrantModule "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/mint"
	mintTypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsClient "github.com/cosmos/cosmos-sdk/x/params/client"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeClient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	ica "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts"
	icaTypes "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts/types"
	"github.com/cosmos/ibc-go/v4/modules/apps/transfer"
	ibcTransferTypes "github.com/cosmos/ibc-go/v4/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v4/modules/core"
	ibcClientClient "github.com/cosmos/ibc-go/v4/modules/core/02-client/client"
	"github.com/strangelove-ventures/packet-forward-middleware/v4/router"

	"github.com/AssetMantle/modules/x/metas"
)

const applicationName = "SimulationApplication"

var ModuleBasicManagers = module.NewBasicManager(
	auth.AppModuleBasic{},
	genutil.AppModuleBasic{},
	bank.AppModuleBasic{},
	capability.AppModuleBasic{},
	staking.AppModuleBasic{},
	mint.AppModuleBasic{},
	distribution.AppModuleBasic{},
	gov.NewAppModuleBasic(
		append(wasmClient.ProposalHandlers,
			paramsClient.ProposalHandler,
			distributionClient.ProposalHandler,
			upgradeClient.ProposalHandler,
			upgradeClient.CancelProposalHandler,
			ibcClientClient.UpdateClientProposalHandler,
			ibcClientClient.UpgradeProposalHandler,
		)...,
	),
	params.AppModuleBasic{},
	crisis.AppModuleBasic{},
	slashing.AppModuleBasic{},
	feeGrantModule.AppModuleBasic{},
	authzModule.AppModuleBasic{},
	ibc.AppModuleBasic{},
	upgrade.AppModuleBasic{},
	evidence.AppModuleBasic{},
	transfer.AppModuleBasic{},
	vesting.AppModuleBasic{},
	router.AppModuleBasic{},
	ica.AppModuleBasic{},

	wasm.AppModuleBasic{},

	assets.Prototype(),
	classifications.Prototype(),
	identities.Prototype(),
	maintainers.Prototype(),
	metas.Prototype(),
	orders.Prototype(),
	splits.Prototype(),
)
var ModuleAccountPermissions = map[string][]string{
	authTypes.FeeCollectorName:         nil,
	distributionTypes.ModuleName:       nil,
	icaTypes.ModuleName:                nil,
	mintTypes.ModuleName:               {authTypes.Minter},
	stakingTypes.BondedPoolName:        {authTypes.Burner, authTypes.Staking},
	stakingTypes.NotBondedPoolName:     {authTypes.Burner, authTypes.Staking},
	govTypes.ModuleName:                {authTypes.Burner},
	ibcTransferTypes.ModuleName:        {authTypes.Minter, authTypes.Burner},
	splits.Prototype().Name():          nil,
	classifications.Prototype().Name(): {authTypes.Burner},
}
var DefaultNodeHome string
var tokenReceiveAllowedModules = map[string]bool{
	distributionTypes.ModuleName: true,
}
