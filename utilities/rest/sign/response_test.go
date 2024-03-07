// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package sign

import (
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
	"testing"

	errorConstants "github.com/AssetMantle/schema/go/errors/constants"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func Test_SignTx_Response(t *testing.T) {
	testFee := legacytx.NewStdFee(12, sdkTypes.NewCoins())

	testStdTx := legacytx.NewStdTx([]sdkTypes.Msg{}, testFee, []legacytx.StdSignature{}, "")
	require.Equal(t, response{Success: true, Error: nil, Tx: testStdTx}, newResponse(testStdTx, nil))
	testResponse := newResponse(testStdTx, errorConstants.IncorrectFormat)
	require.Equal(t, false, testResponse.IsSuccessful())
	require.Equal(t, errorConstants.IncorrectFormat, testResponse.GetError())
}
