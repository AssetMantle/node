// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package applications

import (
	"io"

	"github.com/cosmos/cosmos-sdk/baseapp"
	serverTypes "github.com/cosmos/cosmos-sdk/server/types"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/log"
	tendermintLog "github.com/tendermint/tendermint/libs/log"
	tendermintDB "github.com/tendermint/tm-db"

	"github.com/AssetMantle/modules/helpers"
)

type Application interface {
	serverTypes.Application

	GetDefaultNodeHome() string
	GetDefaultClientHome() string
	GetModuleBasicManager() module.BasicManager
	GetCodec() helpers.Codec

	LoadHeight(int64) error
	ExportApplicationStateAndValidators(bool, []string) (serverTypes.ExportedApp, error)

	Name() string
	Logger() log.Logger
	MountStores(keys ...sdkTypes.StoreKey)
	MountKVStores(keys map[string]*sdkTypes.KVStoreKey)
	MountTransientStores(keys map[string]*sdkTypes.TransientStoreKey)
	MountStore(key sdkTypes.StoreKey, typ sdkTypes.StoreType)
	LastCommitID() sdkTypes.CommitID
	LastBlockHeight() int64
	Router() sdkTypes.Router
	QueryRouter() sdkTypes.QueryRouter
	Seal()
	IsSealed() bool

	AppCreator(log.Logger, tendermintDB.DB, io.Writer, serverTypes.AppOptions) serverTypes.Application
	AppExporter(log.Logger, tendermintDB.DB, io.Writer, int64, bool, []string, serverTypes.AppOptions) (serverTypes.ExportedApp, error)
	ModuleInitFlags(startCmd *cobra.Command)

	Initialize(logger tendermintLog.Logger, db tendermintDB.DB, traceStore io.Writer, loadLatest bool, invCheckPeriod uint, skipUpgradeHeights map[int64]bool, home string, appOptions serverTypes.AppOptions, baseAppOptions ...func(*baseapp.BaseApp)) Application
}
