package keeper

import (
	"fmt"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/x/bonds/types"
)

// creates a querier for REST endpoints
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.QueryDeposit:
			return queryDeposit(ctx, req, k)
		case types.QueryUnbondingDeposit:
			return queryUnbondingDeposit(ctx, req, k)
		case types.QueryPool:
			return queryPool(ctx, k)
		case types.QueryParameters:
			return queryParameters(ctx, k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown query endpoint")
		}
	}
}

func queryDeposit(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryBondsParams

	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	deposit, found := k.GetDeposit(ctx, params.DepositorAddr)
	if !found {
		return nil, types.ErrNoDeposit(types.DefaultCodespace)
	}

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, deposit)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to marshal result to JSON", err.Error()))
	}

	return res, nil
}

func queryUnbondingDeposit(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryBondsParams

	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	unbond, found := k.GetUnbondingDeposit(ctx, params.DepositorAddr)
	if !found {
		return nil, types.ErrNoUnbondingDeposit(types.DefaultCodespace)
	}

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, unbond)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return res, nil
}

func queryPool(ctx sdk.Context, k Keeper) ([]byte, sdk.Error) {
	bondDenom := k.BondDenom(ctx)
	bondedPool := k.GetBondedPool(ctx)
	notBondedPool := k.GetNotBondedPool(ctx)
	if bondedPool == nil || notBondedPool == nil {
		return nil, sdk.ErrInternal("pool accounts haven't been set")
	}

	pool := types.NewPool(
		notBondedPool.GetCoins().AmountOf(bondDenom),
		bondedPool.GetCoins().AmountOf(bondDenom),
	)

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, pool)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return res, nil
}

func queryParameters(ctx sdk.Context, k Keeper) ([]byte, sdk.Error) {
	params := k.GetParams(ctx)

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, params)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return res, nil
}
