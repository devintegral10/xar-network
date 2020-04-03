package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/x/bonds/types"
)

var (
	addrAcc1, addrAcc2 = Addrs[0], Addrs[1]
)

func TestNewQuerier(t *testing.T) {
	cdc := codec.New()
	ctx, _, keeper, _ := CreateTestInput(t, false, sdk.NewInt(1000))
	// Create Depositors
	amts := []sdk.Int{sdk.NewInt(9), sdk.NewInt(8)}
	for i, amt := range amts {
		err := keeper.Bond(ctx, Addrs[i], amt, sdk.Unbonded, true)
		require.NoError(t, err)
	}

	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}

	querier := NewQuerier(keeper)

	bz, err := querier(ctx, []string{"other"}, query)
	require.NotNil(t, err)
	require.Nil(t, bz)

	_, err = querier(ctx, []string{"pool"}, query)
	require.Nil(t, err)

	_, err = querier(ctx, []string{"parameters"}, query)
	require.Nil(t, err)

	queryDepParams := types.NewQueryDepositorParams(addrAcc1)
	bz, errRes := cdc.MarshalJSON(queryDepParams)
	require.Nil(t, errRes)

	// Query deposit
	query = abci.RequestQuery{
		Path: "/custom/bonds/deposit",
		Data: bz,
	}

	deposit1, found := keeper.GetDeposit(ctx, addrAcc1)
	require.True(t, found)

	res, err := queryDeposit(ctx, query, keeper)
	require.Nil(t, err)

	var depositRes1 types.Deposit
	errRes = cdc.UnmarshalJSON(res, &depositRes1)
	require.Nil(t, errRes)

	require.Equal(t, deposit1, depositRes1)

	// unbond
	unbondTokens := sdk.NewInt(1)
	_, err = keeper.Unbond(ctx, addrAcc2, unbondTokens)
	require.NoError(t, err)

	// Query unbondingDeposit

	queryDepParams = types.NewQueryDepositorParams(addrAcc2)
	bz, errRes = cdc.MarshalJSON(queryDepParams)
	require.Nil(t, errRes)

	query = abci.RequestQuery{
		Path: "/custom/bonds/unbondingDeposit",
		Data: bz,
	}

	res, err = queryUnbondingDeposit(ctx, query, keeper)
	require.NoError(t, err)

	var unbondingDepositRes2 types.UnbondingDeposit
	errRes = cdc.UnmarshalJSON(res, &unbondingDepositRes2)
	require.Nil(t, errRes)

	deposit2, found := keeper.GetDeposit(ctx, addrAcc2)
	require.True(t, found)

	require.Equal(t, deposit2.DepositorAddress, unbondingDepositRes2.DepositorAddress)
	require.Equal(t, 1, len(unbondingDepositRes2.Entries))
	require.Equal(t, unbondTokens, unbondingDepositRes2.Entries[0].Balance)
	require.Equal(t, unbondTokens, unbondingDepositRes2.Entries[0].InitialBalance)
}

func TestQueryParametersPool(t *testing.T) {
	cdc := codec.New()
	ctx, _, keeper, _ := CreateTestInput(t, false, sdk.NewInt(1000))

	res, err := queryParameters(ctx, keeper)
	require.Nil(t, err)

	var params types.Params
	errRes := cdc.UnmarshalJSON(res, &params)
	require.Nil(t, errRes)
	require.Equal(t, keeper.GetParams(ctx), params)

	res, err = queryPool(ctx, keeper)
	require.Nil(t, err)

	var pool types.Pool
	bondedPool := keeper.GetBondedPool(ctx)
	notBondedPool := keeper.GetNotBondedPool(ctx)
	errRes = cdc.UnmarshalJSON(res, &pool)
	require.Nil(t, errRes)
	require.Equal(t, bondedPool.GetCoins().AmountOf(Denom), pool.BondedTokens)
	require.Equal(t, notBondedPool.GetCoins().AmountOf(Denom), pool.NotBondedTokens)
}
