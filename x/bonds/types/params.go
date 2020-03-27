package types

import (
	"bytes"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// nolint - Keys for parameter access
var (
	KeyUnbondingTime     = []byte("UnbondingTime")
	KeyMaxEntries        = []byte("KeyMaxEntries")
	KeyBondDenom         = []byte("BondDenom")
	KeyNotBondedPoolName = []byte("NotBondedPoolName")
	KeyBondedPoolName    = []byte("BondedPoolName")
)

var _ params.ParamSet = (*Params)(nil)

// Params defines the high level settings
type Params struct {
	UnbondingTime time.Duration `json:"unbonding_time" yaml:"unbonding_time"` // time duration of unbonding
	MaxEntries    uint16        `json:"max_entries" yaml:"max_entries"`       // max entries for either unbonding deposit or redeposit (per pair/trio)
	// note: we need to be a bit careful about potential overflow here, since this is user-determined
	BondDenom         string `json:"bond_denom" yaml:"bond_denom"`                     // bondable coin denomination
	BondedPoolName    string `json:"bonded_pool_name" yaml:"bonded_pool_name"`         //  pool for bonded tokens
	NotBondedPoolName string `json:"not_bonded_pool_name" yaml:"not_bonded_pool_name"` //  pool for not bonded tokens
}

// NewParams creates a new Params instance
func NewParams(unbondingTime time.Duration, maxEntries uint16,
	bondDenom string, bondedPoolName string, notBondedPoolName string) Params {

	return Params{
		UnbondingTime:     unbondingTime,
		MaxEntries:        maxEntries,
		BondDenom:         bondDenom,
		BondedPoolName:    bondedPoolName,
		NotBondedPoolName: notBondedPoolName,
	}
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyUnbondingTime, Value: &p.UnbondingTime},
		{Key: KeyMaxEntries, Value: &p.MaxEntries},
		{Key: KeyBondDenom, Value: &p.BondDenom},
		{Key: KeyBondedPoolName, Value: &p.BondedPoolName},
		{Key: KeyNotBondedPoolName, Value: &p.NotBondedPoolName},
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
	return fmt.Sprintf(`Params:
  Unbonding Time:    %s
  Max Entries:       %d
  Bonded Coin Denom: %s
  Bonded Pool Name: %s
  Not Bonded Pool Name: %s`, p.UnbondingTime,
		p.MaxEntries, p.BondDenom, p.BondedPoolName, p.NotBondedPoolName)
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
	if p.BondDenom == "" {
		return fmt.Errorf("parameter BondDenom can't be an empty string")
	}
	if p.BondedPoolName == "" {
		return fmt.Errorf("parameter BondedPoolName can't be an empty string")
	}
	if p.NotBondedPoolName == "" {
		return fmt.Errorf("parameter NotBondedPoolName can't be an empty string")
	}
	return nil
}
