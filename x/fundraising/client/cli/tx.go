package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/tendermint/fundraising/x/fundraising/keeper"
	"github.com/tendermint/fundraising/x/fundraising/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		NewCreateFixedPriceAuctionCmd(),
		NewCreateBatchAuctionCmd(),
		NewCancelAuctionCmd(),
		NewPlaceBidCmd(),
		NewModifyBidCmd(),
	)
	if keeper.EnableAddAllowedBidder {
		cmd.AddCommand(NewAddAllowedBidderCmd())
	}
	return cmd
}

func NewCreateFixedPriceAuctionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-fixed-price-auction [file]",
		Args:  cobra.ExactArgs(1),
		Short: "Create a fixed price auction",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Create a fixed price auction.
The auction details must be provided through a JSON file. 
		
Example:
$ %s tx %s create-fixed-price-auction <path/to/auction.json> --from mykey 

Where auction.json contains:

{
  "start_price": "1.000000000000000000",
  "min_bid_price": "0.100000000000000000",
  "selling_coin": {
    "denom": "denom1",
    "amount": "1000000000000"
  },
  "paying_coin_denom": "denom2",
  "vesting_schedules": [
    {
      "release_time": "2022-01-01T00:00:00Z",
      "weight": "0.500000000000000000"
    },
    {
      "release_time": "2022-06-01T00:00:00Z",
      "weight": "0.250000000000000000"
    },
    {
      "release_time": "2022-12-01T00:00:00Z",
      "weight": "0.250000000000000000"
    }
  ],
  "start_time": "2021-11-01T00:00:00Z",
  "end_time": "2021-12-01T00:00:00Z"
}

Description of the parameters:

