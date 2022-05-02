package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

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

func (s *KeeperTestSuite) TestAllowedBidderByAuction() {
	auction := s.createFixedPriceAuction(
		s.addr(0),
		sdk.MustNewDecFromStr("1.0"),
		parseCoin("1000000000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		time.Now().AddDate(0, 6, 0),
		time.Now().AddDate(0, 6, 0).AddDate(0, 1, 0),
		true,
	)
	s.Require().Equal(auction.GetStatus(), types.AuctionStatusStandBy)

	allowedBidders := s.keeper.GetAllowedBiddersByAuction(s.ctx, auction.Id)
	s.Require().Len(allowedBidders, 0)

	// Add new allowed bidders
	newAllowedBidders := []types.AllowedBidder{
		{Bidder: s.addr(1).String(), MaxBidAmount: parseInt("100000")},
		{Bidder: s.addr(2).String(), MaxBidAmount: parseInt("100000")},
		{Bidder: s.addr(3).String(), MaxBidAmount: parseInt("100000")},
	}
	err := s.keeper.AddAllowedBidders(s.ctx, auction.Id, newAllowedBidders)
	s.Require().NoError(err)

	allowedBidders = s.keeper.GetAllowedBiddersByAuction(s.ctx, auction.Id)
	s.Require().Len(allowedBidders, 3)
}

func (s *KeeperTestSuite) TestLastBidId() {
	auction := s.createFixedPriceAuction(
		s.addr(0),
		sdk.OneDec(),
		parseCoin("500000000000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	bidId := s.keeper.GetLastBidId(s.ctx, auction.Id)
	s.Require().Equal(uint64(0), bidId)

	s.placeBidFixedPrice(auction.Id, s.addr(1), sdk.OneDec(), parseCoin("20000000denom2"), true)
	s.placeBidFixedPrice(auction.Id, s.addr(2), sdk.OneDec(), parseCoin("20000000denom2"), true)
	s.placeBidFixedPrice(auction.Id, s.addr(3), sdk.OneDec(), parseCoin("15000000denom2"), true)

	bidsById := s.keeper.GetBidsByAuctionId(s.ctx, auction.GetId())
	s.Require().Len(bidsById, 3)

	nextId := s.keeper.GetNextBidIdWithUpdate(s.ctx, auction.GetId())
	s.Require().Equal(uint64(4), nextId)

	// Create another auction
	auction2 := s.createFixedPriceAuction(
		s.addr(0),
		sdk.MustNewDecFromStr("0.5"),
		parseCoin("1000000000000denom3"),
		"denom4",
		[]types.VestingSchedule{},
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)

	// Bid id must start from 1 with new auction
	bidsById = s.keeper.GetBidsByAuctionId(s.ctx, auction2.GetId())
	s.Require().Len(bidsById, 0)

	nextId = s.keeper.GetNextBidIdWithUpdate(s.ctx, auction2.GetId())
	s.Require().Equal(uint64(1), nextId)
}

func (s *KeeperTestSuite) TestIterateBids() {
	startedAuction := s.createFixedPriceAuction(
		s.addr(0),
		sdk.OneDec(),
		parseCoin("500000000000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)

	auction, found := s.keeper.GetAuction(s.ctx, startedAuction.GetId())
	s.Require().True(found)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	s.placeBidFixedPrice(auction.GetId(), s.addr(1), sdk.OneDec(), parseCoin("20000000denom2"), true)
	s.placeBidFixedPrice(auction.GetId(), s.addr(2), sdk.OneDec(), parseCoin("20000000denom2"), true)
	s.placeBidFixedPrice(auction.GetId(), s.addr(2), sdk.OneDec(), parseCoin("15000000denom2"), true)
	s.placeBidFixedPrice(auction.GetId(), s.addr(3), sdk.OneDec(), parseCoin("35000000denom2"), true)

	bids := s.keeper.GetBids(s.ctx)
	s.Require().Len(bids, 4)

	bidsById := s.keeper.GetBidsByAuctionId(s.ctx, auction.GetId())
	s.Require().Len(bidsById, 4)

	bidsByBidder := s.keeper.GetBidsByBidder(s.ctx, s.addr(2))
	s.Require().Len(bidsByBidder, 2)
}

func (s *KeeperTestSuite) TestVestingQueue() {
	vestingQueue := types.NewVestingQueue(
		1, s.addr(1), parseCoin("100_000_000denom1"),
		types.MustParseRFC3339("2023-01-01T00:00:00Z"), false,
	)

	vq := s.keeper.GetVestingQueue(s.ctx, 1, vestingQueue.ReleaseTime)
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
			Weight:      sdk.MustNewDecFromStr("0.5"),
		},
		{
			ReleaseTime: types.MustParseRFC3339("2023-06-01T00:00:00Z"),
			Weight:      sdk.MustNewDecFromStr("0.5"),
		},
	} {
		payingAmt := reserveCoin.Amount.ToDec().MulTruncate(vs.Weight).TruncateInt()

		s.keeper.SetVestingQueue(s.ctx, types.VestingQueue{
			AuctionId:   uint64(1),
			Auctioneer:  s.addr(1).String(),
			PayingCoin:  sdk.NewCoin(payingCoinDenom, payingAmt),
			ReleaseTime: vs.ReleaseTime,
			Released:    false,
		})
	}

	// Set vesting schedules with 4 vesting queues
	for _, vs := range []types.VestingSchedule{
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
	} {
		payingAmt := reserveCoin.Amount.ToDec().MulTruncate(vs.Weight).TruncateInt()

		s.keeper.SetVestingQueue(s.ctx, types.VestingQueue{
			AuctionId:   uint64(2),
			Auctioneer:  s.addr(2).String(),
			PayingCoin:  sdk.NewCoin(payingCoinDenom, payingAmt),
			ReleaseTime: vs.ReleaseTime,
			Released:    false,
		})
	}

	s.Require().Len(s.keeper.GetVestingQueuesByAuctionId(s.ctx, uint64(1)), 2)
	s.Require().Len(s.keeper.GetVestingQueuesByAuctionId(s.ctx, uint64(2)), 4)
	s.Require().Len(s.keeper.GetVestingQueues(s.ctx), 6)

	totalPayingCoin := sdk.NewInt64Coin(payingCoinDenom, 0)
	for _, vq := range s.keeper.GetVestingQueuesByAuctionId(s.ctx, uint64(2)) {
		totalPayingCoin = totalPayingCoin.Add(vq.PayingCoin)
	}
	s.Require().Equal(reserveCoin, totalPayingCoin)
}
