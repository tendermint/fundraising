package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/tendermint/fundraising/x/fundraising/keeper"
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
				s.addr(0).String(),
				sdk.OneDec(),
				sdk.NewInt64Coin("denom1", 1_000_000_000_000),
				"denom2",
				[]types.VestingSchedule{
					{
						ReleaseTime: types.MustParseRFC3339("2023-01-01T00:00:00Z"),
						Weight:      sdk.MustNewDecFromStr("0.25"),
					},
					{
						ReleaseTime: types.MustParseRFC3339("2023-05-01T00:00:00Z"),
						Weight:      sdk.MustNewDecFromStr("0.25"),
					},
					{
						ReleaseTime: types.MustParseRFC3339("2023-09-01T00:00:00Z"),
						Weight:      sdk.MustNewDecFromStr("0.25"),
					},
					{
						ReleaseTime: types.MustParseRFC3339("2023-12-01T00:00:00Z"),
						Weight:      sdk.MustNewDecFromStr("0.25"),
					},
				},
				types.MustParseRFC3339("2022-05-01T00:00:00Z"),
				types.MustParseRFC3339("2023-06-01T00:00:00Z"),
			),
			nil,
		},
	} {
		s.Run(tc.name, func() {
			params := s.keeper.GetParams(s.ctx)
			s.fundAddr(tc.msg.GetAuctioneer(), params.AuctionCreationFee.Add(tc.msg.SellingCoin))

			_, err := s.msgServer.CreateFixedPriceAuction(ctx, tc.msg)
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

	auctioneer := s.addr(0)

	_, err := s.msgServer.CancelAuction(ctx, types.NewMsgCancelAuction(
		auctioneer.String(),
		uint64(1),
	))
	s.Require().ErrorIs(err, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "auction %d is not found", uint64(1)))

	// Create a fixed price auction that has started status
	startedAuction := s.createFixedPriceAuction(
		auctioneer,
		sdk.MustNewDecFromStr("0.5"),
		sdk.NewInt64Coin("denom1", 500_000_000_000),
		"denom2",
		[]types.VestingSchedule{},
		types.MustParseRFC3339("2022-01-01T00:00:00Z"),
		types.MustParseRFC3339("2022-03-01T00:00:00Z"),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, startedAuction.GetStatus())

	// Try to cancel with an incorrect auctioneer
	_, err = s.msgServer.CancelAuction(ctx, types.NewMsgCancelAuction(
		s.addr(1).String(),
		startedAuction.GetId(),
	))
	s.Require().ErrorIs(err, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "failed to verify ownership of the auction"))

	// Try to cancel with the auction that is already started
	_, err = s.msgServer.CancelAuction(ctx, types.NewMsgCancelAuction(
		startedAuction.GetAuctioneer().String(),
		startedAuction.GetId(),
	))
	s.Require().ErrorIs(err, sdkerrors.Wrap(types.ErrInvalidAuctionStatus, "auction cannot be canceled due to current status"))

	// Create another fixed price auction that is stand by status
	// Create a fixed price auction that has started status
	standByAuction := s.createFixedPriceAuction(
		auctioneer,
		sdk.MustNewDecFromStr("0.5"),
		sdk.NewInt64Coin("denom3", 500_000_000_000),
		"denom4",
		[]types.VestingSchedule{},
		types.MustParseRFC3339("2030-01-01T00:00:00Z"),
		types.MustParseRFC3339("2030-03-01T00:00:00Z"),
		true,
	)
	s.Require().Equal(types.AuctionStatusStandBy, standByAuction.GetStatus())

	// Cancel it successfully
	_, err = s.msgServer.CancelAuction(ctx, types.NewMsgCancelAuction(
		standByAuction.GetAuctioneer().String(),
		standByAuction.GetId(),
	))
	s.Require().NoError(err)

	// The selling reserve account must have zero balance
	sellingReserveAddr := standByAuction.GetSellingReserveAddress()
	s.Require().True(s.getBalance(sellingReserveAddr, standByAuction.GetSellingCoin().Denom).IsZero())
}

