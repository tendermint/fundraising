package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

// Querier is used as Keeper will have duplicate methods if used directly, and gRPC names take precedence over keeper.
type Querier struct {
	Keeper
}

var _ types.QueryServer = Querier{}

// Params queries the parameters of the farming module.
func (k Querier) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	var params types.Params
	k.Keeper.paramSpace.GetParamSet(ctx, &params)
	return &types.QueryParamsResponse{Params: params}, nil
}

func (k Querier) Auctions(c context.Context, _ *types.QueryAuctionsRequest) (*types.QueryAuctionsResponse, error) {
	// TODO: not implemented yet
	return &types.QueryAuctionsResponse{Auctions: nil}, nil
}

func (k Querier) Auction(c context.Context, _ *types.QueryAuctionRequest) (*types.QueryAuctionResponse, error) {
	// TODO: not implemented yet
	return &types.QueryAuctionResponse{Auction: nil}, nil
}

func (k Querier) Bids(c context.Context, _ *types.QueryBidsRequest) (*types.QueryBidsResponse, error) {
	// TODO: not implemented yet
	return &types.QueryBidsResponse{Bids: nil}, nil
}

func (k Querier) Vestings(c context.Context, _ *types.QueryVestingsRequest) (*types.QueryVestingsResponse, error) {
	// TODO: not implemented yet
	return &types.QueryVestingsResponse{Vestings: nil}, nil
}
