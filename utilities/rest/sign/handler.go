// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package sign

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authSigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func handler(context client.Context) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		var request request
		if !rest.ReadRESTReq(responseWriter, httpRequest, context.LegacyAmino, &request) {
			rest.WriteErrorResponse(responseWriter, http.StatusBadRequest, "")
			return
		}

		if request.BaseRequest.ChainID == "" {
			request.BaseRequest.ChainID = viper.GetString(flags.FlagChainID)
			if request.BaseRequest.ChainID == "" {
				rest.WriteErrorResponse(responseWriter, http.StatusBadRequest, "Chain-ID required but not specified")
				return
			}
		}

		context = context.WithFrom(request.BaseRequest.From).WithChainID(request.BaseRequest.ChainID)

		keyInfo, err := context.Keyring.Key(request.BaseRequest.From)
		if err != nil {
			rest.WriteErrorResponse(responseWriter, http.StatusBadRequest, err.Error())
			return
		}
		context = context.WithFromAddress(keyInfo.GetAddress())
		account, err := context.AccountRetriever.GetAccount(context, keyInfo.GetAddress())
		if err != nil {
			rest.WriteErrorResponse(responseWriter, http.StatusBadRequest, err.Error())
			return
		}

		txFactory := tx.Factory{}.WithKeybase(context.Keyring).WithTxConfig(context.TxConfig).WithAccountRetriever(context.AccountRetriever).WithAccountNumber(account.GetAccountNumber()).WithSequence(account.GetSequence()).WithFees(request.BaseRequest.Fees.String())
		txFactory = txFactory.WithChainID(request.BaseRequest.ChainID)
		txBuilder, err := tx.BuildUnsignedTx(txFactory, request.StdTx.GetMsgs()...)
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
		signatureBytes, _, err := context.Keyring.Sign(request.BaseRequest.From, bytesToSign)
		if err != nil {
			rest.WriteErrorResponse(responseWriter, http.StatusBadRequest, err.Error())
			return
		}

		signatureData = signing.SingleSignatureData{
			SignMode:  signMode,
			Signature: signatureBytes,
		}
		signature = signing.SignatureV2{
			PubKey:   keyInfo.GetPubKey(),
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
