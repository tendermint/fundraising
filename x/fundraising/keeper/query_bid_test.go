package keeper_test

import (
	"context"
	"testing"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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

func createNBid(keeper keeper.Keeper, ctx context.Context, n int) ([]types.Bid, error) {
	items := make([]types.Bid, n)
	auctionId := uint64(0)
	for i := range items {
		bidid := uint64(i)
		items[i].AuctionId = auctionId
		items[i].Id = bidid
		items[i].Bidder = sample.Address()
		items[i].Coin = sdk.NewCoin("coin", math.NewInt(int64(i)))
		items[i].Price = math.LegacyNewDec(int64(i))
		items[i].Type = types.BidTypeFixedPrice

		if err := keeper.Bid.Set(ctx, collections.Join(auctionId, bidid), items[i]); err != nil {
			return nil, err
		}
		if err := keeper.BidSeq.Set(ctx, auctionId, items[i].Id); err != nil {
			return nil, err
		}
	}
	return items, nil
}

func TestBidQuerySingle(t *testing.T) {
	k, ctx, _ := keepertest.FundraisingKeeper(t)
	qs := keeper.NewQueryServerImpl(k)
	msgs, err := createNBid(k, ctx, 2)
	require.NoError(t, err)

	tests := []struct {
		desc     string
		request  *types.QueryGetBidRequest
		response *types.QueryGetBidResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetBidRequest{AuctionId: 0, BidId: msgs[0].Id},
			response: &types.QueryGetBidResponse{Bid: msgs[0]},
		},
		{
			desc:     "Second",
			request:  &types.QueryGetBidRequest{AuctionId: 0, BidId: msgs[1].Id},
			response: &types.QueryGetBidResponse{Bid: msgs[1]},
		},
		{
			desc:    "KeyNotFound",
			request: &types.QueryGetBidRequest{AuctionId: 0, BidId: uint64(len(msgs))},
			err:     sdkerrors.ErrKeyNotFound,
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := qs.GetBid(ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t,
					nullify.Fill(tc.response),
					nullify.Fill(response),
				)
			}
		})
	}
}

func TestBidQueryPaginated(t *testing.T) {
	k, ctx, _ := keepertest.FundraisingKeeper(t)
	qs := keeper.NewQueryServerImpl(k)
	msgs, err := createNBid(k, ctx, 5)
	require.NoError(t, err)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllBidRequest {
		return &types.QueryAllBidRequest{
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
			resp, err := qs.ListBid(ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Bid), step)
			require.Subset(t, msgs, resp.Bid)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListBid(ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Bid), step)
			require.Subset(t, msgs, resp.Bid)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := qs.ListBid(ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t, msgs, resp.Bid)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := qs.ListBid(ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
