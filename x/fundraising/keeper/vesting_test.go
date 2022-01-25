package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/fundraising/x/fundraising/types"

	_ "github.com/stretchr/testify/suite"
)

func (s *KeeperTestSuite) TestVestingQueueRemainingCoin() {
	auction := &types.BaseAuction{
		Id:                    1,
		Type:                  types.AuctionTypeFixedPrice,
		Auctioneer:            s.addrs[5].String(),
		SellingReserveAddress: types.SellingReserveAcc(1).String(),
		PayingReserveAddress:  types.PayingReserveAcc(1).String(),
		StartPrice:            sdk.MustNewDecFromStr("0.5"),
		SellingCoin:           sdk.NewInt64Coin(denom3, 1_000_000_000_000),
		PayingCoinDenom:       denom4,
		VestingReserveAddress: types.VestingReserveAcc(1).String(),
		VestingSchedules: []types.VestingSchedule{
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
		WinningPrice:  sdk.ZeroDec(),
		RemainingCoin: sdk.NewInt64Coin(denom3, 1_000_000_000_000),
		StartTime:     types.MustParseRFC3339("2021-12-10T00:00:00Z"),
		EndTimes:      []time.Time{types.MustParseRFC3339("2022-12-20T00:00:00Z")},
		Status:        types.AuctionStatusStarted,
	}

	s.SetAuction(auction)

	bids := []types.Bid{
		{
			AuctionId: auction.Id,
			Sequence:  1,
			Bidder:    s.addrs[1].String(),
			Price:     auction.StartPrice,
			Coin:      sdk.NewInt64Coin(denom4, 666666),
		},
	}

	for _, bid := range bids {
		s.PlaceBidWithCustom(bid.AuctionId, bid.Sequence, bid.Bidder, bid.Price, bid.Coin)
	}

	err := s.keeper.SetVestingSchedules(s.ctx, auction)
	s.Require().NoError(err)

	reserveCoin := s.app.BankKeeper.GetBalance(s.ctx, auction.GetVestingReserveAddress(), denom4)

	for _, vq := range s.keeper.GetVestingQueuesByAuctionId(s.ctx, auction.GetId()) {
		reserveCoin = reserveCoin.Sub(vq.PayingCoin)
	}
	s.Require().True(reserveCoin.IsZero())
}

func (s *KeeperTestSuite) TestVestingQueueIterator() {
	payingReserveAcc, payingCoinDenom := s.addrs[5], denom1 // 100_000_000_000_000denom1
	reserveCoin := s.app.BankKeeper.GetBalance(s.ctx, payingReserveAcc, denom1)

	// vesting schedule contains 2 vesting queues
	for _, vs := range s.sampleVestingSchedules1 {
		payingAmt := reserveCoin.Amount.ToDec().MulTruncate(vs.Weight).TruncateInt()

		s.keeper.SetVestingQueue(s.ctx, uint64(1), vs.ReleaseTime, types.VestingQueue{
			AuctionId:   uint64(1),
			Auctioneer:  s.addrs[0].String(),
			PayingCoin:  sdk.NewCoin(payingCoinDenom, payingAmt),
			ReleaseTime: vs.ReleaseTime,
			Released:    false,
		})
	}

	// vesting schedule contains 4 vesting queues
	for _, vs := range s.sampleVestingSchedules2 {
		payingAmt := reserveCoin.Amount.ToDec().MulTruncate(vs.Weight).TruncateInt()

		s.keeper.SetVestingQueue(s.ctx, uint64(2), vs.ReleaseTime, types.VestingQueue{
			AuctionId:   uint64(2),
			Auctioneer:  s.addrs[1].String(),
			PayingCoin:  sdk.NewCoin(payingCoinDenom, payingAmt),
			ReleaseTime: vs.ReleaseTime,
			Released:    false,
		})
	}

	s.Require().Len(s.keeper.GetVestingQueuesByAuctionId(s.ctx, uint64(1)), 2)
	s.Require().Len(s.keeper.GetVestingQueuesByAuctionId(s.ctx, uint64(2)), 4)
	s.Require().Len(s.keeper.GetVestingQueues(s.ctx), 6)

	totalPayingCoin := sdk.NewInt64Coin(denom1, 0)
	for _, vq := range s.keeper.GetVestingQueuesByAuctionId(s.ctx, uint64(2)) {
		totalPayingCoin = totalPayingCoin.Add(vq.PayingCoin)
	}
	s.Require().Equal(reserveCoin, totalPayingCoin)
}
