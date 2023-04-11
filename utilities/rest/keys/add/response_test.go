// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package add

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/stretchr/testify/require"

	errorConstants "github.com/AssetMantle/schema/x/errors/constants"
)

func Test_Add_Response(t *testing.T) {
	testKeyOutput, _ := keyring.NewKeyOutput("name", "keyType", "address", "pubkey")
	testResponse := newResponse(testKeyOutput, nil)
	require.Equal(t, response{Success: true, Error: nil, KeyOutput: testKeyOutput}, testResponse)
	require.Equal(t, true, testResponse.IsSuccessful())
	require.Equal(t, nil, testResponse.GetError())
	testResponse2 := newResponse(testKeyOutput, errorConstants.IncorrectFormat)
	require.Equal(t, false, testResponse2.IsSuccessful())
	require.Equal(t, errorConstants.IncorrectFormat, testResponse2.GetError())
}
