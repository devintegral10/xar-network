package governors

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	bondskeeper "github.com/xar-network/xar-network/x/bonds/keeper"
	"github.com/xar-network/xar-network/x/governors/keeper"
	"github.com/xar-network/xar-network/x/governors/types"
)

// Hogpodge of all sorts of input required for testing.
// `initTokens` is an initial amount of tokens.
// If `initTokens` is 0, no addrs get created.
func CreateTestInput(t *testing.T, isCheckTx bool, initTokens sdk.Int) (sdk.Context, auth.AccountKeeper, keeper.Keeper, types.SupplyKeeper, []sdk.AccAddress) {
	ctx, a, bondsKeeper, s := bondskeeper.CreateTestInput(t, isCheckTx, initTokens)

	return ctx, a, keeper.Keeper{bondsKeeper}, s, bondskeeper.Addrs
}
