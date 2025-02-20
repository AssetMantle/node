// Copyright [2021] - [2025], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package add

import (
	"encoding/json"
	"fmt"
	"github.com/AssetMantle/modules/utilities/rest"
	"io"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/go-bip39"
	"github.com/gorilla/mux"
)

func handler(context client.Context) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		body, err := io.ReadAll(httpRequest.Body)
		if err != nil {
			rest.WriteErrorResponse(responseWriter, http.StatusBadRequest, err.Error())
			return
		}

		request := request{}
		if err := json.Unmarshal(body, &request); err != nil {
			rest.WriteErrorResponse(responseWriter, http.StatusBadRequest, err.Error())
			return
		}

		if _, err := context.Keyring.Key(request.Name); err == nil {
			rest.WriteErrorResponse(responseWriter, http.StatusInternalServerError, fmt.Sprintf("Account for keyname %v already exists", request.Name))
			return
		}

		if request.Mnemonic != "" && !bip39.IsMnemonicValid(request.Mnemonic) {
			rest.WriteErrorResponse(responseWriter, http.StatusInternalServerError, "invalid mnemonic")
			return
		}

		if request.Mnemonic == "" {
			const mnemonicEntropySize = 256

			entropySeed, err := bip39.NewEntropy(mnemonicEntropySize)
			if err != nil {
				rest.WriteErrorResponse(responseWriter, http.StatusInternalServerError, err.Error())
				return
			}

			request.Mnemonic, err = bip39.NewMnemonic(entropySeed)
			if err != nil {
				rest.WriteErrorResponse(responseWriter, http.StatusInternalServerError, err.Error())
				return
			}
		}

		info, err := context.Keyring.NewAccount(request.Name, request.Mnemonic, keyring.DefaultBIP39Passphrase, sdkTypes.FullFundraiserPath, hd.Secp256k1)
		if err != nil {
			rest.WriteErrorResponse(responseWriter, http.StatusInternalServerError, err.Error())
			return
		}

		var keyOutput keyring.KeyOutput
		keyOutput, err = keyring.MkAccKeyOutput(info)
		if err != nil {
			rest.WriteErrorResponse(responseWriter, http.StatusInternalServerError, err.Error())
			return
		}

		keyOutput.Mnemonic = request.Mnemonic
		rest.PostProcessResponse(responseWriter, context, newResponse(keyOutput, nil))
	}
}

func RegisterRESTRoutes(context client.Context, router *mux.Router) {
	router.HandleFunc("/keys/add", handler(context)).Methods("POST")
}
