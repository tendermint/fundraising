package keeper

import (
	"time"

	gogotypes "github.com/gogo/protobuf/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

// GetLastAuctionId returns the last auction id.
func (k Keeper) GetLastAuctionId(ctx sdk.Context) uint64 {
	var id uint64
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.LastAuctionIdKey)
	if bz == nil {
		id = 0 // initialize the auction id
	} else {
		val := gogotypes.UInt64Value{}
		k.cdc.MustUnmarshal(bz, &val)
		id = val.GetValue()
	}
	return id
}

// SetAuctionId stores the last auction id.
func (k Keeper) SetAuctionId(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&gogotypes.UInt64Value{Value: id})
	store.Set(types.LastAuctionIdKey, bz)
}

// GetAuction returns an auction interface from the given auction id.
func (k Keeper) GetAuction(ctx sdk.Context, id uint64) (auction types.AuctionI, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetAuctionKey(id))
	if bz == nil {
		return auction, false
	}

	auction = types.MustUnmarshalAuction(k.cdc, bz)

	return auction, true
}

// SetAuction sets an auction with the given auction id.
func (k Keeper) SetAuction(ctx sdk.Context, auction types.AuctionI) {
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalAuction(k.cdc, auction)
	store.Set(types.GetAuctionKey(auction.GetId()), bz)
}

// GetAuctions returns all auctions in the store.
func (k Keeper) GetAuctions(ctx sdk.Context) (auctions []types.AuctionI) {
	k.IterateAuctions(ctx, func(auction types.AuctionI) (stop bool) {
		auctions = append(auctions, auction)
		return false
	})
	return auctions
}

// IterateAuctions iterates over all the stored auctions and performs a callback function.
// Stops iteration when callback returns true.
func (k Keeper) IterateAuctions(ctx sdk.Context, cb func(auction types.AuctionI) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.AuctionKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		auction := types.MustUnmarshalAuction(k.cdc, iterator.Value())

		if cb(auction) {
			break
		}
	}
}

// GetLastBidId returns the last bid id for the bid.
func (k Keeper) GetLastBidId(ctx sdk.Context, auctionId uint64) uint64 {
	var id uint64
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetLastBidIdKey(auctionId))
	if bz == nil {
		id = 0 // initialize the bid id
	} else {
		val := gogotypes.UInt64Value{}
		k.cdc.MustUnmarshal(bz, &val)
		id = val.GetValue()
	}
	return id
}

// SetBidId sets the bid id for the auction.
func (k Keeper) SetBidId(ctx sdk.Context, auctionId uint64, bidId uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&gogotypes.UInt64Value{Value: bidId})
	store.Set(types.GetLastBidIdKey(auctionId), bz)
}

// GetBid returns a bid for the given auction id and bid id.
// A bidder can have as many bids as they want, so bid id is required to get the bid.
func (k Keeper) GetBid(ctx sdk.Context, auctionId uint64, bidId uint64) (bid types.Bid, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetBidKey(auctionId, bidId))
	if bz == nil {
		return bid, false
	}
	k.cdc.MustUnmarshal(bz, &bid)
	return bid, true
}

// SetBid sets a bid with the given arguments.
func (k Keeper) SetBid(ctx sdk.Context, bid types.Bid) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&bid)
	store.Set(types.GetBidKey(bid.AuctionId, bid.Id), bz)
	store.Set(types.GetBidIndexKey(bid.GetBidder(), bid.AuctionId, bid.Id), []byte{})
}

// GetBids returns all bids registered in the store.
func (k Keeper) GetBids(ctx sdk.Context) []types.Bid {
	bids := []types.Bid{}
	k.IterateBids(ctx, func(bid types.Bid) (stop bool) {
		bids = append(bids, bid)
		return false
	})
	return bids
}

// GetBidsByAuctionId returns all bids associated with the auction id that are registered in the store.
func (k Keeper) GetBidsByAuctionId(ctx sdk.Context, auctionId uint64) []types.Bid {
	bids := []types.Bid{}
	k.IterateBidsByAuctionId(ctx, auctionId, func(bid types.Bid) (stop bool) {
		bids = append(bids, bid)
		return false
	})
	return bids
}

// GetBidsByBidder returns all bids associated with the bidder that are registered in the store.
func (k Keeper) GetBidsByBidder(ctx sdk.Context, bidderAddr sdk.AccAddress) []types.Bid {
	bids := []types.Bid{}
	k.IterateBidsByBidder(ctx, bidderAddr, func(bid types.Bid) (stop bool) {
		bids = append(bids, bid)
		return false
	})
	return bids
}

// IterateBids iterates through all bids stored in the store and invokes callback function for each item.
// Stops the iteration when the callback function returns true.
func (k Keeper) IterateBids(ctx sdk.Context, cb func(bid types.Bid) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.BidKeyPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var bid types.Bid
		k.cdc.MustUnmarshal(iter.Value(), &bid)
		if cb(bid) {
			break
		}
	}
}

// IterateBidsByAuctionId iterates through all bids associated with the auction id stored in the store
// and invokes callback function for each item.
// Stops the iteration when the callback function returns true.
func (k Keeper) IterateBidsByAuctionId(ctx sdk.Context, auctionId uint64, cb func(bid types.Bid) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetBidByAuctionIdPrefix(auctionId))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var bid types.Bid
		k.cdc.MustUnmarshal(iter.Value(), &bid)
		if cb(bid) {
			break
		}
	}
}

// IterateBidsByBidder iterates through all bids associated with the bidder stored in the store
// and invokes callback function for each item.
// Stops the iteration when the callback function returns true.
func (k Keeper) IterateBidsByBidder(ctx sdk.Context, bidderAddr sdk.AccAddress, cb func(bid types.Bid) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetBidIndexByBidderPrefix(bidderAddr))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		auctionId, bidId := types.ParseBidIndexKey(iter.Key())
		bid, _ := k.GetBid(ctx, auctionId, bidId)
		if cb(bid) {
			break
		}
	}
}

func (k Keeper) GetMatchedBidsLen(ctx sdk.Context, auctionId uint64) int64 {
	var len int64
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetLastMatchedBidsLenKey(auctionId))
	if bz == nil {
		len = 0 // initialize the auction id
	} else {
		val := gogotypes.Int64Value{}
		k.cdc.MustUnmarshal(bz, &val)
		len = val.GetValue()
	}
	return len
}

func (k Keeper) SetMatchedBidsLen(ctx sdk.Context, auctionId uint64, lastWinningBidsLen int64) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&gogotypes.Int64Value{Value: lastWinningBidsLen})
	store.Set(types.GetLastMatchedBidsLenKey(auctionId), bz)
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
func (k Keeper) SetVestingQueue(ctx sdk.Context, queue types.VestingQueue) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&queue)
	store.Set(types.GetVestingQueueKey(queue.AuctionId, queue.ReleaseTime), bz)
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
