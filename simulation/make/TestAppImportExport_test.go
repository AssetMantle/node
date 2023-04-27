// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package make

import (
	"encoding/json"
	"fmt"
	"github.com/AssetMantle/node/application"
	"github.com/AssetMantle/node/application/types/applications/base"
	simulationTypes "github.com/cosmos/cosmos-sdk/types/simulation"

	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	mintTypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	slashingTypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"os"
	"testing"

	"github.com/AssetMantle/modules/x/assets"
	"github.com/AssetMantle/modules/x/classifications"
	"github.com/AssetMantle/modules/x/identities"
	"github.com/AssetMantle/modules/x/maintainers"
	"github.com/AssetMantle/modules/x/metas"
	"github.com/AssetMantle/modules/x/orders"
	"github.com/AssetMantle/modules/x/splits"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramsTypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/stretchr/testify/require"
)

// GetProperty flags every time the simulator is run
func init() {
	simapp.GetSimulatorFlags()
}

type StoreKeysPrefixes struct {
	A        sdk.StoreKey
	B        sdk.StoreKey
	Prefixes [][]byte
}

// fauxMerkleModeOpt returns a BaseApp option to use a dbStoreAdapter instead of
// an IAVLStore for faster simulation speed.
func fauxMerkleModeOpt(baseApplication *baseapp.BaseApp) {
	baseApplication.SetFauxMerkleMode()
}

func interBlockCacheOpt() func(*baseapp.BaseApp) {
	return baseapp.SetInterBlockCache(store.NewCommitKVStoreCacheManager())
}

func TestAppImportExport(t *testing.T) {
	config, db, _, logger, skip, closeFn, err := setupRun(t, "leveldb-app-sim", "Simulation")
	defer closeFn()

	if skip {
		t.Skip("skipping application import/export simulation")
	}

	require.NoError(t, err, "simulation setup failed")

	simulationApplication := base.NewSimulationApplication(logger, db, nil, true, map[int64]bool{}, application.Prototype.GetDefaultNodeHome(), 0, simapp.MakeTestEncodingConfig(), simapp.EmptyAppOptions{}, fauxMerkleModeOpt).(*base.SimulationApplication)

	require.Equal(t, "Simulation", simulationApplication.Name())

	// Run randomized simulation
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

	fmt.Printf("exporting genesis...\n")

	appState, err := simulationApplication.ExportAppStateAndValidators(false, []string{})
	require.NoError(t, err)

	fmt.Printf("importing genesis...\n")

	_, newDB, _, logger, _, newCloseFn, err := setupRun(t, "leveldb-app-sim-2", "Simulation-2")
	defer newCloseFn()

	require.NoError(t, err, "simulation setup failed")

	newSimulationApplication := base.NewSimulationApplication(logger, newDB, nil, true, map[int64]bool{}, application.Prototype.GetDefaultNodeHome(), 0, simapp.MakeTestEncodingConfig(), simapp.EmptyAppOptions{}, fauxMerkleModeOpt).(*base.SimulationApplication)

	require.Equal(t, "SimulationApplication", newSimulationApplication.Name())

	var genesisState simapp.GenesisState
	err = json.Unmarshal(appState.AppState, &genesisState)
	require.NoError(t, err)

	ctxA := newSimulationApplication.NewContext(true, tmproto.Header{Height: newSimulationApplication.LastBlockHeight()})
	ctxB := newSimulationApplication.NewContext(true, tmproto.Header{Height: newSimulationApplication.LastBlockHeight()})

	newSimulationApplication.GetModuleManager().InitGenesis(ctxB, newSimulationApplication.GetAppCodec(), genesisState)

	fmt.Printf("comparing stores...\n")

	storeKeysPrefixes := []StoreKeysPrefixes{
		//{simulationApplication.GetKey(baseapp.MainStoreKey), newSimulationApplication.GetKey(baseapp.MainStoreKey), [][]byte{}},
		{simulationApplication.GetKey(authTypes.StoreKey), newSimulationApplication.GetKey(authTypes.StoreKey), [][]byte{}},
		{simulationApplication.GetKey(stakingTypes.StoreKey), newSimulationApplication.GetKey(stakingTypes.StoreKey),
			[][]byte{
				stakingTypes.UnbondingQueueKey, stakingTypes.RedelegationQueueKey, stakingTypes.ValidatorQueueKey,
			}}, // ordering may change but it doesn't matter
		{simulationApplication.GetKey(slashingTypes.StoreKey), newSimulationApplication.GetKey(slashingTypes.StoreKey), [][]byte{}},
		{simulationApplication.GetKey(mintTypes.StoreKey), newSimulationApplication.GetKey(mintTypes.StoreKey), [][]byte{}},
		{simulationApplication.GetKey(distributionTypes.StoreKey), newSimulationApplication.GetKey(distributionTypes.StoreKey), [][]byte{}},
		{simulationApplication.GetKey(paramsTypes.StoreKey), newSimulationApplication.GetKey(paramsTypes.StoreKey), [][]byte{}},
		{simulationApplication.GetKey(govTypes.StoreKey), newSimulationApplication.GetKey(govTypes.StoreKey), [][]byte{}},
		{simulationApplication.GetKey(assets.Prototype().Name()), newSimulationApplication.GetKey(assets.Prototype().Name()), [][]byte{}},
		{simulationApplication.GetKey(classifications.Prototype().Name()), newSimulationApplication.GetKey(classifications.Prototype().Name()), [][]byte{}},
		{simulationApplication.GetKey(identities.Prototype().Name()), newSimulationApplication.GetKey(identities.Prototype().Name()), [][]byte{}},
		{simulationApplication.GetKey(maintainers.Prototype().Name()), newSimulationApplication.GetKey(maintainers.Prototype().Name()), [][]byte{}},
		{simulationApplication.GetKey(metas.Prototype().Name()), newSimulationApplication.GetKey(metas.Prototype().Name()), [][]byte{}},
		{simulationApplication.GetKey(orders.Prototype().Name()), newSimulationApplication.GetKey(orders.Prototype().Name()), [][]byte{}},
		{simulationApplication.GetKey(splits.Prototype().Name()), newSimulationApplication.GetKey(splits.Prototype().Name()), [][]byte{}},
	}

	for _, skp := range storeKeysPrefixes {
		storeA := ctxA.KVStore(skp.A)
		storeB := ctxB.KVStore(skp.B)

		failedKVAs, failedKVBs := sdk.DiffKVStores(storeA, storeB, skp.Prefixes)
		require.Equal(t, len(failedKVAs), len(failedKVBs), "unequal sets of key-values to compare")

		fmt.Printf("compared %d key/value pairs between %s and %s\n", len(failedKVAs), skp.A, skp.B)
		require.Equal(t, len(failedKVAs), 0, simapp.GetSimulationLog(skp.A.Name(), simulationApplication.SimulationManager().StoreDecoders, failedKVAs, failedKVBs))
	}
}
