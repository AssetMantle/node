// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package make

import (
	"os"
	"testing"

	"github.com/cosmos/cosmos-sdk/simapp"
	simulationTypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/stretchr/testify/require"

	"github.com/AssetMantle/node/application/types/applications/base"
)

// Profile with:
// /usr/local/go/bin/go test -benchmem -run=^$ github.com/AssetMantle/modules/simapp -bench ^BenchmarkFullAppSimulation$ -Commit=true -cpuprofile cpu.out
func BenchmarkFullAppSimulation(b *testing.B) {
	config, db, _, logger, _, closeFn, err := setupRun(b, "goleveldb-app-sim", "Simulation")
	defer closeFn()

	require.NoError(b, err, "simulation setup failed")

	simulationApplication := base.NewSimulationApplication(logger, db, nil, true, map[int64]bool{}, DefaultNodeHome, simapp.FlagPeriodValue, simapp.MakeTestEncodingConfig(), simapp.EmptyAppOptions{}, interBlockCacheOpt())

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
}
