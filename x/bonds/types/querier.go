package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// query endpoints supported by the Querier
const (
	QueryDepositorUnbondingDeposits = "depositorUnbondingDeposits"
	QueryDeposit                    = "deposit"
	QueryUnbondingDeposit           = "unbondingDeposit"
	QueryPool                       = "pool"
	QueryParameters                 = "parameters"
)

// defines the params for the following queries:
// - 'custom/module/depositorUnbondingDeposits'
type QueryDepositorParams struct {
	DepositorAddr sdk.AccAddress
}

func NewQueryDepositorParams(depositorAddr sdk.AccAddress) QueryDepositorParams {
	return QueryDepositorParams{
		DepositorAddr: depositorAddr,
	}
}

// defines the params for the following queries:
// - 'custom/module/deposit'
// - 'custom/module/unbondingDeposit'
type QueryBondsParams struct {
	DepositorAddr sdk.AccAddress
}

func NewQueryBondsParams(depositorAddr sdk.AccAddress) QueryBondsParams {
	return QueryBondsParams{
		DepositorAddr: depositorAddr,
	}
}
