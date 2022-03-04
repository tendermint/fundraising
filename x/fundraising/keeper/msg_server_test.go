package keeper_test

import (
	"time"

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
						ReleaseTime: time.Now().AddDate(1, 0, 0),
						Weight:      sdk.MustNewDecFromStr("0.25"),
					},
					{
						ReleaseTime: time.Now().AddDate(1, 3, 0),
						Weight:      sdk.MustNewDecFromStr("0.25"),
					},
					{
						ReleaseTime: time.Now().AddDate(1, 6, 0),
						Weight:      sdk.MustNewDecFromStr("0.25"),
					},
					{
						ReleaseTime: time.Now().AddDate(1, 9, 0),
						Weight:      sdk.MustNewDecFromStr("0.25"),
					},
				},
				time.Now(),
				time.Now().AddDate(0, 1, 0),
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

func (s *KeeperTestSuite) TestMsgCreateBatchAuction() {
	ctx := sdk.WrapSDKContext(s.ctx)

	for _, tc := range []struct {
		name string
		msg  *types.MsgCreateBatchAuction
		err  error
	}{
		{
			"valid message with the future start time",
			types.NewMsgCreateBatchAuction(
				s.addr(0).String(),
				sdk.MustNewDecFromStr("0.1"),
				sdk.NewInt64Coin("denom1", 1_000_000_000_000),
				"denom2",
				[]types.VestingSchedule{
					{
						ReleaseTime: time.Now().AddDate(1, 0, 0),
						Weight:      sdk.MustNewDecFromStr("0.25"),
					},
					{
						ReleaseTime: time.Now().AddDate(1, 3, 0),
						Weight:      sdk.MustNewDecFromStr("0.25"),
					},
					{
						ReleaseTime: time.Now().AddDate(1, 6, 0),
						Weight:      sdk.MustNewDecFromStr("0.25"),
					},
					{
						ReleaseTime: time.Now().AddDate(1, 9, 0),
						Weight:      sdk.MustNewDecFromStr("0.25"),
					},
				},
				1,
				sdk.MustNewDecFromStr("0.2"),
				time.Now(),
				time.Now().AddDate(0, 1, 0),
			),
			nil,
		},
	} {
		s.Run(tc.name, func() {
			params := s.keeper.GetParams(s.ctx)
			s.fundAddr(tc.msg.GetAuctioneer(), params.AuctionCreationFee.Add(tc.msg.SellingCoin))

			_, err := s.msgServer.CreateBatchAuction(ctx, tc.msg)
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

func (s *KeeperTestSuite) TestMsgCancelAuction() {
	ctx := sdk.WrapSDKContext(s.ctx)

	auctioneer := s.addr(0)

	_, err := s.msgServer.CancelAuction(ctx, types.NewMsgCancelAuction(
		auctioneer.String(),
		uint64(1),
	))
	s.Require().ErrorIs(err, sdkerrors.ErrNotFound)

	// Create a fixed price auction that has started status
	startedAuction := s.createFixedPriceAuction(
		auctioneer,
		sdk.MustNewDecFromStr("0.5"),
		sdk.NewInt64Coin("denom1", 500_000_000_000),
		"denom2",
		[]types.VestingSchedule{},
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 1, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, startedAuction.GetStatus())

	// Try to cancel with an incorrect auctioneer
	_, err = s.msgServer.CancelAuction(ctx, types.NewMsgCancelAuction(
		s.addr(1).String(),
		startedAuction.GetId(),
	))
	s.Require().ErrorIs(err, sdkerrors.ErrInvalidRequest)

	// Try to cancel with the auction that is already started
	_, err = s.msgServer.CancelAuction(ctx, types.NewMsgCancelAuction(
		startedAuction.GetAuctioneer().String(),
		startedAuction.GetId(),
	))
	s.Require().ErrorIs(err, types.ErrInvalidAuctionStatus)

	// Create another fixed price auction that is stand by status
	// Create a fixed price auction that has started status
	standByAuction := s.createFixedPriceAuction(
		auctioneer,
		sdk.MustNewDecFromStr("0.5"),
		sdk.NewInt64Coin("denom3", 500_000_000_000),
		"denom4",
		[]types.VestingSchedule{},
		time.Now().AddDate(0, 3, 0),
		time.Now().AddDate(0, 4, 0),
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
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 1, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	for _, tc := range []struct {
		name string
		msg  *types.MsgPlaceBid
		err  error
	}{
		{
			"valid message",
			types.NewMsgPlaceBid(
				auction.GetId(),
				s.addr(1).String(),
				types.BidTypeFixedPrice,
				auction.StartPrice,
				parseCoin("1000000denom2"),
			),
			nil,
		},
		{
			"invalid start price",
			types.NewMsgPlaceBid(
				auction.GetId(),
				s.addr(1).String(),
				types.BidTypeFixedPrice,
				sdk.MustNewDecFromStr("1.0"),
				parseCoin("1000000denom2"),
			),
			sdkerrors.Wrap(types.ErrInvalidStartPrice, "bid price must be equal to start price"),
		},
		{
			"incorrect coin denom",
			types.NewMsgPlaceBid(
				auction.GetId(),
				s.addr(1).String(),
				types.BidTypeFixedPrice,
				auction.StartPrice,
				parseCoin("1000000denom1"),
			),
			types.ErrIncorrectCoinDenom,
		},
	} {
		s.Run(tc.name, func() {
			s.fundAddr(tc.msg.GetBidder(), sdk.NewCoins(tc.msg.Coin))
			s.addAllowedBidder(auction.Id, tc.msg.GetBidder(), exchangedSellingAmount(tc.msg.Price, tc.msg.Coin))

			_, err := s.msgServer.PlaceBid(ctx, tc.msg)
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
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 1, 0),
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
