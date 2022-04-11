package simulation_test

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	chain "github.com/tendermint/fundraising/app"
	"github.com/tendermint/fundraising/app/params"
	"github.com/tendermint/fundraising/testutil/simapp"
	"github.com/tendermint/fundraising/x/fundraising/simulation"
	"github.com/tendermint/fundraising/x/fundraising/types"
)

// TestWeightedOperations tests the weights of the operations.
func TestWeightedOperations(t *testing.T) {
	app, ctx := createTestApp(false)

	ctx.WithChainID("test-chain")

	cdc := types.ModuleCdc
	appParams := make(simtypes.AppParams)

	weightedOps := simulation.WeightedOperations(appParams, cdc, app.AuthKeeper, app.BankKeeper, app.FundraisingKeeper)

	s := rand.NewSource(1)
	r := rand.New(s)
	accs := getTestingAccounts(t, r, app, ctx, 1)

	expected := []struct {
		weight     int
		opMsgRoute string
		opMsgName  string
	}{
		{params.DefaultWeightMsgCreateFixedPriceAuction, types.ModuleName, types.TypeMsgCreateFixedPriceAuction},
		{params.DefaultWeightMsgCreateBatchAuction, types.ModuleName, types.TypeMsgCreateBatchAuction},
		{params.DefaultWeightMsgCancelAuction, types.ModuleName, types.TypeMsgCancelAuction},
		{params.DefaultWeightMsgPlaceBid, types.ModuleName, types.TypeMsgPlaceBid},
	}

	for i, w := range weightedOps {
		operationMsg, _, _ := w.Op()(r, app.BaseApp, ctx, accs, ctx.ChainID())
		// the following checks are very much dependent from the ordering of the output given
		// by WeightedOperations. if the ordering in WeightedOperations changes some tests
		// will fail
		require.Equal(t, expected[i].weight, w.Weight(), "weight should be the same")
		require.Equal(t, expected[i].opMsgRoute, operationMsg.Route, "route should be the same")
		require.Equal(t, expected[i].opMsgName, operationMsg.Name, "operation Msg name should be the same")
	}
}

func TestSimulateCreateFixedPriceAuction(t *testing.T) {
	app, ctx := createTestApp(false)

	// setup a single account
	s := rand.NewSource(1)
	r := rand.New(s)

	accounts := getTestingAccounts(t, r, app, ctx, 1)

	app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: app.LastBlockHeight() + 1, AppHash: app.LastCommitID().Hash}})

	op := simulation.SimulateMsgCreateFixedPriceAuction(app.AuthKeeper, app.BankKeeper, app.FundraisingKeeper)
	opMsg, futureOps, err := op(r, app.BaseApp, ctx, accounts, "")
	require.NoError(t, err)
	require.True(t, opMsg.OK)
	require.Len(t, futureOps, 0)

	var msg types.MsgCreateFixedPriceAuction
	types.ModuleCdc.MustUnmarshalJSON(opMsg.Msg, &msg)

	require.Equal(t, types.TypeMsgCreateFixedPriceAuction, msg.Type())
	require.Equal(t, types.ModuleName, msg.Route())
	require.Equal(t, "cosmos1tnh2q55v8wyygtt9srz5safamzdengsnqeycj3", msg.Auctioneer)
	require.Equal(t, "denomc", msg.SellingCoin.Denom)
	require.Equal(t, sdk.DefaultBondDenom, msg.PayingCoinDenom)
	require.Equal(t, sdk.MustNewDecFromStr("0.5"), msg.StartPrice)
}

func TestSimulateCreateBatchAuction(t *testing.T) {
	app, ctx := createTestApp(false)

	// setup a single account
	s := rand.NewSource(1)
	r := rand.New(s)

	accounts := getTestingAccounts(t, r, app, ctx, 1)

	app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: app.LastBlockHeight() + 1, AppHash: app.LastCommitID().Hash}})

	op := simulation.SimulateMsgCreateBatchAuction(app.AuthKeeper, app.BankKeeper, app.FundraisingKeeper)
	opMsg, futureOps, err := op(r, app.BaseApp, ctx, accounts, "")
	require.NoError(t, err)
	require.True(t, opMsg.OK)
	require.Len(t, futureOps, 0)

	var msg types.MsgCreateBatchAuction
	types.ModuleCdc.MustUnmarshalJSON(opMsg.Msg, &msg)

	require.Equal(t, types.TypeMsgCreateBatchAuction, msg.Type())
	require.Equal(t, types.ModuleName, msg.Route())
	require.Equal(t, "cosmos1tnh2q55v8wyygtt9srz5safamzdengsnqeycj3", msg.Auctioneer)
	require.Equal(t, "denom10", msg.SellingCoin.Denom)
	require.Equal(t, "stake", msg.PayingCoinDenom)
	require.Equal(t, uint32(3), msg.MaxExtendedRound)
	require.Equal(t, sdk.MustNewDecFromStr("0.04"), msg.ExtendedRoundRate)
}

