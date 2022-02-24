package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/fundraising/x/fundraising/types"

	_ "github.com/stretchr/testify/suite"
)

func (s *KeeperTestSuite) TestFixedPrice_InvalidStartPrice() {
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

	// Correct bid case
	s.addAllowedBidder(auction.Id, s.addr(1), exchangedSellingAmount(parseDec("1"), parseCoin("200000000denom2")))
	s.placeBid(auction.Id, s.addr(1), types.BidTypeFixedPrice, parseDec("1"), parseCoin("200000000denom2"), true)

	// The bid price must be the same as the start price of the auction.
	s.fundAddr(s.addr(2), parseCoins("200000000denom2"))
	s.addAllowedBidder(auction.Id, s.addr(2), exchangedSellingAmount(parseDec("1"), parseCoin("200000000denom2")))

	_, err := s.keeper.PlaceBid(s.ctx, &types.MsgPlaceBid{
		AuctionId: auction.Id,
		Bidder:    s.addr(2).String(),
		BidType:   types.BidTypeFixedPrice,
		Price:     parseDec("0.5"),
		Coin:      parseCoin("200000000denom2"),
	})
	s.Require().ErrorIs(err, types.ErrInvalidStartPrice)
}

func (s *KeeperTestSuite) TestFixedPrice_InsufficientRemainingAmount() {
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

	s.addAllowedBidder(auction.Id, s.addr(1), exchangedSellingAmount(parseDec("1"), parseCoin("200000000denom2")))
	s.addAllowedBidder(auction.Id, s.addr(2), exchangedSellingAmount(parseDec("1"), parseCoin("200000000denom2")))
	s.addAllowedBidder(auction.Id, s.addr(3), exchangedSellingAmount(parseDec("1"), parseCoin("250000000denom2")))
	s.addAllowedBidder(auction.Id, s.addr(4), exchangedSellingAmount(parseDec("1"), parseCoin("250000000denom2")))

	s.placeBid(auction.Id, s.addr(1), types.BidTypeFixedPrice, parseDec("1"), parseCoin("200000000denom2"), true)
	s.placeBid(auction.Id, s.addr(2), types.BidTypeFixedPrice, parseDec("1"), parseCoin("200000000denom2"), true)
	s.placeBid(auction.Id, s.addr(3), types.BidTypeFixedPrice, parseDec("1"), parseCoin("250000000denom2"), true)
	s.placeBid(auction.Id, s.addr(4), types.BidTypeFixedPrice, parseDec("1"), parseCoin("250000000denom2"), true)

	// The remaining coin amount must be insufficient
	s.fundAddr(s.addr(5), parseCoins("300000000denom2"))
	s.addAllowedBidder(auction.Id, s.addr(5), exchangedSellingAmount(parseDec("1"), parseCoin("300000000denom2")))

	_, err := s.keeper.PlaceBid(s.ctx, &types.MsgPlaceBid{
		AuctionId: auction.Id,
		Bidder:    s.addr(5).String(),
		BidType:   types.BidTypeFixedPrice,
		Price:     parseDec("1.0"),
		Coin:      parseCoin("300000000denom2"),
	})
	s.Require().ErrorIs(err, types.ErrInsufficientRemainingAmount)
}

func (s *KeeperTestSuite) TestFixedPrice_OverMaxBidAmountLimit() {
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

	s.addAllowedBidder(auction.Id, s.addr(1), exchangedSellingAmount(parseDec("1"), parseCoin("200000000denom2")))
	s.placeBid(auction.Id, s.addr(1), types.BidTypeFixedPrice, parseDec("1"), parseCoin("100000000denom2"), true)
	s.placeBid(auction.Id, s.addr(1), types.BidTypeFixedPrice, parseDec("1"), parseCoin("100000000denom2"), true)

	// The total amount of bids that a bidder places must be equal to or smaller than MaxBidAmount.
	s.fundAddr(s.addr(2), parseCoins("200000000denom2"))
	s.addAllowedBidder(auction.Id, s.addr(2), exchangedSellingAmount(parseDec("1"), parseCoin("200000000denom2")))
	s.placeBid(auction.Id, s.addr(2), types.BidTypeFixedPrice, parseDec("1"), parseCoin("100000000denom2"), true)

	_, err := s.keeper.PlaceBid(s.ctx, &types.MsgPlaceBid{
		AuctionId: auction.Id,
		Bidder:    s.addr(2).String(),
		BidType:   types.BidTypeFixedPrice,
		Price:     parseDec("1"),
		Coin:      parseCoin("100000001denom2"),
	})
	s.Require().ErrorIs(err, types.ErrOverMaxBidAmountLimit)
}

func (s *KeeperTestSuite) TestBatchAuction_IncorrectCoinDenom() {
	auction := s.createBatchAuction(
		s.addr(1),
		parseDec("0.5"),
		parseCoin("1000000000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		1,
		sdk.MustNewDecFromStr("0.2"),
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	s.fundAddr(s.addr(1), parseCoins("200000000denom2"))
	s.addAllowedBidder(auction.Id, s.addr(1), exchangedSellingAmount(parseDec("0.5"), parseCoin("200000000denom2")))

	// Place a BidTypeBatchWorth bid with an incorrect denom
	_, err := s.keeper.PlaceBid(s.ctx, &types.MsgPlaceBid{
		AuctionId: auction.Id,
		Bidder:    s.addr(1).String(),
		BidType:   types.BidTypeBatchWorth,
		Price:     parseDec("1"),
		Coin:      parseCoin("100000000denom1"),
	})
	s.Require().ErrorIs(err, types.ErrIncorrectCoinDenom)

	// Place a BidTypeBatchMany bid with an incorrect denom
	_, err = s.keeper.PlaceBid(s.ctx, &types.MsgPlaceBid{
		AuctionId: auction.Id,
		Bidder:    s.addr(1).String(),
		BidType:   types.BidTypeBatchMany,
		Price:     parseDec("1"),
		Coin:      parseCoin("100000000denom2"),
	})
	s.Require().ErrorIs(err, types.ErrIncorrectCoinDenom)

}

//func (s *KeeperTestSuite) TestBatchWorth {
//
//}
//
//func (s *KeeperTestSuite) TestBatchMany {
//
//}

func (s *KeeperTestSuite) TestModifyBid_IncorrectAuctionType() {
	// TODO: not implemented yet

}

func (s *KeeperTestSuite) TestModifyBid_IncorrectCoinDenom() {
	// TODO: not implemented yet

}

func (s *KeeperTestSuite) TestModifyBid_IncorrectBidPrice() {
	// TODO: not implemented yet
	// cover a case to modify a bid with higher price
}

func (s *KeeperTestSuite) TestModifyBid_IncorrectCoinAmount() {
	// TODO: not implemented yet
	// cover a case to modify a bid with higher coin amount
}

func (s *KeeperTestSuite) TestModifyBid_IncorrectBidId() {
	// TODO: not implemented yet
	// if Bid Id does not exist for the bidder
	// if Bid type for the bid Id is not the same

}

func (s *KeeperTestSuite) TestModifyBid_InsufficientRemainingAmount() {
	// TODO: not implemented yet

}
