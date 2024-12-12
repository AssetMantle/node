// Copyright [2021] - [2025], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package sign

import (
	"testing"

	"github.com/AssetMantle/modules/helpers/constants"
)

func Test_SignTx_Response(t *testing.T) {
	testFee := legacytx.NewStdFee(12, sdkTypes.NewCoins())

	testStdTx := legacytx.NewStdTx([]sdkTypes.Msg{}, testFee, []legacytx.StdSignature{}, "")
	require.Equal(t, response{Success: true, Error: nil, Tx: testStdTx}, newResponse(testStdTx, nil))
	testResponse := newResponse(testStdTx, constants.IncorrectFormat)
	require.Equal(t, false, testResponse.IsSuccessful())
	require.Equal(t, constants.IncorrectFormat, testResponse.GetError())
}
