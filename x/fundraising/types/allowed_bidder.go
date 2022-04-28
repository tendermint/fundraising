package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// NewAllowedBidder returns a new AllowedBidder.
func NewAllowedBidder(auctionId uint64, bidderAddr sdk.AccAddress, maxBidAmount sdk.Int) AllowedBidder {
	return AllowedBidder{
		AuctionId:    auctionId,
		Bidder:       bidderAddr.String(),
		MaxBidAmount: maxBidAmount,
	}
}

// GetBidder returns the bidder account address.
func (ab AllowedBidder) GetBidder() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(ab.Bidder)
	if err != nil {
		panic(err)
	}
	return addr
}

// GetAllowedBiddersMap returns allowed bidders map.
func GetAllowedBiddersMap(allowedBidders []AllowedBidder) map[string]sdk.Int { // map(bidder => maxBidAmount)
	allowedBiddersMap := make(map[string]sdk.Int)
	for _, allowedBidder := range allowedBidders {
		allowedBiddersMap[allowedBidder.Bidder] = allowedBidder.MaxBidAmount
	}
	return allowedBiddersMap
}
