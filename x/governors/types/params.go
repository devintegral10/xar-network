package types

import (
	"bytes"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/xar-network/xar-network/x/bonds"
)

const (
	// Default governors denom
	DefaultBondDenom = sdk.DefaultBondDenom

	// DefaultUnbondingTime reflects three weeks in seconds as the default
	// unbonding time.
	DefaultUnbondingTime = 24 * time.Hour

	// Default maximum unbonding entries
	DefaultMaxEntries = 7

	BondedPoolName    = "governors_bonded_pool"
	NotBondedPoolName = "governors_unbonded_pool"
)

// Default parameter namespace
const (
	DefaultParamspace = ModuleName
)

// nolint - Keys for parameter access
var (
	KeyBonds = []byte("Bonds")
)

var _ params.ParamSet = (*Params)(nil)

// Params defines the high level settings
type Params struct {
	Bonds bonds.Params `json:"bonds" yaml:"bonds"`
}

// NewParams creates a new Params instance
func NewParams(bonds bonds.Params) Params {
	return Params{
		Bonds: bonds,
	}
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyBonds, Value: &p.Bonds},
	}
}

// Equal returns a boolean determining if two Param types are identical.
// TODO: This is slower than comparing struct fields directly
func (p Params) Equal(p2 Params) bool {
	bz1 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p)
	bz2 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p2)
	return bytes.Equal(bz1, bz2)
}

// String returns a human readable string representation of the parameters.
func (p Params) String() string {
	return p.Bonds.String()
}

// unmarshal the current params value from store key or panic
func MustUnmarshalParams(cdc *codec.Codec, value []byte) Params {
	params, err := UnmarshalParams(cdc, value)
	if err != nil {
		panic(err)
	}
	return params
}

// unmarshal the current params value from store key
func UnmarshalParams(cdc *codec.Codec, value []byte) (params Params, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &params)
	if err != nil {
		return
	}
	return
}

// validate a set of params
func (p Params) Validate() error {
	return p.Bonds.Validate()
}

// DefaultParams default params
func DefaultParams() Params {
	return NewParams(bonds.Params{
		UnbondingTime:     DefaultUnbondingTime,
		MaxEntries:        DefaultMaxEntries,
		BondDenom:         DefaultBondDenom,
		BondedPoolName:    BondedPoolName,
		NotBondedPoolName: NotBondedPoolName,
	})
}
