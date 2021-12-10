package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

// SetVestingSchedules stores vesting queues based on the vesting schedules of the auction and
// sets status to vesting.
func (k Keeper) SetVestingSchedules(ctx sdk.Context, auction types.AuctionI) error {
	payingReserveAcc := types.PayingReserveAcc(auction.GetId())
	vestingReserveAcc := types.VestingReserveAcc(auction.GetId())

	reserveBalance := k.bankKeeper.GetBalance(ctx, payingReserveAcc, auction.GetPayingCoinDenom())
	reserveCoins := sdk.NewCoins(reserveBalance)

	if len(auction.GetVestingSchedules()) == 0 {
		auctioneerAcc, err := sdk.AccAddressFromBech32(auction.GetAuctioneer())
		if err != nil {
			return err
		}

		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, payingReserveAcc.String(), auctioneerAcc, reserveCoins)
		if err != nil {
			return err
		}

		if err := auction.SetStatus(types.AuctionStatusFinished); err != nil {
			return err
		}
		return nil
	}

	for _, vs := range auction.GetVestingSchedules() {
		payingAmt := reserveBalance.Amount.ToDec().Mul(vs.Weight).TruncateInt()

		queue := types.NewVestingQueue(
			auction.GetId(),
			auction.GetAuctioneer(),
			sdk.NewCoin(auction.GetPayingCoinDenom(), payingAmt),
			vs.ReleaseTime,
			false,
		)

		// Store vesting queue
		k.SetVestingQueue(ctx, auction.GetId(), vs.ReleaseTime, queue)
	}

	// Move paying coin to vesting reserve module account
	err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, payingReserveAcc.String(), vestingReserveAcc.String(), reserveCoins)
	if err != nil {
		return err
	}

	if err := auction.SetStatus(types.AuctionStatusVesting); err != nil {
		return err
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

// GetVestingQueuesByAuctionId returns all vesting queues associated with the auction id that are registered in the store.
func (k Keeper) GetVestingQueuesByAuctionId(ctx sdk.Context, auctionId uint64) []types.VestingQueue {
	queues := []types.VestingQueue{}
	k.IterateVestingQueuesByAuctionId(ctx, auctionId, func(queue types.VestingQueue) (stop bool) {
		queues = append(queues, queue)
		return false
	})
	return queues
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
