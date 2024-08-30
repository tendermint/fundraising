package keeper

import (
	"fmt"

	"cosmossdk.io/math"
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
		for _, inv := range []func(Keeper) sdk.Invariant{
			SellingPoolReserveAmountInvariant,
			PayingPoolReserveAmountInvariant,
			VestingPoolReserveAmountInvariant,
		} {
			res, stop := inv(k)(ctx)
			if stop {
				return res, stop
			}
		}
		return "", false
	}
}

// SellingPoolReserveAmountInvariant checks an invariant that the total amount of selling coin for an auction
// must equal or greater than the selling reserve account balance.
func SellingPoolReserveAmountInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		msg := ""
		count := 0

		auctions, err := k.Auctions(ctx)
		if err != nil {
			return "", false
		}
		for _, auction := range auctions {
			if auction.GetStatus() == types.AuctionStatusStarted {
				sellingReserveAddr := auction.GetSellingReserveAddress()
				sellingCoinDenom := auction.GetSellingCoin().Denom
				spendable := k.bankKeeper.SpendableCoins(ctx, sellingReserveAddr)
				sellingReserve := sdk.NewCoin(sellingCoinDenom, spendable.AmountOf(sellingCoinDenom))
				fmt.Println("sellingReserve: ", sellingReserve)
				fmt.Println("auction.GetSellingCoin(): ", auction.GetSellingCoin())
				if !sellingReserve.IsGTE(auction.GetSellingCoin()) {
					msg += fmt.Sprintf("\tselling reserve balance %s\n"+
						"\tselling pool reserve: %v\n"+
						"\ttotal selling coin: %v\n",
						sellingReserveAddr.String(), sellingReserve, auction.GetSellingCoin())
					count++
				}
			}
		}
		broken := count != 0

		return sdk.FormatInvariant(types.ModuleName, "selling pool reserve amount and selling coin amount", msg), broken
	}
}

// PayingPoolReserveAmountInvariant checks an invariant that the total bid amount
// must equal or greater than the paying reserve account balance.
func PayingPoolReserveAmountInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		msg := ""
		count := 0

		auctions, err := k.Auctions(ctx)
		if err != nil {
			return "", false
		}
		for _, auction := range auctions {
			totalBidCoin := sdk.NewCoin(auction.GetPayingCoinDenom(), math.ZeroInt())

			bids, err := k.GetBidsByAuctionId(ctx, auction.GetId())
			if err != nil {
				return "", false
			}
			if auction.GetStatus() == types.AuctionStatusStarted {
				for _, bid := range bids {
					bidAmt := bid.ConvertToPayingAmount(auction.GetPayingCoinDenom())
					totalBidCoin = totalBidCoin.Add(sdk.NewCoin(auction.GetPayingCoinDenom(), bidAmt))
				}
			}

			payingReserveAddr := auction.GetPayingReserveAddress()
			payingCoinDenom := auction.GetPayingCoinDenom()
			spendable := k.bankKeeper.SpendableCoins(ctx, payingReserveAddr)
			payingReserve := sdk.NewCoin(payingCoinDenom, spendable.AmountOf(payingCoinDenom))
			if !payingReserve.IsGTE(totalBidCoin) {
				msg += fmt.Sprintf("\tpaying reserve balance %s\n"+
					"\tpaying pool reserve: %v\n"+
					"\ttotal bid coin: %v\n",
					payingReserveAddr.String(), payingReserve, totalBidCoin)
				count++
			}
		}
		broken := count != 0

		return sdk.FormatInvariant(types.ModuleName, "paying pool reserve amount and total bids amount", msg), broken
	}
}

// VestingPoolReserveAmountInvariant checks an invariant that the total vesting amount
// must be equal or greater than the vesting reserve account balance.
func VestingPoolReserveAmountInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		msg := ""
		count := 0

		auctions, err := k.Auctions(ctx)
		if err != nil {
			return "", false
		}
		for _, auction := range auctions {
			totalPayingCoin := sdk.NewCoin(auction.GetPayingCoinDenom(), math.ZeroInt())

			vestingQueues, err := k.GetVestingQueuesByAuctionId(ctx, auction.GetId())
			if err != nil {
				return "", false
			}
			if auction.GetStatus() == types.AuctionStatusVesting {
				for _, queue := range vestingQueues {
					if !queue.Released {
						totalPayingCoin = totalPayingCoin.Add(queue.PayingCoin)
					}
				}
			}

			vestingReserveAddr := auction.GetVestingReserveAddress()
			payingCoinDenom := auction.GetPayingCoinDenom()
			spendable := k.bankKeeper.SpendableCoins(ctx, vestingReserveAddr)
			vestingReserve := sdk.NewCoin(payingCoinDenom, spendable.AmountOf(payingCoinDenom))
			if !vestingReserve.IsGTE(totalPayingCoin) {
				msg += fmt.Sprintf("\tvesting reserve balance %s\n"+
					"\tvesting pool reserve: %v\n"+
					"\ttotal paying coin: %v\n",
					vestingReserveAddr.String(), vestingReserve, totalPayingCoin)
				count++
			}
		}
		broken := count != 0

		return sdk.FormatInvariant(types.ModuleName, "vesting pool reserve amount and total paying amount", msg), broken
	}
}
