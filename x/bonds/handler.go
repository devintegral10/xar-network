package bonds

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/xar-network/xar-network/x/bonds/keeper"
	"github.com/xar-network/xar-network/x/bonds/types"
)

// Called every block, update deposits set
func EndBlocker(ctx sdk.Context, k keeper.Keeper) []abci.ValidatorUpdate {
	// Remove all mature unbonding deposits from the ubd queue.
	accs := k.DequeueAllMatureUBDQueue(ctx, ctx.BlockHeader().Time)
	for _, acc := range accs {
		err := k.CompleteUnbonding(ctx, acc)
		if err != nil {
			continue
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeCompleteUnbonding,
				sdk.NewAttribute(types.AttributeKeyDepostior, acc.String()),
			),
		)
	}

	return nil
}
