package governors

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/x/bonds"
	"github.com/xar-network/xar-network/x/governors/keeper"
	"github.com/xar-network/xar-network/x/governors/types"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {

		case types.MsgBond:
			return handleMsgBond(ctx, msg, k)

		case types.MsgUnbond:
			return handleMsgUnbond(ctx, msg, k)

		default:
			errMsg := fmt.Sprintf("unrecognized governors message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// These functions assume everything has been authenticated,
// now we just perform action and save

func handleMsgBond(ctx sdk.Context, msg types.MsgBond, k keeper.Keeper) sdk.Result {
	if msg.Amount.Denom != k.BondDenom(ctx) {
		return bonds.ErrBadDenom(k.Codespace()).Result()
	}

	// NOTE: source funds are always unbonded
	err := k.Bond(ctx, msg.DepositorAddress, msg.Amount.Amount, sdk.Unbonded, true)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			bonds.EventTypeBond,
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.Amount.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, bonds.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DepositorAddress.String()),
		),
	})

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgUnbond(ctx sdk.Context, msg types.MsgUnbond, k keeper.Keeper) sdk.Result {
	err := k.ValidateUnbondAmount(
		ctx, msg.DepositorAddress, msg.Amount.Amount,
	)
	if err != nil {
		return err.Result()
	}

	if msg.Amount.Denom != k.BondDenom(ctx) {
		return bonds.ErrBadDenom(k.Codespace()).Result()
	}

	completionTime, err := k.Unbond(ctx, msg.DepositorAddress, msg.Amount.Amount)
	if err != nil {
		return err.Result()
	}

	completionTimeBz := bonds.ModuleCdc.MustMarshalBinaryLengthPrefixed(completionTime)
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			bonds.EventTypeUnbond,
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.Amount.String()),
			sdk.NewAttribute(bonds.AttributeKeyCompletionTime, completionTime.Format(time.RFC3339)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, bonds.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DepositorAddress.String()),
		),
	})

	return sdk.Result{Data: completionTimeBz, Events: ctx.EventManager().Events()}
}
