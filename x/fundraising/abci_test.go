package fundraising_test

import (
	"github.com/tendermint/fundraising/x/fundraising"
	"github.com/tendermint/fundraising/x/fundraising/types"

	_ "github.com/stretchr/testify/suite"
)

func (suite *ModuleTestSuite) TestEndBlockerStandByStatus() {
	suite.keeper.CreateFixedPriceAuction(suite.ctx, suite.sampleFixedPriceAuctions[0])

	auction, found := suite.keeper.GetAuction(suite.ctx, 1)
	suite.Require().True(found)
	suite.Require().Equal(types.AuctionStatusStandBy, auction.GetStatus())

	// Modify start time and block time
	t := types.ParseTime("2021-12-27T00:00:01Z")
	_ = auction.SetStartTime(t)
	suite.keeper.SetAuction(suite.ctx, auction)
	suite.ctx = suite.ctx.WithBlockTime(t)
	fundraising.EndBlocker(suite.ctx, suite.keeper)

	auction, found = suite.keeper.GetAuction(suite.ctx, 1)
	suite.Require().True(found)
	suite.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())
}

func (suite *ModuleTestSuite) TestEndBlockerStartedStatus() {
	suite.keeper.CreateFixedPriceAuction(suite.ctx, suite.sampleFixedPriceAuctions[1])

	auction, found := suite.keeper.GetAuction(suite.ctx, 1)
	suite.Require().True(found)
	suite.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())
	suite.Require().Equal(types.ParseTime("2021-12-20T00:00:00Z"), auction.GetEndTimes()[0])

	sellingReserveAcc := types.SellingReserveAcc(auction.GetId())
	sellingPool := suite.app.BankKeeper.GetBalance(suite.ctx, sellingReserveAcc, auction.GetSellingCoin().Denom)
	suite.Require().Equal(auction.GetSellingCoin(), sellingPool)

	// for _, bid := range suite.sampleFixedPriceBids {
	// 	err := suite.keeper.PlaceBid(suite.ctx, bid)
	// 	suite.Require().NoError(err)
	// }

	// payingReserveAcc := types.PayingReserveAcc(auction.GetId())
	// payingPool := suite.app.BankKeeper.GetBalance(suite.ctx, payingReserveAcc, auction.GetSellingCoin().Denom)
	// fmt.Println("paying: ", payingPool)

	// t := types.ParseTime("2021-12-20T00:00:00Z")
	// suite.ctx = suite.ctx.WithBlockTime(t)
	// fundraising.EndBlocker(suite.ctx, suite.keeper)

	// sellingPool = suite.app.BankKeeper.GetBalance(suite.ctx, sellingReserveAcc, auction.GetSellingCoin().Denom)
	// fmt.Println("sellingPool: ", sellingPool)

	// balances := suite.app.BankKeeper.GetAllBalances(suite.ctx, suite.addrs[0])
	// fmt.Println("balances: ", balances)
}

func (suite *ModuleTestSuite) TestEndBlockerVestingStatus() {
	// TODO: not implemented yet
}
