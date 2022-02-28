package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

func (s *KeeperTestSuite) TestExecuteStartedAuction_BatchAuction() {
	ba := s.createBatchAuction(
		s.addr(1),
		parseDec("1"),
		parseCoin("10000000000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		1,
		sdk.MustNewDecFromStr("0.2"),
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, ba.GetStatus())

	s.placeBidBatchWorth(ba.Id, s.addr(1), parseDec("10"), parseCoin("100000000denom2"), sdk.NewInt(1000000000), true)
	s.placeBidBatchWorth(ba.Id, s.addr(2), parseDec("9"), parseCoin("150000000denom2"), sdk.NewInt(1000000000), true)
	s.placeBidBatchWorth(ba.Id, s.addr(3), parseDec("5.5"), parseCoin("250000000denom2"), sdk.NewInt(1000000000), true)
	s.placeBidBatchMany(ba.Id, s.addr(4), parseDec("6"), parseCoin("400000000denom1"), sdk.NewInt(1000000000), true)
	s.placeBidBatchMany(ba.Id, s.addr(6), parseDec("4.5"), parseCoin("150000000denom1"), sdk.NewInt(1000000000), true)
	s.placeBidBatchMany(ba.Id, s.addr(7), parseDec("3.8"), parseCoin("150000000denom1"), sdk.NewInt(1000000000), true)

	auction, found := s.keeper.GetAuction(s.ctx, ba.Id)
	s.Require().True(found)

	s.keeper.ExecuteStartedStatus(s.ctx, auction)

}
