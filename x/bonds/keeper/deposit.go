package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/x/bonds/types"
)

// return a specific deposit
func (k Keeper) GetDeposit(ctx sdk.Context,
	depAddr sdk.AccAddress) (
	deposit types.Deposit, found bool) {

	store := ctx.KVStore(k.storeKey)
	key := types.GetDepositKey(depAddr)
	value := store.Get(key)
	if value == nil {
		return deposit, false
	}

	deposit = types.MustUnmarshalDeposit(k.cdc, value)
	return deposit, true
}

// IterateAllDeposits iterate through all of the deposits
func (k Keeper) IterateAllDeposits(ctx sdk.Context, cb func(deposit types.Deposit) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.DepositKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		deposit := types.MustUnmarshalDeposit(k.cdc, iterator.Value())
		if cb(deposit) {
			break
		}
	}
}

// GetAllDeposits returns all deposits used during genesis dump
func (k Keeper) GetAllDeposits(ctx sdk.Context) (deposits []types.Deposit) {
	k.IterateAllDeposits(ctx, func(deposit types.Deposit) bool {
		deposits = append(deposits, deposit)
		return false
	})
	return deposits
}

// set a deposit
func (k Keeper) SetDeposit(ctx sdk.Context, deposit types.Deposit) {
	store := ctx.KVStore(k.storeKey)
	b := types.MustMarshalDeposit(k.cdc, deposit)
	store.Set(types.GetDepositKey(deposit.DepositorAddress), b)
}

// remove a deposit
func (k Keeper) RemoveDeposit(ctx sdk.Context, deposit types.Deposit) {
	// TODO: Consider calling hooks outside of the store wrapper functions, it's unobvious.
	k.BeforeDepositRemoved(ctx, deposit.DepositorAddress)
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetDepositKey(deposit.DepositorAddress))
}

// return a unbonding deposit
func (k Keeper) GetUnbondingDeposit(ctx sdk.Context, depAddr sdk.AccAddress) (ubd types.UnbondingDeposit, found bool) {

	store := ctx.KVStore(k.storeKey)
	key := types.GetUBDKey(depAddr)
	value := store.Get(key)
	if value == nil {
		return ubd, false
	}

	ubd = types.MustUnmarshalUBD(k.cdc, value)
	return ubd, true
}

// iterate through all of the unbonding deposits
func (k Keeper) IterateUnbondingDeposits(ctx sdk.Context, fn func(index int64, ubd types.UnbondingDeposit) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.UnbondingDepositKey)
	defer iterator.Close()

	for i := int64(0); iterator.Valid(); iterator.Next() {
		ubd := types.MustUnmarshalUBD(k.cdc, iterator.Value())
		if stop := fn(i, ubd); stop {
			break
		}
		i++
	}
}

// HasMaxUnbondingDepositEntries - check if unbonding deposit has maximum number of entries
func (k Keeper) HasMaxUnbondingDepositEntries(ctx sdk.Context,
	depositorAddr sdk.AccAddress) bool {

	ubd, found := k.GetUnbondingDeposit(ctx, depositorAddr)
	if !found {
		return false
	}
	return len(ubd.Entries) >= int(k.MaxEntries(ctx))
}

// set the unbonding deposit and associated index
func (k Keeper) SetUnbondingDeposit(ctx sdk.Context, ubd types.UnbondingDeposit) {
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalUBD(k.cdc, ubd)
	key := types.GetUBDKey(ubd.DepositorAddress)
	store.Set(key, bz)
}

// remove the unbonding deposit object and associated index
func (k Keeper) RemoveUnbondingDeposit(ctx sdk.Context, ubd types.UnbondingDeposit) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetUBDKey(ubd.DepositorAddress)
	store.Delete(key)
}

// SetUnbondingDepositEntry adds an entry to the unbonding deposit at
// the given addresses. It creates the unbonding deposit if it does not exist
func (k Keeper) SetUnbondingDepositEntry(ctx sdk.Context,
	depositorAddr sdk.AccAddress,
	creationHeight int64, minTime time.Time, balance sdk.Int) types.UnbondingDeposit {

	ubd, found := k.GetUnbondingDeposit(ctx, depositorAddr)
	if found {
		ubd.AddEntry(creationHeight, minTime, balance)
	} else {
		ubd = types.NewUnbondingDeposit(depositorAddr, creationHeight, minTime, balance)
	}
	k.SetUnbondingDeposit(ctx, ubd)
	return ubd
}

