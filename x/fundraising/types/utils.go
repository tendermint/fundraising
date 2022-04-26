package types

import (
	"sort"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	"github.com/tendermint/tendermint/crypto"
)

// MustParseRFC3339 parses string time to time in RFC3339 format.
// This is used only for internal testing purpose.
func MustParseRFC3339(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

// DeriveAddress derives an address with the given address length type, module name, and
// address derivation name. It is used to derive reserve account addresses for selling, paying, and vesting.
func DeriveAddress(addressType AddressType, moduleName, name string) sdk.AccAddress {
	switch addressType {
	case AddressType32Bytes:
		return sdk.AccAddress(address.Module(moduleName, []byte(name)))
	case AddressType20Bytes:
		return sdk.AccAddress(crypto.AddressHash([]byte(moduleName + name)))
	default:
		return sdk.AccAddress{}
	}
}

// SortBids sorts bid array by bid price in descending order.
func SortBids(bids []Bid) []Bid {
	sort.Slice(bids, func(i, j int) bool {
		if bids[i].Price.GT(bids[j].Price) {
			return true
		}
		return bids[i].Id < bids[j].Id
	})
	return bids
}

func BidsByPrice(bids []Bid) (prices []sdk.Dec, bidsByPrice map[string][]Bid) {
	bids = SortBids(bids)

	bidsByPrice = map[string][]Bid{} // price => []Bid

	for _, bid := range bids {
		priceStr := bid.Price.String()
		bidsByPrice[priceStr] = append(bidsByPrice[priceStr], bid)
	}

	// Sort prices in descending order.
	prices = make([]sdk.Dec, len(bidsByPrice))
	i := 0 // TODO: is it too much optimization? we can use append(...)
	for priceStr := range bidsByPrice {
		prices[i] = sdk.MustNewDecFromStr(priceStr)
		i++
	}
	sort.Slice(prices, func(i, j int) bool {
		return prices[i].GT(prices[j])
	})
	return
}

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

func Match(auction AuctionI, matchPrice sdk.Dec, prices []sdk.Dec, bidsByPrice map[string][]Bid) (res *MatchResult, matched bool) {
	biddableAmtByBidder := auction.GetAllowedBiddersMap()
	res = &MatchResult{
		MatchPrice:          matchPrice,
		MatchedAmount:       sdk.ZeroInt(),
		MatchResultByBidder: map[string]*BidderMatchResult{},
	}

	for _, price := range prices {
		if price.LT(matchPrice) {
			break
		}

		for _, bid := range bidsByPrice[price.String()] {
			var bidAmt sdk.Int
			switch bid.Type {
			case BidTypeBatchWorth:
				bidAmt = bid.Coin.Amount.ToDec().QuoTruncate(matchPrice).TruncateInt()
			case BidTypeBatchMany:
				bidAmt = bid.Coin.Amount
			}
			biddableAmt := biddableAmtByBidder[bid.Bidder]
			matchAmt := sdk.MinInt(bidAmt, biddableAmtByBidder[bid.Bidder])

			if res.MatchedAmount.Add(matchAmt).GT(auction.GetSellingCoin().Amount) {
				// Including this bid will exceed the auction's selling amount.
				// Thus, we found the ideal match price.
				return res, false
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

			biddableAmtByBidder[bid.Bidder] = biddableAmt.Sub(matchAmt)
			res.MatchedBids = append(res.MatchedBids, bid)
			res.MatchedAmount = res.MatchedAmount.Add(matchAmt)
		}
	}

	return res, true
}
