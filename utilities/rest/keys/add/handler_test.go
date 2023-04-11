// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package add

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"

	"github.com/AssetMantle/schema/x"
)

func TestHandler(t *testing.T) {
	Codec := codec.NewLegacyAmino()
	x.RegisterLegacyAminoCodec(Codec)
	Codec.RegisterConcrete(request{}, "request", nil)
	Codec.RegisterConcrete(response{}, "response", nil)

	handler := handler(client.Context{})

	viper.Set(flags.FlagKeyringBackend, keyring.BackendTest)
	viper.Set(flags.FlagHome, t.TempDir())

	Keyring, err := keyring.New(sdk.KeyringServiceName(), keyring.BackendTest, t.TempDir(), strings.NewReader(""))
	require.NoError(t, err)

	router := mux.NewRouter()
	RegisterRESTRoutes(client.Context{}, router)

	t.Cleanup(func() {
		_ = Keyring.Delete("keyName1")
		_ = Keyring.Delete("keyName2")
		_ = Keyring.Delete("keyName3")
	})

	getResponse := func(responseBytes []byte) response {
		var Response rest.ResponseWithHeight

		innerErr := Codec.UnmarshalJSON(responseBytes, &Response)
		require.Nil(t, innerErr)

		var ResponseValue response
		innerErr = Codec.UnmarshalJSON(Response.Result, &ResponseValue)
		require.Nil(t, innerErr)

		return ResponseValue
	}

	// create account without mnemonic
	var requestBody1 []byte
	requestBody1, err = Codec.MarshalJSON(request{
		Name: "testKey1",
	})
	require.Nil(t, err)

	var testRequest1 *http.Request
	testRequest1, err = http.NewRequest("POST", "/keys/add", bytes.NewBuffer(requestBody1))
	require.Nil(t, err)

	responseRecorder := httptest.NewRecorder()
	handler.ServeHTTP(responseRecorder, testRequest1)
	require.Equal(t, http.StatusOK, responseRecorder.Code)

	response1 := getResponse(responseRecorder.Body.Bytes())
	require.Nil(t, response1.Error)
	require.Equal(t, true, response1.Success)

	// create account with mnemonic
	var requestBody2 []byte
	requestBody2, err = Codec.MarshalJSON(request{
		Name:     "testKey2",
		Mnemonic: "wage thunder live sense resemble foil apple course spin horse glass mansion midnight laundry acoustic rhythm loan scale talent push green direct brick please",
	})
	require.Nil(t, err)

	var testRequest2 *http.Request
	testRequest2, err = http.NewRequest("POST", "/keys/add", bytes.NewBuffer(requestBody2))
	require.Nil(t, err)

	responseRecorder = httptest.NewRecorder()
	handler.ServeHTTP(responseRecorder, testRequest2)
	require.Equal(t, responseRecorder.Code, http.StatusOK)

	response2 := getResponse(responseRecorder.Body.Bytes())
	require.Nil(t, response2.Error)
	require.Equal(t, true, response2.Success)

	// invalid mnemonic
	var requestBody3 []byte
	requestBody3, err = Codec.MarshalJSON(request{
		Name:     "testKey3",
		Mnemonic: "wage brick please",
	})
	require.Nil(t, err)

	var testRequest3 *http.Request
	testRequest3, err = http.NewRequest("POST", "/keys/add", bytes.NewBuffer(requestBody3))
	require.Nil(t, err)

	responseRecorder = httptest.NewRecorder()
	handler.ServeHTTP(responseRecorder, testRequest3)
	require.Equal(t, http.StatusInternalServerError, responseRecorder.Code)
	require.Equal(t, `{"error":"invalid mnemonic"}`, responseRecorder.Body.String())

	// invalid request
	testRequest4, err := http.NewRequest("POST", "/keys/add", bytes.NewBuffer([]byte{}))
	require.Nil(t, err)

	responseRecorder = httptest.NewRecorder()
	handler.ServeHTTP(responseRecorder, testRequest4)
	require.Equal(t, http.StatusBadRequest, responseRecorder.Code)

	// Retry adding same account
	var requestBody5 []byte
	requestBody5, err = Codec.MarshalJSON(request{
		Name: "testKey1",
	})
	require.Nil(t, err)

	var testRequest5 *http.Request
	testRequest5, err = http.NewRequest("POST", "/keys/add", bytes.NewBuffer(requestBody5))
	require.Nil(t, err)

	responseRecorder = httptest.NewRecorder()
	handler.ServeHTTP(responseRecorder, testRequest5)
	require.Equal(t, http.StatusInternalServerError, responseRecorder.Code)
	require.Equal(t, `{"error":"Account for keyname testKey1 already exists"}`, responseRecorder.Body.String())
}
