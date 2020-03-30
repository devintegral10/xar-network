package governors

import (
	"encoding/json"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	sim "github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/xar-network/xar-network/x/bonds"
	"github.com/xar-network/xar-network/x/governors/client/cli"
	"github.com/xar-network/xar-network/x/governors/client/rest"
	"github.com/xar-network/xar-network/x/governors/keeper"
	"github.com/xar-network/xar-network/x/governors/simulation"
	"github.com/xar-network/xar-network/x/governors/types"
)

var (
	_ module.AppModule           = AppModule{}
	_ module.AppModuleBasic      = AppModuleBasic{}
	_ module.AppModuleSimulation = AppModuleSimulation{}
)

// AppModuleBasic defines the basic application module used by the module.
type AppModuleBasic struct{}

var _ module.AppModuleBasic = AppModuleBasic{}

// Name returns the module's name.
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterCodec registers the module's types for the given codec.
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	types.RegisterCodec(cdc)
}

// DefaultGenesis returns default genesis state as raw bytes for the module
// module.
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return types.ModuleCdc.MustMarshalJSON(types.DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the module.
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var data types.GenesisState
	err := types.ModuleCdc.UnmarshalJSON(bz, &data)
	if err != nil {
		return err
	}
	return ValidateGenesis(data)
}

// RegisterRESTRoutes registers the REST routes for the module.
func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	rest.RegisterRoutes(ctx, rtr)
}

// GetTxCmd returns the root tx command for the module.
func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetTxCmd(types.StoreKey, cdc)
}

// GetQueryCmd returns no root query command for the module.
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetQueryCmd(types.StoreKey, cdc)
}

//____________________________________________________________________________

// AppModuleSimulation defines the module simulation functions used by the module.
type AppModuleSimulation struct{}

// RegisterStoreDecoder registers a decoder for module's types
func (AppModuleSimulation) RegisterStoreDecoder(sdr sdk.StoreDecoderRegistry) {
	sdr[types.StoreKey] = simulation.DecodeStore
}

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModuleSimulation) GenerateGenesisState(simState *module.SimulationState) {
	simulation.RandomizedGenState(simState)
}

// RandomizedParams creates randomized param changes for the simulator.
func (AppModuleSimulation) RandomizedParams(r *rand.Rand) []sim.ParamChange {
	return simulation.ParamChanges(r)
}

//____________________________________________________________________________

// AppModule implements an application module for the governors module.
type AppModule struct {
	AppModuleBasic
	AppModuleSimulation

	keeper        keeper.Keeper
	accountKeeper types.AccountKeeper
	supplyKeeper  types.SupplyKeeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(keeper keeper.Keeper, accountKeeper types.AccountKeeper, supplyKeeper types.SupplyKeeper) AppModule {

	return AppModule{
		AppModuleBasic:      AppModuleBasic{},
		AppModuleSimulation: AppModuleSimulation{},
		keeper:              keeper,
		accountKeeper:       accountKeeper,
		supplyKeeper:        supplyKeeper,
	}
}

// Name returns the governors module's name.
func (AppModule) Name() string {
	return types.ModuleName
}

// RegisterInvariants registers the module invariants.
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {}

// Route returns the message routing key for the module.
func (AppModule) Route() string {
	return types.RouterKey
}

// NewHandler returns an sdk.Handler for the module.
func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.keeper)
}

// QuerierRoute returns the governors module's querier route name.
func (AppModule) QuerierRoute() string {
	return types.QuerierRoute
}

// NewQuerierHandler returns the governors module sdk.Querier.
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return keeper.NewQuerier(am.keeper.Keeper)
}

// InitGenesis performs genesis initialization for the governors module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState
	types.ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	return InitGenesis(ctx, am.keeper, am.supplyKeeper, genesisState)
}

// ExportGenesis returns the exported genesis state as raw bytes for the xar
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return types.ModuleCdc.MustMarshalJSON(gs)
}

// BeginBlock returns the begin blocker for the governors module.
func (AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}

// EndBlock returns the end blocker for the governors module. It returns no validator
// updates.
func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return bonds.EndBlocker(ctx, am.keeper.Keeper)
}
