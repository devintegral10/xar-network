package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the  module
	ModuleName = "bond"

	// StoreKey is the string store representation
	StoreKey = ModuleName

	// QuerierRoute is the querier route for the module
	QuerierRoute = ModuleName

	// RouterKey is the msg router key for the module
	RouterKey = ModuleName
)

//nolint
var (
	// Keys for store prefixes
	// Last* values are constant during a block.
	DepositKey          = []byte{0x21} // key for a deposit
	UnbondingDepositKey = []byte{0x22} // key for an unbonding-deposit

	UnbondingQueueKey = []byte{0x31} // prefix for the timestamps in unbonding queue
)

//______________________________________________________________________________

// gets the prefix for a depositor
func GetDepositKey(depAddr sdk.AccAddress) []byte {
	return append(DepositKey, depAddr.Bytes()...)
}

//______________

// gets the prefix for all unbonding deposits from a depositor
func GetUBDKey(depAddr sdk.AccAddress) []byte {
	return append(UnbondingDepositKey, depAddr.Bytes()...)
}

// gets the prefix for all unbonding deposits from a depositor
func GetUnbondingDepositTimeKey(timestamp time.Time) []byte {
	bz := sdk.FormatTimeBytes(timestamp)
	return append(UnbondingQueueKey, bz...)
}
