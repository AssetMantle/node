// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package queuing

import (
	"testing"

	schemaCodec "github.com/AssetMantle/schema/go/codec"
	"github.com/Shopify/sarama"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/stretchr/testify/require"
)

func TestKafkaProducerDeliverMessage(t *testing.T) {
	testProducer := []string{"testProducer"}
	var legacyAmino = codec.NewLegacyAmino()
	require.Panics(t, func() {
		schemaCodec.RegisterLegacyAminoCodec(legacyAmino)
		std.RegisterLegacyAminoCodec(legacyAmino)

		testKafkaMessage := kafkaMsg{Msg: nil}

		producer, _ := sarama.NewSyncProducer(testProducer, nil)
		// TODO: Add test cases.
		// require.Nil(t, err)

		require.Equal(t, kafkaProducerDeliverMessage(testKafkaMessage, "Topic", producer, legacyAmino), nil)
	})

}

func TestNewProducer(t *testing.T) {
	testProducer := []string{"testProducer"}

	producer, _ := sarama.NewSyncProducer(testProducer, nil)

	// TODO: Add test cases.
	// require.Nil(t, err, "should not happened. err %v", err)

	require.Panics(t, func() {
		require.Equal(t, newProducer(testProducer), producer)
	})
}