func TestSimulateCancelAuction(t *testing.T) {
	app, ctx := createTestApp(false)

	// setup a single account
	s := rand.NewSource(1)
	r := rand.New(s)

	accounts := getTestingAccounts(t, r, app, ctx, 1)

	app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: app.LastBlockHeight() + 1, AppHash: app.LastCommitID().Hash}})

	// Create a fixed price auction
	_, err := app.FundraisingKeeper.CreateFixedPriceAuction(ctx, &types.MsgCreateFixedPriceAuction{
		Auctioneer:       accounts[0].Address.String(),
		StartPrice:       sdk.MustNewDecFromStr("0.5"),
		SellingCoin:      sdk.NewInt64Coin("denom1", 5000000000),
		PayingCoinDenom:  "denom2",
		VestingSchedules: []types.VestingSchedule{},
		StartTime:        ctx.BlockTime().AddDate(0, 1, 0),
		EndTime:          ctx.BlockTime().AddDate(0, 2, 0),
	})
	require.NoError(t, err)

	op := simulation.SimulateMsgCancelAuction(app.AuthKeeper, app.BankKeeper, app.FundraisingKeeper)
	opMsg, futureOps, err := op(r, app.BaseApp, ctx, accounts, "")
	require.NoError(t, err)
	require.True(t, opMsg.OK)
	require.Len(t, futureOps, 0)

	var msg types.MsgCancelAuction
	types.ModuleCdc.MustUnmarshalJSON(opMsg.Msg, &msg)

	require.Equal(t, types.TypeMsgCancelAuction, msg.Type())
	require.Equal(t, types.ModuleName, msg.Route())
	require.Equal(t, "cosmos1tnh2q55v8wyygtt9srz5safamzdengsnqeycj3", msg.Auctioneer)
	require.Equal(t, uint64(1), msg.AuctionId)
}

func TestSimulatePlaceBid(t *testing.T) {
	app, ctx := createTestApp(false)

	// Setup a single account
	s := rand.NewSource(1)
	r := rand.New(s)

	accounts := getTestingAccounts(t, r, app, ctx, 1)

	app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: app.LastBlockHeight() + 1, AppHash: app.LastCommitID().Hash}})

	// Create a fixed price auction
	_, err := app.FundraisingKeeper.CreateFixedPriceAuction(ctx, &types.MsgCreateFixedPriceAuction{
		Auctioneer:       accounts[0].Address.String(),
		StartPrice:       sdk.MustNewDecFromStr("0.5"),
		SellingCoin:      sdk.NewInt64Coin("denom1", 5000000000),
		PayingCoinDenom:  "denom2",
		VestingSchedules: []types.VestingSchedule{},
		StartTime:        ctx.BlockTime(),
		EndTime:          ctx.BlockTime().AddDate(0, 1, 0),
	})
	require.NoError(t, err)

	// Create a batch auction
	_, err = app.FundraisingKeeper.CreateBatchAuction(ctx, &types.MsgCreateBatchAuction{
		Auctioneer:        accounts[0].Address.String(),
		StartPrice:        sdk.MustNewDecFromStr("0.5"),
		MinBidPrice:       sdk.MustNewDecFromStr("0.1"),
		SellingCoin:       sdk.NewInt64Coin("denom3", 5000000000),
		PayingCoinDenom:   "denom4",
		MaxExtendedRound:  3,
		ExtendedRoundRate: sdk.MustNewDecFromStr("0.1"),
		VestingSchedules:  []types.VestingSchedule{},
		StartTime:         ctx.BlockTime(),
		EndTime:           ctx.BlockTime().AddDate(0, 1, 0),
	})
	require.NoError(t, err)

	op := simulation.SimulateMsgPlaceBid(app.AuthKeeper, app.BankKeeper, app.FundraisingKeeper)
	opMsg, futureOps, err := op(r, app.BaseApp, ctx, accounts, "")
	require.NoError(t, err)
	require.True(t, opMsg.OK)
	require.Len(t, futureOps, 0)

	var msg types.MsgPlaceBid
	types.ModuleCdc.MustUnmarshalJSON(opMsg.Msg, &msg)

	require.Equal(t, types.TypeMsgPlaceBid, msg.Type())
	require.Equal(t, types.ModuleName, msg.Route())
	require.Equal(t, "cosmos1tnh2q55v8wyygtt9srz5safamzdengsnqeycj3", msg.Bidder)
	require.Equal(t, uint64(2), msg.AuctionId)
	require.Equal(t, types.BidTypeBatchWorth, msg.BidType)
	require.Equal(t, sdk.MustNewDecFromStr("1.3"), msg.Price)
	require.Equal(t, sdk.NewInt64Coin("denom4", 336222540), msg.Coin)
}

func createTestApp(isCheckTx bool) (*chain.App, sdk.Context) {
	app := simapp.New(chain.DefaultNodeHome)

	ctx := app.BaseApp.NewContext(isCheckTx, tmproto.Header{})
	app.MintKeeper.SetParams(ctx, minttypes.DefaultParams())

	return app, ctx
}

func getTestingAccounts(t *testing.T, r *rand.Rand, app *chain.App, ctx sdk.Context, n int) []simtypes.Account {
	accs := simtypes.RandomAccounts(r, n)

	initAmt := app.StakingKeeper.TokensFromConsensusPower(ctx, 500)
	coins := sdk.NewCoins(
		sdk.NewCoin(sdk.DefaultBondDenom, initAmt),
		sdk.NewInt64Coin("denoma", 1_000_000_000_000_000),
		sdk.NewInt64Coin("denomb", 1_000_000_000_000_000),
		sdk.NewInt64Coin("denomc", 1_000_000_000_000_000),
		sdk.NewInt64Coin("denomd", 1_000_000_000_000_000),
	)

	// add coins to the accounts
	for _, acc := range accs {
		acc := app.AuthKeeper.NewAccountWithAddress(ctx, acc.Address)
		app.AuthKeeper.SetAccount(ctx, acc)

		err := app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, coins)
		require.NoError(t, err)

		err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, acc.GetAddress(), coins)
		require.NoError(t, err)
	}

	return accs
}
