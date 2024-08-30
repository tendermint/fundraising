package keeper_test

import (
	"time"

	"github.com/tendermint/fundraising/x/fundraising/types"

	_ "github.com/stretchr/testify/suite"
)

func (s *KeeperTestSuite) TestApplyVestingSchedules_NoSchedule() {
	startTime := time.Now().AddDate(0, 0, -1)
	endTime := startTime.AddDate(0, 1, 0)

	auction := s.createFixedPriceAuction(
		s.addr(0),
		parseDec("0.5"),
		parseCoin("1_000_000_000_000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		startTime,
		endTime,
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	s.placeBidFixedPrice(auction.GetId(), s.addr(1), parseDec("0.5"), parseCoin("15_000_000denom1"), true)
	s.placeBidFixedPrice(auction.GetId(), s.addr(2), parseDec("0.5"), parseCoin("15_000_000denom2"), true)
	s.placeBidFixedPrice(auction.GetId(), s.addr(3), parseDec("0.5"), parseCoin("55_000_000denom1"), true)
	s.placeBidFixedPrice(auction.GetId(), s.addr(4), parseDec("0.5"), parseCoin("30_000_000denom2"), true)

	// Apply schedules
	err := s.keeper.ApplyVestingSchedules(s.ctx, auction)
	s.Require().NoError(err)

	// Vesting reserve must have zero balance
	vestingReserveAddr := auction.GetVestingReserveAddress()
	vestingReserveCoin := s.getBalance(vestingReserveAddr, auction.PayingCoinDenom)
	s.Require().True(vestingReserveCoin.IsZero())

	// Auctioneer must have all the paying coin amounts in exchange of the selling coin
	auctioneerBalance := s.getBalance(auction.GetAuctioneer(), auction.PayingCoinDenom)
	s.Require().False(auctioneerBalance.IsZero())

	// Status must be finished
	a, err := s.keeper.Auction.Get(s.ctx, auction.GetId())
	s.Require().NoError(err)
	s.Require().Equal(types.AuctionStatusFinished, a.GetStatus())
}

func (s *KeeperTestSuite) TestApplyVestingSchedules_RemainingCoin() {
	startTime := time.Now().AddDate(0, 0, -1)
	endTime := startTime.AddDate(0, 1, 0)

	auction := s.createFixedPriceAuction(
		s.addr(0),
		parseDec("1.0"),
		parseCoin("1_000_000_000_000denom1"),
		"denom2",
		[]types.VestingSchedule{
			{
				ReleaseTime: endTime.AddDate(0, 6, 0),
				Weight:      parseDec("0.3"),
			},
			{
				ReleaseTime: endTime.AddDate(0, 9, 0),
				Weight:      parseDec("0.3"),
			},
			{
				ReleaseTime: endTime.AddDate(1, 0, 0),
				Weight:      parseDec("0.4"),
			},
		},
		startTime,
		endTime,
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	s.placeBidFixedPrice(auction.GetId(), s.addr(1), parseDec("1.0"), parseCoin("20000000denom2"), true)
	s.placeBidFixedPrice(auction.GetId(), s.addr(2), parseDec("1.0"), parseCoin("20000000denom2"), true)
	s.placeBidFixedPrice(auction.GetId(), s.addr(2), parseDec("1.0"), parseCoin("15000000denom2"), true)

	err := s.keeper.ApplyVestingSchedules(s.ctx, auction)
	s.Require().NoError(err)

	vestingReserveAddr := auction.GetVestingReserveAddress()
	vestingReserveCoin := s.getBalance(vestingReserveAddr, auction.PayingCoinDenom)

	vqs, err := s.keeper.GetVestingQueuesByAuctionId(s.ctx, auction.GetId())
	s.Require().NoError(err)

	for _, vq := range vqs {
		vestingReserveCoin = vestingReserveCoin.Sub(vq.PayingCoin)
	}
	s.Require().True(vestingReserveCoin.IsZero())
}
