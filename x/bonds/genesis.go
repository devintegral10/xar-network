package bonds

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/xar-network/xar-network/x/bonds/keeper"
	"github.com/xar-network/xar-network/x/bonds/types"
)

// InitGenesis sets the pool and parameters for the provided keeper.
// Tt also sets any deposits found in
// data.
// Returns final deposits set
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper,
	supplyKeeper types.SupplyKeeper, data types.GenesisState) (res []abci.ValidatorUpdate) {

	bondedTokens := sdk.ZeroInt()
	notBondedTokens := sdk.ZeroInt()

	// We need to pretend to be "n blocks before genesis", where "n" is the
	// validator update delay, so that e.g. slashing periods are correctly
	// initialized for the validator set e.g. with a one-block offset - the
	// first TM block is at height 1, so state updates applied from
	// genesis.json are in block 0.
	ctx = ctx.WithBlockHeight(1 - sdk.ValidatorUpdateDelay)

	keeper.SetParams(ctx, data.Params)

	for _, deposit := range data.Deposits {
		// Call the before-creation hook if not exported
		if !data.Exported {
			keeper.BeforeDepositCreated(ctx, deposit.DepositorAddress)
		}
		keeper.SetDeposit(ctx, deposit)

		// Call the after-modification hook if not exported
		if !data.Exported {
			keeper.AfterDepositModified(ctx, deposit.DepositorAddress)
		}
	}

	for _, ubd := range data.UnbondingDeposits {
		keeper.SetUnbondingDeposit(ctx, ubd)
		for _, entry := range ubd.Entries {
			keeper.InsertUBDQueue(ctx, ubd, entry.CompletionTime)
			notBondedTokens = notBondedTokens.Add(entry.Balance)
		}
	}

	bondedCoins := sdk.NewCoins(sdk.NewCoin(data.Params.BondDenom, bondedTokens))
	notBondedCoins := sdk.NewCoins(sdk.NewCoin(data.Params.BondDenom, notBondedTokens))

	// check if the unbonded and bonded pools accounts exists
	bondedPool := keeper.GetBondedPool(ctx)
	if bondedPool == nil {
		panic(fmt.Sprintf("%s module account has not been set", keeper.BondedPoolName(ctx)))
	}

	// TODO remove with genesis 2-phases refactor https://github.com/cosmos/cosmos-sdk/issues/2862
	// add coins if not provided on genesis
	if bondedPool.GetCoins().IsZero() {
		if err := bondedPool.SetCoins(bondedCoins); err != nil {
			panic(err)
		}
		supplyKeeper.SetModuleAccount(ctx, bondedPool)
	}

	notBondedPool := keeper.GetNotBondedPool(ctx)
	if notBondedPool == nil {
		panic(fmt.Sprintf("%s module account has not been set", keeper.NotBondedPoolName(ctx)))
	}

	if notBondedPool.GetCoins().IsZero() {
		if err := notBondedPool.SetCoins(notBondedCoins); err != nil {
			panic(err)
		}
		supplyKeeper.SetModuleAccount(ctx, notBondedPool)
	}

	return res
}

// ExportGenesis returns a GenesisState for a given context and keeper. The
// GenesisState will contain the pool, params, deposits found in
// the keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) types.GenesisState {
	params := keeper.GetParams(ctx)
	deposits := keeper.GetAllDeposits(ctx)
	var unbondingDeposits []types.UnbondingDeposit
	keeper.IterateUnbondingDeposits(ctx, func(_ int64, ubd types.UnbondingDeposit) (stop bool) {
		unbondingDeposits = append(unbondingDeposits, ubd)
		return false
	})

	return types.GenesisState{
		Params: params,
		GenesisStateData: types.GenesisStateData{
			Deposits:          deposits,
			UnbondingDeposits: unbondingDeposits,
		},
		Exported: true,
	}
}

// ValidateGenesis validates the provided xar genesis state to ensure the
// expected invariants holds. (i.e. params in correct bounds)
func ValidateGenesis(data types.GenesisState) error {
	err := data.Params.Validate()
	if err != nil {
		return err
	}

	return nil
}
