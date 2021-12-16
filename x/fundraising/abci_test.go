package fundraising_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/fundraising/x/fundraising"
	"github.com/tendermint/fundraising/x/fundraising/types"

	_ "github.com/stretchr/testify/suite"
)

func (suite *ModuleTestSuite) TestEndBlockerStandByStatus() {
	suite.keeper.SetAuction(suite.ctx, suite.sampleFixedPriceAuctions[0])

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
	suite.keeper.SetAuction(suite.ctx, suite.sampleFixedPriceAuctions[1])
	err := suite.keeper.ReserveSellingCoin(
		suite.ctx,
		suite.sampleFixedPriceAuctions[1].GetId(),
		suite.sampleFixedPriceAuctions[1].GetAuctioneer(),
		suite.sampleFixedPriceAuctions[1].GetSellingCoin(),
	)
	suite.Require().NoError(err)

	auction, found := suite.keeper.GetAuction(suite.ctx, 2)
	suite.Require().True(found)
	suite.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	totalBidCoin := sdk.NewInt64Coin(suite.sampleFixedPriceAuctions[1].GetPayingCoinDenom(), 0)
	for _, bid := range suite.sampleFixedPriceBids {
		bidderAcc, err := sdk.AccAddressFromBech32(bid.Bidder)
		suite.Require().NoError(err)
		suite.keeper.SetBid(suite.ctx, bid.AuctionId, bid.Sequence, bidderAcc, bid)

		err = suite.keeper.ReservePayingCoin(suite.ctx, auction.GetId(), bidderAcc, bid.Coin)
		suite.Require().NoError(err)

		totalBidCoin = totalBidCoin.Add(bid.Coin)
	}

	receiveCoin := sdk.NewCoin(
		auction.GetSellingCoin().Denom,
		totalBidCoin.Amount.ToDec().Quo(auction.GetStartPrice()).TruncateInt(),
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
	suite.keeper.SetAuction(suite.ctx, suite.sampleFixedPriceAuctions[1])
	err := suite.keeper.ReserveSellingCoin(
		suite.ctx,
		suite.sampleFixedPriceAuctions[1].GetId(),
		suite.sampleFixedPriceAuctions[1].GetAuctioneer(),
		suite.sampleFixedPriceAuctions[1].GetSellingCoin(),
	)
	suite.Require().NoError(err)

	auction, found := suite.keeper.GetAuction(suite.ctx, 2)
	suite.Require().True(found)
	suite.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	totalBidCoin := sdk.NewInt64Coin(suite.sampleFixedPriceAuctions[1].GetPayingCoinDenom(), 0)
	for _, bid := range suite.sampleFixedPriceBids {
		bidderAcc, err := sdk.AccAddressFromBech32(bid.Bidder)
		suite.Require().NoError(err)
		suite.keeper.SetBid(suite.ctx, bid.AuctionId, bid.Sequence, bidderAcc, bid)

		err = suite.keeper.ReservePayingCoin(suite.ctx, auction.GetId(), bidderAcc, bid.Coin)
		suite.Require().NoError(err)

		totalBidCoin = totalBidCoin.Add(bid.Coin)
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

	// auctioneer should have received two vested amounts
	auctioneerBalance := suite.app.BankKeeper.GetBalance(
		suite.ctx,
		suite.addrs[5],
		auction.GetPayingCoinDenom(),
	)
	suite.Require().Equal(totalBidCoin.Amount, auctioneerBalance.Amount.Sub(initialBalances.AmountOf(denom4)))
}
