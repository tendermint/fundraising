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
		sdk.NewInt64Coin(sdk.DefaultBondDenom, 100_000_000),
		sdk.NewInt64Coin(denom1, 1_000_000_000_000),
		sdk.NewInt64Coin(denom2, 1_000_000_000_000),
		sdk.NewInt64Coin(denom3, 1_000_000_000_000),
		sdk.NewInt64Coin(denom4, 1_000_000_000_000),
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
	sampleFixedPriceAuctions []*types.MsgCreateFixedPriceAuction
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
	suite.sampleFixedPriceAuctions = []*types.MsgCreateFixedPriceAuction{
		{
			Auctioneer:       suite.addrs[4].String(),
			StartPrice:       sdk.OneDec(),
			SellingCoin:      sdk.NewInt64Coin(denom1, 1_000_000_000_000),
			PayingCoinDenom:  denom2,
			VestingSchedules: []types.VestingSchedule{}, // no vesting schedules
			StartTime:        types.ParseTime("2030-01-01T00:00:00Z"),
			EndTime:          types.ParseTime("2030-01-10T00:00:00Z"),
		},
		{
			Auctioneer:      suite.addrs[5].String(),
			StartPrice:      sdk.MustNewDecFromStr("0.5"),
			SellingCoin:     sdk.NewInt64Coin(denom3, 1_000_000_000_000),
			PayingCoinDenom: denom4,
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
			StartTime: types.ParseTime("2021-12-10T00:00:00Z"),
			EndTime:   types.ParseTime("2021-12-20T00:00:00Z"),
		},
	}
	suite.sampleFixedPriceBids = []*types.MsgPlaceBid{
		{
			AuctionId: 1,
			Bidder:    suite.addrs[0].String(),
			Price:     sdk.MustNewDecFromStr("0.5"),
			Coin:      sdk.NewInt64Coin(denom4, 30_000_000),
		},
		{
			AuctionId: 1,
			Bidder:    suite.addrs[1].String(),
			Price:     sdk.MustNewDecFromStr("0.5"),
			Coin:      sdk.NewInt64Coin(denom4, 50_000_000),
		},
	}
}
