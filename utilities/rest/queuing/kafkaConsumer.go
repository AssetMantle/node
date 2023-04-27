// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package queuing

import (
	"github.com/Shopify/sarama"
	"github.com/cosmos/cosmos-sdk/codec"
)

// newConsumer : is a consumer which is needed to create child consumers to consume topics
func newConsumer(kafkaNodes []string) sarama.Consumer {
	config := sarama.NewConfig()

	consumer, err := sarama.NewConsumer(kafkaNodes, config)
	if err != nil {
		panic(err)
	}

	return consumer
}

// partitionConsumers : is a child consumer
func partitionConsumers(consumer sarama.Consumer, topic string) sarama.PartitionConsumer {
	// partition and offset defined in configurations.go
	partitionConsumer, err := consumer.ConsumePartition(topic, partition, offset)
	if err != nil {
		panic(err)
	}

	return partitionConsumer
}

// kafkaTopicConsumer : Takes a consumer and makes it consume a topic message at a time
func kafkaTopicConsumer(topic string, consumers map[string]sarama.PartitionConsumer, legacyAmino *codec.LegacyAmino) kafkaMsg {
	partitionConsumer := consumers[topic]

	if len(partitionConsumer.Messages()) == 0 {
		return kafkaMsg{Msg: nil}
	}

	var consumedKafkaMsg kafkaMsg
	err := legacyAmino.UnmarshalJSON((<-partitionConsumer.Messages()).Value, &consumedKafkaMsg)
	if err != nil {
		panic(err)
	}

	return consumedKafkaMsg
}
