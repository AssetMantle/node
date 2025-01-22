// Copyright [2021] - [2025], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package commands

import (
	"github.com/cosmos/cosmos-sdk/client/snapshot"
	"os"

	"github.com/AssetMantle/modules/helpers"
	"github.com/AssetMantle/modules/helpers/base"
	tmcfg "github.com/cometbft/cometbft/config"
	tendermintCLI "github.com/cometbft/cometbft/libs/cli"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/client/debug"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/pruning"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/server"
	serverConfig "github.com/cosmos/cosmos-sdk/server/config"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	sdkClientCLI "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/spf13/cobra"

	"github.com/AssetMantle/node/application"
	"github.com/AssetMantle/node/application/internal/configurations"
)

func RootCommand() (*cobra.Command, helpers.Codec) {
	codec := base.CodecPrototype().Initialize(configurations.ModuleBasicManager)
	context := client.Context{}.
		WithCodec(codec.GetProtoCodec()).
		WithInterfaceRegistry(codec.InterfaceRegistry()).
		WithTxConfig(codec).
		WithLegacyAmino(codec.GetLegacyAmino()).
		WithInput(os.Stdin).
		WithAccountRetriever(types.AccountRetriever{}).
		WithHomeDir(application.Prototype.GetDefaultNodeHome()).
		WithViper("")

	rootCmd := &cobra.Command{
		Use:   "mantleNode",
		Short: "",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			context, err := client.ReadPersistentCommandFlags(context, cmd.Flags())
			if err != nil {
				return err
			}

			context, err = config.ReadFromClientConfig(context)
			if err != nil {
				return err
			}

			if err = client.SetCmdClientContextHandler(context, cmd); err != nil {
				return err
			}

			ServerConfig := serverConfig.DefaultConfig()
			ServerConfig.StateSync.SnapshotInterval = 1000
			ServerConfig.StateSync.SnapshotKeepRecent = 10
			ServerConfig.API.Enable = true
			ServerConfig.API.Swagger = true

			return server.InterceptConfigsPreRunHandler(cmd, serverConfig.DefaultConfigTemplate, &ServerConfig, tmcfg.DefaultConfig())
		},
	}

	configurations.SetAndSealSDKConfig()

	rootCmd.AddCommand(
		sdkClientCLI.InitCmd(configurations.ModuleBasicManager.GetBasicManager(), application.Prototype.GetDefaultNodeHome()),
		tendermintCLI.NewCompletionCmd(rootCmd, true),
		debug.Cmd(),
		config.Cmd(),
		pruning.PruningCmd(application.Prototype.AppCreator),
		TestnetCommand(configurations.ModuleBasicManager.GetBasicManager(), bankTypes.GenesisBalancesIterator{}),
		snapshot.Cmd(application.Prototype.AppCreator),
		version.NewVersionCommand(),
	)

	server.AddCommands(rootCmd, application.Prototype.GetDefaultNodeHome(), application.Prototype.AppCreator, application.Prototype.AppExporter, application.Prototype.ModuleInitFlags)

	rootCmd.AddCommand(
		rpc.StatusCommand(),
		sdkClientCLI.GenesisCoreCommand(context.TxConfig, configurations.ModuleBasicManager.GetBasicManager(), application.Prototype.GetDefaultNodeHome()),
		queryCommand(),
		txCommand(),
		keys.Commands(application.Prototype.GetDefaultNodeHome()),
	)

	return rootCmd, codec
}

func queryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "query",
		Aliases:                    []string{"q"},
		Short:                      "Querying subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		cli.GetAccountCmd(),
		rpc.ValidatorCommand(),
		rpc.BlockCommand(),
		cli.QueryTxsByEventsCmd(),
		cli.QueryTxCmd(),
	)

	configurations.ModuleBasicManager.AddQueryCommands(cmd)
	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

func txCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "tx",
		Short:                      "Transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		cli.GetSignCommand(),
		cli.GetSignBatchCommand(),
		cli.GetMultiSignCommand(),
		cli.GetMultiSignBatchCmd(),
		cli.GetValidateSignaturesCommand(),
		flags.LineBreak,
		cli.GetBroadcastCommand(),
		cli.GetEncodeCommand(),
		cli.GetDecodeCommand(),
	)

	configurations.ModuleBasicManager.AddTxCommands(cmd)
	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}
