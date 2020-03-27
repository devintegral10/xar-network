package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/x/bonds/types"
)

// return all unbonding-deposits for a depositor
func (k Keeper) GetAllUnbondingDeposits(ctx sdk.Context, depositor sdk.AccAddress) []types.UnbondingDeposit {
	unbondingDeposits := make([]types.UnbondingDeposit, 0)

	store := ctx.KVStore(k.storeKey)
	depositorPrefixKey := types.GetUBDKey(depositor)
	iterator := sdk.KVStorePrefixIterator(store, depositorPrefixKey) // smallest to largest
	defer iterator.Close()

	for i := 0; iterator.Valid(); iterator.Next() {
		unbondingDeposit := types.MustUnmarshalUBD(k.cdc, iterator.Value())
		unbondingDeposits = append(unbondingDeposits, unbondingDeposit)
		i++
	}

	return unbondingDeposits
}
