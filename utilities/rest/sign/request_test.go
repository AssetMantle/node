// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package sign

import (
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
	"testing"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/stretchr/testify/require"
)

func Test_SignTx_Request(t *testing.T) {
	fromAddress := "cosmos1pkkayn066msg6kn33wnl5srhdt3tnu2vzasz9c"
	testBaseReq := rest.BaseReq{From: fromAddress, ChainID: "test", Fees: sdkTypes.NewCoins()}

	testFee := legacytx.NewStdFee(12, sdkTypes.NewCoins())

	testStdTx := legacytx.NewStdTx([]sdkTypes.Msg{}, testFee, []legacytx.StdSignature{}, "")
	require.Equal(t, nil, request{BaseRequest: testBaseReq, Type: "type", StdTx: testStdTx}.Validate())
}
