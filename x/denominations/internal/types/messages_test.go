package types

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

type MsgInterface interface{ sdk.Msg }

func validateError(cases []struct {
	valid bool
	tx    MsgInterface
}, t *testing.T) {
	for i, tc := range cases {
		err := tc.tx.ValidateBasic()
		if tc.valid {
			require.Nil(t, err, fmt.Sprintf("transaction [no: %d] [%v] failed but was supposed to be valid", i, tc.tx))
		} else {
			require.NotNil(t, err, fmt.Sprintf("transaction [no: %d] [%v] is valid but is supposed to have an error", i, tc.tx))
		}
	}
}

// Tests

func TestMsgIssueToken(t *testing.T) {
	var (
		name                 = "Zap"
		originalSymbol       = "ZAP"
		symbol               = "zap"
		total          int64 = 1
		max            int64 = 10
		owner                = sdk.AccAddress([]byte("me"))
		nominee              = sdk.AccAddress([]byte("nominee"))
		msg                  = NewMsgIssueToken(nominee, owner, name, symbol, originalSymbol, sdk.NewInt(total), sdk.NewInt(max), false)
	)

	require.Equal(t, msg.Route(), RouterKey)
	require.Equal(t, msg.Type(), "issue_token")
}

func TestMsgIssueTokenValidation(t *testing.T) {
	var (
		name                 = "Zap"
		originalSymbol       = "ZAP"
		symbol               = "zap"
		total          int64 = 1
		max            int64 = 10
		totalInvalid   int64 = 0
		maxInvalid     int64 = 0
		acc                  = sdk.AccAddress([]byte("me"))
		nominee              = sdk.AccAddress([]byte("nominee"))
		name2                = "a"
		total2         int64 = 2
		max2           int64 = 20
		acc2                 = sdk.AccAddress([]byte("you"))
		nominee2             = sdk.AccAddress([]byte("nominee2"))
	)

	cases := []struct {
		valid bool
		tx    MsgInterface
	}{
		{true, NewMsgIssueToken(nominee, acc, name, symbol, originalSymbol, sdk.NewInt(total), sdk.NewInt(max), false)},
		{true, NewMsgIssueToken(nominee, acc, name, symbol, originalSymbol, sdk.NewInt(total), sdk.NewInt(max), false)},
		{false, NewMsgIssueToken(nominee, acc, name, symbol, originalSymbol, sdk.NewInt(totalInvalid), sdk.NewInt(maxInvalid), false)},
		{true, NewMsgIssueToken(nominee2, acc2, name2, symbol, originalSymbol, sdk.NewInt(total2), sdk.NewInt(max2), false)},
		{true, NewMsgIssueToken(nominee2, acc2, name2, symbol, originalSymbol, sdk.NewInt(total), sdk.NewInt(max), false)},
		{true, NewMsgIssueToken(nominee, acc, name2, symbol, originalSymbol, sdk.NewInt(total2), sdk.NewInt(max2), false)},
		{false, NewMsgIssueToken(nominee, nil, name, symbol, originalSymbol, sdk.NewInt(total2), sdk.NewInt(max2), false)},
		{false, NewMsgIssueToken(nominee2, acc2, "", symbol, originalSymbol, sdk.NewInt(total2), sdk.NewInt(max2), false)},
		{false, NewMsgIssueToken(nominee2, acc2, name, symbol, originalSymbol, sdk.NewInt(totalInvalid), sdk.NewInt(maxInvalid), false)},
	}

	validateError(cases, t)
}

func TestMsgIssueTokenGetSignBytes(t *testing.T) {
	var (
		name                 = "Zap"
		originalSymbol       = "ZAP"
		symbol               = "zap"
		total          int64 = 1
		max            int64 = 10
		owner                = sdk.AccAddress([]byte("me"))
		nominee              = sdk.AccAddress([]byte("nominee"))
		msg                  = NewMsgIssueToken(nominee, owner, name, symbol, originalSymbol, sdk.NewInt(total), sdk.NewInt(max), false)
	)
	actual := msg.GetSignBytes()

	expected := `{"type":"denominations/MsgIssueToken",` +
		`"value":{` +
		`"max_supply":"10",` +
		`"mintable":false,` +
		`"name":"Zap",` +
		`"original_symbol":"ZAP",` +
		`"owner":"cosmos1d4js690r9j",` +
		`"source_address":"cosmos1dehk66twv4js5dq8xr",` +
		`"symbol":"` + symbol + `",` +
		`"total_supply":"1"}}`

	require.Equal(t, expected, string(actual))
}

