package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/xar-network/xar-network/x/bonds/types"
)

// keeper of the store
type Keeper struct {
	storeKey     sdk.StoreKey
	cdc          *codec.Codec
	supplyKeeper types.SupplyKeeper
	hooks        types.BondHooks
	paramstore   params.Subspace

	// codespace
	codespace sdk.CodespaceType
}

// NewKeeper creates a new Keeper instance
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, supplyKeeper types.SupplyKeeper,
	paramstore params.Subspace, codespace sdk.CodespaceType) Keeper {

	return Keeper{
		storeKey:     key,
		cdc:          cdc,
		supplyKeeper: supplyKeeper,
		paramstore:   paramstore.WithKeyTable(ParamKeyTable()),
		hooks:        nil,
		codespace:    codespace,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// Set the bond hooks
func (k *Keeper) SetHooks(sh types.BondHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set bond hooks twice")
	}
	k.hooks = sh
	return k
}

// return the codespace
func (k Keeper) Codespace() sdk.CodespaceType {
	return k.codespace
}
