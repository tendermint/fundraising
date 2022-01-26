package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/fundraising/x/fundraising/types"

	_ "github.com/stretchr/testify/suite"
)

func (s *KeeperTestSuite) TestVestingQueueRemainingCoin() {
	auction := s.createFixedPriceAuction(
		s.addr(0),
		sdk.OneDec(),
		sdk.NewInt64Coin("denom1", 1_000_000_000_000),
		"denom2",
		[]types.VestingSchedule{
			{
				ReleaseTime: types.MustParseRFC3339("2022-01-01T22:00:00+00:00"),
				Weight:      sdk.MustNewDecFromStr("0.3"),
			},
			{
				ReleaseTime: types.MustParseRFC3339("2022-04-01T22:00:00+00:00"),
				Weight:      sdk.MustNewDecFromStr("0.3"),
			},
			{
				ReleaseTime: types.MustParseRFC3339("2022-08-01T22:00:00+00:00"),
				Weight:      sdk.MustNewDecFromStr("0.4"),
			},
		},
		types.MustParseRFC3339("2022-01-01T00:00:00Z"),
		types.MustParseRFC3339("2022-03-01T00:00:00Z"),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	s.placeBid(auction.GetId(), s.addr(1), sdk.OneDec(), sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 20_000_000), true)
	s.placeBid(auction.GetId(), s.addr(2), sdk.OneDec(), sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 20_000_000), true)
	s.placeBid(auction.GetId(), s.addr(2), sdk.OneDec(), sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 15_000_000), true)

	err := s.keeper.SetVestingSchedules(s.ctx, auction)
	s.Require().NoError(err)

	vestingReserveAddr := auction.GetVestingReserveAddress()
	vestingReserveCoin := s.getBalance(vestingReserveAddr, auction.PayingCoinDenom)

	for _, vq := range s.keeper.GetVestingQueuesByAuctionId(s.ctx, auction.GetId()) {
		vestingReserveCoin = vestingReserveCoin.Sub(vq.PayingCoin)
	}
	s.Require().True(vestingReserveCoin.IsZero())
}

func (s *KeeperTestSuite) TestVestingQueueIterator() {
	payingReserveAcc := s.addr(0)
	payingCoinDenom := "denom1"
	reserveCoin := s.getBalance(payingReserveAcc, payingCoinDenom)

	// Set vesting schedules with 2 vesting queues
	for _, vs := range []types.VestingSchedule{
		{
			ReleaseTime: types.MustParseRFC3339("2023-01-01T00:00:00Z"),
			Weight:      sdk.MustNewDecFromStr("0.5"),
		},
		{
			ReleaseTime: types.MustParseRFC3339("2023-06-01T00:00:00Z"),
			Weight:      sdk.MustNewDecFromStr("0.5"),
		},
	} {
		payingAmt := reserveCoin.Amount.ToDec().MulTruncate(vs.Weight).TruncateInt()

		s.keeper.SetVestingQueue(s.ctx, uint64(1), vs.ReleaseTime, types.VestingQueue{
			AuctionId:   uint64(1),
			Auctioneer:  s.addr(1).String(),
			PayingCoin:  sdk.NewCoin(payingCoinDenom, payingAmt),
			ReleaseTime: vs.ReleaseTime,
			Released:    false,
		})
	}

	// Set vesting schedules with 4 vesting queues
	for _, vs := range []types.VestingSchedule{
		{
			ReleaseTime: types.MustParseRFC3339("2023-01-01T00:00:00Z"),
			Weight:      sdk.MustNewDecFromStr("0.25"),
		},
		{
			ReleaseTime: types.MustParseRFC3339("2023-05-01T00:00:00Z"),
			Weight:      sdk.MustNewDecFromStr("0.25"),
		},
		{
			ReleaseTime: types.MustParseRFC3339("2023-09-01T00:00:00Z"),
			Weight:      sdk.MustNewDecFromStr("0.25"),
		},
		{
			ReleaseTime: types.MustParseRFC3339("2023-12-01T00:00:00Z"),
			Weight:      sdk.MustNewDecFromStr("0.25"),
		},
	} {
		payingAmt := reserveCoin.Amount.ToDec().MulTruncate(vs.Weight).TruncateInt()

		s.keeper.SetVestingQueue(s.ctx, uint64(2), vs.ReleaseTime, types.VestingQueue{
			AuctionId:   uint64(2),
			Auctioneer:  s.addr(2).String(),
			PayingCoin:  sdk.NewCoin(payingCoinDenom, payingAmt),
			ReleaseTime: vs.ReleaseTime,
			Released:    false,
		})
	}

	s.Require().Len(s.keeper.GetVestingQueuesByAuctionId(s.ctx, uint64(1)), 2)
	s.Require().Len(s.keeper.GetVestingQueuesByAuctionId(s.ctx, uint64(2)), 4)
	s.Require().Len(s.keeper.GetVestingQueues(s.ctx), 6)

	totalPayingCoin := sdk.NewInt64Coin(payingCoinDenom, 0)
	for _, vq := range s.keeper.GetVestingQueuesByAuctionId(s.ctx, uint64(2)) {
		totalPayingCoin = totalPayingCoin.Add(vq.PayingCoin)
	}
	s.Require().Equal(reserveCoin, totalPayingCoin)
}
