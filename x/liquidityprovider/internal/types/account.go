package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

var _ exported.Account = (*LiquidityProviderAccount)(nil)

func init() {
	// Register the LiquidityProviderAccount on the auth module codec
	authtypes.RegisterAccountTypeCodec(&LiquidityProviderAccount{}, "liquidityprovider/LiquidityProviderAccount")
}

type LiquidityProviderAccount struct {
	exported.Account

	Credit sdk.Coins `json:"credit" yaml:"credit"`
}

func NewLiquidityProviderAccount(baseAccount exported.Account, credit sdk.Coins) *LiquidityProviderAccount {
	return &LiquidityProviderAccount{
		Account: baseAccount,
		Credit:  credit,
	}
}

func (acc *LiquidityProviderAccount) IncreaseCredit(increase sdk.Coins) {
	acc.Credit = acc.Credit.Add(increase)
}

// Function panics if resulting credit is negative. Should be checked prior to invocation for cleaner handling.
func (acc *LiquidityProviderAccount) DecreaseCredit(decrease sdk.Coins) {
	if newCredit, anyNegative := acc.Credit.SafeSub(decrease); !anyNegative {
		acc.Credit = newCredit
		return
	}

	panic(fmt.Errorf("credit cannot be negative"))
}

func (acc LiquidityProviderAccount) String() string {
	var pubkey string

	if acc.GetPubKey() != nil {
		pubkey = sdk.MustBech32ifyAccPub(acc.GetPubKey())
	}

	return fmt.Sprintf(`Account:
  Address:       %s
  Pubkey:        %s
  Credit:        %s
  Coins:         %s
  AccountNumber: %d
  Sequence:      %d`,
		acc.GetAddress(), pubkey, acc.Credit, acc.GetCoins(), acc.GetAccountNumber(), acc.GetSequence(),
	)
}