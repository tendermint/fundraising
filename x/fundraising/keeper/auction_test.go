package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/fundraising/x/fundraising"
	"github.com/tendermint/fundraising/x/fundraising/types"

	_ "github.com/stretchr/testify/suite"
)

func (suite *KeeperTestSuite) TestAuctionId() {
	auctionId := suite.keeper.GetLastAuctionId(suite.ctx)
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

func (suite *KeeperTestSuite) TestAuctionStatus() {
	standByAuction := suite.sampleFixedPriceAuctions[0]
	startedAuction := suite.sampleFixedPriceAuctions[1]

	suite.SetAuction(standByAuction)

	auction, found := suite.keeper.GetAuction(suite.ctx, standByAuction.GetId())
	suite.Require().True(found)
	suite.Require().Equal(types.AuctionStatusStandBy, auction.GetStatus())

	suite.SetAuction(startedAuction)

	auction, found = suite.keeper.GetAuction(suite.ctx, startedAuction.GetId())
	suite.Require().True(found)
	suite.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())
}

func (suite *KeeperTestSuite) TestDistributeSellingCoin() {
	suite.SetAuction(suite.sampleFixedPriceAuctions[1])

	auction, found := suite.keeper.GetAuction(suite.ctx, suite.sampleFixedPriceAuctions[1].GetId())
	suite.Require().True(found)

	bidderAcc1 := suite.addrs[0]
	bidderAcc2 := suite.addrs[1]
	bidderAcc3 := suite.addrs[2]

	bids := []types.Bid{
		{
			AuctionId: auction.GetId(),
			Sequence:  1,
			Bidder:    bidderAcc1.String(),
			Price:     auction.GetStartPrice(),
			Coin:      sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 100_000_000),
			Height:    uint64(suite.ctx.BlockHeight()),
			Eligible:  false,
		},
		{
			AuctionId: auction.GetId(),
			Sequence:  2,
			Bidder:    bidderAcc2.String(),
			Price:     auction.GetStartPrice(),
			Coin:      sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 200_000_000),
			Height:    uint64(suite.ctx.BlockHeight()),
			Eligible:  false,
		},
		{
			AuctionId: auction.GetId(),
			Sequence:  3,
			Bidder:    bidderAcc3.String(),
			Price:     auction.GetStartPrice(),
			Coin:      sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 300_000_000),
			Height:    uint64(suite.ctx.BlockHeight()),
			Eligible:  false,
		},
	}

	for _, bid := range bids {
		suite.PlaceBid(bid)
	}

	// selling reserve account must be empty
	err := suite.keeper.DistributeSellingCoin(suite.ctx, auction)
	suite.Require().NoError(err)
	suite.Require().Equal(
		sdk.NewCoin(auction.GetSellingCoin().Denom, sdk.ZeroInt()),
		suite.app.BankKeeper.GetBalance(suite.ctx, auction.GetSellingReserveAddress(), auction.GetSellingCoin().Denom),
	)

	bal1 := suite.app.BankKeeper.GetBalance(suite.ctx, bidderAcc1, auction.GetSellingCoin().Denom)
	bal2 := suite.app.BankKeeper.GetBalance(suite.ctx, bidderAcc2, auction.GetSellingCoin().Denom)
	bal3 := suite.app.BankKeeper.GetBalance(suite.ctx, bidderAcc3, auction.GetSellingCoin().Denom)
	suite.Require().True(bal1.Amount.GT(initialBalances.AmountOf(denom3)))
	suite.Require().True(bal2.Amount.GT(initialBalances.AmountOf(denom3)))
	suite.Require().True(bal3.Amount.GT(initialBalances.AmountOf(denom3)))
}

func (suite *KeeperTestSuite) TestDistributePayingCoin() {
	suite.SetAuction(suite.sampleFixedPriceAuctions[1])

	auction, found := suite.keeper.GetAuction(suite.ctx, suite.sampleFixedPriceAuctions[1].GetId())
	suite.Require().True(found)

	for _, bid := range suite.sampleFixedPriceBids {
		suite.PlaceBid(bid)
	}

	err := suite.keeper.DistributeSellingCoin(suite.ctx, auction)
	suite.Require().NoError(err)

	err = suite.keeper.SetVestingSchedules(suite.ctx, auction)
	suite.Require().NoError(err)

	vqs := suite.keeper.GetVestingQueuesByAuctionId(suite.ctx, auction.GetId())
	suite.Require().Equal(4, len(vqs))

	// all of the vesting queues must not be released yet
	for _, vq := range vqs {
		suite.Require().False(vq.Released)
	}

	suite.ctx = suite.ctx.WithBlockTime(vqs[0].GetReleaseTime().AddDate(0, 4, 1))
	fundraising.EndBlocker(suite.ctx, suite.keeper)
	suite.Require().NoError(suite.keeper.DistributePayingCoin(suite.ctx, auction))

	// first two vesting queues must be released
	for i, vq := range suite.keeper.GetVestingQueuesByAuctionId(suite.ctx, auction.GetId()) {
		if i == 0 || i == 1 {
			suite.Require().True(vq.Released)
		} else {
			suite.Require().False(vq.Released)
		}
	}

	suite.ctx = suite.ctx.WithBlockTime(vqs[3].GetReleaseTime().AddDate(0, 0, 1))
	fundraising.EndBlocker(suite.ctx, suite.keeper)
	suite.Require().NoError(suite.keeper.DistributePayingCoin(suite.ctx, auction))

	// all of the vesting queues must be released
	for _, vq := range suite.keeper.GetVestingQueuesByAuctionId(suite.ctx, auction.GetId()) {
		suite.Require().True(vq.Released)
	}

	auction, found = suite.keeper.GetAuction(suite.ctx, auction.GetId())
	suite.Require().True(found)
	suite.Require().Equal(types.AuctionStatusFinished, auction.GetStatus())
}

func (suite *KeeperTestSuite) TestCancelAuction() {
	standByAuction := suite.sampleFixedPriceAuctions[0]

	suite.SetAuction(standByAuction)

	auction, found := suite.keeper.GetAuction(suite.ctx, standByAuction.GetId())
	suite.Require().True(found)
	suite.Require().Equal(types.AuctionStatusStandBy, auction.GetStatus())

	suite.CancelAuction(auction)

	auction, found = suite.keeper.GetAuction(suite.ctx, standByAuction.GetId())
	suite.Require().True(found)
	suite.Require().Equal(types.AuctionStatusCancelled, auction.GetStatus())

	sellingReserve := suite.app.BankKeeper.GetBalance(suite.ctx, auction.GetSellingReserveAddress(), auction.GetSellingCoin().Denom)
	suite.Require().Equal(sdk.NewCoin(auction.GetSellingCoin().Denom, sdk.ZeroInt()), sellingReserve)
}
