package fundraising_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/fundraising/x/fundraising"
	"github.com/tendermint/fundraising/x/fundraising/types"

	_ "github.com/stretchr/testify/suite"
)

func (suite *ModuleTestSuite) TestEndBlockerStandByStatus() {
	err := suite.keeper.CreateFixedPriceAuction(suite.ctx, suite.sampleFixedPriceAuctions[0])
	suite.Require().NoError(err)

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
	err := suite.keeper.CreateFixedPriceAuction(suite.ctx, suite.sampleFixedPriceAuctions[1])
	suite.Require().NoError(err)

	auction, found := suite.keeper.GetAuction(suite.ctx, 1)
	suite.Require().True(found)
	suite.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	totalBidCoin := sdk.NewInt64Coin(denom4, 0)
	for _, bid := range suite.sampleFixedPriceBids {
		err := suite.keeper.PlaceBid(suite.ctx, bid)
		suite.Require().NoError(err)

		totalBidCoin = totalBidCoin.Add(bid.Coin)
	}

	bidAmt := totalBidCoin.Amount.ToDec().Quo(auction.GetStartPrice()).TruncateInt()
	receiveCoin := sdk.NewCoin(auction.GetSellingCoin().Denom, bidAmt)

	// bids with 30_000_000denom2 and 50_000_000denom2
	payingReserve := suite.app.BankKeeper.GetBalance(
		suite.ctx,
		types.PayingReserveAcc(auction.GetId()),
		auction.GetPayingCoinDenom(),
	)
	suite.Require().Equal(totalBidCoin, payingReserve)

	suite.ctx = suite.ctx.WithBlockTime(auction.GetEndTimes()[0].AddDate(0, 0, -1))
	fundraising.EndBlocker(suite.ctx, suite.keeper)

	// remaining selling coin should have returned to the auctioneer
	auctioneerBalance := suite.app.BankKeeper.GetBalance(
		suite.ctx,
		suite.addrs[5],
		auction.GetSellingCoin().Denom,
	)
	fmt.Println("auction.GetSellingCoin(): ", auction.GetSellingCoin())
	fmt.Println("auctioneerBalance: ", auctioneerBalance)
	suite.Require().Equal(auction.GetSellingCoin(), auctioneerBalance.Add(receiveCoin))
}

func (suite *ModuleTestSuite) TestEndBlockerVestingStatus() {
	err := suite.keeper.CreateFixedPriceAuction(suite.ctx, suite.sampleFixedPriceAuctions[1])
	suite.Require().NoError(err)

	auction, found := suite.keeper.GetAuction(suite.ctx, 1)
	suite.Require().True(found)
	suite.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	totalBidCoin := sdk.NewInt64Coin(denom4, 0)
	for _, bid := range suite.sampleFixedPriceBids {
		totalBidCoin = totalBidCoin.Add(bid.Coin)

		err := suite.keeper.PlaceBid(suite.ctx, bid)
		suite.Require().NoError(err)
	}

	suite.ctx = suite.ctx.WithBlockTime(auction.GetEndTimes()[0].AddDate(0, 0, -1))
	fundraising.EndBlocker(suite.ctx, suite.keeper)

	vestingReserve := suite.app.BankKeeper.GetBalance(
		suite.ctx,
		types.VestingReserveAcc(auction.GetId()),
		auction.GetPayingCoinDenom(),
	)
	suite.Require().Equal(totalBidCoin, vestingReserve)

	suite.ctx = suite.ctx.WithBlockTime(types.ParseTime("2022-04-02T00:00:00Z"))
	fundraising.EndBlocker(suite.ctx, suite.keeper)

	queues := suite.keeper.GetVestingQueuesByAuctionId(suite.ctx, auction.GetId())
	suite.Require().Len(queues, 4)
	suite.Require().True(queues[0].Vested)
	suite.Require().True(queues[1].Vested)

	// Auctioneer should have received two vested amounts
	auctioneerBalance := suite.app.BankKeeper.GetBalance(
		suite.ctx,
		suite.addrs[5],
		auction.GetPayingCoinDenom(),
	)
	suite.Require().Equal(totalBidCoin.Amount, auctioneerBalance.Amount.Sub(initialBalances.AmountOf(denom4)))
}
