package simulation

// DONTCOVER

import (
	"fmt"
	"github.com/xar-network/xar-network/x/bonds"
	"github.com/xar-network/xar-network/x/governors/types"
	"math/rand"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

// Simulation parameter constants
const (
	UnbondingTime = "unbonding_time"
)

// GenUnbondingTime randomized UnbondingTime
func GenUnbondingTime(r *rand.Rand) (ubdTime time.Duration) {
	return time.Duration(simulation.RandIntBetween(r, 60, 60*60*24*3*2)) * time.Second
}

// GenMaxValidators randomized MaxValidators
func GenMaxValidators(r *rand.Rand) (maxValidators uint16) {
	return uint16(r.Intn(250) + 1)
}

// RandomizedGenState generates a random GenesisState for the module
func RandomizedGenState(simState *module.SimulationState) {
	// params
	var unbondTime time.Duration
	simState.AppParams.GetOrGenerate(
		simState.Cdc, UnbondingTime, &unbondTime, simState.Rand,
		func(r *rand.Rand) { unbondTime = GenUnbondingTime(r) },
	)

	// NOTE: the slashing module need to be defined after the module on the
	// NewSimulationManager constructor for this to work
	simState.UnbondTime = unbondTime

	params := types.NewParams(bonds.NewParams(simState.UnbondTime, 7, types.DefaultBondDenom, types.BondedPoolName, types.NotBondedPoolName))

	// deposits
	var (
		deposits []bonds.Deposit
	)

	for i := 0; i < int(simState.NumBonded); i++ {
		deposit := bonds.NewDeposit(simState.Accounts[i].Address, sdk.NewInt(simState.InitialStake))
		deposits = append(deposits, deposit)
	}

	genesis := types.NewGenesisState(params, bonds.NewGenesisState(params.Bonds, deposits).GenesisStateData)

	fmt.Printf("Selected randomly generated governors parameters:\n%s\n", codec.MustMarshalJSONIndent(simState.Cdc, genesis.Params))
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(genesis)
}