// unbonding deposit queue timeslice operations

// gets a specific unbonding queue timeslice. A timeslice is a slice of DVPairs
// corresponding to unbonding deposits that expire at a certain time.
func (k Keeper) GetUBDQueueTimeSlice(ctx sdk.Context, timestamp time.Time) (keys []sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetUnbondingDepositTimeKey(timestamp))
	if bz == nil {
		return []sdk.AccAddress{}
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &keys)
	return keys
}

// Sets a specific unbonding queue timeslice.
func (k Keeper) SetUBDQueueTimeSlice(ctx sdk.Context, timestamp time.Time, keys []sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(keys)
	store.Set(types.GetUnbondingDepositTimeKey(timestamp), bz)
}

// Insert an unbonding deposit to the appropriate timeslice in the unbonding queue
func (k Keeper) InsertUBDQueue(ctx sdk.Context, ubd types.UnbondingDeposit,
	completionTime time.Time) {

	timeSlice := k.GetUBDQueueTimeSlice(ctx, completionTime)
	acc := ubd.DepositorAddress
	if len(timeSlice) == 0 {
		k.SetUBDQueueTimeSlice(ctx, completionTime, []sdk.AccAddress{acc})
	} else {
		timeSlice = append(timeSlice, acc)
		k.SetUBDQueueTimeSlice(ctx, completionTime, timeSlice)
	}
}

// Returns all the unbonding queue timeslices from time 0 until endTime
func (k Keeper) UBDQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator(types.UnbondingQueueKey,
		sdk.InclusiveEndBytes(types.GetUnbondingDepositTimeKey(endTime)))
}

// Returns a concatenated list of all the timeslices inclusively previous to
// currTime, and deletes the timeslices from the queue
func (k Keeper) DequeueAllMatureUBDQueue(ctx sdk.Context,
	currTime time.Time) (matureUnbonds []sdk.AccAddress) {

	store := ctx.KVStore(k.storeKey)
	// gets an iterator for all timeslices from time 0 until the current Blockheader time
	unbondingTimesliceIterator := k.UBDQueueIterator(ctx, ctx.BlockHeader().Time)
	for ; unbondingTimesliceIterator.Valid(); unbondingTimesliceIterator.Next() {
		timeslice := []sdk.AccAddress{}
		value := unbondingTimesliceIterator.Value()
		k.cdc.MustUnmarshalBinaryLengthPrefixed(value, &timeslice)
		matureUnbonds = append(matureUnbonds, timeslice...)
		store.Delete(unbondingTimesliceIterator.Key())
	}
	return matureUnbonds
}

// Perform a deposit, set/update everything necessary within the store.
// tokenSrc indicates the bond status of the incoming funds.
func (k Keeper) Bond(ctx sdk.Context, depAddr sdk.AccAddress, bondAmt sdk.Int, tokenSrc sdk.BondStatus,
	subtractAccount bool) (newShares sdk.Dec, err sdk.Error) {

	// Get or create the deposit object
	deposit, found := k.GetDeposit(ctx, depAddr)
	if !found {
		deposit = types.NewDeposit(depAddr, sdk.NewInt(0))
	}

	// call the appropriate hook if present
	if found {
		k.BeforeDepositTokensModified(ctx, depAddr)
	} else {
		k.BeforeDepositCreated(ctx, depAddr)
	}

	// if subtractAccount is true then we are
	// performing a deposit and not a redeposit, thus the source tokens are
	// all non bonded
	if subtractAccount {
		if tokenSrc == sdk.Bonded {
			panic("deposit token source cannot be bonded")
		}

		coins := sdk.NewCoins(sdk.NewCoin(k.BondDenom(ctx), bondAmt))
		err := k.supplyKeeper.DelegateCoinsFromAccountToModule(ctx, deposit.DepositorAddress, k.BondedPoolName(ctx), coins)
		if err != nil {
			return sdk.Dec{}, err
		}
	} else {
		// potentially transfer tokens between pools, if
		if tokenSrc == sdk.Unbonded || tokenSrc == sdk.Unbonding {
			// transfer pools
			k.notBondedTokensToBonded(ctx, bondAmt)
		}
	}

	// Update deposit
	deposit.Tokens = deposit.Tokens.Add(bondAmt)
	k.SetDeposit(ctx, deposit)

	// Call the after-modification hook
	k.AfterDepositModified(ctx, deposit.DepositorAddress)

	return newShares, nil
}

