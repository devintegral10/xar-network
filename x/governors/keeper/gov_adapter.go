package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/tendermint/tendermint/crypto"
	"github.com/xar-network/xar-network/x/bonds"
)

// Adapter functions to be used with gov module

type (
	GovAdapter struct {
		k *Keeper
	}
	GovernorAdapter struct {
		bonds.Deposit
	}
	FakePubKey struct {
		addr sdk.Address
	}
)

// GovAdapter implements govtypes.StakingKeeper
var _ govtypes.StakingKeeper = GovAdapter{}

// GovernorAdapter implements ValidatorI
var _ staking.ValidatorI = GovernorAdapter{}

// GovernorAdapter implements ValidatorI
var _ crypto.PubKey = FakePubKey{}

func (g GovernorAdapter) IsJailed() bool {
	return false
}
func (g GovernorAdapter) GetMoniker() string {
	return g.DepositorAddress.String()
}
func (g GovernorAdapter) GetStatus() sdk.BondStatus {
	return sdk.Bonded
}
func (g GovernorAdapter) IsBonded() bool {
	return true
}
func (g GovernorAdapter) IsUnbonded() bool {
	return false
}
func (g GovernorAdapter) IsUnbonding() bool {
	return false
}
func (g GovernorAdapter) GetOperator() sdk.ValAddress {
	return g.DepositorAddress.Bytes()
}
func (g GovernorAdapter) GetConsPubKey() crypto.PubKey {
	return FakePubKey{g.DepositorAddress}
}
func (g GovernorAdapter) GetConsAddr() sdk.ConsAddress {
	return g.DepositorAddress.Bytes()
}
func (g GovernorAdapter) GetTokens() sdk.Int {
	return g.Tokens
}
func (g GovernorAdapter) GetBondedTokens() sdk.Int {
	return g.Tokens
}
func (g GovernorAdapter) GetConsensusPower() int64 {
	return 0
}
func (g GovernorAdapter) GetCommission() sdk.Dec {
	return sdk.ZeroDec()
}
func (g GovernorAdapter) GetMinSelfDelegation() sdk.Int {
	return sdk.ZeroInt()
}
func (g GovernorAdapter) GetDelegatorShares() sdk.Dec {
	return sdk.NewDec(1)
}
func (g GovernorAdapter) TokensFromShares(sdk.Dec) sdk.Dec {
	return sdk.ZeroDec()
}
func (g GovernorAdapter) TokensFromSharesTruncated(sdk.Dec) sdk.Dec {
	return sdk.ZeroDec()
}
func (g GovernorAdapter) TokensFromSharesRoundUp(sdk.Dec) sdk.Dec {
	return sdk.ZeroDec()
}
func (g GovernorAdapter) SharesFromTokens(amt sdk.Int) (sdk.Dec, sdk.Error) {
	return sdk.ZeroDec(), nil
}
func (g GovernorAdapter) SharesFromTokensTruncated(amt sdk.Int) (sdk.Dec, sdk.Error) {
	return sdk.ZeroDec(), nil
}

func (p FakePubKey) Address() crypto.Address {
	return p.addr.Bytes()
}
func (p FakePubKey) Bytes() []byte {
	return p.addr.Bytes()
}
func (p FakePubKey) VerifyBytes(msg []byte, sig []byte) bool {
	return false
}
func (p FakePubKey) Equals(p2 crypto.PubKey) bool {
	return false
}

func (k Keeper) AsGovAdapter() GovAdapter {
	return GovAdapter{&k}
}

// IterateBondedValidatorsByPower iterates governors in an undefined order
// (not sorted by power, not gov module doesn't depend on it!)
func (k GovAdapter) IterateBondedValidatorsByPower(ctx sdk.Context,
	f func(index int64, validator staking.ValidatorI) (stop bool)) {

	i := int64(0)
	k.k.IterateAllDeposits(ctx, func(deposit bonds.Deposit) (stop bool) {
		stop = f(i, &GovernorAdapter{deposit})
		i++
		return stop
	})
}

// TotalBondedTokens returns total supply of government token
func (k GovAdapter) TotalBondedTokens(ctx sdk.Context) sdk.Int {
	return k.k.TokenSupply(ctx)
}

// IterateDelegations does nothing
func (k GovAdapter) IterateDelegations(ctx sdk.Context, delegator sdk.AccAddress,
	fn func(index int64, delegation staking.DelegationI) (stop bool)) {

	return
}
