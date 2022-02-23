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

	s.addAllowedBidder(auction.Id, s.addr(1), exchangeToSellingAmount(sdk.OneDec(), parseCoin("20000000denom2")))
	s.addAllowedBidder(auction.Id, s.addr(2), exchangeToSellingAmount(sdk.OneDec(), parseCoin("20000000denom2")))
	s.addAllowedBidder(auction.Id, s.addr(3), exchangeToSellingAmount(sdk.OneDec(), parseCoin("15000000denom2")))

	s.placeBid(auction.Id, s.addr(1), types.BidTypeFixedPrice, sdk.OneDec(), parseCoin("20000000denom2"), true)
	s.placeBid(auction.Id, s.addr(2), types.BidTypeFixedPrice, sdk.OneDec(), parseCoin("20000000denom2"), true)
	s.placeBid(auction.Id, s.addr(3), types.BidTypeFixedPrice, sdk.OneDec(), parseCoin("15000000denom2"), true)

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

	s.addAllowedBidder(auction.GetId(), s.addr(1), exchangeToSellingAmount(sdk.OneDec(), parseCoin("20000000denom2")))
	s.addAllowedBidder(auction.GetId(), s.addr(2), exchangeToSellingAmount(sdk.OneDec(), parseCoin("35000000denom2")))
	s.addAllowedBidder(auction.GetId(), s.addr(3), exchangeToSellingAmount(sdk.OneDec(), parseCoin("35000000denom2")))

	s.placeBid(auction.GetId(), s.addr(1), types.BidTypeFixedPrice, sdk.OneDec(), parseCoin("20000000denom2"), true)
	s.placeBid(auction.GetId(), s.addr(2), types.BidTypeFixedPrice, sdk.OneDec(), parseCoin("20000000denom2"), true)
	s.placeBid(auction.GetId(), s.addr(2), types.BidTypeFixedPrice, sdk.OneDec(), parseCoin("15000000denom2"), true)
	s.placeBid(auction.GetId(), s.addr(3), types.BidTypeFixedPrice, sdk.OneDec(), parseCoin("35000000denom2"), true)

	bids := s.keeper.GetBids(s.ctx)
	s.Require().Len(bids, 4)

	bidsById := s.keeper.GetBidsByAuctionId(s.ctx, auction.GetId())
	s.Require().Len(bidsById, 4)

	bidsByBidder := s.keeper.GetBidsByBidder(s.ctx, s.addr(2))
	s.Require().Len(bidsByBidder, 2)
}
