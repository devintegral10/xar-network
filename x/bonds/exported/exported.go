package exported

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DepositI is a bonded deposit
type DepositI interface {
	GetDepositorAddr() sdk.AccAddress // depositor sdk.AccAddress for the bond
	GetTokens() sdk.Int               // amount of tokens held in this deposit
}
