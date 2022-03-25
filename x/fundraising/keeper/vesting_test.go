package keeper_test

import (
	"time"

	"github.com/tendermint/fundraising/x/fundraising/types"

	_ "github.com/stretchr/testify/suite"
)

func (s *KeeperTestSuite) TestApplyVestingSchedules() {
	// TODO: Normal Vesting Schedules
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

	for _, vq := range s.keeper.GetVestingQueuesByAuctionId(s.ctx, auction.GetId()) {
		vestingReserveCoin = vestingReserveCoin.Sub(vq.PayingCoin)
	}
	s.Require().True(vestingReserveCoin.IsZero())
}
