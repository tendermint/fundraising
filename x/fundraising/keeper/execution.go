package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

// ExecuteStandByStatus simply updates the auction status to AuctionStatusStarted
// if the auction is ready to get started.
func (k Keeper) ExecuteStandByStatus(ctx sdk.Context, auction types.AuctionI) {
	if auction.ShouldAuctionStarted(ctx.BlockTime()) { // BlockTime >= StartTime
		_ = auction.SetStatus(types.AuctionStatusStarted)
		k.SetAuction(ctx, auction)
	}
}

// ExecuteStartedStatus executes operations depending on the auction type.
func (k Keeper) ExecuteStartedStatus(ctx sdk.Context, auction types.AuctionI) {
	ctx, writeCache := ctx.CacheContext()

	// Do nothing when the auction still needs time to pass the end time
	if !auction.ShouldAuctionFinished(ctx.BlockTime()) { // BlockTime < EndTime
		return
	}

	// Finish the auction
	switch auction.GetType() {
	case types.AuctionTypeFixedPrice:
		k.FinishFixedPriceAuction(ctx, auction)

	case types.AuctionTypeBatch:
		k.FinishBatchAuction(ctx, auction)
	}

	writeCache()
}

// ExecuteVestingStatus first gets all vesting queues in the store and
// look up the release time of each vesting queue to see if the module needs to
// distribute the paying coin to the auctioneer.
func (k Keeper) ExecuteVestingStatus(ctx sdk.Context, auction types.AuctionI) {
	if err := k.AllocatePayingCoin(ctx, auction); err != nil {
		panic(err)
	}
}
