/*
 Copyright [2019] - [2020], PERSISTENCE TECHNOLOGIES PTE. LTD. and the assetMantle contributors
 SPDX-License-Identifier: Apache-2.0
*/

package application

import (
	"github.com/persistenceOne/assetMantle/application/internal/configurations"
	"github.com/persistenceOne/persistenceSDK/schema/applications/base"
)

var Prototype = base.NewApplication(
	configurations.Name,
	configurations.ModuleBasicManager,
	configurations.EnabledWasmProposalTypeList,
	configurations.ModuleAccountPermissions,
	configurations.TokenReceiveAllowedModules,
)