func TestMsgMintCoins(t *testing.T) {
	var (
		amount int64 = 10
		symbol       = "ZAP-001"
		owner        = sdk.AccAddress([]byte("me"))
		msg          = NewMsgMintCoins(sdk.NewInt(amount), symbol, owner)
	)

	require.Equal(t, msg.Route(), RouterKey)
	require.Equal(t, msg.Type(), "mint_coins")
}

func TestMsgMintCoinsValidation(t *testing.T) {
	var (
		amount  int64 = 10
		symbol        = "ZAP-001"
		symbol2       = "MNT-500"
		owner         = sdk.AccAddress([]byte("me"))
		owner2        = sdk.AccAddress([]byte("you"))
	)

	cases := []struct {
		valid bool
		tx    MsgInterface
	}{
		{true, NewMsgMintCoins(sdk.NewInt(amount), symbol, owner)},
		{true, NewMsgMintCoins(sdk.NewInt(amount), symbol2, owner2)},
		{false, NewMsgMintCoins(sdk.NewInt(-1), symbol, owner)},
		{false, NewMsgMintCoins(sdk.NewInt(0), symbol, owner)},
		{true, NewMsgMintCoins(sdk.NewInt(1), symbol, owner)},
		{false, NewMsgMintCoins(sdk.NewInt(amount), symbol, nil)},
		{false, NewMsgMintCoins(sdk.NewInt(amount), "", owner)},
	}

	validateError(cases, t)
}

func TestMsgMintCoinsGetSignBytes(t *testing.T) {
	var (
		amount int64 = 10
		symbol       = "ZAP-001"
		owner        = sdk.AccAddress([]byte("me"))
		msg          = NewMsgMintCoins(sdk.NewInt(amount), symbol, owner)
	)
	actual := msg.GetSignBytes()

	expected := `{"type":"denominations/MsgMintCoins","value":{` +
		`"amount":"10",` +
		`"owner":"cosmos1d4js690r9j",` +
		`"symbol":"ZAP-001"}}`

	require.Equal(t, expected, string(actual))
}

func TestMsgBurnCoins(t *testing.T) {
	var (
		amount int64 = 10
		symbol       = "ZAP-001"
		owner        = sdk.AccAddress([]byte("me"))
		msg          = NewMsgBurnCoins(sdk.NewInt(amount), symbol, owner)
	)

	require.Equal(t, msg.Route(), RouterKey)
	require.Equal(t, msg.Type(), "burn_coins")
}

func TestMsgBurnCoinsValidation(t *testing.T) {
	var (
		amount  int64 = 20
		symbol        = "ZAP-001"
		symbol2       = "BRN-002"
		owner         = sdk.AccAddress([]byte("me"))
		owner2        = sdk.AccAddress([]byte("you"))
	)

	cases := []struct {
		valid bool
		tx    MsgInterface
	}{
		{true, NewMsgBurnCoins(sdk.NewInt(amount), symbol, owner)},
		{true, NewMsgBurnCoins(sdk.NewInt(amount), symbol2, owner2)},
		{false, NewMsgBurnCoins(sdk.NewInt(-1), symbol, owner)},
		{false, NewMsgBurnCoins(sdk.NewInt(0), symbol, owner)},
		{true, NewMsgBurnCoins(sdk.NewInt(1), symbol, owner)},
		{false, NewMsgBurnCoins(sdk.NewInt(amount), symbol, nil)},
		{false, NewMsgBurnCoins(sdk.NewInt(amount), "", owner)},
	}

	validateError(cases, t)
}

func TestMsgBurnCoinsGetSignBytes(t *testing.T) {
	var (
		amount int64 = 100
		symbol       = "ZAP-999"
		owner        = sdk.AccAddress([]byte("me"))
		msg          = NewMsgBurnCoins(sdk.NewInt(amount), symbol, owner)
	)
	actual := msg.GetSignBytes()

	expected := `{"type":"denominations/MsgBurnCoins","value":{` +
		`"amount":"100",` +
		`"owner":"cosmos1d4js690r9j",` +
		`"symbol":"ZAP-999"}}`

	require.Equal(t, expected, string(actual))
}

