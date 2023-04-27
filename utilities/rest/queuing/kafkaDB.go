// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package queuing

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/gorilla/mux"
	dbm "github.com/tendermint/tm-db"
)

// setTicketIDtoDB : initiates TicketID in Database
func setTicketIDtoDB(ticket TicketID, kafkaDB *dbm.GoLevelDB, legacyAmino *codec.LegacyAmino, msg []byte) {
	ticketID, err := legacyAmino.MarshalJSON(ticket)
	if err != nil {
		panic(err)
	}

	if err := kafkaDB.Set(ticketID, msg); err != nil {
		panic(err)
	}
}

// addResponseToDB : Updates response to DB
func addResponseToDB(ticket TicketID, response []byte, kafkaDB *dbm.GoLevelDB, legacyAmino *codec.LegacyAmino) {
	ticketID, err := legacyAmino.MarshalJSON(ticket)
	if err != nil {
		panic(err)
	}

	err = kafkaDB.SetSync(ticketID, response)
	if err != nil {
		panic(err)
	}
}

// getResponseFromDB : gives the response from DB
func getResponseFromDB(ticket TicketID, kafkaDB *dbm.GoLevelDB, legacyAmino *codec.LegacyAmino) []byte {
	ticketID, err := legacyAmino.MarshalJSON(ticket)
	if err != nil {
		panic(err)
	}

	val, _ := kafkaDB.Get(ticketID)

	return val
}

// queryDB : REST outputs info from DB
func queryDB(legacyAmino *codec.LegacyAmino) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		vars := mux.Vars(r)

		ticketIDBytes, err := legacyAmino.MarshalJSON(vars["TicketID"])
		if err != nil {
			panic(err)
		}

		var response []byte

		check, _ := KafkaState.KafkaDB.Has(ticketIDBytes)
		if check {
			response = getResponseFromDB(TicketID(vars["TicketID"]), KafkaState.KafkaDB, legacyAmino)
		} else {
			output, err := legacyAmino.MarshalJSON("The ticket ID does not exist, it must have been deleted, Query the chain to know")
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write([]byte(fmt.Sprintf("ticket ID does not exist. Error: %s", err.Error())))
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(output)
			return
		}

		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write(response)
	}
}
