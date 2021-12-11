package fundraising_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/tendermint/fundraising/app"
	"github.com/tendermint/fundraising/testutil/simapp"
	"github.com/tendermint/fundraising/x/fundraising/keeper"
	"github.com/tendermint/fundraising/x/fundraising/types"
)

const (
	denom1 = "denom1" // selling coin denom
	denom2 = "denom2" // paying coin denom
	denom3 = "denom3"
	denom4 = "denom4"
)

var (
	initialBalances = sdk.NewCoins(
		sdk.NewInt64Coin(sdk.DefaultBondDenom, 100_000_000_000_000),
		sdk.NewInt64Coin(denom1, 100_000_000_000_000),
		sdk.NewInt64Coin(denom2, 100_000_000_000_000),
		sdk.NewInt64Coin(denom3, 100_000_000_000_000),
		sdk.NewInt64Coin(denom4, 100_000_000_000_000),
	)
)

type ModuleTestSuite struct {
	suite.Suite

	app                      *app.App
	ctx                      sdk.Context
	keeper                   keeper.Keeper
	querier                  keeper.Querier
	srv                      types.MsgServer
	addrs                    []sdk.AccAddress
	sampleFixedPriceAuctions []types.AuctionI
	sampleFixedPriceBids     []*types.MsgPlaceBid
}

func TestModuleTestSuite(t *testing.T) {
	suite.Run(t, new(ModuleTestSuite))
}

func (suite *ModuleTestSuite) SetupTest() {
	app := simapp.New(app.DefaultNodeHome)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	suite.app = app
	suite.ctx = ctx
	suite.ctx = suite.ctx.WithBlockTime(time.Now()) // set to current time
	suite.keeper = suite.app.FundraisingKeeper
	suite.querier = keeper.Querier{Keeper: suite.keeper}
	suite.srv = keeper.NewMsgServerImpl(suite.keeper)
	suite.addrs = simapp.AddTestAddrs(suite.app, suite.ctx, 6, sdk.ZeroInt())
	for _, addr := range suite.addrs {
		err := simapp.FundAccount(suite.app.BankKeeper, suite.ctx, addr, initialBalances)
		suite.Require().NoError(err)
	}
	suite.sampleFixedPriceAuctions = []types.AuctionI{
		types.NewFixedPriceAuction(
			&types.BaseAuction{
				Id:                 1,
				Type:               types.AuctionTypeFixedPrice,
				Auctioneer:         suite.addrs[4].String(),
				SellingPoolAddress: types.SellingReserveAcc(1).String(),
				PayingPoolAddress:  types.PayingReserveAcc(1).String(),
				StartPrice:         sdk.OneDec(), // start price corresponds to ratio of the paying coin
				SellingCoin:        sdk.NewInt64Coin(denom1, 1_000_000_000_000),
				PayingCoinDenom:    denom2,
				VestingAddress:     types.VestingReserveAcc(1).String(),
				VestingSchedules:   []types.VestingSchedule{},
				WinningPrice:       sdk.ZeroDec(),
				RemainingCoin:      sdk.NewInt64Coin(denom1, 1_000_000_000_000),
				StartTime:          types.ParseTime("2021-12-20T00:00:00Z"),
				EndTimes:           []time.Time{types.ParseTime("2021-12-27T00:00:00Z")},
				Status:             types.AuctionStatusStandBy,
			},
		),
		types.NewFixedPriceAuction(
			&types.BaseAuction{
				Id:                 2,
				Type:               types.AuctionTypeFixedPrice,
				Auctioneer:         suite.addrs[5].String(),
				SellingPoolAddress: types.SellingReserveAcc(1).String(),
				PayingPoolAddress:  types.PayingReserveAcc(1).String(),
				StartPrice:         sdk.MustNewDecFromStr("0.5"),
				SellingCoin:        sdk.NewInt64Coin(denom3, 1_000_000_000_000),
				PayingCoinDenom:    denom4,
				VestingAddress:     types.VestingReserveAcc(1).String(),
				VestingSchedules: []types.VestingSchedule{
					{
						ReleaseTime: types.ParseTime("2022-01-01T00:00:00Z"),
						Weight:      sdk.MustNewDecFromStr("0.25"),
					},
					{
						ReleaseTime: types.ParseTime("2022-04-01T00:00:00Z"),
						Weight:      sdk.MustNewDecFromStr("0.25"),
					},
					{
						ReleaseTime: types.ParseTime("2022-08-01T00:00:00Z"),
						Weight:      sdk.MustNewDecFromStr("0.25"),
					},
					{
						ReleaseTime: types.ParseTime("2022-12-01T00:00:00Z"),
						Weight:      sdk.MustNewDecFromStr("0.25"),
					},
				},
				WinningPrice:  sdk.ZeroDec(),
				RemainingCoin: sdk.NewInt64Coin(denom3, 1_000_000_000_000),
				StartTime:     types.ParseTime("2021-12-01T00:00:00Z"),
				EndTimes:      []time.Time{types.ParseTime("2022-12-20T00:00:00Z")},
				Status:        types.AuctionStatusStarted,
			},
		),
	}
	suite.sampleFixedPriceBids = []*types.MsgPlaceBid{
		types.NewMsgPlaceBid(
			1,
			suite.addrs[0].String(),
			sdk.OneDec(),
			sdk.NewInt64Coin(denom2, 50_000_000),
		),
		types.NewMsgPlaceBid(
			1,
			suite.addrs[1].String(),
			sdk.OneDec(),
			sdk.NewInt64Coin(denom2, 100_000_000),
		),
	}
}

// VestingSchedules is a convenient method to test
func (suite *ModuleTestSuite) VestingSchedules() []types.VestingSchedule {
	return []types.VestingSchedule{
		{
			ReleaseTime: types.ParseTime("2022-01-01T22:00:00+00:00"),
			Weight:      sdk.MustNewDecFromStr("0.25")},
		{
			ReleaseTime: types.ParseTime("2022-04-01T22:00:00+00:00"),
			Weight:      sdk.MustNewDecFromStr("0.25")},
		{
			ReleaseTime: types.ParseTime("2022-08-01T22:00:00+00:00"),
			Weight:      sdk.MustNewDecFromStr("0.25")},
		{
			ReleaseTime: types.ParseTime("2022-12-01T22:00:00+00:00"),
			Weight:      sdk.MustNewDecFromStr("0.25"),
		},
	}
}
