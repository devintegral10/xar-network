package cli

import (
	"bufio"
	"fmt"
	"github.com/xar-network/xar-network/x/governors/types"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Governors transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(client.PostCommands(
		GetCmdBond(cdc),
		GetCmdUnbond(storeKey, cdc),
	)...)

	return txCmd
}

// GetCmdBond implements the bond command.
func GetCmdBond(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "bond [amount]",
		Args:  cobra.ExactArgs(1),
		Short: "Bond liquid tokens",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Bond an amount of liquid coins from your wallet.

Example:
$ %s tx governors bond 1000stake --from mykey
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			depAddr := cliCtx.GetFromAddress()
			amount, err := sdk.ParseCoin(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgBond(depAddr, amount)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdUnbond implements the unbond command.
func GetCmdUnbond(storeName string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "unbond [amount]",
		Short: "Unbond tokens",
		Args:  cobra.ExactArgs(1),
		Long: strings.TrimSpace(
			fmt.Sprintf(`Unbond an amount of bonded tokens.

Example:
$ %s tx governors unbond 100stake --from mykey
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			depAddr := cliCtx.GetFromAddress()

			amount, err := sdk.ParseCoin(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgUnbond(depAddr, amount)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
