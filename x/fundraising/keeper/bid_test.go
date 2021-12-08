package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	_ "github.com/stretchr/testify/suite"
)

func (suite *KeeperTestSuite) TestBidIterators() {
	ctx := sdk.WrapSDKContext(suite.ctx)

	// Create a fixed price auction with already started status
	suite.keeper.SetAuction(suite.ctx, suite.sampleFixedPriceAuctions[0])

	auction, found := suite.keeper.GetAuction(suite.ctx, 1)
	suite.Require().True(found)

	_, err := suite.srv.PlaceBid(ctx, suite.sampleFixedPriceBids[0])
	suite.Require().NoError(err)

	_, err = suite.srv.PlaceBid(ctx, suite.sampleFixedPriceBids[1])
	suite.Require().NoError(err)

	bids := suite.keeper.GetBids(suite.ctx)
	suite.Require().Len(bids, 2)

	bidsById := suite.keeper.GetBidsByAuctionId(suite.ctx, auction.GetId())
	suite.Require().Len(bidsById, 2)

	bidsByBidder := suite.keeper.GetBidsByBidder(suite.ctx, suite.addrs[0])
	suite.Require().Len(bidsByBidder, 1)
}
