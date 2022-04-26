package keeper_test

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/fundraising/x/fundraising/types"

	_ "github.com/stretchr/testify/suite"
)

func (s *KeeperTestSuite) TestExampleFullString() {
	auction := s.createFixedPriceAuction(
		s.addr(0),
		parseDec("0.5"),
		parseCoin("1_000_000_000_000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)

	a, found := s.keeper.GetAuction(s.ctx, auction.GetId())
	s.Require().True(found)

	s.placeBidFixedPrice(a.GetId(), s.addr(1), a.GetStartPrice(), parseCoin("15_000_000denom2"), true)
	s.placeBidFixedPrice(a.GetId(), s.addr(2), a.GetStartPrice(), parseCoin("20_000_000denom2"), true)
	s.placeBidFixedPrice(a.GetId(), s.addr(4), a.GetStartPrice(), parseCoin("10_000_000denom1"), true)
	s.placeBidFixedPrice(a.GetId(), s.addr(6), a.GetStartPrice(), parseCoin("20_000_000denom1"), true)

	mInfo := s.keeper.CalculateFixedPriceAllocation(s.ctx, a)
	fmt.Println(s.fullString(a.GetId(), mInfo))

	// Output:
	// [Bids]
	// +--------------------bidder---------------------+-id-+---------price---------+---------type---------+-----reserve-amount-----+-------bid-amount-------+
	// | cosmos1qgqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqggwm7m |  1 |  0.500000000000000000 | BID_TYPE_FIXED_PRICE |               15000000 |               30000000 |
	// | cosmos1qsqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqv4uhu3 |  2 |  0.500000000000000000 | BID_TYPE_FIXED_PRICE |               20000000 |               40000000 |
	// | cosmos1pqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqvzjng6 |  3 |  0.500000000000000000 | BID_TYPE_FIXED_PRICE |                5000000 |               10000000 |
	// | cosmos1psqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqn5wmnk |  4 |  0.500000000000000000 | BID_TYPE_FIXED_PRICE |               10000000 |               20000000 |
	// +-----------------------------------------------+----+-----------------------+----------------------+------------------------+------------------------+

	// [Allocation]
	// +--------------------bidder---------------------+------allocated-amount------+
	// | cosmos1qgqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqggwm7m |                   30000000 |
	// | cosmos1qsqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqv4uhu3 |                   40000000 |
	// | cosmos1pqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqvzjng6 |                   10000000 |
	// | cosmos1psqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqn5wmnk |                   20000000 |
	// +-----------------------------------------------+----------------------------+

	// [MatchingInfo]
	// +-matched-len-+------matched-price------+------total-matched-amount------+
	// |           4 |    0.500000000000000000 |                      100000000 |
	// +-------------+-------------------------+--------------------------------+
}

func (s *KeeperTestSuite) TestBatchAuction_Many() {
	auction := s.createBatchAuction(
		s.addr(0),
		parseDec("1"),
		parseDec("0.1"),
		parseCoin("1_000_000_000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		1,
		sdk.MustNewDecFromStr("0.2"),
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	s.placeBidBatchMany(auction.Id, s.addr(1), parseDec("1"), parseCoin("500_000_000denom1"), sdk.NewInt(1_000_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(2), parseDec("0.9"), parseCoin("500_000_000denom1"), sdk.NewInt(1_000_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(3), parseDec("0.8"), parseCoin("500_000_000denom1"), sdk.NewInt(1_000_000_000), true)

	a, found := s.keeper.GetAuction(s.ctx, auction.Id)
	s.Require().True(found)

	mInfo := s.keeper.CalculateBatchAllocation(s.ctx, a)

	s.Require().Equal(mInfo.MatchedLen, int64(2))
	s.Require().Equal(mInfo.MatchedPrice, parseDec("0.9"))
	s.Require().Equal(mInfo.TotalMatchedAmount, sdk.NewInt(1_000_000_000))
	s.Require().Equal(mInfo.AllocationMap[s.addr(1).String()], sdk.NewInt(500_000_000))
	s.Require().Equal(mInfo.AllocationMap[s.addr(2).String()], sdk.NewInt(500_000_000))
	s.Require().Equal(mInfo.AllocationMap[s.addr(3).String()], sdk.NewInt(0))
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(1).String()], sdk.NewInt(450_000_000))
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(2).String()], sdk.NewInt(450_000_000))
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(3).String()], sdk.NewInt(0))
	s.Require().Equal(mInfo.RefundMap[s.addr(1).String()], sdk.NewInt(50_000_000))
	s.Require().Equal(mInfo.RefundMap[s.addr(3).String()], sdk.NewInt(400_000_000))
	s.Require().True(mInfo.RefundMap[s.addr(2).String()].Equal(sdk.NewInt(0)))

	err := s.keeper.AllocateSellingCoin(s.ctx, auction, mInfo)
	s.Require().NoError(err)

	sellingReserveAmt := s.getBalance(auction.GetSellingReserveAddress(), auction.SellingCoin.Denom).Amount
	remainingAmt := auction.GetSellingCoin().Amount.Sub(mInfo.TotalMatchedAmount)
	s.Require().True(sellingReserveAmt.Equal(remainingAmt))

	err = s.keeper.RefundRemainingSellingCoin(s.ctx, auction)
	s.Require().NoError(err)
	s.Require().True(s.getBalance(auction.GetSellingReserveAddress(), auction.SellingCoin.Denom).IsZero())

	// The auctioneer must have sellingCoin.Amount - TotalMatchedAmount
	s.Require().Equal(s.getBalance(s.addr(0), auction.GetSellingCoin().Denom).Amount, sdk.NewInt(0))

	// The bidders must have the matched selling coin
	s.Require().Equal(s.getBalance(s.addr(1), auction.GetSellingCoin().Denom).Amount, sdk.NewInt(500_000_000))
	s.Require().Equal(s.getBalance(s.addr(2), auction.GetSellingCoin().Denom).Amount, sdk.NewInt(500_000_000))
	s.Require().Equal(s.getBalance(s.addr(3), auction.GetSellingCoin().Denom).Amount, sdk.NewInt(0))

	// s.Require().True(s.getBalance(s.addr(3), auction.GetSellingCoin().Denom).IsZero())

	// Refund payingCoin
	err = s.keeper.RefundPayingCoin(s.ctx, auction, mInfo)
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) TestBatchAuction_Worth() {
	auction := s.createBatchAuction(
		s.addr(0),
		parseDec("1"),
		parseDec("0.1"),
		parseCoin("1_500_000_000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		1,
		sdk.MustNewDecFromStr("0.2"),
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	s.placeBidBatchWorth(auction.Id, s.addr(1), parseDec("1"), parseCoin("500_000_000denom2"), sdk.NewInt(1500_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(2), parseDec("0.9"), parseCoin("500_000_000denom2"), sdk.NewInt(1500_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(3), parseDec("0.8"), parseCoin("500_000_000denom2"), sdk.NewInt(1500_000_000), true)

	a, found := s.keeper.GetAuction(s.ctx, auction.Id)
	s.Require().True(found)

	mInfo := s.keeper.CalculateBatchAllocation(s.ctx, a)

	// Checking
	s.Require().Equal(int64(2), mInfo.MatchedLen)
	s.Require().Equal(parseDec("0.9"), mInfo.MatchedPrice)
	matchingPrice := parseDec("0.9")
	matchedAmt := sdk.NewInt(500_000_000).ToDec().QuoTruncate(matchingPrice).TruncateInt()

	s.Require().Equal(mInfo.TotalMatchedAmount, matchedAmt.Add(matchedAmt))
	s.Require().Equal(mInfo.AllocationMap[s.addr(1).String()], matchedAmt)
	s.Require().Equal(mInfo.AllocationMap[s.addr(2).String()], matchedAmt)
	s.Require().Equal(mInfo.AllocationMap[s.addr(3).String()], sdk.NewInt(0))
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(1).String()], sdk.NewInt(500_000_000))
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(2).String()], sdk.NewInt(500_000_000))
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(3).String()], sdk.NewInt(0))
	s.Require().True(mInfo.RefundMap[s.addr(1).String()].IsZero())
	s.Require().True(mInfo.RefundMap[s.addr(2).String()].IsZero())
	s.Require().Equal(mInfo.RefundMap[s.addr(3).String()], sdk.NewInt(500_000_000))

	// Distribute selling coin
	err := s.keeper.AllocateSellingCoin(s.ctx, auction, mInfo)
	s.Require().NoError(err)

	err = s.keeper.RefundRemainingSellingCoin(s.ctx, auction)
	s.Require().NoError(err)

	// The selling reserve account balance must be zero
	s.Require().True(s.getBalance(auction.GetSellingReserveAddress(), auction.SellingCoin.Denom).IsZero())

	// The auctioneer must have sellingCoin.Amount - TotalMatchedAmount
	s.Require().True(
		s.getBalance(s.addr(0), auction.GetSellingCoin().Denom).Amount.
			Equal(auction.SellingCoin.Amount.Sub(mInfo.TotalMatchedAmount)),
	)

	// The bidders must have the matched selling coin
	s.Require().Equal(s.getBalance(s.addr(1), auction.GetSellingCoin().Denom).Amount, matchedAmt)
	s.Require().Equal(s.getBalance(s.addr(2), auction.GetSellingCoin().Denom).Amount, matchedAmt)
	s.Require().True(s.getBalance(s.addr(3), auction.GetSellingCoin().Denom).IsZero())

	// Refund payingCoin
	err = s.keeper.RefundPayingCoin(s.ctx, auction, mInfo)
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) TestCalculateAllocation_Mixed() {
	auction := s.createBatchAuction(
		s.addr(0),
		parseDec("1"),
		parseDec("0.1"),
		parseCoin("1700_000_000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		1,
		sdk.MustNewDecFromStr("0.2"),
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	s.placeBidBatchMany(auction.Id, s.addr(1), parseDec("1"), parseCoin("500_000_000denom1"), sdk.NewInt(1500_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(2), parseDec("0.9"), parseCoin("500_000_000denom2"), sdk.NewInt(1500_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(3), parseDec("0.8"), parseCoin("500_000_000denom2"), sdk.NewInt(1500_000_000), true)

	a, found := s.keeper.GetAuction(s.ctx, auction.Id)
	s.Require().True(found)

	mInfo := s.keeper.CalculateBatchAllocation(s.ctx, a)

	// Checking
	s.Require().Equal(mInfo.MatchedLen, int64(2))
	s.Require().Equal(mInfo.MatchedPrice, parseDec("0.9"))
	matchingPrice := parseDec("0.9")
	matchedAmt1 := sdk.NewInt(500_000_000)
	matchedAmt2 := sdk.NewInt(500_000_000).ToDec().QuoTruncate(matchingPrice).TruncateInt()

	s.Require().Equal(mInfo.TotalMatchedAmount, sdk.NewInt(500_000_000).Add(matchedAmt2))
	s.Require().Equal(mInfo.AllocationMap[s.addr(1).String()], matchedAmt1)
	s.Require().Equal(mInfo.AllocationMap[s.addr(2).String()], matchedAmt2)
	s.Require().Equal(mInfo.AllocationMap[s.addr(3).String()], sdk.NewInt(0))
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(1).String()], sdk.NewInt(450_000_000))
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(2).String()], sdk.NewInt(500_000_000))
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(3).String()], sdk.NewInt(0))
	s.Require().Equal(mInfo.RefundMap[s.addr(1).String()], sdk.NewInt(50_000_000))
	s.Require().Equal(mInfo.RefundMap[s.addr(3).String()], sdk.NewInt(500_000_000))
	s.Require().True(mInfo.RefundMap[s.addr(2).String()].IsZero())

	// Distribute selling coin
	err := s.keeper.AllocateSellingCoin(s.ctx, auction, mInfo)
	s.Require().NoError(err)

	err = s.keeper.RefundRemainingSellingCoin(s.ctx, auction)
	s.Require().NoError(err)

	// The selling reserve account balance must be zero
	s.Require().True(s.getBalance(auction.GetSellingReserveAddress(), auction.SellingCoin.Denom).IsZero())

	// The auctioneer must have sellingCoin.Amount - TotalMatchedAmount
	s.Require().True(
		s.getBalance(s.addr(0), auction.GetSellingCoin().Denom).Amount.
			Equal(auction.SellingCoin.Amount.Sub(mInfo.TotalMatchedAmount)),
	)

	// The bidders must have the matched selling coin
	s.Require().Equal(s.getBalance(s.addr(1), auction.GetSellingCoin().Denom).Amount, matchedAmt1)
	s.Require().Equal(s.getBalance(s.addr(2), auction.GetSellingCoin().Denom).Amount, matchedAmt2)
	s.Require().True(s.getBalance(s.addr(3), auction.GetSellingCoin().Denom).IsZero())

	// Refund payingCoin
	err = s.keeper.RefundPayingCoin(s.ctx, auction, mInfo)
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) TestCalculateAllocation_Many_Limited() {
	auction := s.createBatchAuction(
		s.addr(0),
		parseDec("1"),
		parseDec("0.1"),
		parseCoin("1000_000_000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		1,
		sdk.MustNewDecFromStr("0.2"),
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	s.placeBidBatchMany(auction.Id, s.addr(1), parseDec("1"), parseCoin("400_000_000denom1"), sdk.NewInt(400_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(2), parseDec("0.9"), parseCoin("400_000_000denom1"), sdk.NewInt(400_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(3), parseDec("0.8"), parseCoin("400_000_000denom1"), sdk.NewInt(400_000_000), true)

	a, found := s.keeper.GetAuction(s.ctx, auction.Id)
	s.Require().True(found)

	mInfo := s.keeper.CalculateBatchAllocation(s.ctx, a)

	// Checking
	s.Require().Equal(mInfo.MatchedLen, int64(2))
	s.Require().Equal(mInfo.MatchedPrice, parseDec("0.9"))
	s.Require().Equal(mInfo.TotalMatchedAmount, sdk.NewInt(800_000_000))
	s.Require().Equal(mInfo.AllocationMap[s.addr(1).String()], sdk.NewInt(400_000_000))
	s.Require().Equal(mInfo.AllocationMap[s.addr(2).String()], sdk.NewInt(400_000_000))
	s.Require().Equal(mInfo.AllocationMap[s.addr(3).String()], sdk.NewInt(0))
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(1).String()], sdk.NewInt(360_000_000))
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(2).String()], sdk.NewInt(360_000_000))
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(3).String()], sdk.NewInt(0))
	s.Require().Equal(mInfo.RefundMap[s.addr(1).String()], sdk.NewInt(40_000_000))
	s.Require().True(mInfo.RefundMap[s.addr(2).String()].IsZero())
	s.Require().Equal(mInfo.RefundMap[s.addr(3).String()], sdk.NewInt(320_000_000))

	// Distribute selling coin
	err := s.keeper.AllocateSellingCoin(s.ctx, auction, mInfo)
	s.Require().NoError(err)

	err = s.keeper.RefundRemainingSellingCoin(s.ctx, auction)
	s.Require().NoError(err)

	// The selling reserve account balance must be zero
	s.Require().True(s.getBalance(auction.GetSellingReserveAddress(), auction.SellingCoin.Denom).IsZero())

	// The auctioneer must have sellingCoin.Amount - TotalMatchedAmount
	s.Require().Equal(s.getBalance(s.addr(0), auction.GetSellingCoin().Denom).Amount, sdk.NewInt(200_000_000))

	// The bidders must have the matched selling coin
	s.Require().Equal(s.getBalance(s.addr(1), auction.GetSellingCoin().Denom).Amount, sdk.NewInt(400_000_000))
	s.Require().Equal(s.getBalance(s.addr(2), auction.GetSellingCoin().Denom).Amount, sdk.NewInt(400_000_000))
	s.Require().True(s.getBalance(s.addr(3), auction.GetSellingCoin().Denom).IsZero())

	// Refund payingCoin
	err = s.keeper.RefundPayingCoin(s.ctx, auction, mInfo)
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) TestCalculateAllocation_Worth_Limited() {
	auction := s.createBatchAuction(
		s.addr(0),
		parseDec("1"),
		parseDec("0.1"),
		parseCoin("1500_000_000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		1,
		sdk.MustNewDecFromStr("0.2"),
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	s.placeBidBatchWorth(auction.Id, s.addr(1), parseDec("1"), parseCoin("400_000_000denom2"), sdk.NewInt(400_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(2), parseDec("0.9"), parseCoin("360_000_000denom2"), sdk.NewInt(400_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(3), parseDec("0.8"), parseCoin("320_000_000denom2"), sdk.NewInt(400_000_000), true)

	a, found := s.keeper.GetAuction(s.ctx, auction.Id)
	s.Require().True(found)

	mInfo := s.keeper.CalculateBatchAllocation(s.ctx, a)

	// Checking
	s.Require().Equal(mInfo.MatchedLen, int64(3))
	s.Require().Equal(mInfo.MatchedPrice, parseDec("0.8"))
	s.Require().Equal(mInfo.TotalMatchedAmount, sdk.NewInt(1200_000_000))
	s.Require().Equal(mInfo.AllocationMap[s.addr(1).String()], sdk.NewInt(400_000_000))
	s.Require().Equal(mInfo.AllocationMap[s.addr(2).String()], sdk.NewInt(400_000_000))
	s.Require().Equal(mInfo.AllocationMap[s.addr(3).String()], sdk.NewInt(400_000_000))
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(1).String()], sdk.NewInt(320_000_000))
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(2).String()], sdk.NewInt(320_000_000))
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(3).String()], sdk.NewInt(320_000_000))
	s.Require().Equal(mInfo.RefundMap[s.addr(1).String()], sdk.NewInt(80_000_000))
	s.Require().Equal(mInfo.RefundMap[s.addr(2).String()], sdk.NewInt(40_000_000))
	s.Require().True(mInfo.RefundMap[s.addr(3).String()].IsZero())

	// Distribute selling coin
	err := s.keeper.AllocateSellingCoin(s.ctx, auction, mInfo)
	s.Require().NoError(err)

	err = s.keeper.RefundRemainingSellingCoin(s.ctx, auction)
	s.Require().NoError(err)

	// The selling reserve account balance must be zero
	s.Require().True(s.getBalance(auction.GetSellingReserveAddress(), auction.SellingCoin.Denom).IsZero())

	// The auctioneer must have sellingCoin.Amount - TotalMatchedAmount
	s.Require().Equal(s.getBalance(s.addr(0), auction.GetSellingCoin().Denom).Amount, sdk.NewInt(300_000_000))

	// The bidders must have the matched selling coin
	s.Require().Equal(s.getBalance(s.addr(1), auction.GetSellingCoin().Denom).Amount, sdk.NewInt(400_000_000))
	s.Require().Equal(s.getBalance(s.addr(2), auction.GetSellingCoin().Denom).Amount, sdk.NewInt(400_000_000))
	s.Require().Equal(s.getBalance(s.addr(3), auction.GetSellingCoin().Denom).Amount, sdk.NewInt(400_000_000))

	// Refund payingCoin
	err = s.keeper.RefundPayingCoin(s.ctx, auction, mInfo)
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) TestCalculateAllocation_Mixed_Limited() {
	auction := s.createBatchAuction(
		s.addr(0),
		parseDec("1"),
		parseDec("0.1"),
		parseCoin("1700_000_000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		1,
		sdk.MustNewDecFromStr("0.2"),
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	s.placeBidBatchMany(auction.Id, s.addr(1), parseDec("1"), parseCoin("500_000_000denom1"), sdk.NewInt(600_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(2), parseDec("0.9"), parseCoin("500_000_000denom2"), sdk.NewInt(600_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(3), parseDec("0.8"), parseCoin("450_000_000denom2"), sdk.NewInt(600_000_000), true)

	a, found := s.keeper.GetAuction(s.ctx, auction.Id)
	s.Require().True(found)

	mInfo := s.keeper.CalculateBatchAllocation(s.ctx, a)

	// Checking
	s.Require().Equal(mInfo.MatchedLen, int64(3))
	s.Require().Equal(mInfo.MatchedPrice, parseDec("0.8"))
	s.Require().Equal(mInfo.TotalMatchedAmount, sdk.NewInt(1662_500_000))
	s.Require().Equal(mInfo.AllocationMap[s.addr(1).String()], sdk.NewInt(500_000_000))
	s.Require().Equal(mInfo.AllocationMap[s.addr(2).String()], sdk.NewInt(600_000_000))
	s.Require().Equal(mInfo.AllocationMap[s.addr(3).String()], sdk.NewInt(562_500_000))
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(1).String()], sdk.NewInt(400_000_000))
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(2).String()], sdk.NewInt(480_000_000))
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(3).String()], sdk.NewInt(450_000_000))
	s.Require().Equal(mInfo.RefundMap[s.addr(1).String()], sdk.NewInt(100_000_000))
	s.Require().Equal(mInfo.RefundMap[s.addr(2).String()], sdk.NewInt(20_000_000))
	s.Require().True(mInfo.RefundMap[s.addr(3).String()].IsZero())

	// Distribute selling coin
	err := s.keeper.AllocateSellingCoin(s.ctx, auction, mInfo)
	s.Require().NoError(err)

	err = s.keeper.RefundRemainingSellingCoin(s.ctx, auction)
	s.Require().NoError(err)

	// The selling reserve account balance must be zero
	s.Require().True(s.getBalance(auction.GetSellingReserveAddress(), auction.SellingCoin.Denom).IsZero())

	// The auctioneer must have sellingCoin.Amount - TotalMatchedAmount
	s.Require().True(
		s.getBalance(s.addr(0), auction.GetSellingCoin().Denom).Amount.
			Equal(sdk.NewInt(37_500_000)),
	)

	// The bidders must have the matched selling coin
	s.Require().Equal(s.getBalance(s.addr(1), auction.GetSellingCoin().Denom).Amount, sdk.NewInt(500_000_000))
	s.Require().Equal(s.getBalance(s.addr(2), auction.GetSellingCoin().Denom).Amount, sdk.NewInt(600_000_000))
	s.Require().Equal(s.getBalance(s.addr(3), auction.GetSellingCoin().Denom).Amount, sdk.NewInt(562_500_000))

	// Refund payingCoin
	err = s.keeper.RefundPayingCoin(s.ctx, auction, mInfo)
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) TestCalculateAllocation_Mixed2() {
	auction := s.createBatchAuction(
		s.addr(0),
		parseDec("1"),
		parseDec("0.1"),
		parseCoin("5000_000_000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		1,
		sdk.MustNewDecFromStr("0.2"),
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	s.placeBidBatchMany(auction.Id, s.addr(1), parseDec("1"), parseCoin("200_000_000denom1"), sdk.NewInt(5000_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(2), parseDec("0.8"), parseCoin("500_000_000denom2"), sdk.NewInt(5000_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(3), parseDec("0.9"), parseCoin("500_000_000denom1"), sdk.NewInt(5000_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(1), parseDec("1.1"), parseCoin("300_000_000denom2"), sdk.NewInt(0), true)
	s.placeBidBatchMany(auction.Id, s.addr(5), parseDec("1.2"), parseCoin("300_000_000denom1"), sdk.NewInt(5000_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(4), parseDec("0.8"), parseCoin("100_000_000denom1"), sdk.NewInt(5000_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(2), parseDec("0.7"), parseCoin("100_000_000denom1"), sdk.NewInt(0), true)
	s.placeBidBatchMany(auction.Id, s.addr(6), parseDec("0.5"), parseCoin("100_000_000denom1"), sdk.NewInt(5000_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(3), parseDec("0.8"), parseCoin("300_000_000denom2"), sdk.NewInt(0), true)
	s.placeBidBatchWorth(auction.Id, s.addr(7), parseDec("0.6"), parseCoin("500_000_000denom2"), sdk.NewInt(5000_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(8), parseDec("0.8"), parseCoin("500_000_000denom1"), sdk.NewInt(5000_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(9), parseDec("0.6"), parseCoin("600_000_000denom1"), sdk.NewInt(5000_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(6), parseDec("0.5"), parseCoin("500_000_000denom2"), sdk.NewInt(0), true)
	s.placeBidBatchMany(auction.Id, s.addr(10), parseDec("0.6"), parseCoin("100_000_000denom1"), sdk.NewInt(5000_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(3), parseDec("0.7"), parseCoin("800_000_000denom2"), sdk.NewInt(0), true)

	a, found := s.keeper.GetAuction(s.ctx, auction.Id)
	s.Require().True(found)

	mInfo := s.keeper.CalculateBatchAllocation(s.ctx, a)

	// Checking
	s.Require().Equal(mInfo.MatchedLen, int64(10))
	matchingPrice := parseDec("0.7")
	s.Require().Equal(mInfo.MatchedPrice, matchingPrice)

	matchedAmt1 := sdk.NewInt(300_000_000).ToDec().QuoTruncate(matchingPrice).TruncateInt().Add(sdk.NewInt(200_000_000))
	matchedAmt2 := sdk.NewInt(500_000_000).ToDec().QuoTruncate(matchingPrice).TruncateInt().Add(sdk.NewInt(100_000_000))
	tMatchedAmt3 := sdk.NewInt(300_000_000).ToDec().QuoTruncate(matchingPrice).TruncateInt().Add(sdk.NewInt(500_000_000))
	matchedAmt3 := tMatchedAmt3.Add(sdk.NewInt(800_000_000).ToDec().QuoTruncate(matchingPrice).TruncateInt())
	matchedAmt4 := sdk.NewInt(100_000_000)
	matchedAmt5 := sdk.NewInt(300_000_000)
	matchedAmt8 := sdk.NewInt(500_000_000)
	matchedAmt_Zero := sdk.NewInt(0)
	totalMatchedAmt := matchedAmt1.Add(matchedAmt2).Add(matchedAmt3).Add(matchedAmt4).Add(matchedAmt5).Add(matchedAmt8)

	s.Require().Equal(mInfo.TotalMatchedAmount, totalMatchedAmt)
	s.Require().Equal(mInfo.AllocationMap[s.addr(1).String()], matchedAmt1)
	s.Require().Equal(mInfo.AllocationMap[s.addr(2).String()], matchedAmt2)
	s.Require().Equal(mInfo.AllocationMap[s.addr(3).String()], matchedAmt3)
	s.Require().Equal(mInfo.AllocationMap[s.addr(4).String()], matchedAmt4)
	s.Require().Equal(mInfo.AllocationMap[s.addr(5).String()], matchedAmt5)
	s.Require().Equal(mInfo.AllocationMap[s.addr(6).String()], matchedAmt_Zero)
	s.Require().Equal(mInfo.AllocationMap[s.addr(7).String()], matchedAmt_Zero)
	s.Require().Equal(mInfo.AllocationMap[s.addr(8).String()], matchedAmt8)
	s.Require().Equal(mInfo.AllocationMap[s.addr(9).String()], matchedAmt_Zero)
	s.Require().Equal(mInfo.AllocationMap[s.addr(10).String()], matchedAmt_Zero)

	reservedmatchedAmt1 := sdk.NewInt(200_000_000).ToDec().Mul(matchingPrice).Ceil().TruncateInt().Add(sdk.NewInt(300_000_000))
	reservedMatchedAmt2 := sdk.NewInt(100_000_000).ToDec().Mul(matchingPrice).Ceil().TruncateInt().Add(sdk.NewInt(500_000_000))
	reservedMatchedAmt3 := sdk.NewInt(500_000_000).ToDec().Mul(matchingPrice).Ceil().TruncateInt().Add(sdk.NewInt(1100_000_000))
	reservedMatchedAmt4 := sdk.NewInt(100_000_000).ToDec().Mul(matchingPrice).Ceil().TruncateInt()
	reservedMatchedAmt5 := sdk.NewInt(300_000_000).ToDec().Mul(matchingPrice).Ceil().TruncateInt()
	reservedMatchedAmt8 := sdk.NewInt(500_000_000).ToDec().Mul(matchingPrice).Ceil().TruncateInt()
	reservedMatchedAmt_Zero := sdk.NewInt(0)

	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(1).String()], reservedmatchedAmt1)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(2).String()], reservedMatchedAmt2)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(3).String()], reservedMatchedAmt3)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(4).String()], reservedMatchedAmt4)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(5).String()], reservedMatchedAmt5)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(6).String()], reservedMatchedAmt_Zero)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(7).String()], reservedMatchedAmt_Zero)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(8).String()], reservedMatchedAmt8)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(9).String()], reservedMatchedAmt_Zero)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(10).String()], reservedMatchedAmt_Zero)

	refundAmt1 := sdk.NewInt(60_000_000)
	refundAmt2 := sdk.NewInt(0)
	refundAmt3 := sdk.NewInt(100_000_000)
	refundAmt4 := sdk.NewInt(10_000_000)
	refundAmt5 := sdk.NewInt(150_000_000)
	refundAmt6 := sdk.NewInt(550_000_000)
	refundAmt7 := sdk.NewInt(500_000_000)
	refundAmt8 := sdk.NewInt(50_000_000)
	refundAmt9 := sdk.NewInt(360_000_000)
	refundAmt10 := sdk.NewInt(60_000_000)

	s.Require().True(mInfo.RefundMap[s.addr(1).String()].Equal(refundAmt1))
	s.Require().True(mInfo.RefundMap[s.addr(2).String()].Equal(refundAmt2))
	s.Require().True(mInfo.RefundMap[s.addr(3).String()].Equal(refundAmt3))
	s.Require().True(mInfo.RefundMap[s.addr(4).String()].Equal(refundAmt4))
	s.Require().True(mInfo.RefundMap[s.addr(5).String()].Equal(refundAmt5))
	s.Require().True(mInfo.RefundMap[s.addr(6).String()].Equal(refundAmt6))
	s.Require().True(mInfo.RefundMap[s.addr(7).String()].Equal(refundAmt7))
	s.Require().True(mInfo.RefundMap[s.addr(8).String()].Equal(refundAmt8))
	s.Require().True(mInfo.RefundMap[s.addr(9).String()].Equal(refundAmt9))
	s.Require().True(mInfo.RefundMap[s.addr(10).String()].Equal(refundAmt10))

	// Distribute selling coin
	err := s.keeper.AllocateSellingCoin(s.ctx, auction, mInfo)
	s.Require().NoError(err)

	err = s.keeper.RefundRemainingSellingCoin(s.ctx, auction)
	s.Require().NoError(err)

	// The selling reserve account balance must be zero
	s.Require().True(s.getBalance(auction.GetSellingReserveAddress(), auction.SellingCoin.Denom).IsZero())

	// The auctioneer must have sellingCoin.Amount - TotalMatchedAmount
	s.Require().Equal(s.getBalance(s.addr(0), auction.GetSellingCoin().Denom).Amount, auction.SellingCoin.Amount.Sub(mInfo.TotalMatchedAmount))

	// The bidders must have the matched selling coin
	s.Require().Equal(s.getBalance(s.addr(1), auction.GetSellingCoin().Denom).Amount, matchedAmt1)
	s.Require().Equal(s.getBalance(s.addr(2), auction.GetSellingCoin().Denom).Amount, matchedAmt2)
	s.Require().Equal(s.getBalance(s.addr(3), auction.GetSellingCoin().Denom).Amount, matchedAmt3)
	s.Require().Equal(s.getBalance(s.addr(4), auction.GetSellingCoin().Denom).Amount, matchedAmt4)
	s.Require().Equal(s.getBalance(s.addr(5), auction.GetSellingCoin().Denom).Amount, matchedAmt5)
	s.Require().True(s.getBalance(s.addr(6), auction.GetSellingCoin().Denom).Amount.IsZero())
	s.Require().True(s.getBalance(s.addr(7), auction.GetSellingCoin().Denom).Amount.IsZero())
	s.Require().Equal(s.getBalance(s.addr(8), auction.GetSellingCoin().Denom).Amount, matchedAmt8)
	s.Require().True(s.getBalance(s.addr(9), auction.GetSellingCoin().Denom).Amount.IsZero())
	s.Require().True(s.getBalance(s.addr(10), auction.GetSellingCoin().Denom).Amount.IsZero())

	// Refund payingCoin
	err = s.keeper.RefundPayingCoin(s.ctx, auction, mInfo)
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) TestCalculateAllocation_Mixed2_LimitedSame() {
	auction := s.createBatchAuction(
		s.addr(0),
		parseDec("1"),
		parseDec("0.1"),
		parseCoin("5000_000_000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		1,
		sdk.MustNewDecFromStr("0.2"),
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	s.placeBidBatchMany(auction.Id, s.addr(1), parseDec("1"), parseCoin("200_000_000denom1"), sdk.NewInt(700_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(2), parseDec("0.8"), parseCoin("500_000_000denom2"), sdk.NewInt(700_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(3), parseDec("0.9"), parseCoin("500_000_000denom1"), sdk.NewInt(700_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(1), parseDec("1.1"), parseCoin("300_000_000denom2"), sdk.NewInt(0), true)
	s.placeBidBatchMany(auction.Id, s.addr(5), parseDec("1.2"), parseCoin("300_000_000denom1"), sdk.NewInt(700_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(4), parseDec("0.8"), parseCoin("100_000_000denom1"), sdk.NewInt(700_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(2), parseDec("0.7"), parseCoin("100_000_000denom1"), sdk.NewInt(0), true)
	s.placeBidBatchMany(auction.Id, s.addr(6), parseDec("0.5"), parseCoin("100_000_000denom1"), sdk.NewInt(700_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(3), parseDec("0.8"), parseCoin("300_000_000denom2"), sdk.NewInt(0), true)
	s.placeBidBatchWorth(auction.Id, s.addr(7), parseDec("0.6"), parseCoin("400_000_000denom2"), sdk.NewInt(700_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(8), parseDec("0.8"), parseCoin("500_000_000denom1"), sdk.NewInt(700_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(9), parseDec("0.6"), parseCoin("600_000_000denom1"), sdk.NewInt(700_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(6), parseDec("0.5"), parseCoin("350_000_000denom2"), sdk.NewInt(0), true)
	s.placeBidBatchMany(auction.Id, s.addr(10), parseDec("0.6"), parseCoin("100_000_000denom1"), sdk.NewInt(700_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(3), parseDec("0.7"), parseCoin("490_000_000denom2"), sdk.NewInt(0), true)

	a, found := s.keeper.GetAuction(s.ctx, auction.Id)
	s.Require().True(found)

	mInfo := s.keeper.CalculateBatchAllocation(s.ctx, a)

	// Checking
	s.Require().Equal(mInfo.MatchedLen, int64(13))
	matchingPrice := parseDec("0.6")
	s.Require().Equal(mInfo.MatchedPrice, matchingPrice)

	matchedAmt1 := sdk.NewInt(700_000_000)
	matchedAmt2 := sdk.NewInt(700_000_000)
	matchedAmt3 := sdk.NewInt(700_000_000)
	matchedAmt4 := sdk.NewInt(100_000_000)
	matchedAmt5 := sdk.NewInt(300_000_000)
	matchedAmt6 := sdk.NewInt(0)
	matchedAmt7 := sdk.NewInt(400_000_000).ToDec().QuoTruncate(matchingPrice).TruncateInt()
	matchedAmt8 := sdk.NewInt(500_000_000)
	matchedAmt9 := sdk.NewInt(600_000_000)
	matchedAmt10 := sdk.NewInt(100_000_000)

	totalMatchedAmt := sdk.NewInt(3700_000_000).Add(matchedAmt7)

	s.Require().Equal(mInfo.TotalMatchedAmount, totalMatchedAmt)
	s.Require().Equal(mInfo.AllocationMap[s.addr(1).String()], matchedAmt1)
	s.Require().Equal(mInfo.AllocationMap[s.addr(2).String()], matchedAmt2)
	s.Require().Equal(mInfo.AllocationMap[s.addr(3).String()], matchedAmt3)
	s.Require().Equal(mInfo.AllocationMap[s.addr(4).String()], matchedAmt4)
	s.Require().Equal(mInfo.AllocationMap[s.addr(5).String()], matchedAmt5)
	s.Require().Equal(mInfo.AllocationMap[s.addr(6).String()], matchedAmt6)
	s.Require().Equal(mInfo.AllocationMap[s.addr(7).String()], matchedAmt7)
	s.Require().Equal(mInfo.AllocationMap[s.addr(8).String()], matchedAmt8)
	s.Require().Equal(mInfo.AllocationMap[s.addr(9).String()], matchedAmt9)
	s.Require().Equal(mInfo.AllocationMap[s.addr(10).String()], matchedAmt10)

	reservedMatchedAmt1 := sdk.NewInt(420_000_000)
	reservedMatchedAmt2 := sdk.NewInt(420_000_000)
	reservedMatchedAmt3 := sdk.NewInt(420_000_000)
	reservedMatchedAmt4 := sdk.NewInt(60_000_000)
	reservedMatchedAmt5 := sdk.NewInt(180_000_000)
	reservedMatchedAmt6 := sdk.NewInt(0)
	reservedMatchedAmt7 := sdk.NewInt(400_000_000)
	reservedMatchedAmt8 := sdk.NewInt(300_000_000)
	reservedMatchedAmt9 := sdk.NewInt(360_000_000)
	reservedMatchedAmt10 := sdk.NewInt(60_000_000)

	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(1).String()], reservedMatchedAmt1)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(2).String()], reservedMatchedAmt2)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(3).String()], reservedMatchedAmt3)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(4).String()], reservedMatchedAmt4)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(5).String()], reservedMatchedAmt5)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(6).String()], reservedMatchedAmt6)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(7).String()], reservedMatchedAmt7)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(8).String()], reservedMatchedAmt8)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(9).String()], reservedMatchedAmt9)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(10).String()], reservedMatchedAmt10)

	refundAmt1 := sdk.NewInt(80_000_000)
	refundAmt2 := sdk.NewInt(150_000_000)
	refundAmt3 := sdk.NewInt(820_000_000)
	refundAmt4 := sdk.NewInt(20_000_000)
	refundAmt5 := sdk.NewInt(180_000_000)
	refundAmt6 := sdk.NewInt(400_000_000)
	refundAmt7 := sdk.NewInt(0)
	refundAmt8 := sdk.NewInt(100_000_000)
	refundAmt9 := sdk.NewInt(0)
	refundAmt10 := sdk.NewInt(0)

	s.Require().True(mInfo.RefundMap[s.addr(1).String()].Equal(refundAmt1))
	s.Require().True(mInfo.RefundMap[s.addr(2).String()].Equal(refundAmt2))
	s.Require().True(mInfo.RefundMap[s.addr(3).String()].Equal(refundAmt3))
	s.Require().True(mInfo.RefundMap[s.addr(4).String()].Equal(refundAmt4))
	s.Require().True(mInfo.RefundMap[s.addr(5).String()].Equal(refundAmt5))
	s.Require().True(mInfo.RefundMap[s.addr(6).String()].Equal(refundAmt6))
	s.Require().True(mInfo.RefundMap[s.addr(7).String()].Equal(refundAmt7))
	s.Require().True(mInfo.RefundMap[s.addr(8).String()].Equal(refundAmt8))
	s.Require().True(mInfo.RefundMap[s.addr(9).String()].Equal(refundAmt9))
	s.Require().True(mInfo.RefundMap[s.addr(10).String()].Equal(refundAmt10))

	// Distribute selling coin
	err := s.keeper.AllocateSellingCoin(s.ctx, auction, mInfo)
	s.Require().NoError(err)

	err = s.keeper.RefundRemainingSellingCoin(s.ctx, auction)
	s.Require().NoError(err)

	// The selling reserve account balance must be zero
	s.Require().True(s.getBalance(auction.GetSellingReserveAddress(), auction.SellingCoin.Denom).IsZero())

	// The auctioneer must have sellingCoin.Amount - TotalMatchedAmount
	s.Require().Equal(s.getBalance(s.addr(0), auction.GetSellingCoin().Denom).Amount, auction.SellingCoin.Amount.Sub(mInfo.TotalMatchedAmount))

	// The bidders must have the matched selling coin
	s.Require().Equal(s.getBalance(s.addr(1), auction.GetSellingCoin().Denom).Amount, matchedAmt1)
	s.Require().Equal(s.getBalance(s.addr(2), auction.GetSellingCoin().Denom).Amount, matchedAmt2)
	s.Require().Equal(s.getBalance(s.addr(3), auction.GetSellingCoin().Denom).Amount, matchedAmt3)
	s.Require().Equal(s.getBalance(s.addr(4), auction.GetSellingCoin().Denom).Amount, matchedAmt4)
	s.Require().Equal(s.getBalance(s.addr(5), auction.GetSellingCoin().Denom).Amount, matchedAmt5)
	s.Require().True(s.getBalance(s.addr(6), auction.GetSellingCoin().Denom).Amount.Equal(matchedAmt6))
	s.Require().Equal(s.getBalance(s.addr(7), auction.GetSellingCoin().Denom).Amount, matchedAmt7)
	s.Require().Equal(s.getBalance(s.addr(8), auction.GetSellingCoin().Denom).Amount, matchedAmt8)
	s.Require().Equal(s.getBalance(s.addr(9), auction.GetSellingCoin().Denom).Amount, matchedAmt9)
	s.Require().Equal(s.getBalance(s.addr(10), auction.GetSellingCoin().Denom).Amount, matchedAmt10)

	// Refund payingCoin
	err = s.keeper.RefundPayingCoin(s.ctx, auction, mInfo)
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) TestCalculateAllocation_Mixed2_LimitedDifferent() {
	auction := s.createBatchAuction(
		s.addr(0),
		parseDec("1"),
		parseDec("0.1"),
		parseCoin("5000_000_000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		1,
		sdk.MustNewDecFromStr("0.2"),
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	s.placeBidBatchMany(auction.Id, s.addr(1), parseDec("1"), parseCoin("200_000_000denom1"), sdk.NewInt(1000_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(2), parseDec("0.8"), parseCoin("500_000_000denom2"), sdk.NewInt(1000_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(3), parseDec("0.9"), parseCoin("500_000_000denom1"), sdk.NewInt(800_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(1), parseDec("1.1"), parseCoin("300_000_000denom2"), sdk.NewInt(0), true)
	s.placeBidBatchMany(auction.Id, s.addr(5), parseDec("1.2"), parseCoin("300_000_000denom1"), sdk.NewInt(600_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(4), parseDec("0.8"), parseCoin("100_000_000denom1"), sdk.NewInt(800_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(2), parseDec("0.7"), parseCoin("100_000_000denom1"), sdk.NewInt(0), true)
	s.placeBidBatchMany(auction.Id, s.addr(6), parseDec("0.5"), parseCoin("100_000_000denom1"), sdk.NewInt(600_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(3), parseDec("0.8"), parseCoin("300_000_000denom2"), sdk.NewInt(0), true)
	s.placeBidBatchWorth(auction.Id, s.addr(7), parseDec("0.6"), parseCoin("200_000_000denom2"), sdk.NewInt(400_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(8), parseDec("0.8"), parseCoin("400_000_000denom1"), sdk.NewInt(400_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(9), parseDec("0.6"), parseCoin("200_000_000denom1"), sdk.NewInt(200_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(6), parseDec("0.5"), parseCoin("300_000_000denom2"), sdk.NewInt(0), true)
	s.placeBidBatchMany(auction.Id, s.addr(10), parseDec("0.6"), parseCoin("100_000_000denom1"), sdk.NewInt(200_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(3), parseDec("0.7"), parseCoin("560_000_000denom2"), sdk.NewInt(0), true)

	a, found := s.keeper.GetAuction(s.ctx, auction.Id)
	s.Require().True(found)

	mInfo := s.keeper.CalculateBatchAllocation(s.ctx, a)

	// Checking
	s.Require().Equal(mInfo.MatchedLen, int64(15))
	matchingPrice := parseDec("0.5")
	s.Require().Equal(mInfo.MatchedPrice, matchingPrice)

	matchedAmt1 := sdk.NewInt(800_000_000)
	matchedAmt2 := sdk.NewInt(1000_000_000)
	matchedAmt3 := sdk.NewInt(800_000_000)
	matchedAmt4 := sdk.NewInt(100_000_000)
	matchedAmt5 := sdk.NewInt(300_000_000)
	matchedAmt6 := sdk.NewInt(600_000_000)
	matchedAmt7 := sdk.NewInt(400_000_000)
	matchedAmt8 := sdk.NewInt(400_000_000)
	matchedAmt9 := sdk.NewInt(200_000_000)
	matchedAmt10 := sdk.NewInt(100_000_000)

	totalMatchedAmt := sdk.NewInt(4700_000_000)

	s.Require().Equal(mInfo.TotalMatchedAmount, totalMatchedAmt)
	s.Require().Equal(mInfo.AllocationMap[s.addr(1).String()], matchedAmt1)
	s.Require().Equal(mInfo.AllocationMap[s.addr(2).String()], matchedAmt2)
	s.Require().Equal(mInfo.AllocationMap[s.addr(3).String()], matchedAmt3)
	s.Require().Equal(mInfo.AllocationMap[s.addr(4).String()], matchedAmt4)
	s.Require().Equal(mInfo.AllocationMap[s.addr(5).String()], matchedAmt5)
	s.Require().Equal(mInfo.AllocationMap[s.addr(6).String()], matchedAmt6)
	s.Require().Equal(mInfo.AllocationMap[s.addr(7).String()], matchedAmt7)
	s.Require().Equal(mInfo.AllocationMap[s.addr(8).String()], matchedAmt8)
	s.Require().Equal(mInfo.AllocationMap[s.addr(9).String()], matchedAmt9)
	s.Require().Equal(mInfo.AllocationMap[s.addr(10).String()], matchedAmt10)

	reservedMatchedAmt1 := sdk.NewInt(400_000_000)
	reservedMatchedAmt2 := sdk.NewInt(500_000_000)
	reservedMatchedAmt3 := sdk.NewInt(400_000_000)
	reservedMatchedAmt4 := sdk.NewInt(50_000_000)
	reservedMatchedAmt5 := sdk.NewInt(150_000_000)
	reservedMatchedAmt6 := sdk.NewInt(300_000_000)
	reservedMatchedAmt7 := sdk.NewInt(200_000_000)
	reservedMatchedAmt8 := sdk.NewInt(200_000_000)
	reservedMatchedAmt9 := sdk.NewInt(100_000_000)
	reservedMatchedAmt10 := sdk.NewInt(50_000_000)

	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(1).String()], reservedMatchedAmt1)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(2).String()], reservedMatchedAmt2)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(3).String()], reservedMatchedAmt3)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(4).String()], reservedMatchedAmt4)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(5).String()], reservedMatchedAmt5)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(6).String()], reservedMatchedAmt6)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(7).String()], reservedMatchedAmt7)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(8).String()], reservedMatchedAmt8)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(9).String()], reservedMatchedAmt9)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(10).String()], reservedMatchedAmt10)

	refundAmt1 := sdk.NewInt(100_000_000)
	refundAmt2 := sdk.NewInt(70_000_000)
	refundAmt3 := sdk.NewInt(910_000_000)
	refundAmt4 := sdk.NewInt(30_000_000)
	refundAmt5 := sdk.NewInt(210_000_000)
	refundAmt6 := sdk.NewInt(50_000_000)
	refundAmt7 := sdk.NewInt(0)
	refundAmt8 := sdk.NewInt(120_000_000)
	refundAmt9 := sdk.NewInt(20_000_000)
	refundAmt10 := sdk.NewInt(10_000_000)

	s.Require().True(mInfo.RefundMap[s.addr(1).String()].Equal(refundAmt1))
	s.Require().True(mInfo.RefundMap[s.addr(2).String()].Equal(refundAmt2))
	s.Require().True(mInfo.RefundMap[s.addr(3).String()].Equal(refundAmt3))
	s.Require().True(mInfo.RefundMap[s.addr(4).String()].Equal(refundAmt4))
	s.Require().True(mInfo.RefundMap[s.addr(5).String()].Equal(refundAmt5))
	s.Require().True(mInfo.RefundMap[s.addr(6).String()].Equal(refundAmt6))
	s.Require().True(mInfo.RefundMap[s.addr(7).String()].Equal(refundAmt7))
	s.Require().True(mInfo.RefundMap[s.addr(8).String()].Equal(refundAmt8))
	s.Require().True(mInfo.RefundMap[s.addr(9).String()].Equal(refundAmt9))
	s.Require().True(mInfo.RefundMap[s.addr(10).String()].Equal(refundAmt10))

	// Distribute selling coin
	err := s.keeper.AllocateSellingCoin(s.ctx, auction, mInfo)
	s.Require().NoError(err)

	err = s.keeper.RefundRemainingSellingCoin(s.ctx, auction)
	s.Require().NoError(err)

	// The selling reserve account balance must be zero
	s.Require().True(s.getBalance(auction.GetSellingReserveAddress(), auction.SellingCoin.Denom).IsZero())

	// The auctioneer must have sellingCoin.Amount - TotalMatchedAmount
	s.Require().Equal(s.getBalance(s.addr(0), auction.GetSellingCoin().Denom).Amount, sdk.NewInt(300_000_000))

	// The bidders must have the matched selling coin
	s.Require().Equal(s.getBalance(s.addr(1), auction.GetSellingCoin().Denom).Amount, matchedAmt1)
	s.Require().Equal(s.getBalance(s.addr(2), auction.GetSellingCoin().Denom).Amount, matchedAmt2)
	s.Require().Equal(s.getBalance(s.addr(3), auction.GetSellingCoin().Denom).Amount, matchedAmt3)
	s.Require().Equal(s.getBalance(s.addr(4), auction.GetSellingCoin().Denom).Amount, matchedAmt4)
	s.Require().Equal(s.getBalance(s.addr(5), auction.GetSellingCoin().Denom).Amount, matchedAmt5)
	s.Require().Equal(s.getBalance(s.addr(6), auction.GetSellingCoin().Denom).Amount, matchedAmt6)
	s.Require().Equal(s.getBalance(s.addr(7), auction.GetSellingCoin().Denom).Amount, matchedAmt7)
	s.Require().Equal(s.getBalance(s.addr(8), auction.GetSellingCoin().Denom).Amount, matchedAmt8)
	s.Require().Equal(s.getBalance(s.addr(9), auction.GetSellingCoin().Denom).Amount, matchedAmt9)
	s.Require().Equal(s.getBalance(s.addr(10), auction.GetSellingCoin().Denom).Amount, matchedAmt10)

	// Refund payingCoin
	err = s.keeper.RefundPayingCoin(s.ctx, auction, mInfo)
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) TestCalculateAllocation_Mixed3() {
	auction := s.createBatchAuction(
		s.addr(0),
		parseDec("10"),
		parseDec("0.1"),
		parseCoin("2500_000_000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		1,
		sdk.MustNewDecFromStr("0.2"),
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	s.placeBidBatchMany(auction.Id, s.addr(1), parseDec("10"), parseCoin("200_000_000denom1"), sdk.NewInt(2500_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(2), parseDec("11"), parseCoin("2000_000_000denom2"), sdk.NewInt(2500_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(3), parseDec("10.5"), parseCoin("500_000_000denom1"), sdk.NewInt(2500_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(4), parseDec("10.2"), parseCoin("1500_000_000denom2"), sdk.NewInt(2500_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(5), parseDec("10.8"), parseCoin("300_000_000denom1"), sdk.NewInt(2500_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(6), parseDec("11.4"), parseCoin("2500_000_000denom2"), sdk.NewInt(2500_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(7), parseDec("11.3"), parseCoin("100_000_000denom1"), sdk.NewInt(2500_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(8), parseDec("9.9"), parseCoin("2500_000_000denom2"), sdk.NewInt(2500_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(9), parseDec("10.1"), parseCoin("300_000_000denom1"), sdk.NewInt(2500_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(10), parseDec("10.45"), parseCoin("2000_000_000denom2"), sdk.NewInt(2500_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(11), parseDec("10.75"), parseCoin("150_000_000denom1"), sdk.NewInt(2500_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(12), parseDec("10.99"), parseCoin("1500_000_000denom2"), sdk.NewInt(2500_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(13), parseDec("10.2"), parseCoin("200_000_000denom1"), sdk.NewInt(2500_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(14), parseDec("9.87"), parseCoin("2000_000_000denom2"), sdk.NewInt(2500_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(15), parseDec("10.25"), parseCoin("200_000_000denom1"), sdk.NewInt(2500_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(16), parseDec("10.48"), parseCoin("2500_000_000denom2"), sdk.NewInt(2500_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(17), parseDec("10.52"), parseCoin("180_000_000denom1"), sdk.NewInt(2500_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(4), parseDec("10.8"), parseCoin("220_000_000denom1"), sdk.NewInt(0), true)
	s.placeBidBatchWorth(auction.Id, s.addr(5), parseDec("10.5"), parseCoin("1500_000_000denom2"), sdk.NewInt(0), true)
	s.placeBidBatchMany(auction.Id, s.addr(6), parseDec("9.7"), parseCoin("250_000_000denom1"), sdk.NewInt(0), true)

	a, found := s.keeper.GetAuction(s.ctx, auction.Id)
	s.Require().True(found)

	mInfo := s.keeper.CalculateBatchAllocation(s.ctx, a)

	// Checking
	s.Require().Equal(mInfo.MatchedLen, int64(11))
	matchingPrice := parseDec("10.48")
	s.Require().Equal(mInfo.MatchedPrice, matchingPrice)

	matchedAmt1 := sdk.NewInt(0)
	matchedAmt2 := sdk.NewInt(2000_000_000).ToDec().QuoTruncate(matchingPrice).TruncateInt()
	matchedAmt3 := sdk.NewInt(500_000_000)
	matchedAmt4 := sdk.NewInt(220_000_000)
	matchedAmt5 := sdk.NewInt(1500_000_000).ToDec().QuoTruncate(matchingPrice).TruncateInt().Add(sdk.NewInt(300_000_000))
	matchedAmt6 := sdk.NewInt(2500_000_000).ToDec().QuoTruncate(matchingPrice).TruncateInt()
	matchedAmt7 := sdk.NewInt(100_000_000)
	matchedAmt8 := sdk.NewInt(0)
	matchedAmt9 := sdk.NewInt(0)
	matchedAmt10 := sdk.NewInt(0)
	matchedAmt11 := sdk.NewInt(150_000_000)
	matchedAmt12 := sdk.NewInt(1500_000_000).ToDec().QuoTruncate(matchingPrice).TruncateInt()
	matchedAmt13 := sdk.NewInt(0)
	matchedAmt14 := sdk.NewInt(0)
	matchedAmt15 := sdk.NewInt(0)
	matchedAmt16 := sdk.NewInt(2500_000_000).ToDec().QuoTruncate(matchingPrice).TruncateInt()
	matchedAmt17 := sdk.NewInt(180_000_000)

	totalMatchedAmt := matchedAmt2.Add(matchedAmt3).
		Add(matchedAmt4).
		Add(matchedAmt5).
		Add(matchedAmt6).
		Add(matchedAmt7).
		Add(matchedAmt11).
		Add(matchedAmt12).
		Add(matchedAmt16).
		Add(matchedAmt17)

	s.Require().Equal(mInfo.TotalMatchedAmount, totalMatchedAmt)
	s.Require().Equal(mInfo.AllocationMap[s.addr(1).String()], matchedAmt1)
	s.Require().Equal(mInfo.AllocationMap[s.addr(2).String()], matchedAmt2)
	s.Require().Equal(mInfo.AllocationMap[s.addr(3).String()], matchedAmt3)
	s.Require().Equal(mInfo.AllocationMap[s.addr(4).String()], matchedAmt4)
	s.Require().Equal(mInfo.AllocationMap[s.addr(5).String()], matchedAmt5)
	s.Require().Equal(mInfo.AllocationMap[s.addr(6).String()], matchedAmt6)
	s.Require().Equal(mInfo.AllocationMap[s.addr(7).String()], matchedAmt7)
	s.Require().Equal(mInfo.AllocationMap[s.addr(8).String()], matchedAmt8)
	s.Require().Equal(mInfo.AllocationMap[s.addr(9).String()], matchedAmt9)
	s.Require().Equal(mInfo.AllocationMap[s.addr(10).String()], matchedAmt10)
	s.Require().Equal(mInfo.AllocationMap[s.addr(11).String()], matchedAmt11)
	s.Require().Equal(mInfo.AllocationMap[s.addr(12).String()], matchedAmt12)
	s.Require().Equal(mInfo.AllocationMap[s.addr(13).String()], matchedAmt13)
	s.Require().Equal(mInfo.AllocationMap[s.addr(14).String()], matchedAmt14)
	s.Require().Equal(mInfo.AllocationMap[s.addr(15).String()], matchedAmt15)
	s.Require().Equal(mInfo.AllocationMap[s.addr(16).String()], matchedAmt16)
	s.Require().Equal(mInfo.AllocationMap[s.addr(17).String()], matchedAmt17)

	reservedMatchedAmt1 := sdk.NewInt(0)
	reservedMatchedAmt2 := sdk.NewInt(1999_999_994)
	reservedMatchedAmt3 := matchedAmt3.ToDec().Mul(matchingPrice).Ceil().TruncateInt()
	reservedMatchedAmt4 := matchedAmt4.ToDec().Mul(matchingPrice).Ceil().TruncateInt()
	reservedMatchedAmt5 := matchedAmt5.ToDec().Mul(matchingPrice).Ceil().TruncateInt()
	reservedMatchedAmt6 := sdk.NewInt(2499_999_997)
	reservedMatchedAmt7 := matchedAmt7.ToDec().Mul(matchingPrice).Ceil().TruncateInt()
	reservedMatchedAmt8 := sdk.NewInt(0)
	reservedMatchedAmt9 := sdk.NewInt(0)
	reservedMatchedAmt10 := sdk.NewInt(0)
	reservedMatchedAmt11 := matchedAmt11.ToDec().Mul(matchingPrice).Ceil().TruncateInt()
	reservedMatchedAmt12 := sdk.NewInt(1499_999_990)
	reservedMatchedAmt13 := sdk.NewInt(0)
	reservedMatchedAmt14 := sdk.NewInt(0)
	reservedMatchedAmt15 := sdk.NewInt(0)
	reservedMatchedAmt16 := sdk.NewInt(2499_999_997)
	reservedMatchedAmt17 := matchedAmt17.ToDec().Mul(matchingPrice).Ceil().TruncateInt()

	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(1).String()], reservedMatchedAmt1)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(2).String()], reservedMatchedAmt2)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(3).String()], reservedMatchedAmt3)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(4).String()], reservedMatchedAmt4)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(5).String()], reservedMatchedAmt5)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(6).String()], reservedMatchedAmt6)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(7).String()], reservedMatchedAmt7)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(8).String()], reservedMatchedAmt8)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(9).String()], reservedMatchedAmt9)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(10).String()], reservedMatchedAmt10)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(11).String()], reservedMatchedAmt11)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(12).String()], reservedMatchedAmt12)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(13).String()], reservedMatchedAmt13)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(14).String()], reservedMatchedAmt14)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(15).String()], reservedMatchedAmt15)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(16).String()], reservedMatchedAmt16)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(17).String()], reservedMatchedAmt17)

	refundAmt1 := sdk.NewInt(200_000_000).ToDec().Mul(parseDec("10")).Ceil().TruncateInt()
	refundAmt2 := sdk.NewInt(2000_000_000).Sub(reservedMatchedAmt2)
	refundAmt3 := sdk.NewInt(500_000_000).ToDec().Mul(parseDec("10.5")).Ceil().TruncateInt().Sub(reservedMatchedAmt3)
	refundAmt4 := sdk.NewInt(220_000_000).ToDec().Mul(parseDec("10.8")).Ceil().TruncateInt().Add(sdk.NewInt(1500_000_000)).Sub(reservedMatchedAmt4)
	refundAmt5 := sdk.NewInt(300_000_000).ToDec().Mul(parseDec("10.8")).Ceil().TruncateInt().Add(sdk.NewInt(1500_000_000)).Sub(reservedMatchedAmt5)
	refundAmt6 := sdk.NewInt(250_000_000).ToDec().Mul(parseDec("9.7")).Ceil().TruncateInt().Add(sdk.NewInt(2500_000_000)).Sub(reservedMatchedAmt6)
	refundAmt7 := sdk.NewInt(100_000_000).ToDec().Mul(parseDec("11.3")).Ceil().TruncateInt().Sub(reservedMatchedAmt7)
	refundAmt8 := sdk.NewInt(2500_000_000)
	refundAmt9 := sdk.NewInt(300_000_000).ToDec().Mul(parseDec("10.1")).Ceil().TruncateInt()
	refundAmt10 := sdk.NewInt(2000_000_000)
	refundAmt11 := sdk.NewInt(150_000_000).ToDec().Mul(parseDec("10.75")).Ceil().TruncateInt().Sub(reservedMatchedAmt11)
	refundAmt12 := sdk.NewInt(1500_000_000).Sub(reservedMatchedAmt12)
	refundAmt13 := sdk.NewInt(200_000_000).ToDec().Mul(parseDec("10.2")).Ceil().TruncateInt()
	refundAmt14 := sdk.NewInt(2000_000_000)
	refundAmt15 := sdk.NewInt(200_000_000).ToDec().Mul(parseDec("10.25")).Ceil().TruncateInt()
	refundAmt16 := sdk.NewInt(2500_000_000).Sub(reservedMatchedAmt16)
	refundAmt17 := sdk.NewInt(180_000_000).ToDec().Mul(parseDec("10.52")).Ceil().TruncateInt().Sub(reservedMatchedAmt17)

	s.Require().True(mInfo.RefundMap[s.addr(1).String()].Equal(refundAmt1))
	s.Require().True(mInfo.RefundMap[s.addr(2).String()].Equal(refundAmt2))
	s.Require().True(mInfo.RefundMap[s.addr(3).String()].Equal(refundAmt3))
	s.Require().True(mInfo.RefundMap[s.addr(4).String()].Equal(refundAmt4))
	s.Require().True(mInfo.RefundMap[s.addr(5).String()].Equal(refundAmt5))
	s.Require().True(mInfo.RefundMap[s.addr(6).String()].Equal(refundAmt6))
	s.Require().True(mInfo.RefundMap[s.addr(7).String()].Equal(refundAmt7))
	s.Require().True(mInfo.RefundMap[s.addr(8).String()].Equal(refundAmt8))
	s.Require().True(mInfo.RefundMap[s.addr(9).String()].Equal(refundAmt9))
	s.Require().True(mInfo.RefundMap[s.addr(10).String()].Equal(refundAmt10))
	s.Require().True(mInfo.RefundMap[s.addr(11).String()].Equal(refundAmt11))
	s.Require().True(mInfo.RefundMap[s.addr(12).String()].Equal(refundAmt12))
	s.Require().True(mInfo.RefundMap[s.addr(13).String()].Equal(refundAmt13))
	s.Require().True(mInfo.RefundMap[s.addr(14).String()].Equal(refundAmt14))
	s.Require().True(mInfo.RefundMap[s.addr(15).String()].Equal(refundAmt15))
	s.Require().True(mInfo.RefundMap[s.addr(16).String()].Equal(refundAmt16))
	s.Require().True(mInfo.RefundMap[s.addr(17).String()].Equal(refundAmt17))

	// Distribute selling coin
	err := s.keeper.AllocateSellingCoin(s.ctx, auction, mInfo)
	s.Require().NoError(err)

	err = s.keeper.RefundRemainingSellingCoin(s.ctx, auction)
	s.Require().NoError(err)

	// The selling reserve account balance must be zero
	s.Require().True(s.getBalance(auction.GetSellingReserveAddress(), auction.SellingCoin.Denom).IsZero())

	// The auctioneer must have sellingCoin.Amount - TotalMatchedAmount
	s.Require().Equal(s.getBalance(s.addr(0), auction.GetSellingCoin().Denom).Amount, auction.SellingCoin.Amount.Sub(mInfo.TotalMatchedAmount))

	// The bidders must have the matched selling coin
	s.Require().True(s.getBalance(s.addr(1), auction.GetSellingCoin().Denom).Amount.Equal(matchedAmt1))
	s.Require().True(s.getBalance(s.addr(2), auction.GetSellingCoin().Denom).Amount.Equal(matchedAmt2))
	s.Require().True(s.getBalance(s.addr(3), auction.GetSellingCoin().Denom).Amount.Equal(matchedAmt3))
	s.Require().True(s.getBalance(s.addr(4), auction.GetSellingCoin().Denom).Amount.Equal(matchedAmt4))
	s.Require().True(s.getBalance(s.addr(5), auction.GetSellingCoin().Denom).Amount.Equal(matchedAmt5))
	s.Require().True(s.getBalance(s.addr(6), auction.GetSellingCoin().Denom).Amount.Equal(matchedAmt6))
	s.Require().True(s.getBalance(s.addr(7), auction.GetSellingCoin().Denom).Amount.Equal(matchedAmt7))
	s.Require().True(s.getBalance(s.addr(8), auction.GetSellingCoin().Denom).Amount.Equal(matchedAmt8))
	s.Require().True(s.getBalance(s.addr(9), auction.GetSellingCoin().Denom).Amount.Equal(matchedAmt9))
	s.Require().True(s.getBalance(s.addr(10), auction.GetSellingCoin().Denom).Amount.Equal(matchedAmt10))
	s.Require().True(s.getBalance(s.addr(11), auction.GetSellingCoin().Denom).Amount.Equal(matchedAmt11))
	s.Require().True(s.getBalance(s.addr(12), auction.GetSellingCoin().Denom).Amount.Equal(matchedAmt12))
	s.Require().True(s.getBalance(s.addr(13), auction.GetSellingCoin().Denom).Amount.Equal(matchedAmt13))
	s.Require().True(s.getBalance(s.addr(14), auction.GetSellingCoin().Denom).Amount.Equal(matchedAmt14))
	s.Require().True(s.getBalance(s.addr(15), auction.GetSellingCoin().Denom).Amount.Equal(matchedAmt15))
	s.Require().True(s.getBalance(s.addr(16), auction.GetSellingCoin().Denom).Amount.Equal(matchedAmt16))
	s.Require().True(s.getBalance(s.addr(17), auction.GetSellingCoin().Denom).Amount.Equal(matchedAmt17))

	// Refund payingCoin
	err = s.keeper.RefundPayingCoin(s.ctx, auction, mInfo)
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) TestCalculateAllocation_Mixed3_LimitedDifferent() {
	auction := s.createBatchAuction(
		s.addr(0),
		parseDec("10"),
		parseDec("0.1"),
		parseCoin("2500_000_000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		1,
		sdk.MustNewDecFromStr("0.2"),
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	s.placeBidBatchMany(auction.Id, s.addr(1), parseDec("10"), parseCoin("200_000_000denom1"), sdk.NewInt(500_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(2), parseDec("11"), parseCoin("2000_000_000denom2"), sdk.NewInt(500_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(3), parseDec("10.5"), parseCoin("500_000_000denom1"), sdk.NewInt(500_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(4), parseDec("10.2"), parseCoin("1500_000_000denom2"), sdk.NewInt(200_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(5), parseDec("10.8"), parseCoin("200_000_000denom1"), sdk.NewInt(200_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(6), parseDec("11.4"), parseCoin("2200_000_000denom2"), sdk.NewInt(200_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(7), parseDec("11.3"), parseCoin("100_000_000denom1"), sdk.NewInt(200_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(8), parseDec("9.9"), parseCoin("1900_000_000denom2"), sdk.NewInt(200_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(9), parseDec("10.1"), parseCoin("200_000_000denom1"), sdk.NewInt(200_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(10), parseDec("10.45"), parseCoin("2000_000_000denom2"), sdk.NewInt(200_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(11), parseDec("10.75"), parseCoin("100_000_000denom1"), sdk.NewInt(100_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(12), parseDec("10.99"), parseCoin("1050_000_000denom2"), sdk.NewInt(100_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(13), parseDec("10.2"), parseCoin("100_000_000denom1"), sdk.NewInt(100_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(14), parseDec("9.87"), parseCoin("980_000_000denom2"), sdk.NewInt(100_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(15), parseDec("10.25"), parseCoin("100_000_000denom1"), sdk.NewInt(100_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(16), parseDec("10.48"), parseCoin("1000_000_000denom2"), sdk.NewInt(100_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(17), parseDec("10.52"), parseCoin("100_000_000denom1"), sdk.NewInt(100_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(4), parseDec("10.8"), parseCoin("200_000_000denom1"), sdk.NewInt(0), true)
	s.placeBidBatchWorth(auction.Id, s.addr(5), parseDec("10.5"), parseCoin("1500_000_000denom2"), sdk.NewInt(0), true)
	s.placeBidBatchMany(auction.Id, s.addr(6), parseDec("9.7"), parseCoin("200_000_000denom1"), sdk.NewInt(0), true)

	a, found := s.keeper.GetAuction(s.ctx, auction.Id)
	s.Require().True(found)

	mInfo := s.keeper.CalculateBatchAllocation(s.ctx, a)

	// Checking
	s.Require().Equal(mInfo.MatchedLen, int64(16))
	matchingPrice := parseDec("10.1")
	s.Require().Equal(mInfo.MatchedPrice, matchingPrice)

	matchedAmt1 := sdk.NewInt(0)
	matchedAmt2 := sdk.NewInt(2000_000_000).ToDec().QuoTruncate(matchingPrice).TruncateInt()
	matchedAmt3 := sdk.NewInt(500_000_000)
	matchedAmt4 := sdk.NewInt(200_000_000)
	matchedAmt5 := sdk.NewInt(200_000_000)
	matchedAmt6 := sdk.NewInt(200_000_000)
	matchedAmt7 := sdk.NewInt(100_000_000)
	matchedAmt8 := sdk.NewInt(0)
	matchedAmt9 := sdk.NewInt(200_000_000)
	matchedAmt10 := sdk.NewInt(2000_000_000).ToDec().QuoTruncate(matchingPrice).TruncateInt()
	matchedAmt11 := sdk.NewInt(100_000_000)
	matchedAmt12 := sdk.NewInt(100_000_000)
	matchedAmt13 := sdk.NewInt(100_000_000)
	matchedAmt14 := sdk.NewInt(0)
	matchedAmt15 := sdk.NewInt(100_000_000)
	matchedAmt16 := sdk.NewInt(1000_000_000).ToDec().QuoTruncate(matchingPrice).TruncateInt()
	matchedAmt17 := sdk.NewInt(100_000_000)

	totalMatchedAmt := matchedAmt2.Add(matchedAmt3).
		Add(matchedAmt4).
		Add(matchedAmt5).
		Add(matchedAmt6).
		Add(matchedAmt7).
		Add(matchedAmt9).
		Add(matchedAmt10).
		Add(matchedAmt11).
		Add(matchedAmt12).
		Add(matchedAmt13).
		Add(matchedAmt15).
		Add(matchedAmt16).
		Add(matchedAmt17)

	s.Require().Equal(mInfo.TotalMatchedAmount, totalMatchedAmt)
	s.Require().Equal(mInfo.AllocationMap[s.addr(1).String()], matchedAmt1)
	s.Require().Equal(mInfo.AllocationMap[s.addr(2).String()], matchedAmt2)
	s.Require().Equal(mInfo.AllocationMap[s.addr(3).String()], matchedAmt3)
	s.Require().Equal(mInfo.AllocationMap[s.addr(4).String()], matchedAmt4)
	s.Require().Equal(mInfo.AllocationMap[s.addr(5).String()], matchedAmt5)
	s.Require().Equal(mInfo.AllocationMap[s.addr(6).String()], matchedAmt6)
	s.Require().Equal(mInfo.AllocationMap[s.addr(7).String()], matchedAmt7)
	s.Require().Equal(mInfo.AllocationMap[s.addr(8).String()], matchedAmt8)
	s.Require().Equal(mInfo.AllocationMap[s.addr(9).String()], matchedAmt9)
	s.Require().Equal(mInfo.AllocationMap[s.addr(10).String()], matchedAmt10)
	s.Require().Equal(mInfo.AllocationMap[s.addr(11).String()], matchedAmt11)
	s.Require().Equal(mInfo.AllocationMap[s.addr(12).String()], matchedAmt12)
	s.Require().Equal(mInfo.AllocationMap[s.addr(13).String()], matchedAmt13)
	s.Require().Equal(mInfo.AllocationMap[s.addr(14).String()], matchedAmt14)
	s.Require().Equal(mInfo.AllocationMap[s.addr(15).String()], matchedAmt15)
	s.Require().Equal(mInfo.AllocationMap[s.addr(16).String()], matchedAmt16)
	s.Require().Equal(mInfo.AllocationMap[s.addr(17).String()], matchedAmt17)

	reservedMatchedAmt1 := sdk.NewInt(0)
	reservedMatchedAmt2 := sdk.NewInt(1999_999_991)
	reservedMatchedAmt3 := sdk.NewInt(5050_000_000)
	reservedMatchedAmt4 := sdk.NewInt(2020_000_000)
	reservedMatchedAmt5 := sdk.NewInt(2020_000_000)
	reservedMatchedAmt6 := sdk.NewInt(2020_000_000)
	reservedMatchedAmt7 := sdk.NewInt(1010_000_000)
	reservedMatchedAmt8 := sdk.NewInt(0)
	reservedMatchedAmt9 := sdk.NewInt(2020_000_000)
	reservedMatchedAmt10 := sdk.NewInt(1999_999_991)
	reservedMatchedAmt11 := sdk.NewInt(1010_000_000)
	reservedMatchedAmt12 := sdk.NewInt(1010_000_000)
	reservedMatchedAmt13 := sdk.NewInt(1010_000_000)
	reservedMatchedAmt14 := sdk.NewInt(0)
	reservedMatchedAmt15 := sdk.NewInt(1010_000_000)
	reservedMatchedAmt16 := sdk.NewInt(999_999_990)
	reservedMatchedAmt17 := sdk.NewInt(1010_000_000)

	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(1).String()], reservedMatchedAmt1)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(2).String()], reservedMatchedAmt2)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(3).String()], reservedMatchedAmt3)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(4).String()], reservedMatchedAmt4)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(5).String()], reservedMatchedAmt5)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(6).String()], reservedMatchedAmt6)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(7).String()], reservedMatchedAmt7)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(8).String()], reservedMatchedAmt8)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(9).String()], reservedMatchedAmt9)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(10).String()], reservedMatchedAmt10)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(11).String()], reservedMatchedAmt11)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(12).String()], reservedMatchedAmt12)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(13).String()], reservedMatchedAmt13)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(14).String()], reservedMatchedAmt14)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(15).String()], reservedMatchedAmt15)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(16).String()], reservedMatchedAmt16)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(17).String()], reservedMatchedAmt17)

	refundAmt1 := sdk.NewInt(2000_000_000)
	refundAmt2 := sdk.NewInt(9)
	refundAmt3 := sdk.NewInt(200_000_000)
	refundAmt4 := sdk.NewInt(1640_000_000)
	refundAmt5 := sdk.NewInt(1640_000_000)
	refundAmt6 := sdk.NewInt(2120_000_000)
	refundAmt7 := sdk.NewInt(120_000_000)
	refundAmt8 := sdk.NewInt(1900_000_000)
	refundAmt9 := sdk.NewInt(0)
	refundAmt10 := sdk.NewInt(9)
	refundAmt11 := sdk.NewInt(65_000_000)
	refundAmt12 := sdk.NewInt(40_000_000)
	refundAmt13 := sdk.NewInt(10_000_000)
	refundAmt14 := sdk.NewInt(980_000_000)
	refundAmt15 := sdk.NewInt(15_000_000)
	refundAmt16 := sdk.NewInt(10)
	refundAmt17 := sdk.NewInt(42_000_000)

	s.Require().True(mInfo.RefundMap[s.addr(1).String()].Equal(refundAmt1))
	s.Require().True(mInfo.RefundMap[s.addr(2).String()].Equal(refundAmt2))
	s.Require().True(mInfo.RefundMap[s.addr(3).String()].Equal(refundAmt3))
	s.Require().True(mInfo.RefundMap[s.addr(4).String()].Equal(refundAmt4))
	s.Require().True(mInfo.RefundMap[s.addr(5).String()].Equal(refundAmt5))
	s.Require().True(mInfo.RefundMap[s.addr(6).String()].Equal(refundAmt6))
	s.Require().True(mInfo.RefundMap[s.addr(7).String()].Equal(refundAmt7))
	s.Require().True(mInfo.RefundMap[s.addr(8).String()].Equal(refundAmt8))
	s.Require().True(mInfo.RefundMap[s.addr(9).String()].Equal(refundAmt9))
	s.Require().True(mInfo.RefundMap[s.addr(10).String()].Equal(refundAmt10))
	s.Require().True(mInfo.RefundMap[s.addr(11).String()].Equal(refundAmt11))
	s.Require().True(mInfo.RefundMap[s.addr(12).String()].Equal(refundAmt12))
	s.Require().True(mInfo.RefundMap[s.addr(13).String()].Equal(refundAmt13))
	s.Require().True(mInfo.RefundMap[s.addr(14).String()].Equal(refundAmt14))
	s.Require().True(mInfo.RefundMap[s.addr(15).String()].Equal(refundAmt15))
	s.Require().True(mInfo.RefundMap[s.addr(16).String()].Equal(refundAmt16))
	s.Require().True(mInfo.RefundMap[s.addr(17).String()].Equal(refundAmt17))

	// Distribute selling coin
	err := s.keeper.AllocateSellingCoin(s.ctx, auction, mInfo)
	s.Require().NoError(err)

	err = s.keeper.RefundRemainingSellingCoin(s.ctx, auction)
	s.Require().NoError(err)

	// The selling reserve account balance must be zero
	s.Require().True(s.getBalance(auction.GetSellingReserveAddress(), auction.SellingCoin.Denom).IsZero())

	// The auctioneer must have sellingCoin.Amount - TotalMatchedAmount
	s.Require().Equal(s.getBalance(s.addr(0), auction.GetSellingCoin().Denom).Amount, auction.SellingCoin.Amount.Sub(mInfo.TotalMatchedAmount))

	// The bidders must have the matched selling coin
	s.Require().True(s.getBalance(s.addr(1), auction.GetSellingCoin().Denom).Amount.Equal(matchedAmt1))
	s.Require().True(s.getBalance(s.addr(2), auction.GetSellingCoin().Denom).Amount.Equal(matchedAmt2))
	s.Require().True(s.getBalance(s.addr(3), auction.GetSellingCoin().Denom).Amount.Equal(matchedAmt3))
	s.Require().True(s.getBalance(s.addr(4), auction.GetSellingCoin().Denom).Amount.Equal(matchedAmt4))
	s.Require().True(s.getBalance(s.addr(5), auction.GetSellingCoin().Denom).Amount.Equal(matchedAmt5))
	s.Require().True(s.getBalance(s.addr(6), auction.GetSellingCoin().Denom).Amount.Equal(matchedAmt6))
	s.Require().True(s.getBalance(s.addr(7), auction.GetSellingCoin().Denom).Amount.Equal(matchedAmt7))
	s.Require().True(s.getBalance(s.addr(8), auction.GetSellingCoin().Denom).Amount.Equal(matchedAmt8))
	s.Require().True(s.getBalance(s.addr(9), auction.GetSellingCoin().Denom).Amount.Equal(matchedAmt9))
	s.Require().True(s.getBalance(s.addr(10), auction.GetSellingCoin().Denom).Amount.Equal(matchedAmt10))
	s.Require().True(s.getBalance(s.addr(11), auction.GetSellingCoin().Denom).Amount.Equal(matchedAmt11))
	s.Require().True(s.getBalance(s.addr(12), auction.GetSellingCoin().Denom).Amount.Equal(matchedAmt12))
	s.Require().True(s.getBalance(s.addr(13), auction.GetSellingCoin().Denom).Amount.Equal(matchedAmt13))
	s.Require().True(s.getBalance(s.addr(14), auction.GetSellingCoin().Denom).Amount.Equal(matchedAmt14))
	s.Require().True(s.getBalance(s.addr(15), auction.GetSellingCoin().Denom).Amount.Equal(matchedAmt15))
	s.Require().True(s.getBalance(s.addr(16), auction.GetSellingCoin().Denom).Amount.Equal(matchedAmt16))
	s.Require().True(s.getBalance(s.addr(17), auction.GetSellingCoin().Denom).Amount.Equal(matchedAmt17))

	// Refund payingCoin
	err = s.keeper.RefundPayingCoin(s.ctx, auction, mInfo)
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) TestBenchmark_CalculateFixedPriceAllocation() {
	for _, numBids := range []int{100, 500, 1000, 1500, 2000} {
		s.SetupTest()

		auction := s.createFixedPriceAuction(
			s.addr(0),
			parseDec("0.5"),
			parseCoin("1_000_000_000_000denom1"),
			"denom2",
			[]types.VestingSchedule{},
			time.Now().AddDate(0, 0, -1),
			time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
			true,
		)
		s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

		for i := 1; i <= numBids; i++ {
			s.placeBidFixedPrice(auction.Id, s.addr(1), auction.GetStartPrice(), parseCoin("1_000_000denom2"), true)
		}

		a, found := s.keeper.GetAuction(s.ctx, auction.Id)
		s.Require().True(found)

		s.app.Commit()

		now := time.Now()
		s.keeper.CalculateFixedPriceAllocation(s.ctx, a)
		elapsed := time.Since(now)
		fmt.Printf("NumBids: %d, Elapsed Time: %s\n", numBids, elapsed)
	}
}

