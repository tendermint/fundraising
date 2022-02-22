package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/tendermint/fundraising/x/fundraising"
	"github.com/tendermint/fundraising/x/fundraising/types"

	_ "github.com/stretchr/testify/suite"
)

func (s *KeeperTestSuite) TestLastAuctionId() {
	auctionId := s.keeper.GetLastAuctionId(s.ctx)
	s.Require().Equal(uint64(0), auctionId)

	cacheCtx, _ := s.ctx.CacheContext()
	nextAuctionId := s.keeper.GetNextAuctionIdWithUpdate(cacheCtx)
	s.Require().Equal(uint64(1), nextAuctionId)

	s.createFixedPriceAuction(
		s.addr(0),
		sdk.MustNewDecFromStr("1.0"),
		parseCoin("1000000000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		time.Now().AddDate(0, 6, 0),
		time.Now().AddDate(0, 6, 0).AddDate(0, 1, 0),
		true,
	)
	nextAuctionId = s.keeper.GetNextAuctionIdWithUpdate(cacheCtx)
	s.Require().Equal(uint64(2), nextAuctionId)

	auctions := s.keeper.GetAuctions(s.ctx)
	s.Require().Len(auctions, 1)

	s.createFixedPriceAuction(
		s.addr(1),
		sdk.MustNewDecFromStr("0.5"),
		parseCoin("5000000000denom3"),
		"denom4",
		[]types.VestingSchedule{},
		time.Now().AddDate(0, 6, 0),
		time.Now().AddDate(0, 6, 0).AddDate(0, 1, 0),
		true,
	)
	nextAuctionId = s.keeper.GetNextAuctionIdWithUpdate(cacheCtx)
	s.Require().Equal(uint64(3), nextAuctionId)

	auctions = s.keeper.GetAuctions(s.ctx)
	s.Require().Len(auctions, 2)
}

func (s *KeeperTestSuite) TestFixedPriceAuction_AuctionStatus() {
	standByAuction := s.createFixedPriceAuction(
		s.addr(0),
		sdk.MustNewDecFromStr("0.5"),
		parseCoin("5000000000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		time.Now().AddDate(0, 6, 0),
		time.Now().AddDate(0, 6, 0).AddDate(0, 1, 0),
		true,
	)

	auction, found := s.keeper.GetAuction(s.ctx, standByAuction.GetId())
	s.Require().True(found)
	s.Require().Equal(types.AuctionStatusStandBy, auction.GetStatus())

	startedAuction := s.createFixedPriceAuction(
		s.addr(1),
		sdk.MustNewDecFromStr("0.5"),
		parseCoin("1000000000000denom3"),
		"denom4",
		[]types.VestingSchedule{},
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)

	auction, found = s.keeper.GetAuction(s.ctx, startedAuction.GetId())
	s.Require().True(found)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())
}

func (s *KeeperTestSuite) TestFixedPriceAuction_InsufficientRemainingAmount() {
	auction := s.createFixedPriceAuction(
		s.addr(0),
		sdk.OneDec(),
		parseCoin("1000000000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	s.placeBid(auction.GetId(), s.addr(1), types.BidTypeFixedPrice, sdk.OneDec(), parseCoin("200000000denom2"), true)
	s.placeBid(auction.GetId(), s.addr(2), types.BidTypeFixedPrice, sdk.OneDec(), parseCoin("200000000denom2"), true)
	s.placeBid(auction.GetId(), s.addr(3), types.BidTypeFixedPrice, sdk.OneDec(), parseCoin("250000000denom2"), true)
	s.placeBid(auction.GetId(), s.addr(4), types.BidTypeFixedPrice, sdk.OneDec(), parseCoin("250000000denom2"), true)

	// Test insufficient remaining coin amount
	coin := parseCoin("300000000denom2")
	s.fundAddr(s.addr(5), sdk.NewCoins(coin))

	receiveAmt := coin.Amount.ToDec().QuoTruncate(sdk.OneDec()).TruncateInt()
	s.addAllowedBidder(auction.GetId(), s.addr(5), receiveAmt)

	_, err := s.keeper.PlaceBid(s.ctx, &types.MsgPlaceBid{
		AuctionId: auction.GetId(),
		Bidder:    s.addr(5).String(),
		Price:     sdk.OneDec(),
		Coin:      coin,
	})
	s.Require().ErrorIs(err, types.ErrInsufficientRemainingAmount)
}

func (s *KeeperTestSuite) TestBatchAuction_AuctionStatus() {
	standByAuction := s.createBatchAuction(
		s.addr(0),
		sdk.MustNewDecFromStr("0.5"),
		parseCoin("5000000000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		1,
		sdk.MustNewDecFromStr("0.2"),
		time.Now().AddDate(0, 6, 0),
		time.Now().AddDate(0, 6, 0).AddDate(0, 1, 0),
		true,
	)

	auction, found := s.keeper.GetAuction(s.ctx, standByAuction.GetId())
	s.Require().True(found)
	s.Require().Equal(types.AuctionStatusStandBy, auction.GetStatus())

	startedAuction := s.createBatchAuction(
		s.addr(1),
		sdk.MustNewDecFromStr("0.5"),
		parseCoin("5000000000denom3"),
		"denom4",
		[]types.VestingSchedule{},
		1,
		sdk.MustNewDecFromStr("0.2"),
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)

	auction, found = s.keeper.GetAuction(s.ctx, startedAuction.GetId())
	s.Require().True(found)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())
}

func (s *KeeperTestSuite) TestDistributeSellingCoin() {
	auction := s.createFixedPriceAuction(
		s.addr(0),
		sdk.OneDec(),
		parseCoin("1000000000000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)

	_, found := s.keeper.GetAuction(s.ctx, auction.Id)
	s.Require().True(found)

	// Place bids
	s.placeBid(auction.Id, s.addr(1), types.BidTypeFixedPrice, sdk.OneDec(), parseCoin("100000000denom2"), true)
	s.placeBid(auction.Id, s.addr(2), types.BidTypeFixedPrice, sdk.OneDec(), parseCoin("200000000denom2"), true)
	s.placeBid(auction.Id, s.addr(3), types.BidTypeFixedPrice, sdk.OneDec(), parseCoin("200000000denom2"), true)

	// Distribute selling coin
	err := s.keeper.DistributeSellingCoin(s.ctx, auction)
	s.Require().NoError(err)

	// The selling reserve account balance must be zero
	s.Require().True(s.getBalance(auction.GetSellingReserveAddress(), auction.SellingCoin.Denom).IsZero())

	// The bidders must have the selling coin
	s.Require().False(s.getBalance(s.addr(1), auction.GetSellingCoin().Denom).IsZero())
	s.Require().False(s.getBalance(s.addr(2), auction.GetSellingCoin().Denom).IsZero())
	s.Require().False(s.getBalance(s.addr(3), auction.GetSellingCoin().Denom).IsZero())
}

func (s *KeeperTestSuite) TestDistributePayingCoin() {
	auction := s.createFixedPriceAuction(
		s.addr(0),
		sdk.OneDec(),
		parseCoin("1000000000000denom1"),
		"denom2",
		[]types.VestingSchedule{
			{
				ReleaseTime: time.Now().AddDate(0, 6, 0),
				Weight:      sdk.MustNewDecFromStr("0.25"),
			},
			{
				ReleaseTime: time.Now().AddDate(0, 9, 0),
				Weight:      sdk.MustNewDecFromStr("0.25"),
			},
			{
				ReleaseTime: time.Now().AddDate(1, 0, 0),
				Weight:      sdk.MustNewDecFromStr("0.25"),
			},
			{
				ReleaseTime: time.Now().AddDate(1, 3, 0),
				Weight:      sdk.MustNewDecFromStr("0.25"),
			},
		},
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	// Place bids
	s.placeBid(auction.GetId(), s.addr(1), types.BidTypeFixedPrice, sdk.OneDec(), parseCoin("100000000denom2"), true)
	s.placeBid(auction.GetId(), s.addr(1), types.BidTypeFixedPrice, sdk.OneDec(), parseCoin("200000000denom2"), true)
	s.placeBid(auction.GetId(), s.addr(1), types.BidTypeFixedPrice, sdk.OneDec(), parseCoin("200000000denom2"), true)

	// Distribute selling coin
	err := s.keeper.DistributeSellingCoin(s.ctx, auction)
	s.Require().NoError(err)

	// Set vesting schedules
	err = s.keeper.SetVestingSchedules(s.ctx, auction)
	s.Require().NoError(err)

	// All of the vesting queues must not be released yet
	vqs := s.keeper.GetVestingQueuesByAuctionId(s.ctx, auction.GetId())
	s.Require().Equal(4, len(vqs))
	for _, vq := range vqs {
		s.Require().False(vq.Released)
	}

	// Change the block time to release two vesting schedules
	s.ctx = s.ctx.WithBlockTime(vqs[0].GetReleaseTime().AddDate(0, 4, 1))
	fundraising.EndBlocker(s.ctx, s.keeper)

	// Distribute paying coin
	err = s.keeper.DistributePayingCoin(s.ctx, auction)
	s.Require().NoError(err)

	// First two vesting queues must be released
	for i, vq := range s.keeper.GetVestingQueuesByAuctionId(s.ctx, auction.GetId()) {
		if i == 0 || i == 1 {
			s.Require().True(vq.Released)
		} else {
			s.Require().False(vq.Released)
		}
	}

	// Change the block time
	s.ctx = s.ctx.WithBlockTime(vqs[3].GetReleaseTime().AddDate(0, 0, 1))
	fundraising.EndBlocker(s.ctx, s.keeper)
	s.Require().NoError(s.keeper.DistributePayingCoin(s.ctx, auction))

	// All of the vesting queues must be released
	for _, vq := range s.keeper.GetVestingQueuesByAuctionId(s.ctx, auction.GetId()) {
		s.Require().True(vq.Released)
	}

	finishedAuction, found := s.keeper.GetAuction(s.ctx, auction.GetId())
	s.Require().True(found)
	s.Require().Equal(types.AuctionStatusFinished, finishedAuction.GetStatus())
}

func (s *KeeperTestSuite) TestCancelAuction() {
	standByAuction := s.createFixedPriceAuction(
		s.addr(0),
		sdk.MustNewDecFromStr("1.0"),
		parseCoin("500000000000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		time.Now().AddDate(0, 1, 0),
		time.Now().AddDate(0, 1, 0).AddDate(0, 1, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStandBy, standByAuction.GetStatus())

	// Cancel the auction
	auction := s.cancelAuction(standByAuction.GetId(), s.addr(0))
	s.Require().Equal(types.AuctionStatusCancelled, auction.GetStatus())

	// The selling reserve balance must be zero
	sellingReserveAddr := auction.GetSellingReserveAddress()
	sellingCoinDenom := auction.GetSellingCoin().Denom
	s.Require().True(s.getBalance(sellingReserveAddr, sellingCoinDenom).IsZero())
}

func (s *KeeperTestSuite) TestAddAllowedBidders() {
	startedAuction := s.createFixedPriceAuction(
		s.addr(0),
		sdk.MustNewDecFromStr("0.5"),
		sdk.NewInt64Coin("denom1", 500_000_000_000),
		"denom2",
		[]types.VestingSchedule{},
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)

	auction, found := s.keeper.GetAuction(s.ctx, startedAuction.GetId())
	s.Require().True(found)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())
	s.Require().Len(auction.GetAllowedBidders(), 0)

	for _, tc := range []struct {
		name        string
		bidders     []types.AllowedBidder
		expectedErr error
	}{
		{
			"single bidder",
			[]types.AllowedBidder{
				{
					Bidder:       s.addr(1).String(),
					MaxBidAmount: sdk.NewInt(100_000_000),
				},
			},
			nil,
		},
		{
			"multiple bidders",
			[]types.AllowedBidder{
				{
					Bidder:       s.addr(1).String(),
					MaxBidAmount: sdk.NewInt(100_000_000),
				},
				{
					Bidder:       s.addr(2).String(),
					MaxBidAmount: sdk.NewInt(500_000_000),
				},
				{
					Bidder:       s.addr(3).String(),
					MaxBidAmount: sdk.NewInt(800_000_000),
				},
			},
			nil,
		},
		{

			"empty bidders",
			[]types.AllowedBidder{},
			types.ErrEmptyAllowedBidders,
		},
		{
			"zero maximum bid amount value",
			[]types.AllowedBidder{
				{
					Bidder:       s.addr(1).String(),
					MaxBidAmount: sdk.NewInt(0),
				},
			},
			types.ErrInvalidMaxBidAmount,
		},
		{
			"negative maximum bid amount value",
			[]types.AllowedBidder{
				{
					Bidder:       s.addr(1).String(),
					MaxBidAmount: sdk.NewInt(-1),
				},
			},
			types.ErrInvalidMaxBidAmount,
		},
	} {
		s.Run(tc.name, func() {
			err := s.keeper.AddAllowedBidders(s.ctx, auction.GetId(), tc.bidders)
			if tc.expectedErr != nil {
				s.Require().ErrorIs(err, tc.expectedErr)
				return
			}
			s.Require().NoError(err)
		})
	}
}

func (s *KeeperTestSuite) TestAddAllowedBidders_Length() {
	startedAuction := s.createFixedPriceAuction(
		s.addr(0),
		sdk.MustNewDecFromStr("0.5"),
		sdk.NewInt64Coin("denom1", 500_000_000_000),
		"denom2",
		[]types.VestingSchedule{},
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)

	auction, found := s.keeper.GetAuction(s.ctx, startedAuction.GetId())
	s.Require().True(found)
	s.Require().Len(auction.GetAllowedBidders(), 0)

	// Add some bidders
	s.Require().NoError(s.keeper.AddAllowedBidders(s.ctx, auction.GetId(), []types.AllowedBidder{
		{Bidder: s.addr(1).String(), MaxBidAmount: sdk.NewInt(100_000_000)},
		{Bidder: s.addr(2).String(), MaxBidAmount: sdk.NewInt(500_000_000)},
	}))

	auction, found = s.keeper.GetAuction(s.ctx, auction.GetId())
	s.Require().True(found)
	s.Require().Len(auction.GetAllowedBidders(), 2)

	// Add more bidders
	s.Require().NoError(s.keeper.AddAllowedBidders(s.ctx, auction.GetId(), []types.AllowedBidder{
		{Bidder: s.addr(3).String(), MaxBidAmount: sdk.NewInt(100_000_000)},
		{Bidder: s.addr(4).String(), MaxBidAmount: sdk.NewInt(100_000_000)},
		{Bidder: s.addr(5).String(), MaxBidAmount: sdk.NewInt(100_000_000)},
	}))

	auction, found = s.keeper.GetAuction(s.ctx, auction.GetId())
	s.Require().True(found)
	s.Require().Len(auction.GetAllowedBidders(), 5)
}

func (s *KeeperTestSuite) TestUpdateAllowedBidder() {
	startedAuction := s.createFixedPriceAuction(
		s.addr(0),
		sdk.MustNewDecFromStr("0.5"),
		sdk.NewInt64Coin("denom1", 500_000_000_000),
		"denom2",
		[]types.VestingSchedule{},
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)

	auction, found := s.keeper.GetAuction(s.ctx, startedAuction.GetId())
	s.Require().True(found)
	s.Require().Len(auction.GetAllowedBidders(), 0)

	// Add 5 bidders with different maximum bid amount
	s.Require().NoError(s.keeper.AddAllowedBidders(s.ctx, auction.GetId(), []types.AllowedBidder{
		{Bidder: s.addr(1).String(), MaxBidAmount: sdk.NewInt(100_000_000)},
		{Bidder: s.addr(2).String(), MaxBidAmount: sdk.NewInt(200_000_000)},
		{Bidder: s.addr(3).String(), MaxBidAmount: sdk.NewInt(300_000_000)},
		{Bidder: s.addr(4).String(), MaxBidAmount: sdk.NewInt(400_000_000)},
		{Bidder: s.addr(5).String(), MaxBidAmount: sdk.NewInt(500_000_000)},
	}))

	auction, found = s.keeper.GetAuction(s.ctx, startedAuction.GetId())
	s.Require().True(found)
	s.Require().Len(auction.GetAllowedBidders(), 5)

	for _, tc := range []struct {
		name         string
		bidder       sdk.AccAddress
		maxBidAmount sdk.Int
		expectedErr  error
	}{
		{
			"update bidder's maximum bid amount",
			s.addr(1),
			sdk.NewInt(555_000_000_000),
			nil,
		},
		{
			"bidder not found",
			s.addr(10),
			sdk.NewInt(300_000_000),
			sdkerrors.Wrapf(sdkerrors.ErrNotFound, "bidder %s is not found", s.addr(10).String()),
		},
		{
			"zero maximum bid amount value",
			s.addr(1),
			sdk.NewInt(0),
			types.ErrInvalidMaxBidAmount,
		},
		{
			"negative maximum bid amount value",
			s.addr(1),
			sdk.NewInt(-1),
			types.ErrInvalidMaxBidAmount,
		},
	} {
		s.Run(tc.name, func() {
			err := s.keeper.UpdateAllowedBidder(s.ctx, auction.GetId(), tc.bidder, tc.maxBidAmount)
			if tc.expectedErr != nil {
				s.Require().ErrorIs(err, tc.expectedErr)
				return
			}
			s.Require().NoError(err)

			auction, found = s.keeper.GetAuction(s.ctx, auction.GetId())
			s.Require().True(found)
			s.Require().Len(auction.GetAllowedBidders(), 5)

			// Check if it is sucessfully updated
			allowedBiddersMap := make(map[string]sdk.Int)
			for _, bidder := range auction.GetAllowedBidders() {
				allowedBiddersMap[bidder.GetBidder()] = bidder.MaxBidAmount
			}
			s.Require().Equal(tc.maxBidAmount, allowedBiddersMap[tc.bidder.String()])
		})
	}
}
