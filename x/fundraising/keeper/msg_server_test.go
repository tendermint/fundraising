package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/tendermint/fundraising/x/fundraising/types"

	_ "github.com/stretchr/testify/suite"
)

func (s *KeeperTestSuite) TestMsgCreateFixedPriceAuction() {
	ctx := sdk.WrapSDKContext(s.ctx)

	for _, tc := range []struct {
		name string
		msg  *types.MsgCreateFixedPriceAuction
		err  error
	}{
		{
			"valid message with the future start time",
			types.NewMsgCreateFixedPriceAuction(
				s.addrs[0].String(),
				sdk.OneDec(),
				sdk.NewInt64Coin(denom1, 1_000_000_000_000),
				denom2,
				s.sampleVestingSchedules2,
				types.MustParseRFC3339("2030-01-01T00:00:00Z"),
				types.MustParseRFC3339("2030-01-10T00:00:00Z"),
			),
			nil,
		},
	} {
		s.Run(tc.name, func() {
			_, err := s.srv.CreateFixedPriceAuction(ctx, tc.msg)
			if tc.err != nil {
				s.Require().ErrorIs(err, tc.err)
				return
			}
			s.Require().NoError(err)

			_, found := s.keeper.GetAuction(s.ctx, 1)
			s.Require().True(found)
		})
	}
}

func (s *KeeperTestSuite) TestMsgCreateEnglishAuction() {
	// TODO: not implemented yet
}

func (s *KeeperTestSuite) TestMsgCancelAuction() {
	ctx := sdk.WrapSDKContext(s.ctx)

	_, err := s.srv.CancelAuction(ctx, types.NewMsgCancelAuction(
		s.addrs[4].String(),
		uint64(1),
	))
	s.Require().ErrorIs(err, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "auction %d is not found", uint64(1)))

	// create a fixed price auction that is started status
	s.SetAuction(s.sampleFixedPriceAuctions[1])

	auction, found := s.keeper.GetAuction(s.ctx, uint64(2))
	s.Require().True(found)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	// try to cancel with an incorrect auctioneer
	_, err = s.srv.CancelAuction(ctx, types.NewMsgCancelAuction(
		s.addrs[0].String(),
		uint64(2),
	))
	s.Require().ErrorIs(err, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "failed to verify ownership of the auction"))

	// try to cancel with the auction that is already started
	_, err = s.srv.CancelAuction(ctx, types.NewMsgCancelAuction(
		auction.GetAuctioneer().String(),
		auction.GetId(),
	))
	s.Require().ErrorIs(err, sdkerrors.Wrap(types.ErrInvalidAuctionStatus, "auction cannot be canceled due to current status"))

	// create another fixed price auction that is stand by status
	s.SetAuction(s.sampleFixedPriceAuctions[0])

	auction, found = s.keeper.GetAuction(s.ctx, uint64(1))
	s.Require().True(found)
	s.Require().Equal(types.AuctionStatusStandBy, auction.GetStatus())

	// success and the selling coin must be released to the auctioneer
	_, err = s.srv.CancelAuction(ctx, types.NewMsgCancelAuction(
		auction.GetAuctioneer().String(),
		auction.GetId(),
	))
	s.Require().NoError(err)
	s.Require().True(
		s.app.BankKeeper.GetBalance(
			s.ctx, auction.GetSellingReserveAddress(),
			auction.GetSellingCoin().Denom).IsZero(),
	)
}

func (s *KeeperTestSuite) TestMsgPlaceBid() {
	ctx := sdk.WrapSDKContext(s.ctx)

	// Create a fixed price auction that should start right away
	auction := s.sampleFixedPriceAuctions[1]
	s.SetAuction(auction)

	for _, tc := range []struct {
		name string
		msg  *types.MsgPlaceBid
		err  error
	}{
		{
			"valid message",
			types.NewMsgPlaceBid(
				auction.GetId(),
				s.addrs[0].String(),
				sdk.MustNewDecFromStr("0.5"),
				sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 1_000_000),
			),
			nil,
		},
		{
			"invalid start price",
			types.NewMsgPlaceBid(
				auction.GetId(),
				s.addrs[0].String(),
				sdk.MustNewDecFromStr("1.0"),
				sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 1_000_000),
			),
			sdkerrors.Wrap(types.ErrInvalidStartPrice, "bid price must be equal to start price"),
		},
		{
			"invalid coin demo",
			types.NewMsgPlaceBid(
				auction.GetId(),
				s.addrs[0].String(),
				sdk.MustNewDecFromStr("0.5"),
				sdk.NewInt64Coin(auction.GetSellingCoin().Denom, 1_000_000),
			),
			sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "coin denom must match with the paying coin denom"),
		},
		{
			"insufficient funds",
			types.NewMsgPlaceBid(
				auction.GetId(),
				s.addrs[0].String(),
				sdk.MustNewDecFromStr("0.5"),
				sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 500_000_000_000_000_000),
			),
			sdkerrors.ErrInsufficientFunds,
		},
	} {
		s.Run(tc.name, func() {
			_, err := s.srv.PlaceBid(ctx, tc.msg)
			if tc.err != nil {
				s.Require().ErrorIs(err, tc.err)
				return
			}
			s.Require().NoError(err)
		})
	}
}
