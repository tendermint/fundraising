package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/fundraising/x/fundraising/types"

	_ "github.com/stretchr/testify/suite"
)

func (s *KeeperTestSuite) TestBidIterators() {
	startedAuction := s.createFixedPriceAuction(
		s.addr(0),
		sdk.OneDec(),
		sdk.NewInt64Coin("denom1", 1_000_000_000_000),
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

func (s *KeeperTestSuite) TestBidSequence() {
	auction := s.createFixedPriceAuction(
		s.addr(0),
		sdk.OneDec(),
		sdk.NewInt64Coin("denom1", 1_000_000_000_000),
		"denom2",
		[]types.VestingSchedule{},
		types.MustParseRFC3339("2022-01-01T00:00:00Z"),
		types.MustParseRFC3339("2022-03-10T00:00:00Z"),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	sequence := s.keeper.GetLastSequence(s.ctx, auction.GetId())
	s.Require().Equal(uint64(0), sequence)

	s.placeBid(auction.GetId(), s.addr(1), sdk.OneDec(), sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 20_000_000), true)
	s.placeBid(auction.GetId(), s.addr(2), sdk.OneDec(), sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 20_000_000), true)
	s.placeBid(auction.GetId(), s.addr(3), sdk.OneDec(), sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 15_000_000), true)

	bidsById := s.keeper.GetBidsByAuctionId(s.ctx, auction.GetId())
	s.Require().Len(bidsById, 3)

	nextSeq := s.keeper.GetNextSequenceWithUpdate(s.ctx, auction.GetId())
	s.Require().Equal(uint64(4), nextSeq)

	// Create another auction
	auction2 := s.createFixedPriceAuction(
		s.addr(0),
		sdk.MustNewDecFromStr("0.5"),
		sdk.NewInt64Coin("denom3", 500_000_000_000),
		"denom3",
		[]types.VestingSchedule{},
		types.MustParseRFC3339("2022-01-10T00:00:00Z"),
		types.MustParseRFC3339("2022-05-10T00:00:00Z"),
		true,
	)

	// Sequence must start from 1 with new auction
	bidsById = s.keeper.GetBidsByAuctionId(s.ctx, auction2.GetId())
	s.Require().Len(bidsById, 0)

	nextSeq = s.keeper.GetNextSequenceWithUpdate(s.ctx, auction2.GetId())
	s.Require().Equal(uint64(1), nextSeq)
}

func (s *KeeperTestSuite) TestCalculateWinners() {
	auction := s.createEnglishAuction(
		s.addr(0),
		parseDec("0.5"),
		sdk.NewInt64Coin("denom1", 1_000_000_000_000),
		"denom2",
		[]types.VestingSchedule{},
		parseDec("1.0"),
		uint32(1),
		parseDec("0.1"),
		types.MustParseRFC3339("2022-01-01T00:00:00Z"),
		types.MustParseRFC3339("2022-03-10T00:00:00Z"),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	//
	s.placeBid(auction.Id, s.addr(1), parseDec("0.6"), sdk.NewInt64Coin(auction.PayingCoinDenom, 20_000_000), true)
	s.placeBid(auction.Id, s.addr(3), parseDec("0.5"), sdk.NewInt64Coin(auction.PayingCoinDenom, 35_000_000), true)
	s.placeBid(auction.Id, s.addr(2), parseDec("0.7"), sdk.NewInt64Coin(auction.PayingCoinDenom, 15_000_000), true)
	s.placeBid(auction.Id, s.addr(2), parseDec("0.8"), sdk.NewInt64Coin(auction.PayingCoinDenom, 20_000_000), true)

	bids := s.keeper.GetBidsByAuctionId(s.ctx, auction.GetId())
	bids = types.SanitizeReverseBids(bids)

	for _, b := range bids {
		fmt.Println("b: ", b)
	}

	// totalSellingAmt := sdk.ZeroDec()
	// totalCoinAmt := sdk.ZeroDec()

	// for _, bid := range bids {

	// 	totalCoinAmt = totalCoinAmt.Add(bid.Coin.Amount.ToDec())
	// 	totalSellingAmt = totalCoinAmt.QuoTruncate(bid.Price)

	// 	fmt.Println("Coin Amount: ", totalCoinAmt)
	// 	fmt.Println("Selling Amount: ", totalSellingAmt)
	// 	fmt.Println("")
	// }

	// sellingCoinDenom := auction.SellingCoin.Denom
	// remainingCoin := sdk.NewInt64Coin(sellingCoinDenom, 100_000_000)
	// remainingCoin = remainingCoin.Sub(sdk.NewCoin(sellingCoinDenom, totalSellingAmt.TruncateInt()))
	// fmt.Println("remainingCoin: ", remainingCoin)
}
