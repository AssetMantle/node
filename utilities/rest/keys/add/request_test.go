// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package add

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Add_Request(t *testing.T) {
	require.Equal(t, nil, request{Name: "name", Mnemonic: "mnemonic"}.Validate())
}
