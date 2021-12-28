package keeper_test

import (
	_ "github.com/stretchr/testify/suite"
)

func (suite *KeeperTestSuite) TestBidIterators() {
	suite.keeper.SetAuction(suite.ctx, suite.sampleFixedPriceAuctions[1])

	auction, found := suite.keeper.GetAuction(suite.ctx, suite.sampleFixedPriceAuctions[1].GetId())
	suite.Require().True(found)

	for _, bid := range suite.sampleFixedPriceBids {
		suite.PlaceBid(bid)
	}

	bids := suite.keeper.GetBids(suite.ctx)
	suite.Require().Len(bids, 4)

	bidsById := suite.keeper.GetBidsByAuctionId(suite.ctx, auction.GetId())
	suite.Require().Len(bidsById, 4)

	bidsByBidder := suite.keeper.GetBidsByBidder(suite.ctx, suite.addrs[0])
	suite.Require().Len(bidsByBidder, 2)
}
