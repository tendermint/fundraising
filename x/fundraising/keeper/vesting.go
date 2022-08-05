package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

// ApplyVestingSchedules stores vesting queues based on the vesting schedules of the auction and
// sets status to vesting.
func (k Keeper) ApplyVestingSchedules(ctx sdk.Context, auction types.AuctionI) error {
	payingReserveAddr := auction.GetPayingReserveAddress()
	vestingReserveAddr := auction.GetVestingReserveAddress()
	payingCoinDenom := auction.GetPayingCoinDenom()
	spendableCoins := k.bankKeeper.SpendableCoins(ctx, payingReserveAddr)
	reserveCoin := sdk.NewCoin(payingCoinDenom, spendableCoins.AmountOf(payingCoinDenom))

	vsLen := len(auction.GetVestingSchedules())
	if vsLen == 0 {
		// Send reserve coins to the auctioneer from the paying reserve account
		if err := k.bankKeeper.SendCoins(ctx, payingReserveAddr, auction.GetAuctioneer(), sdk.NewCoins(reserveCoin)); err != nil {
			return err
		}

		_ = auction.SetStatus(types.AuctionStatusFinished)
		k.SetAuction(ctx, auction)

	} else {
		// Move reserve coins from the paying reserve to the vesting reserve account
		if err := k.bankKeeper.SendCoins(ctx, payingReserveAddr, vestingReserveAddr, sdk.NewCoins(reserveCoin)); err != nil {
			return err
		}

		remaining := reserveCoin
		for i, schedule := range auction.GetVestingSchedules() {
			payingAmt := sdk.NewDecFromInt(reserveCoin.Amount).MulTruncate(schedule.Weight).TruncateInt()

			// All the remaining paying coin goes to the last vesting queue
			if i == vsLen-1 {
				payingAmt = remaining.Amount
			}

			k.SetVestingQueue(ctx, types.VestingQueue{
				AuctionId:   auction.GetId(),
				Auctioneer:  auction.GetAuctioneer().String(),
				PayingCoin:  sdk.NewCoin(payingCoinDenom, payingAmt),
				ReleaseTime: schedule.ReleaseTime,
				Released:    false,
			})

			remaining = remaining.SubAmount(payingAmt)
		}

		_ = auction.SetStatus(types.AuctionStatusVesting)
		k.SetAuction(ctx, auction)
	}

	return nil
}
