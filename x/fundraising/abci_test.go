package fundraising_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/fundraising/x/fundraising"
	"github.com/tendermint/fundraising/x/fundraising/types"

	_ "github.com/stretchr/testify/suite"
)

func (s *ModuleTestSuite) TestEndBlockerStandByStatus() {
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
	s.Require().Equal(types.AuctionStatusStandBy, standByAuction.GetStatus())

	// Modify current time
	s.ctx = s.ctx.WithBlockTime(standByAuction.StartTime.AddDate(0, 0, 1))
	fundraising.EndBlocker(s.ctx, s.keeper)

	auction, found := s.keeper.GetAuction(s.ctx, standByAuction.GetId())
	s.Require().True(found)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())
}

func (s *ModuleTestSuite) TestEndBlockerStartedStatus() {
	auctioneer := s.addr(0)
	auction := s.createFixedPriceAuction(
		auctioneer,
		parseDec("1"),
		parseCoin("500000000000denom1"),
		"denom2",
		[]types.VestingSchedule{
			{
				ReleaseTime: types.MustParseRFC3339("2024-01-01T00:00:00Z"),
				Weight:      sdk.MustNewDecFromStr("0.5"),
			},
			{
				ReleaseTime: types.MustParseRFC3339("2024-06-01T00:00:00Z"),
				Weight:      sdk.MustNewDecFromStr("0.5"),
			},
		},
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 1, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	bid1 := s.placeBid(auction.GetId(), s.addr(1), types.BidTypeFixedPrice, sdk.OneDec(), parseCoin("20000000denom2"), true)
	bid2 := s.placeBid(auction.GetId(), s.addr(2), types.BidTypeFixedPrice, sdk.OneDec(), parseCoin("20000000denom2"), true)
	bid3 := s.placeBid(auction.GetId(), s.addr(3), types.BidTypeFixedPrice, sdk.OneDec(), parseCoin("20000000denom2"), true)

	totalBidCoin := bid1.Coin.Add(bid2.Coin).Add(bid3.Coin)
	receiveAmt := totalBidCoin.Amount.ToDec().QuoTruncate(auction.GetStartPrice()).TruncateInt()
	receiveCoin := sdk.NewCoin(auction.GetSellingCoin().Denom, receiveAmt)

	payingReserve := s.getBalance(auction.GetPayingReserveAddress(), auction.GetPayingCoinDenom())
	s.Require().True(coinEq(totalBidCoin, payingReserve))

	// Modify the current block time a day after the end time
	s.ctx = s.ctx.WithBlockTime(auction.GetEndTimes()[0].AddDate(0, 0, 1))
	fundraising.EndBlocker(s.ctx, s.keeper)

	// The remaining selling coin must be returned to the auctioneer
	auctioneerBalance := s.getBalance(auctioneer, auction.GetSellingCoin().Denom)
	s.Require().True(coinEq(auction.GetSellingCoin(), auctioneerBalance.Add(receiveCoin)))
}

func (s *ModuleTestSuite) TestEndBlockerVestingStatus() {
	auctioneer := s.addr(0)
	auction := s.createFixedPriceAuction(
		auctioneer,
		parseDec("1"),
		sdk.NewInt64Coin("denom1", 500_000_000_000),
		"denom2",
		[]types.VestingSchedule{
			{
				ReleaseTime: time.Now().AddDate(0, 0, -1).AddDate(0, 6, 0),
				Weight:      sdk.MustNewDecFromStr("0.5"),
			},
			{
				ReleaseTime: time.Now().AddDate(0, 0, -1).AddDate(1, 0, 0),
				Weight:      sdk.MustNewDecFromStr("0.5"),
			},
		},
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 1, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	bid1 := s.placeBid(auction.GetId(), s.addr(1), types.BidTypeFixedPrice, sdk.OneDec(), parseCoin("20000000denom2"), true)
	bid2 := s.placeBid(auction.GetId(), s.addr(2), types.BidTypeFixedPrice, sdk.OneDec(), parseCoin("20000000denom2"), true)
	bid3 := s.placeBid(auction.GetId(), s.addr(3), types.BidTypeFixedPrice, sdk.OneDec(), parseCoin("20000000denom2"), true)

	totalBidCoin := bid1.Coin.Add(bid2.Coin).Add(bid3.Coin)

	// Modify the current block time a day after the end time
	s.ctx = s.ctx.WithBlockTime(auction.GetEndTimes()[0].AddDate(0, 0, 1))
	fundraising.EndBlocker(s.ctx, s.keeper)

	vestingReserve := s.getBalance(auction.GetVestingReserveAddress(), auction.GetPayingCoinDenom())
	s.Require().Equal(totalBidCoin, vestingReserve)

	// Modify the current block time a day after the last vesting schedule
	s.ctx = s.ctx.WithBlockTime(auction.VestingSchedules[len(auction.VestingSchedules)-1].ReleaseTime.AddDate(0, 0, 1))
	fundraising.EndBlocker(s.ctx, s.keeper)

	queues := s.keeper.GetVestingQueuesByAuctionId(s.ctx, auction.GetId())
	s.Require().Len(queues, 2)
	s.Require().True(queues[0].Released)
	s.Require().True(queues[1].Released)

	// The auctioneer must have released the paying coin
	auctioneerBalance := s.getBalance(auctioneer, auction.GetPayingCoinDenom())
	s.Require().True(coinEq(totalBidCoin, auctioneerBalance))
}
