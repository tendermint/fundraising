package keeper_test

import (
	"time"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/stretchr/testify/suite"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

func (s *KeeperTestSuite) TestLastAuctionId() {
	cacheCtx, _ := s.ctx.CacheContext()

	auctionId, err := s.keeper.AuctionSeq.Peek(cacheCtx)
	s.Require().NoError(err)
	s.Require().Equal(uint64(0), auctionId)

	nextAuctionId, err := s.keeper.AuctionSeq.Next(cacheCtx)
	s.Require().NoError(err)
	s.Require().Equal(uint64(0), nextAuctionId)

	nextAuctionId, err = s.keeper.AuctionSeq.Next(cacheCtx)
	s.Require().NoError(err)
	s.Require().Equal(uint64(1), nextAuctionId)

	s.createFixedPriceAuction(
		s.addr(0),
		math.LegacyMustNewDecFromStr("1.0"),
		parseCoin("1000000000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		time.Now().AddDate(0, 6, 0),
		time.Now().AddDate(0, 6, 0).AddDate(0, 1, 0),
		true,
	)
	nextAuctionId, err = s.keeper.AuctionSeq.Next(cacheCtx)
	s.Require().NoError(err)
	s.Require().Equal(uint64(2), nextAuctionId)

	auctions, err := s.keeper.Auctions(s.ctx)
	s.Require().NoError(err)
	s.Require().Len(auctions, 1)

	s.createFixedPriceAuction(
		s.addr(1),
		math.LegacyMustNewDecFromStr("0.5"),
		parseCoin("5000000000denom3"),
		"denom4",
		[]types.VestingSchedule{},
		time.Now().AddDate(0, 6, 0),
		time.Now().AddDate(0, 6, 0).AddDate(0, 1, 0),
		true,
	)
	nextAuctionId, err = s.keeper.AuctionSeq.Next(cacheCtx)
	s.Require().NoError(err)
	s.Require().Equal(uint64(3), nextAuctionId)

	auctions, err = s.keeper.Auctions(s.ctx)
	s.Require().NoError(err)
	s.Require().Len(auctions, 2)
}

func (s *KeeperTestSuite) TestAllowedBidderByAuction() {
	auction := s.createFixedPriceAuction(
		s.addr(0),
		math.LegacyMustNewDecFromStr("1.0"),
		parseCoin("1000000000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		time.Now().AddDate(0, 6, 0),
		time.Now().AddDate(0, 6, 0).AddDate(0, 1, 0),
		true,
	)
	s.Require().Equal(auction.GetStatus(), types.AuctionStatusStandBy)

	allowedBidders, err := s.keeper.GetAllowedBiddersByAuction(s.ctx, auction.Id)
	s.Require().NoError(err)
	s.Require().Len(allowedBidders, 0)

	// Add new allowed bidders
	newAllowedBidders := []types.AllowedBidder{
		{AuctionId: 1, Bidder: s.addr(1).String(), MaxBidAmount: parseInt("100000")},
		{AuctionId: 1, Bidder: s.addr(2).String(), MaxBidAmount: parseInt("100000")},
		{AuctionId: 1, Bidder: s.addr(3).String(), MaxBidAmount: parseInt("100000")},
	}
	err = s.keeper.AddAllowedBidders(s.ctx, auction.Id, newAllowedBidders)
	s.Require().NoError(err)

	allowedBidders, err = s.keeper.GetAllowedBiddersByAuction(s.ctx, auction.Id)
	s.Require().NoError(err)
	s.Require().Len(allowedBidders, 3)
}

func (s *KeeperTestSuite) TestLastBidId() {
	auction := s.createFixedPriceAuction(
		s.addr(0),
		math.LegacyOneDec(),
		parseCoin("500000000000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	bidId, err := s.keeper.BidSeq.Get(s.ctx, auction.Id)
	s.Require().Error(err)
	s.Require().Equal(uint64(0), bidId)

	s.placeBidFixedPrice(auction.Id, s.addr(1), math.LegacyOneDec(), parseCoin("20000000denom2"), true)
	s.placeBidFixedPrice(auction.Id, s.addr(2), math.LegacyOneDec(), parseCoin("20000000denom2"), true)
	s.placeBidFixedPrice(auction.Id, s.addr(3), math.LegacyOneDec(), parseCoin("15000000denom2"), true)

	bidsById, err := s.keeper.GetBidsByAuctionId(s.ctx, auction.GetId())
	s.Require().NoError(err)
	s.Require().Len(bidsById, 3)

	nextId, err := s.keeper.GetNextBidIdWithUpdate(s.ctx, auction.GetId())
	s.Require().NoError(err)
	s.Require().Equal(uint64(4), nextId)

	// Create another auction
	auction2 := s.createFixedPriceAuction(
		s.addr(0),
		math.LegacyMustNewDecFromStr("0.5"),
		parseCoin("1000000000000denom3"),
		"denom4",
		[]types.VestingSchedule{},
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)

	// Bid id must start from 1 with new auction
	bidsById, err = s.keeper.GetBidsByAuctionId(s.ctx, auction2.GetId())
	s.Require().NoError(err)
	s.Require().Len(bidsById, 0)

	nextId, err = s.keeper.GetNextBidIdWithUpdate(s.ctx, auction2.GetId())
	s.Require().NoError(err)
	s.Require().Equal(uint64(1), nextId)
}

func (s *KeeperTestSuite) TestIterateBids() {
	startedAuction := s.createFixedPriceAuction(
		s.addr(0),
		math.LegacyOneDec(),
		parseCoin("500000000000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)

	auction, err := s.keeper.Auction.Get(s.ctx, startedAuction.GetId())
	s.Require().NoError(err)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	s.placeBidFixedPrice(auction.GetId(), s.addr(1), math.LegacyOneDec(), parseCoin("20000000denom2"), true)
	s.placeBidFixedPrice(auction.GetId(), s.addr(2), math.LegacyOneDec(), parseCoin("20000000denom2"), true)
	s.placeBidFixedPrice(auction.GetId(), s.addr(2), math.LegacyOneDec(), parseCoin("15000000denom2"), true)
	s.placeBidFixedPrice(auction.GetId(), s.addr(3), math.LegacyOneDec(), parseCoin("35000000denom2"), true)

	bids, err := s.keeper.Bids(s.ctx)
	s.Require().NoError(err)
	s.Require().Len(bids, 4)

	bidsById, err := s.keeper.GetBidsByAuctionId(s.ctx, auction.GetId())
	s.Require().NoError(err)
	s.Require().Len(bidsById, 4)

	bidsByBidder, err := s.keeper.GetBidsByBidder(s.ctx, s.addr(2))
	s.Require().NoError(err)
	s.Require().Len(bidsByBidder, 2)
}

func (s *KeeperTestSuite) TestVestingQueue() {
	vestingQueue := types.NewVestingQueue(
		1,
		s.addr(1),
		parseCoin("100_000_000denom1"),
		types.MustParseRFC3339("2023-01-01T00:00:00Z"),
		false,
	)
	err := s.keeper.VestingQueue.Set(
		s.ctx,
		collections.Join(
			vestingQueue.AuctionId,
			vestingQueue.ReleaseTime,
		),
		vestingQueue,
	)
	s.Require().NoError(err)

	vq, err := s.keeper.VestingQueue.Get(s.ctx, collections.Join(vestingQueue.AuctionId, vestingQueue.ReleaseTime))
	s.Require().NoError(err)
	s.Require().EqualValues(vestingQueue, vq)
}

func (s *KeeperTestSuite) TestVestingQueueIterator() {
	payingReserveAddress := s.addr(0)
	payingCoinDenom := "denom1"
	reserveCoin := s.getBalance(payingReserveAddress, payingCoinDenom)

	// Set vesting schedules with 2 vesting queues
	for _, vs := range []types.VestingSchedule{
		{
			ReleaseTime: types.MustParseRFC3339("2023-01-01T00:00:00Z"),
			Weight:      math.LegacyMustNewDecFromStr("0.5"),
		},
		{
			ReleaseTime: types.MustParseRFC3339("2023-06-01T00:00:00Z"),
			Weight:      math.LegacyMustNewDecFromStr("0.5"),
		},
	} {
		payingAmt := math.LegacyNewDecFromInt(reserveCoin.Amount).MulTruncate(vs.Weight).TruncateInt()

		vestingQueue := types.VestingQueue{
			AuctionId:   uint64(1),
			Auctioneer:  s.addr(1).String(),
			PayingCoin:  sdk.NewCoin(payingCoinDenom, payingAmt),
			ReleaseTime: vs.ReleaseTime,
			Released:    false,
		}
		err := s.keeper.VestingQueue.Set(
			s.ctx,
			collections.Join(
				vestingQueue.AuctionId,
				vestingQueue.ReleaseTime,
			),
			vestingQueue,
		)
		s.Require().NoError(err)
	}

	// Set vesting schedules with 4 vesting queues
	for _, vs := range []types.VestingSchedule{
		{
			ReleaseTime: types.MustParseRFC3339("2023-01-01T00:00:00Z"),
			Weight:      math.LegacyMustNewDecFromStr("0.25"),
		},
		{
			ReleaseTime: types.MustParseRFC3339("2023-05-01T00:00:00Z"),
			Weight:      math.LegacyMustNewDecFromStr("0.25"),
		},
		{
			ReleaseTime: types.MustParseRFC3339("2023-09-01T00:00:00Z"),
			Weight:      math.LegacyMustNewDecFromStr("0.25"),
		},
		{
			ReleaseTime: types.MustParseRFC3339("2023-12-01T00:00:00Z"),
			Weight:      math.LegacyMustNewDecFromStr("0.25"),
		},
	} {
		payingAmt := math.LegacyNewDecFromInt(reserveCoin.Amount).MulTruncate(vs.Weight).TruncateInt()

		vestingQueue := types.VestingQueue{
			AuctionId:   uint64(2),
			Auctioneer:  s.addr(2).String(),
			PayingCoin:  sdk.NewCoin(payingCoinDenom, payingAmt),
			ReleaseTime: vs.ReleaseTime,
			Released:    false,
		}
		err := s.keeper.VestingQueue.Set(
			s.ctx,
			collections.Join(
				vestingQueue.AuctionId,
				vestingQueue.ReleaseTime,
			),
			vestingQueue,
		)
		s.Require().NoError(err)
	}

	vq1, err := s.keeper.GetVestingQueuesByAuctionId(s.ctx, uint64(1))
	s.Require().NoError(err)
	s.Require().Len(vq1, 2)

	vq2, err := s.keeper.GetVestingQueuesByAuctionId(s.ctx, uint64(2))
	s.Require().NoError(err)
	s.Require().Len(vq2, 4)

	vqs, err := s.keeper.VestingQueues(s.ctx)
	s.Require().NoError(err)
	s.Require().Len(vqs, 6)

	totalPayingCoin := sdk.NewInt64Coin(payingCoinDenom, 0)
	for _, vq := range vq2 {
		totalPayingCoin = totalPayingCoin.Add(vq.PayingCoin)
	}
	s.Require().Equal(reserveCoin, totalPayingCoin)
}
