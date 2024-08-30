package keeper_test

import (
	"fmt"
	"time"

	"cosmossdk.io/math"
	_ "github.com/stretchr/testify/suite"

	"github.com/tendermint/fundraising/x/fundraising/types"
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

	a, err := s.keeper.Auction.Get(s.ctx, auction.GetId())
	s.Require().NoError(err)

	s.placeBidFixedPrice(a.GetId(), s.addr(1), a.GetStartPrice(), parseCoin("15_000_000denom2"), true)
	s.placeBidFixedPrice(a.GetId(), s.addr(2), a.GetStartPrice(), parseCoin("20_000_000denom2"), true)
	s.placeBidFixedPrice(a.GetId(), s.addr(4), a.GetStartPrice(), parseCoin("10_000_000denom1"), true)
	s.placeBidFixedPrice(a.GetId(), s.addr(6), a.GetStartPrice(), parseCoin("20_000_000denom1"), true)

	mInfo, err := s.keeper.CalculateFixedPriceAllocation(s.ctx, a)
	s.Require().NoError(err)
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
		math.LegacyMustNewDecFromStr("0.2"),
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	s.placeBidBatchMany(auction.Id, s.addr(1), parseDec("1"), parseCoin("500_000_000denom1"), math.NewInt(1_000_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(2), parseDec("0.9"), parseCoin("500_000_000denom1"), math.NewInt(1_000_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(3), parseDec("0.8"), parseCoin("500_000_000denom1"), math.NewInt(1_000_000_000), true)

	a, err := s.keeper.Auction.Get(s.ctx, auction.Id)
	s.Require().NoError(err)

	mInfo, err := s.keeper.CalculateBatchAllocation(s.ctx, a)
	s.Require().NoError(err)

	s.Require().Equal(mInfo.MatchedLen, int64(2))
	s.Require().Equal(mInfo.MatchedPrice, parseDec("0.9"))
	s.Require().Equal(mInfo.TotalMatchedAmount, math.NewInt(1_000_000_000))
	s.Require().Equal(mInfo.AllocationMap[s.addr(1).String()], math.NewInt(500_000_000))
	s.Require().Equal(mInfo.AllocationMap[s.addr(2).String()], math.NewInt(500_000_000))
	s.Require().Equal(mInfo.AllocationMap[s.addr(3).String()], math.NewInt(0))
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(1).String()], math.NewInt(450_000_000))
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(2).String()], math.NewInt(450_000_000))
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(3).String()], math.NewInt(0))
	s.Require().Equal(mInfo.RefundMap[s.addr(1).String()], math.NewInt(50_000_000))
	s.Require().Equal(mInfo.RefundMap[s.addr(3).String()], math.NewInt(400_000_000))
	s.Require().True(mInfo.RefundMap[s.addr(2).String()].Equal(math.NewInt(0)))

	err = s.keeper.AllocateSellingCoin(s.ctx, auction, mInfo)
	s.Require().NoError(err)

	sellingReserveAmt := s.getBalance(auction.GetSellingReserveAddress(), auction.SellingCoin.Denom).Amount
	remainingAmt := auction.GetSellingCoin().Amount.Sub(mInfo.TotalMatchedAmount)
	s.Require().True(sellingReserveAmt.Equal(remainingAmt))

	err = s.keeper.RefundRemainingSellingCoin(s.ctx, auction)
	s.Require().NoError(err)
	s.Require().True(s.getBalance(auction.GetSellingReserveAddress(), auction.SellingCoin.Denom).IsZero())

	// The auctioneer must have sellingCoin.Amount - TotalMatchedAmount
	s.Require().Equal(s.getBalance(s.addr(0), auction.GetSellingCoin().Denom).Amount, math.NewInt(0))

	// The bidders must have the matched selling coin
	s.Require().Equal(s.getBalance(s.addr(1), auction.GetSellingCoin().Denom).Amount, math.NewInt(500_000_000))
	s.Require().Equal(s.getBalance(s.addr(2), auction.GetSellingCoin().Denom).Amount, math.NewInt(500_000_000))
	s.Require().Equal(s.getBalance(s.addr(3), auction.GetSellingCoin().Denom).Amount, math.NewInt(0))

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
		math.LegacyMustNewDecFromStr("0.2"),
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	s.placeBidBatchWorth(auction.Id, s.addr(1), parseDec("1"), parseCoin("500_000_000denom2"), math.NewInt(1500_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(2), parseDec("0.9"), parseCoin("500_000_000denom2"), math.NewInt(1500_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(3), parseDec("0.8"), parseCoin("500_000_000denom2"), math.NewInt(1500_000_000), true)

	a, err := s.keeper.Auction.Get(s.ctx, auction.Id)
	s.Require().NoError(err)

	mInfo, err := s.keeper.CalculateBatchAllocation(s.ctx, a)
	s.Require().NoError(err)

	// Checking
	s.Require().Equal(int64(2), mInfo.MatchedLen)
	s.Require().Equal(parseDec("0.9"), mInfo.MatchedPrice)
	matchingPrice := parseDec("0.9")
	matchedAmt := math.LegacyNewDec(500_000_000).QuoTruncate(matchingPrice).TruncateInt()

	s.Require().Equal(mInfo.TotalMatchedAmount, matchedAmt.Add(matchedAmt))
	s.Require().Equal(mInfo.AllocationMap[s.addr(1).String()], matchedAmt)
	s.Require().Equal(mInfo.AllocationMap[s.addr(2).String()], matchedAmt)
	s.Require().Equal(mInfo.AllocationMap[s.addr(3).String()], math.NewInt(0))
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(1).String()], math.NewInt(500_000_000))
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(2).String()], math.NewInt(500_000_000))
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(3).String()], math.NewInt(0))
	s.Require().True(mInfo.RefundMap[s.addr(1).String()].IsZero())
	s.Require().True(mInfo.RefundMap[s.addr(2).String()].IsZero())
	s.Require().Equal(mInfo.RefundMap[s.addr(3).String()], math.NewInt(500_000_000))

	// Distribute selling coin
	err = s.keeper.AllocateSellingCoin(s.ctx, auction, mInfo)
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
		math.LegacyMustNewDecFromStr("0.2"),
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	s.placeBidBatchMany(auction.Id, s.addr(1), parseDec("1"), parseCoin("500_000_000denom1"), math.NewInt(1500_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(2), parseDec("0.9"), parseCoin("500_000_000denom2"), math.NewInt(1500_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(3), parseDec("0.8"), parseCoin("500_000_000denom2"), math.NewInt(1500_000_000), true)

	a, err := s.keeper.Auction.Get(s.ctx, auction.Id)
	s.Require().NoError(err)

	mInfo, err := s.keeper.CalculateBatchAllocation(s.ctx, a)
	s.Require().NoError(err)

	// Checking
	s.Require().Equal(mInfo.MatchedLen, int64(2))
	s.Require().Equal(mInfo.MatchedPrice, parseDec("0.9"))
	matchingPrice := parseDec("0.9")
	matchedAmt1 := math.NewInt(500_000_000)
	matchedAmt2 := math.LegacyNewDec(500_000_000).QuoTruncate(matchingPrice).TruncateInt()

	s.Require().Equal(mInfo.TotalMatchedAmount, math.NewInt(500_000_000).Add(matchedAmt2))
	s.Require().Equal(mInfo.AllocationMap[s.addr(1).String()], matchedAmt1)
	s.Require().Equal(mInfo.AllocationMap[s.addr(2).String()], matchedAmt2)
	s.Require().Equal(mInfo.AllocationMap[s.addr(3).String()], math.NewInt(0))
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(1).String()], math.NewInt(450_000_000))
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(2).String()], math.NewInt(500_000_000))
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(3).String()], math.NewInt(0))
	s.Require().Equal(mInfo.RefundMap[s.addr(1).String()], math.NewInt(50_000_000))
	s.Require().Equal(mInfo.RefundMap[s.addr(3).String()], math.NewInt(500_000_000))
	s.Require().True(mInfo.RefundMap[s.addr(2).String()].IsZero())

	// Distribute selling coin
	err = s.keeper.AllocateSellingCoin(s.ctx, auction, mInfo)
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
		math.LegacyMustNewDecFromStr("0.2"),
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	s.placeBidBatchMany(auction.Id, s.addr(1), parseDec("1"), parseCoin("400_000_000denom1"), math.NewInt(400_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(2), parseDec("0.9"), parseCoin("400_000_000denom1"), math.NewInt(400_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(3), parseDec("0.8"), parseCoin("400_000_000denom1"), math.NewInt(400_000_000), true)

	a, err := s.keeper.Auction.Get(s.ctx, auction.Id)
	s.Require().NoError(err)

	mInfo, err := s.keeper.CalculateBatchAllocation(s.ctx, a)
	s.Require().NoError(err)

	// Checking
	s.Require().Equal(mInfo.MatchedLen, int64(2))
	s.Require().Equal(mInfo.MatchedPrice, parseDec("0.9"))
	s.Require().Equal(mInfo.TotalMatchedAmount, math.NewInt(800_000_000))
	s.Require().Equal(mInfo.AllocationMap[s.addr(1).String()], math.NewInt(400_000_000))
	s.Require().Equal(mInfo.AllocationMap[s.addr(2).String()], math.NewInt(400_000_000))
	s.Require().Equal(mInfo.AllocationMap[s.addr(3).String()], math.NewInt(0))
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(1).String()], math.NewInt(360_000_000))
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(2).String()], math.NewInt(360_000_000))
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(3).String()], math.NewInt(0))
	s.Require().Equal(mInfo.RefundMap[s.addr(1).String()], math.NewInt(40_000_000))
	s.Require().True(mInfo.RefundMap[s.addr(2).String()].IsZero())
	s.Require().Equal(mInfo.RefundMap[s.addr(3).String()], math.NewInt(320_000_000))

	// Distribute selling coin
	err = s.keeper.AllocateSellingCoin(s.ctx, auction, mInfo)
	s.Require().NoError(err)

	err = s.keeper.RefundRemainingSellingCoin(s.ctx, auction)
	s.Require().NoError(err)

	// The selling reserve account balance must be zero
	s.Require().True(s.getBalance(auction.GetSellingReserveAddress(), auction.SellingCoin.Denom).IsZero())

	// The auctioneer must have sellingCoin.Amount - TotalMatchedAmount
	s.Require().Equal(s.getBalance(s.addr(0), auction.GetSellingCoin().Denom).Amount, math.NewInt(200_000_000))

	// The bidders must have the matched selling coin
	s.Require().Equal(s.getBalance(s.addr(1), auction.GetSellingCoin().Denom).Amount, math.NewInt(400_000_000))
	s.Require().Equal(s.getBalance(s.addr(2), auction.GetSellingCoin().Denom).Amount, math.NewInt(400_000_000))
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
		math.LegacyMustNewDecFromStr("0.2"),
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	s.placeBidBatchWorth(auction.Id, s.addr(1), parseDec("1"), parseCoin("400_000_000denom2"), math.NewInt(400_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(2), parseDec("0.9"), parseCoin("360_000_000denom2"), math.NewInt(400_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(3), parseDec("0.8"), parseCoin("320_000_000denom2"), math.NewInt(400_000_000), true)

	a, err := s.keeper.Auction.Get(s.ctx, auction.Id)
	s.Require().NoError(err)

	mInfo, err := s.keeper.CalculateBatchAllocation(s.ctx, a)
	s.Require().NoError(err)

	// Checking
	s.Require().Equal(mInfo.MatchedLen, int64(3))
	s.Require().Equal(mInfo.MatchedPrice, parseDec("0.8"))
	s.Require().Equal(mInfo.TotalMatchedAmount, math.NewInt(1200_000_000))
	s.Require().Equal(mInfo.AllocationMap[s.addr(1).String()], math.NewInt(400_000_000))
	s.Require().Equal(mInfo.AllocationMap[s.addr(2).String()], math.NewInt(400_000_000))
	s.Require().Equal(mInfo.AllocationMap[s.addr(3).String()], math.NewInt(400_000_000))
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(1).String()], math.NewInt(320_000_000))
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(2).String()], math.NewInt(320_000_000))
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(3).String()], math.NewInt(320_000_000))
	s.Require().Equal(mInfo.RefundMap[s.addr(1).String()], math.NewInt(80_000_000))
	s.Require().Equal(mInfo.RefundMap[s.addr(2).String()], math.NewInt(40_000_000))
	s.Require().True(mInfo.RefundMap[s.addr(3).String()].IsZero())

	// Distribute selling coin
	err = s.keeper.AllocateSellingCoin(s.ctx, auction, mInfo)
	s.Require().NoError(err)

	err = s.keeper.RefundRemainingSellingCoin(s.ctx, auction)
	s.Require().NoError(err)

	// The selling reserve account balance must be zero
	s.Require().True(s.getBalance(auction.GetSellingReserveAddress(), auction.SellingCoin.Denom).IsZero())

	// The auctioneer must have sellingCoin.Amount - TotalMatchedAmount
	s.Require().Equal(s.getBalance(s.addr(0), auction.GetSellingCoin().Denom).Amount, math.NewInt(300_000_000))

	// The bidders must have the matched selling coin
	s.Require().Equal(s.getBalance(s.addr(1), auction.GetSellingCoin().Denom).Amount, math.NewInt(400_000_000))
	s.Require().Equal(s.getBalance(s.addr(2), auction.GetSellingCoin().Denom).Amount, math.NewInt(400_000_000))
	s.Require().Equal(s.getBalance(s.addr(3), auction.GetSellingCoin().Denom).Amount, math.NewInt(400_000_000))

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
		math.LegacyMustNewDecFromStr("0.2"),
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	s.placeBidBatchMany(auction.Id, s.addr(1), parseDec("1"), parseCoin("500_000_000denom1"), math.NewInt(600_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(2), parseDec("0.9"), parseCoin("500_000_000denom2"), math.NewInt(600_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(3), parseDec("0.8"), parseCoin("450_000_000denom2"), math.NewInt(600_000_000), true)

	a, err := s.keeper.Auction.Get(s.ctx, auction.Id)
	s.Require().NoError(err)

	mInfo, err := s.keeper.CalculateBatchAllocation(s.ctx, a)
	s.Require().NoError(err)

	// Checking
	s.Require().Equal(mInfo.MatchedLen, int64(3))
	s.Require().Equal(mInfo.MatchedPrice, parseDec("0.8"))
	s.Require().Equal(mInfo.TotalMatchedAmount, math.NewInt(1662_500_000))
	s.Require().Equal(mInfo.AllocationMap[s.addr(1).String()], math.NewInt(500_000_000))
	s.Require().Equal(mInfo.AllocationMap[s.addr(2).String()], math.NewInt(600_000_000))
	s.Require().Equal(mInfo.AllocationMap[s.addr(3).String()], math.NewInt(562_500_000))
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(1).String()], math.NewInt(400_000_000))
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(2).String()], math.NewInt(480_000_000))
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(3).String()], math.NewInt(450_000_000))
	s.Require().Equal(mInfo.RefundMap[s.addr(1).String()], math.NewInt(100_000_000))
	s.Require().Equal(mInfo.RefundMap[s.addr(2).String()], math.NewInt(20_000_000))
	s.Require().True(mInfo.RefundMap[s.addr(3).String()].IsZero())

	// Distribute selling coin
	err = s.keeper.AllocateSellingCoin(s.ctx, auction, mInfo)
	s.Require().NoError(err)

	err = s.keeper.RefundRemainingSellingCoin(s.ctx, auction)
	s.Require().NoError(err)

	// The selling reserve account balance must be zero
	s.Require().True(s.getBalance(auction.GetSellingReserveAddress(), auction.SellingCoin.Denom).IsZero())

	// The auctioneer must have sellingCoin.Amount - TotalMatchedAmount
	s.Require().True(
		s.getBalance(s.addr(0), auction.GetSellingCoin().Denom).Amount.
			Equal(math.NewInt(37_500_000)),
	)

	// The bidders must have the matched selling coin
	s.Require().Equal(s.getBalance(s.addr(1), auction.GetSellingCoin().Denom).Amount, math.NewInt(500_000_000))
	s.Require().Equal(s.getBalance(s.addr(2), auction.GetSellingCoin().Denom).Amount, math.NewInt(600_000_000))
	s.Require().Equal(s.getBalance(s.addr(3), auction.GetSellingCoin().Denom).Amount, math.NewInt(562_500_000))

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
		math.LegacyMustNewDecFromStr("0.2"),
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	s.placeBidBatchMany(auction.Id, s.addr(1), parseDec("1"), parseCoin("200_000_000denom1"), math.NewInt(5000_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(2), parseDec("0.8"), parseCoin("500_000_000denom2"), math.NewInt(5000_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(3), parseDec("0.9"), parseCoin("500_000_000denom1"), math.NewInt(5000_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(1), parseDec("1.1"), parseCoin("300_000_000denom2"), math.NewInt(0), true)
	s.placeBidBatchMany(auction.Id, s.addr(5), parseDec("1.2"), parseCoin("300_000_000denom1"), math.NewInt(5000_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(4), parseDec("0.8"), parseCoin("100_000_000denom1"), math.NewInt(5000_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(2), parseDec("0.7"), parseCoin("100_000_000denom1"), math.NewInt(0), true)
	s.placeBidBatchMany(auction.Id, s.addr(6), parseDec("0.5"), parseCoin("100_000_000denom1"), math.NewInt(5000_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(3), parseDec("0.8"), parseCoin("300_000_000denom2"), math.NewInt(0), true)
	s.placeBidBatchWorth(auction.Id, s.addr(7), parseDec("0.6"), parseCoin("500_000_000denom2"), math.NewInt(5000_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(8), parseDec("0.8"), parseCoin("500_000_000denom1"), math.NewInt(5000_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(9), parseDec("0.6"), parseCoin("600_000_000denom1"), math.NewInt(5000_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(6), parseDec("0.5"), parseCoin("500_000_000denom2"), math.NewInt(0), true)
	s.placeBidBatchMany(auction.Id, s.addr(10), parseDec("0.6"), parseCoin("100_000_000denom1"), math.NewInt(5000_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(3), parseDec("0.7"), parseCoin("800_000_000denom2"), math.NewInt(0), true)

	a, err := s.keeper.Auction.Get(s.ctx, auction.Id)
	s.Require().NoError(err)

	mInfo, err := s.keeper.CalculateBatchAllocation(s.ctx, a)
	s.Require().NoError(err)

	// Checking
	s.Require().Equal(mInfo.MatchedLen, int64(10))
	matchingPrice := parseDec("0.7")
	s.Require().Equal(mInfo.MatchedPrice, matchingPrice)

	matchedAmt1 := math.LegacyNewDec(300_000_000).QuoTruncate(matchingPrice).TruncateInt().Add(math.NewInt(200_000_000))
	matchedAmt2 := math.LegacyNewDec(500_000_000).QuoTruncate(matchingPrice).TruncateInt().Add(math.NewInt(100_000_000))
	tMatchedAmt3 := math.LegacyNewDec(300_000_000).QuoTruncate(matchingPrice).TruncateInt().Add(math.NewInt(500_000_000))
	matchedAmt3 := tMatchedAmt3.Add(math.LegacyNewDec(800_000_000).QuoTruncate(matchingPrice).TruncateInt())
	matchedAmt4 := math.NewInt(100_000_000)
	matchedAmt5 := math.NewInt(300_000_000)
	matchedAmt8 := math.NewInt(500_000_000)
	matchedAmt_Zero := math.NewInt(0)
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

	reservedmatchedAmt1 := math.LegacyNewDec(200_000_000).Mul(matchingPrice).Ceil().TruncateInt().Add(math.NewInt(300_000_000))
	reservedMatchedAmt2 := math.LegacyNewDec(100_000_000).Mul(matchingPrice).Ceil().TruncateInt().Add(math.NewInt(500_000_000))
	reservedMatchedAmt3 := math.LegacyNewDec(500_000_000).Mul(matchingPrice).Ceil().TruncateInt().Add(math.NewInt(1100_000_000))
	reservedMatchedAmt4 := math.LegacyNewDec(100_000_000).Mul(matchingPrice).Ceil().TruncateInt()
	reservedMatchedAmt5 := math.LegacyNewDec(300_000_000).Mul(matchingPrice).Ceil().TruncateInt()
	reservedMatchedAmt8 := math.LegacyNewDec(500_000_000).Mul(matchingPrice).Ceil().TruncateInt()
	reservedMatchedAmt_Zero := math.NewInt(0)

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

	refundAmt1 := math.NewInt(60_000_000)
	refundAmt2 := math.NewInt(0)
	refundAmt3 := math.NewInt(100_000_000)
	refundAmt4 := math.NewInt(10_000_000)
	refundAmt5 := math.NewInt(150_000_000)
	refundAmt6 := math.NewInt(550_000_000)
	refundAmt7 := math.NewInt(500_000_000)
	refundAmt8 := math.NewInt(50_000_000)
	refundAmt9 := math.NewInt(360_000_000)
	refundAmt10 := math.NewInt(60_000_000)

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
	err = s.keeper.AllocateSellingCoin(s.ctx, auction, mInfo)
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
		math.LegacyMustNewDecFromStr("0.2"),
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	s.placeBidBatchMany(auction.Id, s.addr(1), parseDec("1"), parseCoin("200_000_000denom1"), math.NewInt(700_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(2), parseDec("0.8"), parseCoin("500_000_000denom2"), math.NewInt(700_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(3), parseDec("0.9"), parseCoin("500_000_000denom1"), math.NewInt(700_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(1), parseDec("1.1"), parseCoin("300_000_000denom2"), math.NewInt(0), true)
	s.placeBidBatchMany(auction.Id, s.addr(5), parseDec("1.2"), parseCoin("300_000_000denom1"), math.NewInt(700_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(4), parseDec("0.8"), parseCoin("100_000_000denom1"), math.NewInt(700_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(2), parseDec("0.7"), parseCoin("100_000_000denom1"), math.NewInt(0), true)
	s.placeBidBatchMany(auction.Id, s.addr(6), parseDec("0.5"), parseCoin("100_000_000denom1"), math.NewInt(700_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(3), parseDec("0.8"), parseCoin("300_000_000denom2"), math.NewInt(0), true)
	s.placeBidBatchWorth(auction.Id, s.addr(7), parseDec("0.6"), parseCoin("400_000_000denom2"), math.NewInt(700_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(8), parseDec("0.8"), parseCoin("500_000_000denom1"), math.NewInt(700_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(9), parseDec("0.6"), parseCoin("600_000_000denom1"), math.NewInt(700_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(6), parseDec("0.5"), parseCoin("350_000_000denom2"), math.NewInt(0), true)
	s.placeBidBatchMany(auction.Id, s.addr(10), parseDec("0.6"), parseCoin("100_000_000denom1"), math.NewInt(700_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(3), parseDec("0.7"), parseCoin("490_000_000denom2"), math.NewInt(0), true)

	a, err := s.keeper.Auction.Get(s.ctx, auction.Id)
	s.Require().NoError(err)

	mInfo, err := s.keeper.CalculateBatchAllocation(s.ctx, a)
	s.Require().NoError(err)

	// Checking
	s.Require().Equal(int64(11), mInfo.MatchedLen)
	matchingPrice := parseDec("0.6")
	s.Require().Equal(mInfo.MatchedPrice, matchingPrice)

	matchedAmt1 := math.NewInt(700_000_000)
	matchedAmt2 := math.NewInt(700_000_000)
	matchedAmt3 := math.NewInt(700_000_000)
	matchedAmt4 := math.NewInt(100_000_000)
	matchedAmt5 := math.NewInt(300_000_000)
	matchedAmt6 := math.NewInt(0)
	matchedAmt7 := math.LegacyNewDec(400_000_000).QuoTruncate(matchingPrice).TruncateInt()
	matchedAmt8 := math.NewInt(500_000_000)
	matchedAmt9 := math.NewInt(600_000_000)
	matchedAmt10 := math.NewInt(100_000_000)

	totalMatchedAmt := math.NewInt(3700_000_000).Add(matchedAmt7)

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

	reservedMatchedAmt1 := math.NewInt(420_000_000)
	reservedMatchedAmt2 := math.NewInt(420_000_000)
	reservedMatchedAmt3 := math.NewInt(420_000_000)
	reservedMatchedAmt4 := math.NewInt(60_000_000)
	reservedMatchedAmt5 := math.NewInt(180_000_000)
	reservedMatchedAmt6 := math.NewInt(0)
	reservedMatchedAmt7 := math.NewInt(400_000_000)
	reservedMatchedAmt8 := math.NewInt(300_000_000)
	reservedMatchedAmt9 := math.NewInt(360_000_000)
	reservedMatchedAmt10 := math.NewInt(60_000_000)

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

	refundAmt1 := math.NewInt(80_000_000)
	refundAmt2 := math.NewInt(150_000_000)
	refundAmt3 := math.NewInt(820_000_000)
	refundAmt4 := math.NewInt(20_000_000)
	refundAmt5 := math.NewInt(180_000_000)
	refundAmt6 := math.NewInt(400_000_000)
	refundAmt7 := math.NewInt(0)
	refundAmt8 := math.NewInt(100_000_000)
	refundAmt9 := math.NewInt(0)
	refundAmt10 := math.NewInt(0)

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
	err = s.keeper.AllocateSellingCoin(s.ctx, auction, mInfo)
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
		math.LegacyMustNewDecFromStr("0.2"),
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	s.placeBidBatchMany(auction.Id, s.addr(1), parseDec("1"), parseCoin("200_000_000denom1"), math.NewInt(1000_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(2), parseDec("0.8"), parseCoin("500_000_000denom2"), math.NewInt(1000_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(3), parseDec("0.9"), parseCoin("500_000_000denom1"), math.NewInt(800_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(1), parseDec("1.1"), parseCoin("300_000_000denom2"), math.NewInt(0), true)
	s.placeBidBatchMany(auction.Id, s.addr(5), parseDec("1.2"), parseCoin("300_000_000denom1"), math.NewInt(600_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(4), parseDec("0.8"), parseCoin("100_000_000denom1"), math.NewInt(800_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(2), parseDec("0.7"), parseCoin("100_000_000denom1"), math.NewInt(0), true)
	s.placeBidBatchMany(auction.Id, s.addr(6), parseDec("0.5"), parseCoin("100_000_000denom1"), math.NewInt(600_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(3), parseDec("0.8"), parseCoin("300_000_000denom2"), math.NewInt(0), true)
	s.placeBidBatchWorth(auction.Id, s.addr(7), parseDec("0.6"), parseCoin("200_000_000denom2"), math.NewInt(400_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(8), parseDec("0.8"), parseCoin("400_000_000denom1"), math.NewInt(400_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(9), parseDec("0.6"), parseCoin("200_000_000denom1"), math.NewInt(200_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(6), parseDec("0.5"), parseCoin("300_000_000denom2"), math.NewInt(0), true)
	s.placeBidBatchMany(auction.Id, s.addr(10), parseDec("0.6"), parseCoin("100_000_000denom1"), math.NewInt(200_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(3), parseDec("0.7"), parseCoin("560_000_000denom2"), math.NewInt(0), true)

	a, err := s.keeper.Auction.Get(s.ctx, auction.Id)
	s.Require().NoError(err)

	mInfo, err := s.keeper.CalculateBatchAllocation(s.ctx, a)
	s.Require().NoError(err)

	// Checking
	s.Require().Equal(int64(13), mInfo.MatchedLen)
	matchingPrice := parseDec("0.5")
	s.Require().Equal(mInfo.MatchedPrice, matchingPrice)

	matchedAmt1 := math.NewInt(800_000_000)
	matchedAmt2 := math.NewInt(1000_000_000)
	matchedAmt3 := math.NewInt(800_000_000)
	matchedAmt4 := math.NewInt(100_000_000)
	matchedAmt5 := math.NewInt(300_000_000)
	matchedAmt6 := math.NewInt(600_000_000)
	matchedAmt7 := math.NewInt(400_000_000)
	matchedAmt8 := math.NewInt(400_000_000)
	matchedAmt9 := math.NewInt(200_000_000)
	matchedAmt10 := math.NewInt(100_000_000)

	totalMatchedAmt := math.NewInt(4700_000_000)

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

	reservedMatchedAmt1 := math.NewInt(400_000_000)
	reservedMatchedAmt2 := math.NewInt(500_000_000)
	reservedMatchedAmt3 := math.NewInt(400_000_000)
	reservedMatchedAmt4 := math.NewInt(50_000_000)
	reservedMatchedAmt5 := math.NewInt(150_000_000)
	reservedMatchedAmt6 := math.NewInt(300_000_000)
	reservedMatchedAmt7 := math.NewInt(200_000_000)
	reservedMatchedAmt8 := math.NewInt(200_000_000)
	reservedMatchedAmt9 := math.NewInt(100_000_000)
	reservedMatchedAmt10 := math.NewInt(50_000_000)

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

	refundAmt1 := math.NewInt(100_000_000)
	refundAmt2 := math.NewInt(70_000_000)
	refundAmt3 := math.NewInt(910_000_000)
	refundAmt4 := math.NewInt(30_000_000)
	refundAmt5 := math.NewInt(210_000_000)
	refundAmt6 := math.NewInt(50_000_000)
	refundAmt7 := math.NewInt(0)
	refundAmt8 := math.NewInt(120_000_000)
	refundAmt9 := math.NewInt(20_000_000)
	refundAmt10 := math.NewInt(10_000_000)

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
	err = s.keeper.AllocateSellingCoin(s.ctx, auction, mInfo)
	s.Require().NoError(err)

	err = s.keeper.RefundRemainingSellingCoin(s.ctx, auction)
	s.Require().NoError(err)

	// The selling reserve account balance must be zero
	s.Require().True(s.getBalance(auction.GetSellingReserveAddress(), auction.SellingCoin.Denom).IsZero())

	// The auctioneer must have sellingCoin.Amount - TotalMatchedAmount
	s.Require().Equal(s.getBalance(s.addr(0), auction.GetSellingCoin().Denom).Amount, math.NewInt(300_000_000))

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
		math.LegacyMustNewDecFromStr("0.2"),
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	s.placeBidBatchMany(auction.Id, s.addr(1), parseDec("10"), parseCoin("200_000_000denom1"), math.NewInt(2500_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(2), parseDec("11"), parseCoin("2000_000_000denom2"), math.NewInt(2500_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(3), parseDec("10.5"), parseCoin("500_000_000denom1"), math.NewInt(2500_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(4), parseDec("10.2"), parseCoin("1500_000_000denom2"), math.NewInt(2500_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(5), parseDec("10.8"), parseCoin("300_000_000denom1"), math.NewInt(2500_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(6), parseDec("11.4"), parseCoin("2500_000_000denom2"), math.NewInt(2500_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(7), parseDec("11.3"), parseCoin("100_000_000denom1"), math.NewInt(2500_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(8), parseDec("9.9"), parseCoin("2500_000_000denom2"), math.NewInt(2500_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(9), parseDec("10.1"), parseCoin("300_000_000denom1"), math.NewInt(2500_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(10), parseDec("10.45"), parseCoin("2000_000_000denom2"), math.NewInt(2500_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(11), parseDec("10.75"), parseCoin("150_000_000denom1"), math.NewInt(2500_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(12), parseDec("10.99"), parseCoin("1500_000_000denom2"), math.NewInt(2500_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(13), parseDec("10.2"), parseCoin("200_000_000denom1"), math.NewInt(2500_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(14), parseDec("9.87"), parseCoin("2000_000_000denom2"), math.NewInt(2500_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(15), parseDec("10.25"), parseCoin("200_000_000denom1"), math.NewInt(2500_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(16), parseDec("10.48"), parseCoin("2500_000_000denom2"), math.NewInt(2500_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(17), parseDec("10.52"), parseCoin("180_000_000denom1"), math.NewInt(2500_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(4), parseDec("10.8"), parseCoin("220_000_000denom1"), math.NewInt(0), true)
	s.placeBidBatchWorth(auction.Id, s.addr(5), parseDec("10.5"), parseCoin("1500_000_000denom2"), math.NewInt(0), true)
	s.placeBidBatchMany(auction.Id, s.addr(6), parseDec("9.7"), parseCoin("250_000_000denom1"), math.NewInt(0), true)

	a, err := s.keeper.Auction.Get(s.ctx, auction.Id)
	s.Require().NoError(err)

	mInfo, err := s.keeper.CalculateBatchAllocation(s.ctx, a)
	s.Require().NoError(err)

	// Checking
	s.Require().Equal(mInfo.MatchedLen, int64(11))
	matchingPrice := parseDec("10.48")
	s.Require().Equal(mInfo.MatchedPrice, matchingPrice)

	matchedAmt1 := math.NewInt(0)
	matchedAmt2 := math.LegacyNewDec(2000_000_000).QuoTruncate(matchingPrice).TruncateInt()
	matchedAmt3 := math.NewInt(500_000_000)
	matchedAmt4 := math.NewInt(220_000_000)
	matchedAmt5 := math.LegacyNewDec(1500_000_000).QuoTruncate(matchingPrice).TruncateInt().Add(math.NewInt(300_000_000))
	matchedAmt6 := math.LegacyNewDec(2500_000_000).QuoTruncate(matchingPrice).TruncateInt()
	matchedAmt7 := math.NewInt(100_000_000)
	matchedAmt8 := math.NewInt(0)
	matchedAmt9 := math.NewInt(0)
	matchedAmt10 := math.NewInt(0)
	matchedAmt11 := math.NewInt(150_000_000)
	matchedAmt12 := math.LegacyNewDec(1500_000_000).QuoTruncate(matchingPrice).TruncateInt()
	matchedAmt13 := math.NewInt(0)
	matchedAmt14 := math.NewInt(0)
	matchedAmt15 := math.NewInt(0)
	matchedAmt16 := math.LegacyNewDec(2500_000_000).QuoTruncate(matchingPrice).TruncateInt()
	matchedAmt17 := math.NewInt(180_000_000)

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

	reservedMatchedAmt1 := math.NewInt(0)
	reservedMatchedAmt2 := math.NewInt(1999_999_994)
	reservedMatchedAmt3 := math.LegacyNewDecFromInt(matchedAmt3).Mul(matchingPrice).Ceil().TruncateInt()
	reservedMatchedAmt4 := math.LegacyNewDecFromInt(matchedAmt4).Mul(matchingPrice).Ceil().TruncateInt()
	reservedMatchedAmt5 := math.LegacyNewDecFromInt(matchedAmt5).Mul(matchingPrice).Ceil().TruncateInt()
	reservedMatchedAmt6 := math.NewInt(2499_999_997)
	reservedMatchedAmt7 := math.LegacyNewDecFromInt(matchedAmt7).Mul(matchingPrice).Ceil().TruncateInt()
	reservedMatchedAmt8 := math.NewInt(0)
	reservedMatchedAmt9 := math.NewInt(0)
	reservedMatchedAmt10 := math.NewInt(0)
	reservedMatchedAmt11 := math.LegacyNewDecFromInt(matchedAmt11).Mul(matchingPrice).Ceil().TruncateInt()
	reservedMatchedAmt12 := math.NewInt(1499_999_990)
	reservedMatchedAmt13 := math.NewInt(0)
	reservedMatchedAmt14 := math.NewInt(0)
	reservedMatchedAmt15 := math.NewInt(0)
	reservedMatchedAmt16 := math.NewInt(2499_999_997)
	reservedMatchedAmt17 := math.LegacyNewDecFromInt(matchedAmt17).Mul(matchingPrice).Ceil().TruncateInt()

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

	refundAmt1 := math.LegacyNewDec(200_000_000).Mul(parseDec("10")).Ceil().TruncateInt()
	refundAmt2 := math.NewInt(2000_000_000).Sub(reservedMatchedAmt2)
	refundAmt3 := math.LegacyNewDec(500_000_000).Mul(parseDec("10.5")).Ceil().TruncateInt().Sub(reservedMatchedAmt3)
	refundAmt4 := math.LegacyNewDec(220_000_000).Mul(parseDec("10.8")).Ceil().TruncateInt().Add(math.NewInt(1500_000_000)).Sub(reservedMatchedAmt4)
	refundAmt5 := math.LegacyNewDec(300_000_000).Mul(parseDec("10.8")).Ceil().TruncateInt().Add(math.NewInt(1500_000_000)).Sub(reservedMatchedAmt5)
	refundAmt6 := math.LegacyNewDec(250_000_000).Mul(parseDec("9.7")).Ceil().TruncateInt().Add(math.NewInt(2500_000_000)).Sub(reservedMatchedAmt6)
	refundAmt7 := math.LegacyNewDec(100_000_000).Mul(parseDec("11.3")).Ceil().TruncateInt().Sub(reservedMatchedAmt7)
	refundAmt8 := math.NewInt(2500_000_000)
	refundAmt9 := math.LegacyNewDec(300_000_000).Mul(parseDec("10.1")).Ceil().TruncateInt()
	refundAmt10 := math.NewInt(2000_000_000)
	refundAmt11 := math.LegacyNewDec(150_000_000).Mul(parseDec("10.75")).Ceil().TruncateInt().Sub(reservedMatchedAmt11)
	refundAmt12 := math.NewInt(1500_000_000).Sub(reservedMatchedAmt12)
	refundAmt13 := math.LegacyNewDec(200_000_000).Mul(parseDec("10.2")).Ceil().TruncateInt()
	refundAmt14 := math.NewInt(2000_000_000)
	refundAmt15 := math.LegacyNewDec(200_000_000).Mul(parseDec("10.25")).Ceil().TruncateInt()
	refundAmt16 := math.NewInt(2500_000_000).Sub(reservedMatchedAmt16)
	refundAmt17 := math.LegacyNewDec(180_000_000).Mul(parseDec("10.52")).Ceil().TruncateInt().Sub(reservedMatchedAmt17)

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
	err = s.keeper.AllocateSellingCoin(s.ctx, auction, mInfo)
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
		math.LegacyMustNewDecFromStr("0.2"),
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	s.placeBidBatchMany(auction.Id, s.addr(1), parseDec("10"), parseCoin("200_000_000denom1"), math.NewInt(500_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(2), parseDec("11"), parseCoin("2000_000_000denom2"), math.NewInt(500_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(3), parseDec("10.5"), parseCoin("500_000_000denom1"), math.NewInt(500_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(4), parseDec("10.2"), parseCoin("1500_000_000denom2"), math.NewInt(200_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(5), parseDec("10.8"), parseCoin("200_000_000denom1"), math.NewInt(200_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(6), parseDec("11.4"), parseCoin("2200_000_000denom2"), math.NewInt(200_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(7), parseDec("11.3"), parseCoin("100_000_000denom1"), math.NewInt(200_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(8), parseDec("9.9"), parseCoin("1900_000_000denom2"), math.NewInt(200_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(9), parseDec("10.1"), parseCoin("200_000_000denom1"), math.NewInt(200_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(10), parseDec("10.45"), parseCoin("2000_000_000denom2"), math.NewInt(200_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(11), parseDec("10.75"), parseCoin("100_000_000denom1"), math.NewInt(100_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(12), parseDec("10.99"), parseCoin("1050_000_000denom2"), math.NewInt(100_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(13), parseDec("10.2"), parseCoin("100_000_000denom1"), math.NewInt(100_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(14), parseDec("9.87"), parseCoin("980_000_000denom2"), math.NewInt(100_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(15), parseDec("10.25"), parseCoin("100_000_000denom1"), math.NewInt(100_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(16), parseDec("10.48"), parseCoin("1000_000_000denom2"), math.NewInt(100_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(17), parseDec("10.52"), parseCoin("100_000_000denom1"), math.NewInt(100_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(4), parseDec("10.8"), parseCoin("200_000_000denom1"), math.NewInt(0), true)
	s.placeBidBatchWorth(auction.Id, s.addr(5), parseDec("10.5"), parseCoin("1500_000_000denom2"), math.NewInt(0), true)
	s.placeBidBatchMany(auction.Id, s.addr(6), parseDec("9.7"), parseCoin("200_000_000denom1"), math.NewInt(0), true)

	a, err := s.keeper.Auction.Get(s.ctx, auction.Id)
	s.Require().NoError(err)

	mInfo, err := s.keeper.CalculateBatchAllocation(s.ctx, a)
	s.Require().NoError(err)

	// Checking
	s.Require().Equal(int64(14), mInfo.MatchedLen)
	matchingPrice := parseDec("10.1")
	s.Require().Equal(mInfo.MatchedPrice, matchingPrice)

	matchedAmt1 := math.NewInt(0)
	matchedAmt2 := math.LegacyNewDec(2000_000_000).QuoTruncate(matchingPrice).TruncateInt()
	matchedAmt3 := math.NewInt(500_000_000)
	matchedAmt4 := math.NewInt(200_000_000)
	matchedAmt5 := math.NewInt(200_000_000)
	matchedAmt6 := math.NewInt(200_000_000)
	matchedAmt7 := math.NewInt(100_000_000)
	matchedAmt8 := math.NewInt(0)
	matchedAmt9 := math.NewInt(200_000_000)
	matchedAmt10 := math.LegacyNewDec(2000_000_000).QuoTruncate(matchingPrice).TruncateInt()
	matchedAmt11 := math.NewInt(100_000_000)
	matchedAmt12 := math.NewInt(100_000_000)
	matchedAmt13 := math.NewInt(100_000_000)
	matchedAmt14 := math.NewInt(0)
	matchedAmt15 := math.NewInt(100_000_000)
	matchedAmt16 := math.LegacyNewDec(1000_000_000).QuoTruncate(matchingPrice).TruncateInt()
	matchedAmt17 := math.NewInt(100_000_000)

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

	reservedMatchedAmt1 := math.NewInt(0)
	reservedMatchedAmt2 := math.NewInt(1999_999_991)
	reservedMatchedAmt3 := math.NewInt(5050_000_000)
	reservedMatchedAmt4 := math.NewInt(2020_000_000)
	reservedMatchedAmt5 := math.NewInt(2020_000_000)
	reservedMatchedAmt6 := math.NewInt(2020_000_000)
	reservedMatchedAmt7 := math.NewInt(1010_000_000)
	reservedMatchedAmt8 := math.NewInt(0)
	reservedMatchedAmt9 := math.NewInt(2020_000_000)
	reservedMatchedAmt10 := math.NewInt(1999_999_991)
	reservedMatchedAmt11 := math.NewInt(1010_000_000)
	reservedMatchedAmt12 := math.NewInt(1010_000_000)
	reservedMatchedAmt13 := math.NewInt(1010_000_000)
	reservedMatchedAmt14 := math.NewInt(0)
	reservedMatchedAmt15 := math.NewInt(1010_000_000)
	reservedMatchedAmt16 := math.NewInt(999_999_990)
	reservedMatchedAmt17 := math.NewInt(1010_000_000)

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

	refundAmt1 := math.NewInt(2000_000_000)
	refundAmt2 := math.NewInt(9)
	refundAmt3 := math.NewInt(200_000_000)
	refundAmt4 := math.NewInt(1640_000_000)
	refundAmt5 := math.NewInt(1640_000_000)
	refundAmt6 := math.NewInt(2120_000_000)
	refundAmt7 := math.NewInt(120_000_000)
	refundAmt8 := math.NewInt(1900_000_000)
	refundAmt9 := math.NewInt(0)
	refundAmt10 := math.NewInt(9)
	refundAmt11 := math.NewInt(65_000_000)
	refundAmt12 := math.NewInt(40_000_000)
	refundAmt13 := math.NewInt(10_000_000)
	refundAmt14 := math.NewInt(980_000_000)
	refundAmt15 := math.NewInt(15_000_000)
	refundAmt16 := math.NewInt(10)
	refundAmt17 := math.NewInt(42_000_000)

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
	err = s.keeper.AllocateSellingCoin(s.ctx, auction, mInfo)
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
