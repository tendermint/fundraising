package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

// ExecuteStandByStatus simply updates the auction status to AuctionStatusStarted
// if the auction is ready to get started.
func (k Keeper) ExecuteStandByStatus(ctx context.Context, auction types.AuctionI) error {
	blockTime := sdk.UnwrapSDKContext(ctx).BlockTime()
	if auction.ShouldAuctionStarted(blockTime) { // BlockTime >= StartTime
		if err := auction.SetStatus(types.AuctionStatusStarted); err != nil {
			return err
		}
		if err := k.Auction.Set(ctx, auction.GetId(), auction); err != nil {
			return err
		}
	}
	return nil
}

// ExecuteStartedStatus executes operations depending on the auction type.
func (k Keeper) ExecuteStartedStatus(ctx context.Context, auction types.AuctionI) error {
	blockTime := sdk.UnwrapSDKContext(ctx).BlockTime()
	if auction.ShouldAuctionClosed(blockTime) { // BlockTime >= EndTime
		switch auction.GetType() {
		case types.AuctionTypeFixedPrice:
			if err := k.CloseFixedPriceAuction(ctx, auction); err != nil {
				return err
			}

		case types.AuctionTypeBatch:
			if err := k.CloseBatchAuction(ctx, auction); err != nil {
				return err
			}
		}
	}
	return nil
}

// ExecuteVestingStatus first gets all vesting queues in the store and
// look up the release time of each vesting queue to see if the module needs to
// distribute the paying coin to the auctioneer.
func (k Keeper) ExecuteVestingStatus(ctx context.Context, auction types.AuctionI) error {
	return k.ReleaseVestingPayingCoin(ctx, auction)
}
