package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	keepertest "github.com/tendermint/fundraising/testutil/keeper"
	"github.com/tendermint/fundraising/x/fundraising/keeper"
	"github.com/tendermint/fundraising/x/fundraising/types"
)

func TestParamsQuery(t *testing.T) {
	k, ctx, _ := keepertest.FundraisingKeeper(t)

	qs := keeper.NewQueryServerImpl(k)
	params := types.DefaultParams()
	require.NoError(t, k.Params.Set(ctx, params))

	response, err := qs.Params(ctx, &types.QueryParamsRequest{})
	require.NoError(t, err)

	// Prevents from nil slice
	if len(response.Params.AuctionCreationFee) == 0 {
		response.Params.AuctionCreationFee = sdk.Coins{}
	}
	if len(response.Params.PlaceBidFee) == 0 {
		response.Params.PlaceBidFee = sdk.Coins{}
	}

	require.Equal(t, &types.QueryParamsResponse{Params: params}, response)
}
