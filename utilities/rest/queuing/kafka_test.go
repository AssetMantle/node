// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package queuing

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/std"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/stretchr/testify/require"

	codecUtilities "github.com/AssetMantle/schema/utilities/codec"
	"github.com/AssetMantle/schema/utilities/random"
	"github.com/AssetMantle/schema/x"
)

type testMessage struct {
	Name string `json:"name"`
}

var _ sdkTypes.Msg = testMessage{}

func (message testMessage) Route() string { return "testModule" }
func (message testMessage) Type() string  { return "" }
func (message testMessage) ValidateBasic() error {
	return nil
}
func (message testMessage) GetSignBytes() []byte {
	return []byte{}
}
func (message testMessage) GetSigners() []sdkTypes.AccAddress {
	fromAddress := "cosmos1pkkayn066msg6kn33wnl5srhdt3tnu2vzasz9c"
	fromAccAddress, _ := sdkTypes.AccAddressFromBech32(fromAddress)
	return []sdkTypes.AccAddress{fromAccAddress}
}
func (testMessage) RegisterLegacyAminoCodec(legacyAmino *codec.LegacyAmino) {
	codecUtilities.RegisterModuleConcrete(legacyAmino, testMessage{})
}

func Test_Kafka(t *testing.T) {
	var legacyAmino = codec.NewLegacyAmino()
	x.RegisterLegacyAminoCodec(legacyAmino)
	std.RegisterLegacyAminoCodec(legacyAmino)

	fromAddress := "cosmos1pkkayn066msg6kn33wnl5srhdt3tnu2vzasz9c"
	fromAccAddress, err := sdkTypes.AccAddressFromBech32(fromAddress)
	require.Nil(t, err)
	testBaseReq := rest.BaseReq{From: fromAddress, ChainID: "test"}
	ticketID := TicketID(random.GenerateUniqueIdentifier("ticket"))
	kafkaPorts := []string{"localhost:9092"}
	require.Panics(t, func() {
		testKafkaState := NewKafkaState(kafkaPorts)
		bank.RegisterCodec(legacyAmino)
		message := bank.NewMsgSend(fromAccAddress, fromAccAddress, sdkTypes.NewCoins(sdkTypes.NewCoin("stake", sdkTypes.NewInt(123))))

		testKafkaMsg := NewKafkaMsgFromRest(message, ticketID, testBaseReq, context)
		SendToKafka(testKafkaMsg, legacyAmino)

		kafkaMsg := kafkaTopicConsumer("Topic", testKafkaState.Consumers, legacyAmino)
		require.Equal(t, testKafkaMsg.TicketID, kafkaMsg.TicketID)
		require.Equal(t, testKafkaMsg.BaseRequest, kafkaMsg.BaseRequest)
	})

}
