// nolint
package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CodeType = sdk.CodeType

const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeInvalidDeposit CodeType = 102
	CodeInvalidInput   CodeType = 103
	CodeInvalidAddress CodeType = sdk.CodeInvalidAddress
	CodeUnauthorized   CodeType = sdk.CodeUnauthorized
	CodeInternal       CodeType = sdk.CodeInternal
	CodeUnknownRequest CodeType = sdk.CodeUnknownRequest
)

func ErrNilDepositorAddr(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, "depositor address is nil")
}

func ErrBadDenom(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidDeposit, "invalid coin denomination")
}

func ErrBadDepositAddr(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, "unexpected address length for this address")
}

func ErrBadDepositAmount(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidDeposit, "amount must be > 0")
}

func ErrNoDeposit(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidDeposit, "no deposit for this address")
}

func ErrBadDepositorAddr(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidDeposit, "depositor does not exist for that address")
}

func ErrNoDepositorForAddress(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidDeposit, "depositor does not contain this deposit")
}

func ErrNotEnoughDepositTokens(codespace sdk.CodespaceType, tokens string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidDeposit, fmt.Sprintf("not enough tokens only have %v", tokens))
}

func ErrNotMature(codespace sdk.CodespaceType, operation, descriptor string, got, min time.Time) sdk.Error {
	msg := fmt.Sprintf("%v is not mature requires a min %v of %v, currently it is %v",
		operation, descriptor, got, min)
	return sdk.NewError(codespace, CodeUnauthorized, msg)
}

func ErrNoUnbondingDeposit(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidDeposit, "no unbonding deposit found")
}

func ErrMaxUnbondingDepositEntries(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidDeposit,
		"too many unbonding deposit entries in this depositor, please wait for some entries to mature")
}
