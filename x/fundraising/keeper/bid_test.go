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

	s.placeBidFixedPrice(auction.Id, s.addr(1), parseDec("1"), parseCoin("200000000denom2"), true)
	s.placeBidFixedPrice(auction.Id, s.addr(2), parseDec("1"), parseCoin("200000000denom2"), true)
	s.placeBidFixedPrice(auction.Id, s.addr(3), parseDec("1"), parseCoin("250000000denom2"), true)
	s.placeBidFixedPrice(auction.Id, s.addr(4), parseDec("1"), parseCoin("250000000denom2"), true)

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

	s.placeBidFixedPrice(auction.Id, s.addr(1), parseDec("1"), parseCoin("100000000denom2"), true)
	s.placeBidFixedPrice(auction.Id, s.addr(1), parseDec("1"), parseCoin("100000000denom2"), true)

	_, err := s.keeper.PlaceBid(s.ctx, &types.MsgPlaceBid{
		AuctionId: auction.Id,
		Bidder:    s.addr(1).String(),
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

	s.fundAddr(s.addr(1), parseCoins("200000000denom1, 200000000denom2"))
	s.addAllowedBidder(auction.Id, s.addr(1), parseCoin("200000000denom1").Amount)
	s.addAllowedBidder(auction.Id, s.addr(1), parseCoin("200000000denom2").Amount)

	// Place a BidTypeBatchWorth bid with an incorrect denom (SellingCoinDenom)
	_, err := s.keeper.PlaceBid(s.ctx, &types.MsgPlaceBid{
		AuctionId: auction.Id,
		Bidder:    s.addr(1).String(),
		BidType:   types.BidTypeBatchWorth,
		Price:     parseDec("1"),
		Coin:      parseCoin("100000000denom1"),
	})
	s.Require().ErrorIs(err, types.ErrIncorrectCoinDenom)

	// Place a BidTypeBatchMany bid with an incorrect denom (PayingCoinDenom)
	_, err = s.keeper.PlaceBid(s.ctx, &types.MsgPlaceBid{
		AuctionId: auction.Id,
		Bidder:    s.addr(1).String(),
		BidType:   types.BidTypeBatchMany,
		Price:     parseDec("1"),
		Coin:      parseCoin("100000000denom2"),
	})
	s.Require().ErrorIs(err, types.ErrIncorrectCoinDenom)

}

func (s *KeeperTestSuite) TestModifyBid_IncorrectCoinDenom() {
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

	s.fundAddr(s.addr(1), parseCoins("200000000denom1, 200000000denom2"))
	s.addAllowedBidder(auction.Id, s.addr(1), parseCoin("200000000denom1").Amount)
	s.addAllowedBidder(auction.Id, s.addr(1), parseCoin("200000000denom2").Amount)

	s.fundAddr(s.addr(2), parseCoins("200000000denom1, 200000000denom2"))
	s.addAllowedBidder(auction.Id, s.addr(2), parseCoin("200000000denom1").Amount)
	s.addAllowedBidder(auction.Id, s.addr(2), parseCoin("200000000denom2").Amount)

	bid1 := s.placeBidBatchWorth(auction.Id, s.addr(1), parseDec("1"), parseCoin("100000000denom2"), true)
	bid2 := s.placeBidBatchMany(auction.Id, s.addr(1), parseDec("1"), parseCoin("100000000denom1"), true)

	// Place a BidTypeBatchWorth bid with an incorrect denom (SellingCoinDenom)
	_, err := s.keeper.ModifyBid(s.ctx, &types.MsgModifyBid{
		AuctionId: auction.Id,
		Bidder:    s.addr(1).String(),
		BidId:     bid1.Id,
		Price:     parseDec("1"),
		Coin:      parseCoin("100000000denom1"),
	})
	s.Require().ErrorIs(err, types.ErrIncorrectCoinDenom)

	// Place a BidTypeBatchMany bid with an incorrect denom (PayingCoinDenom)
	_, err = s.keeper.ModifyBid(s.ctx, &types.MsgModifyBid{
		AuctionId: auction.Id,
		Bidder:    s.addr(2).String(),
		BidId:     bid2.Id,
		Price:     parseDec("1"),
		Coin:      parseCoin("100000000denom2"),
	})
	s.Require().ErrorIs(err, types.ErrIncorrectCoinDenom)
}

//func (s *KeeperTestSuite) TestModifyBid_IncorrectBidPrice() {
//	auction := s.createBatchAuction(
//		s.addr(1),
//		parseDec("0.5"),
//		parseCoin("1000000000denom1"),
//		"denom2",
//		[]types.VestingSchedule{},
//		1,
//		sdk.MustNewDecFromStr("0.2"),
//		time.Now().AddDate(0, 0, -1),
//		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
//		true,
//	)
//	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())
//
//	s.fundAddr(s.addr(1), parseCoins("200000000denom1, 200000000denom2"))
//	s.addAllowedBidder(auction.Id, s.addr(1), parseCoin("200000000denom1").Amount)
//	s.addAllowedBidder(auction.Id, s.addr(1), parseCoin("200000000denom2").Amount)
//
//	s.fundAddr(s.addr(2), parseCoins("200000000denom1, 200000000denom2"))
//	s.addAllowedBidder(auction.Id, s.addr(2), parseCoin("200000000denom1").Amount)
//	s.addAllowedBidder(auction.Id, s.addr(2), parseCoin("200000000denom2").Amount)
//
//	bid1 := s.placeBidBatchWorth(auction.Id, s.addr(1), parseDec("1"), parseCoin("100000000denom2"), true)
//	bid2 := s.placeBidBatchMany(auction.Id, s.addr(1), parseDec("1"), parseCoin("100000000denom1"), true)
//
//	// Place a BidTypeBatchWorth bid with an incorrect price
//	_, err := s.keeper.ModifyBid(s.ctx, &types.MsgModifyBid{
//		AuctionId: auction.Id,
//		Bidder:    s.addr(1).String(),
//		BidId:     bid1.Id,
//		Price:     parseDec("0.9"),
//		Coin:      parseCoin("100000000denom2"),
//	})
//	s.Require().ErrorIs(err, types.ErrInvalidModify)
//
//	// Place a BidTypeBatchMany bid with an incorrect price
//	_, err = s.keeper.ModifyBid(s.ctx, &types.MsgModifyBid{
//		AuctionId: auction.Id,
//		Bidder:    s.addr(2).String(),
//		BidId:     bid2.Id,
//		Price:     parseDec("0.9"),
//		Coin:      parseCoin("100000000denom1"),
//	})
//	s.Require().ErrorIs(err, types.ErrInvalidModify)
//}

//func (s *KeeperTestSuite) TestModifyBid_IncorrectCoinAmount() {
//	auction := s.createBatchAuction(
//		s.addr(1),
//		parseDec("0.5"),
//		parseCoin("1000000000denom1"),
//		"denom2",
//		[]types.VestingSchedule{},
//		1,
//		sdk.MustNewDecFromStr("0.2"),
//		time.Now().AddDate(0, 0, -1),
//		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
//		true,
//	)
//	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())
//
//	s.fundAddr(s.addr(1), parseCoins("200000000denom1, 200000000denom2"))
//	s.addAllowedBidder(auction.Id, s.addr(1), parseCoin("200000000denom1").Amount)
//	s.addAllowedBidder(auction.Id, s.addr(1), parseCoin("200000000denom2").Amount)
//
//	s.fundAddr(s.addr(2), parseCoins("200000000denom1, 200000000denom2"))
//	s.addAllowedBidder(auction.Id, s.addr(2), parseCoin("200000000denom1").Amount)
//	s.addAllowedBidder(auction.Id, s.addr(2), parseCoin("200000000denom2").Amount)
//
//	bid1 := s.placeBidBatchWorth(auction.Id, s.addr(1), parseDec("1"), parseCoin("100000000denom2"), true)
//	bid2 := s.placeBidBatchMany(auction.Id, s.addr(1), parseDec("1"), parseCoin("100000000denom1"), true)
//
//	// Place a BidTypeBatchWorth bid with an incorrect amount
//	_, err := s.keeper.ModifyBid(s.ctx, &types.MsgModifyBid{
//		AuctionId: auction.Id,
//		Bidder:    s.addr(1).String(),
//		BidId:     bid1.Id,
//		Price:     parseDec("1.0"),
//		Coin:      parseCoin("50000000denom2"),
//	})
//	s.Require().ErrorIs(err, types.ErrInvalidModify)
//
//	// Place a BidTypeBatchMany bid with an incorrect amount
//	_, err = s.keeper.ModifyBid(s.ctx, &types.MsgModifyBid{
//		AuctionId: auction.Id,
//		Bidder:    s.addr(2).String(),
//		BidId:     bid2.Id,
//		Price:     parseDec("1.0"),
//		Coin:      parseCoin("99999999denom1"),
//	})
//	s.Require().ErrorIs(err, types.ErrInvalidModify)
//}

//
//func (s *KeeperTestSuite) TestModifyBid_IncorrectBidId() {
//	auction := s.createBatchAuction(
//		s.addr(1),
//		parseDec("0.5"),
//		parseCoin("1000000000denom1"),
//		"denom2",
//		[]types.VestingSchedule{},
//		1,
//		sdk.MustNewDecFromStr("0.2"),
//		time.Now().AddDate(0, 0, -1),
//		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
//		true,
//	)
//	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())
//
//	s.fundAddr(s.addr(1), parseCoins("200000000denom1, 200000000denom2"))
//	s.addAllowedBidder(auction.Id, s.addr(1), parseCoin("200000000denom1").Amount)
//	s.addAllowedBidder(auction.Id, s.addr(1), parseCoin("200000000denom2").Amount)
//
//	s.fundAddr(s.addr(2), parseCoins("200000000denom1, 200000000denom2"))
//	s.addAllowedBidder(auction.Id, s.addr(2), parseCoin("200000000denom1").Amount)
//	s.addAllowedBidder(auction.Id, s.addr(2), parseCoin("200000000denom2").Amount)
//
//	bid1 := s.placeBidBatchWorth(auction.Id, s.addr(1), parseDec("1"), parseCoin("100000000denom2"), true)
//	bid2 := s.placeBidBatchMany(auction.Id, s.addr(1), parseDec("1"), parseCoin("100000000denom1"), true)
//
//	// Place a BidTypeBatchWorth bid with an incorrect bid id
//	_, err := s.keeper.ModifyBid(s.ctx, &types.MsgModifyBid{
//		AuctionId: auction.Id,
//		Bidder:    s.addr(1).String(),
//		BidId:     bid2.Id,
//		Price:     parseDec("1.0"),
//		Coin:      parseCoin("100000000denom2"),
//	})
//	s.Require().ErrorIs(err, types.ErrInvalidBidId)
//
//	// Place a BidTypeBatchMany bid with an incorrect bid id
//	_, err = s.keeper.ModifyBid(s.ctx, &types.MsgModifyBid{
//		AuctionId: auction.Id,
//		Bidder:    s.addr(2).String(),
//		BidId:     uint64(3),
//		Price:     parseDec("1.0"),
//		Coin:      parseCoin("100000000denom1"),
//	})
//	s.Require().ErrorIs(err, types.ErrInvalidBidId)
//}
