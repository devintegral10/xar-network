/**

Baseline from Kava Cosmos Module

**/

package csdt

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/x/csdt/internal/keeper"
	"github.com/xar-network/xar-network/x/csdt/internal/types"
)

type CoinU struct {
	Denom  string   `json:"denom"`
	Amount sdk.Uint `json:"amount"`
}

type CoinUs []CoinU

// GenesisState is the state that must be provided at genesis.
type GenesisState struct {
	Params        types.Params `json:"params"`
	GlobalDebt    sdk.Int      `json:"global_debt"`
	CSDTs         types.CSDTs  `json:"csdts" yaml:"csdts"`
	TotalBorrows  CoinUs       `json:"total_borrows"`
	TotalSupplies CoinUs       `json:"total_supplies"`
	TotalReserves CoinUs       `json:"total_reserves"`
	// don't need to setup CollateralStates as they are created as needed
}

var (
	DefaultBaseRatePerYear   = sdk.NewUint(10)
	DefaultMultiplierPerYear = sdk.NewUint(1)
)

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{
		types.Params{
			GlobalDebtLimit: sdk.NewCoins(sdk.NewCoin(types.StableDenom, sdk.NewInt(25000000000000))),
			CollateralParams: types.CollateralParams{
				{
					Denom:            "ubtc",
					LiquidationRatio: sdk.MustNewDecFromStr("1.5"),
					DebtLimit:        sdk.NewCoins(sdk.NewCoin(types.StableDenom, sdk.NewInt(500000000000))),
					InterestModel:    types.NewCsdtInterest(DefaultBaseRatePerYear, DefaultMultiplierPerYear),
				},
				{
					Denom:            "ubnb",
					LiquidationRatio: sdk.MustNewDecFromStr("1.5"),
					DebtLimit:        sdk.NewCoins(sdk.NewCoin(types.StableDenom, sdk.NewInt(500000000000))),
					InterestModel:    types.NewCsdtInterest(DefaultBaseRatePerYear, DefaultMultiplierPerYear),
				},
				{
					Denom:            "ueth",
					LiquidationRatio: sdk.MustNewDecFromStr("1.5"),
					DebtLimit:        sdk.NewCoins(sdk.NewCoin(types.StableDenom, sdk.NewInt(500000000000))),
					InterestModel:    types.NewCsdtInterest(DefaultBaseRatePerYear, DefaultMultiplierPerYear),
				},
				{
					Denom:            "uftm",
					LiquidationRatio: sdk.MustNewDecFromStr("1.5"),
					DebtLimit:        sdk.NewCoins(sdk.NewCoin(types.StableDenom, sdk.NewInt(500000000000))),
					InterestModel:    types.NewCsdtInterest(DefaultBaseRatePerYear, DefaultMultiplierPerYear),
				},
				{
					Denom:            "uzar",
					LiquidationRatio: sdk.MustNewDecFromStr("1.3"),
					DebtLimit:        sdk.NewCoins(sdk.NewCoin(types.StableDenom, sdk.NewInt(500000000000))),
					InterestModel:    types.NewCsdtInterest(DefaultBaseRatePerYear, DefaultMultiplierPerYear),
				},
			},
		},
		sdk.ZeroInt(),
		types.CSDTs{},
		CoinUs{},
		CoinUs{},
		CoinUs{},
	}
}

func NewGenesisState(params types.Params, globalDebt sdk.Int) GenesisState {
	return GenesisState{
		Params:     params,
		GlobalDebt: globalDebt,
		CSDTs:      types.CSDTs{},
	}
}

// InitGenesis sets the genesis state in the keeper.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data GenesisState) {
	k.SetParams(ctx, data.Params)

	for _, csdt := range data.CSDTs {
		k.SetCSDT(ctx, csdt)
	}

	k.SetGlobalDebt(ctx, data.GlobalDebt)
}

// ValidateGenesis performs basic validation of genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
	// TODO implement this
	// validate denoms
	// check collateral debt limits sum to global limit?
	// check limits are > 0
	// check ratios are > 1
	// check no repeated denoms

	// check global debt is zero - force the chain to always start with zero stable coin, otherwise collateralStatus's will need to be set up as well.
	return nil
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) GenesisState {
	params := k.GetParams(ctx)
	csdts := types.CSDTs{}
	borrows := CoinUs{}
	supplies := CoinUs{}
	reserves := CoinUs{}
	for _, param := range params.CollateralParams {
		l, err := k.GetCSDTs(ctx, param.Denom, sdk.Dec{})
		if err != nil {
			panic(err)
		} else {
			csdts = append(csdts, l...)
		}

		borrow, ok := k.GetTotalBorrows(ctx, param.Denom)
		if !ok {
			panic(fmt.Sprintf("Failed to retrieve total borrows for '%s'", param.Denom))
		}
		borrows = append(borrows, CoinU{
			Denom:  param.Denom,
			Amount: borrow,
		})

		supply, ok := k.GetTotalCash(ctx, param.Denom)
		if !ok {
			panic(fmt.Sprintf("Failed to retrieve total cash/supply for '%s'", param.Denom))
		}
		supplies = append(supplies, CoinU{
			Denom:  param.Denom,
			Amount: supply,
		})

		reserve, ok := k.GetTotalReserve(ctx, param.Denom)
		if !ok {
			panic(fmt.Sprintf("Failed to retrieve total reserve for '%s'", param.Denom))
		}
		reserves = append(reserves, CoinU{
			Denom:  param.Denom,
			Amount: reserve,
		})
	}
	debt := k.GetGlobalDebt(ctx)

	return GenesisState{
		Params:        params,
		GlobalDebt:    debt,
		CSDTs:         csdts,
		TotalBorrows:  borrows,
		TotalSupplies: supplies,
		TotalReserves: reserves,
	}
}
