package cli

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/version"

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
		NewQueryParamsCmd(),
		NewQueryAuctionsCmd(),
		NewQueryAuctionCmd(),
		NewQueryBidsCmd(),
		NewQueryVestingsCmd(),
	)

	return cmd
}

func NewQueryParamsCmd() *cobra.Command {
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

func NewQueryAuctionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auctions",
		Args:  cobra.NoArgs,
		Short: "Query for all auctions",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query for all auctions with the given optional flags.
Example:
$ %s query %s auctions
$ %s query %s auctions --status AUCTION_STATUS_STANDBY
$ %s query %s auctions --type AUCTION_TYPE_FIXED_PRICE

Auction statuses: AUCTION_STATUS_STANDBY, AUCTION_STATUS_STARTED, AUCTION_STATUS_VESTING, AUCTION_STATUS_FINISHED, and AUCTION_STATUS_CANCELLED
Auction types: AUCTION_TYPE_FIXED_PRICE and AUCTION_TYPE_ENGLISH
`,
				version.AppName, types.ModuleName,
				version.AppName, types.ModuleName,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			status, _ := cmd.Flags().GetString(FlagAuctionStatus)
			typ, _ := cmd.Flags().GetString(FlagAuctionType)

			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			req := &types.QueryAuctionsRequest{
				Status:     status,
				Type:       typ,
				Pagination: pageReq,
			}

			resp, err := queryClient.Auctions(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(resp)
		},
	}

	cmd.Flags().AddFlagSet(flagSetAuctions())
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func NewQueryAuctionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auction [auction-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query a specific auction",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details about a speicif auction.
Example:
$ %s query %s auction 1
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			auctionId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "auction-id %s is not valid", args[0])
			}

			queryClient := types.NewQueryClient(clientCtx)

			resp, err := queryClient.Auction(cmd.Context(), &types.QueryAuctionRequest{
				AuctionId: auctionId,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(resp)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func NewQueryBidsCmd() *cobra.Command {
	bech32PrefixAccAddr := sdk.GetConfig().GetBech32AccountAddrPrefix()

	cmd := &cobra.Command{
		Use:   "bids [auction-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query all bids for the auction",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all bids for the auction.
Example:
$ %s query %s bids 1
$ %s query %s bids 1 --bidder %s1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj
$ %s query %s bids 1 --winner 
`,
				version.AppName, types.ModuleName,
				version.AppName, types.ModuleName, bech32PrefixAccAddr,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			bidderAddr, _ := cmd.Flags().GetString(FlagBidderAddr)
			eligible, _ := cmd.Flags().GetString(FlagIsWinner)

			auctionId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "auction-id %s is not valid", args[0])
			}

			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			req := &types.QueryBidsRequest{
				AuctionId:  auctionId,
				Bidder:     bidderAddr,
				IsWinner:   eligible,
				Pagination: pageReq,
			}

			resp, err := queryClient.Bids(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(resp)
		},
	}

	cmd.Flags().AddFlagSet(flagSetBids())
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func NewQueryVestingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vestings [auction-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query vesting schedules about the auction",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query vesting schedules for the auction.
Example:
$ %s query %s vestings 1
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			auctionId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "auction-id %s is not valid", args[0])
			}

			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryVestingsRequest{
				AuctionId: auctionId,
			}

			resp, err := queryClient.Vestings(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(resp)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
