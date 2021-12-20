package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/fundraising/x/fundraising/types"
)

// RegisterInvariants registers all fundraising invariants.
func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, "selling-pool-reserve-amount",
		SellingPoolReserveAmountInvariant(k))
	ir.RegisterRoute(types.ModuleName, "paying-pool-reserve-amount",
		PayingPoolReserveAmountInvariant(k))
	ir.RegisterRoute(types.ModuleName, "vesting-pool-reserve-amount",
		VestingPoolReserveAmountInvariant(k))
}

// AllInvariants runs all invariants of the fundraising module.
func AllInvariants(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		for _, inv := range []func(Keeper) sdk.Invariant{} {
			res, stop := inv(k)(ctx)
			if stop {
				return res, stop
			}
		}
		return "", false
	}
}

// SellingPoolReserveAmountInvariant checks an invariant that the total amount of selling coin for an auction
// must equal to the selling reserve account balance.
func SellingPoolReserveAmountInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		msg := ""
		count := 0

		for _, auction := range k.GetAuctions(ctx) {
			if auction.GetStatus() == types.AuctionStatusStarted {
				sellingPoolAcc := auction.GetSellingPoolAddress()
				sellingReserve := k.bankKeeper.GetBalance(ctx, sellingPoolAcc, auction.GetSellingCoin().Denom)
				if !sellingReserve.Equal(auction.GetSellingCoin()) {
					msg += fmt.Sprintf("\tselling reserve account %s\n"+
						"\tselling pool reserve: %v\n"+
						"\ttotal selling coin: %v",
						sellingPoolAcc.String(), sellingReserve, auction.GetSellingCoin())
					count++
				}
			}
		}
		broken := count != 0

		return sdk.FormatInvariant(types.ModuleName, "selling pool reserve amount and selling coin amount", msg), broken
	}
}

// PayingPoolReserveAmountInvariant checks an invariant that the total bid amount
// must equal to the paying reserve account balance.
func PayingPoolReserveAmountInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		msg := ""
		count := 0

		for _, auction := range k.GetAuctions(ctx) {
			totalBidCoin := sdk.NewCoin(auction.GetPayingCoinDenom(), sdk.ZeroInt())

			if auction.GetStatus() == types.AuctionStatusStarted {
				for _, bid := range k.GetBidsByAuctionId(ctx, auction.GetId()) {
					totalBidCoin = totalBidCoin.Add(bid.Coin)
				}
			}

			payingPoolAcc := auction.GetPayingPoolAddress()
			payingReserve := k.bankKeeper.GetBalance(ctx, payingPoolAcc, auction.GetPayingCoinDenom())
			fmt.Println("payingReserve: ", payingReserve)
			fmt.Println("totalBidCoin: ", totalBidCoin)
			if !payingReserve.Equal(totalBidCoin) {
				msg += fmt.Sprintf("\tpaying reserve account %s\n"+
					"\tpaying pool reserve: %v\n"+
					"\ttotal paying coin: %v",
					payingPoolAcc.String(), payingReserve, totalBidCoin)
				count++
			}
		}
		broken := count != 0

		return sdk.FormatInvariant(types.ModuleName, "paying pool reserve amount and all bids amount", msg), broken
	}
}

func VestingPoolReserveAmountInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var (
			msg    string
			broken bool
		)
		// TODO: not implemented yet
		return sdk.FormatInvariant(types.ModuleName, "paying pool reserve amount and selling coin amount", msg), broken
	}
}
