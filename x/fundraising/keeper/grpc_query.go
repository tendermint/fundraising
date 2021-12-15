package keeper

import (
	"context"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

// Querier is used as Keeper will have duplicate methods if used directly, and gRPC names take precedence over keeper.
type Querier struct {
	Keeper
}

var _ types.QueryServer = Querier{}

// Params queries the parameters of the fundraising module.
func (k Querier) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	var params types.Params
	k.Keeper.paramSpace.GetParamSet(ctx, &params)
	return &types.QueryParamsResponse{Params: params}, nil
}

// Auctions queries all auctions.
func (k Querier) Auctions(c context.Context, req *types.QueryAuctionsRequest) (*types.QueryAuctionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.Type != "" && !(req.Type == types.AuctionTypeFixedPrice.String() || req.Type == types.AuctionTypeEnglish.String()) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid auction type %s", req.Type)
	}

	if req.Status != "" && !(req.Status == types.AuctionStatusStandBy.String() || req.Status == types.AuctionStatusStarted.String() ||
		req.Status == types.AuctionStatusVesting.String() || req.Status == types.AuctionStatusFinished.String() ||
		req.Status == types.AuctionStatusCancelled.String()) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid auction status %s", req.Status)
	}

	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(k.storeKey)
	auctionStore := prefix.NewStore(store, types.AuctionKeyPrefix)

	var auctions []*codectypes.Any
	pageRes, err := query.FilteredPaginate(auctionStore, req.Pagination, func(key, value []byte, accumulate bool) (bool, error) {
		auction, err := k.Keeper.UnmarshalAuction(value)
		if err != nil {
			return false, err
		}

		auctionAny, err := types.PackAuction(auction)
		if err != nil {
			return false, err
		}

		if req.Type != "" && auction.GetType().String() != req.Type {
			return false, nil
		}

		if req.Status != "" && auction.GetStatus().String() != req.Status {
			return false, nil
		}

		if accumulate {
			auctions = append(auctions, auctionAny)
		}

		return true, nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAuctionsResponse{Auctions: auctions, Pagination: pageRes}, nil
}

// Auction queries the specific auction.
func (k Querier) Auction(c context.Context, req *types.QueryAuctionRequest) (*types.QueryAuctionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	auction, found := k.Keeper.GetAuction(ctx, req.AuctionId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "auction %d not found", req.AuctionId)
	}

	auctionAny, err := types.PackAuction(auction)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAuctionResponse{Auction: auctionAny}, nil
}

// Bids queries all bids for the auction.
func (k Querier) Bids(c context.Context, req *types.QueryBidsRequest) (*types.QueryBidsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	_, found := k.Keeper.GetAuction(ctx, req.AuctionId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "auction %d not found", req.AuctionId)
	}

	if req.Bidder != "" {
		if _, err := sdk.AccAddressFromBech32(req.Bidder); err != nil {
			return nil, err
		}
	}

	var eligible bool
	if req.Eligible != "" {
		var err error
		eligible, err = strconv.ParseBool(req.Eligible)
		if err != nil {
			return nil, err
		}
	}

	store := ctx.KVStore(k.storeKey)
	bidStore := prefix.NewStore(store, types.BidKeyPrefix)

	var bids []types.Bid
	pageRes, err := query.FilteredPaginate(bidStore, req.Pagination, func(key, value []byte, accumulate bool) (bool, error) {
		bid, err := types.UnmarshalBid(k.cdc, value)
		if err != nil {
			return false, nil
		}

		if bid.AuctionId != req.AuctionId {
			return false, nil
		}

		if req.Bidder != "" && bid.GetBidder() != req.Bidder {
			return false, nil
		}

		if req.Eligible != "" {
			if bid.Eligible != eligible {
				return false, nil
			}
		}

		if accumulate {
			bids = append(bids, bid)
		}

		return true, nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryBidsResponse{Bids: bids, Pagination: pageRes}, nil
}

// Vestings queries all vesting queues for the auction.
func (k Querier) Vestings(c context.Context, req *types.QueryVestingsRequest) (*types.QueryVestingsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	auction, found := k.Keeper.GetAuction(ctx, req.AuctionId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "auction %d not found", req.AuctionId)
	}

	queues := k.Keeper.GetVestingQueuesByAuctionId(ctx, auction.GetId())

	return &types.QueryVestingsResponse{Vestings: queues}, nil
}
