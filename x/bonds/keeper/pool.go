package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply/exported"
)

// GetBondedPool returns the bonded tokens pool's module account
func (k Keeper) GetBondedPool(ctx sdk.Context) (bondedPool exported.ModuleAccountI) {
	return k.supplyKeeper.GetModuleAccount(ctx, k.BondedPoolName(ctx))
}

// GetNotBondedPool returns the not bonded tokens pool's module account
func (k Keeper) GetNotBondedPool(ctx sdk.Context) (notBondedPool exported.ModuleAccountI) {
	return k.supplyKeeper.GetModuleAccount(ctx, k.NotBondedPoolName(ctx))
}

// bondedTokensToNotBonded transfers coins from the bonded to the not bonded pool
func (k Keeper) bondedTokensToNotBonded(ctx sdk.Context, tokens sdk.Int) {
	coins := sdk.NewCoins(sdk.NewCoin(k.BondDenom(ctx), tokens))
	err := k.supplyKeeper.SendCoinsFromModuleToModule(ctx, k.BondedPoolName(ctx), k.NotBondedPoolName(ctx), coins)
	if err != nil {
		panic(err)
	}
}

// notBondedTokensToBonded transfers coins from the not bonded to the bonded pool
func (k Keeper) notBondedTokensToBonded(ctx sdk.Context, tokens sdk.Int) {
	coins := sdk.NewCoins(sdk.NewCoin(k.BondDenom(ctx), tokens))
	err := k.supplyKeeper.SendCoinsFromModuleToModule(ctx, k.NotBondedPoolName(ctx), k.BondedPoolName(ctx), coins)
	if err != nil {
		panic(err)
	}
}

// burnBondedTokens removes coins from the bonded pool module account
func (k Keeper) burnBondedTokens(ctx sdk.Context, amt sdk.Int) sdk.Error {
	if !amt.IsPositive() {
		// skip as no coins need to be burned
		return nil
	}
	coins := sdk.NewCoins(sdk.NewCoin(k.BondDenom(ctx), amt))
	return k.supplyKeeper.BurnCoins(ctx, k.BondedPoolName(ctx), coins)
}

// burnNotBondedTokens removes coins from the not bonded pool module account
func (k Keeper) burnNotBondedTokens(ctx sdk.Context, amt sdk.Int) sdk.Error {
	if !amt.IsPositive() {
		// skip as no coins need to be burned
		return nil
	}
	coins := sdk.NewCoins(sdk.NewCoin(k.BondDenom(ctx), amt))
	return k.supplyKeeper.BurnCoins(ctx, k.NotBondedPoolName(ctx), coins)
}

// TotalBondedTokens total tokens supply which is bonded
func (k Keeper) TotalBondedTokens(ctx sdk.Context) sdk.Int {
	bondedPool := k.GetBondedPool(ctx)
	return bondedPool.GetCoins().AmountOf(k.BondDenom(ctx))
}

// TokenSupply tokens from the total supply
func (k Keeper) TokenSupply(ctx sdk.Context) sdk.Int {
	return k.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf(k.BondDenom(ctx))
}

// BondedRatio the fraction of the tokens which are currently bonded
func (k Keeper) BondedRatio(ctx sdk.Context) sdk.Dec {
	bondedPool := k.GetBondedPool(ctx)

	supply := k.TokenSupply(ctx)
	if supply.IsPositive() {
		return bondedPool.GetCoins().AmountOf(k.BondDenom(ctx)).ToDec().QuoInt(supply)
	}
	return sdk.ZeroDec()
}
