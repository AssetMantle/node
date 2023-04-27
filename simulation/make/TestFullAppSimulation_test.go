// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package make

import (
	"github.com/AssetMantle/node/application"
	simulationTypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"os"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/stretchr/testify/require"

	"github.com/AssetMantle/node/application/types/applications/base"
)

func TestFullAppSimulation(t *testing.T) {
	simapp.FlagEnabledValue = true
	simapp.FlagCommitValue = true
	simapp.FlagSeedValue = time.Now().UnixNano()

	config, db, _, logger, skip, closeFn, err := setupRun(t, "leveldb-app-sim", "Simulation")
	defer closeFn()
	if skip {
		t.Skip("skipping application simulation")
	}
	require.NoError(t, err, "simulation setup failed")

	simulationApplication := base.NewSimulationApplication(logger, db, nil, true, map[int64]bool{}, application.Prototype.GetDefaultNodeHome(), 0, simapp.MakeTestEncodingConfig(), simapp.EmptyAppOptions{}, fauxMerkleModeOpt).(*base.SimulationApplication)

	require.Equal(t, "Simulation", simulationApplication.Name())

	// run randomized simulation
	_, simParams, simErr := simulation.SimulateFromSeed(
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

	// export state and simParams before the simulation error is checked
	err = simapp.CheckExportSimulation(simulationApplication, config, simParams)
	require.NoError(t, err)
	require.NoError(t, simErr)

	if config.Commit {
		simapp.PrintStats(db)
	}
}
