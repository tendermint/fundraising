package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

// SetVestingSchedules stores vesting queues based on the vesting schedules of the auction and
// sets status to vesting.
func (k Keeper) SetVestingSchedules(ctx sdk.Context, auction types.AuctionI) error {
	payingReserveAddr := auction.GetPayingReserveAddress()
	vestingReserveAddr := auction.GetVestingReserveAddress()
	reserveCoin := k.bankKeeper.GetBalance(ctx, payingReserveAddr, auction.GetPayingCoinDenom())
	reserveCoins := sdk.NewCoins(reserveCoin)

	vsLen := len(auction.GetVestingSchedules())
	if vsLen == 0 {
		if err := k.bankKeeper.SendCoins(ctx, payingReserveAddr, auction.GetAuctioneer(), reserveCoins); err != nil {
			return err
		}

		if err := auction.SetStatus(types.AuctionStatusFinished); err != nil {
			return err
		}
		k.SetAuction(ctx, auction)

	} else {
		if err := k.bankKeeper.SendCoins(ctx, payingReserveAddr, vestingReserveAddr, reserveCoins); err != nil {
			return err
		}

		remaining := reserveCoin

		for i, vs := range auction.GetVestingSchedules() {
			payingAmt := reserveCoin.Amount.ToDec().MulTruncate(vs.Weight).TruncateInt()

			// Store remaining to the paying coin in the last queue
			if i == vsLen-1 {
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
