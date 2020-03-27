package cli

import (
	"fmt"
	"github.com/xar-network/xar-network/x/bonds"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/xar-network/xar-network/x/governors/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	txQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the governors module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txQueryCmd.AddCommand(client.GetCommands(
		GetCmdQueryDeposit(queryRoute, cdc),
		GetCmdQueryUnbondingDeposit(queryRoute, cdc),
		GetCmdQueryUnbondingDeposits(queryRoute, cdc),
		GetCmdQueryParams(queryRoute, cdc),
		GetCmdQueryPool(queryRoute, cdc))...)

	return txQueryCmd

}

// GetCmdQueryDeposit the query deposit command.
func GetCmdQueryDeposit(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "deposit [depositor-addr]",
		Short: "Query a deposit based on address",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query deposits for an individual depositor.

Example:
$ %s query governors deposit cosmos1gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			depAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			bz, err := cdc.MarshalJSON(bonds.NewQueryBondsParams(depAddr))
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, bonds.QueryDeposit)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var resp bonds.Deposit
			if err := cdc.UnmarshalJSON(res, &resp); err != nil {
				return err
			}

			return cliCtx.PrintOutput(resp)
		},
	}
}

// GetCmdQueryUnbondingDeposit implements the command to query a single
// unbonding-deposit record.
func GetCmdQueryUnbondingDeposit(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "unbonding-deposit [depositor-addr]",
		Short: "Query an unbonding-deposit record based on depositor",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query unbonding deposits for an individual depositor.

Example:
$ %s query governors unbonding-deposit cosmos1gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			depAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			bz, err := cdc.MarshalJSON(bonds.NewQueryBondsParams(depAddr))
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, bonds.QueryUnbondingDeposit)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var ubd bonds.UnbondingDeposit
			if err = cdc.UnmarshalJSON(res, &ubd); err != nil {
				return err
			}

			return cliCtx.PrintOutput(ubd)
		},
	}
}

// GetCmdQueryUnbondingDeposits implements the command to query all the
// unbonding-deposit records for a depositor.
func GetCmdQueryUnbondingDeposits(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "unbonding-deposits [depositor-addr]",
		Short: "Query all unbonding-deposits records for one depositor",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query unbonding deposits for an individual depositor.

Example:
$ %s query governors unbonding-deposit cosmos1gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			depositorAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			bz, err := cdc.MarshalJSON(bonds.NewQueryDepositorParams(depositorAddr))
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, bonds.QueryDepositorUnbondingDeposits)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var ubds bonds.UnbondingDeposits
			if err = cdc.UnmarshalJSON(res, &ubds); err != nil {
				return err
			}

			return cliCtx.PrintOutput(ubds)
		},
	}
}

// GetCmdQueryPool implements the pool query command.
func GetCmdQueryPool(storeName string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "pool",
		Args:  cobra.NoArgs,
		Short: "Query the current governors pool values",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query values for amounts stored in the governors pool.

Example:
$ %s query governors pool
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			bz, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/pool", storeName), nil)
			if err != nil {
				return err
			}

			var pool bonds.Pool
			if err := cdc.UnmarshalJSON(bz, &pool); err != nil {
				return err
			}

			return cliCtx.PrintOutput(pool)
		},
	}
}

// GetCmdQueryParams implements the params query command.
func GetCmdQueryParams(storeName string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query the current governors parameters information",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query values set as governors parameters.

Example:
$ %s query governors params
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", storeName, bonds.QueryParameters)
			bz, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var params types.Params
			cdc.MustUnmarshalJSON(bz, &params)
			return cliCtx.PrintOutput(params)
		},
	}
}
