package keeper_test

import (
	_ "github.com/stretchr/testify/suite"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

func (suite *KeeperTestSuite) TestAuctionId() {
	auctionId := suite.keeper.GetAuctionId(suite.ctx)
	suite.Require().Equal(uint64(0), auctionId)

	cacheCtx, _ := suite.ctx.CacheContext()
	nextAuctionId := suite.keeper.GetNextAuctionIdWithUpdate(cacheCtx)
	suite.Require().Equal(uint64(1), nextAuctionId)

	// TODO: not implemented yet
}

func (suite *KeeperTestSuite) TestCreateAuctionStatus() {
	suite.ctx = suite.ctx.WithBlockTime(types.ParseTime("2021-12-01T00:00:00Z"))

	// Create a fixed price auction with the future start time
	suite.keeper.CreateFixedPriceAuction(suite.ctx, types.NewMsgCreateFixedPriceAuction(
		suite.addrs[0].String(),
		suite.StartPrice("0.5"),
		suite.SellingCoin(denom2, 100_000_000_000),
		suite.PayingCoinDenom("denom1"),
		[]types.VestingSchedule{},
		types.ParseTime("2022-12-10T00:00:00Z"),
		types.ParseTime("2022-12-20T00:00:00Z"),
	))

	auction, found := suite.keeper.GetAuction(suite.ctx, 1)
	suite.Require().True(found)
	suite.Require().Equal(types.AuctionStatusStandBy, auction.GetStatus())

	// Create a fixed price auction with the past start time
	suite.keeper.CreateFixedPriceAuction(suite.ctx, types.NewMsgCreateFixedPriceAuction(
		suite.addrs[0].String(),
		suite.StartPrice("0.5"),
		suite.SellingCoin(denom2, 100_000_000_000),
		suite.PayingCoinDenom("denom1"),
		[]types.VestingSchedule{},
		types.ParseTime("2021-11-01T00:00:00Z"),
		types.ParseTime("2021-12-10T00:00:00Z"),
	))

	auction, found = suite.keeper.GetAuction(suite.ctx, 2)
	suite.Require().True(found)
	suite.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())
}
