package keeper_test

import (
	_ "github.com/stretchr/testify/suite"
)

func (suite *KeeperTestSuite) TestAuctionId() {
	auctionID := suite.keeper.GetAuctionId(suite.ctx)
	suite.Require().Equal(uint64(0), auctionID)

	cacheCtx, _ := suite.ctx.CacheContext()
	nextAuctionID := suite.keeper.GetNextAuctionIdWithUpdate(cacheCtx)
	suite.Require().Equal(uint64(1), nextAuctionID)

	// TODO: not implemented yet
}
