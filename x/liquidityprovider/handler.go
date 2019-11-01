package liquidityprovider

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/x/liquidityprovider/keeper"
	types "github.com/xar-network/xar-network/x/liquidityprovider/types"
)

// TODO Accept Keeper argument
func newHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.MsgMintTokens:
			return k.MintTokensFromCredit(ctx, msg.LiquidityProvider, msg.Amount)
		case types.MsgBurnTokens:
			return k.BurnTokensFromBalance(ctx, msg.LiquidityProvider, msg.Amount)
		default:
			errMsg := fmt.Sprintf("Unrecognized issuance Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}
