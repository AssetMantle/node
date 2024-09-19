// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package base

import (
	"encoding/json"
	"errors"
	rateLimitKeeper "github.com/Stride-Labs/ibc-rate-limiting/ratelimit/keeper"
	rateLimitTypes "github.com/Stride-Labs/ibc-rate-limiting/ratelimit/types"
	cometbftTypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/client/grpc/node"
	snapshotsTypes "github.com/cosmos/cosmos-sdk/snapshots/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	consensusParamKeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	consensusParamTypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	ibcFeeKeeper "github.com/cosmos/ibc-go/v7/modules/apps/29-fee/keeper"
	ibcFeeTypes "github.com/cosmos/ibc-go/v7/modules/apps/29-fee/types"
	ibcExported "github.com/cosmos/ibc-go/v7/modules/core/exported"

	"github.com/cometbft/cometbft/libs/log"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/AssetMantle/modules/helpers"
	"github.com/AssetMantle/modules/helpers/base"
	documentIDGetters "github.com/AssetMantle/modules/utilities/rest/id_getters/docs"
	"github.com/AssetMantle/modules/x/assets"
	"github.com/AssetMantle/modules/x/classifications"
	"github.com/AssetMantle/modules/x/classifications/auxiliaries/bond"
	"github.com/AssetMantle/modules/x/classifications/auxiliaries/burn"
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
	burnSplits "github.com/AssetMantle/modules/x/splits/auxiliaries/burn"
	splitsMint "github.com/AssetMantle/modules/x/splits/auxiliaries/mint"
	"github.com/AssetMantle/modules/x/splits/auxiliaries/purge"
	"github.com/AssetMantle/modules/x/splits/auxiliaries/renumerate"
	"github.com/AssetMantle/modules/x/splits/auxiliaries/transfer"
	tendermintDB "github.com/cometbft/cometbft-db"
	abciTypes "github.com/cometbft/cometbft/abci/types"
	tendermintJSON "github.com/cometbft/cometbft/libs/json"
	tendermintLog "github.com/cometbft/cometbft/libs/log"
	tendermintOS "github.com/cometbft/cometbft/libs/os"
	protoTendermintTypes "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
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
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authKeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
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
	govTypesV1Beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
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
	router "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v7/packetforward"
	routerKeeper "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v7/packetforward/keeper"
	routerTypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v7/packetforward/types"
	ica "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts"
	icaHost "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host"
	icaHostKeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host/keeper"
	icaHostTypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host/types"
	icaTypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	ibcTransfer "github.com/cosmos/ibc-go/v7/modules/apps/transfer"
	ibcTransferKeeper "github.com/cosmos/ibc-go/v7/modules/apps/transfer/keeper"
	ibcTransferTypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v7/modules/core"
	ibcClient "github.com/cosmos/ibc-go/v7/modules/core/02-client"
	ibcClientTypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	ibcPortTypes "github.com/cosmos/ibc-go/v7/modules/core/05-port/types"
	ibcKeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"
	"github.com/rakyll/statik/fs"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/AssetMantle/node/application/types/applications"
	"github.com/AssetMantle/node/application/types/applications/constants"
	"github.com/AssetMantle/node/utilities/rest"
)

type application struct {
	name string

	moduleManager helpers.ModuleManager

	moduleAccountPermissions   map[string][]string
	tokenReceiveAllowedModules map[string]bool

	keys map[string]*storeTypes.KVStoreKey

	stakingKeeper      *stakingKeeper.Keeper
	crisisKeeper       *crisisKeeper.Keeper
	slashingKeeper     slashingKeeper.Keeper
	distributionKeeper distributionKeeper.Keeper

	*baseapp.BaseApp
}

var _ applications.Application = (*application)(nil)

