package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/fundraising/x/fundraising/types"

	_ "github.com/stretchr/testify/suite"
)

func (suite *KeeperTestSuite) TestVestingQueues() {
	payingReserveAcc := suite.addrs[5] // 100_000_000_000_000denom1
	payingCoinDenom := denom1
	reserveBalance := suite.app.BankKeeper.GetBalance(suite.ctx, payingReserveAcc, denom1)

	// vesting schedule contains 4 vesting queues
	for _, vs := range suite.VestingSchedules() {
		payingAmt := reserveBalance.Amount.ToDec().Mul(vs.Weight).TruncateInt()

		suite.keeper.SetVestingQueue(suite.ctx, uint64(1), vs.ReleaseTime, types.VestingQueue{
			AuctionId:   uint64(1),
			Auctioneer:  suite.addrs[0].String(),
			PayingCoin:  sdk.NewCoin(payingCoinDenom, payingAmt),
			ReleaseTime: vs.ReleaseTime,
			Vested:      false,
		})
	}

	// vesting schedule contains 2 vesting queues
	for _, vs := range suite.VestingSchedules2() {
		payingAmt := reserveBalance.Amount.ToDec().Mul(vs.Weight).TruncateInt()

		suite.keeper.SetVestingQueue(suite.ctx, uint64(2), vs.ReleaseTime, types.VestingQueue{
			AuctionId:   uint64(2),
			Auctioneer:  suite.addrs[1].String(),
			PayingCoin:  sdk.NewCoin(payingCoinDenom, payingAmt),
			ReleaseTime: vs.ReleaseTime,
			Vested:      false,
		})
	}

	suite.Require().Len(suite.keeper.GetVestingQueuesByAuctionId(suite.ctx, uint64(1)), 4)
	suite.Require().Len(suite.keeper.GetVestingQueuesByAuctionId(suite.ctx, uint64(2)), 2)
	suite.Require().Len(suite.keeper.GetVestingQueues(suite.ctx), 6)

	totalPayingCoin := sdk.NewInt64Coin(denom1, 0)
	for _, vq := range suite.keeper.GetVestingQueuesByAuctionId(suite.ctx, uint64(1)) {
		totalPayingCoin = totalPayingCoin.Add(vq.PayingCoin)
	}
	suite.Require().Equal(reserveBalance, totalPayingCoin)
}
