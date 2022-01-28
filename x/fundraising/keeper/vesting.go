package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

// SetVestingSchedules stores vesting queues based on the vesting schedules of the auction and
// sets status to vesting.
func (k Keeper) SetVestingSchedules(ctx sdk.Context, auction types.AuctionI) error {
	payingReserveAddress := auction.GetPayingReserveAddress()
	vestingReserveAddress := auction.GetVestingReserveAddress()

	reserveCoin := k.bankKeeper.GetBalance(ctx, payingReserveAddress, auction.GetPayingCoinDenom())
	reserveCoins := sdk.NewCoins(reserveCoin)

	lenVestingSchedules := len(auction.GetVestingSchedules())

	if lenVestingSchedules == 0 {
		if err := k.bankKeeper.SendCoins(ctx, payingReserveAddress, auction.GetAuctioneer(), reserveCoins); err != nil {
			return err
		}

		if err := auction.SetStatus(types.AuctionStatusFinished); err != nil {
			return err
		}

		k.SetAuction(ctx, auction)

	} else {
		if err := k.bankKeeper.SendCoins(ctx, payingReserveAddress, vestingReserveAddress, reserveCoins); err != nil {
			return err
		}

		remaining := reserveCoin

		for i, vs := range auction.GetVestingSchedules() {
			payingAmt := reserveCoin.Amount.ToDec().MulTruncate(vs.Weight).TruncateInt()

			// Store remaining to the paying coin in the last queue
			if i == lenVestingSchedules-1 {
				payingAmt = remaining.Amount
			}

			k.SetVestingQueue(ctx, auction.GetId(), vs.ReleaseTime, types.VestingQueue{
				AuctionId:   auction.GetId(),
				Auctioneer:  auction.GetAuctioneer().String(),
				PayingCoin:  sdk.NewCoin(auction.GetPayingCoinDenom(), payingAmt),
				ReleaseTime: vs.ReleaseTime,
				Released:    false,
			})

			remaining = remaining.SubAmount(payingAmt)
		}

		if err := auction.SetStatus(types.AuctionStatusVesting); err != nil {
			return err
		}

		k.SetAuction(ctx, auction)
	}

	return nil
}
