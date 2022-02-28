package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/fundraising/x/fundraising"
	"github.com/tendermint/fundraising/x/fundraising/types"

	_ "github.com/stretchr/testify/suite"
)

func (s *KeeperTestSuite) TestDefaultGenesis() {
	genState := types.DefaultGenesisState()

	s.keeper.InitGenesis(s.ctx, *genState)
	got := s.keeper.ExportGenesis(s.ctx)
	s.Require().Equal(genState, got)
}

func (s *KeeperTestSuite) TestGenesisState() {
	fixedAuction := s.createFixedPriceAuction(
		s.addr(0),
		parseDec("1.0"),
		parseDec("0.1"),
		parseCoin("200000000000denom1"),
		"denom2",
		[]types.VestingSchedule{
			{
				ReleaseTime: time.Now().AddDate(0, 0, -1).AddDate(0, 6, 0),
				Weight:      sdk.MustNewDecFromStr("0.25"),
			},
			{
				ReleaseTime: time.Now().AddDate(0, 0, -1).AddDate(0, 9, 0),
				Weight:      sdk.MustNewDecFromStr("0.25"),
			},
			{
				ReleaseTime: time.Now().AddDate(0, 0, -1).AddDate(1, 0, 0),
				Weight:      sdk.MustNewDecFromStr("0.25"),
			},
			{
				ReleaseTime: time.Now().AddDate(0, 0, -1).AddDate(1, 3, 0),
				Weight:      sdk.MustNewDecFromStr("0.25"),
			},
		},
		time.Now().AddDate(0, -6, 0),
		time.Now().AddDate(0, 0, 10),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, fixedAuction.GetStatus())

	// Place bids
	s.placeBidFixedPrice(fixedAuction.Id, s.addr(1), sdk.OneDec(), parseCoin("20000000denom2"), true)
	s.placeBidFixedPrice(fixedAuction.Id, s.addr(2), sdk.OneDec(), parseCoin("30000000denom2"), true)
	s.placeBidFixedPrice(fixedAuction.Id, s.addr(3), sdk.OneDec(), parseCoin("15000000denom2"), true)
	s.placeBidFixedPrice(fixedAuction.Id, s.addr(4), sdk.OneDec(), parseCoin("35000000denom2"), true)

	batchAuction := s.createBatchAuction(
		s.addr(1),
		parseDec("0.1"),
		parseDec("0.1"),
		parseCoin("1000000000000denom3"),
		"denom4",
		[]types.VestingSchedule{
			{
				ReleaseTime: time.Now().AddDate(0, 1, 0).AddDate(0, 6, 0),
				Weight:      sdk.MustNewDecFromStr("0.5"),
			},
			{
				ReleaseTime: time.Now().AddDate(0, 1, 0).AddDate(1, 0, 0),
				Weight:      sdk.MustNewDecFromStr("0.5"),
			},
		},
		3,
		parseDec("0.3"),
		time.Now().AddDate(0, -1, 0),
		time.Now().AddDate(0, 1, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, batchAuction.GetStatus())

	s.placeBidBatchWorth(batchAuction.Id, s.addr(2), parseDec("0.5"), parseCoin("100000000denom4"), sdk.NewInt(1000000000), true)
	s.placeBidBatchWorth(batchAuction.Id, s.addr(2), parseDec("0.4"), parseCoin("150000000denom4"), sdk.NewInt(1000000000), true)
	s.placeBidBatchWorth(batchAuction.Id, s.addr(3), parseDec("0.66"), parseCoin("250000000denom4"), sdk.NewInt(1000000000), true)
	s.placeBidBatchMany(batchAuction.Id, s.addr(4), parseDec("0.35"), parseCoin("400000000denom3"), sdk.NewInt(1000000000), true)
	s.placeBidBatchMany(batchAuction.Id, s.addr(5), parseDec("0.8"), parseCoin("150000000denom3"), sdk.NewInt(1000000000), true)
	s.placeBidBatchMany(batchAuction.Id, s.addr(6), parseDec("0.2"), parseCoin("150000000denom3"), sdk.NewInt(1000000000), true)

	// Modify the current block time a day after the end time
	s.ctx = s.ctx.WithBlockTime(fixedAuction.GetEndTimes()[0].AddDate(0, 0, 1))
	fundraising.EndBlocker(s.ctx, s.keeper)

	// Modify the time to make the first and second vesting queues over
	s.ctx = s.ctx.WithBlockTime(fixedAuction.VestingSchedules[1].ReleaseTime.AddDate(0, 0, 1))
	fundraising.EndBlocker(s.ctx, s.keeper)

	queues := s.keeper.GetVestingQueuesByAuctionId(s.ctx, 1)
	s.Require().Len(queues, 4)

	for i, queue := range queues {
		if i == 0 || i == 1 {
			s.Require().True(queue.Released)
		} else {
			s.Require().False(queue.Released)
		}
	}

	var genState *types.GenesisState
	s.Require().NotPanics(func() {
		genState = s.keeper.ExportGenesis(s.ctx)
	})
	s.Require().NoError(genState.Validate())

	s.Require().NotPanics(func() {
		s.keeper.InitGenesis(s.ctx, *genState)
	})
	s.Require().Equal(genState, s.keeper.ExportGenesis(s.ctx))
}
