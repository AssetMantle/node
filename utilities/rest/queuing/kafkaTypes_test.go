// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package queuing

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/std"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/stretchr/testify/require"

	"github.com/AssetMantle/schema/utilities/random"
	"github.com/AssetMantle/schema/x"
)

func Test_Kafka_Types(t *testing.T) {
	var legacyAmino = codec.NewLegacyAmino()
	x.RegisterLegacyAminoCodec(legacyAmino)
	std.RegisterLegacyAminoCodec(legacyAmino)
	legacyAmino.Seal()

	fromAddress := "cosmos1pkkayn066msg6kn33wnl5srhdt3tnu2vzasz9c"
	testBaseReq := rest.BaseReq{From: fromAddress, ChainID: "test", Fees: sdkTypes.NewCoins()}

	testMessage := sdkTypes.NewTestMsg()

	ticketID := TicketID(random.GenerateUniqueIdentifier("name"))
	testKafkaMsg := NewKafkaMsgFromRest(testMessage, ticketID, testBaseReq, context)
	kafkaCliCtx := kafkaCliCtx{
		OutputFormat:  context.OutputFormat,
		ChainID:       context.ChainID,
		Height:        context.Height,
		HomeDir:       context.HomeDir,
		NodeURI:       context.NodeURI,
		From:          context.From,
		TrustNode:     context.TrustNode,
		UseLedger:     context.UseLedger,
		BroadcastMode: context.BroadcastMode,
		Simulate:      context.Simulate,
		GenerateOnly:  context.GenerateOnly,
		FromAddress:   context.FromAddress,
		FromName:      context.FromName,
		Indent:        context.Indent,
		SkipConfirm:   context.SkipConfirm,
	}
	require.Equal(t, kafkaMsg{Msg: testMessage, TicketID: ticketID, BaseRequest: testBaseReq, KafkaCliCtx: kafkaCliCtx}, testKafkaMsg)
	require.Equal(t, context, cliCtxFromKafkaMsg(testKafkaMsg, context))
	// require
}
