package fundraising_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/tendermint/fundraising/app"
	"github.com/tendermint/fundraising/testutil/simapp"
	"github.com/tendermint/fundraising/x/fundraising"
	"github.com/tendermint/fundraising/x/fundraising/types"

	_ "github.com/stretchr/testify/suite"
)

func (s *ModuleTestSuite) TestGenesisState() {
	auction := s.createFixedPriceAuction(
		s.addr(0),
		sdk.OneDec(),
		sdk.NewInt64Coin("denom1", 200_000_000_000),
		"denom2",
		[]types.VestingSchedule{
			{
				ReleaseTime: types.MustParseRFC3339("2023-01-01T00:00:00Z"),
				Weight:      sdk.MustNewDecFromStr("0.25"),
			},
			{
				ReleaseTime: types.MustParseRFC3339("2023-06-01T00:00:00Z"),
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
		},
		types.MustParseRFC3339("2022-01-01T00:00:00Z"),
		types.MustParseRFC3339("2022-05-21T00:00:00Z"),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	// Place bids
	s.placeBid(auction.GetId(), s.addr(1), sdk.OneDec(), sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 20_000_000), true)
	s.placeBid(auction.GetId(), s.addr(2), sdk.OneDec(), sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 30_000_000), true)
	s.placeBid(auction.GetId(), s.addr(3), sdk.OneDec(), sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 15_000_000), true)
	s.placeBid(auction.GetId(), s.addr(4), sdk.OneDec(), sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 35_000_000), true)

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
		genState = fundraising.ExportGenesis(s.ctx, s.keeper)
	})
	s.Require().NoError(genState.Validate())

	s.Require().NotPanics(func() {
		fundraising.InitGenesis(s.ctx, s.keeper, *genState)
	})
	s.Require().Equal(genState, fundraising.ExportGenesis(s.ctx, s.keeper))
}

func (s *ModuleTestSuite) TestMarshalUnmarshalDefaultGenesis() {
	genState := fundraising.ExportGenesis(s.ctx, s.keeper)
	bz, err := s.app.AppCodec().MarshalJSON(genState)
	s.Require().NoError(err)

	genState2 := types.GenesisState{}
	err = s.app.AppCodec().UnmarshalJSON(bz, &genState2)
	s.Require().NoError(err)
	s.Require().Equal(*genState, genState2)

	app2 := simapp.New(app.DefaultNodeHome)
	ctx2 := app2.BaseApp.NewContext(false, tmproto.Header{})
	fundraising.InitGenesis(ctx2, app2.FundraisingKeeper, genState2)

	genState3 := fundraising.ExportGenesis(ctx2, app2.FundraisingKeeper)
	s.Require().Equal(genState2, *genState3)
}
