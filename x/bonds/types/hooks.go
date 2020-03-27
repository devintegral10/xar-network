package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// combine multiple hooks, all hook functions are run in array sequence
type MultiBondHooks []BondHooks

func NewMultiBondHooks(hooks ...BondHooks) MultiBondHooks {
	return hooks
}

// nolint
func (h MultiBondHooks) BeforeDepositCreated(ctx sdk.Context, depAddr sdk.AccAddress) {
	for i := range h {
		h[i].BeforeDepositCreated(ctx, depAddr)
	}
}
func (h MultiBondHooks) BeforeDepositTokensModified(ctx sdk.Context, depAddr sdk.AccAddress) {
	for i := range h {
		h[i].BeforeDepositTokensModified(ctx, depAddr)
	}
}
func (h MultiBondHooks) BeforeDepositRemoved(ctx sdk.Context, depAddr sdk.AccAddress) {
	for i := range h {
		h[i].BeforeDepositRemoved(ctx, depAddr)
	}
}
func (h MultiBondHooks) AfterDepositModified(ctx sdk.Context, depAddr sdk.AccAddress) {
	for i := range h {
		h[i].AfterDepositModified(ctx, depAddr)
	}
}
