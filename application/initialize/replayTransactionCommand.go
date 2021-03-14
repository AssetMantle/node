/*
 Copyright [2019] - [2020], PERSISTENCE TECHNOLOGIES PTE. LTD. and the assetMantle contributors
 SPDX-License-Identifier: Apache-2.0
*/

package initialize

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/otiai10/copy"
	"github.com/persistenceOne/assetMantle/application"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tendermintABCITypes "github.com/tendermint/tendermint/abci/types"
	tendermintOS "github.com/tendermint/tendermint/libs/os"
	"github.com/tendermint/tendermint/proxy"
	tendermintState "github.com/tendermint/tendermint/state"
	tendermintStore "github.com/tendermint/tendermint/store"
	tendermintTypes "github.com/tendermint/tendermint/types"
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

		if err := copy.Copy(oldRootDir, rootDir); err != nil {
			return err
		}
	}

	configDir := filepath.Join(rootDir, "config")
	dataDir := filepath.Join(rootDir, "data")
	ctx := server.NewDefaultContext()

	appDB, err := sdkTypes.NewLevelDB("application", dataDir)
	if err != nil {
		return err
	}

	tmDB, err := sdkTypes.NewLevelDB("state", dataDir)
	if err != nil {
		return err
	}

	bcDB, err := sdkTypes.NewLevelDB("blockstore", dataDir)
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

	myapp := application.Prototype.Initialize(
		ctx.Logger,
		appDB,
		traceStoreWriter,
		true,
		uint(1),
		map[int64]bool{},
		viper.GetString(flags.FlagHome),
	)

	// Genesis
	var genDocPath = filepath.Join(configDir, "genesis.json")
	genDoc, err := tendermintTypes.GenesisDocFromFile(genDocPath)

	if err != nil {
		return err
	}

	genState, err := tendermintState.MakeGenesisState(genDoc)
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

	state := tendermintState.LoadState(tmDB)
	if state.LastBlockHeight == 0 {
		validators := tendermintTypes.TM2PB.ValidatorUpdates(genState.Validators)
		csParams := tendermintTypes.TM2PB.ConsensusParams(genDoc.ConsensusParams)
		req := tendermintABCITypes.RequestInitChain{
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

		newValidatorList, err := tendermintTypes.PB2TM.ValidatorUpdates(res.Validators)
		if err != nil {
			return err
		}

		validatorSet := tendermintTypes.NewValidatorSet(newValidatorList)

		// Take the genesis state.
		state = genState
		state.Validators = validatorSet
		state.NextValidators = validatorSet
	}

	blockExec := tendermintState.NewBlockExecutor(tmDB, ctx.Logger, proxyApp.Consensus(), nil, tendermintState.MockEvidencePool{})

	blockStore := tendermintStore.NewBlockStore(bcDB)

	tz := []time.Duration{0, 0, 0}

	for i := int(state.LastBlockHeight) + 1; ; i++ {
		t1 := time.Now()

		blockMeta := blockStore.LoadBlockMeta(int64(i))
		if blockMeta == nil {
			return nil
		}

		block := blockStore.LoadBlock(int64(i))
		if block == nil {
			return fmt.Errorf("couldn't find block %d", i)
		}

		t2 := time.Now()

		state, _, err = blockExec.ApplyBlock(state, blockMeta.BlockID, block)
		if err != nil {
			return err
		}

		t3 := time.Now()
		tz[0] += t2.Sub(t1)
		tz[1] += t3.Sub(t2)
	}
}
