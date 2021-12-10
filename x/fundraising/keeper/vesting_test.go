package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/fundraising/x/fundraising/types"

	_ "github.com/stretchr/testify/suite"
)

func (suite *KeeperTestSuite) TestVestingQueues() {
	auctionId := uint64(1)
	auctioneerAcc := suite.addrs[0]
	payingReserveAcc := suite.addrs[5] // 100_000_000_000_000denom1
	payingCoinDenom := denom1
	reserveBalance := suite.app.BankKeeper.GetBalance(suite.ctx, payingReserveAcc, payingCoinDenom)

	for _, vs := range suite.VestingSchedules() {
		payingAmt := reserveBalance.Amount.ToDec().Mul(vs.Weight).TruncateInt()

		queue := types.NewVestingQueue(
			auctionId,
			auctioneerAcc.String(),
			sdk.NewCoin(payingCoinDenom, payingAmt),
			vs.ReleaseTime,
			false,
		)
		suite.keeper.SetVestingQueue(suite.ctx, auctionId, vs.ReleaseTime, queue)
	}

	vestingQueues := suite.keeper.GetVestingQueuesByAuctionId(suite.ctx, auctionId)
	suite.Require().Len(vestingQueues, 4)

	totalPayingCoin := sdk.NewInt64Coin(denom1, 0)
	for _, vq := range vestingQueues {
		totalPayingCoin = totalPayingCoin.Add(vq.PayingCoin)
	}
	suite.Require().Equal(sdk.NewInt64Coin(denom1, 100_000_000_000_000), totalPayingCoin)
}
