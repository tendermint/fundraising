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
	reserveCoin := k.bankKeeper.GetBalance(ctx, payingReserveAddr, auction.GetPayingCoinDenom())
	reserveCoins := sdk.NewCoins(reserveCoin)

	vsLen := len(auction.GetVestingSchedules())
	if vsLen == 0 {
		// Send reserve coins to the auctioneer from the paying reserve account
		if err := k.bankKeeper.SendCoins(ctx, payingReserveAddr, auction.GetAuctioneer(), reserveCoins); err != nil {
			return err
		}

		_ = auction.SetStatus(types.AuctionStatusFinished)
		k.SetAuction(ctx, auction)

	} else {
		// Send reserve coins to the vesting reserve account from the paying reserve account
		if err := k.bankKeeper.SendCoins(ctx, payingReserveAddr, vestingReserveAddr, reserveCoins); err != nil {
			return err
		}

		remaining := reserveCoin

		for i, schedule := range auction.GetVestingSchedules() {
			payingAmt := reserveCoin.Amount.ToDec().MulTruncate(schedule.Weight).TruncateInt()

			// All the remaining paying coin goes to the last vesting queue
			if i == vsLen-1 {
				payingAmt = remaining.Amount
			}

			k.SetVestingQueue(ctx, types.VestingQueue{
				AuctionId:   auction.GetId(),
				Auctioneer:  auction.GetAuctioneer().String(),
				PayingCoin:  sdk.NewCoin(auction.GetPayingCoinDenom(), payingAmt),
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
