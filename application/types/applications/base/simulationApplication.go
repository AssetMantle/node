// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package base

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	serverTypes "github.com/cosmos/cosmos-sdk/server/types"
	simAppParams "github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/cosmos/cosmos-sdk/store/streaming"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
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
	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"
	"github.com/spf13/cast"
	abci "github.com/tendermint/tendermint/abci/types"
	tmLog "github.com/tendermint/tendermint/libs/log"
	tmos "github.com/tendermint/tendermint/libs/os"
	tmProto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/AssetMantle/node/application/types/applications"
	simulationMake "github.com/AssetMantle/node/simulation/make"
)

type GenesisState map[string]json.RawMessage

func NewDefaultGenesisState(cdc codec.JSONCodec) GenesisState {
	return simulationMake.ModuleBasicManagers.DefaultGenesis(cdc)
}

type SimulationApplication struct {
	*baseapp.BaseApp
	legacyAmino       *codec.LegacyAmino
	appCodec          codec.Codec
	interfaceRegistry types.InterfaceRegistry

	invCheckPeriod uint

	keys          map[string]*sdkTypes.KVStoreKey
	transientKeys map[string]*sdkTypes.TransientStoreKey
	memoryKeys    map[string]*sdkTypes.MemoryStoreKey

	AccountKeeper      authKeeper.AccountKeeper
	BankKeeper         bankKeeper.Keeper
	CapabilityKeeper   *capabilityKeeper.Keeper
	StakingKeeper      stakingKeeper.Keeper
	SlashingKeeper     slashingKeeper.Keeper
	MintKeeper         mintKeeper.Keeper
	DistributionKeeper distributionKeeper.Keeper
	GovKeeper          govKeeper.Keeper
	CrisisKeeper       crisisKeeper.Keeper
	UpgradeKeeper      upgradeKeeper.Keeper
	ParamsKeeper       paramsKeeper.Keeper
	AuthzKeeper        authzKeeper.Keeper
	EvidenceKeeper     evidenceKeeper.Keeper
	FeeGrantKeeper     feeGrantKeeper.Keeper

	moduleManager     *module.Manager
	simulationManager *module.SimulationManager
	configurator      module.Configurator
}

var _ applications.SimulationApplication = (*SimulationApplication)(nil)

func (app *SimulationApplication) ExportAppStateAndValidators(forZeroHeight bool, jailAllowedAddrs []string) (serverTypes.ExportedApp, error) {
	ctx := app.NewContext(true, tmProto.Header{Height: app.LastBlockHeight()})

	height := app.LastBlockHeight() + 1
	if forZeroHeight {
		height = 0
		app.prepForZeroHeightGenesis(ctx, jailAllowedAddrs)
	}

	genState := app.moduleManager.ExportGenesis(ctx, app.appCodec)
	appState, err := json.MarshalIndent(genState, "", "  ")
	if err != nil {
		return serverTypes.ExportedApp{}, err
	}

	validators, err := staking.WriteValidators(ctx, app.StakingKeeper)
	return serverTypes.ExportedApp{
		AppState:        appState,
		Validators:      validators,
		Height:          height,
		ConsensusParams: app.BaseApp.GetConsensusParams(ctx),
	}, err
}

func (app *SimulationApplication) GetBaseApp() *baseapp.BaseApp {
	return app.BaseApp
}

func (app *SimulationApplication) GetCrisisKeeper() crisisKeeper.Keeper {
	return app.CrisisKeeper
}

