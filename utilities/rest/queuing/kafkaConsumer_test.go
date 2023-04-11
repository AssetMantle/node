// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package queuing

import (
	"testing"

	"github.com/Shopify/sarama"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/stretchr/testify/require"

	"github.com/AssetMantle/schema/x"
)

func TestKafkaTopicConsumer(t *testing.T) {
	testConsumers := []string{"testConsumers"}

	var legacyAmino = codec.NewLegacyAmino()

	x.RegisterLegacyAminoCodec(legacyAmino)
	std.RegisterLegacyAminoCodec(legacyAmino)

	require.Panics(t, func() {
		testKafkaState := NewKafkaState(testConsumers)
		partitionConsumer := testKafkaState.Consumers["Topic"]

		var kafkaStore kafkaMsg
		if len(partitionConsumer.Messages()) == 0 {
			kafkaStore = kafkaMsg{Msg: nil}
		}

		kafkaMsg := <-partitionConsumer.Messages()

		err := legacyAmino.UnmarshalJSON(kafkaMsg.Value, &kafkaStore)
		if err != nil {
			panic(err)
		}

		require.Equal(t, kafkaTopicConsumer("Topic", testKafkaState.Consumers, legacyAmino), kafkaStore)
	})
}

func TestNewConsumer(t *testing.T) {
	consumers := []string{"testConsumers"}
	config := sarama.NewConfig()

	consumer, _ := sarama.NewConsumer(consumers, config)

	// TODO: Add test cases.
	// require.Nil(t, err, "should not happened. err %v", err)

	require.Panics(t, func() {
		require.Equal(t, newConsumer(consumers), consumer)
	})
}
