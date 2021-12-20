package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	fundraisingkeeper "github.com/tendermint/fundraising/x/fundraising/keeper"

	_ "github.com/stretchr/testify/suite"
)

func (suite *KeeperTestSuite) TestSellingPoolReserveAmountInvariant() {
	k, ctx, auction := suite.keeper, suite.ctx, suite.sampleFixedPriceAuctions[1]

	k.SetAuction(suite.ctx, auction)

	_, broken := fundraisingkeeper.SellingPoolReserveAmountInvariant(k)(ctx)
	suite.Require().True(broken)

	err := k.ReserveSellingCoin(
		ctx,
		auction.GetId(),
		auction.GetAuctioneer(),
		auction.GetSellingCoin(),
	)
	suite.Require().NoError(err)

	_, broken = fundraisingkeeper.SellingPoolReserveAmountInvariant(k)(ctx)
	suite.Require().False(broken)
}

func (suite *KeeperTestSuite) TestPayingPoolReserveAmountInvariant() {
	k, ctx, auction := suite.keeper, suite.ctx, suite.sampleFixedPriceAuctions[1]

	k.SetAuction(suite.ctx, auction)
	err := k.ReserveSellingCoin(
		ctx,
		auction.GetId(),
		auction.GetAuctioneer(),
		auction.GetSellingCoin(),
	)
	suite.Require().NoError(err)

	for _, bid := range suite.sampleFixedPriceBids {
		bidderAcc, err := sdk.AccAddressFromBech32(bid.Bidder)
		suite.Require().NoError(err)
		suite.keeper.SetBid(suite.ctx, bid.AuctionId, bid.Sequence, bidderAcc, bid)

		err = suite.keeper.ReservePayingCoin(suite.ctx, bid.GetAuctionId(), bidderAcc, bid.Coin)
		suite.Require().NoError(err)
	}

	_, broken := fundraisingkeeper.PayingPoolReserveAmountInvariant(k)(ctx)
	suite.Require().False(broken)
}

func (suite *KeeperTestSuite) TestVestingPoolReserveAmountInvariant() {

}
