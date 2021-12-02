package cli

import (
	"context"
	"fmt"
	"strings"

	// "strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"

	// "github.com/cosmos/cosmos-sdk/client/flags"
	// sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group fundraising queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	// this line is used by starport scaffolding # 1
	cmd.AddCommand(
		QueryParams(),
		QueryAuctions(),
		QueryAuction(),
		QueryBids(),
		QueryVestings(),
	)

	return cmd
}

func QueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query the current fundraising parameters information",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query values set as fundraising parameters.
Example:
$ %s query %s params
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			resp, err := queryClient.Params(context.Background(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&resp.Params)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func QueryAuctions() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auctions",
		Args:  cobra.NoArgs,
		Short: "Query the current fundraising parameters information",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query values set as fundraising parameters.
Example:
$ %s query %s params
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			fmt.Println(queryClient)

			// TODO: not implemented yet

			// return clientCtx.PrintProto(&resp.Params)
			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func QueryAuction() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auction",
		Args:  cobra.NoArgs,
		Short: "Query the current fundraising parameters information",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query values set as fundraising parameters.
Example:
$ %s query %s params
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			fmt.Println(queryClient)

			// TODO: not implemented yet

			// return clientCtx.PrintProto(&resp.Params)
			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func QueryBids() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bids",
		Args:  cobra.NoArgs,
		Short: "Query the current fundraising parameters information",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query values set as fundraising parameters.
Example:
$ %s query %s params
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			fmt.Println(queryClient)

			// TODO: not implemented yet

			// return clientCtx.PrintProto(&resp.Params)
			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func QueryVestings() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vestings",
		Args:  cobra.NoArgs,
		Short: "Query the current fundraising parameters information",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query values set as fundraising parameters.
Example:
$ %s query %s params
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			fmt.Println(queryClient)

			// TODO: not implemented yet

			// return clientCtx.PrintProto(&resp.Params)
			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
