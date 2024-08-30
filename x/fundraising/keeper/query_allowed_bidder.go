package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

func (q queryServer) ListAllowedBidder(ctx context.Context, req *types.QueryAllAllowedBidderRequest) (*types.QueryAllAllowedBidderResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	allowedBidders, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.AllowedBidder,
		req.Pagination,
		func(_ collections.Pair[uint64, sdk.AccAddress], value types.AllowedBidder) (types.AllowedBidder, error) {
			return value, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllAllowedBidderResponse{AllowedBidder: allowedBidders, Pagination: pageRes}, nil
}

func (q queryServer) GetAllowedBidder(ctx context.Context, req *types.QueryGetAllowedBidderRequest) (*types.QueryGetAllowedBidderResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	bidder, err := sdk.AccAddressFromBech32(req.Bidder)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid bidder")
	}

	val, err := q.k.AllowedBidder.Get(ctx, collections.Join(req.AuctionId, bidder))
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryGetAllowedBidderResponse{AllowedBidder: val}, nil
}
