package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ sdk.Msg = MsgIncreaseCredit{}
	_ sdk.Msg = MsgDecreaseCredit{}
	_ sdk.Msg = MsgRevokeLiquidityProvider{}
	_ sdk.Msg = MsgSetInflation{}
)

// Increase the credit of a liquidity provider. If the account is not previously an LP, it will be made one.
type (
	MsgIncreaseCredit struct {
		CreditIncrease    sdk.Coins
		LiquidityProvider sdk.AccAddress
		Issuer            sdk.AccAddress
	}

	MsgDecreaseCredit struct {
		CreditDecrease    sdk.Coins
		LiquidityProvider sdk.AccAddress
		Issuer            sdk.AccAddress
	}

	MsgRevokeLiquidityProvider struct {
		LiquidityProvider sdk.AccAddress
		Issuer            sdk.AccAddress
	}

	MsgSetInflation struct {
		Denom         string
		InflationRate sdk.Dec
		Issuer        sdk.AccAddress
	}
)

func (msg MsgSetInflation) Route() string { return ModuleName }

func (msg MsgSetInflation) Type() string { return "setInflation" }

func (msg MsgSetInflation) ValidateBasic() sdk.Error {
	if msg.InflationRate.IsNegative() {
		return ErrNegativeInflation()
	}

	if msg.Issuer.Empty() {
		return sdk.ErrInvalidAddress("missing issuer address")
	}

	return nil
}

func (msg MsgSetInflation) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgSetInflation) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Issuer}
}

func (msg MsgRevokeLiquidityProvider) Route() string { return ModuleName }

func (msg MsgRevokeLiquidityProvider) Type() string { return "revokeLiquidityProvider" }

func (msg MsgRevokeLiquidityProvider) ValidateBasic() sdk.Error {
	if msg.LiquidityProvider.Empty() {
		return sdk.ErrInvalidAddress("missing liquidity provider address")
	}

	if msg.Issuer.Empty() {
		return sdk.ErrInvalidAddress("missing issuer address")
	}

	return nil
}

func (msg MsgRevokeLiquidityProvider) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgRevokeLiquidityProvider) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Issuer}
}

func (msg MsgDecreaseCredit) Route() string { return ModuleName }

func (msg MsgDecreaseCredit) Type() string { return "decreaseCredit" }

func (msg MsgDecreaseCredit) ValidateBasic() sdk.Error {
	if msg.LiquidityProvider.Empty() {
		return sdk.ErrInvalidAddress("missing liquidity provider address")
	}

	if msg.Issuer.Empty() {
		return sdk.ErrInvalidAddress("missing issuer address")
	}

	if !msg.CreditDecrease.IsValid() {
		return sdk.ErrInvalidCoins("credit decrease is invalid: " + msg.CreditDecrease.String())
	}

	return nil
}

func (msg MsgDecreaseCredit) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgDecreaseCredit) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Issuer}
}

func (msg MsgIncreaseCredit) Route() string { return ModuleName }

func (msg MsgIncreaseCredit) Type() string { return "increaseCredit" }

func (msg MsgIncreaseCredit) ValidateBasic() sdk.Error {
	if msg.LiquidityProvider.Empty() {
		return sdk.ErrInvalidAddress("missing liquidity provider address")
	}

	if msg.Issuer.Empty() {
		return sdk.ErrInvalidAddress("missing issuer address")
	}

	if !msg.CreditIncrease.IsValid() {
		return sdk.ErrInvalidCoins("credit increase is invalid: " + msg.CreditIncrease.String())
	}

	return nil
}

func (msg MsgIncreaseCredit) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgIncreaseCredit) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Issuer}
}
