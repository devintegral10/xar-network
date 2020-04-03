package keeper

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/x/bonds/types"

	"github.com/stretchr/testify/require"
)

// tests GetDeposit, SetDeposit, RemoveDeposit
func TestDeposit(t *testing.T) {
	ctx, accountKeeper, keeper, _ := CreateTestInput(t, false, sdk.NewInt(20))

	var deposits [3]types.Deposit
	amts := []sdk.Int{sdk.NewInt(9), sdk.NewInt(8), sdk.NewInt(1)}
	for i, amt := range amts {
		// Balance before bonding
		coinsBefore := accountKeeper.GetAccount(ctx, Addrs[i]).GetCoins()

		// try zero deposit
		err := keeper.Bond(ctx, Addrs[i], sdk.ZeroInt(), sdk.Unbonded, true)
		require.Error(t, err, types.ErrZeroAmount(keeper.Codespace()))

		// Create Depositor
		err = keeper.Bond(ctx, Addrs[i], amt, sdk.Unbonded, true)
		require.NoError(t, err)
		deposits[i], _ = keeper.GetDeposit(ctx, Addrs[i])

		// test record
		require.Equal(t, amt, deposits[i].Tokens)
		require.Equal(t, Addrs[i], deposits[i].DepositorAddress)

		// Balance after bonding
		coinsAfter := accountKeeper.GetAccount(ctx, Addrs[i]).GetCoins()

		// test balances
		require.Equal(t, 1, coinsBefore.Len())
		require.Equal(t, 1, coinsAfter.Len())
		require.Equal(t, coinsBefore.AmountOf(Denom).Sub(amt).String(), coinsAfter.AmountOf(Denom).String())
	}

	// Increase deposit
	var icnreasedDeposits [3]types.Deposit
	amts = []sdk.Int{sdk.NewInt(3), sdk.NewInt(2), sdk.NewInt(1)}
	for i, amt := range amts {
		// Balance before bonding
		coinsBefore := accountKeeper.GetAccount(ctx, Addrs[i]).GetCoins()

		// try zero deposit
		err := keeper.Bond(ctx, Addrs[i], sdk.ZeroInt(), sdk.Unbonded, true)
		require.Error(t, err, types.ErrZeroAmount(keeper.Codespace()))

		// Create Depositor
		err = keeper.Bond(ctx, Addrs[i], amt, sdk.Unbonded, true)
		require.NoError(t, err)
		icnreasedDeposits[i], _ = keeper.GetDeposit(ctx, Addrs[i])

		// test record
		require.Equal(t, deposits[i].Tokens.Add(amt), icnreasedDeposits[i].Tokens)
		require.Equal(t, Addrs[i], icnreasedDeposits[i].DepositorAddress)

		// Balance after bonding
		coinsAfter := accountKeeper.GetAccount(ctx, Addrs[i]).GetCoins()

		// test balances
		require.Equal(t, 1, coinsBefore.Len())
		require.Equal(t, 1, coinsAfter.Len())
		require.Equal(t, coinsBefore.AmountOf(Denom).Sub(amt).String(), coinsAfter.AmountOf(Denom).String())
	}

	// delete a record
	keeper.RemoveDeposit(ctx, icnreasedDeposits[0])
	_, found := keeper.GetUnbondingDeposit(ctx, Addrs[0])
	require.False(t, found)
}

// tests Get/Set/Remove UnbondingDeposit
func TestUnbondingDeposit(t *testing.T) {
	ctx, _, keeper, _ := CreateTestInput(t, false, sdk.NewInt(0))

	ubd := types.NewUnbondingDeposit(Addrs[0], 0,
		time.Unix(0, 0), sdk.NewInt(5))

	// set and retrieve a record
	keeper.SetUnbondingDeposit(ctx, ubd)
	resUnbond, found := keeper.GetUnbondingDeposit(ctx, Addrs[0])
	require.True(t, found)
	require.True(t, ubd.Equal(resUnbond))

	// modify a records, save, and retrieve
	ubd.Entries[0].Balance = sdk.NewInt(21)
	keeper.SetUnbondingDeposit(ctx, ubd)

	resUnbond, found = keeper.GetUnbondingDeposit(ctx, Addrs[0])
	require.True(t, found)

	require.True(t, found)
	require.True(t, ubd.Equal(resUnbond))

	// delete a record
	keeper.RemoveUnbondingDeposit(ctx, ubd)
	_, found = keeper.GetUnbondingDeposit(ctx, Addrs[0])
	require.False(t, found)
}

