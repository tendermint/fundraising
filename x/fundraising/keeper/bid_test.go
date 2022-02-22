package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/fundraising/x/fundraising/types"

	_ "github.com/stretchr/testify/suite"
)

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

	bidId := s.keeper.GetLastBidId(s.ctx, auction.GetId())
	s.Require().Equal(uint64(0), bidId)

	s.placeBid(auction.GetId(), s.addr(1), types.BidTypeFixedPrice, sdk.OneDec(), parseCoin("20000000denom2"), true)
	s.placeBid(auction.GetId(), s.addr(2), types.BidTypeFixedPrice, sdk.OneDec(), parseCoin("20000000denom2"), true)
	s.placeBid(auction.GetId(), s.addr(3), types.BidTypeFixedPrice, sdk.OneDec(), parseCoin("15000000denom2"), true)

	bidsById := s.keeper.GetBidsByAuctionId(s.ctx, auction.GetId())
	s.Require().Len(bidsById, 3)

	nextSeq := s.keeper.GetNextBidIdWithUpdate(s.ctx, auction.GetId())
	s.Require().Equal(uint64(4), nextSeq)

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

	nextSeq = s.keeper.GetNextBidIdWithUpdate(s.ctx, auction2.GetId())
	s.Require().Equal(uint64(1), nextSeq)
}

func (s *KeeperTestSuite) TestBidIterators() {
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

func (s *KeeperTestSuite) TestPlaceBid() {
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

	allowedBidder := s.addr(1)
	notAllowedBidder := s.addr(2)

	// Add allowed bidder to allowed bidder list
	err := s.keeper.AddAllowedBidders(s.ctx, auction.GetId(), []types.AllowedBidder{
		{Bidder: s.addr(1).String(), MaxBidAmount: sdk.NewInt(10_000_000)},
	})
	s.Require().NoError(err)

	// The bidder is not allowed
	s.fundAddr(notAllowedBidder, sdk.NewCoins(sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 5_000_000)))
	_, err = s.keeper.PlaceBid(s.ctx, &types.MsgPlaceBid{
		AuctionId: auction.GetId(),
		Bidder:    notAllowedBidder.String(),
		Price:     sdk.OneDec(),
		Coin:      sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 5_000_000),
	})
	s.Require().Error(err)

	// Maximum bid amount limit
	s.fundAddr(allowedBidder, sdk.NewCoins(sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 20_000_000)))
	_, err = s.keeper.PlaceBid(s.ctx, &types.MsgPlaceBid{
		AuctionId: auction.GetId(),
		Bidder:    allowedBidder.String(),
		Price:     sdk.OneDec(),
		Coin:      sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 20_000_000),
	})
	s.Require().Error(err)

	// Happy case
	s.fundAddr(allowedBidder, sdk.NewCoins(sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 10_000_000)))
	_, err = s.keeper.PlaceBid(s.ctx, &types.MsgPlaceBid{
		AuctionId: auction.GetId(),
		Bidder:    allowedBidder.String(),
		Price:     sdk.OneDec(),
		Coin:      sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 10_000_000),
	})
	s.Require().NoError(err)
}
