package keeper

import (
	"fmt"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) {
	if err := genState.Validate(); err != nil {
		panic(err)
	}

	// Prevents from nil slice
	if len(genState.Params.AuctionCreationFee) == 0 {
		genState.Params.AuctionCreationFee = sdk.Coins{}
	}
	if len(genState.Params.PlaceBidFee) == 0 {
		genState.Params.PlaceBidFee = sdk.Coins{}
	}

	k.SetParams(ctx, genState.Params)

	for _, auction := range genState.Auctions {
		auction, err := types.UnpackAuction(auction)
		if err != nil {
			panic(err)
		}
		k.GetNextAuctionIdWithUpdate(ctx)
		k.SetAuction(ctx, auction)
	}

	for _, record := range genState.AllowedBidderRecords {
		k.SetAllowedBidder(ctx, record.AuctionId, record.AllowedBidder)
	}

	for _, bid := range genState.Bids {
		_, found := k.GetAuction(ctx, bid.AuctionId)
		if !found {
			panic(fmt.Sprintf("auction %d is not found", bid.AuctionId))
		}
		k.GetNextBidIdWithUpdate(ctx, bid.AuctionId)
		k.SetBid(ctx, bid)
	}

	for _, queue := range genState.VestingQueues {
		_, found := k.GetAuction(ctx, queue.AuctionId)
		if !found {
			panic(fmt.Sprintf("auction %d is not found", queue.AuctionId))
		}
		k.SetVestingQueue(ctx, queue)
	}
}

// ExportGenesis returns the module's exported genesis state.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	params := k.GetParams(ctx)
	bids := k.GetBids(ctx)
	queues := k.GetVestingQueues(ctx)

	// Prevents from nil slice
	if len(params.AuctionCreationFee) == 0 {
		params.AuctionCreationFee = sdk.Coins{}
	}
	if len(params.PlaceBidFee) == 0 {
		params.PlaceBidFee = sdk.Coins{}
	}

	auctions := []*codectypes.Any{}
	allowedBidderRecords := []types.AllowedBidderRecord{}
	for _, auction := range k.GetAuctions(ctx) {
		auctionAny, err := types.PackAuction(auction)
		if err != nil {
			panic(err)
		}
		auctions = append(auctions, auctionAny)

		if err := k.IterateAllowedBiddersByAuction(ctx, auction.GetId(), func(ab types.AllowedBidder) (stop bool, err error) {
			allowedBidderRecords = append(allowedBidderRecords, types.AllowedBidderRecord{
				AuctionId:     auction.GetId(),
				AllowedBidder: ab,
			})
			return false, nil
		}); err != nil {
			panic(err)
		}
	}

	return &types.GenesisState{
		Params:               params,
		Auctions:             auctions,
		AllowedBidderRecords: allowedBidderRecords,
		Bids:                 bids,
		VestingQueues:        queues,
	}
}
