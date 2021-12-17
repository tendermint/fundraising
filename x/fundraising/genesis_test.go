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

func (suite *ModuleTestSuite) TestGenesisState() {
	// create auctions and reserve selling coin to reserve account
	for _, auction := range suite.sampleFixedPriceAuctions {
		suite.keeper.SetAuction(suite.ctx, auction)

		err := suite.keeper.ReserveSellingCoin(
			suite.ctx,
			auction.GetId(),
			auction.GetAuctioneer(),
			auction.GetSellingCoin(),
		)
		suite.Require().NoError(err)
	}
	suite.Require().Len(suite.keeper.GetAuctions(suite.ctx), 2)

	// make bids and reserve paying coin to reserve account
	for _, bid := range suite.sampleFixedPriceBids {
		bidderAcc, err := sdk.AccAddressFromBech32(bid.Bidder)
		suite.Require().NoError(err)
		suite.keeper.SetBid(suite.ctx, bid.AuctionId, bid.Sequence, bidderAcc, bid)

		err = suite.keeper.ReservePayingCoin(suite.ctx, bid.GetAuctionId(), bidderAcc, bid.Coin)
		suite.Require().NoError(err)
	}
	suite.Require().Len(suite.keeper.GetBids(suite.ctx), 4)

	// set the current block time a day before second auction so that it gets finished
	suite.ctx = suite.ctx.WithBlockTime(suite.sampleFixedPriceAuctions[1].GetEndTimes()[0].AddDate(0, 0, -1))
	fundraising.EndBlocker(suite.ctx, suite.keeper)

	// make first and second vesting queues over
	suite.ctx = suite.ctx.WithBlockTime(types.ParseTime("2022-04-02T00:00:00Z"))
	fundraising.EndBlocker(suite.ctx, suite.keeper)

	queues := suite.keeper.GetVestingQueuesByAuctionId(suite.ctx, 2)
	suite.Require().Len(queues, 4)

	var genState *types.GenesisState
	suite.Require().NotPanics(func() {
		genState = fundraising.ExportGenesis(suite.ctx, suite.keeper)
	})
	suite.Require().NoError(genState.Validate())

	suite.Require().NotPanics(func() {
		fundraising.InitGenesis(suite.ctx, suite.keeper, *genState)
	})
	suite.Require().Equal(genState, fundraising.ExportGenesis(suite.ctx, suite.keeper))
}

func (suite *ModuleTestSuite) TestMarshalUnmarshalDefaultGenesis() {
	genState := fundraising.ExportGenesis(suite.ctx, suite.keeper)
	bz, err := suite.app.AppCodec().MarshalJSON(genState)
	suite.Require().NoError(err)

	genState2 := types.GenesisState{}
	err = suite.app.AppCodec().UnmarshalJSON(bz, &genState2)
	suite.Require().NoError(err)
	suite.Require().Equal(*genState, genState2)

	app2 := simapp.New(app.DefaultNodeHome)
	ctx2 := app2.BaseApp.NewContext(false, tmproto.Header{})
	fundraising.InitGenesis(ctx2, app2.FundraisingKeeper, genState2)

	genState3 := fundraising.ExportGenesis(ctx2, app2.FundraisingKeeper)
	suite.Require().Equal(genState2, *genState3)
}
