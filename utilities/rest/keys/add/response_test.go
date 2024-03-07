// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package add

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys/multisig"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"testing"

	errorConstants "github.com/AssetMantle/schema/go/errors/constants"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/stretchr/testify/require"
)

func Test_Add_Response(t *testing.T) {
	sk := secp256k1.PrivKey{Key: []byte{154, 49, 3, 117, 55, 232, 249, 20, 205, 216, 102, 7, 136, 72, 177, 2, 131, 202, 234, 81, 31, 208, 46, 244, 179, 192, 167, 163, 142, 117, 246, 13}}
	tmpKey := sk.PubKey()
	multisigPk := multisig.NewLegacyAminoPubKey(1, []types.PubKey{tmpKey})

	info, err := keyring.NewMultiInfo("multisig", multisigPk)
	require.NoError(t, err)
	accAddr := sdk.AccAddress(info.GetPubKey().Address())
	testKeyOutput, err := keyring.NewKeyOutput(info.GetName(), info.GetType(), accAddr, multisigPk)
	require.NoError(t, err)
	testResponse := newResponse(testKeyOutput, nil)
	require.Equal(t, response{Success: true, Error: nil, KeyOutput: testKeyOutput}, testResponse)
	require.Equal(t, true, testResponse.IsSuccessful())
	require.Equal(t, nil, testResponse.GetError())
	testResponse2 := newResponse(testKeyOutput, errorConstants.IncorrectFormat)
	require.Equal(t, false, testResponse2.IsSuccessful())
	require.Equal(t, errorConstants.IncorrectFormat, testResponse2.GetError())
}
