package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/xar-network/xar-network/x/bonds"
	"github.com/xar-network/xar-network/x/governors/types"
)

type (
	Keeper struct {
		bonds.Keeper
	}
)

// NewKeeper creates a new Keeper instance
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, supplyKeeper types.SupplyKeeper,
	paramstore params.Subspace, codespace sdk.CodespaceType) Keeper {

	return Keeper{
		Keeper: bonds.NewKeeper(cdc, key, supplyKeeper, paramstore, codespace),
	}
}
