// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package application

import (
	"github.com/AssetMantle/node/application/internal/configurations"
	"github.com/AssetMantle/node/application/types/applications/base"
)

var Prototype = base.NewApplication(
	configurations.Name,
	configurations.ModuleBasicManager,
	configurations.ModuleAccountPermissions,
	configurations.TokenReceiveAllowedModules,
)
