package types

import sdk "github.com/cosmos/cosmos-sdk/types"

/*
Message types for starting various auctions.
Note: these message types are not final and will likely change.
Design options and problems:
 - msgs that only start auctions
	- senders have to pay fees
	- these msgs cannot be bundled into a tx with a PlaceBid msg because PlaceBid requires an auction ID
 - msgs that start auctions and place an initial bid
	- place bid can fail, leaving auction without bids which is similar to first case
 - no msgs, auctions started automatically
	- running this as an endblocker adds complexity and potential vulnerabilities
*/

type MsgSeizeAndStartCollateralAuction struct {
	Sender          sdk.AccAddress `json:"sender" yaml:"sender"`
	CsdtOwner       sdk.AccAddress `json:"owner" yaml:"owner"`
	CollateralDenom string         `json:"collateral_denom" yaml:"collateral_denom"`
}

// Route return the message type used for routing the message.
func (msg MsgSeizeAndStartCollateralAuction) Route() string { return "liquidator" }

// Type returns a human-readable string for the message, intended for utilization within tags.
func (msg MsgSeizeAndStartCollateralAuction) Type() string { return "seize_and_start_auction" } // TODO snake case?

// ValidateBasic does a simple validation check that doesn't require access to any other information.
func (msg MsgSeizeAndStartCollateralAuction) ValidateBasic() sdk.Error {
	if msg.Sender.Empty() {
		return sdk.ErrInternal("invalid (empty) sender address")
	}
	if msg.CsdtOwner.Empty() {
		return sdk.ErrInternal("invalid (empty) CSDT owner address")
	}
	// TODO check coin denoms
	return nil
}

// GetSignBytes gets the canonical byte representation of the Msg.
func (msg MsgSeizeAndStartCollateralAuction) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners returns the addresses of signers that must sign.
func (msg MsgSeizeAndStartCollateralAuction) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

type MsgStartDebtAuction struct {
	Sender sdk.AccAddress `json:"sender" yaml:"sender"`
}

func (msg MsgStartDebtAuction) Route() string { return "liquidator" }
func (msg MsgStartDebtAuction) Type() string  { return "start_debt_auction" } // TODO snake case?
func (msg MsgStartDebtAuction) ValidateBasic() sdk.Error {
	if msg.Sender.Empty() {
		return sdk.ErrInternal("invalid (empty) sender address")
	}
	return nil
}
func (msg MsgStartDebtAuction) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}
func (msg MsgStartDebtAuction) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.Sender} }
