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

	// this line is used by starport scaffolding # 1
	cmd.AddCommand(
		NewCreateFixedPriceAuction(),
		NewCreateEnglishAuction(),
		NewCancelAuction(),
		NewPlaceBid(),
	)
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

			auctionID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgCancelAuction(
				clientCtx.GetFromAddress().String(),
				auctionID,
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
$ %s tx %s bid 1 1.0 100000000ugdex--from mykey 

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

			auctionID, err := strconv.ParseUint(args[0], 10, 64)
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
				auctionID,
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
