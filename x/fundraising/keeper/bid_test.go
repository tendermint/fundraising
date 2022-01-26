package keeper_test

// import (
// 	"fmt"
// 	"sort"

// 	sdk "github.com/cosmos/cosmos-sdk/types"

// 	"github.com/tendermint/fundraising/x/fundraising/types"

// 	_ "github.com/stretchr/testify/suite"
// )

// func (s *KeeperTestSuite) TestBidIterators() {
// 	startedAuction := s.createFixedPriceAuction(
// 		s.addr(0),
// 		sdk.OneDec(),
// 		sdk.NewInt64Coin("denom1", 500_000_000_000),
// 		"denom2",
// 		[]types.VestingSchedule{},
// 		types.MustParseRFC3339("2022-01-01T00:00:00Z"),
// 		types.MustParseRFC3339("2022-03-10T00:00:00Z"),
// 		true,
// 	)

// 	auction, found := s.keeper.GetAuction(s.ctx, startedAuction.GetId())
// 	s.Require().True(found)
// 	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

// 	s.placeBid(auction.GetId(), s.addr(1), sdk.OneDec(), sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 20_000_000), true)
// 	s.placeBid(auction.GetId(), s.addr(2), sdk.OneDec(), sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 20_000_000), true)
// 	s.placeBid(auction.GetId(), s.addr(2), sdk.OneDec(), sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 15_000_000), true)
// 	s.placeBid(auction.GetId(), s.addr(3), sdk.OneDec(), sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 35_000_000), true)

// 	bids := s.keeper.GetBids(s.ctx)
// 	s.Require().Len(bids, 4)

// 	bidsById := s.keeper.GetBidsByAuctionId(s.ctx, auction.GetId())
// 	s.Require().Len(bidsById, 4)

// 	bidsByBidder := s.keeper.GetBidsByBidder(s.ctx, s.addr(2))
// 	s.Require().Len(bidsByBidder, 2)
// }

// func (s *KeeperTestSuite) TestBidSequence() {
// 	auction := s.createFixedPriceAuction(
// 		s.addr(0),
// 		sdk.OneDec(),
// 		sdk.NewInt64Coin("denom1", 500_000_000_000),
// 		"denom2",
// 		[]types.VestingSchedule{},
// 		types.MustParseRFC3339("2022-01-01T00:00:00Z"),
// 		types.MustParseRFC3339("2022-03-10T00:00:00Z"),
// 		true,
// 	)
// 	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

// 	sequence := s.keeper.GetLastSequence(s.ctx, auction.GetId())
// 	s.Require().Equal(uint64(0), sequence)

// 	s.placeBid(auction.GetId(), s.addr(1), sdk.OneDec(), sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 20_000_000), true)
// 	s.placeBid(auction.GetId(), s.addr(2), sdk.OneDec(), sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 20_000_000), true)
// 	s.placeBid(auction.GetId(), s.addr(3), sdk.OneDec(), sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 15_000_000), true)

// 	bidsById := s.keeper.GetBidsByAuctionId(s.ctx, auction.GetId())
// 	s.Require().Len(bidsById, 3)

// 	nextSeq := s.keeper.GetNextSequenceWithUpdate(s.ctx, auction.GetId())
// 	s.Require().Equal(uint64(4), nextSeq)

// 	// Create another auction
// 	auction2 := s.createFixedPriceAuction(
// 		s.addr(0),
// 		sdk.MustNewDecFromStr("0.5"),
// 		sdk.NewInt64Coin("denom3", 500_000_000_000),
// 		"denom3",
// 		[]types.VestingSchedule{},
// 		types.MustParseRFC3339("2022-01-10T00:00:00Z"),
// 		types.MustParseRFC3339("2022-02-10T00:00:00Z"),
// 		true,
// 	)

// 	// Sequence must start from 1 with new auction
// 	bidsById = s.keeper.GetBidsByAuctionId(s.ctx, auction2.GetId())
// 	s.Require().Len(bidsById, 0)

// 	nextSeq = s.keeper.GetNextSequenceWithUpdate(s.ctx, auction2.GetId())
// 	s.Require().Equal(uint64(1), nextSeq)
// }

// func (suite *KeeperTestSuite) TestCalculateWinners() {
// 	sellingCoinDenom := denom2
// 	payingCoinDenom := denom1
// 	remainingCoin := sdk.NewInt64Coin(sellingCoinDenom, 100_000_000)

// 	bids := []types.Bid{
// 		{
// 			AuctionId: 1,
// 			Sequence:  1,
// 			Bidder:    suite.addrs[0].String(),
// 			Price:     sdk.MustNewDecFromStr("0.85"),
// 			Coin:      sdk.NewInt64Coin(payingCoinDenom, 1_000_000),
// 		},
// 		{
// 			AuctionId: 1,
// 			Sequence:  2,
// 			Bidder:    suite.addrs[1].String(),
// 			Price:     sdk.MustNewDecFromStr("1.0"),
// 			Coin:      sdk.NewInt64Coin(payingCoinDenom, 1_000_000),
// 		},
// 		{
// 			AuctionId: 1,
// 			Sequence:  3,
// 			Bidder:    suite.addrs[2].String(),
// 			Price:     sdk.MustNewDecFromStr("0.95"),
// 			Coin:      sdk.NewInt64Coin(payingCoinDenom, 1_000_000),
// 		},
// 		{
// 			AuctionId: 1,
// 			Sequence:  4,
// 			Bidder:    suite.addrs[3].String(),
// 			Price:     sdk.MustNewDecFromStr("0.7"),
// 			Coin:      sdk.NewInt64Coin(payingCoinDenom, 1_000_000),
// 		},
// 	}

// 	// Sort in descending order
// 	sort.SliceStable(bids, func(i, j int) bool {
// 		return bids[i].Price.GTE(bids[j].Price)
// 	})

// 	totalSellingAmt := sdk.ZeroDec()
// 	totalCoinAmt := sdk.ZeroDec()

// 	for _, bid := range bids {

// 		totalCoinAmt = totalCoinAmt.Add(bid.Coin.Amount.ToDec())
// 		totalSellingAmt = totalCoinAmt.QuoTruncate(bid.Price)

// 		fmt.Println("Coin Amount: ", totalCoinAmt)
// 		fmt.Println("Selling Amount: ", totalSellingAmt)
// 		fmt.Println("")
// 	}

// 	remainingCoin = remainingCoin.Sub(sdk.NewCoin(sellingCoinDenom, totalSellingAmt.TruncateInt()))
// 	fmt.Println("remainingCoin: ", remainingCoin)
// }
