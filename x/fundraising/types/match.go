package types

import (
	"cosmossdk.io/math"
)

type MatchResult struct {
	MatchPrice          math.LegacyDec
	MatchedAmount       math.Int
	MatchedBids         []Bid
	MatchResultByBidder map[string]*BidderMatchResult
}

type BidderMatchResult struct {
	PayingAmount  math.Int
	MatchedAmount math.Int
}

// Match returns the match result for all bids that correspond with the auction.
func Match(matchPrice math.LegacyDec, prices []math.LegacyDec, bidsByPrice map[string][]Bid, sellingAmt math.Int, allowedBidders []AllowedBidder) (res *MatchResult, matched bool) {
	res = &MatchResult{
		MatchPrice:          matchPrice,
		MatchedAmount:       math.ZeroInt(),
		MatchResultByBidder: map[string]*BidderMatchResult{},
	}

	biddableAmtByBidder := map[string]math.Int{}
	for _, allowedBidder := range allowedBidders {
		biddableAmtByBidder[allowedBidder.Bidder] = allowedBidder.MaxBidAmount
	}

	for _, price := range prices {
		if price.LT(matchPrice) {
			break
		}

		for _, bid := range bidsByPrice[price.String()] {
			var bidAmt math.Int
			switch bid.Type {
			case BidTypeBatchWorth:
				bidAmt = math.LegacyNewDecFromInt(bid.Coin.Amount).QuoTruncate(matchPrice).TruncateInt()
			case BidTypeBatchMany:
				bidAmt = bid.Coin.Amount
			}
			biddableAmt := biddableAmtByBidder[bid.Bidder]
			matchAmt := math.MinInt(bidAmt, biddableAmtByBidder[bid.Bidder])

			if res.MatchedAmount.Add(matchAmt).GT(sellingAmt) {
				// Including this bid will exceed the auction's selling amount.
				return nil, false
			}

			payingAmt := matchPrice.MulInt(matchAmt).Ceil().TruncateInt()

			bidderRes, ok := res.MatchResultByBidder[bid.Bidder]
			if !ok {
				bidderRes = &BidderMatchResult{
					PayingAmount:  math.ZeroInt(),
					MatchedAmount: math.ZeroInt(),
				}
				res.MatchResultByBidder[bid.Bidder] = bidderRes
			}
			bidderRes.MatchedAmount = bidderRes.MatchedAmount.Add(matchAmt)
			bidderRes.PayingAmount = bidderRes.PayingAmount.Add(payingAmt)

			if matchAmt.IsPositive() {
				biddableAmtByBidder[bid.Bidder] = biddableAmt.Sub(matchAmt)
				res.MatchedBids = append(res.MatchedBids, bid)
				res.MatchedAmount = res.MatchedAmount.Add(matchAmt)
				matched = true
			}
		}
	}

	return res, matched
}
