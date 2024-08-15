package keeper

import (
	"context"
	"errors"
	"sort"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

// MatchingInfo holds information about an auction matching information.
type MatchingInfo struct {
	MatchedLen         int64               // the length of matched bids
	MatchedPrice       math.LegacyDec      // the final matched price
	TotalMatchedAmount math.Int            // the total sold amount
	AllocationMap      map[string]math.Int // the map that holds allocate amount information for each bidder
	ReservedMatchedMap map[string]math.Int // the map that holds each bidder's matched amount out of their total reserved amount
	RefundMap          map[string]math.Int // the map that holds refund amount information for each bidder
}

func (k Keeper) GetLastMatchedBidsLen(ctx context.Context, auctionId uint64) (int64, error) {
	matchedBidsLen, err := k.MatchedBidsLen.Get(ctx, auctionId)
	if errors.Is(err, collections.ErrNotFound) {
		return 0, nil
	}
	return matchedBidsLen, err
}

func (k Keeper) SetMatchedBidsLen(ctx context.Context, auctionId uint64, matchedLen int64) error {
	return k.MatchedBidsLen.Set(ctx, auctionId, matchedLen)
}

// CalculateFixedPriceAllocation loops through all bids for the auction and calculate matching information.
func (k Keeper) CalculateFixedPriceAllocation(ctx context.Context, auction types.AuctionI) (MatchingInfo, error) {
	mInfo := MatchingInfo{
		MatchedPrice:       auction.GetStartPrice(),
		TotalMatchedAmount: math.ZeroInt(),
		AllocationMap:      map[string]math.Int{},
	}

	bids, err := k.GetBidsByAuctionId(ctx, auction.GetId())
	if err != nil {
		return mInfo, err
	}

	// All bids for the auction are already matched in message level
	// Loop through all bids and calculate allocated amount
	// Accumulate the allocated amount if a bidder placed multiple bids
	for _, bid := range bids {
		bidAmt := bid.ConvertToSellingAmount(auction.GetPayingCoinDenom())

		allocatedAmt, ok := mInfo.AllocationMap[bid.Bidder]
		if !ok {
			allocatedAmt = math.ZeroInt()
		}
		mInfo.AllocationMap[bid.Bidder] = allocatedAmt.Add(bidAmt)
		mInfo.TotalMatchedAmount = mInfo.TotalMatchedAmount.Add(bidAmt)
		mInfo.MatchedLen++
	}

	return mInfo, nil
}

func (k Keeper) CalculateBatchAllocation(ctx context.Context, auction types.AuctionI) (MatchingInfo, error) {
	mInfo := MatchingInfo{
		AllocationMap:      map[string]math.Int{},
		ReservedMatchedMap: map[string]math.Int{},
		RefundMap:          map[string]math.Int{},
	}

	bids, err := k.GetBidsByAuctionId(ctx, auction.GetId())
	if err != nil {
		return mInfo, err
	}
	prices, bidsByPrice := types.BidsByPrice(bids)
	sellingAmt := auction.GetSellingCoin().Amount

	allowedBidders, err := k.GetAllowedBiddersByAuction(ctx, auction.GetId())
	if err != nil {
		return mInfo, err
	}

	matchRes := &types.MatchResult{
		MatchPrice:          math.LegacyDec{},
		MatchedAmount:       math.ZeroInt(),
		MatchResultByBidder: map[string]*types.BidderMatchResult{},
	}

	// We use binary search to find the best(the lowest possible) matching price.
	// Note that the returned index from sort.Search is not used, since
	// we're already storing the match result inside the closure.
	// In this way, we can reduce redundant calculation for the matching price
	// after finding it.
	sort.Search(len(prices), func(i int) bool {
		// Reverse the index, since prices are sorted in descending order.
		// Note that our goal is to find the first true(matched) condition, starting
		// from the lowest price.
		i = (len(prices) - 1) - i
		res, matched := types.Match(prices[i], prices, bidsByPrice, sellingAmt, allowedBidders)
		if matched { // If we found a valid matching price, store the result
			matchRes = res
		}
		return matched
	})

	mInfo.MatchedLen = int64(len(matchRes.MatchedBids))
	mInfo.MatchedPrice = matchRes.MatchPrice
	mInfo.TotalMatchedAmount = matchRes.MatchedAmount

	reservedAmtByBidder := map[string]math.Int{}
	for _, bid := range bids {
		bidderReservedAmt, ok := reservedAmtByBidder[bid.Bidder]
		if !ok {
			bidderReservedAmt = math.ZeroInt()
		}
		reservedAmtByBidder[bid.Bidder] = bidderReservedAmt.Add(bid.ConvertToPayingAmount(auction.GetPayingCoinDenom()))
	}

	for bidder, reservedAmt := range reservedAmtByBidder {
		mInfo.AllocationMap[bidder] = math.ZeroInt()
		mInfo.ReservedMatchedMap[bidder] = math.ZeroInt()
		mInfo.RefundMap[bidder] = reservedAmt
	}

	for bidder, bidderRes := range matchRes.MatchResultByBidder {
		mInfo.AllocationMap[bidder] = bidderRes.MatchedAmount
		mInfo.ReservedMatchedMap[bidder] = bidderRes.PayingAmount
		mInfo.RefundMap[bidder] = reservedAmtByBidder[bidder].Sub(bidderRes.PayingAmount)
	}

	for _, bid := range matchRes.MatchedBids {
		bid.SetMatched(true)
		if err := k.Bid.Set(ctx, collections.Join(bid.AuctionId, bid.Id), bid); err != nil {
			return mInfo, err
		}
	}

	return mInfo, k.SetMatchedBidsLen(ctx, auction.GetId(), mInfo.MatchedLen)
}
