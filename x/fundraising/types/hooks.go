package types

import (
	time "time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MultiFundraisingHooks combines multiple fundraising hooks.
// All hook functions are run in array sequence
type MultiFundraisingHooks []FundraisingHooks

func NewMultiFundraisingHooks(hooks ...FundraisingHooks) MultiFundraisingHooks {
	return hooks
}

func (h MultiFundraisingHooks) BeforeFixedPriceAuctionCreated(
	ctx sdk.Context,
	auctioneer string,
	startPrice,
	minBidPrice sdk.Dec,
	sellingCoin sdk.Coin,
	payingCoinDenom string,
	vestingSchedules []VestingSchedule,
	startTime,
	endTime time.Time,
) {
	for i := range h {
		h[i].BeforeFixedPriceAuctionCreated(
			ctx,
			auctioneer,
			startPrice,
			minBidPrice,
			sellingCoin,
			payingCoinDenom,
			vestingSchedules,
			startTime,
			endTime,
		)
	}
}

func (h MultiFundraisingHooks) BeforeAuctionCanceled(
	ctx sdk.Context,
	auctionId uint64,
	auctioneer string,
) {
	for i := range h {
		h[i].BeforeAuctionCanceled(ctx, auctionId, auctioneer)
	}
}

func (h MultiFundraisingHooks) BeforeBidPlaced(
	ctx sdk.Context,
	auctionId uint64,
	bidder string,
	bidType BidType,
	price sdk.Dec,
	coin sdk.Coin,
) {
	for i := range h {
		h[i].BeforeBidPlaced(ctx, auctionId, bidder, bidType, price, coin)
	}
}

func (h MultiFundraisingHooks) BeforeBidModified(
	ctx sdk.Context,
	auctionId uint64,
	bidder string,
	bidType BidType,
	price sdk.Dec,
	coin sdk.Coin,
) {
	for i := range h {
		h[i].BeforeBidModified(ctx, auctionId, bidder, bidType, price, coin)
	}
}

func (h MultiFundraisingHooks) BeforeAllowedBidderAdded(
	ctx sdk.Context,
	auctionId uint64,
	allowedBidder AllowedBidder,
) {
	for i := range h {
		h[i].BeforeAllowedBidderAdded(ctx, auctionId, allowedBidder)
	}
}
