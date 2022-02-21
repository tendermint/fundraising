package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/fundraising/x/fundraising/types"

	_ "github.com/stretchr/testify/suite"
)

func (s *KeeperTestSuite) TestBidIterators() {
	startedAuction := s.createFixedPriceAuction(
		s.addr(0),
		sdk.OneDec(),
		sdk.NewInt64Coin("denom1", 500_000_000_000),
		"denom2",
		[]types.VestingSchedule{},
		types.MustParseRFC3339("2022-01-01T00:00:00Z"),
		types.MustParseRFC3339("2022-03-10T00:00:00Z"),
		true,
	)

	auction, found := s.keeper.GetAuction(s.ctx, startedAuction.GetId())
	s.Require().True(found)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	s.placeBid(auction.GetId(), s.addr(1), sdk.OneDec(), sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 20_000_000), true)
	s.placeBid(auction.GetId(), s.addr(2), sdk.OneDec(), sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 20_000_000), true)
	s.placeBid(auction.GetId(), s.addr(2), sdk.OneDec(), sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 15_000_000), true)
	s.placeBid(auction.GetId(), s.addr(3), sdk.OneDec(), sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 35_000_000), true)

	bids := s.keeper.GetBids(s.ctx)
	s.Require().Len(bids, 4)

	bidsById := s.keeper.GetBidsByAuctionId(s.ctx, auction.GetId())
	s.Require().Len(bidsById, 4)

	bidsByBidder := s.keeper.GetBidsByBidder(s.ctx, s.addr(2))
	s.Require().Len(bidsByBidder, 2)
}

func (s *KeeperTestSuite) TestBidId() {
	auction := s.createFixedPriceAuction(
		s.addr(0),
		sdk.OneDec(),
		sdk.NewInt64Coin("denom1", 500_000_000_000),
		"denom2",
		[]types.VestingSchedule{},
		types.MustParseRFC3339("2022-01-01T00:00:00Z"),
		types.MustParseRFC3339("2022-03-10T00:00:00Z"),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	bidId := s.keeper.GetLastBidId(s.ctx, auction.GetId())
	s.Require().Equal(uint64(0), bidId)

	s.placeBid(auction.GetId(), s.addr(1), sdk.OneDec(), sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 20_000_000), true)
	s.placeBid(auction.GetId(), s.addr(2), sdk.OneDec(), sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 20_000_000), true)
	s.placeBid(auction.GetId(), s.addr(3), sdk.OneDec(), sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 15_000_000), true)

	bidsById := s.keeper.GetBidsByAuctionId(s.ctx, auction.GetId())
	s.Require().Len(bidsById, 3)

	nextSeq := s.keeper.GetNextBidIdWithUpdate(s.ctx, auction.GetId())
	s.Require().Equal(uint64(4), nextSeq)

	// Create another auction
	auction2 := s.createFixedPriceAuction(
		s.addr(0),
		sdk.MustNewDecFromStr("0.5"),
		sdk.NewInt64Coin("denom3", 500_000_000_000),
		"denom3",
		[]types.VestingSchedule{},
		types.MustParseRFC3339("2022-01-10T00:00:00Z"),
		types.MustParseRFC3339("2022-12-10T00:00:00Z"),
		true,
	)

	// Bid id must start from 1 with new auction
	bidsById = s.keeper.GetBidsByAuctionId(s.ctx, auction2.GetId())
	s.Require().Len(bidsById, 0)

	nextSeq = s.keeper.GetNextBidIdWithUpdate(s.ctx, auction2.GetId())
	s.Require().Equal(uint64(1), nextSeq)
}

func (s *KeeperTestSuite) TestPlaceBid() {
	startedAuction := s.createFixedPriceAuction(
		s.addr(0),
		sdk.OneDec(),
		sdk.NewInt64Coin("denom1", 500_000_000_000),
		"denom2",
		[]types.VestingSchedule{},
		types.MustParseRFC3339("2022-01-01T00:00:00Z"),
		types.MustParseRFC3339("2022-06-10T00:00:00Z"),
		true,
	)
	auction, found := s.keeper.GetAuction(s.ctx, startedAuction.GetId())
	s.Require().True(found)

	allowedBidder := s.addr(1)
	notAllowedBidder := s.addr(2)

	// Add allowed bidder to allowed bidder list
	err := s.keeper.AddAllowedBidders(s.ctx, auction.GetId(), []*types.AllowedBidder{
		{Bidder: s.addr(1).String(), MaxBidAmount: sdk.NewInt(10_000_000)},
	})
	s.Require().NoError(err)

	// The bidder is not allowed
	s.fundAddr(notAllowedBidder, sdk.NewCoins(sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 5_000_000)))
	_, err = s.keeper.PlaceBid(s.ctx, &types.MsgPlaceBid{
		AuctionId: auction.GetId(),
		Bidder:    notAllowedBidder.String(),
		BidPrice:  sdk.OneDec(),
		BidCoin:   sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 5_000_000),
	})
	s.Require().Error(err)

	// Maximum bid amount limit
	s.fundAddr(allowedBidder, sdk.NewCoins(sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 20_000_000)))
	_, err = s.keeper.PlaceBid(s.ctx, &types.MsgPlaceBid{
		AuctionId: auction.GetId(),
		Bidder:    allowedBidder.String(),
		BidPrice:  sdk.OneDec(),
		BidCoin:   sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 20_000_000),
	})
	s.Require().Error(err)

	// Happy case
	s.fundAddr(allowedBidder, sdk.NewCoins(sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 10_000_000)))
	_, err = s.keeper.PlaceBid(s.ctx, &types.MsgPlaceBid{
		AuctionId: auction.GetId(),
		Bidder:    allowedBidder.String(),
		BidPrice:  sdk.OneDec(),
		BidCoin:   sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 10_000_000),
	})
	s.Require().NoError(err)
}
