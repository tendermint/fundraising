package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/fundraising/x/fundraising/types"

	_ "github.com/stretchr/testify/suite"
)

func (suite *KeeperTestSuite) TestGRPCParams() {
	resp, err := suite.querier.Params(sdk.WrapSDKContext(suite.ctx), &types.QueryParamsRequest{})
	suite.Require().NoError(err)

	suite.Require().Equal(suite.keeper.GetParams(suite.ctx), resp.Params)
}

func (suite *KeeperTestSuite) TestGRPCAuctions() {
	// TODO: not implemented yet
}

func (suite *KeeperTestSuite) TestGRPCAuction() {
	// TODO: not implemented yet
}

func (suite *KeeperTestSuite) TestGRPCBids() {
	// TODO: not implemented yet
}

func (suite *KeeperTestSuite) TestGRPCVestings() {
	// TODO: not implemented yet
}
