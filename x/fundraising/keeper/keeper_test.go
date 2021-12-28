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

type KeeperTestSuite struct {
	suite.Suite

	app                      *app.App
	ctx                      sdk.Context
	keeper                   keeper.Keeper
	querier                  keeper.Querier
	srv                      types.MsgServer
	addrs                    []sdk.AccAddress
	sampleVestingSchedules1  []types.VestingSchedule
	sampleVestingSchedules2  []types.VestingSchedule
	sampleFixedPriceAuctions []types.AuctionI
	sampleFixedPriceBids     []types.Bid
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
	suite.sampleVestingSchedules1 = []types.VestingSchedule{
		{
			ReleaseTime: types.ParseTime("2030-01-31T22:00:00+00:00"),
			Weight:      sdk.MustNewDecFromStr("0.5"),
		},
		{
			ReleaseTime: types.ParseTime("2030-12-01T22:00:00+00:00"),
			Weight:      sdk.MustNewDecFromStr("0.5"),
		},
	}
	suite.sampleVestingSchedules2 = []types.VestingSchedule{
		{
			ReleaseTime: types.ParseTime("2022-01-01T22:00:00+00:00"),
			Weight:      sdk.MustNewDecFromStr("0.25"),
		},
		{
			ReleaseTime: types.ParseTime("2022-04-01T22:00:00+00:00"),
			Weight:      sdk.MustNewDecFromStr("0.25"),
		},
		{
			ReleaseTime: types.ParseTime("2022-08-01T22:00:00+00:00"),
			Weight:      sdk.MustNewDecFromStr("0.25"),
		},
		{
			ReleaseTime: types.ParseTime("2022-12-01T22:00:00+00:00"),
			Weight:      sdk.MustNewDecFromStr("0.25"),
		},
	}
	suite.sampleFixedPriceAuctions = []types.AuctionI{
		types.NewFixedPriceAuction(
			&types.BaseAuction{
				Id:                    1,
				Type:                  types.AuctionTypeFixedPrice,
				Auctioneer:            suite.addrs[4].String(),
				SellingReserveAddress: types.SellingReserveAcc(1).String(),
				PayingReserveAddress:  types.PayingReserveAcc(1).String(),
				StartPrice:            sdk.OneDec(), // start price corresponds to the ratio of the paying coin
				SellingCoin:           sdk.NewInt64Coin(denom1, 1_000_000_000_000),
				PayingCoinDenom:       denom2,
				VestingReserveAddress: types.VestingReserveAcc(1).String(),
				VestingSchedules:      suite.sampleVestingSchedules1,
				WinningPrice:          sdk.ZeroDec(),
				RemainingCoin:         sdk.NewInt64Coin(denom1, 1_000_000_000_000),
				StartTime:             types.ParseTime("2022-01-01T00:00:00Z"),
				EndTimes:              []time.Time{types.ParseTime("2022-01-10T00:00:00Z")},
				Status:                types.AuctionStatusStandBy,
			},
		),
		types.NewFixedPriceAuction(
			&types.BaseAuction{
				Id:                    2,
				Type:                  types.AuctionTypeFixedPrice,
				Auctioneer:            suite.addrs[5].String(),
				SellingReserveAddress: types.SellingReserveAcc(2).String(),
				PayingReserveAddress:  types.PayingReserveAcc(2).String(),
				StartPrice:            sdk.MustNewDecFromStr("0.5"),
				SellingCoin:           sdk.NewInt64Coin(denom3, 1_000_000_000_000),
				PayingCoinDenom:       denom4,
				VestingReserveAddress: types.VestingReserveAcc(2).String(),
				VestingSchedules:      suite.sampleVestingSchedules2,
				WinningPrice:          sdk.ZeroDec(),
				RemainingCoin:         sdk.NewInt64Coin(denom3, 1_000_000_000_000),
				StartTime:             types.ParseTime("2021-12-10T00:00:00Z"),
				EndTimes:              []time.Time{types.ParseTime("2021-12-24T00:00:00Z")},
				Status:                types.AuctionStatusStarted,
			},
		),
	}
	suite.sampleFixedPriceBids = []types.Bid{
		{
			AuctionId: 2,
			Sequence:  1,
			Bidder:    suite.addrs[0].String(),
			Price:     sdk.MustNewDecFromStr("0.5"),
			Coin:      sdk.NewInt64Coin(denom4, 20_000_000),
			Height:    uint64(suite.ctx.BlockHeight()),
			Eligible:  false,
		},
		{
			AuctionId: 2,
			Sequence:  2,
			Bidder:    suite.addrs[0].String(),
			Price:     sdk.MustNewDecFromStr("0.5"),
			Coin:      sdk.NewInt64Coin(denom4, 30_000_000),
			Height:    uint64(suite.ctx.BlockHeight()),
			Eligible:  false,
		},
		{
			AuctionId: 2,
			Sequence:  3,
			Bidder:    suite.addrs[1].String(),
			Price:     sdk.MustNewDecFromStr("0.5"),
			Coin:      sdk.NewInt64Coin(denom4, 50_000_000),
			Height:    uint64(suite.ctx.BlockHeight()),
			Eligible:  false,
		},
		{
			AuctionId: 2,
			Sequence:  4,
			Bidder:    suite.addrs[1].String(),
			Price:     sdk.MustNewDecFromStr("0.5"),
			Coin:      sdk.NewInt64Coin(denom4, 50_000_000),
			Height:    uint64(suite.ctx.BlockHeight()),
			Eligible:  true,
		},
	}
}