func (app *SimulationApplication) prepForZeroHeightGenesis(ctx sdkTypes.Context, jailAllowedAddrs []string) {
	applyAllowedAddrs := false

	if len(jailAllowedAddrs) > 0 {
		applyAllowedAddrs = true
	}

	allowedAddrsMap := make(map[string]bool)

	for _, addr := range jailAllowedAddrs {
		_, err := sdkTypes.ValAddressFromBech32(addr)
		if err != nil {
			log.Fatal(err)
		}
		allowedAddrsMap[addr] = true
	}

	app.CrisisKeeper.AssertInvariants(ctx)

	app.StakingKeeper.IterateValidators(ctx, func(_ int64, val stakingTypes.ValidatorI) (stop bool) {
		_, _ = app.DistributionKeeper.WithdrawValidatorCommission(ctx, val.GetOperator())
		return false
	})

	dels := app.StakingKeeper.GetAllDelegations(ctx)
	for _, delegation := range dels {
		valAddr, err := sdkTypes.ValAddressFromBech32(delegation.ValidatorAddress)
		if err != nil {
			panic(err)
		}

		delAddr := sdkTypes.MustAccAddressFromBech32(delegation.DelegatorAddress)

		_, _ = app.DistributionKeeper.WithdrawDelegationRewards(ctx, delAddr, valAddr)
	}

	app.DistributionKeeper.DeleteAllValidatorSlashEvents(ctx)

	app.DistributionKeeper.DeleteAllValidatorHistoricalRewards(ctx)

	height := ctx.BlockHeight()
	ctx = ctx.WithBlockHeight(0)

	app.StakingKeeper.IterateValidators(ctx, func(_ int64, val stakingTypes.ValidatorI) (stop bool) {
		scraps := app.DistributionKeeper.GetValidatorOutstandingRewardsCoins(ctx, val.GetOperator())
		feePool := app.DistributionKeeper.GetFeePool(ctx)
		feePool.CommunityPool = feePool.CommunityPool.Add(scraps...)
		app.DistributionKeeper.SetFeePool(ctx, feePool)

		app.DistributionKeeper.Hooks().AfterValidatorCreated(ctx, val.GetOperator())
		return false
	})

	for _, del := range dels {
		valAddr, err := sdkTypes.ValAddressFromBech32(del.ValidatorAddress)
		if err != nil {
			panic(err)
		}
		delAddr := sdkTypes.MustAccAddressFromBech32(del.DelegatorAddress)
		app.DistributionKeeper.Hooks().BeforeDelegationCreated(ctx, delAddr, valAddr)
		app.DistributionKeeper.Hooks().AfterDelegationModified(ctx, delAddr, valAddr)
	}

	ctx = ctx.WithBlockHeight(height)

	app.StakingKeeper.IterateRedelegations(ctx, func(_ int64, red stakingTypes.Redelegation) (stop bool) {
		for i := range red.Entries {
			red.Entries[i].CreationHeight = 0
		}
		app.StakingKeeper.SetRedelegation(ctx, red)
		return false
	})

	app.StakingKeeper.IterateUnbondingDelegations(ctx, func(_ int64, ubd stakingTypes.UnbondingDelegation) (stop bool) {
		for i := range ubd.Entries {
			ubd.Entries[i].CreationHeight = 0
		}
		app.StakingKeeper.SetUnbondingDelegation(ctx, ubd)
		return false
	})

	store := ctx.KVStore(app.keys[stakingTypes.StoreKey])
	iter := sdkTypes.KVStoreReversePrefixIterator(store, stakingTypes.ValidatorsKey)
	counter := int16(0)

	for ; iter.Valid(); iter.Next() {
		addr := sdkTypes.ValAddress(stakingTypes.AddressFromValidatorsKey(iter.Key()))
		validator, found := app.StakingKeeper.GetValidator(ctx, addr)
		if !found {
			panic("expected validator, not found")
		}

		validator.UnbondingHeight = 0
		if applyAllowedAddrs && !allowedAddrsMap[addr.String()] {
			validator.Jailed = true
		}

		app.StakingKeeper.SetValidator(ctx, validator)
		counter++
	}

	iter.Close()

	_, err := app.StakingKeeper.ApplyAndReturnValidatorSetUpdates(ctx)
	if err != nil {
		log.Fatal(err)
	}

	app.SlashingKeeper.IterateValidatorSigningInfos(
		ctx,
		func(addr sdkTypes.ConsAddress, info slashingTypes.ValidatorSigningInfo) (stop bool) {
			info.StartHeight = 0
			app.SlashingKeeper.SetValidatorSigningInfo(ctx, addr, info)
			return false
		},
	)
}

func (app *SimulationApplication) Name() string { return app.BaseApp.Name() }

func (app *SimulationApplication) BeginBlocker(ctx sdkTypes.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.moduleManager.BeginBlock(ctx, req)
}

func (app *SimulationApplication) EndBlocker(ctx sdkTypes.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.moduleManager.EndBlock(ctx, req)
}

func (app *SimulationApplication) InitChainer(ctx sdkTypes.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState GenesisState
	if err := json.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}
	app.UpgradeKeeper.SetModuleVersionMap(ctx, app.moduleManager.GetVersionMap())
	return app.moduleManager.InitGenesis(ctx, app.appCodec, genesisState)
}

func (app *SimulationApplication) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

