package types

import (
	"github.com/xar-network/xar-network/x/bonds"
)

// GenesisState - all state that must be provided at genesis
type GenesisState struct {
	Params   Params                 `json:"params" yaml:"params"`
	Bonds    bonds.GenesisStateData `json:"bonds" yaml:"bonds"`
	Exported bool                   `json:"exported" yaml:"exported"`
}

func NewGenesisState(params Params, bonds bonds.GenesisStateData) GenesisState {
	return GenesisState{
		Params: params,
		Bonds:  bonds,
	}
}

func (g *GenesisState) BondsGenesis() bonds.GenesisState {
	return bonds.GenesisState{
		GenesisStateData: g.Bonds,
		Params:           g.Params.Bonds,
		Exported:         g.Exported,
	}
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params: DefaultParams(),
	}
}
