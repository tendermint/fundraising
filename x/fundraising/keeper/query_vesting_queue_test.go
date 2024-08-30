package keeper_test

import (
	"context"
	"strconv"
	"testing"
	"time"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/tendermint/fundraising/testutil/keeper"
	"github.com/tendermint/fundraising/testutil/nullify"
	"github.com/tendermint/fundraising/testutil/sample"
	"github.com/tendermint/fundraising/x/fundraising/keeper"
	"github.com/tendermint/fundraising/x/fundraising/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNVestingQueue(keeper keeper.Keeper, ctx context.Context, n int) []types.VestingQueue {
	items := make([]types.VestingQueue, n)
	for i := range items {
		items[i].AuctionId = uint64(i)
		items[i].ReleaseTime = time.Now().UTC()
		items[i].PayingCoin = sdk.NewCoin("coin", math.NewInt(int64(i)))
		items[i].Auctioneer = sample.Address()

		_ = keeper.VestingQueue.Set(ctx, collections.Join(items[i].AuctionId, items[i].ReleaseTime), items[i])
	}
	return items
}

func TestVestingQueueQueryPaginated(t *testing.T) {
	k, ctx, _ := keepertest.FundraisingKeeper(t)
	qs := keeper.NewQueryServerImpl(k)
	msgs := createNVestingQueue(k, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllVestingQueueRequest {
		return &types.QueryAllVestingQueueRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListVestingQueue(ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.VestingQueue), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.VestingQueue),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListVestingQueue(ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.VestingQueue), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.VestingQueue),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := qs.ListVestingQueue(ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.VestingQueue),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := qs.ListVestingQueue(ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
