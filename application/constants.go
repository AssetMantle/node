// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package application

import (
	"os"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

const (
	Name             = "AssetMantle"
	Bech32MainPrefix = "mantle"

	CoinType           = 18
	FullFundraiserPath = "44'/118'/0'/0/0"

	Bech32PrefixAccAddr  = Bech32MainPrefix
	Bech32PrefixAccPub   = Bech32MainPrefix + sdkTypes.PrefixPublic
	Bech32PrefixValAddr  = Bech32MainPrefix + sdkTypes.PrefixValidator + sdkTypes.PrefixOperator
	Bech32PrefixValPub   = Bech32MainPrefix + sdkTypes.PrefixValidator + sdkTypes.PrefixOperator + sdkTypes.PrefixPublic
	Bech32PrefixConsAddr = Bech32MainPrefix + sdkTypes.PrefixValidator + sdkTypes.PrefixConsensus
	Bech32PrefixConsPub  = Bech32MainPrefix + sdkTypes.PrefixValidator + sdkTypes.PrefixConsensus + sdkTypes.PrefixPublic
)

var DefaultClientHome = os.ExpandEnv("$HOME/.assetClient")
var DefaultNodeHome = os.ExpandEnv("$HOME/.assetNode")
