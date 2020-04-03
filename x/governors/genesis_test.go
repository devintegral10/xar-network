package governors

import (
	"testing"

	"github.com/stretchr/testify/assert"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/x/bonds"
	bondskeeper "github.com/xar-network/xar-network/x/bonds/keeper"
	"github.com/xar-network/xar-network/x/governors/types"
)

func genGenesisState() types.GenesisState {
	genesisState := types.DefaultGenesisState()

	for i := 0; i < 10; i++ {
		genesisState.Bonds.Deposits = append(genesisState.Bonds.Deposits, bonds.Deposit{
			DepositorAddress: bondskeeper.Addrs[i],
			Tokens:           sdk.NewInt(10),
		})
	}
	return genesisState
}

func TestValidateGenesis(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(*types.GenesisState)
		wantErr bool
	}{
		{"default", func(*types.GenesisState) {}, false},
		// validate genesis validators
		{"duplicate depositor", func(data *types.GenesisState) {
			data.Bonds.Deposits = append(data.Bonds.Deposits, data.Bonds.Deposits[0])
		}, true},
		{"no delegator tokens", func(data *types.GenesisState) {
			data.Bonds.Deposits[0].Tokens = sdk.ZeroInt()
		}, true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			genesisState := genGenesisState()
			tt.mutate(&genesisState)
			if tt.wantErr {
				assert.Error(t, ValidateGenesis(genesisState))
			} else {
				assert.NoError(t, ValidateGenesis(genesisState))
			}
		})
	}
}
