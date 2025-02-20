// Copyright [2021] - [2025], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package make

import (
	"os"

	db "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	simulationTypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

type testingLogging interface {
	Error(...interface{})
	Fatal(...interface{})
}

func setupRun(t testingLogging, dirPrefix, dbName string) (config simulationTypes.Config, db db.DB, dir string, logger log.Logger, skip bool, closeFn func(), err error) {
	closeFn = func() {}

	config, db, dir, logger, skip, err = simapp.SetupSimulation(dirPrefix, dbName)
	if err != nil {
		return
	}

	closeFn = func() {
		var gotError bool

		if db != nil {
			err = db.Close()
			if err != nil {
				gotError = true

				t.Error(err)
			}
		}

		err = os.RemoveAll(dir)
		if err != nil {
			gotError = true

			t.Error(err)
		}

		if gotError {
			t.Fatal()
		}
	}

	return
}
