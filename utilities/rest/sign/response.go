// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package sign

import (
	"github.com/AssetMantle/modules/helpers"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
)

type response struct {
	Success bool       `json:"success"`
	Error   error      `json:"error"`
	Tx      signing.Tx `json:"tx"`
}

var _ helpers.Response = response{}

func (response response) IsSuccessful() bool {
	return response.Success
}
func (response response) GetError() error {
	return response.Error
}

func newResponse(tx signing.Tx, error error) helpers.Response {
	success := true
	if error != nil {
		success = false
	}

	return response{
		Success: success,
		Error:   error,
		Tx:      tx,
	}
}
