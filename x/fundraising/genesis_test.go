package fundraising_test

import (
	"testing"

	"github.com/tendermint/fundraising/x/fundraising"
	"github.com/tendermint/fundraising/x/fundraising/types"

	_ "github.com/stretchr/testify/suite"
)

func TestGenesis(t *testing.T) {
	// genesisState := types.GenesisState{
	// this line is used by starport scaffolding # genesis/test/state
	// }

	// ctx, k := keepertest.Fundraising(t)
	// fundraising.InitGenesis(ctx, *k, genesisState)
	// got := fundraising.ExportGenesis(ctx, *k)
	// require.NotNil(t, got)

	// this line is used by starport scaffolding # genesis/test/assert
}

func (suite *ModuleTestSuite) TestInitGenesis() {

}

func (suite *ModuleTestSuite) TestExportGenesis() {
	for _, auction := range suite.sampleFixedPriceAuctions {
		suite.keeper.SetAuction(suite.ctx, auction)
	}

	auction, found := suite.keeper.GetAuction(suite.ctx, 2)
	suite.Require().True(found)
	suite.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	// for _, bid := range suite.sampleFixedPriceBids {
	// 	err := suite.keeper.PlaceBid(suite.ctx, bid)
	// 	suite.Require().NoError(err)
	// }

	genState := fundraising.ExportGenesis(suite.ctx, suite.keeper)
	bz, err := suite.app.AppCodec().MarshalJSON(genState)
	suite.Require().NoError(err)

	*genState = types.GenesisState{}
	err = suite.app.AppCodec().UnmarshalJSON(bz, genState)
	suite.Require().NoError(err)
}