func (app *SimulationApplication) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range simulationMake.ModuleAccountPermissions {
		modAccAddrs[authTypes.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

func (app *SimulationApplication) LegacyAmino() *codec.LegacyAmino {
	return app.legacyAmino
}

func (app *SimulationApplication) GetAppCodec() codec.Codec {
	return app.appCodec
}

func (app *SimulationApplication) InterfaceRegistry() types.InterfaceRegistry {
	return app.interfaceRegistry
}

func (app *SimulationApplication) GetKey(storeKey string) *sdkTypes.KVStoreKey {
	return app.keys[storeKey]
}

func (app *SimulationApplication) GetTKey(storeKey string) *sdkTypes.TransientStoreKey {
	return app.transientKeys[storeKey]
}

func (app *SimulationApplication) GetMemKey(storeKey string) *sdkTypes.MemoryStoreKey {
	return app.memoryKeys[storeKey]
}

func (app *SimulationApplication) GetSubspace(moduleName string) paramsTypes.Subspace {
	subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

func (app *SimulationApplication) SimulationManager() *module.SimulationManager {
	return app.simulationManager
}

func (app *SimulationApplication) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiSvr.ClientCtx
	rpc.RegisterRoutes(clientCtx, apiSvr.Router)
	authRest.RegisterTxRoutes(clientCtx, apiSvr.Router)
	authTx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	tmservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	simulationMake.ModuleBasicManagers.RegisterRESTRoutes(clientCtx, apiSvr.Router)
	simulationMake.ModuleBasicManagers.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	if apiConfig.Swagger {
		RegisterSwaggerAPI(clientCtx, apiSvr.Router)
	}
}

func (app *SimulationApplication) RegisterTxService(clientCtx client.Context) {
	authTx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
}

func (app *SimulationApplication) RegisterTendermintService(clientCtx client.Context) {
	tmservice.RegisterTendermintService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.interfaceRegistry)
}

func RegisterSwaggerAPI(ctx client.Context, rtr *mux.Router) {
	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}

	staticServer := http.FileServer(statikFS)
	rtr.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", staticServer))
}