func (s *KeeperTestSuite) TestMsgPlaceBid() {
	ctx := sdk.WrapSDKContext(s.ctx)

	// Create a fixed price auction that has started status
	auction := s.createFixedPriceAuction(
		s.addr(0),
		sdk.MustNewDecFromStr("0.5"),
		sdk.NewInt64Coin("denom1", 500_000_000_000),
		"denom2",
		[]types.VestingSchedule{},
		types.MustParseRFC3339("2022-01-01T00:00:00Z"),
		types.MustParseRFC3339("2022-06-01T00:00:00Z"),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	// Fund the bidder
	bidder := s.addr(1)
	s.fundAddr(bidder, sdk.NewCoins(sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 5_000_000)))

	for _, tc := range []struct {
		name string
		msg  *types.MsgPlaceBid
		err  error
	}{
		{
			"valid message",
			types.NewMsgPlaceBid(
				auction.GetId(),
				bidder.String(),
				sdk.MustNewDecFromStr("0.5"),
				sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 1_000_000),
			),
			nil,
		},
		{
			"invalid start price",
			types.NewMsgPlaceBid(
				auction.GetId(),
				bidder.String(),
				sdk.MustNewDecFromStr("1.0"),
				sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 1_000_000),
			),
			sdkerrors.Wrap(types.ErrInvalidStartPrice, "bid price must be equal to start price"),
		},
		{
			"invalid paying coin denom",
			types.NewMsgPlaceBid(
				auction.GetId(),
				bidder.String(),
				sdk.MustNewDecFromStr("0.5"),
				sdk.NewInt64Coin(auction.GetSellingCoin().Denom, 1_000_000),
			),
			types.ErrInvalidPayingCoinDenom,
		},
		{
			"insufficient funds",
			types.NewMsgPlaceBid(
				auction.GetId(),
				bidder.String(),
				sdk.MustNewDecFromStr("0.5"),
				sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 500_000_000_000_000_000),
			),
			sdkerrors.ErrInsufficientFunds,
		},
	} {
		s.Run(tc.name, func() {
			receiveAmt := tc.msg.Coin.Amount.ToDec().QuoTruncate(tc.msg.Price).TruncateInt()

			err := s.keeper.AddAllowedBidders(s.ctx, tc.msg.AuctionId, []types.AllowedBidder{
				{Bidder: bidder.String(), MaxBidAmount: receiveAmt},
			})
			s.Require().NoError(err)

			_, err = s.msgServer.PlaceBid(ctx, tc.msg)
			if tc.err != nil {
				s.Require().ErrorIs(err, tc.err)
				return
			}
			s.Require().NoError(err)
		})
	}
}

func (s *KeeperTestSuite) TestMsgAddAllowedBidder() {
	ctx := sdk.WrapSDKContext(s.ctx)

	auction := s.createFixedPriceAuction(
		s.addr(0),
		sdk.MustNewDecFromStr("0.5"),
		sdk.NewInt64Coin("denom1", 500_000_000_000),
		"denom2",
		[]types.VestingSchedule{},
		types.MustParseRFC3339("2022-01-01T00:00:00Z"),
		types.MustParseRFC3339("2022-06-01T00:00:00Z"),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	for _, tc := range []struct {
		name                   string
		msg                    *types.MsgAddAllowedBidder
		enableAddAllowedBidder bool
		err                    error
	}{
		{
			"valid",
			types.NewAddAllowedBidder(
				auction.GetId(),
				types.AllowedBidder{
					Bidder:       s.addr(1).String(),
					MaxBidAmount: sdk.NewInt(100000000),
				},
			),
			true,
			nil,
		},
		{
			"invalid",
			types.NewAddAllowedBidder(
				auction.GetId(),
				types.AllowedBidder{
					Bidder:       s.addr(1).String(),
					MaxBidAmount: sdk.NewInt(100000000),
				},
			),
			false,
			sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "EnableAddAllowedBidder is disabled"),
		},
	} {
		s.Run(tc.name, func() {
			keeper.EnableAddAllowedBidder = tc.enableAddAllowedBidder

			_, err := s.msgServer.AddAllowedBidder(ctx, tc.msg)
			if tc.err != nil {
				s.Require().ErrorIs(err, tc.err)
				return
			}
			s.Require().NoError(err)
		})
	}
}
