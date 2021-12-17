package fundraising

import (
	"fmt"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/fundraising/x/fundraising/keeper"
	"github.com/tendermint/fundraising/x/fundraising/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	if err := genState.Validate(); err != nil {
		panic(err)
	}

	k.SetParams(ctx, genState.Params)

	for i, auction := range genState.Auctions {
		auction, err := types.UnpackAuction(auction)
		if err != nil {
			panic(err)
		}
		k.SetAuction(ctx, auction)

		if i == len(genState.Auctions)-1 {
			k.SetAuctionId(ctx, auction.GetId())
		}
	}

	for _, bid := range genState.Bids {
		bidderAcc, err := sdk.AccAddressFromBech32(bid.Bidder)
		if err != nil {
			panic(err)
		}

		_, found := k.GetAuction(ctx, bid.AuctionId)
		if !found {
			panic(fmt.Sprintf("auction %d is not found", bid.AuctionId))
		}

		k.SetBid(ctx, bid.AuctionId, bid.Sequence, bidderAcc, bid)
	}

	for _, queue := range genState.VestingQueues {
		_, found := k.GetAuction(ctx, queue.AuctionId)
		if !found {
			panic(fmt.Sprintf("auction %d is not found", queue.AuctionId))
		}

		k.SetVestingQueue(ctx, queue.AuctionId, queue.ReleaseTime, queue)
	}
}

// ExportGenesis returns the module's exported genesis state.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	params := k.GetParams(ctx)
	bids := k.GetBids(ctx)
	queues := k.GetVestingQueues(ctx)

	auctions := []*codectypes.Any{}
	for _, auction := range k.GetAuctions(ctx) {
		auctionAny, err := types.PackAuction(auction)
		if err != nil {
			panic(err)
		}
		auctions = append(auctions, auctionAny)
	}

	return &types.GenesisState{
		Params:        params,
		Auctions:      auctions,
		Bids:          bids,
		VestingQueues: queues,
	}
}
