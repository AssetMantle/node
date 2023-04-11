// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package queuing

import (
	"time"

	"github.com/cosmos/cosmos-sdk/client"
)

// kafkaConsumerMessages : messages to consume 5-second delay
func kafkaConsumerMessages(context client.Context) {
	quit := make(chan bool)

	var kafkaMsgList []kafkaMsg

	go func() {
		var msg kafkaMsg

		for {
			select {
			case <-quit:
				return
			default:
				msg = kafkaTopicConsumer("Topic", KafkaState.Consumers, context.LegacyAmino)
				if msg.Msg != nil {
					kafkaMsgList = append(kafkaMsgList, msg)
				}
			}
		}
	}()

	time.Sleep(sleepTimer)
	quit <- true

	if len(kafkaMsgList) == 0 {
		return
	}

	output, err := signAndBroadcastMultiple(kafkaMsgList, context)
	if err != nil {
		jsonError, e := context.LegacyAmino.MarshalJSON(struct {
			Error string `json:"error"`
		}{Error: err.Error()})
		if e != nil {
			panic(e)
		}

		for _, kafkaMsg := range kafkaMsgList {
			addResponseToDB(kafkaMsg.TicketID, jsonError, KafkaState.KafkaDB, context.LegacyAmino)
		}

		return
	}

	for _, kafkaMsg := range kafkaMsgList {
		addResponseToDB(kafkaMsg.TicketID, output, KafkaState.KafkaDB, context.LegacyAmino)
	}
}
