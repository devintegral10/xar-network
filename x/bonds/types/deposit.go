package types

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/x/bonds/exported"
)

// Implements Deposit interface
var _ exported.DepositI = Deposit{}

// Deposit represents the bond with tokens held by an account.
type Deposit struct {
	DepositorAddress sdk.AccAddress `json:"depositor_address" yaml:"depositor_address"`
	Tokens           sdk.Int        `json:"tokens" yaml:"tokens"`
}

// NewDeposit creates a new deposit object
func NewDeposit(depositorAddr sdk.AccAddress, tokens sdk.Int) Deposit {

	return Deposit{
		DepositorAddress: depositorAddr,
		Tokens:           tokens,
	}
}

// return the deposit
func MustMarshalDeposit(cdc *codec.Codec, deposit Deposit) []byte {
	return cdc.MustMarshalBinaryLengthPrefixed(deposit)
}

// return the deposit
func MustUnmarshalDeposit(cdc *codec.Codec, value []byte) Deposit {
	deposit, err := UnmarshalDeposit(cdc, value)
	if err != nil {
		panic(err)
	}
	return deposit
}

// return the deposit
func UnmarshalDeposit(cdc *codec.Codec, value []byte) (deposit Deposit, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &deposit)
	return deposit, err
}

// nolint
func (d Deposit) Equal(d2 Deposit) bool {
	return bytes.Equal(d.DepositorAddress, d2.DepositorAddress) &&
		d.Tokens.Equal(d2.Tokens)
}

// nolint - for Deposit
func (d Deposit) GetDepositorAddr() sdk.AccAddress { return d.DepositorAddress }
func (d Deposit) GetTokens() sdk.Int               { return d.Tokens }

// String returns a human readable string representation of a Deposit.
func (d Deposit) String() string {
	return fmt.Sprintf(`Deposit:
  Depositor: %s
  Tokens: %s`, d.DepositorAddress,
		d.Tokens)
}

// Deposits is a collection of deposits
type Deposits []Deposit

func (d Deposits) String() (out string) {
	for _, dep := range d {
		out += dep.String() + "\n"
	}
	return strings.TrimSpace(out)
}

// UnbondingDeposit stores all of a single depositor's unbonding deposits
type UnbondingDeposit struct {
	DepositorAddress sdk.AccAddress          `json:"depositor_address" yaml:"depositor_address"` // depositor
	Entries          []UnbondingDepositEntry `json:"entries" yaml:"entries"`                     // unbonding deposit entries
}

// UnbondingDepositEntry - entry to an UnbondingDeposit
type UnbondingDepositEntry struct {
	CreationHeight int64     `json:"creation_height" yaml:"creation_height"` // height which the unbonding took place
	CompletionTime time.Time `json:"completion_time" yaml:"completion_time"` // time at which the unbonding deposit will complete
	InitialBalance sdk.Int   `json:"initial_balance" yaml:"initial_balance"` // atoms initially scheduled to receive at completion
	Balance        sdk.Int   `json:"balance" yaml:"balance"`                 // atoms to receive at completion
}

// IsMature - is the current entry mature
func (e UnbondingDepositEntry) IsMature(currentTime time.Time) bool {
	return !e.CompletionTime.After(currentTime)
}

// NewUnbondingDeposit - create a new unbonding deposit object
func NewUnbondingDeposit(depositorAddr sdk.AccAddress,
	creationHeight int64, minTime time.Time,
	balance sdk.Int) UnbondingDeposit {

	entry := NewUnbondingDepositEntry(creationHeight, minTime, balance)
	return UnbondingDeposit{
		DepositorAddress: depositorAddr,
		Entries:          []UnbondingDepositEntry{entry},
	}
}

// NewUnbondingDeposit - create a new unbonding deposit object
func NewUnbondingDepositEntry(creationHeight int64, completionTime time.Time,
	balance sdk.Int) UnbondingDepositEntry {

	return UnbondingDepositEntry{
		CreationHeight: creationHeight,
		CompletionTime: completionTime,
		InitialBalance: balance,
		Balance:        balance,
	}
}

// AddEntry - append entry to the unbonding deposit
func (d *UnbondingDeposit) AddEntry(creationHeight int64,
	minTime time.Time, balance sdk.Int) {

	entry := NewUnbondingDepositEntry(creationHeight, minTime, balance)
	d.Entries = append(d.Entries, entry)
}

// RemoveEntry - remove entry at index i to the unbonding deposit
func (d *UnbondingDeposit) RemoveEntry(i int64) {
	d.Entries = append(d.Entries[:i], d.Entries[i+1:]...)
}

// return the unbonding deposit
func MustMarshalUBD(cdc *codec.Codec, ubd UnbondingDeposit) []byte {
	return cdc.MustMarshalBinaryLengthPrefixed(ubd)
}

// unmarshal a unbonding deposit from a store value
func MustUnmarshalUBD(cdc *codec.Codec, value []byte) UnbondingDeposit {
	ubd, err := UnmarshalUBD(cdc, value)
	if err != nil {
		panic(err)
	}
	return ubd
}

// unmarshal a unbonding deposit from a store value
func UnmarshalUBD(cdc *codec.Codec, value []byte) (ubd UnbondingDeposit, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &ubd)
	return ubd, err
}

// nolint
// inefficient but only used in testing
func (d UnbondingDeposit) Equal(d2 UnbondingDeposit) bool {
	bz1 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&d)
	bz2 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&d2)
	return bytes.Equal(bz1, bz2)
}

// String returns a human readable string representation of an UnbondingDeposit.
func (d UnbondingDeposit) String() string {
	out := fmt.Sprintf(`Unbonding Deposits between:
  Depositor:                 %s
	Entries:`, d.DepositorAddress)
	for i, entry := range d.Entries {
		out += fmt.Sprintf(`    Unbonding Deposit %d:
      Creation Height:           %v
      Min time to unbond (unix): %v
      Expected balance:          %s`, i, entry.CreationHeight,
			entry.CompletionTime, entry.Balance)
	}
	return out
}

// UnbondingDeposits is a collection of UnbondingDeposit
type UnbondingDeposits []UnbondingDeposit

func (ubds UnbondingDeposits) String() (out string) {
	for _, u := range ubds {
		out += u.String() + "\n"
	}
	return strings.TrimSpace(out)
}
