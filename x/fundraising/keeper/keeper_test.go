package keeper_test

import (
	"encoding/binary"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/tendermint/fundraising/app"
	"github.com/tendermint/fundraising/testutil/simapp"
	"github.com/tendermint/fundraising/x/fundraising/keeper"
	"github.com/tendermint/fundraising/x/fundraising/types"
)

type KeeperTestSuite struct {
	suite.Suite

	app       *app.App
	ctx       sdk.Context
	keeper    keeper.Keeper
	querier   keeper.Querier
	msgServer types.MsgServer
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) SetupTest() {
	app := simapp.New(app.DefaultNodeHome)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	s.app = app
	s.ctx = ctx
	s.ctx = s.ctx.WithBlockTime(time.Now()) // set to current time
	s.keeper = s.app.FundraisingKeeper
	s.querier = keeper.Querier{Keeper: s.keeper}
	s.msgServer = keeper.NewMsgServerImpl(s.keeper)
}

//
// Below are just shortcuts to frequently-used functions.
//

func (s *KeeperTestSuite) getBalance(addr sdk.AccAddress, denom string) sdk.Coin {
	return s.app.BankKeeper.GetBalance(s.ctx, addr, denom)
}

func (s *KeeperTestSuite) sendCoins(fromAddr, toAddr sdk.AccAddress, coins sdk.Coins, fund bool) {
	if fund {
		s.fundAddr(fromAddr, coins)
	}
	err := s.app.BankKeeper.SendCoins(s.ctx, fromAddr, toAddr, coins)
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) createFixedPriceAuction(
	auctioneer sdk.AccAddress,
	startPrice sdk.Dec,
	sellingCoin sdk.Coin,
	payingCoinDenom string,
	vestingSchedules []types.VestingSchedule,
	startTime time.Time,
	endTime time.Time,
	fund bool,
) *types.FixedPriceAuction {
	params := s.keeper.GetParams(s.ctx)
	if fund {
		s.fundAddr(auctioneer, params.AuctionCreationFee.Add(sellingCoin))
	}
	auction, err := s.keeper.CreateFixedPriceAuction(s.ctx, &types.MsgCreateFixedPriceAuction{
		Auctioneer:       auctioneer.String(),
		StartPrice:       startPrice,
		SellingCoin:      sellingCoin,
		PayingCoinDenom:  payingCoinDenom,
		VestingSchedules: vestingSchedules,
		StartTime:        startTime,
		EndTime:          endTime,
	})
	s.Require().NoError(err)

	return auction
}

func (s *KeeperTestSuite) placeBid(auctionId uint64, bidder sdk.AccAddress, price sdk.Dec, coin sdk.Coin, fund bool) types.Bid {
	if fund {
		s.fundAddr(bidder, sdk.NewCoins(coin))
	}

	receiveAmt := coin.Amount.ToDec().QuoTruncate(price).TruncateInt()
	s.addAllowedBidder(auctionId, bidder, receiveAmt)

	bid, err := s.keeper.PlaceBid(s.ctx, &types.MsgPlaceBid{
		AuctionId: auctionId,
		Bidder:    bidder.String(),
		BidPrice:  price,
		BidCoin:   coin,
	})
	s.Require().NoError(err)

	return bid
}

func (s *KeeperTestSuite) cancelAuction(auctionId uint64, auctioneer sdk.AccAddress) types.AuctionI {
	auction, err := s.keeper.CancelAuction(s.ctx, &types.MsgCancelAuction{
		Auctioneer: auctioneer.String(),
		AuctionId:  auctionId,
	})
	s.Require().NoError(err)

	return auction
}

//
// Below are useful helpers to write test code easily.
//

func (s *KeeperTestSuite) addr(addrNum int) sdk.AccAddress {
	addr := make(sdk.AccAddress, 20)
	binary.PutVarint(addr, int64(addrNum))
	return addr
}

func (s *KeeperTestSuite) fundAddr(addr sdk.AccAddress, coins sdk.Coins) {
	err := simapp.FundAccount(s.app.BankKeeper, s.ctx, addr, coins)
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) addAllowedBidder(auctionId uint64, bidder sdk.AccAddress, maxBidAmount sdk.Int) {
	err := s.keeper.AddAllowedBidders(s.ctx, auctionId, []types.AllowedBidder{
		{
			Bidder:       bidder.String(),
			MaxBidAmount: maxBidAmount,
		},
	})
	s.Require().NoError(err)
}

// coinEq is a convenient method to test expected and got values of sdk.Coin.
func coinEq(exp, got sdk.Coin) (bool, string, string, string) {
	return exp.IsEqual(got), "expected:\t%v\ngot:\t\t%v", exp.String(), got.String()
}
