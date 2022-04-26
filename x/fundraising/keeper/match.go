package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

// MatchingInfo holds information about a batch auction matching info.
type MatchingInfo struct {
	MatchedLen         int64              // the length of matched bids
	MatchedPrice       sdk.Dec            // the final matched price
	TotalMatchedAmount sdk.Int            // the total sold amount
	AllocationMap      map[string]sdk.Int // the map that holds allocate amount information for each bidder
	ReservedMatchedMap map[string]sdk.Int // the map that holds each bidder's matched amount out of their total reserved amount
	RefundMap          map[string]sdk.Int // the map that holds refund amount information for each bidder
}

func (k Keeper) CalculateFixedPriceAllocation(ctx sdk.Context, auction types.AuctionI) MatchingInfo {
	mInfo := MatchingInfo{
		MatchedPrice:       sdk.ZeroDec(),
		TotalMatchedAmount: sdk.ZeroInt(),
		AllocationMap:      map[string]sdk.Int{},
	}

	totalMatchedAmt := sdk.ZeroInt()
	allocMap := map[string]sdk.Int{}

	for _, b := range k.GetBidsByAuctionId(ctx, auction.GetId()) {
		bidAmt := b.ConvertToSellingAmount(auction.GetPayingCoinDenom())

		// Accumulate bid amount if the bidder has other bid(s)
		if allocatedAmt, ok := allocMap[b.Bidder]; ok {
			allocMap[b.Bidder] = allocatedAmt.Add(bidAmt)
		} else {
			allocMap[b.Bidder] = bidAmt
		}
		totalMatchedAmt = totalMatchedAmt.Add(bidAmt)
		mInfo.MatchedLen = mInfo.MatchedLen + 1
	}

	mInfo.MatchedPrice = auction.GetStartPrice()
	mInfo.TotalMatchedAmount = totalMatchedAmt
	mInfo.AllocationMap = allocMap

	return mInfo
}

func (k Keeper) CalculateBatchAllocation(ctx sdk.Context, auction types.AuctionI) MatchingInfo {
	matchingInfo := MatchingInfo{
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

	matchingInfo.MatchedLen = int64(len(matchRes.MatchedBids))
	matchingInfo.MatchedPrice = matchRes.MatchPrice
	matchingInfo.TotalMatchedAmount = matchRes.MatchedAmount

	reservedAmtByBidder := map[string]sdk.Int{}
	for _, bid := range bids {
		bidderReservedAmt, ok := reservedAmtByBidder[bid.Bidder]
		if !ok {
			bidderReservedAmt = sdk.ZeroInt()
		}
		reservedAmtByBidder[bid.Bidder] = bidderReservedAmt.Add(bid.ConvertToPayingAmount(auction.GetPayingCoinDenom()))
	}

	for bidder, reservedAmt := range reservedAmtByBidder {
		matchingInfo.AllocationMap[bidder] = sdk.ZeroInt()
		matchingInfo.ReservedMatchedMap[bidder] = sdk.ZeroInt()
		matchingInfo.RefundMap[bidder] = reservedAmt
	}

	for bidder, bidderRes := range matchRes.MatchResultByBidder {
		matchingInfo.AllocationMap[bidder] = bidderRes.MatchedAmount
		matchingInfo.ReservedMatchedMap[bidder] = bidderRes.PayingAmount
		matchingInfo.RefundMap[bidder] = reservedAmtByBidder[bidder].Sub(bidderRes.PayingAmount)
	}

	for _, bid := range matchRes.MatchedBids {
		bid.SetMatched(true)
		k.SetBid(ctx, bid)
	}
	k.SetMatchedBidsLen(ctx, auction.GetId(), matchingInfo.MatchedLen)

	return matchingInfo
}
