package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

// ExecuteStandByStatus simply updates the auction status to AuctionStatusStarted
// if the auction is ready to get started.
func (k Keeper) ExecuteStandByStatus(ctx sdk.Context, auction types.AuctionI) {
	if auction.ShouldAuctionStarted(ctx.BlockTime()) {
		_ = auction.SetStatus(types.AuctionStatusStarted)
		k.SetAuction(ctx, auction)
	}
}

// ExecuteStartedStatus executes operations depending on the auction type.
// For FixedPriceAuction, it distributes the allocated paying coin to the bidders  and
// sets vesting schedules if they are defined.
func (k Keeper) ExecuteStartedStatus(ctx sdk.Context, auction types.AuctionI) {
	ctx, writeCache := ctx.CacheContext()

	// Do nothing when the auction is still in started status
	if !auction.ShouldAuctionFinished(ctx.BlockTime()) {
		return
	}

	switch auction.GetType() {
	case types.AuctionTypeFixedPrice:
		if err := k.DistributeSellingCoin(ctx, auction); err != nil {
			panic(err)
		}

		if err := k.SetVestingSchedules(ctx, auction); err != nil {
			panic(err)
		}

	case types.AuctionTypeBatch:
		k.CalculateAllocation(ctx, auction)

		k.FinishBatchAuction(ctx, auction)
	}

	writeCache()
}

// ExecuteVestingStatus first gets all vesting queues in the store and
// look up the release time of each vesting queue to see if the module needs to
// distribute the paying coin to the auctioneer.
func (k Keeper) ExecuteVestingStatus(ctx sdk.Context, auction types.AuctionI) {
	if err := k.DistributePayingCoin(ctx, auction); err != nil {
		panic(err)
	}
}
