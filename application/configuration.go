/*
 Copyright [2019] - [2020], PERSISTENCE TECHNOLOGIES PTE. LTD. and the assetMantle contributors
 SPDX-License-Identifier: Apache-2.0
*/

package application

import sdkTypes "github.com/cosmos/cosmos-sdk/types"

func SetConfiguration() {
	configuration := sdkTypes.GetConfig()
	configuration.SetBech32PrefixForAccount(Bech32PrefixAccAddr, Bech32PrefixAccPub)
	configuration.SetBech32PrefixForValidator(Bech32PrefixValAddr, Bech32PrefixValPub)
	configuration.SetBech32PrefixForConsensusNode(Bech32PrefixConsAddr, Bech32PrefixConsPub)
	configuration.SetCoinType(CoinType)
	configuration.SetFullFundraiserPath(FullFundraiserPath)
	configuration.Seal()
}
