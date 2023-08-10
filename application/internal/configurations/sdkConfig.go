// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package configurations

import (
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

func SetAndSealSDKConfig() {
	sdkConfig := sdkTypes.GetConfig()
	sdkConfig.SetBech32PrefixForAccount(bech32Prefix, bech32Prefix+sdkTypes.PrefixPublic)
	sdkConfig.SetBech32PrefixForValidator(bech32Prefix+sdkTypes.PrefixValidator+sdkTypes.PrefixOperator, bech32Prefix+sdkTypes.PrefixValidator+sdkTypes.PrefixOperator+sdkTypes.PrefixPublic)
	sdkConfig.SetBech32PrefixForConsensusNode(bech32Prefix+sdkTypes.PrefixValidator+sdkTypes.PrefixConsensus, bech32Prefix+sdkTypes.PrefixValidator+sdkTypes.PrefixConsensus+sdkTypes.PrefixPublic)
	sdkConfig.SetCoinType(coinType)
	sdkConfig.SetPurpose(purpose)
	sdkConfig.Seal()
}

func coinDenomRegex() string {
	return `[a-zA-Z][a-zA-Z0-9+_=/-]{2,127}`
}

func SetCoinDenomRegex() {
	sdkTypes.SetCoinDenomRegex(coinDenomRegex)
}
