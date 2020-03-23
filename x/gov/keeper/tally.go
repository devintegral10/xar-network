package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/types"
)

var (
	// TotalSupply is total supply of government tokens
	TotalSupply = sdk.NewDec(1000000000) // TODO move into genesis
	// Denom is government tokens
	Denom = "uxar" // TODO move into genesis
)

// TODO: Break into several smaller functions for clarity

// Tally iterates over the votes and updates the tally of a proposal based on the voting power of the
// voters
func (keeper Keeper) Tally(ctx sdk.Context, proposal types.Proposal) (passes bool, burnDeposits bool, tallyResults types.TallyResult) {
	results := make(map[types.VoteOption]sdk.Dec)
	results[types.OptionYes] = sdk.ZeroDec()
	results[types.OptionAbstain] = sdk.ZeroDec()
	results[types.OptionNo] = sdk.ZeroDec()
	results[types.OptionNoWithVeto] = sdk.ZeroDec()

	totalVotingPower := sdk.ZeroDec()

	keeper.IterateVotes(ctx, proposal.ProposalID, func(vote types.Vote) bool {
		if vote.Option == types.OptionEmpty {
			return false
		}

		votingPower := keeper.bank.GetCoins(ctx, vote.Voter).AmountOf(Denom).ToDec()

		results[vote.Option] = results[vote.Option].Add(votingPower)
		totalVotingPower = totalVotingPower.Add(votingPower)

		keeper.deleteVote(ctx, vote.ProposalID, vote.Voter)
		return false
	})

	tallyParams := keeper.GetTallyParams(ctx)
	tallyResults = types.NewTallyResultFromMap(results)

	// If there is not enough quorum of votes, the proposal fails
	percentVoting := totalVotingPower.Quo(TotalSupply)
	if percentVoting.LT(tallyParams.Quorum) {
		return false, true, tallyResults
	}

	// If no one votes (everyone abstains), proposal fails
	if totalVotingPower.Sub(results[types.OptionAbstain]).Equal(sdk.ZeroDec()) {
		return false, false, tallyResults
	}

	// If more than 1/3 of voters veto, proposal fails
	if results[types.OptionNoWithVeto].Quo(totalVotingPower).GT(tallyParams.Veto) {
		return false, true, tallyResults
	}

	// If more than 1/2 of non-abstaining voters vote Yes, proposal passes
	if results[types.OptionYes].Quo(totalVotingPower.Sub(results[types.OptionAbstain])).GT(tallyParams.Threshold) {
		return true, false, tallyResults
	}

	// If more than 1/2 of non-abstaining voters vote No, proposal fails
	return false, false, tallyResults
}