func GetModuleAccountPermissions() map[string][]string {
	dupMaccPerms := make(map[string][]string)
	for k, v := range simulationMake.ModuleAccountPermissions {
		dupMaccPerms[k] = v
	}
	return dupMaccPerms
}
func NewSimulationApplication(logger tmLog.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool, skipUpgradeHeights map[int64]bool,
	homePath string, invCheckPeriod uint, encodingConfig simAppParams.EncodingConfig,
	appOpts serverTypes.AppOptions, baseAppOptions ...func(*baseapp.BaseApp),
) applications.SimulationApplication {
	appCodec := encodingConfig.Marshaler
	legacyAmino := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry

	bApp := baseapp.NewBaseApp("Simulation", logger, db, encodingConfig.TxConfig.TxDecoder(), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion("Simulation")
	bApp.SetInterfaceRegistry(interfaceRegistry)

	keys := sdkTypes.NewKVStoreKeys(
		authTypes.StoreKey, bankTypes.StoreKey, stakingTypes.StoreKey,
		mintTypes.StoreKey, distributionTypes.StoreKey, slashingTypes.StoreKey,
		govTypes.StoreKey, paramsTypes.StoreKey, upgradeTypes.StoreKey, feegrant.StoreKey,
		evidenceTypes.StoreKey, capabilityTypes.StoreKey,
		authzKeeper.StoreKey,
	)
	tkeys := sdkTypes.NewTransientStoreKeys(paramsTypes.TStoreKey)
	memKeys := sdkTypes.NewMemoryStoreKeys(capabilityTypes.MemStoreKey, "testingkey")

	if _, _, err := streaming.LoadStreamingServices(bApp, appOpts, appCodec, keys); err != nil {
		tmos.Exit(err.Error())
	}

	app := &SimulationApplication{
		BaseApp:           bApp,
		legacyAmino:       legacyAmino,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
		invCheckPeriod:    invCheckPeriod,
		keys:              keys,
		transientKeys:     tkeys,
		memoryKeys:        memKeys,
	}
	app.ParamsKeeper = initParamsKeeper(appCodec, legacyAmino, keys[paramsTypes.StoreKey], tkeys[paramsTypes.TStoreKey])

	bApp.SetParamStore(app.ParamsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramsKeeper.ConsensusParamsKeyTable()))

	app.CapabilityKeeper = capabilityKeeper.NewKeeper(appCodec, keys[capabilityTypes.StoreKey], memKeys[capabilityTypes.MemStoreKey])
	app.CapabilityKeeper.Seal()

	app.AccountKeeper = authKeeper.NewAccountKeeper(
		appCodec, keys[authTypes.StoreKey], app.GetSubspace(authTypes.ModuleName), authTypes.ProtoBaseAccount, simulationMake.ModuleAccountPermissions,
	)
	app.BankKeeper = bankKeeper.NewBaseKeeper(
		appCodec, keys[bankTypes.StoreKey], app.AccountKeeper, app.GetSubspace(bankTypes.ModuleName), app.ModuleAccountAddrs(),
	)
	stakingKeeper := stakingKeeper.NewKeeper(
		appCodec, keys[stakingTypes.StoreKey], app.AccountKeeper, app.BankKeeper, app.GetSubspace(stakingTypes.ModuleName),
	)
	app.MintKeeper = mintKeeper.NewKeeper(
		appCodec, keys[mintTypes.StoreKey], app.GetSubspace(mintTypes.ModuleName), &stakingKeeper,
		app.AccountKeeper, app.BankKeeper, authTypes.FeeCollectorName,
	)
	app.DistributionKeeper = distributionKeeper.NewKeeper(
		appCodec, keys[distributionTypes.StoreKey], app.GetSubspace(distributionTypes.ModuleName), app.AccountKeeper, app.BankKeeper,
		&stakingKeeper, authTypes.FeeCollectorName, app.ModuleAccountAddrs(),
	)
	app.SlashingKeeper = slashingKeeper.NewKeeper(
		appCodec, keys[slashingTypes.StoreKey], &stakingKeeper, app.GetSubspace(slashingTypes.ModuleName),
	)
	app.CrisisKeeper = crisisKeeper.NewKeeper(
		app.GetSubspace(crisisTypes.ModuleName), invCheckPeriod, app.BankKeeper, authTypes.FeeCollectorName,
	)

	app.FeeGrantKeeper = feeGrantKeeper.NewKeeper(appCodec, keys[feegrant.StoreKey], app.AccountKeeper)
	app.UpgradeKeeper = upgradeKeeper.NewKeeper(skipUpgradeHeights, keys[upgradeTypes.StoreKey], appCodec, homePath, app.BaseApp)

	app.StakingKeeper = *stakingKeeper.SetHooks(
		stakingTypes.NewMultiStakingHooks(app.DistributionKeeper.Hooks(), app.SlashingKeeper.Hooks()),
	)

	app.AuthzKeeper = authzKeeper.NewKeeper(keys[authzKeeper.StoreKey], appCodec, app.BaseApp.MsgServiceRouter())

	govRouter := govTypes.NewRouter()
	govRouter.AddRoute(govTypes.RouterKey, govTypes.ProposalHandler).
		AddRoute(paramsProposal.RouterKey, params.NewParamChangeProposalHandler(app.ParamsKeeper)).
		AddRoute(distributionTypes.RouterKey, distribution.NewCommunityPoolSpendProposalHandler(app.DistributionKeeper)).
		AddRoute(upgradeTypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(app.UpgradeKeeper))
	govKeeper := govKeeper.NewKeeper(
		appCodec, keys[govTypes.StoreKey], app.GetSubspace(govTypes.ModuleName), app.AccountKeeper, app.BankKeeper,
		&stakingKeeper, govRouter,
	)

	app.GovKeeper = *govKeeper.SetHooks(
		govTypes.NewMultiGovHooks(),
	)

	evidenceKeeper := evidenceKeeper.NewKeeper(
		appCodec, keys[evidenceTypes.StoreKey], &app.StakingKeeper, app.SlashingKeeper,
	)
	app.EvidenceKeeper = *evidenceKeeper

	skipGenesisInvariants := cast.ToBool(appOpts.Get(crisis.FlagSkipGenesisInvariants))

	app.moduleManager = module.NewManager(
		genutil.NewAppModule(
			app.AccountKeeper, app.StakingKeeper, app.BaseApp.DeliverTx,
			encodingConfig.TxConfig,
		),
		auth.NewAppModule(appCodec, app.AccountKeeper, authSimulation.RandomGenesisAccounts),
		vesting.NewAppModule(app.AccountKeeper, app.BankKeeper),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper),
		crisis.NewAppModule(&app.CrisisKeeper, skipGenesisInvariants),
		feeGrantModule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, app.interfaceRegistry),
		gov.NewAppModule(appCodec, app.GovKeeper, app.AccountKeeper, app.BankKeeper),
		mint.NewAppModule(appCodec, app.MintKeeper, app.AccountKeeper),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		distribution.NewAppModule(appCodec, app.DistributionKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		staking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
		upgrade.NewAppModule(app.UpgradeKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		params.NewAppModule(app.ParamsKeeper),
		authzModule.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
	)

	app.moduleManager.SetOrderBeginBlockers(
		upgradeTypes.ModuleName, capabilityTypes.ModuleName, mintTypes.ModuleName, distributionTypes.ModuleName, slashingTypes.ModuleName,
		evidenceTypes.ModuleName, stakingTypes.ModuleName,
		authTypes.ModuleName, bankTypes.ModuleName, govTypes.ModuleName, crisisTypes.ModuleName, genutilTypes.ModuleName,
		authz.ModuleName, feegrant.ModuleName,
		paramsTypes.ModuleName, vestingTypes.ModuleName,
	)
	app.moduleManager.SetOrderEndBlockers(
		crisisTypes.ModuleName, govTypes.ModuleName, stakingTypes.ModuleName,
		capabilityTypes.ModuleName, authTypes.ModuleName, bankTypes.ModuleName, distributionTypes.ModuleName,
		slashingTypes.ModuleName, mintTypes.ModuleName,
		genutilTypes.ModuleName, evidenceTypes.ModuleName, authz.ModuleName,
		feegrant.ModuleName,
		paramsTypes.ModuleName, upgradeTypes.ModuleName, vestingTypes.ModuleName,
	)

	app.moduleManager.SetOrderInitGenesis(
		capabilityTypes.ModuleName, authTypes.ModuleName, bankTypes.ModuleName, distributionTypes.ModuleName, stakingTypes.ModuleName,
		slashingTypes.ModuleName, govTypes.ModuleName, mintTypes.ModuleName, crisisTypes.ModuleName,
		genutilTypes.ModuleName, evidenceTypes.ModuleName, authz.ModuleName,
		feegrant.ModuleName,
		paramsTypes.ModuleName, upgradeTypes.ModuleName, vestingTypes.ModuleName,
	)

	app.moduleManager.RegisterInvariants(&app.CrisisKeeper)
	app.moduleManager.RegisterRoutes(app.Router(), app.QueryRouter(), encodingConfig.Amino)
	app.configurator = module.NewConfigurator(app.appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())
	app.moduleManager.RegisterServices(app.configurator)

	testdata.RegisterQueryServer(app.GRPCQueryRouter(), testdata.QueryImpl{})

	overrideModules := map[string]module.AppModuleSimulation{
		authTypes.ModuleName: auth.NewAppModule(app.appCodec, app.AccountKeeper, authSimulation.RandomGenesisAccounts),
	}
	app.simulationManager = module.NewSimulationManagerFromAppModules(app.moduleManager.Modules, overrideModules)

	app.simulationManager.RegisterStoreDecoders()
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)
	app.MountMemoryStores(memKeys)

	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)

	anteHandler, err := ante.NewAnteHandler(
		ante.HandlerOptions{
			AccountKeeper:   app.AccountKeeper,
			BankKeeper:      app.BankKeeper,
			SignModeHandler: encodingConfig.TxConfig.SignModeHandler(),
			FeegrantKeeper:  app.FeeGrantKeeper,
			SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
		},
	)
	if err != nil {
		panic(err)
	}

	app.SetAnteHandler(anteHandler)
	app.SetEndBlocker(app.EndBlocker)

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			tmos.Exit(err.Error())
		}
	}

	return app
}

func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey sdkTypes.StoreKey) paramsKeeper.Keeper {
	parametersKeeper := paramsKeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	parametersKeeper.Subspace(authTypes.ModuleName)
	parametersKeeper.Subspace(bankTypes.ModuleName)
	parametersKeeper.Subspace(stakingTypes.ModuleName)
	parametersKeeper.Subspace(mintTypes.ModuleName)
	parametersKeeper.Subspace(distributionTypes.ModuleName)
	parametersKeeper.Subspace(slashingTypes.ModuleName)
	parametersKeeper.Subspace(govTypes.ModuleName).WithKeyTable(govTypes.ParamKeyTable())
	parametersKeeper.Subspace(crisisTypes.ModuleName)

	return parametersKeeper
}
