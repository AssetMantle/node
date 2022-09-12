// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package application

import (
	sdkTypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/AssetMantle/node/application/internal/configurations"
)

func SetConfiguration() {
	configuration := sdkTypes.GetConfig()
	configuration.SetBech32PrefixForAccount(configurations.Bech32PrefixAccAddr, configurations.Bech32PrefixAccPub)
	configuration.SetBech32PrefixForValidator(configurations.Bech32PrefixValAddr, configurations.Bech32PrefixValPub)
	configuration.SetBech32PrefixForConsensusNode(configurations.Bech32PrefixConsAddr, configurations.Bech32PrefixConsPub)
	configuration.SetCoinType(configurations.CoinType)
	configuration.SetFullFundraiserPath(configurations.FullFundraiserPath)
	configuration.Seal()
}
