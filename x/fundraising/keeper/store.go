package keeper

import (
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
		err := k.cdc.Unmarshal(bz, &val)
		if err != nil {
			panic(err)
		}
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

// GetAuctions returns all auctions in the store.
func (k Keeper) GetAuctions(ctx sdk.Context) (auctions []types.AuctionI) {
	k.IterateAuctions(ctx, func(auction types.AuctionI) (stop bool) {
		auctions = append(auctions, auction)
		return false
	})

	return auctions
}

// SetAuction sets an auction with the given auction id.
func (k Keeper) SetAuction(ctx sdk.Context, auction types.AuctionI) {
	id := auction.GetId()
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalAuction(k.cdc, auction)
	store.Set(types.GetAuctionKey(id), bz)
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

// GetLastSequence returns the last sequence for the bid.
func (k Keeper) GetLastSequence(ctx sdk.Context, auctionId uint64) uint64 {
	var seq uint64
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetSequenceKey(auctionId))
	if bz == nil {
		seq = 0 // initialize the sequence
	} else {
		val := gogotypes.UInt64Value{}
		err := k.cdc.Unmarshal(bz, &val)
		if err != nil {
			panic(err)
		}
		seq = val.GetValue()
	}
	return seq
}

// SetSequence sets the sequence number for the auction.
func (k Keeper) SetSequence(ctx sdk.Context, auctionId uint64, seq uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&gogotypes.UInt64Value{Value: seq})
	store.Set(types.GetSequenceKey(auctionId), bz)
}

// GetBid returns a bid for the given auction id and sequence number.
// A bidder can have as many bids as they want, so sequence is required to get the bid.
func (k Keeper) GetBid(ctx sdk.Context, auctionId uint64, sequence uint64) (bid types.Bid, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetBidKey(auctionId, sequence))
	if bz == nil {
		return bid, false
	}
	k.cdc.MustUnmarshal(bz, &bid)
	return bid, true
}

// SetBid sets a bid with the given arguments.
func (k Keeper) SetBid(ctx sdk.Context, auctionId uint64, sequence uint64, bidderAddr sdk.AccAddress, bid types.Bid) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&bid)
	store.Set(types.GetBidKey(auctionId, sequence), bz)
	store.Set(types.GetBidIndexKey(bidderAddr, auctionId, sequence), []byte{})
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
	iter := sdk.KVStorePrefixIterator(store, types.GetBidAuctionIDKey(auctionId))
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
		auctionId, sequence := types.ParseBidIndexKey(iter.Key())
		bid, _ := k.GetBid(ctx, auctionId, sequence)
		if cb(bid) {
			break
		}
	}
}
