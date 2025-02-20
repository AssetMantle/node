// Copyright [2021] - [2025], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package sign

import (
	"encoding/json"
	"github.com/AssetMantle/modules/utilities/rest"
	"io"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"

	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authSigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
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

		if request.CommonTransactionRequest.GetChainID() == "" {
			request.CommonTransactionRequest = request.CommonTransactionRequest.SetChainID(viper.GetString(flags.FlagChainID))
			if request.CommonTransactionRequest.GetChainID() == "" {
				rest.WriteErrorResponse(responseWriter, http.StatusBadRequest, "Chain-ID required but not specified")
				return
			}
		}

		context = context.WithFrom(request.CommonTransactionRequest.GetFrom()).WithChainID(request.CommonTransactionRequest.GetChainID())

		keyInfo, err := context.Keyring.Key(request.CommonTransactionRequest.GetFrom())
		if err != nil {
			rest.WriteErrorResponse(responseWriter, http.StatusBadRequest, err.Error())
			return
		}

		address, err := keyInfo.GetAddress()
		if err != nil {
			rest.WriteErrorResponse(responseWriter, http.StatusBadRequest, err.Error())
			return
		}
		context = context.WithFromAddress(address)

		account, err := context.AccountRetriever.GetAccount(context, address)
		if err != nil {
			rest.WriteErrorResponse(responseWriter, http.StatusBadRequest, err.Error())
			return
		}

		txFactory := tx.Factory{}.WithKeybase(context.Keyring).WithTxConfig(context.TxConfig).WithAccountRetriever(context.AccountRetriever).WithAccountNumber(account.GetAccountNumber()).WithSequence(account.GetSequence()).WithFees(request.CommonTransactionRequest.GetFees().String())
		txFactory = txFactory.WithChainID(request.CommonTransactionRequest.GetChainID())
		txBuilder, err := txFactory.BuildUnsignedTx(request.StdTx.GetMsgs()...)
		if err != nil {
			rest.WriteErrorResponse(responseWriter, http.StatusBadRequest, err.Error())
			return
		}

		signerData := authSigning.SignerData{
			ChainID:       context.ChainID,
			AccountNumber: account.GetAccountNumber(),
			Sequence:      account.GetSequence(),
		}
		signatureData := signing.SingleSignatureData{
			SignMode:  context.TxConfig.SignModeHandler().DefaultMode(),
			Signature: nil,
		}
		signature := signing.SignatureV2{
			PubKey:   account.GetPubKey(),
			Data:     &signatureData,
			Sequence: txFactory.Sequence(),
		}

		signMode := context.TxConfig.SignModeHandler().DefaultMode()
		bytesToSign, err := context.TxConfig.SignModeHandler().GetSignBytes(signMode, signerData, txBuilder.GetTx())
		if err != nil {
			rest.WriteErrorResponse(responseWriter, http.StatusBadRequest, err.Error())
			return
		}
		signatureBytes, _, err := context.Keyring.Sign(request.CommonTransactionRequest.GetFrom(), bytesToSign)
		if err != nil {
			rest.WriteErrorResponse(responseWriter, http.StatusBadRequest, err.Error())
			return
		}

		signatureData = signing.SingleSignatureData{
			SignMode:  signMode,
			Signature: signatureBytes,
		}

		pubKey, err := keyInfo.GetPubKey()
		if err != nil {
			rest.WriteErrorResponse(responseWriter, http.StatusBadRequest, err.Error())
			return
		}
		signature = signing.SignatureV2{
			PubKey:   pubKey,
			Data:     &signatureData,
			Sequence: txFactory.Sequence(),
		}

		if err := txBuilder.SetSignatures(signature); err != nil {
			rest.WriteErrorResponse(responseWriter, http.StatusBadRequest, err.Error())
			return
		}

		rest.PostProcessResponse(responseWriter, context, newResponse(txBuilder.GetTx(), nil))
	}
}

func prepareFactory(clientCtx client.Context, txf tx.Factory) (tx.Factory, error) {
	from := clientCtx.GetFromAddress()

	if err := txf.AccountRetriever().EnsureExists(clientCtx, from); err != nil {
		return txf, err
	}

	initNum, initSeq := txf.AccountNumber(), txf.Sequence()
	if initNum == 0 || initSeq == 0 {
		num, seq, err := txf.AccountRetriever().GetAccountNumberSequence(clientCtx, from)
		if err != nil {
			return txf, err
		}

		if initNum == 0 {
			txf = txf.WithAccountNumber(num)
		}

		if initSeq == 0 {
			txf = txf.WithSequence(seq)
		}
	}

	return txf, nil
}

func RegisterRESTRoutes(context client.Context, router *mux.Router) {
	router.HandleFunc("/sign", handler(context)).Methods("POST")
}
