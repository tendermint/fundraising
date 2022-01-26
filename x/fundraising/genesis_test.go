package fundraising_test

import (
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/tendermint/fundraising/app"
	"github.com/tendermint/fundraising/testutil/simapp"
	"github.com/tendermint/fundraising/x/fundraising"
	"github.com/tendermint/fundraising/x/fundraising/types"

	_ "github.com/stretchr/testify/suite"
)

func (s *ModuleTestSuite) TestGenesisState() {
	// // create auctions and reserve selling coin to reserve account
	// for _, auction := range s.sampleFixedPriceAuctions {
	// 	s.SetAuction(auction)
	// }
	// s.Require().Len(s.keeper.GetAuctions(s.ctx), 2)

	// // make bids and reserve paying coin to reserve account
	// for _, bid := range s.sampleFixedPriceBids {
	// 	s.PlaceBid(bid)
	// }
	// s.Require().Len(s.keeper.GetBids(s.ctx), 4)

	// // set the current block time a day before second auction so that it gets finished
	// s.ctx = s.ctx.WithBlockTime(s.sampleFixedPriceAuctions[1].GetEndTimes()[0].AddDate(0, 0, -1))
	// fundraising.EndBlocker(s.ctx, s.keeper)

	// // make first and second vesting queues over
	// s.ctx = s.ctx.WithBlockTime(types.MustParseRFC3339("2022-04-02T00:00:00Z"))
	// fundraising.EndBlocker(s.ctx, s.keeper)

	// queues := s.keeper.GetVestingQueuesByAuctionId(s.ctx, 2)
	// s.Require().Len(queues, 4)

	// var genState *types.GenesisState
	// s.Require().NotPanics(func() {
	// 	genState = fundraising.ExportGenesis(s.ctx, s.keeper)
	// })
	// s.Require().NoError(genState.Validate())

	// s.Require().NotPanics(func() {
	// 	fundraising.InitGenesis(s.ctx, s.keeper, *genState)
	// })
	// s.Require().Equal(genState, fundraising.ExportGenesis(s.ctx, s.keeper))
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
