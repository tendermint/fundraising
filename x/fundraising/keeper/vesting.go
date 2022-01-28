package keeper

import (
	"time"

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

// GetVestingQueue returns a slice of vesting queues that the auction is complete and
// waiting in a queue to release the vesting amount of coin at the respective release time.
func (k Keeper) GetVestingQueue(ctx sdk.Context, auctionId uint64, releaseTime time.Time) types.VestingQueue {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetVestingQueueKey(auctionId, releaseTime))
	if bz == nil {
		return types.VestingQueue{}
	}

	queue := types.VestingQueue{}
	k.cdc.MustUnmarshal(bz, &queue)

	return queue
}

// SetVestingQueue sets vesting queue into with the given release time and auction id.
func (k Keeper) SetVestingQueue(ctx sdk.Context, auctionId uint64, releaseTime time.Time, queue types.VestingQueue) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&queue)
	store.Set(types.GetVestingQueueKey(auctionId, releaseTime), bz)
}

// GetVestingQueues returns all vesting queues registered in the store.
func (k Keeper) GetVestingQueues(ctx sdk.Context) []types.VestingQueue {
	queues := []types.VestingQueue{}
	k.IterateVestingQueues(ctx, func(queue types.VestingQueue) (stop bool) {
		queues = append(queues, queue)
		return false
	})
	return queues
}

// GetVestingQueuesByAuctionId returns all vesting queues associated with the auction id that are registered in the store.
func (k Keeper) GetVestingQueuesByAuctionId(ctx sdk.Context, auctionId uint64) []types.VestingQueue {
	queues := []types.VestingQueue{}
	k.IterateVestingQueuesByAuctionId(ctx, auctionId, func(queue types.VestingQueue) (stop bool) {
		queues = append(queues, queue)
		return false
	})
	return queues
}

// IterateVestingQueues iterates through all VestingQueues and invokes callback function for each item.
// Stops the iteration when the callback function returns true.
func (k Keeper) IterateVestingQueues(ctx sdk.Context, cb func(queue types.VestingQueue) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.VestingQueueKeyPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var queue types.VestingQueue
		k.cdc.MustUnmarshal(iter.Value(), &queue)
		if cb(queue) {
			break
		}
	}
}

// IterateVestingQueuesByAuctionId iterates through all VestingQueues associated with the auction id stored in the store
// and invokes callback function for each item.
// Stops the iteration when the callback function returns true.
func (k Keeper) IterateVestingQueuesByAuctionId(ctx sdk.Context, auctionId uint64, cb func(queue types.VestingQueue) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetVestingQueueByAuctionIdPrefix(auctionId))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var queue types.VestingQueue
		k.cdc.MustUnmarshal(iter.Value(), &queue)
		if cb(queue) {
			break
		}
	}
}
