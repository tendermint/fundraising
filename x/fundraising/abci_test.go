package fundraising_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/fundraising/x/fundraising"
	"github.com/tendermint/fundraising/x/fundraising/types"

	_ "github.com/stretchr/testify/suite"
)

func (suite *ModuleTestSuite) TestEndBlockerStandByStatus() {
	suite.SetAuction(suite.ctx, suite.sampleFixedPriceAuctions[0])

	auction, found := suite.keeper.GetAuction(suite.ctx, 1)
	suite.Require().True(found)
	suite.Require().Equal(types.AuctionStatusStandBy, auction.GetStatus())

	t := types.ParseTime("2021-12-27T00:00:01Z")

	// modify start time and block time
	_ = auction.SetStartTime(t)
	suite.keeper.SetAuction(suite.ctx, auction)
	suite.ctx = suite.ctx.WithBlockTime(t)
	fundraising.EndBlocker(suite.ctx, suite.keeper)

	auction, found = suite.keeper.GetAuction(suite.ctx, 1)
	suite.Require().True(found)
	suite.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())
}

func (suite *ModuleTestSuite) TestEndBlockerStartedStatus() {
	suite.SetAuction(suite.ctx, suite.sampleFixedPriceAuctions[1])

	auction, found := suite.keeper.GetAuction(suite.ctx, 2)
	suite.Require().True(found)
	suite.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	totalBidCoin := sdk.NewInt64Coin(suite.sampleFixedPriceAuctions[1].GetPayingCoinDenom(), 0)
	for _, bid := range suite.sampleFixedPriceBids {
		suite.PlaceBid(suite.ctx, bid)

		totalBidCoin = totalBidCoin.Add(bid.Coin)
	}

	receiveAmt := totalBidCoin.Amount.ToDec().Quo(auction.GetStartPrice()).TruncateInt()
	receiveCoin := sdk.NewCoin(
		auction.GetSellingCoin().Denom,
		receiveAmt,
	)

	// total bid amounts must be 150_000_000denom4
	payingReserve := suite.app.BankKeeper.GetBalance(
		suite.ctx,
		types.PayingReserveAcc(auction.GetId()),
		auction.GetPayingCoinDenom(),
	)
	suite.Require().True(coinEq(totalBidCoin, payingReserve))

	suite.ctx = suite.ctx.WithBlockTime(auction.GetEndTimes()[0].AddDate(0, 0, -1))
	fundraising.EndBlocker(suite.ctx, suite.keeper)

	// remaining selling coin should have returned to the auctioneer
	auctioneerBalance := suite.app.BankKeeper.GetBalance(
		suite.ctx,
		suite.addrs[5],
		auction.GetSellingCoin().Denom,
	)
	suite.Require().Equal(auction.GetSellingCoin(), auctioneerBalance.Add(receiveCoin))
}

func (suite *ModuleTestSuite) TestEndBlockerVestingStatus() {
	suite.SetAuction(suite.ctx, suite.sampleFixedPriceAuctions[1])

	auction, found := suite.keeper.GetAuction(suite.ctx, 2)
	suite.Require().True(found)
	suite.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	totalBidCoin := sdk.NewInt64Coin(suite.sampleFixedPriceAuctions[1].GetPayingCoinDenom(), 0)
	for _, bid := range suite.sampleFixedPriceBids {
		suite.PlaceBid(suite.ctx, bid)

		totalBidCoin = totalBidCoin.Add(bid.Coin)
	}

	// set the current block time a day before so that it gets finished
	suite.ctx = suite.ctx.WithBlockTime(auction.GetEndTimes()[0].AddDate(0, 0, -1))
	fundraising.EndBlocker(suite.ctx, suite.keeper)

	vestingReserve := suite.app.BankKeeper.GetBalance(
		suite.ctx,
		auction.GetVestingReserveAddress(),
		auction.GetPayingCoinDenom(),
	)
	suite.Require().Equal(totalBidCoin, vestingReserve)

	suite.ctx = suite.ctx.WithBlockTime(types.ParseTime("2022-04-02T00:00:00Z"))
	fundraising.EndBlocker(suite.ctx, suite.keeper)

	queues := suite.keeper.GetVestingQueuesByAuctionId(suite.ctx, auction.GetId())
	suite.Require().Len(queues, 4)
	suite.Require().True(queues[0].Released)
	suite.Require().True(queues[1].Released)
	suite.Require().False(queues[2].Released)
	suite.Require().False(queues[3].Released)

	// auctioneer should have received two released amounts
	auctioneerBalance := suite.app.BankKeeper.GetBalance(
		suite.ctx,
		suite.addrs[5],
		auction.GetPayingCoinDenom(),
	)
	suite.Require().Equal(
		totalBidCoin.Amount.Quo(sdk.NewInt(2)),
		auctioneerBalance.Amount.Sub(initialBalances.AmountOf(denom4)),
	)
}
