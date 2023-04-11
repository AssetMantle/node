// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package queuing

import (
	"os"
	"time"
)

// sleepTimer : the time the kafka messages are to be taken in
const sleepTimer = time.Duration(1000000000)

// sleepRoutine : the time the kafka messages are to be taken in
const sleepRoutine = time.Duration(2500000000)

// defaultCLIHome : is the home path
var defaultCLIHome = os.ExpandEnv("$HOME/.kafka")

const partition = int32(0)
const offset = int64(0)

// topics : is list of topics
var topics = []string{
	"Topic",
}
