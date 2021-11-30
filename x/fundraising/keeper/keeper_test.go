package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/tendermint/fundraising/app"
	"github.com/tendermint/fundraising/testutil/simapp"
	"github.com/tendermint/fundraising/x/fundraising/keeper"
	"github.com/tendermint/fundraising/x/fundraising/types"
)

const (
	denom1 = "denom1"
	denom2 = "denom2"
)

var (
	initialBalances = sdk.NewCoins(
		sdk.NewInt64Coin(sdk.DefaultBondDenom, 100_000_000_000_000),
		sdk.NewInt64Coin(denom1, 100_000_000_000_000),
		sdk.NewInt64Coin(denom2, 100_000_000_000_000),
	)
)

type KeeperTestSuite struct {
	suite.Suite

	app     *app.App
	ctx     sdk.Context
	keeper  keeper.Keeper
	querier keeper.Querier
	srv     types.MsgServer
	addrs   []sdk.AccAddress
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) SetupTest() {
	app := simapp.New(app.DefaultNodeHome)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	suite.app = app
	suite.ctx = ctx
	suite.keeper = suite.app.FundraisingKeeper
	suite.querier = keeper.Querier{Keeper: suite.keeper}
	suite.srv = keeper.NewMsgServerImpl(suite.keeper)
	suite.addrs = simapp.AddTestAddrs(suite.app, suite.ctx, 5, sdk.ZeroInt())
	for _, addr := range suite.addrs {
		err := simapp.FundAccount(suite.app.BankKeeper, suite.ctx, addr, initialBalances)
		suite.Require().NoError(err)
	}

	// TODO: not implemented yet
}

func (suite *KeeperTestSuite) StartPrice(price string) sdk.Dec {
	return sdk.MustNewDecFromStr(price)
}

func (suite *KeeperTestSuite) SellingCoin(amount int64) sdk.Coin {
	return sdk.NewInt64Coin(denom2, amount)
}

func (suite *KeeperTestSuite) PayingCoinDenom() string {
	return denom1
}

func (suite *KeeperTestSuite) VestingSchedules() []types.VestingSchedule {
	return []types.VestingSchedule{
		types.NewVestingSchedule(types.ParseTime("2022-01-01T22:00:00+00:00"), sdk.MustNewDecFromStr("0.25")),
		types.NewVestingSchedule(types.ParseTime("2022-04-01T22:00:00+00:00"), sdk.MustNewDecFromStr("0.25")),
		types.NewVestingSchedule(types.ParseTime("2022-08-01T22:00:00+00:00"), sdk.MustNewDecFromStr("0.25")),
		types.NewVestingSchedule(types.ParseTime("2022-12-01T22:00:00+00:00"), sdk.MustNewDecFromStr("0.25")),
	}
}