[start_price]: the start price of the selling coin that is proportional to the paying coin denom 
[selling_coin]: the selling amount of coin for the auction
[paying_coin_denom]: the paying coin denom that the auctioneer wants to exchange with
[vesting_schedules]: the vesting schedules that release the paying coins to the autioneer
[start_time]: the start time of the auction
[end_time]: the end time of the auction
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			auction, err := ParseFixedPriceAuctionRequest(args[0])
			if err != nil {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "failed to parse %s file due to %v", args[0], err)
			}

			msg := types.NewMsgCreateFixedPriceAuction(
				clientCtx.GetFromAddress().String(),
				auction.StartPrice,
				auction.SellingCoin,
				auction.PayingCoinDenom,
				auction.VestingSchedules,
				auction.StartTime,
				auction.EndTime,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewCreateBatchAuctionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-batch-auction [file]",
		Args:  cobra.ExactArgs(1),
		Short: "Create a batch auction",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Create a batch auction.
The auction details must be provided through a JSON file. 
		
Example:
$ %s tx %s create-batch-auction <path/to/auction.json> --from mykey 

Where auction.json contains:
{
  "start_price": "0.500000000000000000",
  "min_bid_price": "0.100000000000000000",
  "selling_coin": {
    "denom": "denom1",
    "amount": "1000000000000"
  },
  "paying_coin_denom": "denom2",
  "vesting_schedules": [
    {
      "release_time": "2023-06-01T00:00:00Z",
      "weight": "0.500000000000000000"
    },
    {
      "release_time": "2023-12-01T00:00:00Z",
      "weight": "0.500000000000000000"
    }
  ],
  "max_extended_round": 2,
  "extended_round_rate": "0.150000000000000000",
  "start_time": "2022-02-01T00:00:00Z",
  "end_time": "2022-06-20T00:00:00Z"
}

Description of the parameters:

[start_price]: the start price of the selling coin that is proportional to the paying coin denom 
[selling_coin]: the selling amount of coin for the auction
[paying_coin_denom]: the paying coin denom that the auctioneer wants to exchange with
[vesting_schedules]: the vesting schedules that release the paying coins to the autioneer
[max_extended_round]: the number of extended rounds
[extended_round_rate]: the rate that determines if the auction needs to run another round
[start_time]: the start time of the auction
[end_time]: the end time of the auction
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			auction, err := ParseBatchAuctionRequest(args[0])
			if err != nil {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "failed to parse %s file due to %v", args[0], err)
			}

			msg := types.NewMsgCreateBatchAuction(
				clientCtx.GetFromAddress().String(),
				auction.StartPrice,
				auction.MinBidPrice,
				auction.SellingCoin,
				auction.PayingCoinDenom,
				auction.VestingSchedules,
				auction.MaxExtendedRound,
				auction.ExtendedRoundRate,
				auction.StartTime,
				auction.EndTime,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewCancelAuctionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel [auction-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Cancel the auction",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Cancel the auction with the id. 
		
Example:
$ %s tx %s cancel 1 --from mykey 
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
				return err
			}

			msg := types.NewMsgCancelAuction(
				clientCtx.GetFromAddress().String(),
				auctionId,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewPlaceBidCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bid [auction-id] [bid-type] [price] [coin]",
		Args:  cobra.ExactArgs(4),
		Short: "Bid for the auction",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Bid for the auction with what price and amount of coin you want to bid for. 

Bid Type Options:
1. fixed-price (fp or f)
2. batch-worth (bw or w) 
3. batch-many  (bm or m)

Example:
$ %s tx %s bid 1 fixed-price 0.55 100000000denom2 --from mykey 
$ %s tx %s bid 1 batch-worth 0.55 100000000denom2 --from mykey 
$ %s tx %s bid 1 batch-many 0.55 100000000denom1 --from mykey 
$ %s tx %s bid 1 fp 0.55 100000000denom2 --from mykey 
$ %s tx %s bid 1 bw 0.55 100000000denom2 --from mykey 
$ %s tx %s bid 1 bm 0.55 100000000denom1 --from mykey 

Note:
In case of placing a bid for a fixed price auction, you must provide [price] argument with the same price of the auction. 
In case of placing a bid for a batch auction, there are two bid type options; batch-worth and batch-many, which you can find more information
in our technical spec docs. https://github.com/tendermint/fundraising/blob/main/x/fundraising/spec/01_concepts.md
`,
				version.AppName, types.ModuleName,
				version.AppName, types.ModuleName,
				version.AppName, types.ModuleName,
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

			auctionId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			bidType, err := parseBidType(args[1])
			if err != nil {
				return fmt.Errorf("parse order direction: %w", err)
			}

			price, err := sdk.NewDecFromStr(args[2])
			if err != nil {
				return err
			}

			coin, err := sdk.ParseCoinNormalized(args[3])
			if err != nil {
				return err
			}

			msg := types.NewMsgPlaceBid(
				auctionId,
				clientCtx.GetFromAddress().String(),
				bidType,
				price,
				coin,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewModifyBidCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "modify-bid [auction-id] [bid-id] [price] [coin]",
		Args:  cobra.ExactArgs(4),
		Short: "Modify the bid",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Modify the bid with new price and coin.
Either price or coin must be higher than the existing bid.

Example:
$ %s tx %s bid 1 1 1.0 100000000denom2 --from mykey
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
				return err
			}

			bidId, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			price, err := sdk.NewDecFromStr(args[2])
			if err != nil {
				return err
			}

			coin, err := sdk.ParseCoinNormalized(args[3])
			if err != nil {
				return err
			}

			msg := types.NewMsgModifyBid(
				auctionId,
				clientCtx.GetFromAddress().String(),
				bidId,
				price,
				coin,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewAddAllowedBidderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-allowed-bidder [auction-id] [bidder] [max-bid-amount]",
		Args:  cobra.ExactArgs(3),
		Short: "(Testing) Add an allowed bidder for the auction",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Add an allowed bidder for the auction.
This message is available for testing purpose and it is only accessible when you build the binary with testing mode.
		
Example:
$ %s tx %s add-allowed-bidder 1 10000000000 --from mykey 
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
				return err
			}

			bidderAddr, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			maxBidAmt, ok := sdk.NewIntFromString(args[2])
			if !ok {
				return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "maxium bid price must be a positive integer")
			}

			msg := types.NewAddAllowedBidder(
				auctionId,
				types.AllowedBidder{
					Bidder:       bidderAddr.String(),
					MaxBidAmount: maxBidAmt,
				},
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
