package keeper

import (
	"context"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

// Implements FundraisingHooks interface
var _ types.FundraisingHooks = Keeper{}

// SetHooks sets the fundraising hooks.
func (k *Keeper) SetHooks(fk types.FundraisingHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set fundraising hooks twice")
	}
	k.hooks = fk
	return k
}

// BeforeFixedPriceAuctionCreated - call hook if registered
func (k Keeper) BeforeFixedPriceAuctionCreated(
	ctx context.Context,
	auctioneer string,
	startPrice math.LegacyDec,
	sellingCoin sdk.Coin,
	payingCoinDenom string,
	vestingSchedules []types.VestingSchedule,
	startTime time.Time,
	endTime time.Time,
) error {
	if k.hooks != nil {
		if err := k.hooks.BeforeFixedPriceAuctionCreated(
			ctx,
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

// AfterFixedPriceAuctionCreated - call hook if registered
func (k Keeper) AfterFixedPriceAuctionCreated(
	ctx context.Context,
	auctionId uint64,
	auctioneer string,
	startPrice math.LegacyDec,
	sellingCoin sdk.Coin,
	payingCoinDenom string,
	vestingSchedules []types.VestingSchedule,
	startTime time.Time,
	endTime time.Time,
) error {
	if k.hooks != nil {
		if err := k.hooks.AfterFixedPriceAuctionCreated(
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

// BeforeBatchAuctionCreated - call hook if registered
func (k Keeper) BeforeBatchAuctionCreated(
	ctx context.Context,
	auctioneer string,
	startPrice math.LegacyDec,
	minBidPrice math.LegacyDec,
	sellingCoin sdk.Coin,
	payingCoinDenom string,
	vestingSchedules []types.VestingSchedule,
	maxExtendedRound uint32,
	extendedRoundRate math.LegacyDec,
	startTime time.Time,
	endTime time.Time,
) error {
	if k.hooks != nil {
		if err := k.hooks.BeforeBatchAuctionCreated(
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

// AfterBatchAuctionCreated - call hook if registered
func (k Keeper) AfterBatchAuctionCreated(
	ctx context.Context,
	auctionId uint64,
	auctioneer string,
	startPrice math.LegacyDec,
	minBidPrice math.LegacyDec,
	sellingCoin sdk.Coin,
	payingCoinDenom string,
	vestingSchedules []types.VestingSchedule,
	maxExtendedRound uint32,
	extendedRoundRate math.LegacyDec,
	startTime time.Time,
	endTime time.Time,
) error {
	if k.hooks != nil {
		if err := k.hooks.AfterBatchAuctionCreated(
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

// BeforeAuctionCanceled - call hook if registered
func (k Keeper) BeforeAuctionCanceled(
	ctx context.Context,
	auctionId uint64,
	auctioneer string,
) error {
	if k.hooks != nil {
		if err := k.hooks.BeforeAuctionCanceled(ctx, auctionId, auctioneer); err != nil {
			return err
		}
	}
	return nil
}

// BeforeBidPlaced - call hook if registered
func (k Keeper) BeforeBidPlaced(
	ctx context.Context,
	auctionId uint64,
	bidId uint64,
	bidder string,
	bidType types.BidType,
	price math.LegacyDec,
	coin sdk.Coin,
) error {
	if k.hooks != nil {
		if err := k.hooks.BeforeBidPlaced(ctx, auctionId, bidId, bidder, bidType, price, coin); err != nil {
			return err
		}
	}
	return nil
}

// BeforeBidModified - call hook if registered
func (k Keeper) BeforeBidModified(
	ctx context.Context,
	auctionId uint64,
	bidId uint64,
	bidder string,
	bidType types.BidType,
	price math.LegacyDec,
	coin sdk.Coin,
) error {
	if k.hooks != nil {
		if err := k.hooks.BeforeBidModified(ctx, auctionId, bidId, bidder, bidType, price, coin); err != nil {
			return err
		}
	}
	return nil
}

// BeforeAllowedBiddersAdded - call hook if registered
func (k Keeper) BeforeAllowedBiddersAdded(
	ctx context.Context,
	allowedBidders []types.AllowedBidder,
) error {
	if k.hooks != nil {
		if err := k.hooks.BeforeAllowedBiddersAdded(ctx, allowedBidders); err != nil {
			return err
		}
	}
	return nil
}

// BeforeAllowedBidderUpdated - call hook if registered
func (k Keeper) BeforeAllowedBidderUpdated(
	ctx context.Context,
	auctionId uint64,
	bidder sdk.AccAddress,
	maxBidAmount math.Int,
) error {
	if k.hooks != nil {
		if err := k.hooks.BeforeAllowedBidderUpdated(ctx, auctionId, bidder, maxBidAmount); err != nil {
			return err
		}
	}
	return nil
}

// BeforeSellingCoinsAllocated - call hook if registered
func (k Keeper) BeforeSellingCoinsAllocated(
	ctx context.Context,
	auctionId uint64,
	allocationMap map[string]math.Int,
	refundMap map[string]math.Int,
) error {
	if k.hooks != nil {
		if err := k.hooks.BeforeSellingCoinsAllocated(ctx, auctionId, allocationMap, refundMap); err != nil {
			return err
		}
	}
	return nil
}
