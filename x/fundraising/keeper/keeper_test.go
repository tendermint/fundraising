package keeper_test

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
	denom3 = "denom3" // selling coin denom
	denom4 = "denom4" // paying coin denom
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

type KeeperTestSuite struct {
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

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) SetupTest() {
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
			types.NewBaseAuction(
				1,
				types.AuctionTypeFixedPrice,
				suite.addrs[4].String(),
				types.SellingReserveAcc(1).String(),
				types.PayingReserveAcc(1).String(),
				suite.StartPrice("1.0"), // 1:1 price of paying coin
				suite.SellingCoin(denom1, 1_000_000_000_000),
				suite.PayingCoinDenom(denom2),
				types.VestingReserveAcc(1).String(),
				[]types.VestingSchedule{}, // no vesting schedules
				sdk.ZeroDec(),
				suite.TotalSellingCoin(1_000_000_000_000),
				types.ParseTime("2021-12-01T00:00:00Z"),
				[]time.Time{types.ParseTime("2022-01-01T00:00:00Z")},
				types.AuctionStatusStandBy,
			),
		),
		types.NewFixedPriceAuction(
			types.NewBaseAuction(
				2,
				types.AuctionTypeFixedPrice,
				suite.addrs[5].String(),
				types.SellingReserveAcc(1).String(),
				types.PayingReserveAcc(1).String(),
				suite.StartPrice("0.5"), // half price of paying coin
				suite.SellingCoin(denom3, 1_000_000_000_000),
				suite.PayingCoinDenom(denom4),
				types.VestingReserveAcc(1).String(),
				[]types.VestingSchedule{
					types.NewVestingSchedule(types.ParseTime("2022-01-01T00:00:00Z"), sdk.MustNewDecFromStr("0.25")),
					types.NewVestingSchedule(types.ParseTime("2022-04-01T00:00:00Z"), sdk.MustNewDecFromStr("0.25")),
					types.NewVestingSchedule(types.ParseTime("2022-08-01T00:00:00Z"), sdk.MustNewDecFromStr("0.25")),
					types.NewVestingSchedule(types.ParseTime("2022-12-01T00:00:00Z"), sdk.MustNewDecFromStr("0.25")),
				},
				sdk.ZeroDec(),
				suite.TotalSellingCoin(1_000_000_000_000),
				types.ParseTime("2021-12-01T00:00:00Z"),
				[]time.Time{types.ParseTime("2022-12-12T00:00:00Z")},
				types.AuctionStatusStandBy,
			),
		),
	}
	suite.sampleFixedPriceBids = []*types.MsgPlaceBid{
		types.NewMsgPlaceBid(
			1,
			suite.addrs[0].String(),
			suite.Price("1.0"),
			suite.Coin(denom2, 50_000_000),
		),
		types.NewMsgPlaceBid(
			1,
			suite.addrs[1].String(),
			suite.Price("1.0"),
			suite.Coin(denom2, 50_000_000),
		),
	}
}

func (suite *KeeperTestSuite) StartPrice(price string) sdk.Dec {
	return sdk.MustNewDecFromStr(price)
}

func (suite *KeeperTestSuite) SellingCoin(denom string, amount int64) sdk.Coin {
	return sdk.NewInt64Coin(denom, amount)
}

func (suite *KeeperTestSuite) PayingCoinDenom(denom string) string {
	return denom
}

func (suite *KeeperTestSuite) VestingSchedules() []types.VestingSchedule {
	return []types.VestingSchedule{
		types.NewVestingSchedule(types.ParseTime("2022-01-01T22:00:00+00:00"), sdk.MustNewDecFromStr("0.25")),
		types.NewVestingSchedule(types.ParseTime("2022-04-01T22:00:00+00:00"), sdk.MustNewDecFromStr("0.25")),
		types.NewVestingSchedule(types.ParseTime("2022-08-01T22:00:00+00:00"), sdk.MustNewDecFromStr("0.25")),
		types.NewVestingSchedule(types.ParseTime("2022-12-01T22:00:00+00:00"), sdk.MustNewDecFromStr("0.25")),
	}
}

func (suite *KeeperTestSuite) TotalSellingCoin(amount int64) sdk.Coin {
	return sdk.NewInt64Coin(denom1, amount)
}

func (suite *KeeperTestSuite) Price(price string) sdk.Dec {
	return sdk.MustNewDecFromStr(price)
}

func (suite *KeeperTestSuite) Coin(denom string, amount int64) sdk.Coin {
	return sdk.NewInt64Coin(denom, amount)
}