// unbond a particular deposit and perform associated store operations
func (k Keeper) unbond(ctx sdk.Context, depAddr sdk.AccAddress, tokens sdk.Int) (err sdk.Error) {

	// check if a deposit object exists in the store
	deposit, found := k.GetDeposit(ctx, depAddr)
	if !found {
		return types.ErrNoDepositorForAddress(k.Codespace())
	}

	// call the before-deposit-modified hook
	k.BeforeDepositTokensModified(ctx, depAddr)

	// ensure that we have enough tokens to remove
	if deposit.Tokens.LT(tokens) {
		return types.ErrNotEnoughDepositTokens(k.Codespace(), deposit.Tokens.String())
	}

	// subtract tokens from deposit
	deposit.Tokens = deposit.Tokens.Sub(tokens)

	// remove the deposit
	if deposit.Tokens.IsZero() {
		k.RemoveDeposit(ctx, deposit)
	} else {
		k.SetDeposit(ctx, deposit)
		// call the after deposit modification hook
		k.AfterDepositModified(ctx, deposit.DepositorAddress)
	}

	return nil
}

// Unbond unbonds an amount of depositor tokens from a given depositor by creating
// an unbonding object and inserting it into the unbonding queue which will be
// processed during the EndBlocker.
func (k Keeper) Unbond(
	ctx sdk.Context, depAddr sdk.AccAddress, tokens sdk.Int,
) (time.Time, sdk.Error) {

	if k.HasMaxUnbondingDepositEntries(ctx, depAddr) {
		return time.Time{}, types.ErrMaxUnbondingDepositEntries(k.Codespace())
	}

	err := k.unbond(ctx, depAddr, tokens)
	if err != nil {
		return time.Time{}, err
	}

	// transfer the tokens to the not bonded pool
	k.bondedTokensToNotBonded(ctx, tokens)

	completionTime := ctx.BlockHeader().Time.Add(k.UnbondingTime(ctx))
	ubd := k.SetUnbondingDepositEntry(ctx, depAddr, ctx.BlockHeight(), completionTime, tokens)
	k.InsertUBDQueue(ctx, ubd, completionTime)

	return completionTime, nil
}

// CompleteUnbonding completes the unbonding of all mature entries in the
// retrieved unbonding deposit object.
func (k Keeper) CompleteUnbonding(ctx sdk.Context, depAddr sdk.AccAddress) sdk.Error {

	ubd, found := k.GetUnbondingDeposit(ctx, depAddr)
	if !found {
		return types.ErrNoUnbondingDeposit(k.Codespace())
	}

	ctxTime := ctx.BlockHeader().Time

	// loop through all the entries and complete unbonding mature entries
	for i := 0; i < len(ubd.Entries); i++ {
		entry := ubd.Entries[i]
		if entry.IsMature(ctxTime) {
			ubd.RemoveEntry(int64(i))
			i--

			// track undeposit only when remaining or truncated tokens are non-zero
			if !entry.Balance.IsZero() {
				amt := sdk.NewCoins(sdk.NewCoin(k.GetParams(ctx).BondDenom, entry.Balance))
				err := k.supplyKeeper.UndelegateCoinsFromModuleToAccount(ctx, k.NotBondedPoolName(ctx), ubd.DepositorAddress, amt)
				if err != nil {
					return err
				}
			}
		}
	}

	// set the unbonding deposit or remove it if there are no more entries
	if len(ubd.Entries) == 0 {
		k.RemoveUnbondingDeposit(ctx, ubd)
	} else {
		k.SetUnbondingDeposit(ctx, ubd)
	}

	return nil
}

// ValidateUnbondAmount validates that a given unbond or deposit amount is
// valied based on upon the deposited tokens. If the amount is valid, the total
// amount of respective tokens is returned, otherwise an error is returned.
func (k Keeper) ValidateUnbondAmount(
	ctx sdk.Context, depAddr sdk.AccAddress, amt sdk.Int,
) (err sdk.Error) {
	dep, found := k.GetDeposit(ctx, depAddr)
	if !found {
		return types.ErrNoDeposit(k.Codespace())
	}

	if amt.GT(dep.Tokens) {
		return types.ErrNotEnoughDepositTokens(k.Codespace(), dep.Tokens.String())
	}

	return nil
}
