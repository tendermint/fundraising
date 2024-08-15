package types

// DONTCOVER

import (
	"context"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ FundraisingHooks = MultiFundraisingHooks{}

// MultiFundraisingHooks combines multiple fundraising hooks.
// All hook functions are run in array sequence
type MultiFundraisingHooks []FundraisingHooks

func NewMultiFundraisingHooks(hooks ...FundraisingHooks) MultiFundraisingHooks {
	return hooks
}

func (h MultiFundraisingHooks) BeforeFixedPriceAuctionCreated(
	ctx context.Context,
	auctioneer string,
	startPrice math.LegacyDec,
	sellingCoin sdk.Coin,
	payingCoinDenom string,
	vestingSchedules []VestingSchedule,
	startTime,
	endTime time.Time,
) error {
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
	return nil
}

func (h MultiFundraisingHooks) AfterFixedPriceAuctionCreated(
	ctx context.Context,
	auctionId uint64,
	auctioneer string,
	startPrice math.LegacyDec,
	sellingCoin sdk.Coin,
	payingCoinDenom string,
	vestingSchedules []VestingSchedule,
	startTime,
	endTime time.Time,
) error {
	for i := range h {
		if err := h[i].AfterFixedPriceAuctionCreated(
			ctx,
			auctionId,
			auctioneer,
			startPrice,
			sellingCoin,
			payingCoinDenom,
			vestingSchedules,
			startTime,
			endTime,
		); err != nil {
			return err
		}
	}
	return nil
}

func (h MultiFundraisingHooks) BeforeBatchAuctionCreated(
	ctx context.Context,
	auctioneer string,
	startPrice math.LegacyDec,
	minBidPrice math.LegacyDec,
	sellingCoin sdk.Coin,
	payingCoinDenom string,
	vestingSchedules []VestingSchedule,
	maxExtendedRound uint32,
	extendedRoundRate math.LegacyDec,
	startTime time.Time,
	endTime time.Time,
) error {
	for i := range h {
		if err := h[i].BeforeBatchAuctionCreated(
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
		); err != nil {
			return err
		}
	}
	return nil
}

func (h MultiFundraisingHooks) AfterBatchAuctionCreated(
	ctx context.Context,
	auctionId uint64,
	auctioneer string,
	startPrice math.LegacyDec,
	minBidPrice math.LegacyDec,
	sellingCoin sdk.Coin,
	payingCoinDenom string,
	vestingSchedules []VestingSchedule,
	maxExtendedRound uint32,
	extendedRoundRate math.LegacyDec,
	startTime time.Time,
	endTime time.Time,
) error {
	for i := range h {
		if err := h[i].AfterBatchAuctionCreated(
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
		); err != nil {
			return err
		}
	}
	return nil
}

func (h MultiFundraisingHooks) BeforeAuctionCanceled(
	ctx context.Context,
	auctionId uint64,
	auctioneer string,
) error {
	for i := range h {
		if err := h[i].BeforeAuctionCanceled(ctx, auctionId, auctioneer); err != nil {
			return err
		}
	}
	return nil
}

func (h MultiFundraisingHooks) BeforeBidPlaced(
	ctx context.Context,
	auctionId uint64,
	bidId uint64,
	bidder string,
	bidType BidType,
	price math.LegacyDec,
	coin sdk.Coin,
) error {
	for i := range h {
		if err := h[i].BeforeBidPlaced(ctx, auctionId, bidId, bidder, bidType, price, coin); err != nil {
			return err
		}
	}
	return nil
}

func (h MultiFundraisingHooks) BeforeBidModified(
	ctx context.Context,
	auctionId uint64,
	bidId uint64,
	bidder string,
	bidType BidType,
	price math.LegacyDec,
	coin sdk.Coin,
) error {
	for i := range h {
		if err := h[i].BeforeBidModified(ctx, auctionId, bidId, bidder, bidType, price, coin); err != nil {
			return err
		}
	}
	return nil
}

func (h MultiFundraisingHooks) BeforeAllowedBiddersAdded(
	ctx context.Context,
	allowedBidders []AllowedBidder,
) error {
	for i := range h {
		if err := h[i].BeforeAllowedBiddersAdded(ctx, allowedBidders); err != nil {
			return err
		}
	}
	return nil
}

func (h MultiFundraisingHooks) BeforeAllowedBidderUpdated(
	ctx context.Context,
	auctionId uint64,
	bidder sdk.AccAddress,
	maxBidAmount math.Int,
) error {
	for i := range h {
		if err := h[i].BeforeAllowedBidderUpdated(ctx, auctionId, bidder, maxBidAmount); err != nil {
			return err
		}
	}
	return nil
}

func (h MultiFundraisingHooks) BeforeSellingCoinsAllocated(
	ctx context.Context,
	auctionId uint64,
	allocationMap map[string]math.Int,
	refundMap map[string]math.Int,
) error {
	for i := range h {
		if err := h[i].BeforeSellingCoinsAllocated(ctx, auctionId, allocationMap, refundMap); err != nil {
			return err
		}
	}
	return nil
}
