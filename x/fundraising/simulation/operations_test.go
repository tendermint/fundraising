package simulation_test

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	chain "github.com/tendermint/fundraising/app"
	"github.com/tendermint/fundraising/testutil/simapp"
	"github.com/tendermint/fundraising/x/fundraising/keeper"
	"github.com/tendermint/fundraising/x/fundraising/simulation"
	"github.com/tendermint/fundraising/x/fundraising/types"
)

type SimTestSuite struct {
	suite.Suite

	app    *chain.App
	ctx    sdk.Context
	keeper keeper.Keeper
}

func (s *SimTestSuite) SetupTest() {
	s.app = simapp.New(chain.DefaultNodeHome)
	s.ctx = s.app.BaseApp.NewContext(false, tmproto.Header{})
	s.keeper = s.app.FundraisingKeeper
}

func TestSimTestSuite(t *testing.T) {
	suite.Run(t, new(SimTestSuite))
}

func (s *SimTestSuite) TestSimulateCreateFixedPriceAuction() {
	r := rand.New(rand.NewSource(0))
	accs := s.getTestingAccounts(r, 1)

	s.app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: s.app.LastBlockHeight() + 1, AppHash: s.app.LastCommitID().Hash}})

	op := simulation.SimulateMsgCreateFixedPriceAuction(s.app.AuthKeeper, s.app.BankKeeper, s.app.FundraisingKeeper)
	opMsg, futureOps, err := op(r, s.app.BaseApp, s.ctx, accs, "")
	s.Require().NoError(err)
	s.Require().True(opMsg.OK)
	s.Require().Len(futureOps, 0)

	var msg types.MsgCreateFixedPriceAuction
	types.ModuleCdc.MustUnmarshalJSON(opMsg.Msg, &msg)

	s.Require().Equal(types.TypeMsgCreateFixedPriceAuction, msg.Type())
	s.Require().Equal(types.ModuleName, msg.Route())
	s.Require().Equal("cosmos1tp4es44j4vv8m59za3z0tm64dkmlnm8wg2frhc", msg.Auctioneer)
	s.Require().Equal("denom10", msg.SellingCoin.Denom)
	s.Require().Equal("stake", msg.PayingCoinDenom)
	s.Require().Equal(sdk.MustNewDecFromStr("5"), msg.StartPrice)
}

func (s *SimTestSuite) TestSimulateCreateBatchAuction() {
	r := rand.New(rand.NewSource(0))
	accs := s.getTestingAccounts(r, 1)

	s.app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: s.app.LastBlockHeight() + 1, AppHash: s.app.LastCommitID().Hash}})

	op := simulation.SimulateMsgCreateBatchAuction(s.app.AuthKeeper, s.app.BankKeeper, s.app.FundraisingKeeper)
	opMsg, futureOps, err := op(r, s.app.BaseApp, s.ctx, accs, "")
	s.Require().NoError(err)
	s.Require().True(opMsg.OK)
	s.Require().Len(futureOps, 0)

	var msg types.MsgCreateBatchAuction
	types.ModuleCdc.MustUnmarshalJSON(opMsg.Msg, &msg)

	s.Require().Equal(types.TypeMsgCreateBatchAuction, msg.Type())
	s.Require().Equal(types.ModuleName, msg.Route())
	s.Require().Equal("cosmos1tp4es44j4vv8m59za3z0tm64dkmlnm8wg2frhc", msg.Auctioneer)
	s.Require().Equal("denom10", msg.SellingCoin.Denom)
	s.Require().Equal("stake", msg.PayingCoinDenom)
}

func (s *SimTestSuite) TestSimulateCancelAuction() {
	r := rand.New(rand.NewSource(0))
	accs := s.getTestingAccounts(r, 1)

	s.app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: s.app.LastBlockHeight() + 1, AppHash: s.app.LastCommitID().Hash}})

	// Create a fixed price auction
	_, err := s.app.FundraisingKeeper.CreateFixedPriceAuction(s.ctx, &types.MsgCreateFixedPriceAuction{
		Auctioneer:       accs[0].Address.String(),
		StartPrice:       sdk.MustNewDecFromStr("0.5"),
		SellingCoin:      sdk.NewInt64Coin("denom1", 5000000000),
		PayingCoinDenom:  "denom2",
		VestingSchedules: []types.VestingSchedule{},
		StartTime:        s.ctx.BlockTime().AddDate(0, 1, 0),
		EndTime:          s.ctx.BlockTime().AddDate(0, 2, 0),
	})
	s.Require().NoError(err)

	op := simulation.SimulateMsgCancelAuction(s.app.AuthKeeper, s.app.BankKeeper, s.app.FundraisingKeeper)
	opMsg, futureOps, err := op(r, s.app.BaseApp, s.ctx, accs, "")
	s.Require().NoError(err)
	s.Require().True(opMsg.OK)
	s.Require().Len(futureOps, 0)

	var msg types.MsgCancelAuction
	types.ModuleCdc.MustUnmarshalJSON(opMsg.Msg, &msg)

	s.Require().Equal(types.TypeMsgCancelAuction, msg.Type())
	s.Require().Equal(types.ModuleName, msg.Route())
	s.Require().Equal("cosmos1tp4es44j4vv8m59za3z0tm64dkmlnm8wg2frhc", msg.Auctioneer)
	s.Require().Equal(uint64(1), msg.AuctionId)
}

func (s *SimTestSuite) getTestingAccounts(r *rand.Rand, n int) []simtypes.Account {
	accs := simtypes.RandomAccounts(r, n)

	initAmt := s.app.StakingKeeper.TokensFromConsensusPower(s.ctx, 200)
	coins := sdk.NewCoins(
		sdk.NewCoin(sdk.DefaultBondDenom, initAmt),
		sdk.NewInt64Coin("denom1", 1_000_000_000_000_000),
		sdk.NewInt64Coin("denom2", 1_000_000_000_000_000),
	)

	// add coins to the accounts
	for _, acc := range accs {
		acc := s.app.AuthKeeper.NewAccountWithAddress(s.ctx, acc.Address)
		s.app.AuthKeeper.SetAccount(s.ctx, acc)

		err := s.app.BankKeeper.MintCoins(s.ctx, minttypes.ModuleName, coins)
		s.Require().NoError(err)

		err = s.app.BankKeeper.SendCoinsFromModuleToAccount(s.ctx, minttypes.ModuleName, acc.GetAddress(), coins)
		s.Require().NoError(err)
	}

	return accs
}
