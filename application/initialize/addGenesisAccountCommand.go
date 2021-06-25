/*
 Copyright [2019] - [2020], PERSISTENCE TECHNOLOGIES PTE. LTD. and the assetMantle contributors
 SPDX-License-Identifier: Apache-2.0
*/

package initialize

import (
	"bufio"
	"errors"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/server"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"
)

const (
	flagClientHome   = "home-client"
	flagVestingStart = "vesting-start-time"
	flagVestingEnd   = "vesting-end-time"
	flagVestingAmt   = "vesting-amount"
)

// AddGenesisAccountCommand returns add-genesis-account cobra Command.
func AddGenesisAccountCommand(context *server.Context, codec *codec.Codec, defaultNodeHome, defaultClientHome string) *cobra.Command {
	command := &cobra.Command{
		Use:   "add-genesis-account [address_or_key_name] [coin][,[coin]]",
		Short: "Add a genesis account to genesis.json",
		Long: `Add a genesis account to genesis.json. The provided account must specify
the account address or key name and a list of initial coins. If a key name is given,
the address will be looked up in the local Keybase. The list of initial tokens must
contain valid denominations. Accounts may optionally be supplied with vesting parameters.
`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			config := context.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			addr, err := sdkTypes.AccAddressFromBech32(args[0])
			inBuf := bufio.NewReader(cmd.InOrStdin())
			if err != nil {
				// attempt to lookup address from Keybase if no address was provided
				kb, err := keys.NewKeyring(
					sdkTypes.KeyringServiceName(),
					viper.GetString(flags.FlagKeyringBackend),
					viper.GetString(flagClientHome),
					inBuf,
				)
				if err != nil {
					return err
				}

				info, err := kb.Get(args[0])
				if err != nil {
					return fmt.Errorf("failed to get address from Keybase: %w", err)
				}

				addr = info.GetAddress()
			}

			coins, err := sdkTypes.ParseCoins(args[1])
			if err != nil {
				return fmt.Errorf("failed to parse coins: %w", err)
			}

			vestingStart := viper.GetInt64(flagVestingStart)
			vestingEnd := viper.GetInt64(flagVestingEnd)
			vestingAmt, err := sdkTypes.ParseCoins(viper.GetString(flagVestingAmt))
			if err != nil {
				return fmt.Errorf("failed to parse vesting amount: %w", err)
			}

			// create concrete account type based on input parameters
			var genAccount exported.GenesisAccount

			baseAccount := auth.NewBaseAccount(addr, coins.Sort(), nil, 0, 0)
			if !vestingAmt.IsZero() {
				baseVestingAccount, err := vesting.NewBaseVestingAccount(baseAccount, vestingAmt.Sort(), vestingEnd)
				if err != nil {
					return fmt.Errorf("failed to create base vesting account: %w", err)
				}

				switch {
				case vestingStart != 0 && vestingEnd != 0:
					genAccount = vesting.NewContinuousVestingAccountRaw(baseVestingAccount, vestingStart)

				case vestingEnd != 0:
					genAccount = vesting.NewDelayedVestingAccountRaw(baseVestingAccount)

				default:
					return errors.New("invalid vesting parameters; must supply start and end time or end time")
				}
			} else {
				genAccount = baseAccount
			}

			if err := genAccount.Validate(); err != nil {
				return fmt.Errorf("failed to validate new genesis account: %w", err)
			}

			genFile := config.GenesisFile()
			appState, genDoc, err := genutil.GenesisStateFromGenFile(codec, genFile)
			if err != nil {
				return fmt.Errorf("failed to unmarshal genesis state: %w", err)
			}

			authGenState := auth.GetGenesisStateFromAppState(codec, appState)
			if authGenState.Accounts.Contains(addr) {
				return fmt.Errorf("cannot add account at existing address %s", addr)
			}

			// Add the new account to the set of genesis accounts and sanitize the
			// accounts afterwards.
			authGenState.Accounts = append(authGenState.Accounts, genAccount)
			authGenState.Accounts = auth.SanitizeGenesisAccounts(authGenState.Accounts)

			authGenStateBz, err := codec.MarshalJSON(authGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal auth genesis state: %w", err)
			}

			appState[auth.ModuleName] = authGenStateBz

			appStateJSON, err := codec.MarshalJSON(appState)
			if err != nil {
				return fmt.Errorf("failed to marshal application genesis state: %w", err)
			}

			genDoc.AppState = appStateJSON
			return genutil.ExportGenesisFile(genDoc, genFile)
		},
	}

	command.Flags().String(cli.HomeFlag, defaultNodeHome, "node's home directory")
	command.Flags().String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend, "Select keyring backend (os|file|test)")
	command.Flags().String(flagClientHome, defaultClientHome, "client's home directory")
	command.Flags().String(flagVestingAmt, "", "amount of coins for vesting accounts")
	command.Flags().Uint64(flagVestingStart, 0, "schedule start time (unix epoch) for vesting accounts")
	command.Flags().Uint64(flagVestingEnd, 0, "schedule end time (unix epoch) for vesting accounts")

	return command
}
