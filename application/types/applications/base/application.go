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

	"github.com/AssetMantle/modules/x/maintainers/auxiliaries/authorize"
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
	icaHost "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts/host"
	icaHostKeeper "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts/host/keeper"
	icaHostTypes "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts/host/types"
	icaTypes "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts/types"
	ibcTransfer "github.com/cosmos/ibc-go/v4/modules/apps/transfer"
	ibcTransferKeeper "github.com/cosmos/ibc-go/v4/modules/apps/transfer/keeper"
	ibcTransferTypes "github.com/cosmos/ibc-go/v4/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v4/modules/core"
	ibcClient "github.com/cosmos/ibc-go/v4/modules/core/02-client"
	ibcClientTypes "github.com/cosmos/ibc-go/v4/modules/core/02-client/types"
	ibcPortTypes "github.com/cosmos/ibc-go/v4/modules/core/05-port/types"
	ibcHost "github.com/cosmos/ibc-go/v4/modules/core/24-host"
	ibcAnte "github.com/cosmos/ibc-go/v4/modules/core/ante"
	ibcKeeper "github.com/cosmos/ibc-go/v4/modules/core/keeper"
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
	"github.com/AssetMantle/modules/x/classifications/auxiliaries/burn"
	"github.com/AssetMantle/modules/x/classifications/auxiliaries/conform"
	"github.com/AssetMantle/modules/x/classifications/auxiliaries/define"
	"github.com/AssetMantle/modules/x/classifications/auxiliaries/member"
	"github.com/AssetMantle/modules/x/classifications/auxiliaries/unbond"
	"github.com/AssetMantle/modules/x/identities"
	"github.com/AssetMantle/modules/x/identities/auxiliaries/authenticate"
	"github.com/AssetMantle/modules/x/maintainers"
	"github.com/AssetMantle/modules/x/maintainers/auxiliaries/deputize"
	"github.com/AssetMantle/modules/x/maintainers/auxiliaries/maintain"
	"github.com/AssetMantle/modules/x/maintainers/auxiliaries/revoke"
	"github.com/AssetMantle/modules/x/maintainers/auxiliaries/super"
	"github.com/AssetMantle/modules/x/metas"
	"github.com/AssetMantle/modules/x/metas/auxiliaries/supplement"
	"github.com/AssetMantle/modules/x/orders"
	"github.com/AssetMantle/modules/x/splits"
	splitsBurn "github.com/AssetMantle/modules/x/splits/auxiliaries/burn"
	splitsMint "github.com/AssetMantle/modules/x/splits/auxiliaries/mint"
	"github.com/AssetMantle/modules/x/splits/auxiliaries/renumerate"
	"github.com/AssetMantle/modules/x/splits/auxiliaries/transfer"

	"github.com/AssetMantle/node/application/types/applications"
	"github.com/AssetMantle/node/application/types/applications/constants"
	"github.com/AssetMantle/node/utilities/rest"
)

type application struct {
	name string

	moduleBasicManager module.BasicManager

	codec helpers.Codec

	moduleAccountPermissions   map[string][]string
	tokenReceiveAllowedModules map[string]bool

	keys map[string]*sdkTypes.KVStoreKey

	stakingKeeper      stakingKeeper.Keeper
	slashingKeeper     slashingKeeper.Keeper
	distributionKeeper distributionKeeper.Keeper
	crisisKeeper       crisisKeeper.Keeper

	moduleManager *module.Manager

	baseapp.BaseApp
}

var _ applications.Application = (*application)(nil)

