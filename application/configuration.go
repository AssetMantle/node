// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package application

import (
	"strconv"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/AssetMantle/node/application/internal/configurations"
)

func SetConfiguration() {
	configuration := sdkTypes.GetConfig()
	configuration.SetBech32PrefixForAccount(configurations.Bech32PrefixAccAddr, configurations.Bech32PrefixAccAddr+sdkTypes.PrefixPublic)
	configuration.SetBech32PrefixForValidator(configurations.Bech32PrefixAccAddr+sdkTypes.PrefixValidator+sdkTypes.PrefixOperator, configurations.Bech32PrefixAccAddr+sdkTypes.PrefixValidator+sdkTypes.PrefixOperator+sdkTypes.PrefixPublic)
	configuration.SetBech32PrefixForConsensusNode(configurations.Bech32PrefixAccAddr+sdkTypes.PrefixValidator+sdkTypes.PrefixConsensus, configurations.Bech32PrefixAccAddr+sdkTypes.PrefixValidator+sdkTypes.PrefixConsensus+sdkTypes.PrefixPublic)
	configuration.SetCoinType(configurations.CoinType)
	configuration.SetFullFundraiserPath("44'/" + strconv.Itoa(configurations.CoinType) + "'/0'/0/0")
	configuration.Seal()
}
