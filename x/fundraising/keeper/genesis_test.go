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
	auction := s.createFixedPriceAuction(
		s.addr(0),
		sdk.OneDec(),
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
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 1, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	// Place bids
	s.placeBidFixedPrice(auction.GetId(), s.addr(1), sdk.OneDec(), parseCoin("20000000denom2"), true)
	s.placeBidFixedPrice(auction.GetId(), s.addr(2), sdk.OneDec(), parseCoin("30000000denom2"), true)
	s.placeBidFixedPrice(auction.GetId(), s.addr(3), sdk.OneDec(), parseCoin("15000000denom2"), true)
	s.placeBidFixedPrice(auction.GetId(), s.addr(4), sdk.OneDec(), parseCoin("35000000denom2"), true)

	// Modify the current block time a day after the end time
	s.ctx = s.ctx.WithBlockTime(auction.GetEndTimes()[0].AddDate(0, 0, 1))
	fundraising.EndBlocker(s.ctx, s.keeper)

	// Modify the time to make the first and second vesting queues over
	s.ctx = s.ctx.WithBlockTime(auction.VestingSchedules[1].ReleaseTime.AddDate(0, 0, 1))
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
