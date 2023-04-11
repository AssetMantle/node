// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package queuing

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/stretchr/testify/require"

	"github.com/AssetMantle/schema/x"
)

func Test_Rest_Utils(t *testing.T) {
	value, err := parseGasAdjustment("")
	require.Equal(t, flags.DefaultGasAdjustment, value)
	require.Equal(t, nil, err)

	value2, error2 := parseGasAdjustment("test")
	require.Equal(t, float64(0), value2)
	require.NotNil(t, error2)

	value3, error3 := parseGasAdjustment("0.3")
	require.Equal(t, 0.3, value3)
	require.Equal(t, nil, error3)

	var legacyAmino = codec.NewLegacyAmino()
	x.RegisterLegacyAminoCodec(legacyAmino)
	std.RegisterLegacyAminoCodec(legacyAmino)
	legacyAmino.Seal()

	gas := uint64(123)
	response, err := simulationResponse(legacyAmino, gas)
	gasEst := rest.GasEstimateResponse{GasEstimate: gas}
	resp, _ := legacyAmino.MarshalJSON(gasEst)
	require.Equal(t, resp, response)
	require.Equal(t, nil, err)

}
