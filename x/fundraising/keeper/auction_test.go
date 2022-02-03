package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/tendermint/fundraising/x/fundraising"
	"github.com/tendermint/fundraising/x/fundraising/types"

	_ "github.com/stretchr/testify/suite"
)

func (s *KeeperTestSuite) TestAuctionId() {
	auctionId := s.keeper.GetLastAuctionId(s.ctx)
	s.Require().Equal(uint64(0), auctionId)

	cacheCtx, _ := s.ctx.CacheContext()
	nextAuctionId := s.keeper.GetNextAuctionIdWithUpdate(cacheCtx)
	s.Require().Equal(uint64(1), nextAuctionId)

	s.createFixedPriceAuction(
		s.addr(0),
		sdk.MustNewDecFromStr("1.0"),
		sdk.NewInt64Coin("denom1", 1_000_000_000_000),
		"denom2",
		[]types.VestingSchedule{},
		types.MustParseRFC3339("2023-01-01T00:00:00Z"),
		types.MustParseRFC3339("2023-01-10T00:00:00Z"),
		true,
	)
	nextAuctionId = s.keeper.GetNextAuctionIdWithUpdate(cacheCtx)
	s.Require().Equal(uint64(2), nextAuctionId)

	auctions := s.keeper.GetAuctions(s.ctx)
	s.Require().Len(auctions, 1)

	s.createFixedPriceAuction(
		s.addr(1),
		sdk.MustNewDecFromStr("0.5"),
		sdk.NewInt64Coin("denom3", 500_000_000_000),
		"denom4",
		[]types.VestingSchedule{},
		types.MustParseRFC3339("2023-02-01T00:00:00Z"),
		types.MustParseRFC3339("2023-02-10T00:00:00Z"),
		true,
	)
	nextAuctionId = s.keeper.GetNextAuctionIdWithUpdate(cacheCtx)
	s.Require().Equal(uint64(3), nextAuctionId)

	auctions = s.keeper.GetAuctions(s.ctx)
	s.Require().Len(auctions, 2)
}

func (s *KeeperTestSuite) TestAuctionStatus() {
	standByAuction := s.createFixedPriceAuction(
		s.addr(0),
		sdk.MustNewDecFromStr("0.5"),
		sdk.NewInt64Coin("denom3", 500_000_000_000),
		"denom4",
		[]types.VestingSchedule{},
		types.MustParseRFC3339("2023-01-01T00:00:00Z"),
		types.MustParseRFC3339("2023-01-10T00:00:00Z"),
		true,
	)
	auction, found := s.keeper.GetAuction(s.ctx, standByAuction.GetId())
	s.Require().True(found)
	s.Require().Equal(types.AuctionStatusStandBy, auction.GetStatus())

	startedAuction := s.createFixedPriceAuction(
		s.addr(1),
		sdk.MustNewDecFromStr("0.5"),
		sdk.NewInt64Coin("denom3", 500_000_000_000),
		"denom4",
		[]types.VestingSchedule{},
		types.MustParseRFC3339("2022-01-01T00:00:00Z"),
		types.MustParseRFC3339("2022-03-10T00:00:00Z"),
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
		sdk.NewInt64Coin("denom1", 1_000_000_000_000),
		"denom2",
		[]types.VestingSchedule{},
		types.MustParseRFC3339("2022-01-01T00:00:00Z"),
		types.MustParseRFC3339("2022-03-01T00:00:00Z"),
		true,
	)
	_, found := s.keeper.GetAuction(s.ctx, auction.Id)
	s.Require().True(found)

	bidder1 := s.addr(1)
	bidder2 := s.addr(2)
	bidder3 := s.addr(3)

	// Place bids
	s.placeBid(auction.Id, bidder1, sdk.OneDec(), sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 100_000_000), true)
	s.placeBid(auction.Id, bidder2, sdk.OneDec(), sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 200_000_000), true)
	s.placeBid(auction.Id, bidder3, sdk.OneDec(), sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 200_000_000), true)

	// Distribute the selling coin and the selling reserve account must be empty afterwards
	err := s.keeper.DistributeSellingCoin(s.ctx, auction)
	s.Require().NoError(err)
	s.Require().Equal(
		sdk.NewCoin(auction.GetSellingCoin().Denom, sdk.ZeroInt()),
		s.app.BankKeeper.GetBalance(s.ctx, auction.GetSellingReserveAddress(), auction.GetSellingCoin().Denom))

	// The bidders must have the selling coin
	s.Require().False(s.getBalance(bidder1, auction.GetSellingCoin().Denom).IsZero())
	s.Require().False(s.getBalance(bidder2, auction.GetSellingCoin().Denom).IsZero())
	s.Require().False(s.getBalance(bidder3, auction.GetSellingCoin().Denom).IsZero())
}

