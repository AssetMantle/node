/*
 Copyright [2019] - [2020], PERSISTENCE TECHNOLOGIES PTE. LTD. and the assetMantle contributors
 SPDX-License-Identifier: Apache-2.0
*/

package application

import (
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"os"
)

const (
	Name             = "AssetMantle"
	Bech32MainPrefix = "mantle"

	CoinType           = 750
	FullFundraiserPath = "44'/750'/0'/0/0"

	Bech32PrefixAccAddr  = Bech32MainPrefix
	Bech32PrefixAccPub   = Bech32MainPrefix + sdkTypes.PrefixPublic
	Bech32PrefixValAddr  = Bech32MainPrefix + sdkTypes.PrefixValidator + sdkTypes.PrefixOperator
	Bech32PrefixValPub   = Bech32MainPrefix + sdkTypes.PrefixValidator + sdkTypes.PrefixOperator + sdkTypes.PrefixPublic
	Bech32PrefixConsAddr = Bech32MainPrefix + sdkTypes.PrefixValidator + sdkTypes.PrefixConsensus
	Bech32PrefixConsPub  = Bech32MainPrefix + sdkTypes.PrefixValidator + sdkTypes.PrefixConsensus + sdkTypes.PrefixPublic
)

var DefaultClientHome = os.ExpandEnv("$HOME/.assetClient")
var DefaultNodeHome = os.ExpandEnv("$HOME/.assetNode")
