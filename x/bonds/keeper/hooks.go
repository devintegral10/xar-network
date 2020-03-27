package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/x/bonds/types"
)

// Implements BondHooks interface
var _ types.BondHooks = Keeper{}

// BeforeDepositCreated - call hook if registered
func (k Keeper) BeforeDepositCreated(ctx sdk.Context, depAddr sdk.AccAddress) {
	if k.hooks != nil {
		k.hooks.BeforeDepositCreated(ctx, depAddr)
	}
}

// BeforeDepositTokensModified - call hook if registered
func (k Keeper) BeforeDepositTokensModified(ctx sdk.Context, depAddr sdk.AccAddress) {
	if k.hooks != nil {
		k.hooks.BeforeDepositTokensModified(ctx, depAddr)
	}
}

// BeforeDepositRemoved - call hook if registered
func (k Keeper) BeforeDepositRemoved(ctx sdk.Context, depAddr sdk.AccAddress) {
	if k.hooks != nil {
		k.hooks.BeforeDepositRemoved(ctx, depAddr)
	}
}

// AfterDepositModified - call hook if registered
func (k Keeper) AfterDepositModified(ctx sdk.Context, depAddr sdk.AccAddress) {
	if k.hooks != nil {
		k.hooks.AfterDepositModified(ctx, depAddr)
	}
}
