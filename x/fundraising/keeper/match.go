package keeper

import (
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

	var matchRes *types.MatchResult
	for i, price := range prices {
		res, found := types.Match(auction, price, prices, bidsByPrice)
		if found || (matchRes == nil && i == len(prices)-1) {
			matchRes = res
		}
		if !found {
			break
		}
	}

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
