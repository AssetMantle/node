// Copyright [2021] - [2025], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package configurations

import (
	"github.com/AssetMantle/modules/x/assets"
	"github.com/AssetMantle/modules/x/classifications"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	mintTypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	icaTypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	ibcTransferTypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
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
	assets.Prototype().Name():          nil,
	classifications.Prototype().Name(): {authTypes.Burner},
}
var TokenReceiveAllowedModules = map[string]bool{
	distributionTypes.ModuleName: true,
}
