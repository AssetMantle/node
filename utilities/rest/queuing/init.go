// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package queuing

import (
	"time"

	"github.com/cosmos/cosmos-sdk/client"
)

func InitializeKafka(nodeList []string, context client.Context) {
	KafkaState = *NewKafkaState(nodeList)
	if KafkaState.IsEnabled {
		go func() {
			for {
				kafkaConsumerMessages(context)
				time.Sleep(sleepRoutine)
			}
		}()
	}
}