func (application application) RegisterNodeService(context client.Context) {
	node.RegisterNodeService(context, application.GRPCQueryRouter())
}
func (application application) GetDefaultNodeHome() string {
	return os.ExpandEnv("$HOME/." + application.name)
}
func (application application) GetDefaultClientHome() string {
	return os.ExpandEnv("$HOME/." + application.name)
}
func (application application) GetModuleManager() helpers.ModuleManager {
	return application.moduleManager
}
func (application application) GetCodec() helpers.Codec {
	return base.CodecPrototype().Initialize(application.GetModuleManager())
}
func (application application) LoadHeight(height int64) error {
	return application.LoadVersion(height)
}
func (application application) ExportApplicationStateAndValidators(forZeroHeight bool, jailWhiteList []string, modulesToExport []string) (serverTypes.ExportedApp, error) {
	context := application.NewContext(true, protoTendermintTypes.Header{Height: application.LastBlockHeight()})

	height := application.LastBlockHeight() + 1
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

		application.crisisKeeper.AssertInvariants(context)

		application.stakingKeeper.IterateValidators(context, func(_ int64, val stakingTypes.ValidatorI) (stop bool) {
			_, _ = application.distributionKeeper.WithdrawValidatorCommission(context, val.GetOperator())
			return false
		})

		delegations := application.stakingKeeper.GetAllDelegations(context)
		for _, delegation := range delegations {
			validatorAddress, err := sdkTypes.ValAddressFromBech32(delegation.ValidatorAddress)
			if err != nil {
				panic(err)
			}

			delegatorAddress, err := sdkTypes.AccAddressFromBech32(delegation.DelegatorAddress)
			if err != nil {
				panic(err)
			}
			_, _ = application.distributionKeeper.WithdrawDelegationRewards(context, delegatorAddress, validatorAddress)
		}

		application.distributionKeeper.DeleteAllValidatorSlashEvents(context)

		application.distributionKeeper.DeleteAllValidatorHistoricalRewards(context)

		height := context.BlockHeight()
		context = context.WithBlockHeight(0)

		application.stakingKeeper.IterateValidators(context, func(_ int64, val stakingTypes.ValidatorI) (stop bool) {
			scraps := application.distributionKeeper.GetValidatorOutstandingRewardsCoins(context, val.GetOperator())
			feePool := application.distributionKeeper.GetFeePool(context)
			feePool.CommunityPool = feePool.CommunityPool.Add(scraps...)
			application.distributionKeeper.SetFeePool(context, feePool)

			if err := application.distributionKeeper.Hooks().AfterValidatorCreated(context, val.GetOperator()); err != nil {
				panic(err)
			}
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
			if err := application.distributionKeeper.Hooks().BeforeDelegationCreated(context, delegatorAddress, validatorAddress); err != nil {
				panic(err)
			}
			if err := application.distributionKeeper.Hooks().AfterDelegationModified(context, delegatorAddress, validatorAddress); err != nil {
				panic(err)
			}
		}

		context = context.WithBlockHeight(height)

		application.stakingKeeper.IterateRedelegations(context, func(_ int64, redelegation stakingTypes.Redelegation) (stop bool) {
			for i := range redelegation.Entries {
				redelegation.Entries[i].CreationHeight = 0
			}
			application.stakingKeeper.SetRedelegation(context, redelegation)
			return false
		})

		application.stakingKeeper.IterateUnbondingDelegations(context, func(_ int64, unbondingDelegation stakingTypes.UnbondingDelegation) (stop bool) {
			for i := range unbondingDelegation.Entries {
				unbondingDelegation.Entries[i].CreationHeight = 0
			}
			application.stakingKeeper.SetUnbondingDelegation(context, unbondingDelegation)
			return false
		})

		kvStore := context.KVStore(application.keys[stakingTypes.StoreKey])
		kvStoreReversePrefixIterator := sdkTypes.KVStoreReversePrefixIterator(kvStore, stakingTypes.ValidatorsKey)
		counter := int16(0)

		func() {
			defer func(kvStoreReversePrefixIterator sdkTypes.Iterator) {
				err := kvStoreReversePrefixIterator.Close()
				if err != nil {
					panic(err)
				}
			}(kvStoreReversePrefixIterator)
			for ; kvStoreReversePrefixIterator.Valid(); kvStoreReversePrefixIterator.Next() {
				addr := sdkTypes.ValAddress(stakingTypes.AddressFromValidatorsKey(kvStoreReversePrefixIterator.Key()))
				validator, found := application.stakingKeeper.GetValidator(context, addr)

				if !found {
					panic("validator not found")
				}

				validator.UnbondingHeight = 0

				if applyWhiteList && !whiteListMap[addr.String()] {
					validator.Jailed = true
				}

				application.stakingKeeper.SetValidator(context, validator)
				counter++
			}
		}()

		_, err := application.stakingKeeper.ApplyAndReturnValidatorSetUpdates(context)
		if err != nil {
			panic(err)
		}

		application.slashingKeeper.IterateValidatorSigningInfos(
			context,
			func(validatorConsAddress sdkTypes.ConsAddress, validatorSigningInfo slashingTypes.ValidatorSigningInfo) (stop bool) {
				validatorSigningInfo.StartHeight = 0
				application.slashingKeeper.SetValidatorSigningInfo(context, validatorConsAddress, validatorSigningInfo)
				return false
			},
		)
	}

	genesisState := application.moduleManager.ExportGenesisForModules(context, application.GetCodec(), modulesToExport)
	applicationState, err := json.MarshalIndent(genesisState, "", "  ")
	if err != nil {
		return serverTypes.ExportedApp{}, err
	}

	validators, err := staking.WriteValidators(context, application.stakingKeeper)
	return serverTypes.ExportedApp{
		AppState:        applicationState,
		Validators:      validators,
		Height:          height,
		ConsensusParams: application.GetConsensusParams(context),
	}, err
}
func (application application) RegisterAPIRoutes(server *api.Server, apiConfig config.APIConfig) {
	clientCtx := server.ClientCtx
	authTx.RegisterGRPCGatewayRoutes(clientCtx, server.GRPCGatewayRouter)
	tmservice.RegisterGRPCGatewayRoutes(clientCtx, server.GRPCGatewayRouter)
	rest.RegisterRESTRoutes(clientCtx, server.Router)
	application.moduleManager.RegisterRESTRoutes(clientCtx, server.Router)
	application.moduleManager.RegisterGRPCGatewayRoutes(clientCtx, server.GRPCGatewayRouter)
	sdkGRPCNode.RegisterGRPCGatewayRoutes(clientCtx, server.GRPCGatewayRouter)
	if err := sdkServer.RegisterSwaggerAPI(server.ClientCtx, server.Router, apiConfig.Swagger); err != nil {
		panic(err)
	}
}
func (application application) RegisterTxService(context client.Context) {
	authTx.RegisterTxService(application.GRPCQueryRouter(), context, application.Simulate, context.InterfaceRegistry)
}
func (application application) RegisterTendermintService(context client.Context) {
	tmservice.RegisterTendermintService(context, application.GRPCQueryRouter(), context.InterfaceRegistry, application.Query)
}
func (application application) AppCreator(logger tendermintLog.Logger, db tendermintDB.DB, writer io.Writer, appOptions serverTypes.AppOptions) serverTypes.Application {
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
	snapshotDB, err := tendermintDB.NewDB("metadata", server.GetAppDBBackend(appOptions), snapshotDir)
	if err != nil {
		panic(err)
	}
	snapshotStore, err := snapshots.NewStore(snapshotDB, snapshotDir)
	if err != nil {
		panic(err)
	}

	snapshotOptions := snapshotsTypes.NewSnapshotOptions(
		cast.ToUint64(appOptions.Get(server.FlagStateSyncSnapshotInterval)),
		cast.ToUint32(appOptions.Get(server.FlagStateSyncSnapshotKeepRecent)),
	)

	chainID := cast.ToString(appOptions.Get(flags.FlagChainID))
	if chainID == "" {
		if appGenesis, err := cometbftTypes.GenesisDocFromFile(filepath.Join(cast.ToString(appOptions.Get(flags.FlagHome)), cast.ToString(appOptions.Get("genesis_file")))); err != nil {
			panic(err)
		} else {
			chainID = appGenesis.ChainID
		}
	}

	return application.Initialize(
		logger,
		db,
		writer,
		true,
		cast.ToUint(appOptions.Get(server.FlagInvCheckPeriod)),
		skipUpgradeHeights,
		cast.ToString(appOptions.Get(flags.FlagHome)),
		appOptions,
		baseapp.SetChainID(chainID),
		baseapp.SetPruning(pruningOpts),
		baseapp.SetMinGasPrices(cast.ToString(appOptions.Get(server.FlagMinGasPrices))),
		baseapp.SetHaltHeight(cast.ToUint64(appOptions.Get(server.FlagHaltHeight))),
		baseapp.SetHaltTime(cast.ToUint64(appOptions.Get(server.FlagHaltTime))),
		baseapp.SetMinRetainBlocks(cast.ToUint64(appOptions.Get(server.FlagMinRetainBlocks))),
		baseapp.SetInterBlockCache(multiStorePersistentCache),
		baseapp.SetTrace(cast.ToBool(appOptions.Get(server.FlagTrace))),
		baseapp.SetIndexEvents(cast.ToStringSlice(appOptions.Get(server.FlagIndexEvents))),
		baseapp.SetSnapshot(snapshotStore, snapshotOptions),
		baseapp.SetIAVLCacheSize(cast.ToInt(appOptions.Get(server.FlagIAVLCacheSize))))
}
func (application application) AppExporter(logger log.Logger, db tendermintDB.DB, writer io.Writer, height int64, forZeroHeight bool, jailAllowedAddresses []string, appOptions serverTypes.AppOptions, modulesToExport []string) (serverTypes.ExportedApp, error) {
	home, ok := appOptions.Get(flags.FlagHome).(string)
	if !ok || home == "" {
		return serverTypes.ExportedApp{}, errors.New("application home is not set")
	}

	var loadLatest bool
	if height == -1 {
		loadLatest = true
	}

	Application := application.Initialize(
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
		if err := application.LoadHeight(height); err != nil {
			return serverTypes.ExportedApp{}, err
		}
	}

	return Application.ExportApplicationStateAndValidators(forZeroHeight, jailAllowedAddresses, modulesToExport)
}
func (application application) ModuleInitFlags(command *cobra.Command) {
	crisis.AddModuleInitFlags(command)
}
func (application application) Initialize(logger tendermintLog.Logger, db tendermintDB.DB, writer io.Writer, loadLatest bool, invCheckPeriod uint, skipUpgradeHeights map[int64]bool, home string, appOptions serverTypes.AppOptions, baseAppOptions ...func(*baseapp.BaseApp)) applications.Application {
	application.BaseApp = baseapp.NewBaseApp(application.name, logger, db, application.GetCodec().TxDecoder(), baseAppOptions...)

	application.BaseApp.SetCommitMultiStoreTracer(writer)
	application.BaseApp.SetVersion(version.Version)
	application.BaseApp.SetInterfaceRegistry(application.GetCodec().InterfaceRegistry())

	application.keys = sdkTypes.NewKVStoreKeys(
		authTypes.StoreKey,
		bankTypes.StoreKey,
		stakingTypes.StoreKey,
		crisisTypes.StoreKey,
		mintTypes.StoreKey,
		distributionTypes.StoreKey,
		slashingTypes.StoreKey,
		govTypes.StoreKey,
		paramsTypes.StoreKey,
		ibcExported.StoreKey,
		upgradeTypes.StoreKey,
		evidenceTypes.StoreKey,
		ibcTransferTypes.StoreKey,
		icaHostTypes.StoreKey,
		capabilityTypes.StoreKey,
		feegrant.StoreKey,
		authzKeeper.StoreKey,
		routerTypes.StoreKey,
		consensusParamTypes.StoreKey,

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
		application.GetCodec(),
		application.GetCodec().GetLegacyAmino(),
		application.keys[paramsTypes.StoreKey],
		transientStoreKeys[paramsTypes.TStoreKey],
	)

	ConsensusParamsKeeper := consensusParamKeeper.NewKeeper(
		application.GetCodec(),
		application.keys[consensusParamTypes.StoreKey],
		authTypes.NewModuleAddress(govTypes.ModuleName).String(),
	)
	application.SetParamStore(&ConsensusParamsKeeper)

	CapabilityKeeper := capabilityKeeper.NewKeeper(application.GetCodec(), application.keys[capabilityTypes.StoreKey], memoryStoreKeys[capabilityTypes.MemStoreKey])
	scopedIBCKeeper := CapabilityKeeper.ScopeToModule(ibcExported.ModuleName)
	scopedTransferKeeper := CapabilityKeeper.ScopeToModule(ibcTransferTypes.ModuleName)
	scopedICAHostKeeper := CapabilityKeeper.ScopeToModule(icaHostTypes.SubModuleName)
	CapabilityKeeper.Seal()

	AccountKeeper := authKeeper.NewAccountKeeper(
		application.GetCodec(),
		application.keys[authTypes.StoreKey],
		authTypes.ProtoBaseAccount,
		application.moduleAccountPermissions,
		sdkTypes.GetConfig().GetBech32AccountAddrPrefix(),
		authTypes.NewModuleAddress(govTypes.ModuleName).String(),
	)

	blacklistedAddresses := make(map[string]bool)
	for account := range application.moduleAccountPermissions {
		blacklistedAddresses[authTypes.NewModuleAddress(account).String()] = !application.tokenReceiveAllowedModules[account]
	}

	BankKeeper := bankKeeper.NewBaseKeeper(
		application.GetCodec(),
		application.keys[bankTypes.StoreKey],
		AccountKeeper,
		blacklistedAddresses,
		authTypes.NewModuleAddress(govTypes.ModuleName).String(),
	)

	AuthzKeeper := authzKeeper.NewKeeper(
		application.keys[authzKeeper.StoreKey],
		application.GetCodec(),
		application.MsgServiceRouter(),
		AccountKeeper,
	)

	FeeGrantKeeper := feeGrantKeeper.NewKeeper(
		application.GetCodec(),
		application.keys[feegrant.StoreKey],
		AccountKeeper,
	)

	application.stakingKeeper = stakingKeeper.NewKeeper(
		application.GetCodec(),
		application.keys[stakingTypes.StoreKey],
		AccountKeeper,
		BankKeeper,
		authTypes.NewModuleAddress(govTypes.ModuleName).String(),
	)

	MintKeeper := mintKeeper.NewKeeper(
		application.GetCodec(),
		application.keys[mintTypes.StoreKey],
		application.stakingKeeper,
		AccountKeeper,
		BankKeeper,
		authTypes.FeeCollectorName,
		authTypes.NewModuleAddress(govTypes.ModuleName).String(),
	)

	application.distributionKeeper = distributionKeeper.NewKeeper(
		application.GetCodec(),
		application.keys[distributionTypes.StoreKey],
		AccountKeeper,
		BankKeeper,
		application.stakingKeeper,
		authTypes.FeeCollectorName,
		authTypes.NewModuleAddress(govTypes.ModuleName).String(),
	)

	application.slashingKeeper = slashingKeeper.NewKeeper(
		application.GetCodec(),
		application.GetCodec().GetLegacyAmino(),
		application.keys[slashingTypes.StoreKey],
		application.stakingKeeper,
		authTypes.NewModuleAddress(govTypes.ModuleName).String(),
	)

	application.crisisKeeper = crisisKeeper.NewKeeper(
		application.GetCodec(),
		application.keys[crisisTypes.StoreKey],
		invCheckPeriod,
		BankKeeper,
		authTypes.FeeCollectorName,
		authTypes.NewModuleAddress(govTypes.ModuleName).String(),
	)

	// UpgradeKeeper must be created before IBCKeeper
	UpgradeKeeper := upgradeKeeper.NewKeeper(
		skipUpgradeHeights,
		application.keys[upgradeTypes.StoreKey],
		application.GetCodec(),
		home,
		application.BaseApp,
		authTypes.NewModuleAddress(govTypes.ModuleName).String(),
	)

	application.stakingKeeper.SetHooks(stakingTypes.NewMultiStakingHooks(application.distributionKeeper.Hooks(), application.slashingKeeper.Hooks()))

	IBCKeeper := ibcKeeper.NewKeeper(
		application.GetCodec(),
		application.keys[ibcExported.StoreKey],
		ParamsKeeper.Subspace(ibcExported.ModuleName),
		application.stakingKeeper,
		UpgradeKeeper,
		scopedIBCKeeper,
	)

	govConfig := govTypes.DefaultConfig()
	govConfig.MaxMetadataLen = 10200

	GovKeeper := govKeeper.NewKeeper(
		application.GetCodec(),
		application.keys[govTypes.StoreKey],
		AccountKeeper,
		BankKeeper,
		application.stakingKeeper,
		application.MsgServiceRouter(),
		govConfig,
		authTypes.NewModuleAddress(govTypes.ModuleName).String(),
	)

	GovKeeper.SetLegacyRouter(
		govTypesV1Beta1.NewRouter().
			AddRoute(govTypes.RouterKey, govTypesV1Beta1.ProposalHandler).
			AddRoute(paramsProposal.RouterKey, params.NewParamChangeProposalHandler(ParamsKeeper)).
			AddRoute(upgradeTypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(UpgradeKeeper)).
			AddRoute(ibcClientTypes.RouterKey, ibcClient.NewClientProposalHandler(IBCKeeper.ClientKeeper)))

	IBCFeeKeeper := ibcFeeKeeper.NewKeeper(
		application.GetCodec(),
		application.keys[ibcFeeTypes.StoreKey],
		IBCKeeper.ChannelKeeper, // may be replaced with IBC middleware
		IBCKeeper.ChannelKeeper,
		&IBCKeeper.PortKeeper,
		AccountKeeper,
		BankKeeper,
	)

	RateLimitKeeper := *rateLimitKeeper.NewKeeper(
		application.GetCodec(),                    // BinaryCodec
		application.keys[rateLimitTypes.StoreKey], // StoreKey
		ParamsKeeper.Subspace(rateLimitTypes.ModuleName),
		authTypes.NewModuleAddress(govTypes.ModuleName).String(),
		BankKeeper,
		IBCKeeper.ChannelKeeper, // ChannelKeeper
		IBCFeeKeeper,            // ICS4Wrapper
	)

	// From gaia: RouterKeeper must be created before TransferKeeper
	RouterKeeper := routerKeeper.NewKeeper(
		application.GetCodec(),
		application.keys[routerTypes.StoreKey],
		nil,
		IBCKeeper.ChannelKeeper,
		application.distributionKeeper,
		BankKeeper,
		RateLimitKeeper,
		authTypes.NewModuleAddress(govTypes.ModuleName).String(),
	)

	IBCTransferKeeper := ibcTransferKeeper.NewKeeper(
		application.GetCodec(),
		application.keys[ibcTransferTypes.StoreKey],
		ParamsKeeper.Subspace(ibcTransferTypes.ModuleName),
		RouterKeeper,
		IBCKeeper.ChannelKeeper,
		&IBCKeeper.PortKeeper,
		AccountKeeper,
		BankKeeper,
		scopedTransferKeeper,
	)

	RouterKeeper.SetTransferKeeper(IBCTransferKeeper)

	ICAHostKeeper := icaHostKeeper.NewKeeper(
		application.GetCodec(),
		application.keys[icaHostTypes.StoreKey],
		ParamsKeeper.Subspace(icaHostTypes.SubModuleName).WithKeyTable(icaHostTypes.ParamKeyTable()),
		IBCKeeper.ChannelKeeper,
		&IBCKeeper.PortKeeper,
		AccountKeeper,
		scopedICAHostKeeper,
		application.MsgServiceRouter(),
	)

	icaHostIBCModule := icaHost.NewIBCModule(ICAHostKeeper)

	var ibcStack ibcPortTypes.IBCModule
	ibcStack = ibcTransfer.NewIBCModule(IBCTransferKeeper)
	ibcStack = router.NewIBCMiddleware(
		ibcStack,
		RouterKeeper,
		0,
		routerKeeper.DefaultForwardTransferPacketTimeoutTimestamp,
		routerKeeper.DefaultRefundTransferPacketTimeoutTimestamp,
	)

	ibcRouter := ibcPortTypes.NewRouter()
	ibcRouter.AddRoute(icaHostTypes.SubModuleName, icaHostIBCModule).
		AddRoute(ibcTransferTypes.ModuleName, ibcStack)

	IBCKeeper.SetRouter(ibcRouter)

	EvidenceKeeper := *evidenceKeeper.NewKeeper(
		application.GetCodec(),
		application.keys[evidenceTypes.StoreKey],
		&application.stakingKeeper,
		application.slashingKeeper,
	)

	metasModule := metas.Prototype().Initialize(
		application.keys[metas.Prototype().Name()],
		ParamsKeeper.Subspace(metas.Prototype().Name()),
	)

	classificationsModule := classifications.Prototype().Initialize(
		application.keys[classifications.Prototype().Name()],
		ParamsKeeper.Subspace(classifications.Prototype().Name()),
		BankKeeper,
		application.stakingKeeper,
	)

	maintainersModule := maintainers.Prototype().Initialize(
		application.keys[maintainers.Prototype().Name()],
		ParamsKeeper.Subspace(maintainers.Prototype().Name()),
		classificationsModule.GetAuxiliary(member.Auxiliary.GetName()),
	)
	identitiesModule := identities.Prototype().Initialize(
		application.keys[identities.Prototype().Name()],
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
		application.keys[splits.Prototype().Name()],
		ParamsKeeper.Subspace(splits.Prototype().Name()),
	)
	assetsModule := assets.Prototype().Initialize(
		application.keys[assets.Prototype().Name()],
		ParamsKeeper.Subspace(assets.Prototype().Name()),
		BankKeeper,
		classificationsModule.GetAuxiliary(conform.Auxiliary.GetName()),
		classificationsModule.GetAuxiliary(define.Auxiliary.GetName()),
		classificationsModule.GetAuxiliary(bond.Auxiliary.GetName()),
		classificationsModule.GetAuxiliary(unbond.Auxiliary.GetName()),
		identitiesModule.GetAuxiliary(authenticate.Auxiliary.GetName()),
		maintainersModule.GetAuxiliary(deputize.Auxiliary.GetName()),
		maintainersModule.GetAuxiliary(maintain.Auxiliary.GetName()),
		maintainersModule.GetAuxiliary(revoke.Auxiliary.GetName()),
		maintainersModule.GetAuxiliary(super.Auxiliary.GetName()),
		maintainersModule.GetAuxiliary(authorize.Auxiliary.GetName()),
		metasModule.GetAuxiliary(supplement.Auxiliary.GetName()),
		splitsModule.GetAuxiliary(transfer.Auxiliary.GetName()),
		splitsModule.GetAuxiliary(renumerate.Auxiliary.GetName()),
		splitsModule.GetAuxiliary(burnSplits.Auxiliary.GetName()),
		splitsModule.GetAuxiliary(purge.Auxiliary.GetName()),
		splitsModule.GetAuxiliary(splitsMint.Auxiliary.GetName()),
	)
	ordersModule := orders.Prototype().Initialize(
		application.keys[orders.Prototype().Name()],
		ParamsKeeper.Subspace(orders.Prototype().Name()),
		identitiesModule.GetAuxiliary(authenticate.Auxiliary.GetName()),
		classificationsModule.GetAuxiliary(bond.Auxiliary.GetName()),
		classificationsModule.GetAuxiliary(burn.Auxiliary.GetName()),
		classificationsModule.GetAuxiliary(unbond.Auxiliary.GetName()),
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

	application.moduleManager = module.NewManager(
		genutil.NewAppModule(AccountKeeper, application.stakingKeeper, application.DeliverTx, application.GetCodec()),
		auth.NewAppModule(application.GetCodec(), AccountKeeper, nil),
		vesting.NewAppModule(AccountKeeper, BankKeeper),
		bank.NewAppModule(application.GetCodec(), BankKeeper, AccountKeeper),
		capability.NewAppModule(application.GetCodec(), *CapabilityKeeper),
		crisis.NewAppModule(&application.crisisKeeper, cast.ToBool(appOptions.Get(crisis.FlagSkipGenesisInvariants))),
		gov.NewAppModule(application.GetCodec(), GovKeeper, AccountKeeper, BankKeeper),
		mint.NewAppModule(application.GetCodec(), MintKeeper, AccountKeeper),
		slashing.NewAppModule(application.GetCodec(), application.slashingKeeper, AccountKeeper, BankKeeper, application.stakingKeeper),
		distribution.NewAppModule(application.GetCodec(), application.distributionKeeper, AccountKeeper, BankKeeper, application.stakingKeeper),
		staking.NewAppModule(application.GetCodec(), application.stakingKeeper, AccountKeeper, BankKeeper),
		upgrade.NewAppModule(UpgradeKeeper),
		evidence.NewAppModule(EvidenceKeeper),
		feeGrantModule.NewAppModule(application.GetCodec(), AccountKeeper, BankKeeper, FeeGrantKeeper, application.GetCodec().InterfaceRegistry()),
		authzModule.NewAppModule(application.GetCodec(), AuthzKeeper, AccountKeeper, BankKeeper, application.GetCodec().InterfaceRegistry()),
		ibc.NewAppModule(IBCKeeper),
		params.NewAppModule(ParamsKeeper),
		ibcTransfer.NewAppModule(IBCTransferKeeper),
		ica.NewAppModule(nil, &ICAHostKeeper),
		router.NewAppModule(RouterKeeper),

		assetsModule,
		classificationsModule,
		identitiesModule,
		maintainersModule,
		metasModule,
		ordersModule,
		splitsModule,
	)

	application.moduleManager.SetOrderBeginBlockers(
		upgradeTypes.ModuleName,
		capabilityTypes.ModuleName,
		mintTypes.ModuleName,
		distributionTypes.ModuleName,
		slashingTypes.ModuleName,
		evidenceTypes.ModuleName,
		stakingTypes.ModuleName,
		authTypes.ModuleName,
		bankTypes.ModuleName,
		govTypes.ModuleName,
		crisisTypes.ModuleName,
		ibcTransferTypes.ModuleName,
		ibcExported.ModuleName,
		icaTypes.ModuleName,
		routerTypes.ModuleName,
		genutilTypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramsTypes.ModuleName,
		vestingTypes.ModuleName,

		// Order doesn't matter here currently since all the BeginBlock functions are empty
		assets.Prototype().Name(),
		classifications.Prototype().Name(),
		identities.Prototype().Name(),
		maintainers.Prototype().Name(),
		metas.Prototype().Name(),
		orders.Prototype().Name(),
		splits.Prototype().Name(),
	)
	application.moduleManager.SetOrderEndBlockers(
		crisisTypes.ModuleName,
		govTypes.ModuleName,
		stakingTypes.ModuleName,
		ibcTransferTypes.ModuleName,
		ibcExported.ModuleName,
		icaTypes.ModuleName,
		routerTypes.ModuleName,
		capabilityTypes.ModuleName,
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
		upgradeTypes.ModuleName,
		vestingTypes.ModuleName,

		orders.Prototype().Name(),
		assets.Prototype().Name(),
		classifications.Prototype().Name(),
		identities.Prototype().Name(),
		maintainers.Prototype().Name(),
		metas.Prototype().Name(),
		splits.Prototype().Name(),
	)
	application.moduleManager.SetOrderInitGenesis(
		capabilityTypes.ModuleName,
		authTypes.ModuleName,
		bankTypes.ModuleName,
		distributionTypes.ModuleName,
		govTypes.ModuleName,
		stakingTypes.ModuleName,
		slashingTypes.ModuleName,
		mintTypes.ModuleName,
		crisisTypes.ModuleName,
		genutilTypes.ModuleName,
		ibcTransferTypes.ModuleName,
		ibcExported.ModuleName,
		icaTypes.ModuleName,
		evidenceTypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		routerTypes.ModuleName,
		paramsTypes.ModuleName,
		upgradeTypes.ModuleName,
		vestingTypes.ModuleName,

		// meta should be initialized first because rest of the modules depends on it
		metas.Prototype().Name(),
		splits.Prototype().Name(),
		classifications.Prototype().Name(),
		maintainers.Prototype().Name(),
		identities.Prototype().Name(),
		assets.Prototype().Name(),
		orders.Prototype().Name(),
	)
	application.moduleManager.RegisterInvariants(&application.crisisKeeper)
	application.moduleManager.RegisterRoutes(application.Router(), application.QueryRouter(), application.GetCodec().GetLegacyAmino())

	configurator := module.NewConfigurator(application.GetCodec(), application.MsgServiceRouter(), application.GRPCQueryRouter())
	application.moduleManager.RegisterServices(configurator)

	module.NewSimulationManager(
		auth.NewAppModule(application.GetCodec(), AccountKeeper, authSimulation.RandomGenesisAccounts),
		bank.NewAppModule(application.GetCodec(), BankKeeper, AccountKeeper),
		capability.NewAppModule(application.GetCodec(), *CapabilityKeeper),
		feeGrantModule.NewAppModule(application.GetCodec(), AccountKeeper, BankKeeper, FeeGrantKeeper, application.GetCodec().InterfaceRegistry()),
		authzModule.NewAppModule(application.GetCodec(), AuthzKeeper, AccountKeeper, BankKeeper, application.GetCodec().InterfaceRegistry()),
		gov.NewAppModule(application.GetCodec(), GovKeeper, AccountKeeper, BankKeeper),
		mint.NewAppModule(application.GetCodec(), MintKeeper, AccountKeeper),
		staking.NewAppModule(application.GetCodec(), application.stakingKeeper, AccountKeeper, BankKeeper),
		distribution.NewAppModule(application.GetCodec(), application.distributionKeeper, AccountKeeper, BankKeeper, application.stakingKeeper),
		slashing.NewAppModule(application.GetCodec(), application.slashingKeeper, AccountKeeper, BankKeeper, application.stakingKeeper),
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

	application.MountKVStores(application.keys)
	application.MountTransientStores(transientStoreKeys)
	application.MountMemoryStores(memoryStoreKeys)
	application.SetAnteHandler(sdkTypes.ChainAnteDecorators(
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
		ante.NewSigVerificationDecorator(AccountKeeper, application.GetCodec().SignModeHandler()),
		ante.NewIncrementSequenceDecorator(AccountKeeper),
		ibcAnte.NewAnteDecorator(IBCKeeper),
	))
	application.SetBeginBlocker(application.moduleManager.BeginBlock)
	application.SetEndBlocker(application.moduleManager.EndBlock)
	application.SetInitChainer(func(context sdkTypes.Context, requestInitChain abciTypes.RequestInitChain) abciTypes.ResponseInitChain {
		var genesisState map[string]json.RawMessage
		if err := tendermintJSON.Unmarshal(requestInitChain.AppStateBytes, &genesisState); err != nil {
			panic(err)
		}

		UpgradeKeeper.SetModuleVersionMap(context, application.moduleManager.GetVersionMap())

		return application.moduleManager.InitGenesis(context, application.GetCodec(), genesisState)
	})

	UpgradeKeeper.SetUpgradeHandler(
		constants.UpgradeName,
		func(ctx sdkTypes.Context, _ upgradeTypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			return application.moduleManager.RunMigrations(ctx, configurator, fromVM)
		},
	)

	upgradeInfo, err := UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}

	if upgradeInfo.Name == constants.UpgradeName && !UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := storeTypes.StoreUpgrades{
			Added: []string{assets.Prototype().Name(), classifications.Prototype().Name(), identities.Prototype().Name(), maintainers.Prototype().Name(), metas.Prototype().Name(), orders.Prototype().Name(), splits.Prototype().Name()},
		}

		application.SetStoreLoader(upgradeTypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}

	if loadLatest {
		err := application.LoadLatestVersion()
		if err != nil {
			tendermintOS.Exit(err.Error())
		}
	}

	return &application
}

func NewApplication(name string, moduleManager helpers.ModuleManager, moduleAccountPermissions map[string][]string, tokenReceiveAllowedModules map[string]bool) applications.Application {
	return &application{
		name:                       name,
		moduleManager:              moduleManager,
		moduleAccountPermissions:   moduleAccountPermissions,
		tokenReceiveAllowedModules: tokenReceiveAllowedModules,
	}
}
