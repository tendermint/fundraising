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
	// TODO: not implemented yet
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
