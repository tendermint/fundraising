package keeper_test

import (
	"context"
	"testing"
	"time"

	"cosmossdk.io/math"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/tendermint/fundraising/testutil/keeper"
	"github.com/tendermint/fundraising/x/fundraising/keeper"
	"github.com/tendermint/fundraising/x/fundraising/types"
)

func createNAuction(keeper keeper.Keeper, ctx context.Context, n int) ([]*codectypes.Any, error) {
	var err error
	items := make([]*codectypes.Any, n)
	for i := range items {
		iu := uint64(i)
		auction := &types.FixedPriceAuction{
			BaseAuction: &types.BaseAuction{
				Id:          iu,
				Auctioneer:  "",
				StartPrice:  math.LegacyMustNewDecFromStr("10"),
				SellingCoin: sdk.NewCoin("coin", math.NewInt(5)),
				StartTime:   time.Now(),
			},
			RemainingSellingCoin: sdk.NewCoin("coin", math.NewInt(1)),
		}
		items[i], err = types.PackAuction(auction)
		if err != nil {
			return nil, err
		}
		if err := keeper.Auction.Set(ctx, iu, auction); err != nil {
			return nil, err
		}
		if err := keeper.AuctionSeq.Set(ctx, iu); err != nil {
			return nil, err
		}
	}
	return items, nil
}

func auctionsToString(auctionsAny []*codectypes.Any) (auctions []string) {
	for _, auction := range auctionsAny {
		auctions = append(auctions, auction.String())
	}
	return
}

func TestAuctionQuerySingle(t *testing.T) {
	k, ctx, _ := keepertest.FundraisingKeeper(t)
	qs := keeper.NewQueryServerImpl(k)
	msgs, err := createNAuction(k, ctx, 2)
	require.NoError(t, err)

	tests := []struct {
		desc     string
		request  *types.QueryGetAuctionRequest
		response *types.QueryGetAuctionResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetAuctionRequest{AuctionId: 0},
			response: &types.QueryGetAuctionResponse{Auction: msgs[0]},
		},
		{
			desc:     "Second",
			request:  &types.QueryGetAuctionRequest{AuctionId: 1},
			response: &types.QueryGetAuctionResponse{Auction: msgs[1]},
		},
		{
			desc:    "KeyNotFound",
			request: &types.QueryGetAuctionRequest{AuctionId: uint64(len(msgs))},
			err:     sdkerrors.ErrKeyNotFound,
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := qs.GetAuction(ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.response.String(), response.String())
			}
		})
	}
}

func TestAuctionQueryPaginated(t *testing.T) {
	k, ctx, _ := keepertest.FundraisingKeeper(t)
	qs := keeper.NewQueryServerImpl(k)
	msgs, err := createNAuction(k, ctx, 5)
	require.NoError(t, err)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllAuctionRequest {
		return &types.QueryAllAuctionRequest{
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
			resp, err := qs.ListAuction(ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Auction), step)
			require.Subset(t,
				auctionsToString(msgs),
				auctionsToString(resp.Auction),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListAuction(ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Auction), step)
			require.Subset(t,
				auctionsToString(msgs),
				auctionsToString(resp.Auction),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := qs.ListAuction(ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			auctionsToString(msgs),
			auctionsToString(resp.Auction),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := qs.ListAuction(ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
