package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	testkeeper "github.com/tendermint/fundraising/testutil/keeper"
	"github.com/tendermint/fundraising/x/fundraising/keeper"
	"github.com/tendermint/fundraising/x/fundraising/types"
)

func setupMsgServer(t testing.TB) (
	sdk.Context,
	*keeper.Keeper,
	types.MsgServer,
) {
	ctx, fundraisingKeeper := testkeeper.Fundraising(t)

	return ctx, fundraisingKeeper, keeper.NewMsgServerImpl(*fundraisingKeeper)
}

func TestMsgCreateFixedPriceAuction(t *testing.T) {
	sdkCtx, k, srv := setupMsgServer(t)
	ctx := sdk.WrapSDKContext(sdkCtx)

	for _, tc := range []struct {
		name string
		msg  *types.MsgCreateFixedPriceAuction
		expectedErr error
	}{
		{
			"valid message",
			types.NewMsgCreateFixedPriceAuction(
				sample.Address(), 
				sample.StartPrice("1.0"), 
				sample.SellingCoin("ugex", 1_000_000_000_000), 
				sample.PayingCoinDenom(), 
				vestingSchedules []types.VestingSchedule, 
				startTime time.Time, 
				endTime time.Time,
			),
			nil
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// TODO: not implemented yet
			srv.CreateFixedPriceAuction(ctx, tc.msg)
		})
	}

}

func TestMsgCreateEnglishAuction(t *testing.T) {
	// TODO: not implemented yet
}

func TestMsgCancelAuction(t *testing.T) {
	// TODO: not implemented yet
}

func TestMsgPlaceBid(t *testing.T) {
	// TODO: not implemented yet
}
