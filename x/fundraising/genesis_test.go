package fundraising_test

import (
	"time"

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
	s.placeBid(auction.GetId(), s.addr(1), types.BidTypeFixedPrice, sdk.OneDec(), parseCoin("20000000denom2"), true)
	s.placeBid(auction.GetId(), s.addr(2), types.BidTypeFixedPrice, sdk.OneDec(), parseCoin("30000000denom2"), true)
	s.placeBid(auction.GetId(), s.addr(3), types.BidTypeFixedPrice, sdk.OneDec(), parseCoin("15000000denom2"), true)
	s.placeBid(auction.GetId(), s.addr(4), types.BidTypeFixedPrice, sdk.OneDec(), parseCoin("35000000denom2"), true)

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
