// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package add

import (
	"github.com/asaskevich/govalidator"

	"github.com/AssetMantle/modules/helpers"
)

type request struct {
	Name     string `json:"name" valid:"required~required field to missing, matches(.+?)~invalid field name"`
	Mnemonic string `json:"mnemonic" valid:"optional"`
}

var _ helpers.Request = request{}

func (request request) Validate() error {
	_, err := govalidator.ValidateStruct(request)
	return err
}
