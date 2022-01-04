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
				suite.sampleVestingSchedules2,
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
	ctx := sdk.WrapSDKContext(suite.ctx)

	for _, tc := range []struct {
		name string
		msg  *types.MsgCreateEnglishAuction
		err  error
	}{
		{
			"valid message with the future start time",
			types.NewMsgCreateEnglishAuction(
				suite.addrs[0].String(),
				sdk.OneDec(),
				sdk.NewInt64Coin(denom1, 1_000_000_000_000),
				denom2,
				suite.sampleVestingSchedules2,
				sdk.MustNewDecFromStr("1.0"),
				sdk.MustNewDecFromStr("0.5"),
				types.ParseTime("2030-01-01T00:00:00Z"),
				types.ParseTime("2030-01-10T00:00:00Z"),
			),
			nil,
		},
	} {
		suite.Run(tc.name, func() {
			_, err := suite.srv.CreateEnglishAuction(ctx, tc.msg)
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

func (suite *KeeperTestSuite) TestMsgCancelAuction() {
	ctx := sdk.WrapSDKContext(suite.ctx)

	_, err := suite.srv.CancelAuction(ctx, types.NewMsgCancelAuction(
		suite.addrs[4].String(),
		uint64(1),
	))
	suite.Require().ErrorIs(err, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "auction %d is not found", uint64(1)))

	// create a fixed price auction that is started status
	suite.SetAuction(suite.sampleFixedPriceAuctions[1])

	auction, found := suite.keeper.GetAuction(suite.ctx, uint64(2))
	suite.Require().True(found)
	suite.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	// try to cancel with an incorrect auctioneer
	_, err = suite.srv.CancelAuction(ctx, types.NewMsgCancelAuction(
		suite.addrs[0].String(),
		uint64(2),
	))
	suite.Require().ErrorIs(err, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "failed to verify ownership of the auction"))

	// try to cancel with the auction that is already started
	_, err = suite.srv.CancelAuction(ctx, types.NewMsgCancelAuction(
		auction.GetAuctioneer().String(),
		auction.GetId(),
	))
	suite.Require().ErrorIs(err, sdkerrors.Wrap(types.ErrInvalidAuctionStatus, "auction cannot be canceled due to current status"))

	// create another fixed price auction that is stand by status
	suite.SetAuction(suite.sampleFixedPriceAuctions[0])

	auction, found = suite.keeper.GetAuction(suite.ctx, uint64(1))
	suite.Require().True(found)
	suite.Require().Equal(types.AuctionStatusStandBy, auction.GetStatus())

	// success and the selling coin must be released to the auctioneer
	_, err = suite.srv.CancelAuction(ctx, types.NewMsgCancelAuction(
		auction.GetAuctioneer().String(),
		auction.GetId(),
	))
	suite.Require().NoError(err)
	suite.Require().True(
		suite.app.BankKeeper.GetBalance(
			suite.ctx, auction.GetSellingReserveAddress(),
			auction.GetSellingCoin().Denom).IsZero(),
	)
}

func (suite *KeeperTestSuite) TestMsgPlaceBid() {
	ctx := sdk.WrapSDKContext(suite.ctx)

	// Create a fixed price auction that should start right away
	auction := suite.sampleFixedPriceAuctions[1]
	suite.SetAuction(auction)

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
				sdk.MustNewDecFromStr("0.5"),
				sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 1_000_000),
			),
			nil,
		},
		{
			"invalid start price",
			types.NewMsgPlaceBid(
				auction.GetId(),
				suite.addrs[0].String(),
				sdk.MustNewDecFromStr("1.0"),
				sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 1_000_000),
			),
			sdkerrors.Wrap(types.ErrInvalidStartPrice, "bid price must be equal to start price"),
		},
		{
			"invalid coin demo",
			types.NewMsgPlaceBid(
				auction.GetId(),
				suite.addrs[0].String(),
				sdk.MustNewDecFromStr("0.5"),
				sdk.NewInt64Coin(auction.GetSellingCoin().Denom, 1_000_000),
			),
			sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "coin denom must match with the paying coin denom"),
		},
		{
			"insufficient funds",
			types.NewMsgPlaceBid(
				auction.GetId(),
				suite.addrs[0].String(),
				sdk.MustNewDecFromStr("0.5"),
				sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 500_000_000_000_000_000),
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
		})
	}
}
