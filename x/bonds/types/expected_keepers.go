package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	supplyexported "github.com/cosmos/cosmos-sdk/x/supply/exported"
)

// AccountKeeper defines the expected account keeper (noalias)
type AccountKeeper interface {
	IterateAccounts(ctx sdk.Context, process func(authexported.Account) (stop bool))
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authexported.Account // only used for simulation
}

// SupplyKeeper defines the expected supply Keeper (noalias)
type SupplyKeeper interface {
	GetSupply(ctx sdk.Context) supplyexported.SupplyI

	GetModuleAddress(name string) sdk.AccAddress
	GetModuleAccount(ctx sdk.Context, moduleName string) supplyexported.ModuleAccountI

	SetModuleAccount(sdk.Context, supplyexported.ModuleAccountI)

	SendCoinsFromModuleToModule(ctx sdk.Context, senderPool, recipientPool string, amt sdk.Coins) sdk.Error
	UndelegateCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) sdk.Error
	DelegateCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) sdk.Error

	BurnCoins(ctx sdk.Context, name string, amt sdk.Coins) sdk.Error
}

//_______________________________________________________________________________
// Event Hooks
// These can be utilized to communicate between a keeper and another
// keeper which must take particular actions when depositors change
// state. The second keeper must implement this interface, which then the
// keeper can call.

// BondHooks event hooks for bond object (noalias)
type BondHooks interface {
	BeforeDepositCreated(ctx sdk.Context, depAddr sdk.AccAddress)        // Must be called when a deposit is created
	BeforeDepositTokensModified(ctx sdk.Context, depAddr sdk.AccAddress) // Must be called when a deposit's tokens are modified
	BeforeDepositRemoved(ctx sdk.Context, depAddr sdk.AccAddress)        // Must be called when a deposit is removed
	AfterDepositModified(ctx sdk.Context, depAddr sdk.AccAddress)
}
