// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package queuing

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	"github.com/stretchr/testify/require"
	dbm "github.com/tendermint/tm-db"

	"github.com/AssetMantle/schema/utilities/random"
	"github.com/AssetMantle/schema/x"
	baseIDs "github.com/AssetMantle/schema/x/ids/base"
)

func Test_Kafka_DB(t *testing.T) {
	require.Panics(t, func() {
		var legacyAmino = codec.NewLegacyAmino()
		x.RegisterLegacyAminoCodec(legacyAmino)
		sdkTypes.RegisterCodec(legacyAmino)
		codec.RegisterCrypto(legacyAmino)
		codec.RegisterEvidences(legacyAmino)
		vesting.RegisterCodec(legacyAmino)
		legacyAmino.Seal()
		ticketID := TicketID(random.GenerateUniqueIdentifier("name"))
		kafkaDB, _ := dbm.NewGoLevelDB("KafkaDB", defaultCLIHome)
		setTicketIDtoDB(ticketID, kafkaDB, legacyAmino, []byte{})
		addResponseToDB(ticketID, baseIDs.NewStringID("").Bytes(), kafkaDB, legacyAmino)
		require.Equal(t, baseIDs.NewStringID("").Bytes(), getResponseFromDB(ticketID, kafkaDB, legacyAmino))
	})
}
