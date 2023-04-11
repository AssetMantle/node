// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package docs

import (
	"encoding/json"

	"github.com/asaskevich/govalidator"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/AssetMantle/modules/helpers"
)

type request struct {
	BaseReq                 rest.BaseReq `json:"baseReq"`
	FromID                  string       `json:"fromID" valid:"optional"`
	ImmutableMetaProperties string       `json:"immutableMetaProperties" valid:"optional"`
	ImmutableProperties     string       `json:"immutableProperties" valid:"optional"`
	MutableMetaProperties   string       `json:"mutableMetaProperties" valid:"optional"`
	MutableProperties       string       `json:"mutableProperties" valid:"optional"`
	ClassificationID        string       `json:"classificationID" valid:"optional"`
	MakerOwnableID          string       `json:"makerOwnableID" valid:"optional"`
	TakerOwnableID          string       `json:"takerOwnableID" valid:"optional"`
	MakerOwnableSplit       string       `json:"makerOwnableSplit" valid:"optional"`
	TakerOwnableSplit       string       `json:"takerOwnableSplit" valid:"optional"`
	ExpiresIn               string       `json:"expiresIn" valid:"optional"`
	Height                  string       `json:"height" valid:"optional"`
	TakerID                 string       `json:"takerID" valid:"optional"`
	NubID                   string       `json:"nubID" valid:"optional"`
	Coins                   string       `json:"coins" valid:"optional"`
}

var _ helpers.TransactionRequest = &request{}

func (request request) FromCLI(command helpers.CLICommand, context client.Context) (helpers.TransactionRequest, error) {
	// TODO implement me
	panic("implement me")
}

func (request) RegisterLegacyAminoCodec(legacyAmino *codec.LegacyAmino) {
	// TODO implement me
	panic("implement me")
}

func (request request) FromJSON(message json.RawMessage) (helpers.TransactionRequest, error) {
	// TODO implement me
	panic("implement me")
}

func (request request) MakeMsg() (sdkTypes.Msg, error) {
	// TODO implement me
	panic("implement me")
}

func (request request) Validate() error {
	_, err := govalidator.ValidateStruct(request)
	return err
}

func (request request) GetBaseReq() rest.BaseReq {
	return request.BaseReq
}

func Prototype() helpers.TransactionRequest {
	return request{}
}
