package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/fundraising/x/fundraising/types"

	_ "github.com/stretchr/testify/suite"
)

func (s *KeeperTestSuite) TestFixedPriceAuction() {
	startedAuction := s.createFixedPriceAuction(
		s.addr(0),
		parseDec("0.5"),
		parseCoin("1000000000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)

	auction, found := s.keeper.GetAuction(s.ctx, startedAuction.GetId())
	s.Require().True(found)

	s.addAllowedBidder(auction.GetId(), s.addr(1), exchangeToSellingAmount(parseDec("0.5"), parseCoin("200000000denom2")))
	s.placeBid(auction.GetId(), s.addr(1), types.BidTypeFixedPrice, parseDec("0.5"), parseCoin("200000000denom2"), true)
}

func (s *KeeperTestSuite) TestFixedPriceAuction_InvalidStartPrice() {
	// TODO: not implemented yet
}

func (s *KeeperTestSuite) TestFixedPriceAuction_InsufficientRemainingAmount() {
	auction := s.createFixedPriceAuction(
		s.addr(0),
		parseDec("1"),
		parseCoin("1000000000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	s.addAllowedBidder(auction.Id, s.addr(1), exchangeToSellingAmount(parseDec("1"), parseCoin("200000000denom2")))
	s.addAllowedBidder(auction.Id, s.addr(2), exchangeToSellingAmount(parseDec("1"), parseCoin("200000000denom2")))
	s.addAllowedBidder(auction.Id, s.addr(3), exchangeToSellingAmount(parseDec("1"), parseCoin("250000000denom2")))
	s.addAllowedBidder(auction.Id, s.addr(4), exchangeToSellingAmount(parseDec("1"), parseCoin("250000000denom2")))

	s.placeBid(auction.Id, s.addr(1), types.BidTypeFixedPrice, parseDec("1"), parseCoin("200000000denom2"), true)
	s.placeBid(auction.Id, s.addr(2), types.BidTypeFixedPrice, parseDec("1"), parseCoin("200000000denom2"), true)
	s.placeBid(auction.Id, s.addr(3), types.BidTypeFixedPrice, parseDec("1"), parseCoin("250000000denom2"), true)
	s.placeBid(auction.Id, s.addr(4), types.BidTypeFixedPrice, parseDec("1"), parseCoin("250000000denom2"), true)

	// The remaining coin amount must be insufficient
	s.fundAddr(s.addr(5), parseCoins("300000000denom2"))
	s.addAllowedBidder(auction.Id, s.addr(5), exchangeToSellingAmount(parseDec("1"), parseCoin("300000000denom2")))

	_, err := s.keeper.PlaceBid(s.ctx, &types.MsgPlaceBid{
		AuctionId: auction.Id,
		Bidder:    s.addr(5).String(),
		BidType:   types.BidTypeFixedPrice,
		Price:     parseDec("1.0"),
		Coin:      parseCoin("300000000denom2"),
	})
	s.Require().ErrorIs(err, types.ErrInsufficientRemainingAmount)
}

func (s *KeeperTestSuite) TestFixedPriceAuction_OverMaxBidAmountLimit() {
	// TODO: not implemented yet
}

func (s *KeeperTestSuite) TestBatchAuction() {
	// TODO: not implemented yet
}

func (s *KeeperTestSuite) TestModifyBid() {
	// TODO: not implemented yet
	// cover a case to modify a bid with higher price
	// cover a case to modify a bid with higher coin amount
}

func (s *KeeperTestSuite) TestHandleBatchWorthBid() {
	// TODO: not implemented yet
}

func (s *KeeperTestSuite) TestHandleBatchManyBid() {
	// TODO: not done yet : to do
	auction := s.createBatchAuction(
		s.addr(1),
		parseDec("0.5"),
		parseCoin("5000000000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		1,
		sdk.MustNewDecFromStr("0.2"),
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	s.addAllowedBidder(auction.Id, s.addr(1), sdk.NewInt(2_000_000_000))

	s.placeBid(auction.Id, s.addr(1), types.BidTypeBatchMany, parseDec("0.2"), parseCoin("100000000denom1"), true)
	s.placeBid(auction.Id, s.addr(1), types.BidTypeBatchMany, parseDec("0.3"), parseCoin("200000000denom1"), true)
	s.placeBid(auction.Id, s.addr(1), types.BidTypeBatchMany, parseDec("0.5"), parseCoin("500000000denom1"), true)
	s.placeBid(auction.Id, s.addr(1), types.BidTypeBatchMany, parseDec("1.0"), parseCoin("500000000denom1"), true)
}
