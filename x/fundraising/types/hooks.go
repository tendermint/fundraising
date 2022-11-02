package types

// DONTCOVER

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
	startPrice sdk.Dec,
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
			sellingCoin,
			payingCoinDenom,
			vestingSchedules,
			startTime,
			endTime,
		)
	}
}

func (h MultiFundraisingHooks) AfterFixedPriceAuctionCreated(
	ctx sdk.Context,
	auctionId uint64,
	auctioneer string,
	startPrice sdk.Dec,
	sellingCoin sdk.Coin,
	payingCoinDenom string,
	vestingSchedules []VestingSchedule,
	startTime,
	endTime time.Time,
) {
	for i := range h {
		h[i].AfterFixedPriceAuctionCreated(
			ctx,
			auctionId,
			auctioneer,
			startPrice,
			sellingCoin,
			payingCoinDenom,
			vestingSchedules,
			startTime,
			endTime,
		)
	}
}

func (h MultiFundraisingHooks) BeforeBatchAuctionCreated(
	ctx sdk.Context,
	auctioneer string,
	startPrice sdk.Dec,
	minBidPrice sdk.Dec,
	sellingCoin sdk.Coin,
	payingCoinDenom string,
	vestingSchedules []VestingSchedule,
	maxExtendedRound uint32,
	extendedRoundRate sdk.Dec,
	startTime time.Time,
	endTime time.Time,
) {
	for i := range h {
		h[i].BeforeBatchAuctionCreated(
			ctx,
			auctioneer,
			startPrice,
			minBidPrice,
			sellingCoin,
			payingCoinDenom,
			vestingSchedules,
			maxExtendedRound,
			extendedRoundRate,
			startTime,
			endTime,
		)
	}
}

func (h MultiFundraisingHooks) AfterBatchAuctionCreated(
	ctx sdk.Context,
	auctionId uint64,
	auctioneer string,
	startPrice sdk.Dec,
	minBidPrice sdk.Dec,
	sellingCoin sdk.Coin,
	payingCoinDenom string,
	vestingSchedules []VestingSchedule,
	maxExtendedRound uint32,
	extendedRoundRate sdk.Dec,
	startTime time.Time,
	endTime time.Time,
) {
	for i := range h {
		h[i].AfterBatchAuctionCreated(
			ctx,
			auctionId,
			auctioneer,
			startPrice,
			minBidPrice,
			sellingCoin,
			payingCoinDenom,
			vestingSchedules,
			maxExtendedRound,
			extendedRoundRate,
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
	bidId uint64,
	bidder string,
	bidType BidType,
	price sdk.Dec,
	coin sdk.Coin,
) {
	for i := range h {
		h[i].BeforeBidPlaced(ctx, auctionId, bidId, bidder, bidType, price, coin)
	}
}

func (h MultiFundraisingHooks) BeforeBidModified(
	ctx sdk.Context,
	auctionId uint64,
	bidId uint64,
	bidder string,
	bidType BidType,
	price sdk.Dec,
	coin sdk.Coin,
) {
	for i := range h {
		h[i].BeforeBidModified(ctx, auctionId, bidId, bidder, bidType, price, coin)
	}
}

func (h MultiFundraisingHooks) BeforeAllowedBiddersAdded(
	ctx sdk.Context,
	allowedBidders []AllowedBidder,
) {
	for i := range h {
		h[i].BeforeAllowedBiddersAdded(ctx, allowedBidders)
	}
}

func (h MultiFundraisingHooks) BeforeAllowedBidderUpdated(
	ctx sdk.Context,
	auctionId uint64,
	bidder sdk.AccAddress,
	maxBidAmount sdk.Int,
) {
	for i := range h {
		h[i].BeforeAllowedBidderUpdated(ctx, auctionId, bidder, maxBidAmount)
	}
}

func (h MultiFundraisingHooks) BeforeSellingCoinsAllocated(
	ctx sdk.Context,
	auctionId uint64,
	allocationMap map[string]sdk.Int,
	refundMap map[string]sdk.Int,
) {
	for i := range h {
		h[i].BeforeSellingCoinsAllocated(ctx, auctionId, allocationMap, refundMap)
	}
}
