/*
 Copyright [2019] - [2020], PERSISTENCE TECHNOLOGIES PTE. LTD. and the assetMantle contributors
 SPDX-License-Identifier: Apache-2.0
*/

package initialize

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/persistenceOne/assetMantle/application"
	"github.com/spf13/viper"
	"io"
	"os"
	"path/filepath"
	"time"

	cpm "github.com/otiai10/copy"
	"github.com/spf13/cobra"

	abci "github.com/tendermint/tendermint/abci/types"
	tendermintOS "github.com/tendermint/tendermint/libs/os"
	"github.com/tendermint/tendermint/proxy"
	tmsm "github.com/tendermint/tendermint/state"
	tmstore "github.com/tendermint/tendermint/store"
	tm "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func ReplayTransactionsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "replay <root-dir>",
		Short: "Replay hub transactions",
		RunE: func(_ *cobra.Command, args []string) error {
			return replayTransactions(args[0])
		},
		Args: cobra.ExactArgs(1),
	}
}

func replayTransactions(rootDir string) error {

	if false {
		oldRootDir := rootDir
		rootDir = oldRootDir + "_replay"
		if tendermintOS.FileExists(rootDir) {
			tendermintOS.Exit(fmt.Sprintf("temporary copy dir %v already exists", rootDir))
		}
		if err := cpm.Copy(oldRootDir, rootDir); err != nil {
			return err
		}
	}

	configDir := filepath.Join(rootDir, "config")
	dataDir := filepath.Join(rootDir, "data")
	ctx := server.NewDefaultContext()

	appDB, err := sdk.NewLevelDB("application", dataDir)
	if err != nil {
		return err
	}

	tmDB, err := sdk.NewLevelDB("state", dataDir)
	if err != nil {
		return err
	}

	bcDB, err := sdk.NewLevelDB("blockstore", dataDir)
	if err != nil {
		return err
	}

	// TraceStore
	var traceStoreWriter io.Writer
	var traceStoreDir = filepath.Join(dataDir, "trace.log")
	traceStoreWriter, err = os.OpenFile(
		traceStoreDir,
		os.O_WRONLY|os.O_APPEND|os.O_CREATE,
		0666,
	)
	if err != nil {
		return err
	}

	myapp := application.NewApplication(
		ctx.Logger,
		appDB,
		traceStoreWriter,
		true,
		uint(1),
		map[int64]bool{},
		viper.GetString(flags.FlagHome),
	)

	// Genesis
	var genDocPath = filepath.Join(configDir, "pCore_local_genesis.json")
	genDoc, err := tm.GenesisDocFromFile(genDocPath)
	if err != nil {
		return err
	}
	genState, err := tmsm.MakeGenesisState(genDoc)
	if err != nil {
		return err
	}
	// tmsm.SaveState(tmDB, genState)

	cc := proxy.NewLocalClientCreator(myapp)
	proxyApp := proxy.NewAppConns(cc)
	err = proxyApp.Start()
	if err != nil {
		return err
	}
	defer func() {
		_ = proxyApp.Stop()
	}()

	state := tmsm.LoadState(tmDB)
	if state.LastBlockHeight == 0 {
		validators := tm.TM2PB.ValidatorUpdates(genState.Validators)
		csParams := tm.TM2PB.ConsensusParams(genDoc.ConsensusParams)
		req := abci.RequestInitChain{
			Time:            genDoc.GenesisTime,
			ChainId:         genDoc.ChainID,
			ConsensusParams: csParams,
			Validators:      validators,
			AppStateBytes:   genDoc.AppState,
		}
		res, err := proxyApp.Consensus().InitChainSync(req)
		if err != nil {
			return err
		}
		newValidatorz, err := tm.PB2TM.ValidatorUpdates(res.Validators)
		if err != nil {
			return err
		}
		newValidators := tm.NewValidatorSet(newValidatorz)

		// Take the genesis state.
		state = genState
		state.Validators = newValidators
		state.NextValidators = newValidators
	}

	blockExec := tmsm.NewBlockExecutor(tmDB, ctx.Logger, proxyApp.Consensus(), nil, tmsm.MockEvidencePool{})

	blockStore := tmstore.NewBlockStore(bcDB)

	tz := []time.Duration{0, 0, 0}
	for i := int(state.LastBlockHeight) + 1; ; i++ {
		t1 := time.Now()

		blockmeta := blockStore.LoadBlockMeta(int64(i))
		if blockmeta == nil {
			return nil
		}
		block := blockStore.LoadBlock(int64(i))
		if block == nil {
			return fmt.Errorf("couldn't find block %d", i)
		}

		t2 := time.Now()

		state, _, err = blockExec.ApplyBlock(state, blockmeta.BlockID, block)
		if err != nil {
			return err
		}

		t3 := time.Now()
		tz[0] += t2.Sub(t1)
		tz[1] += t3.Sub(t2)

	}
}
