package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type MatchResult struct {
	MatchPrice          sdk.Dec
	MatchedAmount       sdk.Int
	MatchedBids         []Bid
	MatchResultByBidder map[string]*BidderMatchResult
}

type BidderMatchResult struct {
	PayingAmount  sdk.Int
	MatchedAmount sdk.Int
}

// Match returns the match result for all bids that correspond with the auction.
func Match(matchPrice sdk.Dec, prices []sdk.Dec, bidsByPrice map[string][]Bid, sellingAmt sdk.Int, allowedBidders []AllowedBidder) (res *MatchResult, matched bool) {
	res = &MatchResult{
		MatchPrice:          matchPrice,
		MatchedAmount:       sdk.ZeroInt(),
		MatchResultByBidder: map[string]*BidderMatchResult{},
	}

	biddableAmtByBidder := map[string]sdk.Int{}
	for _, allowedBidder := range allowedBidders {
		biddableAmtByBidder[allowedBidder.Bidder] = allowedBidder.MaxBidAmount
	}

	for _, price := range prices {
		if price.LT(matchPrice) {
			break
		}

		for _, bid := range bidsByPrice[price.String()] {
			var bidAmt sdk.Int
			switch bid.Type {
			case BidTypeBatchWorth:
				bidAmt = sdk.NewDecFromInt(bid.Coin.Amount).QuoTruncate(matchPrice).TruncateInt()
			case BidTypeBatchMany:
				bidAmt = bid.Coin.Amount
			}
			biddableAmt := biddableAmtByBidder[bid.Bidder]
			matchAmt := sdk.MinInt(bidAmt, biddableAmtByBidder[bid.Bidder])

			if res.MatchedAmount.Add(matchAmt).GT(sellingAmt) {
				// Including this bid will exceed the auction's selling amount.
				return nil, false
			}

			payingAmt := matchPrice.MulInt(matchAmt).Ceil().TruncateInt()

			bidderRes, ok := res.MatchResultByBidder[bid.Bidder]
			if !ok {
				bidderRes = &BidderMatchResult{
					PayingAmount:  sdk.ZeroInt(),
					MatchedAmount: sdk.ZeroInt(),
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
