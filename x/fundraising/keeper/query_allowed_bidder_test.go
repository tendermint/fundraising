package keeper_test

import (
	"context"
	"strconv"
	"testing"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
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

func createNAllowedBidder(keeper keeper.Keeper, ctx context.Context, n int) []types.AllowedBidder {
	items := make([]types.AllowedBidder, n)
	for i := range items {
		bidder := sample.AccAddress()
		items[i].AuctionId = uint64(i)
		items[i].Bidder = bidder.String()
		items[i].MaxBidAmount = math.ZeroInt()

		_ = keeper.AllowedBidder.Set(ctx, collections.Join(items[i].AuctionId, bidder), items[i])
	}
	return items
}

func TestAllowedBidderQuerySingle(t *testing.T) {
	k, ctx, _ := keepertest.FundraisingKeeper(t)
	qs := keeper.NewQueryServerImpl(k)
	msgs := createNAllowedBidder(k, ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetAllowedBidderRequest
		response *types.QueryGetAllowedBidderResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetAllowedBidderRequest{
				AuctionId: msgs[0].AuctionId,
				Bidder:    msgs[0].Bidder,
			},
			response: &types.QueryGetAllowedBidderResponse{AllowedBidder: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetAllowedBidderRequest{
				AuctionId: msgs[1].AuctionId,
				Bidder:    msgs[1].Bidder,
			},
			response: &types.QueryGetAllowedBidderResponse{AllowedBidder: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetAllowedBidderRequest{
				AuctionId: 100000,
				Bidder:    sample.Address(),
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := qs.GetAllowedBidder(ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.response, response)
			}
		})
	}
}

func TestAllowedBidderQueryPaginated(t *testing.T) {
	k, ctx, _ := keepertest.FundraisingKeeper(t)
	qs := keeper.NewQueryServerImpl(k)
	msgs := createNAllowedBidder(k, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllAllowedBidderRequest {
		return &types.QueryAllAllowedBidderRequest{
			AuctionId: 0,
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
			resp, err := qs.ListAllowedBidder(ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.AllowedBidder), step)
			require.Subset(t,
				msgs,
				nullify.Fill(resp.AllowedBidder),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListAllowedBidder(ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.AllowedBidder), step)
			require.Subset(t,
				msgs,
				nullify.Fill(resp.AllowedBidder),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := qs.ListAllowedBidder(ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.AllowedBidder),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := qs.ListAllowedBidder(ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
