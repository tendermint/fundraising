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

		for _, auction := range k.GetAuctions(ctx) {
			if auction.GetStatus() == types.AuctionStatusStarted {
				sellingPoolAcc := auction.GetSellingReserveAddress()
				sellingReserve := k.bankKeeper.GetBalance(ctx, sellingPoolAcc, auction.GetSellingCoin().Denom)
				if !sellingReserve.IsGTE(auction.GetSellingCoin()) {
					msg += fmt.Sprintf("\tselling reserve balance %s\n"+
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
// must equal or greater than the paying reserve account balance.
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

			payingPoolAcc := auction.GetPayingReserveAddress()
			payingReserve := k.bankKeeper.GetBalance(ctx, payingPoolAcc, auction.GetPayingCoinDenom())
			if !payingReserve.IsGTE(totalBidCoin) {
				msg += fmt.Sprintf("\tpaying reserve balance %s\n"+
					"\tpaying pool reserve: %v\n"+
					"\ttotal bid coin: %v",
					payingPoolAcc.String(), payingReserve, totalBidCoin)
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

		for _, auction := range k.GetAuctions(ctx) {
			totalPayingCoin := sdk.NewCoin(auction.GetPayingCoinDenom(), sdk.ZeroInt())

			if auction.GetStatus() == types.AuctionStatusVesting {
				for _, queue := range k.GetVestingQueuesByAuctionId(ctx, auction.GetId()) {
					if !queue.Released {
						totalPayingCoin = totalPayingCoin.Add(queue.PayingCoin)
					}
				}
			}

			vestingPoolAcc := auction.GetVestingReserveAddress()
			vestingReserve := k.bankKeeper.GetBalance(ctx, vestingPoolAcc, auction.GetPayingCoinDenom())
			if !vestingReserve.IsGTE(totalPayingCoin) {
				msg += fmt.Sprintf("\tvesting reserve balance %s\n"+
					"\tvesting pool reserve: %v\n"+
					"\ttotal paying coin: %v",
					vestingPoolAcc.String(), vestingReserve, totalPayingCoin)
				count++
			}
		}
		broken := count != 0

		return sdk.FormatInvariant(types.ModuleName, "vesting pool reserve amount and total paying amount", msg), broken
	}
}

// AuctionStatusInvariant checks an invariant
func AuctionStatusInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {

		// TODO: depending on auction status, reserve account must not have any balance (they can have some balances due to transfer)
		// for _, auction := range k.GetAuctions(ctx) {

		// }

		msg := ""
		broken := false
		return sdk.FormatInvariant(types.ModuleName, "vesting pool reserve amount and total paying amount", msg), broken
	}
}