func (application application) GetDefaultNodeHome() string {
	return os.ExpandEnv("$HOME/." + application.name)
}
func (application application) GetDefaultClientHome() string {
	return os.ExpandEnv("$HOME/." + application.name)
}
func (application application) GetModuleBasicManager() module.BasicManager {
	return application.moduleBasicManager
}
func (application application) GetCodec() helpers.Codec {
	return application.codec
}
func (application application) LoadHeight(height int64) error {
	return application.BaseApp.LoadVersion(height)
}
func (application application) ExportApplicationStateAndValidators(forZeroHeight bool, jailWhiteList []string) (serverTypes.ExportedApp, error) {
	context := application.BaseApp.NewContext(true, protoTendermintTypes.Header{Height: application.BaseApp.LastBlockHeight()})

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

			application.distributionKeeper.Hooks().AfterValidatorCreated(context, val.GetOperator())
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
			application.distributionKeeper.Hooks().BeforeDelegationCreated(context, delegatorAddress, validatorAddress)
			application.distributionKeeper.Hooks().AfterDelegationModified(context, delegatorAddress, validatorAddress)
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

		for ; kvStoreReversePrefixIterator.Valid(); kvStoreReversePrefixIterator.Next() {
			addr := sdkTypes.ValAddress(stakingTypes.AddressFromValidatorsKey(kvStoreReversePrefixIterator.Key()))
			validator, found := application.stakingKeeper.GetValidator(context, addr)

			if !found {
				panic("Validator not found!")
			}

			validator.UnbondingHeight = 0

			if applyWhiteList && !whiteListMap[addr.String()] {
				validator.Jailed = true
			}

			application.stakingKeeper.SetValidator(context, validator)
			counter++
		}

		if err := kvStoreReversePrefixIterator.Close(); err != nil {
			log.Fatal(err)
		}

		if _, err := application.stakingKeeper.ApplyAndReturnValidatorSetUpdates(context); err != nil {
			log.Fatal(err)
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

	genesisState := application.moduleManager.ExportGenesis(context, application.GetCodec())
	applicationState, err := json.MarshalIndent(genesisState, "", "  ")
	if err != nil {
		return serverTypes.ExportedApp{}, err
	}

	validators, err := staking.WriteValidators(context, application.stakingKeeper)
	return serverTypes.ExportedApp{
		AppState:        applicationState,
		Validators:      validators,
		Height:          height,
		ConsensusParams: application.BaseApp.GetConsensusParams(context),
	}, err
}
func (application application) RegisterAPIRoutes(server *api.Server, apiConfig config.APIConfig) {
	clientCtx := server.ClientCtx
	rpc.RegisterRoutes(clientCtx, server.Router)
	authRest.RegisterTxRoutes(clientCtx, server.Router)
	authTx.RegisterGRPCGatewayRoutes(clientCtx, server.GRPCGatewayRouter)
	tmservice.RegisterGRPCGatewayRoutes(clientCtx, server.GRPCGatewayRouter)
	documentIDGetters.RegisterRESTRoutes(clientCtx, server.Router)
	application.moduleBasicManager.RegisterRESTRoutes(clientCtx, server.Router)
	application.moduleBasicManager.RegisterGRPCGatewayRoutes(clientCtx, server.GRPCGatewayRouter)
	rest.RegisterRESTRoutes(clientCtx, server.Router)
	if apiConfig.Swagger {
		Fs, err := fs.New()
		if err != nil {
			panic(err)
		}

		server.Router.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", http.FileServer(Fs)))
	}
}
func (application application) RegisterTxService(context client.Context) {
	authTx.RegisterTxService(application.BaseApp.GRPCQueryRouter(), context, application.BaseApp.Simulate, context.InterfaceRegistry)
}
func (application application) RegisterTendermintService(context client.Context) {
	tmservice.RegisterTendermintService(application.BaseApp.GRPCQueryRouter(), context, context.InterfaceRegistry)
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
	snapshotDB, err := sdkTypes.NewLevelDB("metadata", snapshotDir)
	if err != nil {
		panic(err)
	}
	snapshotStore, err := snapshots.NewStore(snapshotDB, snapshotDir)
	if err != nil {
		panic(err)
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
func (application application) AppExporter(logger tendermintLog.Logger, db tendermintDB.DB, writer io.Writer, height int64, forZeroHeight bool, jailAllowedAddrs []string, appOptions serverTypes.AppOptions) (serverTypes.ExportedApp, error) {
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

	return Application.ExportApplicationStateAndValidators(forZeroHeight, jailAllowedAddrs)
}
func (application application) ModuleInitFlags(command *cobra.Command) {
	crisis.AddModuleInitFlags(command)
}
func (application application) Initialize(logger tendermintLog.Logger, db tendermintDB.DB, writer io.Writer, loadLatest bool, invCheckPeriod uint, skipUpgradeHeights map[int64]bool, home string, appOptions serverTypes.AppOptions, baseAppOptions ...func(*baseapp.BaseApp)) applications.Application {
	application.BaseApp = *baseapp.NewBaseApp(application.name, logger, db, application.GetCodec().TxDecoder(), baseAppOptions...)
	application.BaseApp.SetCommitMultiStoreTracer(writer)
	application.BaseApp.SetVersion(version.Version)
	application.BaseApp.SetInterfaceRegistry(application.GetCodec().InterfaceRegistry())

	application.keys = sdkTypes.NewKVStoreKeys(
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

	application.BaseApp.SetParamStore(ParamsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramsKeeper.ConsensusParamsKeyTable()))

	CapabilityKeeper := capabilityKeeper.NewKeeper(application.GetCodec(), application.keys[capabilityTypes.StoreKey], memoryStoreKeys[capabilityTypes.MemStoreKey])
	scopedIBCKeeper := CapabilityKeeper.ScopeToModule(ibcHost.ModuleName)
	scopedTransferKeeper := CapabilityKeeper.ScopeToModule(ibcTransferTypes.ModuleName)
	scopedICAHostKeeper := CapabilityKeeper.ScopeToModule(icaHostTypes.SubModuleName)
	CapabilityKeeper.Seal()

	AccountKeeper := authKeeper.NewAccountKeeper(
		application.GetCodec(),
		application.keys[authTypes.StoreKey],
		ParamsKeeper.Subspace(authTypes.ModuleName),
		authTypes.ProtoBaseAccount,
		application.moduleAccountPermissions,
	)

	blacklistedAddresses := make(map[string]bool)
	for account := range application.moduleAccountPermissions {
		blacklistedAddresses[authTypes.NewModuleAddress(account).String()] = !application.tokenReceiveAllowedModules[account]
	}

	BankKeeper := bankKeeper.NewBaseKeeper(
		application.GetCodec(),
		application.keys[bankTypes.StoreKey],
		AccountKeeper,
		ParamsKeeper.Subspace(bankTypes.ModuleName),
		blacklistedAddresses,
	)

	AuthzKeeper := authzKeeper.NewKeeper(
		application.keys[authzKeeper.StoreKey],
		application.GetCodec(),
		application.BaseApp.MsgServiceRouter(),
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
		ParamsKeeper.Subspace(stakingTypes.ModuleName),
	)

	MintKeeper := mintKeeper.NewKeeper(
		application.GetCodec(),
		application.keys[mintTypes.StoreKey],
		ParamsKeeper.Subspace(mintTypes.ModuleName),
		&application.stakingKeeper,
		AccountKeeper,
		BankKeeper,
		authTypes.FeeCollectorName,
	)

	application.distributionKeeper = distributionKeeper.NewKeeper(
		application.GetCodec(),
		application.keys[distributionTypes.StoreKey],
		ParamsKeeper.Subspace(distributionTypes.ModuleName),
		AccountKeeper,
		BankKeeper,
		&application.stakingKeeper,
		authTypes.FeeCollectorName,
		blacklistedAddresses,
	)

	application.slashingKeeper = slashingKeeper.NewKeeper(
		application.GetCodec(),
		application.keys[slashingTypes.StoreKey],
		&application.stakingKeeper,
		ParamsKeeper.Subspace(slashingTypes.ModuleName),
	)

	application.crisisKeeper = crisisKeeper.NewKeeper(
		ParamsKeeper.Subspace(crisisTypes.ModuleName),
		invCheckPeriod,
		BankKeeper,
		authTypes.FeeCollectorName,
	)

	UpgradeKeeper := upgradeKeeper.NewKeeper(
		skipUpgradeHeights,
		application.keys[upgradeTypes.StoreKey],
		application.GetCodec(),
		home,
		&application.BaseApp,
	)

	application.stakingKeeper = *application.stakingKeeper.SetHooks(stakingTypes.NewMultiStakingHooks(application.distributionKeeper.Hooks(), application.slashingKeeper.Hooks()))

	IBCKeeper := ibcKeeper.NewKeeper(
		application.GetCodec(),
		application.keys[ibcHost.StoreKey],
		ParamsKeeper.Subspace(ibcHost.ModuleName),
		application.stakingKeeper,
		UpgradeKeeper,
		scopedIBCKeeper,
	)

	govRouter := govTypes.NewRouter().
		AddRoute(govTypes.RouterKey, govTypes.ProposalHandler).
		AddRoute(paramsProposal.RouterKey, params.NewParamChangeProposalHandler(ParamsKeeper)).
		AddRoute(distributionTypes.RouterKey, distribution.NewCommunityPoolSpendProposalHandler(application.distributionKeeper)).
		AddRoute(upgradeTypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(UpgradeKeeper)).
		AddRoute(ibcClientTypes.RouterKey, ibcClient.NewClientProposalHandler(IBCKeeper.ClientKeeper))

	GovKeeper := govKeeper.NewKeeper(
		application.GetCodec(),
		application.keys[govTypes.StoreKey],
		ParamsKeeper.Subspace(govTypes.ModuleName).WithKeyTable(govTypes.ParamKeyTable()),
		AccountKeeper,
		BankKeeper,
		&application.stakingKeeper,
		govRouter,
	)

	IBCTransferKeeper := ibcTransferKeeper.NewKeeper(
		application.GetCodec(),
		application.keys[ibcTransferTypes.StoreKey],
		ParamsKeeper.Subspace(ibcTransferTypes.ModuleName),
		IBCKeeper.ChannelKeeper,
		IBCKeeper.ChannelKeeper,
		&IBCKeeper.PortKeeper,
		AccountKeeper,
		BankKeeper,
		scopedTransferKeeper,
	)

	ICAHostKeeper := icaHostKeeper.NewKeeper(
		application.GetCodec(),
		application.keys[icaHostTypes.StoreKey],
		ParamsKeeper.Subspace(icaHostTypes.SubModuleName),
		IBCKeeper.ChannelKeeper,
		&IBCKeeper.PortKeeper,
		AccountKeeper,
		scopedICAHostKeeper,
		application.MsgServiceRouter(),
	)

	icaHostIBCModule := icaHost.NewIBCModule(ICAHostKeeper)

	var ibcStack ibcPortTypes.IBCModule
	ibcStack = ibcTransfer.NewIBCModule(IBCTransferKeeper)

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
		BankKeeper,
		identitiesModule.GetAuxiliary(authenticate.Auxiliary.GetName()),
	)
	assetsModule := assets.Prototype().Initialize(
		application.keys[assets.Prototype().Name()],
		ParamsKeeper.Subspace(assets.Prototype().Name()),
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
		splitsModule.GetAuxiliary(renumerate.Auxiliary.GetName()),
		splitsModule.GetAuxiliary(splitsBurn.Auxiliary.GetName()),
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
		genutil.NewAppModule(AccountKeeper, application.stakingKeeper, application.BaseApp.DeliverTx, application.GetCodec()),
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
	application.moduleManager.SetOrderEndBlockers(
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
	application.moduleManager.SetOrderInitGenesis(
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
	application.moduleManager.RegisterInvariants(&application.crisisKeeper)
	application.moduleManager.RegisterRoutes(application.BaseApp.Router(), application.BaseApp.QueryRouter(), application.GetCodec().GetLegacyAmino())

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

	application.BaseApp.MountKVStores(application.keys)
	application.BaseApp.MountTransientStores(transientStoreKeys)
	application.BaseApp.MountMemoryStores(memoryStoreKeys)
	application.BaseApp.SetAnteHandler(sdkTypes.ChainAnteDecorators(
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
	application.BaseApp.SetBeginBlocker(application.moduleManager.BeginBlock)
	application.BaseApp.SetEndBlocker(application.moduleManager.EndBlock)
	application.BaseApp.SetInitChainer(func(context sdkTypes.Context, requestInitChain abciTypes.RequestInitChain) abciTypes.ResponseInitChain {
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
		err := application.BaseApp.LoadLatestVersion()
		if err != nil {
			tendermintOS.Exit(err.Error())
		}
	}

	return &application
}

func NewApplication(name string, moduleBasicManager module.BasicManager, moduleAccountPermissions map[string][]string, tokenReceiveAllowedModules map[string]bool) applications.Application {
	return &application{
		name:                       name,
		moduleBasicManager:         moduleBasicManager,
		codec:                      base.CodecPrototype().Initialize(moduleBasicManager),
		moduleAccountPermissions:   moduleAccountPermissions,
		tokenReceiveAllowedModules: tokenReceiveAllowedModules,
	}
}
