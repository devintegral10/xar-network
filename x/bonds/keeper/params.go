package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/xar-network/xar-network/x/bonds/types"
)

// Default parameter namespace
const (
	DefaultParamspace = types.ModuleName
)

// ParamTable for module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&types.Params{})
}

// UnbondingTime
func (k Keeper) UnbondingTime(ctx sdk.Context) (res time.Duration) {
	k.paramstore.Get(ctx, types.KeyUnbondingTime, &res)
	return
}

// MaxEntries - Maximum number of simultaneous unbonding
// deposits or redeposits (per pair/trio)
func (k Keeper) MaxEntries(ctx sdk.Context) (res uint16) {
	k.paramstore.Get(ctx, types.KeyMaxEntries, &res)
	return
}

// BondDenom - Bondable coin denomination
func (k Keeper) BondDenom(ctx sdk.Context) (res string) {
	k.paramstore.Get(ctx, types.KeyBondDenom, &res)
	return
}

// BondedPoolName - pool for bonded tokens
func (k Keeper) BondedPoolName(ctx sdk.Context) (res string) {
	k.paramstore.Get(ctx, types.KeyBondedPoolName, &res)
	return
}

// NotBondedPoolName - pool for not bonded tokens
func (k Keeper) NotBondedPoolName(ctx sdk.Context) (res string) {
	k.paramstore.Get(ctx, types.KeyNotBondedPoolName, &res)
	return
}

// Get all parameteras as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		k.UnbondingTime(ctx),
		k.MaxEntries(ctx),
		k.BondDenom(ctx),
		k.BondedPoolName(ctx),
		k.NotBondedPoolName(ctx),
	)
}

// set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}
