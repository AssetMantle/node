// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package add

import (
	"github.com/AssetMantle/modules/helpers"
)

type request struct {
	Name     string `json:"name"`
	Mnemonic string `json:"mnemonic"`
}

var _ helpers.Request = request{}

func (request request) Validate() error {
	return nil
}
