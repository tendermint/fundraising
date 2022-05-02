package keeper

import (
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

// MatchingInfo holds information about an auction matching information.
type MatchingInfo struct {
	MatchedLen         int64              // the length of matched bids
	MatchedPrice       sdk.Dec            // the final matched price
	TotalMatchedAmount sdk.Int            // the total sold amount
	AllocationMap      map[string]sdk.Int // the map that holds allocate amount information for each bidder
	ReservedMatchedMap map[string]sdk.Int // the map that holds each bidder's matched amount out of their total reserved amount
	RefundMap          map[string]sdk.Int // the map that holds refund amount information for each bidder
}

// CalculateFixedPriceAllocation loops through all bids for the auction and calculate matching information.
func (k Keeper) CalculateFixedPriceAllocation(ctx sdk.Context, auction types.AuctionI) MatchingInfo {
	mInfo := MatchingInfo{
		MatchedPrice:       auction.GetStartPrice(),
		TotalMatchedAmount: sdk.ZeroInt(),
		AllocationMap:      map[string]sdk.Int{},
	}

	bids := k.GetBidsByAuctionId(ctx, auction.GetId())

	// All bids for the auction are already matched in message level
	// Loop through all bids and calculate allocated amount
	// Accumulate the allocated amount if a bidder placed multiple bids
	for _, bid := range bids {
		bidAmt := bid.ConvertToSellingAmount(auction.GetPayingCoinDenom())

		allocatedAmt, ok := mInfo.AllocationMap[bid.Bidder]
		if !ok {
			allocatedAmt = sdk.ZeroInt()
		}
		mInfo.AllocationMap[bid.Bidder] = allocatedAmt.Add(bidAmt)
		mInfo.TotalMatchedAmount = mInfo.TotalMatchedAmount.Add(bidAmt)
		mInfo.MatchedLen++
	}

	return mInfo
}

func (k Keeper) CalculateBatchAllocation(ctx sdk.Context, auction types.AuctionI) MatchingInfo {
	mInfo := MatchingInfo{
		AllocationMap:      map[string]sdk.Int{},
		ReservedMatchedMap: map[string]sdk.Int{},
		RefundMap:          map[string]sdk.Int{},
	}

	bids := k.GetBidsByAuctionId(ctx, auction.GetId())
	prices, bidsByPrice := types.BidsByPrice(bids)
	sellingAmt := auction.GetSellingCoin().Amount

	allowedBidders := k.GetAllowedBiddersByAuction(ctx, auction.GetId())

	matchRes := &types.MatchResult{
		MatchPrice:          sdk.Dec{},
		MatchedAmount:       sdk.ZeroInt(),
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

	reservedAmtByBidder := map[string]sdk.Int{}
	for _, bid := range bids {
		bidderReservedAmt, ok := reservedAmtByBidder[bid.Bidder]
		if !ok {
			bidderReservedAmt = sdk.ZeroInt()
		}
		reservedAmtByBidder[bid.Bidder] = bidderReservedAmt.Add(bid.ConvertToPayingAmount(auction.GetPayingCoinDenom()))
	}

	for bidder, reservedAmt := range reservedAmtByBidder {
		mInfo.AllocationMap[bidder] = sdk.ZeroInt()
		mInfo.ReservedMatchedMap[bidder] = sdk.ZeroInt()
		mInfo.RefundMap[bidder] = reservedAmt
	}

	for bidder, bidderRes := range matchRes.MatchResultByBidder {
		mInfo.AllocationMap[bidder] = bidderRes.MatchedAmount
		mInfo.ReservedMatchedMap[bidder] = bidderRes.PayingAmount
		mInfo.RefundMap[bidder] = reservedAmtByBidder[bidder].Sub(bidderRes.PayingAmount)
	}

	for _, bid := range matchRes.MatchedBids {
		bid.SetMatched(true)
		k.SetBid(ctx, bid)
	}
	k.SetMatchedBidsLen(ctx, auction.GetId(), mInfo.MatchedLen)

	return mInfo
}
