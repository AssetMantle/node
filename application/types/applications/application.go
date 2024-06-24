// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package applications

import (
	storeTypes "github.com/cosmos/cosmos-sdk/store/types"
	"io"

	tendermintDB "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tendermintLog "github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/baseapp"
	serverTypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/spf13/cobra"

	"github.com/AssetMantle/modules/helpers"
)

type Application interface {
	serverTypes.Application

	GetDefaultNodeHome() string
	GetDefaultClientHome() string
	GetModuleBasicManager() module.BasicManager
	GetCodec() helpers.Codec

	LoadHeight(int64) error
	ExportApplicationStateAndValidators(bool, []string, []string) (serverTypes.ExportedApp, error)

	Name() string
	Logger() log.Logger
	MountStores(keys ...storeTypes.StoreKey)
	MountKVStores(keys map[string]*storeTypes.KVStoreKey)
	MountTransientStores(keys map[string]*storeTypes.TransientStoreKey)
	MountStore(key storeTypes.StoreKey, typ storeTypes.StoreType)
	LastCommitID() storeTypes.CommitID
	LastBlockHeight() int64
	//Router() storeTypes.Router
	//QueryRouter() sdkTypes.QueryRouter
	Seal()
	IsSealed() bool

	AppCreator(log.Logger, tendermintDB.DB, io.Writer, serverTypes.AppOptions) serverTypes.Application
	AppExporter(log.Logger, tendermintDB.DB, io.Writer, int64, bool, []string, serverTypes.AppOptions, []string) (serverTypes.ExportedApp, error)
	ModuleInitFlags(startCmd *cobra.Command)

	Initialize(logger tendermintLog.Logger, db tendermintDB.DB, traceStore io.Writer, loadLatest bool, invCheckPeriod uint, skipUpgradeHeights map[int64]bool, home string, appOptions serverTypes.AppOptions, baseAppOptions ...func(*baseapp.BaseApp)) Application
}
