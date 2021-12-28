package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/fundraising/x/fundraising"
	"github.com/tendermint/fundraising/x/fundraising/types"

	_ "github.com/stretchr/testify/suite"
)

func (suite *KeeperTestSuite) TestAuctionId() {
	auctionId := suite.keeper.GetAuctionId(suite.ctx)
	suite.Require().Equal(uint64(0), auctionId)

	cacheCtx, _ := suite.ctx.CacheContext()
	nextAuctionId := suite.keeper.GetNextAuctionIdWithUpdate(cacheCtx)
	suite.Require().Equal(uint64(1), nextAuctionId)

	// set an auction
	suite.keeper.SetAuction(suite.ctx, suite.sampleFixedPriceAuctions[0])
	nextAuctionId = suite.keeper.GetNextAuctionIdWithUpdate(cacheCtx)
	suite.Require().Equal(uint64(2), nextAuctionId)

	auctions := suite.keeper.GetAuctions(suite.ctx)
	suite.Require().Len(auctions, 1)

	// set another auction
	suite.keeper.SetAuction(suite.ctx, suite.sampleFixedPriceAuctions[1])
	nextAuctionId = suite.keeper.GetNextAuctionIdWithUpdate(cacheCtx)
	suite.Require().Equal(uint64(3), nextAuctionId)

	auctions = suite.keeper.GetAuctions(suite.ctx)
	suite.Require().Len(auctions, 2)
}

func (suite *KeeperTestSuite) TestCreateAuctionStatus() {
	suite.ctx = suite.ctx.WithBlockTime(types.ParseTime("2021-12-01T00:00:00Z"))

	// create a fixed price auction with the future start time
	err := suite.keeper.CreateFixedPriceAuction(suite.ctx, types.NewMsgCreateFixedPriceAuction(
		suite.addrs[0].String(),
		sdk.MustNewDecFromStr("0.5"),
		sdk.NewInt64Coin(denom2, 100_000_000_000),
		denom1,
		[]types.VestingSchedule{},
		types.ParseTime("2022-12-10T00:00:00Z"),
		types.ParseTime("2022-12-20T00:00:00Z"),
	))
	suite.Require().NoError(err)

	auction, found := suite.keeper.GetAuction(suite.ctx, 1)
	suite.Require().True(found)
	suite.Require().Equal(types.AuctionStatusStandBy, auction.GetStatus())

	// create a fixed price auction with the past start time
	err = suite.keeper.CreateFixedPriceAuction(suite.ctx, types.NewMsgCreateFixedPriceAuction(
		suite.addrs[0].String(),
		sdk.MustNewDecFromStr("0.5"),
		sdk.NewInt64Coin(denom2, 100_000_000_000),
		denom1,
		[]types.VestingSchedule{},
		types.ParseTime("2021-11-01T00:00:00Z"),
		types.ParseTime("2021-12-10T00:00:00Z"),
	))
	suite.Require().NoError(err)

	auction, found = suite.keeper.GetAuction(suite.ctx, 2)
	suite.Require().True(found)
	suite.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())
}

func (suite *KeeperTestSuite) TestDistributePayingCoin() {
	ctx, k, auction := suite.ctx, suite.keeper, suite.sampleFixedPriceAuctions[1]

	suite.SetAuction(ctx, auction)

	auction, found := k.GetAuction(ctx, auction.GetId())
	suite.Require().True(found)

	for _, bid := range suite.sampleFixedPriceBids {
		suite.PlaceBid(ctx, bid)
	}

	suite.Require().NoError(k.DistributeSellingCoin(ctx, auction))
	suite.Require().NoError(k.SetVestingSchedules(ctx, auction))

	vqs := k.GetVestingQueuesByAuctionId(ctx, auction.GetId())
	suite.Require().Equal(4, len(vqs))

	// all of the vesting queues must not be released yet
	for _, vq := range vqs {
		suite.Require().False(vq.Released)
	}

	ctx = ctx.WithBlockTime(vqs[0].GetReleaseTime().AddDate(0, 4, 1))
	fundraising.EndBlocker(ctx, k)
	suite.Require().NoError(k.DistributePayingCoin(ctx, auction))

	// first two vesting queues must be released
	for i, vq := range k.GetVestingQueuesByAuctionId(ctx, auction.GetId()) {
		if i == 0 || i == 1 {
			suite.Require().True(vq.Released)
		} else {
			suite.Require().False(vq.Released)
		}
	}

	ctx = ctx.WithBlockTime(vqs[3].GetReleaseTime().AddDate(0, 0, 1))
	fundraising.EndBlocker(ctx, k)
	suite.Require().NoError(k.DistributePayingCoin(ctx, auction))

	// all of the vesting queues must be released
	for _, vq := range k.GetVestingQueuesByAuctionId(ctx, auction.GetId()) {
		suite.Require().True(vq.Released)
	}

	auction, found = k.GetAuction(ctx, auction.GetId())
	suite.Require().True(found)
	suite.Require().Equal(types.AuctionStatusFinished, auction.GetStatus())
}
