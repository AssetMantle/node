// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/AssetMantle/modules/schema/helpers/constants"
	keysAdd "github.com/AssetMantle/modules/utilities/rest/keys/add"
	"github.com/AssetMantle/modules/utilities/rest/queuing"
	"github.com/AssetMantle/modules/utilities/rest/sign"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authCLI "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authREST "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankCLI "github.com/cosmos/cosmos-sdk/x/bank/client/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"

	"github.com/AssetMantle/node/application"
)

func main() {
	cobra.EnableCommandSorting = false

	config := sdkTypes.GetConfig()
	config.SetBech32PrefixForAccount(sdkTypes.Bech32PrefixAccAddr, sdkTypes.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(sdkTypes.Bech32PrefixValAddr, sdkTypes.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(sdkTypes.Bech32PrefixConsAddr, sdkTypes.Bech32PrefixConsPub)
	config.Seal()

	rootCommand := &cobra.Command{
		Use:   "client",
		Short: "Command line interface for interacting with node",
	}

	rootCommand.PersistentFlags().String(flags.FlagChainID, "", "Chain ID of tendermint node")
	rootCommand.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		return initializeConfiguration(rootCommand)
	}

	rootCommand.AddCommand(
		rpc.StatusCommand(),
		client.ConfigCmd(application.Prototype.GetDefaultClientHome()),
		queryCommand(application.Prototype.GetCodec()),
		transactionCommand(application.Prototype.GetCodec()),
		flags.LineBreak,
		ServeCommand(application.Prototype.GetCodec()),
		flags.LineBreak,
		keys.Commands(),
		flags.LineBreak,
		version.Cmd,
		flags.NewCompletionCmd(rootCommand, true),
	)

	executor := cli.PrepareMainCmd(rootCommand, "HC", application.Prototype.GetDefaultClientHome())

	err := executor.Execute()
	if err != nil {
		fmt.Printf("Failed executing CLI command: %s, exiting...\n", err)
		os.Exit(1)
	}
}

func registerRoutes(restServer *lcd.RestServer) {
	client.RegisterRoutes(restServer.CliCtx, restServer.Mux)
	authREST.RegisterTxRoutes(restServer.CliCtx, restServer.Mux)
	application.Prototype.GetModuleBasicManager().RegisterRESTRoutes(restServer.CliCtx, restServer.Mux)
	keysAdd.RegisterRESTRoutes(restServer.CliCtx, restServer.Mux)
	sign.RegisterRESTRoutes(restServer.CliCtx, restServer.Mux)
	queuing.RegisterRoutes(restServer.CliCtx, restServer.Mux)
}

func queryCommand(codec *codec.Codec) *cobra.Command {
	queryCommand := &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "Root command for querying.",
	}

	queryCommand.AddCommand(
		authCLI.GetAccountCmd(codec),
		flags.LineBreak,
		rpc.ValidatorCommand(codec),
		rpc.BlockCommand(),
		authCLI.QueryTxsByEventsCmd(codec),
		authCLI.QueryTxCmd(codec),
		flags.LineBreak,
	)

	application.Prototype.GetModuleBasicManager().AddQueryCommands(queryCommand, codec)

	return queryCommand
}

// ServeCommand will start the application REST service as a blocking process. It
// takes a codec to create a RestServer object and a function to register all
// necessary routes.
func ServeCommand(codec *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rest-server",
		Short: "Start LCD (light-client daemon), a local REST server",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			restServer := lcd.NewRestServer(codec)
			registerRoutes(restServer)

			if constants.Queuing.ReadCLIValue().(bool) {
				queuing.InitializeKafka(strings.Split(strings.Trim(constants.KafkaNodes.ReadCLIValue().(string), "\" "), " "), restServer.CliCtx)
			}
			return restServer.Start(
				viper.GetString(flags.FlagListenAddr),
				viper.GetInt(flags.FlagMaxOpenConnections),
				uint(viper.GetInt(flags.FlagRPCReadTimeout)),
				uint(viper.GetInt(flags.FlagRPCWriteTimeout)),
				viper.GetBool(flags.FlagUnsafeCORS),
			)
		},
	}
	constants.Queuing.Register(cmd)
	constants.KafkaNodes.Register(cmd)
	cmd.Flags().String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend, "Select keyring's backend (os|file|kwallet|pass|test)")
	cmd.Flags().Bool(flags.FlagGenerateOnly, false, "Build an unsigned transaction and write it as response to rest (when enabled, the local Keybase is not accessible and the node operates offline)")
	cmd.Flags().StringP(flags.FlagBroadcastMode, "b", flags.BroadcastSync, "Transaction broadcasting mode (sync|async|block)")

	return flags.RegisterRestServerFlags(cmd)
}

func transactionCommand(codec *codec.Codec) *cobra.Command {
	transactionCommand := &cobra.Command{
		Use:   "tx",
		Short: "Transactions subcommands",
	}

	transactionCommand.AddCommand(
		bankCLI.SendTxCmd(codec),
		flags.LineBreak,
		authCLI.GetSignCommand(codec),
		authCLI.GetMultiSignCommand(codec),
		flags.LineBreak,
		authCLI.GetBroadcastCommand(codec),
		authCLI.GetEncodeCommand(codec),
		flags.LineBreak,
	)

	application.Prototype.GetModuleBasicManager().AddTxCommands(transactionCommand, codec)

	var commandListToRemove []*cobra.Command

	for _, cmd := range transactionCommand.Commands() {
		if cmd.Use == auth.ModuleName || cmd.Use == bank.ModuleName {
			commandListToRemove = append(commandListToRemove, cmd)
		}
	}

	transactionCommand.RemoveCommand(commandListToRemove...)

	return transactionCommand
}

func initializeConfiguration(command *cobra.Command) error {
	home, err := command.PersistentFlags().GetString(cli.HomeFlag)
	if err != nil {
		return err
	}

	configurationFile := path.Join(home, "configuration", "configuration.toml")
	if _, err := os.Stat(configurationFile); err == nil {
		viper.SetConfigFile(configurationFile)

		if err := viper.ReadInConfig(); err != nil {
			return err
		}
	}
	if err := viper.BindPFlag(flags.FlagChainID, command.PersistentFlags().Lookup(flags.FlagChainID)); err != nil {
		return err
	}
	if err := viper.BindPFlag(cli.EncodingFlag, command.PersistentFlags().Lookup(cli.EncodingFlag)); err != nil {
		return err
	}
	return viper.BindPFlag(cli.OutputFlag, command.PersistentFlags().Lookup(cli.OutputFlag))
}