func (s *KeeperTestSuite) TestBenchmark_CalculateBatchAllocation() {
	// for _, numBids := range []int{100, 500, 1000, 1500, 2000} {
	for _, numBids := range []int{100} {
		s.SetupTest()

		auction := s.createBatchAuction(
			s.addr(0),
			parseDec("1"),
			parseDec("0.1"),
			parseCoin("1_000_000_000_000denom1"),
			"denom2",
			[]types.VestingSchedule{},
			1,
			sdk.MustNewDecFromStr("0.2"),
			time.Now().AddDate(0, 0, -1),
			time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
			true,
		)
		s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

		for i := 1; i <= numBids; i++ {
			price := parseDec("1")
			coin := parseCoin("1_000_000denom2")
			s.placeBidBatchWorth(auction.GetId(), s.addr(i), price, coin, parseInt("1_000_000_000"), true)
		}

		a, found := s.keeper.GetAuction(s.ctx, auction.Id)
		s.Require().True(found)

		s.app.Commit()

		now := time.Now()
		mInfo := s.keeper.CalculateBatchAllocation(s.ctx, a)
		fmt.Println(s.fullString(auction.Id, mInfo))
		elapsed := time.Since(now)
		fmt.Printf("NumBids: %d, Elapsed TIme: %s\n", numBids, elapsed)
	}
}
