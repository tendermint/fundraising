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
	AllocationMap      map[string]sdk.Int // the allocation map that holds allocate amount information for each bidder
	ReservedMatchedMap map[string]sdk.Int
	RefundMap          map[string]sdk.Int // the refund map that holds refund amount information for each bidder
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
		bidAmt := b.GetBidSellingAmount(auction.GetPayingCoinDenom())

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
	bids := k.GetBidsByAuctionId(ctx, auction.GetId())
	bids = types.SortByBidPrice(bids)

	mInfo := MatchingInfo{
		MatchedLen:         0,
		MatchedPrice:       sdk.ZeroDec(),
		TotalMatchedAmount: sdk.ZeroInt(),
		AllocationMap:      map[string]sdk.Int{},
		ReservedMatchedMap: map[string]sdk.Int{},
		RefundMap:          map[string]sdk.Int{},
	}
	allowedBiddersMap := auction.GetAllowedBiddersMap() // map(bidder => maxBidAmt)
	allocationMap := map[string]sdk.Int{}               // map(bidder => allocatedAmt)
	reservedMatchedMap := map[string]sdk.Int{}          // map(bidder => reservedMatchedAmt)
	reservedMap := map[string]sdk.Int{}                 // map(bidder => reservedAmt)
	refundMap := map[string]sdk.Int{}                   // map(bidder => refundAmt)

	// Initialize value for all maps
	for _, ab := range auction.GetAllowedBidders() {
		mInfo.AllocationMap[ab.Bidder] = sdk.ZeroInt()
		mInfo.ReservedMatchedMap[ab.Bidder] = sdk.ZeroInt()
		refundMap[ab.Bidder] = sdk.ZeroInt()
		reservedMap[ab.Bidder] = sdk.ZeroInt()
	}

	// Iterate from the highest matching bid price and stop until it finds
	// the matching information to store them into MatchingInfo object
	for _, bid := range bids {
		matchingPrice := bid.Price
		totalMatchedAmt := sdk.ZeroInt()

		// Add all allowed bidders for the matching price
		for _, ab := range auction.GetAllowedBidders() {
			allocationMap[ab.Bidder] = sdk.ZeroInt()
			reservedMatchedMap[ab.Bidder] = sdk.ZeroInt()
		}

		// Iterate all bids and execute the logics when the bid price is
		// higher than the current matching price
		for _, b := range bids {
			if b.Price.LT(matchingPrice) {
				continue
			}

			maxBidAmt := allowedBiddersMap[b.Bidder]
			allocateAmt := allocationMap[b.Bidder]

			// Uses minimum of the two amounts to prevent from exceeding the bidder's maximum bid amount
			if b.Type == types.BidTypeBatchWorth {
				bidAmt := b.Coin.Amount.ToDec().QuoTruncate(matchingPrice).TruncateInt()

				// MinInt(BidAmt, MaxBidAmt-AccumulatedBidAmt)
				matchingAmt := sdk.MinInt(bidAmt, maxBidAmt.Sub(allocateAmt))

				// Accumulate matching amount since a bidder can have multiple bids
				if alloc, ok := allocationMap[b.Bidder]; ok {
					allocationMap[b.Bidder] = alloc.Add(matchingAmt)
				}

				// Accumulate how much reserved paying coin amount is matched
				if reservedMatchedAmt, ok := reservedMatchedMap[b.Bidder]; ok {
					reserveAmt := sdk.ZeroInt()
					if matchingAmt.LT(bidAmt) {
						reserveAmt = matchingAmt.ToDec().Mul(matchingPrice).Ceil().TruncateInt()
					} else {
						reserveAmt = b.Coin.Amount
					}
					reservedMatchedMap[b.Bidder] = reservedMatchedAmt.Add(reserveAmt)
				}

				totalMatchedAmt = totalMatchedAmt.Add(matchingAmt)
			} else {
				bidAmt := b.Coin.Amount

				// MinInt(BidAmt, MaxBidAmount-AccumulatedBidAmount)
				matchingAmt := sdk.MinInt(bidAmt, maxBidAmt.Sub(allocateAmt))

				// Accumulate matching amount since a bidder can have multiple bids
				if alloc, ok := allocationMap[b.Bidder]; ok {
					allocationMap[b.Bidder] = alloc.Add(matchingAmt)
				}

				// Accumulate how much reserved paying coin amount is matched
				if reservedMatchedAmt, ok := reservedMatchedMap[b.Bidder]; ok {
					reserveAmt := matchingAmt.ToDec().Mul(matchingPrice).Ceil().TruncateInt()
					reservedMatchedMap[b.Bidder] = reservedMatchedAmt.Add(reserveAmt)
				}

				totalMatchedAmt = totalMatchedAmt.Add(matchingAmt)
			}
		}

		// Exit the iteration when the total matched amount is greater than the total selling coin amount
		if totalMatchedAmt.GT(auction.GetSellingCoin().Amount) {
			break
		}

		mInfo.MatchedLen = mInfo.MatchedLen + 1
		mInfo.MatchedPrice = matchingPrice
		mInfo.TotalMatchedAmount = totalMatchedAmt

		for _, ab := range auction.GetAllowedBidders() {
			mInfo.AllocationMap[ab.Bidder] = allocationMap[ab.Bidder]
			mInfo.ReservedMatchedMap[ab.Bidder] = reservedMatchedMap[ab.Bidder]
		}

		bid.SetMatched(true)
		k.SetBid(ctx, bid)
	}

	// Iterate all bids to get refund amount for each bidder
	// Calculate the refund amount by substracting allocate amount from
	// how much a bidder reserved to place a bid for the auction
	for _, b := range bids {
		if b.Type == types.BidTypeBatchWorth {
			reservedMap[b.Bidder] = reservedMap[b.Bidder].Add(b.Coin.Amount)
		} else {
			bidAmt := b.Coin.Amount.ToDec().Mul(b.Price).Ceil().TruncateInt()
			reservedMap[b.Bidder] = reservedMap[b.Bidder].Add(bidAmt)
		}
	}

	for bidder, reservedAmt := range reservedMap {
		reservedMatchedAmt, ok := mInfo.ReservedMatchedMap[bidder]
		if ok {
			refundMap[bidder] = reservedAmt.Sub(reservedMatchedAmt)
			continue
		}
		refundMap[bidder] = reservedAmt
	}

	mInfo.RefundMap = refundMap

	k.SetMatchedBidsLen(ctx, auction.GetId(), mInfo.MatchedLen)

	return mInfo
}
