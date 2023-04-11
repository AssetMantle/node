// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package sign

import (
	"testing"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"

	errorConstants "github.com/AssetMantle/schema/x/errors/constants"
)

func Test_SignTx_Response(t *testing.T) {
	testFee := authTypes.NewStdFee(12, sdkTypes.NewCoins())

	testStdTx := authTypes.NewStdTx([]sdkTypes.Msg{}, testFee, []authTypes.StdSignature{}, "")
	require.Equal(t, response{Success: true, Error: nil, StdTx: testStdTx}, newResponse(testStdTx, nil))
	testResponse := newResponse(testStdTx, errorConstants.IncorrectFormat)
	require.Equal(t, false, testResponse.IsSuccessful())
	require.Equal(t, errorConstants.IncorrectFormat, testResponse.GetError())
}
