package keeper_test

import (
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/stretchr/testify/suite"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

func (s *KeeperTestSuite) TestEndBlockerStandByStatus() {
	standByAuction := s.createFixedPriceAuction(
		s.addr(0),
		parseDec("0.5"),
		parseCoin("500000000000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		time.Now().AddDate(0, 3, 0),
		time.Now().AddDate(0, 5, 0),
		true,
	)

	auction, err := s.keeper.Auction.Get(s.ctx, standByAuction.GetId())
	s.Require().NoError(err)
	s.Require().Equal(types.AuctionStatusStandBy, auction.GetStatus())

	// Modify current time and call end blocker
	s.ctx = s.ctx.WithBlockTime(standByAuction.StartTime.AddDate(0, 0, 1))
	s.Require().NoError(s.keeper.BeginBlocker(s.ctx))

	auction, err = s.keeper.Auction.Get(s.ctx, standByAuction.GetId())
	s.Require().NoError(err)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())
}

func (s *KeeperTestSuite) TestEndBlockerStartedStatus() {
	auctioneer := s.addr(0)
	auction := s.createFixedPriceAuction(
		auctioneer,
		parseDec("1"),
		parseCoin("500000000000denom1"),
		"denom2",
		[]types.VestingSchedule{
			{
				ReleaseTime: types.MustParseRFC3339("2024-01-01T00:00:00Z"),
				Weight:      math.LegacyMustNewDecFromStr("0.5"),
			},
			{
				ReleaseTime: types.MustParseRFC3339("2024-06-01T00:00:00Z"),
				Weight:      math.LegacyMustNewDecFromStr("0.5"),
			},
		},
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 1, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	bid1 := s.placeBidFixedPrice(auction.GetId(), s.addr(1), math.LegacyOneDec(), parseCoin("20000000denom2"), true)
	bid2 := s.placeBidFixedPrice(auction.GetId(), s.addr(2), math.LegacyOneDec(), parseCoin("20000000denom2"), true)
	bid3 := s.placeBidFixedPrice(auction.GetId(), s.addr(3), math.LegacyOneDec(), parseCoin("20000000denom2"), true)

	totalBidCoin := bid1.Coin.Add(bid2.Coin).Add(bid3.Coin)
	receiveAmt := math.LegacyNewDecFromInt(totalBidCoin.Amount).QuoTruncate(auction.GetStartPrice()).TruncateInt()
	receiveCoin := sdk.NewCoin(auction.GetSellingCoin().Denom, receiveAmt)

	payingReserve := s.getBalance(auction.GetPayingReserveAddress(), auction.GetPayingCoinDenom())
	s.Require().True(coinEq(totalBidCoin, payingReserve))

	// Modify the current block time a day after the end time
	s.ctx = s.ctx.WithBlockTime(auction.GetEndTimes()[0].AddDate(0, 0, 1))
	s.Require().NoError(s.keeper.BeginBlocker(s.ctx))

	// The remaining selling coin must be returned to the auctioneer
	auctioneerBalance := s.getBalance(auctioneer, auction.GetSellingCoin().Denom)
	s.Require().True(coinEq(auction.GetSellingCoin(), auctioneerBalance.Add(receiveCoin)))
}

func (s *KeeperTestSuite) TestEndBlockerVestingStatus() {
	auctioneer := s.addr(0)
	auction := s.createFixedPriceAuction(
		auctioneer,
		parseDec("1"),
		sdk.NewInt64Coin("denom1", 500_000_000_000),
		"denom2",
		[]types.VestingSchedule{
			{
				ReleaseTime: time.Now().AddDate(0, 0, -1).AddDate(0, 6, 0),
				Weight:      math.LegacyMustNewDecFromStr("0.5"),
			},
			{
				ReleaseTime: time.Now().AddDate(0, 0, -1).AddDate(1, 0, 0),
				Weight:      math.LegacyMustNewDecFromStr("0.5"),
			},
		},
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 1, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	bid1 := s.placeBidFixedPrice(auction.GetId(), s.addr(1), math.LegacyOneDec(), parseCoin("20000000denom2"), true)
	bid2 := s.placeBidFixedPrice(auction.GetId(), s.addr(2), math.LegacyOneDec(), parseCoin("20000000denom2"), true)
	bid3 := s.placeBidFixedPrice(auction.GetId(), s.addr(3), math.LegacyOneDec(), parseCoin("20000000denom2"), true)

	totalBidCoin := bid1.Coin.Add(bid2.Coin).Add(bid3.Coin)

	// Modify the current block time a day after the end time
	s.ctx = s.ctx.WithBlockTime(auction.GetEndTimes()[0].AddDate(0, 0, 1))
	s.Require().NoError(s.keeper.BeginBlocker(s.ctx))

	vestingReserve := s.getBalance(auction.GetVestingReserveAddress(), auction.GetPayingCoinDenom())
	s.Require().Equal(totalBidCoin, vestingReserve)

	// Modify the current block time a day after the last vesting schedule
	s.ctx = s.ctx.WithBlockTime(auction.VestingSchedules[len(auction.VestingSchedules)-1].ReleaseTime.AddDate(0, 0, 1))
	s.Require().NoError(s.keeper.BeginBlocker(s.ctx))

	queues, err := s.keeper.GetVestingQueuesByAuctionId(s.ctx, auction.GetId())
	s.Require().NoError(err)
	s.Require().Len(queues, 2)
	s.Require().True(queues[0].Released)
	s.Require().True(queues[1].Released)

	// The auctioneer must have released the paying coin
	auctioneerBalance := s.getBalance(auctioneer, auction.GetPayingCoinDenom())
	s.Require().True(coinEq(totalBidCoin, auctioneerBalance))
}

func (s *KeeperTestSuite) TestExecuteStartedAuction_BatchAuction() {
	ba := s.createBatchAuction(
		s.addr(1),
		parseDec("1"),
		parseDec("0.1"),
		parseCoin("10000000000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		1,
		math.LegacyMustNewDecFromStr("0.2"),
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, ba.GetStatus())

	s.placeBidBatchWorth(ba.Id, s.addr(1), parseDec("10"), parseCoin("100000000denom2"), math.NewInt(1000000000), true)
	s.placeBidBatchWorth(ba.Id, s.addr(2), parseDec("9"), parseCoin("150000000denom2"), math.NewInt(1000000000), true)
	s.placeBidBatchWorth(ba.Id, s.addr(3), parseDec("5.5"), parseCoin("250000000denom2"), math.NewInt(1000000000), true)
	s.placeBidBatchMany(ba.Id, s.addr(4), parseDec("6"), parseCoin("400000000denom1"), math.NewInt(1000000000), true)
	s.placeBidBatchMany(ba.Id, s.addr(6), parseDec("4.5"), parseCoin("150000000denom1"), math.NewInt(1000000000), true)
	s.placeBidBatchMany(ba.Id, s.addr(7), parseDec("3.8"), parseCoin("150000000denom1"), math.NewInt(1000000000), true)

	auction, err := s.keeper.Auction.Get(s.ctx, ba.Id)
	s.Require().NoError(err)

	err = s.keeper.ExecuteStartedStatus(s.ctx, auction)
	s.Require().NoError(err)
}
