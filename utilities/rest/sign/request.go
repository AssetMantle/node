// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package sign

import (
	"github.com/AssetMantle/modules/helpers"
	"github.com/AssetMantle/modules/utilities/rest"

	"github.com/cosmos/cosmos-sdk/x/auth/signing"
)

type request struct {
	BaseRequest rest.BaseReq `json:"baseReq"`
	Type        string       `json:"type"`
	StdTx       signing.Tx   `json:"value"`
}

var _ helpers.Request = request{}

func (request request) Validate() error {
	return nil
}
