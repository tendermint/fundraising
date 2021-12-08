package keeper_test

import (
	_ "github.com/stretchr/testify/suite"
)

func (suite *KeeperTestSuite) TestAuctionId() {
	auctionId := suite.keeper.GetAuctionId(suite.ctx)
	suite.Require().Equal(uint64(0), auctionId)

	cacheCtx, _ := suite.ctx.CacheContext()
	nextAuctionId := suite.keeper.GetNextAuctionIdWithUpdate(cacheCtx)
	suite.Require().Equal(uint64(1), nextAuctionId)

	// TODO: not implemented yet
}
