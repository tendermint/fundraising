package fundraising

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cosmossdk.io/collections"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/fundraising/x/fundraising/keeper"
	"github.com/tendermint/fundraising/x/fundraising/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx context.Context, k keeper.Keeper, genState types.GenesisState) error {
	// Prevents from nil slice
	if len(genState.Params.AuctionCreationFee) == 0 {
		genState.Params.AuctionCreationFee = sdk.Coins{}
	}
	if len(genState.Params.PlaceBidFee) == 0 {
		genState.Params.PlaceBidFee = sdk.Coins{}
	}

	// Set all the auction
	for _, elem := range genState.AuctionList {
		auction, err := types.UnpackAuction(elem)
		if err != nil {
			return err
		}
		auctionID, err := k.AuctionSeq.Next(ctx)
		if err != nil {
			return err
		}
		if err := auction.SetId(auctionID); err != nil {
			return err
		}
		if err := k.Auction.Set(ctx, auctionID, auction); err != nil {
			return err
		}
	}

	// Set all the allowedBidder
	for _, elem := range genState.AllowedBidderList {
		bidder, err := sdk.AccAddressFromBech32(elem.Bidder)
		if err != nil {
			return err
		}
		if err := k.AllowedBidder.Set(ctx, collections.Join(elem.AuctionId, bidder), elem); err != nil {
			return err
		}
	}

	// Set all the bid
	for _, elem := range genState.BidList {
		_, err := k.Auction.Get(ctx, elem.AuctionId)
		if errors.Is(err, collections.ErrNotFound) {
			return fmt.Errorf("bid auction %d is not found", elem.AuctionId)
		}

		bidID, err := k.GetNextBidIdWithUpdate(ctx, elem.AuctionId)
		if err != nil {
			return err
		}
		elem.Id = bidID
		if err := k.Bid.Set(ctx, collections.Join(elem.AuctionId, elem.Id), elem); err != nil {
			return err
		}
	}

	// Set all the vestingQueue
	for _, elem := range genState.VestingQueueList {
		_, err := k.Auction.Get(ctx, elem.AuctionId)
		if errors.Is(err, collections.ErrNotFound) {
			return fmt.Errorf("vesting queue auction %d is not found", elem.AuctionId)
		}

		if err := k.VestingQueue.Set(
			ctx,
			collections.Join(
				elem.AuctionId,
				elem.ReleaseTime,
			),
			elem); err != nil {
			return err
		}
	}

	// this line is used by starport scaffolding # genesis/module/init

	return k.Params.Set(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx context.Context, k keeper.Keeper) (*types.GenesisState, error) {
	var err error

	genesis := types.DefaultGenesis()
	genesis.Params, err = k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}
	// Prevents from nil slice
	if len(genesis.Params.AuctionCreationFee) == 0 {
		genesis.Params.AuctionCreationFee = sdk.Coins{}
	}
	if len(genesis.Params.PlaceBidFee) == 0 {
		genesis.Params.PlaceBidFee = sdk.Coins{}
	}

	if err := k.AllowedBidder.Walk(ctx, nil, func(_ collections.Pair[uint64, sdk.AccAddress], val types.AllowedBidder) (bool, error) {
		genesis.AllowedBidderList = append(genesis.AllowedBidderList, val)
		return false, nil
	}); err != nil {
		return nil, err
	}
	if err := k.VestingQueue.Walk(ctx, nil, func(key collections.Pair[uint64, time.Time], val types.VestingQueue) (bool, error) {
		genesis.VestingQueueList = append(genesis.VestingQueueList, val)
		return false, nil
	}); err != nil {
		return nil, err
	}

	err = k.Bid.Walk(ctx, nil, func(key collections.Pair[uint64, uint64], val types.Bid) (bool, error) {
		genesis.BidList = append(genesis.BidList, val)
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	genesis.AuctionList = make([]*codectypes.Any, 0)
	err = k.Auction.Walk(ctx, nil, func(key uint64, elem types.AuctionI) (bool, error) {
		auctionAny, err := types.PackAuction(elem)
		if err != nil {
			panic(err)
		}
		genesis.AuctionList = append(genesis.AuctionList, auctionAny)
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	// this line is used by starport scaffolding # genesis/module/export

	return genesis, nil
}
