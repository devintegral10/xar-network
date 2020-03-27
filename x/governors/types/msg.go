package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xar-network/xar-network/x/bonds"
)

// ensure Msg interface compliance at compile time
var (
	_ sdk.Msg = &MsgBond{}
	_ sdk.Msg = &MsgUnbond{}
)

// MsgBond - struct for bonding transactions
type MsgBond struct {
	DepositorAddress sdk.AccAddress `json:"depositor_address" yaml:"depositor_address"`
	Amount           sdk.Coin       `json:"amount" yaml:"amount"`
}

func NewMsgBond(depAddr sdk.AccAddress, amount sdk.Coin) MsgBond {
	return MsgBond{
		DepositorAddress: depAddr,
		Amount:           amount,
	}
}

//nolint
func (msg MsgBond) Route() string { return RouterKey }
func (msg MsgBond) Type() string  { return "bond" }
func (msg MsgBond) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.DepositorAddress}
}

// get the bytes for the message signer to sign on
func (msg MsgBond) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// quick validity check
func (msg MsgBond) ValidateBasic() sdk.Error {
	if msg.DepositorAddress.Empty() {
		return bonds.ErrNilDepositorAddr(DefaultCodespace)
	}
	if msg.Amount.Amount.LTE(sdk.ZeroInt()) {
		return bonds.ErrBadDepositAmount(DefaultCodespace)
	}
	return nil
}

// MsgUnbond - struct for unbonding transactions
type MsgUnbond struct {
	DepositorAddress sdk.AccAddress `json:"depositor_address" yaml:"depositor_address"`
	Amount           sdk.Coin       `json:"amount" yaml:"amount"`
}

func NewMsgUnbond(depAddr sdk.AccAddress, amount sdk.Coin) MsgUnbond {
	return MsgUnbond{
		DepositorAddress: depAddr,
		Amount:           amount,
	}
}

//nolint
func (msg MsgUnbond) Route() string                { return RouterKey }
func (msg MsgUnbond) Type() string                 { return "begin_unbonding" }
func (msg MsgUnbond) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.DepositorAddress} }

// get the bytes for the message signer to sign on
func (msg MsgUnbond) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// quick validity check
func (msg MsgUnbond) ValidateBasic() sdk.Error {
	if msg.DepositorAddress.Empty() {
		return bonds.ErrNilDepositorAddr(DefaultCodespace)
	}
	if msg.Amount.Amount.LTE(sdk.ZeroInt()) {
		return bonds.ErrBadDepositAmount(DefaultCodespace)
	}
	return nil
}