func (s *KeeperTestSuite) TestDistributePayingCoin() {
	auction := s.createFixedPriceAuction(
		s.addr(0),
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
		types.MustParseRFC3339("2022-01-01T00:00:00Z"),
		types.MustParseRFC3339("2022-03-01T00:00:00Z"),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	// Place bids
	s.placeBid(auction.GetId(), s.addr(1), sdk.OneDec(), sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 100_000_000), true)
	s.placeBid(auction.GetId(), s.addr(1), sdk.OneDec(), sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 200_000_000), true)
	s.placeBid(auction.GetId(), s.addr(1), sdk.OneDec(), sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 200_000_000), true)

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
	auctioneer := s.addr(0)

	standByAuction := s.createFixedPriceAuction(
		auctioneer,
		sdk.MustNewDecFromStr("1.0"),
		sdk.NewInt64Coin("denom1", 500_000_000_000),
		"denom2",
		[]types.VestingSchedule{},
		types.MustParseRFC3339("2023-01-01T00:00:00Z"),
		types.MustParseRFC3339("2023-02-01T00:00:00Z"),
		true,
	)
	s.Require().Equal(types.AuctionStatusStandBy, standByAuction.GetStatus())

	// Cancel the auction since it is not started yet
	auction := s.cancelAuction(standByAuction.GetId(), auctioneer)
	s.Require().Equal(types.AuctionStatusCancelled, auction.GetStatus())

	// The selling reserve balance must be zero
	sellingCoinDenom := auction.GetSellingCoin().Denom
	sellingReserve := s.getBalance(auction.GetSellingReserveAddress(), sellingCoinDenom)
	s.Require().True(coinEq(sdk.NewCoin(sellingCoinDenom, sdk.ZeroInt()), sellingReserve))
}

func (s *KeeperTestSuite) TestAddAllowedBidder() {
	startedAuction := s.createFixedPriceAuction(
		s.addr(0),
		sdk.MustNewDecFromStr("0.5"),
		sdk.NewInt64Coin("denom1", 500_000_000_000),
		"denom2",
		[]types.VestingSchedule{},
		types.MustParseRFC3339("2022-01-01T00:00:00Z"),
		types.MustParseRFC3339("2022-06-10T00:00:00Z"),
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

func (s *KeeperTestSuite) TestUpdateAllowedBidder() {
	startedAuction := s.createFixedPriceAuction(
		s.addr(0),
		sdk.MustNewDecFromStr("0.5"),
		sdk.NewInt64Coin("denom1", 500_000_000_000),
		"denom2",
		[]types.VestingSchedule{},
		types.MustParseRFC3339("2022-01-01T00:00:00Z"),
		types.MustParseRFC3339("2022-06-10T00:00:00Z"),
		true,
	)

	auction, found := s.keeper.GetAuction(s.ctx, startedAuction.GetId())
	s.Require().True(found)
	s.Require().Len(auction.GetAllowedBidders(), 0)

	// Add 5 bidders with different maximum bid amount
	err := s.keeper.AddAllowedBidders(s.ctx, auction.GetId(), []types.AllowedBidder{
		{Bidder: s.addr(1).String(), MaxBidAmount: sdk.NewInt(100_000_000)},
		{Bidder: s.addr(2).String(), MaxBidAmount: sdk.NewInt(200_000_000)},
		{Bidder: s.addr(3).String(), MaxBidAmount: sdk.NewInt(300_000_000)},
		{Bidder: s.addr(4).String(), MaxBidAmount: sdk.NewInt(400_000_000)},
		{Bidder: s.addr(5).String(), MaxBidAmount: sdk.NewInt(500_000_000)},
	})
	s.Require().NoError(err)

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

			// Check if it is sucessfully updated
			allowedBiddersMap := make(map[string]sdk.Int)
			for _, bidder := range auction.GetAllowedBidders() {
				allowedBiddersMap[bidder.GetBidder()] = bidder.MaxBidAmount
			}
			s.Require().Equal(tc.maxBidAmount, allowedBiddersMap[tc.bidder.String()])
		})
	}
}
