package governors

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/xar-network/xar-network/x/bonds"
	"github.com/xar-network/xar-network/x/governors/keeper"
	"github.com/xar-network/xar-network/x/governors/types"
)

// InitGenesis sets the pool and parameters for the provided keeper.
// Tt also sets any deposits found in
// data.
// Returns final deposits set
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, accountKeeper types.AccountKeeper,
	supplyKeeper types.SupplyKeeper, data types.GenesisState) (res []abci.ValidatorUpdate) {

	return bonds.InitGenesis(ctx, keeper.Keeper, accountKeeper, supplyKeeper, data.BondsGenesis())
}

// ExportGenesis returns a GenesisState for a given context and keeper. The
// GenesisState will contain the pool, params, deposits found in
// the keeper.
func ExportGenesis(ctx sdk.Context, bondsKeeper keeper.Keeper) types.GenesisState {
	g := bonds.ExportGenesis(ctx, bondsKeeper.Keeper)
	return types.GenesisState{
		Params: types.Params{
			Bonds: g.Params,
		},
		Bonds:    g.GenesisStateData,
		Exported: g.Exported,
	}
}

// ValidateGenesis validates the provided genesis state to ensure the
// expected invariants holds. (i.e. params in correct bounds)
func ValidateGenesis(data types.GenesisState) error {
	err := data.Params.Validate()
	if err != nil {
		return err
	}

	return nil
}
