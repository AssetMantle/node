// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package sign

import (
	"testing"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"
)

func Test_SignTx_Request(t *testing.T) {
	fromAddress := "cosmos1pkkayn066msg6kn33wnl5srhdt3tnu2vzasz9c"
	testBaseReq := rest.BaseReq{From: fromAddress, ChainID: "test", Fees: sdkTypes.NewCoins()}

	testFee := authTypes.NewStdFee(12, sdkTypes.NewCoins())

	testStdTx := authTypes.NewStdTx([]sdkTypes.Msg{}, testFee, []authTypes.StdSignature{}, "")
	require.Equal(t, nil, request{BaseRequest: testBaseReq, Type: "type", StdTx: testStdTx}.Validate())
}
