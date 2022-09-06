// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"encoding/json"
	"io"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/debug"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/store"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tendermintABCITypes "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	tendermintTypes "github.com/tendermint/tendermint/types"
	tendermintDB "github.com/tendermint/tm-db"

	"github.com/AssetMantle/node/application"
	"github.com/AssetMantle/node/application/initialize"
)

const flagInvariantsCheckPeriod = "invariants-check-period"

var invariantsCheckPeriod uint

func main() {
	serverContext := server.NewDefaultContext()

	configuration := sdkTypes.GetConfig()
	configuration.SetBech32PrefixForAccount(sdkTypes.Bech32PrefixAccAddr, sdkTypes.Bech32PrefixAccPub)
	configuration.SetBech32PrefixForValidator(sdkTypes.Bech32PrefixValAddr, sdkTypes.Bech32PrefixValPub)
	configuration.SetBech32PrefixForConsensusNode(sdkTypes.Bech32PrefixConsAddr, sdkTypes.Bech32PrefixConsPub)
	configuration.Seal()

	cobra.EnableCommandSorting = false

	rootCommand := &cobra.Command{
		Use:               "assetNode",
		Short:             "Persistence Hub Node Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(serverContext),
	}

	rootCommand.AddCommand(initialize.Command(
		serverContext,
		application.Prototype.GetCodec(),
		application.Prototype.GetModuleBasicManager(),
		application.Prototype.GetDefaultNodeHome(),
	))
	rootCommand.AddCommand(initialize.CollectGenesisTransactionsCommand(
		serverContext,
		application.Prototype.GetCodec(),
		auth.GenesisAccountIterator{},
		application.Prototype.GetDefaultNodeHome(),
	))
	rootCommand.AddCommand(initialize.MigrateGenesisCommand(
		serverContext,
		application.Prototype.GetCodec(),
	))
	rootCommand.AddCommand(initialize.GenesisTransactionCommand(
		serverContext,
		application.Prototype.GetCodec(),
		application.Prototype.GetModuleBasicManager(),
		staking.AppModuleBasic{},
		auth.GenesisAccountIterator{},
		application.Prototype.GetDefaultNodeHome(),
		application.Prototype.GetDefaultClientHome(),
	))
	rootCommand.AddCommand(initialize.ValidateGenesisCommand(
		serverContext,
		application.Prototype.GetCodec(),
		application.Prototype.GetModuleBasicManager(),
	))
	rootCommand.AddCommand(initialize.AddGenesisAccountCommand(
		serverContext,
		application.Prototype.GetCodec(),
		application.Prototype.GetDefaultNodeHome(),
		application.Prototype.GetDefaultClientHome(),
	))
	rootCommand.AddCommand(flags.NewCompletionCmd(rootCommand, true))
	rootCommand.AddCommand(initialize.ReplayTransactionsCommand())
	rootCommand.AddCommand(debug.Cmd(application.Prototype.GetCodec()))
	rootCommand.AddCommand(version.Cmd)
	rootCommand.PersistentFlags().UintVar(
		&invariantsCheckPeriod,
		flagInvariantsCheckPeriod,
		0,
		"Assert registered invariants every N blocks",
	)

	appCreator := func(
		logger log.Logger,
		db tendermintDB.DB,
		traceStore io.Writer,
	) tendermintABCITypes.Application {
		var cache sdkTypes.MultiStorePersistentCache

		if viper.GetBool(server.FlagInterBlockCache) {
			cache = store.NewCommitKVStoreCacheManager()
		}

		skipUpgradeHeights := make(map[int64]bool)
		for _, h := range viper.GetIntSlice(server.FlagUnsafeSkipUpgrades) {
			skipUpgradeHeights[int64(h)] = true
		}

		pruningOpts, err := server.GetPruningOptionsFromFlags()
		if err != nil {
			panic(err)
		}

		return application.Prototype.Initialize(
			logger,
			db,
			traceStore,
			true,
			invariantsCheckPeriod,
			skipUpgradeHeights,
			viper.GetString(flags.FlagHome),
			baseapp.SetPruning(pruningOpts),
			baseapp.SetMinGasPrices(viper.GetString(server.FlagMinGasPrices)),
			baseapp.SetHaltHeight(viper.GetUint64(server.FlagHaltHeight)),
			baseapp.SetHaltTime(viper.GetUint64(server.FlagHaltTime)),
			baseapp.SetInterBlockCache(cache),
		)
	}

	appExporter := func(
		logger log.Logger,
		db tendermintDB.DB,
		traceStore io.Writer,
		height int64,
		forZeroHeight bool,
		jailWhiteList []string,
	) (json.RawMessage, []tendermintTypes.GenesisValidator, error) {
		if height != -1 {
			genesisApplication := application.Prototype.Initialize(
				logger,
				db,
				traceStore,
				false,
				uint(1),
				map[int64]bool{},
				"",
			)
			err := genesisApplication.LoadHeight(height)

			if err != nil {
				return nil, nil, err
			}

			return genesisApplication.ExportApplicationStateAndValidators(forZeroHeight, jailWhiteList)
		}

		genesisApplication := application.Prototype.Initialize(
			logger,
			db,
			traceStore,
			true,
			uint(1),
			map[int64]bool{},
			"",
		)

		return genesisApplication.ExportApplicationStateAndValidators(forZeroHeight, jailWhiteList)
	}

	server.AddCommands(
		serverContext,
		application.Prototype.GetCodec(),
		rootCommand,
		appCreator,
		appExporter,
	)

	executor := cli.PrepareBaseCmd(rootCommand, "CA", application.Prototype.GetDefaultNodeHome())
	err := executor.Execute()

	if err != nil {
		panic(err)
	}
}
