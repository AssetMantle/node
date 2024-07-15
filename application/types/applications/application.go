// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package applications

import (
	"github.com/AssetMantle/modules/helpers"
	db "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/baseapp"
	serverTypes "github.com/cosmos/cosmos-sdk/server/types"
	storeTypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/spf13/cobra"
	"io"
)

type Application interface {
	serverTypes.Application

	GetDefaultNodeHome() string
	GetDefaultClientHome() string
	GetModuleManager() helpers.ModuleManager
	GetCodec() helpers.Codec

	LoadHeight(int64) error
	ExportApplicationStateAndValidators(bool, []string, []string) (serverTypes.ExportedApp, error)

	Name() string
	Logger() log.Logger
	MountStores(...storeTypes.StoreKey)
	MountKVStores(map[string]*storeTypes.KVStoreKey)
	MountTransientStores(map[string]*storeTypes.TransientStoreKey)
	MountStore(storeTypes.StoreKey, storeTypes.StoreType)
	LastCommitID() storeTypes.CommitID
	LastBlockHeight() int64
	//Router() storeTypes.Router
	//QueryRouter() sdkTypes.QueryRouter
	Seal()
	IsSealed() bool

	AppCreator(log.Logger, db.DB, io.Writer, serverTypes.AppOptions) serverTypes.Application
	AppExporter(log.Logger, db.DB, io.Writer, int64, bool, []string, serverTypes.AppOptions, []string) (serverTypes.ExportedApp, error)
	ModuleInitFlags(*cobra.Command)

	Initialize(log.Logger, db.DB, io.Writer, bool, uint, map[int64]bool, string, serverTypes.AppOptions, ...func(*baseapp.BaseApp)) Application
}
