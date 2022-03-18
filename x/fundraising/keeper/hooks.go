package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

// Implements FundraisingHooks interface
var _ types.FundraisingHooks = Keeper{}

func (k Keeper) BeforeFixedPriceAuctionCreated(
	ctx sdk.Context,
	auctioneer string,
	startPrice sdk.Dec,
	sellingCoin sdk.Coin,
	payingCoinDenom string,
	vestingSchedules []types.VestingSchedule,
	startTime,
	endTime time.Time,
) {
	if k.hooks != nil {
		k.hooks.BeforeFixedPriceAuctionCreated(
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

func (k Keeper) BeforeBatchAuctionCreated(
	ctx sdk.Context,
	auctioneer string,
	startPrice sdk.Dec,
	minBidPrice sdk.Dec,
	sellingCoin sdk.Coin,
	payingCoinDenom string,
	vestingSchedules []types.VestingSchedule,
	maxExtendedRound uint32,
	extendedRoundRate sdk.Dec,
	startTime time.Time,
	endTime time.Time,
) {
	if k.hooks != nil {
		k.hooks.BeforeBatchAuctionCreated(
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

func (k Keeper) BeforeAuctionCanceled(
	ctx sdk.Context,
	auctionId uint64,
	auctioneer string,
) {
	if k.hooks != nil {
		k.hooks.BeforeAuctionCanceled(ctx, auctionId, auctioneer)
	}
}

func (k Keeper) BeforeBidPlaced(
	ctx sdk.Context,
	auctionId uint64,
	bidder string,
	bidType types.BidType,
	price sdk.Dec,
	coin sdk.Coin,
) {
	if k.hooks != nil {
		k.hooks.BeforeBidPlaced(ctx, auctionId, bidder, bidType, price, coin)
	}
}

func (k Keeper) BeforeBidModified(
	ctx sdk.Context,
	auctionId uint64,
	bidder string,
	bidType types.BidType,
	price sdk.Dec,
	coin sdk.Coin,
) {
	if k.hooks != nil {
		k.hooks.BeforeBidModified(ctx, auctionId, bidder, bidType, price, coin)
	}
}

func (k Keeper) BeforeAllowedBidderAdded(
	ctx sdk.Context,
	auctionId uint64,
	allowedBidder types.AllowedBidder,
) {
	if k.hooks != nil {
		k.hooks.BeforeAllowedBidderAdded(ctx, auctionId, allowedBidder)
	}
}
