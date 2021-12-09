package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/tendermint/fundraising/x/fundraising/types"

	_ "github.com/stretchr/testify/suite"
)

func (suite *KeeperTestSuite) TestMsgCreateFixedPriceAuction() {
	ctx := sdk.WrapSDKContext(suite.ctx)

	for _, tc := range []struct {
		name string
		msg  *types.MsgCreateFixedPriceAuction
		err  error
	}{
		{
			"valid message with the future start time",
			types.NewMsgCreateFixedPriceAuction(
				suite.addrs[0].String(),
				sdk.OneDec(),
				sdk.NewInt64Coin(denom1, 1_000_000_000_000),
				denom2,
				suite.VestingSchedules(),
				types.ParseTime("2030-01-01T00:00:00Z"),
				types.ParseTime("2030-01-10T00:00:00Z"),
			),
			nil,
		},
	} {
		suite.Run(tc.name, func() {
			_, err := suite.srv.CreateFixedPriceAuction(ctx, tc.msg)
			if tc.err != nil {
				suite.Require().ErrorIs(err, tc.err)
				return
			}
			suite.Require().NoError(err)

			_, found := suite.keeper.GetAuction(suite.ctx, 1)
			suite.Require().True(found)
		})
	}
}

func (suite *KeeperTestSuite) TestMsgCreateEnglishAuction() {
	// TODO: not implemented yet
}

func (suite *KeeperTestSuite) TestMsgCancelAuction() {
	ctx := sdk.WrapSDKContext(suite.ctx)

	_, err := suite.srv.CancelAuction(ctx, types.NewMsgCancelAuction(
		suite.addrs[4].String(),
		uint64(1),
	))
	suite.Require().ErrorIs(err, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "auction %d is not found", uint64(1)))

	// Create a fixed price auction that is started status
	suite.keeper.SetAuction(suite.ctx, suite.sampleFixedPriceAuctions[0])

	auction, found := suite.keeper.GetAuction(suite.ctx, uint64(1))
	suite.Require().True(found)
	suite.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	// Try to cancel with an incorrect auctioneer
	_, err = suite.srv.CancelAuction(ctx, types.NewMsgCancelAuction(
		suite.addrs[0].String(),
		uint64(1),
	))
	suite.Require().ErrorIs(err, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "failed to verify ownership of the auction"))

	// Try to cancel with the auction that is already started
	_, err = suite.srv.CancelAuction(ctx, types.NewMsgCancelAuction(
		auction.GetAuctioneer(),
		auction.GetId(),
	))
	suite.Require().ErrorIs(err, sdkerrors.Wrap(types.ErrInvalidAuctionStatus, "auction cannot be canceled due to current status"))

	// Create another fixed price auction that is stand by status
	suite.keeper.SetAuction(suite.ctx, suite.sampleFixedPriceAuctions[1])

	auction, found = suite.keeper.GetAuction(suite.ctx, uint64(2))
	suite.Require().True(found)
	suite.Require().Equal(types.AuctionStatusStandBy, auction.GetStatus())

	// Success
	_, err = suite.srv.CancelAuction(ctx, types.NewMsgCancelAuction(
		auction.GetAuctioneer(),
		auction.GetId(),
	))
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) TestMsgPlaceBid() {
	ctx := sdk.WrapSDKContext(suite.ctx)

	// Create a fixed price auction that should start right away
	auction := suite.sampleFixedPriceAuctions[0]
	suite.keeper.SetAuction(suite.ctx, auction)

	for _, tc := range []struct {
		name string
		msg  *types.MsgPlaceBid
		err  error
	}{
		{
			"valid message",
			types.NewMsgPlaceBid(
				auction.GetId(),
				suite.addrs[0].String(),
				sdk.MustNewDecFromStr("1.0"),
				sdk.NewInt64Coin(auction.GetSellingCoin().Denom, 1_000_000),
			),
			nil,
		},
		{
			"invalid start price",
			types.NewMsgPlaceBid(
				auction.GetId(),
				suite.addrs[0].String(),
				sdk.MustNewDecFromStr("0.5"),
				sdk.NewInt64Coin(auction.GetSellingCoin().Denom, 1_000_000),
			),
			sdkerrors.Wrap(types.ErrInvalidStartPrice, "bid price must be equal to start price"),
		},
		{
			"insufficient funds",
			types.NewMsgPlaceBid(
				auction.GetId(),
				suite.addrs[0].String(),
				sdk.MustNewDecFromStr("1.0"),
				sdk.NewInt64Coin(auction.GetSellingCoin().Denom, 500_000_000_000_000_000),
			),
			sdkerrors.ErrInsufficientFunds,
		},
	} {
		suite.Run(tc.name, func() {
			_, err := suite.srv.PlaceBid(ctx, tc.msg)
			if tc.err != nil {
				suite.Require().ErrorIs(err, tc.err)
				return
			}
			suite.Require().NoError(err)

			a, found := suite.keeper.GetAuction(suite.ctx, uint64(1))
			suite.Require().True(found)

			// Verify the total selling coin
			original := auction.GetRemainingCoin()
			remaining := a.GetRemainingCoin()
			suite.Require().Equal(original.Sub(remaining), tc.msg.Coin)
		})
	}
}
