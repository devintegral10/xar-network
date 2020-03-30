package simulation

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/xar-network/xar-network/x/governors/keeper"
	"github.com/xar-network/xar-network/x/governors/types"
)

// SimulateMsgBond generates a MsgBond with random values
// nolint: funlen
func SimulateMsgBond(ak types.AccountKeeper, k keeper.Keeper) simulation.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string,
	) (simulation.OperationMsg, []simulation.FutureOperation, error) {

		denom := k.GetParams(ctx).BondDenom

		simAccount, _ := simulation.RandomAcc(r, accs)

		amount := ak.GetAccount(ctx, simAccount.Address).GetCoins().AmountOf(denom)
		if !amount.IsPositive() {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		amount, err := simulation.RandPositiveInt(r, amount)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		bondAmt := sdk.NewCoin(denom, amount)

		account := ak.GetAccount(ctx, simAccount.Address)
		coins := account.SpendableCoins(ctx.BlockTime())

		var fees sdk.Coins
		coins, hasNeg := coins.SafeSub(sdk.Coins{bondAmt})
		if !hasNeg {
			fees, err = simulation.RandomFees(r, ctx, coins)
			if err != nil {
				return simulation.NoOpMsg(types.ModuleName), nil, err
			}
		}

		msg := types.NewMsgBond(simAccount.Address, bondAmt)

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)

		res := app.Deliver(tx)
		if !res.IsOK() {
			return simulation.NoOpMsg(types.ModuleName), nil, errors.New(res.Log)
		}

		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgUnbond generates a MsgUnbond with random values
// nolint: funlen
func SimulateMsgUnbond(ak types.AccountKeeper, k keeper.Keeper) simulation.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simulation.Account, chainID string,
	) (simulation.OperationMsg, []simulation.FutureOperation, error) {

		deposits := k.GetAllDeposits(ctx)

		// get random depositor
		deposit := deposits[r.Intn(len(deposits))]
		depAddr := deposit.GetDepositorAddr()

		if k.HasMaxUnbondingDepositEntries(ctx, depAddr) {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		unbondAmt, err := simulation.RandPositiveInt(r, deposit.Tokens)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		if unbondAmt.IsZero() {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		msg := types.NewMsgUnbond(
			depAddr, sdk.NewCoin(k.BondDenom(ctx), unbondAmt),
		)

		// need to retrieve the simulation account associated with deposit to retrieve PrivKey
		var simAccount simulation.Account
		for _, simAcc := range accs {
			if simAcc.Address.Equals(depAddr) {
				simAccount = simAcc
				break
			}
		}
		// if simaccount.PrivKey == nil, deposit address does not exist in accs. Return error
		if simAccount.PrivKey == nil {
			return simulation.NoOpMsg(types.ModuleName), nil, fmt.Errorf("deposit addr: %s does not exist in simulation accounts", depAddr)
		}

		account := ak.GetAccount(ctx, depAddr)
		fees, err := simulation.RandomFees(r, ctx, account.SpendableCoins(ctx.BlockTime()))
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		tx := helpers.GenTx(
			[]sdk.Msg{msg},
			fees,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)

		res := app.Deliver(tx)
		if !res.IsOK() {
			return simulation.NoOpMsg(types.ModuleName), nil, errors.New(res.Log)
		}

		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}
