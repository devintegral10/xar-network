package governors

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	keep "github.com/xar-network/xar-network/x/governors/keeper"
	"github.com/xar-network/xar-network/x/governors/types"
)

func TestInvalidMsg(t *testing.T) {
	k := keep.Keeper{}
	h := NewHandler(k)

	res := h(sdk.NewContext(nil, abci.Header{}, false, nil), sdk.NewTestMsg())
	require.False(t, res.IsOK())
	require.True(t, strings.Contains(res.Log, "unrecognized governors message type"))
}

func TestInvalidCoinDenom(t *testing.T) {
	tokens := sdk.NewInt(100)
	ctx, _, k, _, addrs := CreateTestInput(t, false, tokens)

	invalidCoin := sdk.NewCoin("churros", tokens)
	validCoin := sdk.NewCoin(sdk.DefaultBondDenom, tokens)
	oneCoin := sdk.NewCoin(sdk.DefaultBondDenom, sdk.OneInt())

	msgDelegate := types.NewMsgBond(addrs[0], invalidCoin)
	got := handleMsgBond(ctx, msgDelegate, k)
	require.False(t, got.IsOK())
	msgDelegate = types.NewMsgBond(addrs[0], validCoin)
	got = handleMsgBond(ctx, msgDelegate, k)
	require.True(t, got.IsOK())

	msgUndelegate := types.NewMsgUnbond(addrs[0], invalidCoin)
	got = handleMsgUnbond(ctx, msgUndelegate, k)
	require.False(t, got.IsOK())
	msgUndelegate = types.NewMsgUnbond(addrs[0], oneCoin)
	got = handleMsgUnbond(ctx, msgUndelegate, k)
	require.True(t, got.IsOK())
}