// SetAuction is a convenient method to set an auction and reserve selling coin to the selling reserve account.
func (suite *KeeperTestSuite) SetAuction(ctx sdk.Context, auction types.AuctionI) {
	suite.keeper.SetAuction(suite.ctx, auction)
	err := suite.keeper.ReserveSellingCoin(
		ctx,
		auction.GetId(),
		auction.GetAuctioneer(),
		auction.GetSellingCoin(),
	)
	suite.Require().NoError(err)
}

// PlaceBid is a convenient method to bid and reserve paying coin to the paying reserve account.
func (suite *KeeperTestSuite) PlaceBid(ctx sdk.Context, bid types.Bid) {
	bidderAcc, err := sdk.AccAddressFromBech32(bid.Bidder)
	suite.Require().NoError(err)

	suite.keeper.SetBid(suite.ctx, bid.AuctionId, bid.Sequence, bidderAcc, bid)

	err = suite.keeper.ReservePayingCoin(
		suite.ctx,
		bid.GetAuctionId(),
		bidderAcc,
		bid.Coin,
	)
	suite.Require().NoError(err)
}

// PlaceBidWithCustom is a convenient method to bid with custom fields and
// reserve paying coin to the paying reserve account.
func (suite *KeeperTestSuite) PlaceBidWithCustom(
	ctx sdk.Context,
	auctionId uint64,
	sequence uint64,
	bidder string,
	price sdk.Dec,
	coin sdk.Coin,
) {
	bidderAcc, err := sdk.AccAddressFromBech32(bidder)
	suite.Require().NoError(err)

	suite.keeper.SetBid(suite.ctx, auctionId, sequence, bidderAcc, types.Bid{
		AuctionId: auctionId,
		Sequence:  sequence,
		Bidder:    bidderAcc.String(),
		Price:     price,
		Coin:      coin,
	})

	err = suite.keeper.ReservePayingCoin(
		suite.ctx,
		auctionId,
		bidderAcc,
		coin,
	)
	suite.Require().NoError(err)
}

// coinEq is a convenient method to test expected and got values of sdk.Coin.
func coinEq(exp, got sdk.Coin) (bool, string, string, string) {
	return exp.IsEqual(got), "expected:\t%v\ngot:\t\t%v", exp.String(), got.String()
}
