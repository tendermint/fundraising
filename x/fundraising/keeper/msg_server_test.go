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
				suite.StartPrice("1.0"),
				suite.SellingCoin(denom1, 1_000_000_000_000),
				suite.PayingCoinDenom(denom2),
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

	auctionId := uint64(1)
	auctioneerAddr := suite.addrs[4].String()

	_, err := suite.srv.CancelAuction(ctx, types.NewMsgCancelAuction(
		auctioneerAddr,
		auctionId,
	))
	suite.Require().ErrorIs(err, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "auction %d is not found", auctionId))

	// Create a fixed price auction
	suite.keeper.SetAuction(suite.ctx, suite.sampleFixedPriceAuctions[0])

	_, found := suite.keeper.GetAuction(suite.ctx, auctionId)
	suite.Require().True(found)

	// Try to cancel with an incorrect address
	_, err = suite.srv.CancelAuction(ctx, types.NewMsgCancelAuction(
		suite.addrs[0].String(),
		auctionId,
	))
	suite.Require().ErrorIs(err, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "failed to verify ownership of the auction"))

	_, err = suite.srv.CancelAuction(ctx, types.NewMsgCancelAuction(
		auctioneerAddr,
		auctionId,
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