func TestMsgFreezeCoins(t *testing.T) {
	var (
		amount int64 = 10
		symbol       = "ZAP-001"
		owner        = sdk.AccAddress([]byte("me"))
		msg          = NewMsgFreezeCoins(sdk.NewInt(amount), symbol, owner, owner)
	)

	require.Equal(t, msg.Route(), RouterKey)
	require.Equal(t, msg.Type(), "freeze_coins")
}

func TestMsgFreezeCoinsValidation(t *testing.T) {
	var (
		amount  int64 = 15
		symbol        = "ZAP-001"
		symbol2       = "FRZ-112"
		owner         = sdk.AccAddress([]byte("me"))
		owner2        = sdk.AccAddress([]byte("you"))
	)

	cases := []struct {
		valid bool
		tx    MsgInterface
	}{
		{true, NewMsgFreezeCoins(sdk.NewInt(amount), symbol, owner, owner)},
		{true, NewMsgFreezeCoins(sdk.NewInt(amount), symbol2, owner2, owner2)},
		{false, NewMsgFreezeCoins(sdk.NewInt(-1), symbol, owner, owner)},
		{false, NewMsgFreezeCoins(sdk.NewInt(0), symbol, owner, owner)},
		{true, NewMsgFreezeCoins(sdk.NewInt(1), symbol, owner, owner)},
		{false, NewMsgFreezeCoins(sdk.NewInt(amount), symbol, nil, nil)},
		{false, NewMsgFreezeCoins(sdk.NewInt(amount), "", owner, owner)},
	}

	validateError(cases, t)
}

func TestMsgFreezeCoinsGetSignBytes(t *testing.T) {
	var (
		amount int64 = 100
		symbol       = "FRZ-999"
		owner        = sdk.AccAddress([]byte("me"))
		msg          = NewMsgFreezeCoins(sdk.NewInt(amount), symbol, owner, owner)
	)
	actual := msg.GetSignBytes()

	expected := `{"type":"denominations/MsgFreezeCoins","value":{` +
		`"address":"cosmos1d4js690r9j",` +
		`"amount":"100",` +
		`"owner":"cosmos1d4js690r9j",` +
		`"symbol":"FRZ-999"}}`

	require.Equal(t, expected, string(actual))
}

func TestMsgUnfreezeCoins(t *testing.T) {
	var (
		amount int64 = 10
		symbol       = "UFZ-001"
		owner        = sdk.AccAddress([]byte("me"))
		msg          = NewMsgUnfreezeCoins(sdk.NewInt(amount), symbol, owner, owner)
	)

	require.Equal(t, msg.Route(), RouterKey)
	require.Equal(t, msg.Type(), "unfreeze_coins")
}

func TestMsgUnfreezeCoinsValidation(t *testing.T) {
	var (
		amount  int64 = 15
		symbol        = "ZAP-001"
		symbol2       = "UFZ-130"
		owner         = sdk.AccAddress([]byte("me"))
		owner2        = sdk.AccAddress([]byte("you"))
	)

	cases := []struct {
		valid bool
		tx    MsgInterface
	}{
		{true, NewMsgUnfreezeCoins(sdk.NewInt(amount), symbol, owner, owner)},
		{true, NewMsgUnfreezeCoins(sdk.NewInt(amount), symbol2, owner2, owner2)},
		{false, NewMsgUnfreezeCoins(sdk.NewInt(-1), symbol, owner, owner)},
		{false, NewMsgUnfreezeCoins(sdk.NewInt(0), symbol, owner, owner)},
		{true, NewMsgUnfreezeCoins(sdk.NewInt(1), symbol, owner, owner)},
		{false, NewMsgUnfreezeCoins(sdk.NewInt(amount), symbol, nil, nil)},
		{false, NewMsgUnfreezeCoins(sdk.NewInt(amount), "", owner, owner)},
	}

	validateError(cases, t)
}

func TestMsgUnfreezeCoinsGetSignBytes(t *testing.T) {
	var (
		amount int64 = 100
		symbol       = "UFZ-999"
		owner        = sdk.AccAddress([]byte("me"))
		msg          = NewMsgUnfreezeCoins(sdk.NewInt(amount), symbol, owner, owner)
	)
	actual := msg.GetSignBytes()

	expected := `{"type":"denominations/MsgUnfreezeCoins","value":{` +
		`"address":"cosmos1d4js690r9j",` +
		`"amount":"100",` +
		`"owner":"cosmos1d4js690r9j",` +
		`"symbol":"UFZ-999"}}`

	require.Equal(t, expected, string(actual))
}
