package keeper

import (
	"context"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

// GetAllowedBiddersByAuction returns allowed bidders list for the auction.
func (k Keeper) GetAllowedBiddersByAuction(ctx context.Context, auctionId uint64) ([]types.AllowedBidder, error) {
	allowedBidders := make([]types.AllowedBidder, 0)
	rng := collections.NewPrefixedPairRange[uint64, sdk.AccAddress](auctionId)
	err := k.AllowedBidder.Walk(ctx, rng, func(_ collections.Pair[uint64, sdk.AccAddress], allowedBidder types.AllowedBidder) (bool, error) {
		allowedBidders = append(allowedBidders, allowedBidder)
		return false, nil
	})
	return allowedBidders, err
}

// AllowedBidders returns all AllowedBidder.
func (k Keeper) AllowedBidders(ctx context.Context) ([]types.AllowedBidder, error) {
	allowedBidders := make([]types.AllowedBidder, 0)
	err := k.IterateAllowedBidders(ctx, func(_ collections.Pair[uint64, sdk.AccAddress], allowedBidder types.AllowedBidder) (bool, error) {
		allowedBidders = append(allowedBidders, allowedBidder)
		return false, nil
	})
	return allowedBidders, err
}

// IterateAllowedBidders iterates over all the AllowedBidders and performs a callback function.
func (k Keeper) IterateAllowedBidders(ctx context.Context, cb func(collections.Pair[uint64, sdk.AccAddress], types.AllowedBidder) (bool, error)) error {
	err := k.AllowedBidder.Walk(ctx, nil, cb)
	if err != nil {
		return err
	}
	return nil
}
