// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package make

import (
	"fmt"
	"github.com/AssetMantle/node/application"
	"github.com/AssetMantle/node/application/types/applications/base"
	abci "github.com/cometbft/cometbft/abci/types"
	simulationTypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestAppSimulationAfterImport(t *testing.T) {
	config, db, _, logger, skip, closeFn, err := setupRun(t, "leveldb-app-sim", "Simulation")
	defer closeFn()

	if skip {
		t.Skip("skipping application simulation after import")
	}

	require.NoError(t, err, "simulation setup failed")

	simulationApplication := base.NewSimulationApplication(logger, db, nil, true, map[int64]bool{}, application.Prototype.GetDefaultNodeHome(), 0, simapp.MakeTestEncodingConfig(), simapp.EmptyAppOptions{}, fauxMerkleModeOpt).(*base.SimulationApplication)

	require.Equal(t, "SimulationApplication", simulationApplication.Name())

	// Run randomized simulation
	stopEarly, simParams, simErr := simulation.SimulateFromSeed(
		t,
		os.Stdout,
		simulationApplication.GetBaseApp(),
		simapp.AppStateFn(simulationApplication.GetAppCodec(), simulationApplication.SimulationManager()),
		nil,
		simapp.SimulationOperations(simulationApplication, simulationApplication.GetAppCodec(), config),
		simulationApplication.ModuleAccountAddrs(),
		config,
		simulationApplication.GetAppCodec(),
	)

	// export state and simParams before the simulation error is checked
	err = simapp.CheckExportSimulation(simulationApplication, config, simParams)
	require.NoError(t, err)
	require.NoError(t, simErr)

	if config.Commit {
		simapp.PrintStats(db)
	}

	if stopEarly {
		fmt.Printf("can't export or import a zero-validator genesis, exiting test...")
		return
	}

	fmt.Printf("exporting genesis...\n")

	appState, err := simulationApplication.ExportAppStateAndValidators(true, []string{})
	require.NoError(t, err)

	fmt.Printf("importing genesis...\n")

	_, newDB, _, _, _, newCloseFn, err := setupRun(t, "leveldb-app-sim-2", "Simulation-2")
	defer newCloseFn()

	require.NoError(t, err, "simulation setup failed")

	newSimulationApplication := base.NewSimulationApplication(logger, newDB, nil, true, map[int64]bool{}, application.Prototype.GetDefaultNodeHome(), 0, simapp.MakeTestEncodingConfig(), simapp.EmptyAppOptions{}, fauxMerkleModeOpt).(*base.SimulationApplication)

	require.Equal(t, "SimulationApplication", newSimulationApplication.Name())

	newSimulationApplication.InitChain(abci.RequestInitChain{
		AppStateBytes: appState.AppState,
	})

	_, _, err = simulation.SimulateFromSeed(
		t,
		os.Stdout,
		newSimulationApplication.GetBaseApp(),
		simapp.AppStateFn(simulationApplication.GetAppCodec(), simulationApplication.SimulationManager()),
		simulationTypes.RandomAccounts,
		simapp.SimulationOperations(newSimulationApplication, newSimulationApplication.GetAppCodec(), config),
		newSimulationApplication.ModuleAccountAddrs(),
		config,
		simulationApplication.GetAppCodec(),
	)
	require.NoError(t, err)
}
