// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package base

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/CosmWasm/wasmd/x/wasm"
	wasmKeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	serverTypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/snapshots"
	"github.com/cosmos/cosmos-sdk/store"
	storeTypes "github.com/cosmos/cosmos-sdk/store/types"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authRest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	authKeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authSimulation "github.com/cosmos/cosmos-sdk/x/auth/simulation"
	authTx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	vestingTypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzKeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	authzModule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	capabilityKeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilityTypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisisKeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisisTypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	distributionKeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	evidenceKeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	evidenceTypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	feeGrantKeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	feeGrantModule "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutilTypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govKeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/mint"
	mintKeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	mintTypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsKeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramsTypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramsProposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingKeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingTypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingKeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeKeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradeTypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ica "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts"
	icaControllerTypes "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts/controller/types"
	icaHostKeeper "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts/host/keeper"
	icaHostTypes "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts/host/types"
	icaTypes "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts/types"
	ibcTransfer "github.com/cosmos/ibc-go/v4/modules/apps/transfer"
	ibcTransferKeeper "github.com/cosmos/ibc-go/v4/modules/apps/transfer/keeper"
	ibcTransferTypes "github.com/cosmos/ibc-go/v4/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v4/modules/core"
	ibcClient "github.com/cosmos/ibc-go/v4/modules/core/02-client"
	ibcClientTypes "github.com/cosmos/ibc-go/v4/modules/core/02-client/types"
	ibcHost "github.com/cosmos/ibc-go/v4/modules/core/24-host"
	ibcAnte "github.com/cosmos/ibc-go/v4/modules/core/ante"
	ibcKeeper "github.com/cosmos/ibc-go/v4/modules/core/keeper"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rakyll/statik/fs"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	abciTypes "github.com/tendermint/tendermint/abci/types"
	tendermintJSON "github.com/tendermint/tendermint/libs/json"
	tendermintLog "github.com/tendermint/tendermint/libs/log"
	tendermintOS "github.com/tendermint/tendermint/libs/os"
	protoTendermintTypes "github.com/tendermint/tendermint/proto/tendermint/types"
	tendermintDB "github.com/tendermint/tm-db"

	"github.com/AssetMantle/modules/helpers"
	"github.com/AssetMantle/modules/helpers/base"
	documentIDGetters "github.com/AssetMantle/modules/utilities/rest/id_getters/docs"
	"github.com/AssetMantle/modules/x/assets"
	"github.com/AssetMantle/modules/x/classifications"
	"github.com/AssetMantle/modules/x/classifications/auxiliaries/bond"
	"github.com/AssetMantle/modules/x/classifications/auxiliaries/conform"
	"github.com/AssetMantle/modules/x/classifications/auxiliaries/define"
	"github.com/AssetMantle/modules/x/classifications/auxiliaries/member"
	"github.com/AssetMantle/modules/x/classifications/auxiliaries/unbond"
	"github.com/AssetMantle/modules/x/identities"
	"github.com/AssetMantle/modules/x/identities/auxiliaries/authenticate"
	"github.com/AssetMantle/modules/x/maintainers"
	"github.com/AssetMantle/modules/x/maintainers/auxiliaries/authorize"
	"github.com/AssetMantle/modules/x/maintainers/auxiliaries/deputize"
	"github.com/AssetMantle/modules/x/maintainers/auxiliaries/maintain"
	"github.com/AssetMantle/modules/x/maintainers/auxiliaries/revoke"
	"github.com/AssetMantle/modules/x/maintainers/auxiliaries/super"
	"github.com/AssetMantle/modules/x/metas"
	"github.com/AssetMantle/modules/x/metas/auxiliaries/supplement"
	"github.com/AssetMantle/modules/x/orders"
	"github.com/AssetMantle/modules/x/splits"
	splitsMint "github.com/AssetMantle/modules/x/splits/auxiliaries/mint"
	"github.com/AssetMantle/modules/x/splits/auxiliaries/renumerate"
	"github.com/AssetMantle/modules/x/splits/auxiliaries/transfer"

	"github.com/AssetMantle/node/application/types/applications"
	"github.com/AssetMantle/node/application/types/applications/constants"
	"github.com/AssetMantle/node/utilities/rest"
)

type wasmApplication struct {
	name string

	moduleBasicManager module.BasicManager

	codec helpers.Codec

	enabledWasmProposalTypeList []wasm.ProposalType
	moduleAccountPermissions    map[string][]string
	tokenReceiveAllowedModules  map[string]bool

	keys map[string]*sdkTypes.KVStoreKey

	stakingKeeper      stakingKeeper.Keeper
	slashingKeeper     slashingKeeper.Keeper
	distributionKeeper distributionKeeper.Keeper
	crisisKeeper       crisisKeeper.Keeper

	moduleManager *module.Manager

	baseapp.BaseApp
}

var _ applications.Application = (*wasmApplication)(nil)

