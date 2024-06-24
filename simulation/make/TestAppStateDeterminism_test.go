// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package make

import (
	"encoding/json"
	"fmt"
	"github.com/AssetMantle/node/application"
	simulationTypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"math/rand"
	"os"
	"testing"

	dbm "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/stretchr/testify/require"

	"github.com/AssetMantle/node/application/types/applications/base"
)

func TestAppStateDeterminism(t *testing.T) {
	if !simapp.FlagEnabledValue {
		t.Skip("skipping application simulation")
	}

	config := simapp.NewConfigFromFlags()
	config.InitialBlockHeight = 1
	config.ExportParamsPath = ""
	config.OnOperation = false
	config.AllInvariants = false
	config.ChainID = helpers.SimAppChainID

	numSeeds := 3
	numTimesToRunPerSeed := 5
	appHashList := make([]json.RawMessage, numTimesToRunPerSeed)

	for i := 0; i < numSeeds; i++ {
		config.Seed = rand.Int63()

		for j := 0; j < numTimesToRunPerSeed; j++ {
			var logger log.Logger
			if simapp.FlagVerboseValue {
				logger = log.TestingLogger()
			} else {
				logger = log.NewNopLogger()
			}

			db := dbm.NewMemDB()

			simulationApplication := base.NewSimulationApplication(logger, db, nil, true, map[int64]bool{}, application.Prototype.GetDefaultNodeHome(), 0, simapp.MakeTestEncodingConfig(), simapp.EmptyAppOptions{}, fauxMerkleModeOpt).(*base.SimulationApplication)

			fmt.Printf(
				"running non-determinism simulation; seed %d: %d/%d, attempt: %d/%d\n",
				config.Seed, i+1, numSeeds, j+1, numTimesToRunPerSeed,
			)

			_, _, err := simulation.SimulateFromSeed(
				t,
				os.Stdout,
				simulationApplication.GetBaseApp(),
				simapp.AppStateFn(simulationApplication.GetAppCodec(), simulationApplication.SimulationManager()),
				simulationTypes.RandomAccounts,
				simapp.SimulationOperations(simulationApplication, simulationApplication.GetAppCodec(), config),
				simulationApplication.ModuleAccountAddrs(),
				config,
				simulationApplication.GetAppCodec(),
			)
			require.NoError(t, err)

			if config.Commit {
				simapp.PrintStats(db)
			}

			appHash := simulationApplication.GetBaseApp().LastCommitID().Hash
			appHashList[j] = appHash

			if j != 0 {
				require.Equal(
					t, appHashList[0], appHashList[j],
					"non-determinism in seed %d: %d/%d, attempt: %d/%d\n", config.Seed, i+1, numSeeds, j+1, numTimesToRunPerSeed,
				)
			}
		}
	}
}
