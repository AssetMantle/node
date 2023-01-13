package commands

import (
	"os"

	"github.com/AssetMantle/modules/schema/helpers"
	"github.com/AssetMantle/modules/schema/helpers/base"
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
	tendermintCLI "github.com/tendermint/tendermint/libs/cli"

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

			return server.InterceptConfigsPreRunHandler(cmd, serverConfig.DefaultConfigTemplate, &ServerConfig)
		},
	}

	configurations.SetAndSealSDKConfig()

	rootCmd.AddCommand(
		sdkClientCLI.InitCmd(configurations.ModuleBasicManager, application.Prototype.GetDefaultNodeHome()),
		sdkClientCLI.CollectGenTxsCmd(bankTypes.GenesisBalancesIterator{}, application.Prototype.GetDefaultNodeHome()),
		sdkClientCLI.MigrateGenesisCmd(),
		sdkClientCLI.GenTxCmd(configurations.ModuleBasicManager, codec, bankTypes.GenesisBalancesIterator{}, application.Prototype.GetDefaultNodeHome()),
		sdkClientCLI.ValidateGenesisCmd(configurations.ModuleBasicManager),
		AddGenesisAccountCommand(application.Prototype.GetDefaultNodeHome()),
		tendermintCLI.NewCompletionCmd(rootCmd, true),
		debug.Cmd(),
		config.Cmd(),
		pruning.PruningCmd(application.Prototype.AppCreator),
		TestnetCommand(configurations.ModuleBasicManager, bankTypes.GenesisBalancesIterator{}),
		version.NewVersionCommand(),
	)

	server.AddCommands(rootCmd, application.Prototype.GetDefaultNodeHome(), application.Prototype.AppCreator, application.Prototype.AppExporter, application.Prototype.ModuleInitFlags)

	rootCmd.AddCommand(
		rpc.StatusCommand(),
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
