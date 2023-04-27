// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package applications

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	crisisKeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	paramsTypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

type SimulationApplication interface {
	simapp.App

	GetBaseApp() *baseapp.BaseApp
	GetAppCodec() codec.Codec
	InterfaceRegistry() codecTypes.InterfaceRegistry
	GetKey(storeKey string) *sdkTypes.KVStoreKey
	GetTKey(storeKey string) *sdkTypes.TransientStoreKey
	GetMemKey(storeKey string) *sdkTypes.MemoryStoreKey
	GetSubspace(moduleName string) paramsTypes.Subspace
	SimulationManager() *module.SimulationManager
	RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig)
	RegisterTxService(clientCtx client.Context)
	RegisterTendermintService(clientCtx client.Context)

	GetCrisisKeeper() crisisKeeper.Keeper
}
