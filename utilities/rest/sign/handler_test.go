// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package sign

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	schemaCodec "github.com/AssetMantle/schema/go/codec"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/std"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestHandler(t *testing.T) {
	Codec := codec.NewLegacyAmino()
	schemaCodec.RegisterLegacyAminoCodec(Codec)
	std.RegisterLegacyAminoCodec(Codec)
	Codec.RegisterConcrete(request{}, "request", nil)
	Codec.RegisterConcrete(response{}, "response", nil)
	TestMessagePrototype().RegisterLegacyAminoCodec(Codec)

	handler := handler(client.Context{})
	viper.Set(flags.FlagKeyringBackend, keyring.BackendTest)

	Keyring, err := keyring.New(sdk.KeyringServiceName(), keyring.BackendTest, viper.GetString(flags.FlagHome), strings.NewReader(""))
	require.NoError(t, err)

	router := mux.NewRouter()
	RegisterRESTRoutes(client.Context{}, router)

	t.Cleanup(func() {
		_ = Keyring.Delete("keyName1")
		_ = Keyring.Delete("keyName2")
		_ = Keyring.Delete("keyName3")
	})
	_, err = Keyring.NewAccount("keyName1", "wage thunder live sense resemble foil apple course spin horse glass mansion midnight laundry acoustic rhythm loan scale talent push green direct brick please",
		keyring.DefaultBIP39Passphrase, sdkTypes.FullFundraiserPath, hd.Secp256k1)
	require.Nil(t, err)

	address := "cosmos1pkkayn066msg6kn33wnl5srhdt3tnu2vzasz9c"
	var sdkAddress sdkTypes.AccAddress
	sdkAddress, err = sdkTypes.AccAddressFromBech32(address)

	// signWithout chainID
	requestBody1, err := Codec.MarshalJSON(request{
		BaseRequest: rest.BaseReq{From: address},
		Type:        "cosmos-sdk/StdTx",
		StdTx:       legacytx.NewStdTx([]sdkTypes.Msg{base.NewTestMessage(sdkAddress, "id")}, legacytx.NewStdFee(10, sdkTypes.NewCoins()), nil, ""),
	})
	require.Nil(t, err)
	testRequest1, err := http.NewRequest("POST", "/sign", bytes.NewBuffer(requestBody1))
	require.Nil(t, err)
	responseRecorder := httptest.NewRecorder()
	handler.ServeHTTP(responseRecorder, testRequest1)
	require.Equal(t, responseRecorder.Code, http.StatusBadRequest)
	require.Equal(t, `{"error":"Chain-ID required but not specified"}`, responseRecorder.Body.String())

	// with wrong key
	requestBody2, err := Codec.MarshalJSON(request{
		BaseRequest: rest.BaseReq{From: "address", ChainID: "test"},
		Type:        "cosmos-sdk/StdTx",
		StdTx:       legacytx.NewStdTx([]sdkTypes.Msg{base.NewTestMessage(sdkAddress, "id")}, legacytx.NewStdFee(20, sdkTypes.NewCoins()), nil, ""),
	})
	require.Nil(t, err)
	testRequest2, err := http.NewRequest("POST", "/sign", bytes.NewBuffer(requestBody2))
	require.Nil(t, err)
	responseRecorder = httptest.NewRecorder()
	handler.ServeHTTP(responseRecorder, testRequest2)
	require.Equal(t, responseRecorder.Code, http.StatusBadRequest)
	require.Equal(t, `{"error":"The specified item could not be found in the Keyring"}`, responseRecorder.Body.String())

	// RPC client offline
	requestBody3, err := Codec.MarshalJSON(request{
		BaseRequest: rest.BaseReq{From: address, ChainID: "test"},
		Type:        "cosmos-sdk/StdTx",
		StdTx:       legacytx.NewStdTx([]sdkTypes.Msg{base.NewTestMessage(sdkAddress, "id")}, legacytx.NewStdFee(30, sdkTypes.NewCoins()), nil, ""),
	})
	require.Nil(t, err)
	testRequest3, err := http.NewRequest("POST", "/sign", bytes.NewBuffer(requestBody3))
	require.Nil(t, err)
	responseRecorder = httptest.NewRecorder()
	handler.ServeHTTP(responseRecorder, testRequest3)
	require.Equal(t, responseRecorder.Code, http.StatusBadRequest)
	require.Equal(t, `{"error":"no RPC client is defined in offline mode"}`, responseRecorder.Body.String())

}
