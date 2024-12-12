// Copyright [2021] - [2025], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package make

import (
	"fmt"
	"github.com/AssetMantle/node/application/types/applications/constants"
	"os"
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	simulationTypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/stretchr/testify/require"

	"github.com/AssetMantle/node/application/types/applications/base"
)

func BenchmarkInvariants(b *testing.B) {
	config, db, _, logger, _, closeFn, err := setupRun(b, "leveldb-app-invariant-bench", "Simulation")
	defer closeFn()

	require.NoError(b, err, "simulation setup failed")

	config.AllInvariants = false

	simulationApplication := base.NewSimulationApplication(logger, db, nil, true, map[int64]bool{}, constants.DefaultNodeHome, simapp.FlagPeriodValue, simapp.MakeTestEncodingConfig(), simapp.EmptyAppOptions{}, interBlockCacheOpt())

	// run randomized simulation
	_, simParams, simErr := simulation.SimulateFromSeed(
		b,
		os.Stdout,
		simulationApplication.GetBaseApp(),
		simapp.AppStateFn(simulationApplication.GetAppCodec(), simulationApplication.SimulationManager()),
		simulationTypes.RandomAccounts,
		simapp.SimulationOperations(simulationApplication, simulationApplication.GetAppCodec(), config),
		simulationApplication.ModuleAccountAddrs(),
		config,
		simulationApplication.GetAppCodec(),
	)
	// export state and simParams before the simulation error is checked
	if err = simapp.CheckExportSimulation(simulationApplication, config, simParams); err != nil {
		b.Fatal(err)
	}

	if simErr != nil {
		b.Fatal(simErr)
	}

	if config.Commit {
		simapp.PrintStats(db)
	}

	ctx := simulationApplication.GetBaseApp().NewContext(true, tmproto.Header{Height: simulationApplication.GetBaseApp().LastBlockHeight() + 1})

	// 3. Benchmark each invariant separately
	//
	// NOTE: We use the crisis keeper as it has all the invariants registered with
	// their respective metadata which makes it useful for testing/benchmarking.
	for _, cr := range simulationApplication.GetCrisisKeeper().Routes() {
		cr := cr

		b.Run(fmt.Sprintf("%s/%s", cr.ModuleName, cr.Route), func(b *testing.B) {
			if res, stop := cr.Invar(ctx); stop {
				b.Fatalf(
					"broken invariant at block %d of %d\n%s",
					ctx.BlockHeight()-1, config.NumBlocks, res,
				)
			}
		})
	}
}
