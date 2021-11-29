package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/fundraising/app"
	"github.com/tendermint/fundraising/x/fundraising/keeper"
)

type KeeperTestSuite struct {
	suite.Suite

	app     *app.App
	ctx     sdk.Context
	keeper  keeper.Keeper
	querier keeper.Querier
	// TODO: not implemented yet
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) SetupTest() {
	// TODO: not implemented yet
}
