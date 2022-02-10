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
		NewCreateFixedPriceAuction(),
		NewCreateEnglishAuction(),
		NewCancelAuction(),
		NewPlaceBid(),
	)
	if keeper.EnableAddAllowedBidder {
		cmd.AddCommand(NewAddAllowedBidderCmd())
	}
	return cmd
}

func NewCreateFixedPriceAuction() *cobra.Command {
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

[start_price]: starting price of the selling coin proportional to the paying coin
[selling_coin]: selling amount of coin for the auction
[paying_coin_denom]: paying coin denom that bidders need to bid for
[vesting_schedules]: vesting schedules that release the paying amount of coins to the autioneer
[start_time]: start time of the auction
[end_time]: end time of the auction
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

func NewCreateEnglishAuction() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-english-auction [file]",
		Args:  cobra.ExactArgs(1),
		Short: "Create a english auction",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Create a english auction.
The auction details must be provided through a JSON file. 
		
Example:
$ %s tx %s create-english-auction <path/to/auction.json> --from mykey 

Where auction.json contains:

{}
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			fmt.Println(clientCtx)

			// TODO: not implemented yet

			return nil
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewCancelAuction() *cobra.Command {
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

func NewPlaceBid() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bid [auction-id] [price] [coin]",
		Args:  cobra.ExactArgs(3),
		Short: "Bid for the auction",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Bid for the auction with what price and amount of coin you want to bid for. 
		
Example:
$ %s tx %s bid 1 1.0 100000000denom2--from mykey 

Note that [price] argument specifies the price of the selling coin. For a fixed price auction, you must use the same start price of the auction.
For an english auction, it is up to you for how much price you want to bid for. Moreover, you must have sufficient balance of the paying coin denom
in order to bid for the amount of coin you bid for the auction.
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

			price, err := sdk.NewDecFromStr(args[1])
			if err != nil {
				return err
			}

			coin, err := sdk.ParseCoinNormalized(args[2])
			if err != nil {
				return err
			}

			msg := types.NewMsgPlaceBid(
				auctionId,
				clientCtx.GetFromAddress().String(),
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
		Use:   "add-allowed-bidder [auction-id] [max-bid-amount]",
		Args:  cobra.ExactArgs(2),
		Short: "Add an allowed bidder for the auction",
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

			maxBidAmt, ok := sdk.NewIntFromString(args[1])
			if !ok {
				return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "maxium bid price must be a positive integer")
			}

			msg := types.NewAddAllowedBidder(
				auctionId,
				types.AllowedBidder{
					Bidder:       clientCtx.GetFromAddress().String(),
					MaxBidAmount: maxBidAmt,
				},
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
