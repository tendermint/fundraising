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

	k.SetParams(ctx, genState.Params)

	for _, auction := range genState.Auctions {
		auction, err := types.UnpackAuction(auction)
		if err != nil {
			panic(err)
		}
		k.GetNextAuctionIdWithUpdate(ctx)
		k.SetAuction(ctx, auction)
	}

	for _, allowedBidder := range genState.AllowedBidders {
		k.SetAllowedBidder(ctx, allowedBidder)
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

	auctions := []*codectypes.Any{}
	allowedBidders := []types.AllowedBidder{}
	for _, auction := range k.GetAuctions(ctx) {
		auctionAny, err := types.PackAuction(auction)
		if err != nil {
			panic(err)
		}
		auctions = append(auctions, auctionAny)

		allowedBidders = append(allowedBidders, k.GetAllowedBiddersByAuction(ctx, auction.GetId())...)
	}

	return &types.GenesisState{
		Params:         params,
		Auctions:       auctions,
		AllowedBidders: allowedBidders,
		Bids:           bids,
		VestingQueues:  queues,
	}
}
