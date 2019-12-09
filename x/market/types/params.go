/*

Copyright 2019 All in Bits, Inc
Copyright 2019 Xar Network

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

*/

package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params/subspace"
)

// Parameter keys
var (
	// ParamStoreKeyAuctionParams Param store key for auction params
	KeyMarkets  = []byte(ModuleName)
	KeyNominees = []byte(ModuleNominees)
)

type Params struct {
	Markets  []Market
	Nominees []string
}

// ParamKeyTable Key declaration for parameters
func ParamKeyTable() subspace.KeyTable {
	return subspace.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of auth module's parameters.
// nolint
func (p *Params) ParamSetPairs() subspace.ParamSetPairs {
	return subspace.ParamSetPairs{
		subspace.NewParamSetPair(KeyMarkets, &p.Markets),
		subspace.NewParamSetPair(KeyNominees, &p.Nominees),
	}
}

// NewParams creates a new Params object
func NewParams(markets []Market, nominees []string) Params {
	return Params{
		Markets:  markets,
		Nominees: nominees,
	}
}

// DefaultParams default params
func DefaultParams() Params {
	return NewParams(Markets{}, []string{})
}

// String implements fmt.stringer
func (p Params) String() string {
	out := "Params:\n"
	for _, a := range p.Markets {
		out += a.String()
	}
	for _, n := range p.Nominees {
		out += n
	}
	return strings.TrimSpace(out)
}

// ParamSubspace defines the expected Subspace interface for parameters
type ParamSubspace interface {
	Get(ctx sdk.Context, key []byte, ptr interface{})
	Set(ctx sdk.Context, key []byte, param interface{})
}

// Validate ensure that params have valid values
func (p Params) Validate() error {
	// iterate over assets and verify them
	for _, market := range p.Markets {
		if !market.ID.IsDefined() {
			return fmt.Errorf("invalid id: %s. missing id", market.String())
		}
	}
	return nil
}