// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package configurations

import (
	"github.com/AssetMantle/modules/x/assets"
	"github.com/AssetMantle/modules/x/classifications"
	"github.com/AssetMantle/modules/x/identities"
	"github.com/AssetMantle/modules/x/maintainers"
	"github.com/AssetMantle/modules/x/metas"
	"github.com/AssetMantle/modules/x/orders"
	"github.com/AssetMantle/modules/x/splits"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	authzModule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/capability"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	distributionClient "github.com/cosmos/cosmos-sdk/x/distribution/client"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	feegrantModule "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsClient "github.com/cosmos/cosmos-sdk/x/params/client"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeClient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	ica "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts"
	"github.com/cosmos/ibc-go/v4/modules/apps/transfer"
	ibc "github.com/cosmos/ibc-go/v4/modules/core"
	ibcClientClient "github.com/cosmos/ibc-go/v4/modules/core/02-client/client"
)

var ModuleBasicManager = module.NewBasicManager(
	auth.AppModuleBasic{},
	genutil.AppModuleBasic{},
	bank.AppModuleBasic{},
	capability.AppModuleBasic{},
	staking.AppModuleBasic{},
	mint.AppModuleBasic{},
	distribution.AppModuleBasic{},
	gov.NewAppModuleBasic(
		paramsClient.ProposalHandler,
		distributionClient.ProposalHandler,
		upgradeClient.ProposalHandler,
		upgradeClient.CancelProposalHandler,
		ibcClientClient.UpdateClientProposalHandler,
		ibcClientClient.UpgradeProposalHandler,
	),
	params.AppModuleBasic{},
	crisis.AppModuleBasic{},
	slashing.AppModuleBasic{},
	feegrantModule.AppModuleBasic{},
	authzModule.AppModuleBasic{},
	ibc.AppModuleBasic{},
	upgrade.AppModuleBasic{},
	evidence.AppModuleBasic{},
	transfer.AppModuleBasic{},
	vesting.AppModuleBasic{},
	ica.AppModuleBasic{},

	assets.Prototype(),
	classifications.Prototype(),
	identities.Prototype(),
	maintainers.Prototype(),
	metas.Prototype(),
	orders.Prototype(),
	splits.Prototype(),
)
