package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/fundraising/x/fundraising/types"

	_ "github.com/stretchr/testify/suite"
)

func (suite *KeeperTestSuite) TestVestingQueueRemainingCoin() {
	auction := &types.BaseAuction{
		Id:                    1,
		Type:                  types.AuctionTypeFixedPrice,
		Auctioneer:            suite.addrs[5].String(),
		SellingReserveAddress: types.SellingReserveAcc(1).String(),
		PayingReserveAddress:  types.PayingReserveAcc(1).String(),
		StartPrice:            sdk.MustNewDecFromStr("0.5"),
		SellingCoin:           sdk.NewInt64Coin(denom3, 1_000_000_000_000),
		PayingCoinDenom:       denom4,
		VestingReserveAddress: types.VestingReserveAcc(1).String(),
		VestingSchedules: []types.VestingSchedule{
			{
				ReleaseTime: types.ParseTime("2022-01-01T22:00:00+00:00"),
				Weight:      sdk.MustNewDecFromStr("0.3"),
			},
			{
				ReleaseTime: types.ParseTime("2022-04-01T22:00:00+00:00"),
				Weight:      sdk.MustNewDecFromStr("0.3"),
			},
			{
				ReleaseTime: types.ParseTime("2022-08-01T22:00:00+00:00"),
				Weight:      sdk.MustNewDecFromStr("0.4"),
			},
		},
		WinningPrice:  sdk.ZeroDec(),
		RemainingCoin: sdk.NewInt64Coin(denom3, 1_000_000_000_000),
		StartTime:     types.ParseTime("2021-12-10T00:00:00Z"),
		EndTimes:      []time.Time{types.ParseTime("2022-12-20T00:00:00Z")},
		Status:        types.AuctionStatusStarted,
	}
	suite.SetAuction(suite.ctx, auction)

	bids := []types.Bid{
		{
			AuctionId: auction.Id,
			Sequence:  1,
			Bidder:    suite.addrs[1].String(),
			Price:     auction.StartPrice,
			Coin:      sdk.NewInt64Coin(denom4, 666666),
		},
	}

	for _, bid := range bids {
		suite.PlaceBidWithCustom(suite.ctx, bid.AuctionId, bid.Sequence, bid.Bidder, bid.Price, bid.Coin)
	}

	err := suite.keeper.SetVestingSchedules(suite.ctx, auction)
	suite.Require().NoError(err)

	reserveCoin := suite.app.BankKeeper.GetBalance(suite.ctx, auction.GetVestingReserveAddress(), denom4)

	for _, vq := range suite.keeper.GetVestingQueuesByAuctionId(suite.ctx, auction.GetId()) {
		reserveCoin = reserveCoin.Sub(vq.PayingCoin)
	}
	suite.Require().True(reserveCoin.IsZero())
}

func (suite *KeeperTestSuite) TestVestingQueueIterator() {
	payingReserveAcc, payingCoinDenom := suite.addrs[5], denom1 // 100_000_000_000_000denom1
	reserveCoin := suite.app.BankKeeper.GetBalance(suite.ctx, payingReserveAcc, denom1)

	// vesting schedule contains 2 vesting queues
	for _, vs := range suite.sampleVestingSchedules1 {
		payingAmt := reserveCoin.Amount.ToDec().Mul(vs.Weight).TruncateInt()

		suite.keeper.SetVestingQueue(suite.ctx, uint64(1), vs.ReleaseTime, types.VestingQueue{
			AuctionId:   uint64(1),
			Auctioneer:  suite.addrs[0].String(),
			PayingCoin:  sdk.NewCoin(payingCoinDenom, payingAmt),
			ReleaseTime: vs.ReleaseTime,
			Released:    false,
		})
	}

	// vesting schedule contains 4 vesting queues
	for _, vs := range suite.sampleVestingSchedules2 {
		payingAmt := reserveCoin.Amount.ToDec().Mul(vs.Weight).TruncateInt()

		suite.keeper.SetVestingQueue(suite.ctx, uint64(2), vs.ReleaseTime, types.VestingQueue{
			AuctionId:   uint64(2),
			Auctioneer:  suite.addrs[1].String(),
			PayingCoin:  sdk.NewCoin(payingCoinDenom, payingAmt),
			ReleaseTime: vs.ReleaseTime,
			Released:    false,
		})
	}

	suite.Require().Len(suite.keeper.GetVestingQueuesByAuctionId(suite.ctx, uint64(1)), 2)
	suite.Require().Len(suite.keeper.GetVestingQueuesByAuctionId(suite.ctx, uint64(2)), 4)
	suite.Require().Len(suite.keeper.GetVestingQueues(suite.ctx), 6)

	totalPayingCoin := sdk.NewInt64Coin(denom1, 0)
	for _, vq := range suite.keeper.GetVestingQueuesByAuctionId(suite.ctx, uint64(2)) {
		totalPayingCoin = totalPayingCoin.Add(vq.PayingCoin)
	}
	suite.Require().Equal(reserveCoin, totalPayingCoin)
}
