// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package applications

import (
	"github.com/AssetMantle/modules/helpers"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	storeTypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	crisisKeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	paramsTypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/ibc-go/v7/testing/simapp"
)

type SimulationApplication interface {
	simapp.App

	GetBaseApp() *baseapp.BaseApp
	GetAppCodec() helpers.Codec
	InterfaceRegistry() codecTypes.InterfaceRegistry
	GetKey(storeKey string) *storeTypes.KVStoreKey
	GetTKey(storeKey string) *storeTypes.TransientStoreKey
	GetMemKey(storeKey string) *storeTypes.MemoryStoreKey
	GetSubspace(moduleName string) paramsTypes.Subspace
	SimulationManager() *module.SimulationManager
	RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig)
	RegisterTxService(clientCtx client.Context)
	RegisterTendermintService(clientCtx client.Context)

	GetCrisisKeeper() crisisKeeper.Keeper
}
