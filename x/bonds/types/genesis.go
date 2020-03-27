package types

// GenesisState - all state that must be provided at genesis
type GenesisState struct {
	GenesisStateData
	Params   Params `json:"params" yaml:"params"`
	Exported bool   `json:"exported" yaml:"exported"`
}

// GenesisState - all state that must be provided at genesis
type GenesisStateData struct {
	Deposits          Deposits           `json:"deposits" yaml:"deposits"`
	UnbondingDeposits []UnbondingDeposit `json:"unbonding_deposits" yaml:"unbonding_deposits"`
}

func NewGenesisState(params Params, deposits []Deposit) GenesisState {
	return GenesisState{
		Params: params,
		GenesisStateData: GenesisStateData{
			Deposits:          deposits,
			UnbondingDeposits: nil,
		},
	}
}