func TestUnbondDeposit(t *testing.T) {
	ctx, _, keeper, _ := CreateTestInput(t, false, sdk.NewInt(0))

	startTokens := sdk.NewInt(20)

	notBondedPool := keeper.GetNotBondedPool(ctx)
	err := notBondedPool.SetCoins(sdk.NewCoins(sdk.NewCoin(Denom, startTokens)))
	require.NoError(t, err)
	keeper.supplyKeeper.SetModuleAccount(ctx, notBondedPool)

	// create a deposit
	deposit := types.NewDeposit(Addrs[0], startTokens)
	keeper.SetDeposit(ctx, deposit)

	unbondTokens := sdk.NewInt(4)
	err = keeper.unbond(ctx, Addrs[0], unbondTokens)
	require.NoError(t, err)

	deposit, found := keeper.GetDeposit(ctx, Addrs[0])
	require.True(t, found)

	remainingTokens := startTokens.Sub(unbondTokens)
	require.Equal(t, remainingTokens, deposit.Tokens)
}

func TestUnbondingDepositsMaxEntries(t *testing.T) {
	startTokens := sdk.NewInt(1000)
	ctx, _, keeper, _ := CreateTestInput(t, false, startTokens)

	notBondedPool := keeper.GetNotBondedPool(ctx)
	err := notBondedPool.SetCoins(sdk.NewCoins(sdk.NewCoin(Denom, startTokens)))
	require.NoError(t, err)
	keeper.supplyKeeper.SetModuleAccount(ctx, notBondedPool)

	// create a deposit
	deposit := types.NewDeposit(Addrs[0], startTokens)
	err = keeper.Bond(ctx, Addrs[0], startTokens, sdk.Unbonded, true)
	require.NoError(t, err)
	require.True(sdk.IntEq(t, startTokens, deposit.Tokens))

	maxEntries := keeper.MaxEntries(ctx)

	oldBonded := keeper.GetBondedPool(ctx).GetCoins().AmountOf(Denom)
	oldNotBonded := keeper.GetNotBondedPool(ctx).GetCoins().AmountOf(Denom)

	// should all pass
	var completionTime time.Time
	for i := uint16(0); i < maxEntries; i++ {
		var err error
		completionTime, err = keeper.Unbond(ctx, Addrs[0], sdk.NewInt(1))
		require.NoError(t, err)
	}

	bondedPool := keeper.GetBondedPool(ctx)
	notBondedPool = keeper.GetNotBondedPool(ctx)
	require.True(sdk.IntEq(t, bondedPool.GetCoins().AmountOf(Denom), oldBonded.SubRaw(int64(maxEntries))))
	require.True(sdk.IntEq(t, notBondedPool.GetCoins().AmountOf(Denom), oldNotBonded.AddRaw(int64(maxEntries))))

	oldBonded = bondedPool.GetCoins().AmountOf(Denom)
	oldNotBonded = notBondedPool.GetCoins().AmountOf(Denom)

	// an additional unbond should fail due to max entries
	_, err = keeper.Unbond(ctx, Addrs[0], sdk.NewInt(1))
	require.Error(t, err)

	bondedPool = keeper.GetBondedPool(ctx)
	notBondedPool = keeper.GetNotBondedPool(ctx)
	require.True(sdk.IntEq(t, bondedPool.GetCoins().AmountOf(Denom), oldBonded))
	require.True(sdk.IntEq(t, notBondedPool.GetCoins().AmountOf(Denom), oldNotBonded))

	// mature unbonding deposits
	ctx = ctx.WithBlockTime(completionTime)
	err = keeper.CompleteUnbonding(ctx, Addrs[0])
	require.NoError(t, err)

	bondedPool = keeper.GetBondedPool(ctx)
	notBondedPool = keeper.GetNotBondedPool(ctx)
	require.True(sdk.IntEq(t, bondedPool.GetCoins().AmountOf(Denom), oldBonded))
	require.True(sdk.IntEq(t, notBondedPool.GetCoins().AmountOf(Denom), oldNotBonded.SubRaw(int64(maxEntries))))

	oldNotBonded = notBondedPool.GetCoins().AmountOf(Denom)

	// unbonding should work again
	_, err = keeper.Unbond(ctx, Addrs[0], sdk.NewInt(1))
	require.NoError(t, err)

	bondedPool = keeper.GetBondedPool(ctx)

	notBondedPool = keeper.GetNotBondedPool(ctx)
	require.True(sdk.IntEq(t, bondedPool.GetCoins().AmountOf(Denom), oldBonded.SubRaw(1)))
	require.True(sdk.IntEq(t, notBondedPool.GetCoins().AmountOf(Denom), oldNotBonded.AddRaw(1)))
}
