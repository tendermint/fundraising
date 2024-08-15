package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

func (q queryServer) ListAuction(ctx context.Context, req *types.QueryAllAuctionRequest) (*types.QueryAllAuctionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.Type != "" && !(req.Type == types.AuctionTypeFixedPrice.String() || req.Type == types.AuctionTypeBatch.String()) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid auction type %s", req.Type)
	}

	if req.Status != "" && !(req.Status == types.AuctionStatusStandBy.String() || req.Status == types.AuctionStatusStarted.String() ||
		req.Status == types.AuctionStatusVesting.String() || req.Status == types.AuctionStatusFinished.String() ||
		req.Status == types.AuctionStatusCancelled.String()) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid auction status %s", req.Status)
	}

	auctions, pageRes, err := query.CollectionFilteredPaginate(
		ctx,
		q.k.Auction,
		req.Pagination,
		func(_ uint64, auction types.AuctionI) (bool, error) {
			if req.Type != "" && auction.GetType().String() != req.Type {
				return false, nil
			}

			if req.Status != "" && auction.GetStatus().String() != req.Status {
				return false, nil
			}

			return true, nil
		},
		func(_ uint64, auction types.AuctionI) (*codectypes.Any, error) {
			auctionAny, err := types.PackAuction(auction)
			if err != nil {
				return nil, err
			}
			return auctionAny, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllAuctionResponse{Auction: auctions, Pagination: pageRes}, nil
}

func (q queryServer) GetAuction(ctx context.Context, req *types.QueryGetAuctionRequest) (*types.QueryGetAuctionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	auction, err := q.k.Auction.Get(ctx, req.AuctionId)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, sdkerrors.ErrKeyNotFound
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	auctionAny, err := types.PackAuction(auction)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryGetAuctionResponse{Auction: auctionAny}, nil
}
