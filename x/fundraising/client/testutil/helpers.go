package testutil

// DONTCOVER

import (
	"fmt"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/testutil"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/fundraising/x/fundraising/client/cli"
)

var commonArgs = []string{
	fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
	fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
	fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10))).String()),
}

// MsgCreateFixedPriceAuctionExec creates a transaction for creating a fixed price auction.
func MsgCreateFixedPriceAuctionExec(clientCtx client.Context,
	from string,
	file string,
	extraArgs ...string,
) (testutil.BufferWriter, error) {

	args := append([]string{
		file,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}, commonArgs...)

	args = append(args, commonArgs...)
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, cli.NewCreateFixedPriceAuctionCmd(), args)
}

// MsgCreateBatchAuctionExec creates a transaction for creating a batch auction.
func MsgCreateBatchAuctionExec(clientCtx client.Context,
	from string,
	file string,
	extraArgs ...string,
) (testutil.BufferWriter, error) {

	args := append([]string{
		file,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}, commonArgs...)

	args = append(args, commonArgs...)
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, cli.NewCreateBatchAuctionCmd(), args)
}

// MsgPlaceBidExec creates a transaction for placing a bid.
func MsgPlaceBidExec(clientCtx client.Context,
	from string,
	auctionId uint64,
	bidType string,
	price sdk.Dec,
	coin sdk.Coin,
	extraArgs ...string,
) (testutil.BufferWriter, error) {

	args := append([]string{
		fmt.Sprint(auctionId),
		bidType,
		price.String(),
		coin.String(),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}, commonArgs...)

	args = append(args, commonArgs...)
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, cli.NewPlaceBidCmd(), args)
}

// MsgAddAllowedBidderExec creates a transaction for adding an allowed bidder.
func MsgAddAllowedBidderExec(clientCtx client.Context,
	from string,
	auctionId uint64,
	maxBidAmt math.Int,
	extraArgs ...string,
) (testutil.BufferWriter, error) {

	args := append([]string{
		fmt.Sprint(auctionId),
		from,
		maxBidAmt.String(),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}, commonArgs...)

	args = append(args, commonArgs...)
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, cli.NewAddAllowedBidderCmd(), args)
}