func (wasmApplication wasmApplication) GetDefaultNodeHome() string {
	return os.ExpandEnv("$HOME/." + wasmApplication.name)
}
func (wasmApplication wasmApplication) GetDefaultClientHome() string {
	return os.ExpandEnv("$HOME/." + wasmApplication.name)
}
func (wasmApplication wasmApplication) GetModuleBasicManager() module.BasicManager {
	return wasmApplication.moduleBasicManager
}
func (wasmApplication wasmApplication) GetCodec() helpers.Codec {
	return wasmApplication.codec
}
func (wasmApplication wasmApplication) LoadHeight(height int64) error {
	return wasmApplication.BaseApp.LoadVersion(height)
}
func (wasmApplication wasmApplication) ExportApplicationStateAndValidators(forZeroHeight bool, jailWhiteList []string) (serverTypes.ExportedApp, error) {
	context := wasmApplication.BaseApp.NewContext(true, protoTendermintTypes.Header{Height: wasmApplication.BaseApp.LastBlockHeight()})

	height := wasmApplication.LastBlockHeight() + 1
	if forZeroHeight {
		applyWhiteList := false

		if len(jailWhiteList) > 0 {
			applyWhiteList = true
		}

		whiteListMap := make(map[string]bool)

		for _, address := range jailWhiteList {
			if _, err := sdkTypes.ValAddressFromBech32(address); err != nil {
				panic(err)
			}

			whiteListMap[address] = true
		}

		wasmApplication.crisisKeeper.AssertInvariants(context)

		wasmApplication.stakingKeeper.IterateValidators(context, func(_ int64, val stakingTypes.ValidatorI) (stop bool) {
			_, _ = wasmApplication.distributionKeeper.WithdrawValidatorCommission(context, val.GetOperator())
			return false
		})

		delegations := wasmApplication.stakingKeeper.GetAllDelegations(context)
		for _, delegation := range delegations {
			validatorAddress, err := sdkTypes.ValAddressFromBech32(delegation.ValidatorAddress)
			if err != nil {
				panic(err)
			}

			delegatorAddress, err := sdkTypes.AccAddressFromBech32(delegation.DelegatorAddress)
			if err != nil {
				panic(err)
			}
			_, _ = wasmApplication.distributionKeeper.WithdrawDelegationRewards(context, delegatorAddress, validatorAddress)
		}

		wasmApplication.distributionKeeper.DeleteAllValidatorSlashEvents(context)

		wasmApplication.distributionKeeper.DeleteAllValidatorHistoricalRewards(context)

		height := context.BlockHeight()
		context = context.WithBlockHeight(0)

		wasmApplication.stakingKeeper.IterateValidators(context, func(_ int64, val stakingTypes.ValidatorI) (stop bool) {
			scraps := wasmApplication.distributionKeeper.GetValidatorOutstandingRewardsCoins(context, val.GetOperator())
			feePool := wasmApplication.distributionKeeper.GetFeePool(context)
			feePool.CommunityPool = feePool.CommunityPool.Add(scraps...)
			wasmApplication.distributionKeeper.SetFeePool(context, feePool)

			wasmApplication.distributionKeeper.Hooks().AfterValidatorCreated(context, val.GetOperator())
			return false
		})

		for _, delegation := range delegations {
			validatorAddress, err := sdkTypes.ValAddressFromBech32(delegation.ValidatorAddress)
			if err != nil {
				panic(err)
			}

			delegatorAddress, err := sdkTypes.AccAddressFromBech32(delegation.DelegatorAddress)
			if err != nil {
				panic(err)
			}
			wasmApplication.distributionKeeper.Hooks().BeforeDelegationCreated(context, delegatorAddress, validatorAddress)
			wasmApplication.distributionKeeper.Hooks().AfterDelegationModified(context, delegatorAddress, validatorAddress)
		}

		context = context.WithBlockHeight(height)

		wasmApplication.stakingKeeper.IterateRedelegations(context, func(_ int64, redelegation stakingTypes.Redelegation) (stop bool) {
			for i := range redelegation.Entries {
				redelegation.Entries[i].CreationHeight = 0
			}
			wasmApplication.stakingKeeper.SetRedelegation(context, redelegation)
			return false
		})

		wasmApplication.stakingKeeper.IterateUnbondingDelegations(context, func(_ int64, unbondingDelegation stakingTypes.UnbondingDelegation) (stop bool) {
			for i := range unbondingDelegation.Entries {
				unbondingDelegation.Entries[i].CreationHeight = 0
			}
			wasmApplication.stakingKeeper.SetUnbondingDelegation(context, unbondingDelegation)
			return false
		})

		kvStore := context.KVStore(wasmApplication.keys[stakingTypes.StoreKey])
		kvStoreReversePrefixIterator := sdkTypes.KVStoreReversePrefixIterator(kvStore, stakingTypes.ValidatorsKey)
		counter := int16(0)

		for ; kvStoreReversePrefixIterator.Valid(); kvStoreReversePrefixIterator.Next() {
			addr := sdkTypes.ValAddress(stakingTypes.AddressFromValidatorsKey(kvStoreReversePrefixIterator.Key()))
			validator, found := wasmApplication.stakingKeeper.GetValidator(context, addr)

			if !found {
				panic("Validator not found!")
			}

			validator.UnbondingHeight = 0

			if applyWhiteList && !whiteListMap[addr.String()] {
				validator.Jailed = true
			}

			wasmApplication.stakingKeeper.SetValidator(context, validator)
			counter++
		}

		if err := kvStoreReversePrefixIterator.Close(); err != nil {
			log.Fatal(err)
		}

		if _, err := wasmApplication.stakingKeeper.ApplyAndReturnValidatorSetUpdates(context); err != nil {
			log.Fatal(err)
		}

		wasmApplication.slashingKeeper.IterateValidatorSigningInfos(
			context,
			func(validatorConsAddress sdkTypes.ConsAddress, validatorSigningInfo slashingTypes.ValidatorSigningInfo) (stop bool) {
				validatorSigningInfo.StartHeight = 0
				wasmApplication.slashingKeeper.SetValidatorSigningInfo(context, validatorConsAddress, validatorSigningInfo)
				return false
			},
		)
	}

	genesisState := wasmApplication.moduleManager.ExportGenesis(context, wasmApplication.GetCodec())
	applicationState, err := json.MarshalIndent(genesisState, "", "  ")
	if err != nil {
		return serverTypes.ExportedApp{}, err
	}

	validators, err := staking.WriteValidators(context, wasmApplication.stakingKeeper)
	return serverTypes.ExportedApp{
		AppState:        applicationState,
		Validators:      validators,
		Height:          height,
		ConsensusParams: wasmApplication.BaseApp.GetConsensusParams(context),
	}, err
}
func (wasmApplication wasmApplication) RegisterAPIRoutes(server *api.Server, apiConfig config.APIConfig) {
	clientCtx := server.ClientCtx
	rpc.RegisterRoutes(clientCtx, server.Router)
	authRest.RegisterTxRoutes(clientCtx, server.Router)
	authTx.RegisterGRPCGatewayRoutes(clientCtx, server.GRPCGatewayRouter)
	tmservice.RegisterGRPCGatewayRoutes(clientCtx, server.GRPCGatewayRouter)
	documentIDGetters.RegisterRESTRoutes(clientCtx, server.Router)
	wasmApplication.moduleBasicManager.RegisterRESTRoutes(clientCtx, server.Router)
	wasmApplication.moduleBasicManager.RegisterGRPCGatewayRoutes(clientCtx, server.GRPCGatewayRouter)
	rest.RegisterRESTRoutes(clientCtx, server.Router)
	if apiConfig.Swagger {
		Fs, err := fs.New()
		if err != nil {
			panic(err)
		}

		server.Router.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", http.FileServer(Fs)))
	}
}
func (wasmApplication wasmApplication) RegisterTxService(context client.Context) {
	authTx.RegisterTxService(wasmApplication.BaseApp.GRPCQueryRouter(), context, wasmApplication.BaseApp.Simulate, context.InterfaceRegistry)
}
func (wasmApplication wasmApplication) RegisterTendermintService(context client.Context) {
	tmservice.RegisterTendermintService(wasmApplication.BaseApp.GRPCQueryRouter(), context, context.InterfaceRegistry)
}
func (wasmApplication wasmApplication) AppCreator(logger tendermintLog.Logger, db tendermintDB.DB, writer io.Writer, appOptions serverTypes.AppOptions) serverTypes.Application {
	var multiStorePersistentCache sdkTypes.MultiStorePersistentCache

	if cast.ToBool(appOptions.Get(server.FlagInterBlockCache)) {
		multiStorePersistentCache = store.NewCommitKVStoreCacheManager()
	}

	skipUpgradeHeights := make(map[int64]bool)
	for _, h := range cast.ToIntSlice(appOptions.Get(server.FlagUnsafeSkipUpgrades)) {
		skipUpgradeHeights[int64(h)] = true
	}

	pruningOpts, err := server.GetPruningOptionsFromFlags(appOptions)
	if err != nil {
		panic(err)
	}

	snapshotDir := filepath.Join(cast.ToString(appOptions.Get(flags.FlagHome)), "data", "snapshots")
	snapshotDB, err := sdkTypes.NewLevelDB("metadata", snapshotDir)
	if err != nil {
		panic(err)
	}
	snapshotStore, err := snapshots.NewStore(snapshotDB, snapshotDir)
	if err != nil {
		panic(err)
	}

	var wasmOpts []wasm.Option
	if cast.ToBool(appOptions.Get("telemetry.enabled")) {
		wasmOpts = append(wasmOpts, wasmKeeper.WithVMCacheMetrics(prometheus.DefaultRegisterer))
	}

	return wasmApplication.Initialize(
		logger,
		db,
		writer,
		true,
		cast.ToUint(appOptions.Get(server.FlagInvCheckPeriod)),
		skipUpgradeHeights,
		cast.ToString(appOptions.Get(flags.FlagHome)),
		appOptions,
		baseapp.SetPruning(pruningOpts),
		baseapp.SetMinGasPrices(cast.ToString(appOptions.Get(server.FlagMinGasPrices))),
		baseapp.SetHaltHeight(cast.ToUint64(appOptions.Get(server.FlagHaltHeight))),
		baseapp.SetHaltTime(cast.ToUint64(appOptions.Get(server.FlagHaltTime))),
		baseapp.SetMinRetainBlocks(cast.ToUint64(appOptions.Get(server.FlagMinRetainBlocks))),
		baseapp.SetInterBlockCache(multiStorePersistentCache),
		baseapp.SetTrace(cast.ToBool(appOptions.Get(server.FlagTrace))),
		baseapp.SetIndexEvents(cast.ToStringSlice(appOptions.Get(server.FlagIndexEvents))),
		baseapp.SetSnapshotStore(snapshotStore),
		baseapp.SetSnapshotInterval(cast.ToUint64(appOptions.Get(server.FlagStateSyncSnapshotInterval))),
		baseapp.SetSnapshotKeepRecent(cast.ToUint32(appOptions.Get(server.FlagStateSyncSnapshotKeepRecent))))
}
func (wasmApplication wasmApplication) AppExporter(logger tendermintLog.Logger, db tendermintDB.DB, writer io.Writer, height int64, forZeroHeight bool, jailAllowedAddrs []string, appOptions serverTypes.AppOptions) (serverTypes.ExportedApp, error) {
	home, ok := appOptions.Get(flags.FlagHome).(string)
	if !ok || home == "" {
		return serverTypes.ExportedApp{}, errors.New("application home is not set")
	}

	var loadLatest bool
	if height == -1 {
		loadLatest = true
	}

	Application := wasmApplication.Initialize(
		logger,
		db,
		writer,
		loadLatest,
		cast.ToUint(appOptions.Get(server.FlagInvCheckPeriod)),
		map[int64]bool{},
		home,
		appOptions,
	)

	if height != -1 {
		if err := wasmApplication.LoadHeight(height); err != nil {
			return serverTypes.ExportedApp{}, err
		}
	}

	return Application.ExportApplicationStateAndValidators(forZeroHeight, jailAllowedAddrs)
}
func (wasmApplication wasmApplication) ModuleInitFlags(command *cobra.Command) {
	crisis.AddModuleInitFlags(command)
}
func (wasmApplication wasmApplication) Initialize(logger tendermintLog.Logger, db tendermintDB.DB, writer io.Writer, loadLatest bool, invCheckPeriod uint, skipUpgradeHeights map[int64]bool, home string, appOptions serverTypes.AppOptions, baseAppOptions ...func(*baseapp.BaseApp)) applications.Application {
	wasmApplication.BaseApp = *baseapp.NewBaseApp(wasmApplication.name, logger, db, wasmApplication.GetCodec().TxDecoder(), baseAppOptions...)
	wasmApplication.BaseApp.SetCommitMultiStoreTracer(writer)
	wasmApplication.BaseApp.SetVersion(version.Version)
	wasmApplication.BaseApp.SetInterfaceRegistry(wasmApplication.GetCodec().InterfaceRegistry())

	wasmApplication.keys = sdkTypes.NewKVStoreKeys(
		authTypes.StoreKey,
		bankTypes.StoreKey,
		stakingTypes.StoreKey,
		mintTypes.StoreKey,
		distributionTypes.StoreKey,
		slashingTypes.StoreKey,
		govTypes.StoreKey,
		paramsTypes.StoreKey,
		ibcHost.StoreKey,
		upgradeTypes.StoreKey,
		evidenceTypes.StoreKey,
		ibcTransferTypes.StoreKey,
		capabilityTypes.StoreKey,
		feegrant.StoreKey,
		authzKeeper.StoreKey,
		icaHostTypes.StoreKey,

		wasm.StoreKey,

		assets.Prototype().Name(),
		classifications.Prototype().Name(),
		identities.Prototype().Name(),
		maintainers.Prototype().Name(),
		metas.Prototype().Name(),
		orders.Prototype().Name(),
		splits.Prototype().Name(),
	)

	transientStoreKeys := sdkTypes.NewTransientStoreKeys(paramsTypes.TStoreKey)
	memoryStoreKeys := sdkTypes.NewMemoryStoreKeys(capabilityTypes.MemStoreKey)

	ParamsKeeper := paramsKeeper.NewKeeper(
		wasmApplication.GetCodec(),
		wasmApplication.GetCodec().GetLegacyAmino(),
		wasmApplication.keys[paramsTypes.StoreKey],
		transientStoreKeys[paramsTypes.TStoreKey],
	)

	wasmApplication.BaseApp.SetParamStore(ParamsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramsKeeper.ConsensusParamsKeyTable()))

	CapabilityKeeper := capabilityKeeper.NewKeeper(wasmApplication.GetCodec(), wasmApplication.keys[capabilityTypes.StoreKey], memoryStoreKeys[capabilityTypes.MemStoreKey])
	scopedIBCKeeper := CapabilityKeeper.ScopeToModule(ibcHost.ModuleName)
	scopedTransferKeeper := CapabilityKeeper.ScopeToModule(ibcTransferTypes.ModuleName)
	scopedICAHostKeeper := CapabilityKeeper.ScopeToModule(icaHostTypes.SubModuleName)
	CapabilityKeeper.Seal()

	AccountKeeper := authKeeper.NewAccountKeeper(
		wasmApplication.GetCodec(),
		wasmApplication.keys[authTypes.StoreKey],
		ParamsKeeper.Subspace(authTypes.ModuleName),
		authTypes.ProtoBaseAccount,
		wasmApplication.moduleAccountPermissions,
	)

	blacklistedAddresses := make(map[string]bool)
	for account := range wasmApplication.moduleAccountPermissions {
		blacklistedAddresses[authTypes.NewModuleAddress(account).String()] = !wasmApplication.tokenReceiveAllowedModules[account]
	}

	BankKeeper := bankKeeper.NewBaseKeeper(
		wasmApplication.GetCodec(),
		wasmApplication.keys[bankTypes.StoreKey],
		AccountKeeper,
		ParamsKeeper.Subspace(bankTypes.ModuleName),
		blacklistedAddresses,
	)

	AuthzKeeper := authzKeeper.NewKeeper(
		wasmApplication.keys[authzKeeper.StoreKey],
		wasmApplication.GetCodec(),
		wasmApplication.BaseApp.MsgServiceRouter(),
	)

	FeeGrantKeeper := feeGrantKeeper.NewKeeper(
		wasmApplication.GetCodec(),
		wasmApplication.keys[feegrant.StoreKey],
		AccountKeeper,
	)

	wasmApplication.stakingKeeper = stakingKeeper.NewKeeper(
		wasmApplication.GetCodec(),
		wasmApplication.keys[stakingTypes.StoreKey],
		AccountKeeper,
		BankKeeper,
		ParamsKeeper.Subspace(stakingTypes.ModuleName),
	)

	MintKeeper := mintKeeper.NewKeeper(
		wasmApplication.GetCodec(),
		wasmApplication.keys[mintTypes.StoreKey],
		ParamsKeeper.Subspace(mintTypes.ModuleName),
		&wasmApplication.stakingKeeper,
		AccountKeeper,
		BankKeeper,
		authTypes.FeeCollectorName,
	)

	wasmApplication.distributionKeeper = distributionKeeper.NewKeeper(
		wasmApplication.GetCodec(),
		wasmApplication.keys[distributionTypes.StoreKey],
		ParamsKeeper.Subspace(distributionTypes.ModuleName),
		AccountKeeper,
		BankKeeper,
		&wasmApplication.stakingKeeper,
		authTypes.FeeCollectorName,
		blacklistedAddresses,
	)

	wasmApplication.slashingKeeper = slashingKeeper.NewKeeper(
		wasmApplication.GetCodec(),
		wasmApplication.keys[slashingTypes.StoreKey],
		&wasmApplication.stakingKeeper,
		ParamsKeeper.Subspace(slashingTypes.ModuleName),
	)

	wasmApplication.crisisKeeper = crisisKeeper.NewKeeper(
		ParamsKeeper.Subspace(crisisTypes.ModuleName),
		invCheckPeriod,
		BankKeeper,
		authTypes.FeeCollectorName,
	)

	UpgradeKeeper := upgradeKeeper.NewKeeper(
		skipUpgradeHeights,
		wasmApplication.keys[upgradeTypes.StoreKey],
		wasmApplication.GetCodec(),
		home,
		&wasmApplication.BaseApp,
	)

	wasmApplication.stakingKeeper = *wasmApplication.stakingKeeper.SetHooks(stakingTypes.NewMultiStakingHooks(wasmApplication.distributionKeeper.Hooks(), wasmApplication.slashingKeeper.Hooks()))

	IBCKeeper := ibcKeeper.NewKeeper(
		wasmApplication.GetCodec(),
		wasmApplication.keys[ibcHost.StoreKey],
		ParamsKeeper.Subspace(ibcHost.ModuleName),
		wasmApplication.stakingKeeper,
		UpgradeKeeper,
		scopedIBCKeeper,
	)

	govRouter := govTypes.NewRouter().
		AddRoute(govTypes.RouterKey, govTypes.ProposalHandler).
		AddRoute(paramsProposal.RouterKey, params.NewParamChangeProposalHandler(ParamsKeeper)).
		AddRoute(distributionTypes.RouterKey, distribution.NewCommunityPoolSpendProposalHandler(wasmApplication.distributionKeeper)).
		AddRoute(upgradeTypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(UpgradeKeeper)).
		AddRoute(ibcClientTypes.RouterKey, ibcClient.NewClientProposalHandler(IBCKeeper.ClientKeeper))

	GovKeeper := govKeeper.NewKeeper(
		wasmApplication.GetCodec(),
		wasmApplication.keys[govTypes.StoreKey],
		ParamsKeeper.Subspace(govTypes.ModuleName).WithKeyTable(govTypes.ParamKeyTable()),
		AccountKeeper,
		BankKeeper,
		&wasmApplication.stakingKeeper,
		govRouter,
	)

	IBCTransferKeeper := ibcTransferKeeper.NewKeeper(
		wasmApplication.GetCodec(),
		wasmApplication.keys[ibcTransferTypes.StoreKey],
		ParamsKeeper.Subspace(ibcTransferTypes.ModuleName),
		IBCKeeper.ChannelKeeper,
		IBCKeeper.ChannelKeeper,
		&IBCKeeper.PortKeeper,
		AccountKeeper,
		BankKeeper,
		scopedTransferKeeper,
	)

	ICAHostKeeper := icaHostKeeper.NewKeeper(
		wasmApplication.GetCodec(),
		wasmApplication.keys[icaHostTypes.StoreKey],
		ParamsKeeper.Subspace(icaHostTypes.SubModuleName),
		IBCKeeper.ChannelKeeper,
		&IBCKeeper.PortKeeper,
		AccountKeeper,
		scopedICAHostKeeper,
		wasmApplication.MsgServiceRouter(),
	)

	EvidenceKeeper := *evidenceKeeper.NewKeeper(
		wasmApplication.GetCodec(),
		wasmApplication.keys[evidenceTypes.StoreKey],
		&wasmApplication.stakingKeeper,
		wasmApplication.slashingKeeper,
	)

	// ******************************
	// WASM Module Boilerplate Code
	// ******************************

	metasModule := metas.Prototype().Initialize(
		wasmApplication.keys[metas.Prototype().Name()],
		ParamsKeeper.Subspace(metas.Prototype().Name()),
	)

	classificationsModule := classifications.Prototype().Initialize(
		wasmApplication.keys[classifications.Prototype().Name()],
		ParamsKeeper.Subspace(classifications.Prototype().Name()),
		BankKeeper,
		wasmApplication.stakingKeeper,
	)

	maintainersModule := maintainers.Prototype().Initialize(
		wasmApplication.keys[metas.Prototype().Name()],
		ParamsKeeper.Subspace(maintainers.Prototype().Name()),
		classificationsModule.GetAuxiliary(member.Auxiliary.GetName()),
	)
	identitiesModule := identities.Prototype().Initialize(
		wasmApplication.keys[identities.Prototype().Name()],
		ParamsKeeper.Subspace(identities.Prototype().Name()),
		classificationsModule.GetAuxiliary(bond.Auxiliary.GetName()),
		classificationsModule.GetAuxiliary(conform.Auxiliary.GetName()),
		classificationsModule.GetAuxiliary(define.Auxiliary.GetName()),
		classificationsModule.GetAuxiliary(member.Auxiliary.GetName()),
		classificationsModule.GetAuxiliary(unbond.Auxiliary.GetName()),
		maintainersModule.GetAuxiliary(deputize.Auxiliary.GetName()),
		maintainersModule.GetAuxiliary(maintain.Auxiliary.GetName()),
		maintainersModule.GetAuxiliary(revoke.Auxiliary.GetName()),
		maintainersModule.GetAuxiliary(super.Auxiliary.GetName()),
		maintainersModule.GetAuxiliary(authorize.Auxiliary.GetName()),
		metasModule.GetAuxiliary(supplement.Auxiliary.GetName()),
	)
	splitsModule := splits.Prototype().Initialize(
		wasmApplication.keys[splits.Prototype().Name()],
		ParamsKeeper.Subspace(splits.Prototype().Name()),
		BankKeeper,
		identitiesModule.GetAuxiliary(authenticate.Auxiliary.GetName()),
	)
	assetsModule := assets.Prototype().Initialize(
		wasmApplication.keys[assets.Prototype().Name()],
		ParamsKeeper.Subspace(assets.Prototype().Name()),
		identitiesModule.GetAuxiliary(authenticate.Auxiliary.GetName()),
		classificationsModule.GetAuxiliary(conform.Auxiliary.GetName()),
		classificationsModule.GetAuxiliary(define.Auxiliary.GetName()),
		classificationsModule.GetAuxiliary(bond.Auxiliary.GetName()),
		classificationsModule.GetAuxiliary(unbond.Auxiliary.GetName()),
		maintainersModule.GetAuxiliary(deputize.Auxiliary.GetName()),
		maintainersModule.GetAuxiliary(maintain.Auxiliary.GetName()),
		splitsModule.GetAuxiliary(renumerate.Auxiliary.GetName()),
		maintainersModule.GetAuxiliary(revoke.Auxiliary.GetName()),
		maintainersModule.GetAuxiliary(super.Auxiliary.GetName()),
		metasModule.GetAuxiliary(supplement.Auxiliary.GetName()),
		splitsModule.GetAuxiliary(splitsMint.Auxiliary.GetName()),
		maintainersModule.GetAuxiliary(authorize.Auxiliary.GetName()),
	)
	ordersModule := orders.Prototype().Initialize(
		wasmApplication.keys[orders.Prototype().Name()],
		ParamsKeeper.Subspace(orders.Prototype().Name()),
		identitiesModule.GetAuxiliary(authenticate.Auxiliary.GetName()),
		classificationsModule.GetAuxiliary(conform.Auxiliary.GetName()),
		classificationsModule.GetAuxiliary(define.Auxiliary.GetName()),
		maintainersModule.GetAuxiliary(deputize.Auxiliary.GetName()),
		maintainersModule.GetAuxiliary(maintain.Auxiliary.GetName()),
		maintainersModule.GetAuxiliary(revoke.Auxiliary.GetName()),
		maintainersModule.GetAuxiliary(super.Auxiliary.GetName()),
		metasModule.GetAuxiliary(supplement.Auxiliary.GetName()),
		splitsModule.GetAuxiliary(transfer.Auxiliary.GetName()),
		maintainersModule.GetAuxiliary(authorize.Auxiliary.GetName()),
	)

	wasmApplication.moduleManager = module.NewManager(
		genutil.NewAppModule(AccountKeeper, wasmApplication.stakingKeeper, wasmApplication.BaseApp.DeliverTx, wasmApplication.GetCodec()),
		auth.NewAppModule(wasmApplication.GetCodec(), AccountKeeper, nil),
		vesting.NewAppModule(AccountKeeper, BankKeeper),
		bank.NewAppModule(wasmApplication.GetCodec(), BankKeeper, AccountKeeper),
		capability.NewAppModule(wasmApplication.GetCodec(), *CapabilityKeeper),
		crisis.NewAppModule(&wasmApplication.crisisKeeper, cast.ToBool(appOptions.Get(crisis.FlagSkipGenesisInvariants))),
		gov.NewAppModule(wasmApplication.GetCodec(), GovKeeper, AccountKeeper, BankKeeper),
		mint.NewAppModule(wasmApplication.GetCodec(), MintKeeper, AccountKeeper),
		slashing.NewAppModule(wasmApplication.GetCodec(), wasmApplication.slashingKeeper, AccountKeeper, BankKeeper, wasmApplication.stakingKeeper),
		distribution.NewAppModule(wasmApplication.GetCodec(), wasmApplication.distributionKeeper, AccountKeeper, BankKeeper, wasmApplication.stakingKeeper),
		staking.NewAppModule(wasmApplication.GetCodec(), wasmApplication.stakingKeeper, AccountKeeper, BankKeeper),
		upgrade.NewAppModule(UpgradeKeeper),
		evidence.NewAppModule(EvidenceKeeper),
		feeGrantModule.NewAppModule(wasmApplication.GetCodec(), AccountKeeper, BankKeeper, FeeGrantKeeper, wasmApplication.GetCodec().InterfaceRegistry()),
		authzModule.NewAppModule(wasmApplication.GetCodec(), AuthzKeeper, AccountKeeper, BankKeeper, wasmApplication.GetCodec().InterfaceRegistry()),
		ibc.NewAppModule(IBCKeeper),
		params.NewAppModule(ParamsKeeper),
		ibcTransfer.NewAppModule(IBCTransferKeeper),
		ica.NewAppModule(nil, &ICAHostKeeper),

		assetsModule,
		classificationsModule,
		identitiesModule,
		maintainersModule,
		metasModule,
		ordersModule,
		splitsModule,
	)

	wasmApplication.moduleManager.SetOrderBeginBlockers(
		upgradeTypes.ModuleName,
		capabilityTypes.ModuleName,
		crisisTypes.ModuleName,
		govTypes.ModuleName,
		stakingTypes.ModuleName,
		ibcTransferTypes.ModuleName,
		ibcHost.ModuleName,
		icaTypes.ModuleName,
		authTypes.ModuleName,
		bankTypes.ModuleName,
		distributionTypes.ModuleName,
		slashingTypes.ModuleName,
		mintTypes.ModuleName,
		genutilTypes.ModuleName,
		evidenceTypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramsTypes.ModuleName,
		vestingTypes.ModuleName,

		assets.Prototype().Name(),
		classifications.Prototype().Name(),
		identities.Prototype().Name(),
		maintainers.Prototype().Name(),
		metas.Prototype().Name(),
		orders.Prototype().Name(),
		splits.Prototype().Name(),
	)
	wasmApplication.moduleManager.SetOrderEndBlockers(
		crisisTypes.ModuleName,
		govTypes.ModuleName,
		stakingTypes.ModuleName,
		ibcTransferTypes.ModuleName,
		ibcHost.ModuleName,
		icaTypes.ModuleName,
		feegrant.ModuleName,
		authz.ModuleName,
		capabilityTypes.ModuleName,
		authTypes.ModuleName,
		bankTypes.ModuleName,
		distributionTypes.ModuleName,
		slashingTypes.ModuleName,
		mintTypes.ModuleName,
		genutilTypes.ModuleName,
		evidenceTypes.ModuleName,
		paramsTypes.ModuleName,
		upgradeTypes.ModuleName,
		vestingTypes.ModuleName,

		assets.Prototype().Name(),
		classifications.Prototype().Name(),
		identities.Prototype().Name(),
		maintainers.Prototype().Name(),
		metas.Prototype().Name(),
		orders.Prototype().Name(),
		splits.Prototype().Name(),
	)
	wasmApplication.moduleManager.SetOrderInitGenesis(
		capabilityTypes.ModuleName,
		bankTypes.ModuleName,
		distributionTypes.ModuleName,
		stakingTypes.ModuleName,
		slashingTypes.ModuleName,
		govTypes.ModuleName,
		mintTypes.ModuleName,
		crisisTypes.ModuleName,
		ibcTransferTypes.ModuleName,
		ibcHost.ModuleName,
		icaTypes.ModuleName,
		evidenceTypes.ModuleName,
		feegrant.ModuleName,
		authz.ModuleName,
		authTypes.ModuleName,
		genutilTypes.ModuleName,
		paramsTypes.ModuleName,
		upgradeTypes.ModuleName,
		vestingTypes.ModuleName,

		assets.Prototype().Name(),
		classifications.Prototype().Name(),
		identities.Prototype().Name(),
		maintainers.Prototype().Name(),
		metas.Prototype().Name(),
		orders.Prototype().Name(),
		splits.Prototype().Name(),
	)
	wasmApplication.moduleManager.RegisterInvariants(&wasmApplication.crisisKeeper)
	wasmApplication.moduleManager.RegisterRoutes(wasmApplication.BaseApp.Router(), wasmApplication.BaseApp.QueryRouter(), wasmApplication.GetCodec().GetLegacyAmino())

	configurator := module.NewConfigurator(wasmApplication.GetCodec(), wasmApplication.MsgServiceRouter(), wasmApplication.GRPCQueryRouter())
	wasmApplication.moduleManager.RegisterServices(configurator)

	module.NewSimulationManager(
		auth.NewAppModule(wasmApplication.GetCodec(), AccountKeeper, authSimulation.RandomGenesisAccounts),
		bank.NewAppModule(wasmApplication.GetCodec(), BankKeeper, AccountKeeper),
		capability.NewAppModule(wasmApplication.GetCodec(), *CapabilityKeeper),
		feeGrantModule.NewAppModule(wasmApplication.GetCodec(), AccountKeeper, BankKeeper, FeeGrantKeeper, wasmApplication.GetCodec().InterfaceRegistry()),
		authzModule.NewAppModule(wasmApplication.GetCodec(), AuthzKeeper, AccountKeeper, BankKeeper, wasmApplication.GetCodec().InterfaceRegistry()),
		gov.NewAppModule(wasmApplication.GetCodec(), GovKeeper, AccountKeeper, BankKeeper),
		mint.NewAppModule(wasmApplication.GetCodec(), MintKeeper, AccountKeeper),
		staking.NewAppModule(wasmApplication.GetCodec(), wasmApplication.stakingKeeper, AccountKeeper, BankKeeper),
		distribution.NewAppModule(wasmApplication.GetCodec(), wasmApplication.distributionKeeper, AccountKeeper, BankKeeper, wasmApplication.stakingKeeper),
		slashing.NewAppModule(wasmApplication.GetCodec(), wasmApplication.slashingKeeper, AccountKeeper, BankKeeper, wasmApplication.stakingKeeper),
		params.NewAppModule(ParamsKeeper),
		evidence.NewAppModule(EvidenceKeeper),
		ibc.NewAppModule(IBCKeeper),
		ibcTransfer.NewAppModule(IBCTransferKeeper),

		assets.Prototype(),
		classifications.Prototype(),
		identities.Prototype(),
		maintainers.Prototype(),
		metas.Prototype(),
		orders.Prototype(),
		splits.Prototype(),
	).RegisterStoreDecoders()

	wasmApplication.BaseApp.MountKVStores(wasmApplication.keys)
	wasmApplication.BaseApp.MountTransientStores(transientStoreKeys)
	wasmApplication.BaseApp.MountMemoryStores(memoryStoreKeys)
	wasmApplication.BaseApp.SetAnteHandler(sdkTypes.ChainAnteDecorators(
		ante.NewSetUpContextDecorator(),
		ante.NewRejectExtensionOptionsDecorator(),
		ante.NewValidateBasicDecorator(),
		ante.NewTxTimeoutHeightDecorator(),
		ante.NewValidateMemoDecorator(AccountKeeper),
		ante.NewConsumeGasForTxSizeDecorator(AccountKeeper),
		ante.NewDeductFeeDecorator(AccountKeeper, BankKeeper, FeeGrantKeeper),
		ante.NewSetPubKeyDecorator(AccountKeeper),
		ante.NewValidateSigCountDecorator(AccountKeeper),
		ante.NewSigGasConsumeDecorator(AccountKeeper, ante.DefaultSigVerificationGasConsumer),
		ante.NewSigVerificationDecorator(AccountKeeper, wasmApplication.GetCodec().SignModeHandler()),
		ante.NewIncrementSequenceDecorator(AccountKeeper),
		ibcAnte.NewAnteDecorator(IBCKeeper),
	))
	wasmApplication.BaseApp.SetBeginBlocker(wasmApplication.moduleManager.BeginBlock)
	wasmApplication.BaseApp.SetEndBlocker(wasmApplication.moduleManager.EndBlock)
	wasmApplication.BaseApp.SetInitChainer(func(context sdkTypes.Context, requestInitChain abciTypes.RequestInitChain) abciTypes.ResponseInitChain {
		var genesisState map[string]json.RawMessage
		if err := tendermintJSON.Unmarshal(requestInitChain.AppStateBytes, &genesisState); err != nil {
			panic(err)
		}

		UpgradeKeeper.SetModuleVersionMap(context, wasmApplication.moduleManager.GetVersionMap())

		return wasmApplication.moduleManager.InitGenesis(context, wasmApplication.GetCodec(), genesisState)
	})

	UpgradeKeeper.SetUpgradeHandler(
		constants.UpgradeName,
		func(ctx sdkTypes.Context, _ upgradeTypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			fromVM[icaTypes.ModuleName] = ica.NewAppModule(nil, &ICAHostKeeper).ConsensusVersion()
			controllerParams := icaControllerTypes.Params{}
			hostParams := icaHostTypes.Params{
				HostEnabled: true,
				AllowMessages: []string{
					constants.AuthzMsgExec,
					constants.AuthzMsgGrant,
					constants.AuthzMsgRevoke,
					constants.BankMsgSend,
					constants.BankMsgMultiSend,
					constants.DistrMsgSetWithdrawAddr,
					constants.DistrMsgWithdrawValidatorCommission,
					constants.DistrMsgFundCommunityPool,
					constants.DistrMsgWithdrawDelegatorReward,
					constants.FeegrantMsgGrantAllowance,
					constants.FeegrantMsgRevokeAllowance,
					constants.GovMsgVoteWeighted,
					constants.GovMsgSubmitProposal,
					constants.GovMsgDeposit,
					constants.GovMsgVote,
					constants.StakingMsgEditValidator,
					constants.StakingMsgDelegate,
					constants.StakingMsgUndelegate,
					constants.StakingMsgBeginRedelegate,
					constants.StakingMsgCreateValidator,
					constants.VestingMsgCreateVestingAccount,
					constants.TransferMsgTransfer,
					constants.LiquidityMsgCreatePool,
					constants.LiquidityMsgSwapWithinBatch,
					constants.LiquidityMsgDepositWithinBatch,
					constants.LiquidityMsgWithdrawWithinBatch,
				},
			}
			ica.NewAppModule(nil, &ICAHostKeeper).InitModule(ctx, controllerParams, hostParams)

			return wasmApplication.moduleManager.RunMigrations(ctx, configurator, fromVM)
		},
	)

	upgradeInfo, err := UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}

	if upgradeInfo.Name == constants.UpgradeName && !UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := storeTypes.StoreUpgrades{
			Added: []string{icaHostTypes.StoreKey},
		}

		wasmApplication.SetStoreLoader(upgradeTypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}

	if loadLatest {
		err := wasmApplication.BaseApp.LoadLatestVersion()
		if err != nil {
			tendermintOS.Exit(err.Error())
		}
	}

	return &wasmApplication
}

func NewWasmApplication(name string, moduleBasicManager module.BasicManager, enabledWasmProposalTypeList []wasm.ProposalType, moduleAccountPermissions map[string][]string, tokenReceiveAllowedModules map[string]bool) applications.Application {
	return &wasmApplication{
		name:                        name,
		moduleBasicManager:          moduleBasicManager,
		codec:                       base.CodecPrototype().Initialize(moduleBasicManager),
		enabledWasmProposalTypeList: enabledWasmProposalTypeList,
		moduleAccountPermissions:    moduleAccountPermissions,
		tokenReceiveAllowedModules:  tokenReceiveAllowedModules,
	}
}
