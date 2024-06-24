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
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	authzModule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/capability"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	feeGrantModule "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govClient "github.com/cosmos/cosmos-sdk/x/gov/client"
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
	router "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v7/packetforward"
	ica "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts"
	icaTypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	"github.com/cosmos/ibc-go/v7/modules/apps/transfer"
	ibcTransferTypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v7/modules/core"
	ibcClientClient "github.com/cosmos/ibc-go/v7/modules/core/02-client/client"

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
	gov.NewAppModuleBasic([]govClient.ProposalHandler{paramsClient.ProposalHandler, upgradeClient.LegacyProposalHandler, upgradeClient.LegacyCancelProposalHandler, ibcClientClient.UpdateClientProposalHandler, ibcClientClient.UpgradeProposalHandler}),
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
