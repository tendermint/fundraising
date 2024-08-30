package keeper

import (
	"context"
	"time"

	"cosmossdk.io/collections"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

func (q queryServer) ListVestingQueue(ctx context.Context, req *types.QueryAllVestingQueueRequest) (*types.QueryAllVestingQueueResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	vestingQueues, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.VestingQueue,
		req.Pagination,
		func(_ collections.Pair[uint64, time.Time], value types.VestingQueue) (types.VestingQueue, error) {
			return value, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllVestingQueueResponse{VestingQueue: vestingQueues, Pagination: pageRes}, nil
}
