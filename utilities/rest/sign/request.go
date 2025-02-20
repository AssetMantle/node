// Copyright [2021] - [2025], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package sign

import (
	"github.com/AssetMantle/modules/helpers"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
)

type request struct {
	helpers.CommonTransactionRequest `json:"commonTransactionRequest"`
	Type                             string     `json:"type"`
	StdTx                            signing.Tx `json:"value"`
}
