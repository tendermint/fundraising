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

// const (
// 	denom1 = "denom1" // selling coin denom
// 	denom2 = "denom2" // paying coin denom
// 	denom3 = "denom3"
// 	denom4 = "denom4"
// )

// var (
// 	initialBalances = sdk.NewCoins(
// 		sdk.NewInt64Coin(sdk.DefaultBondDenom, 100_000_000_000_000),
// 		sdk.NewInt64Coin(denom1, 100_000_000_000_000),
// 		sdk.NewInt64Coin(denom2, 100_000_000_000_000),
// 		sdk.NewInt64Coin(denom3, 100_000_000_000_000),
// 		sdk.NewInt64Coin(denom4, 100_000_000_000_000),
// 	)
// )

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

	// s.msgServer = keeper.NewMsgServerImpl(s.keeper)
	// s.addrs = simapp.AddTestAddrs(s.app, s.ctx, 6, sdk.ZeroInt())
	// for _, addr := range s.addrs {
	// 	err := simapp.FundAccount(s.app.BankKeeper, s.ctx, addr, initialBalances)
	// 	s.Require().NoError(err)
	// }
	// s.sampleVestingSchedules1 = []types.VestingSchedule{
	// 	{
	// 		ReleaseTime: types.MustParseRFC3339("2030-01-31T22:00:00+00:00"),
	// 		Weight:      sdk.MustNewDecFromStr("0.5"),
	// 	},
	// 	{
	// 		ReleaseTime: types.MustParseRFC3339("2030-12-01T22:00:00+00:00"),
	// 		Weight:      sdk.MustNewDecFromStr("0.5"),
	// 	},
	// }
	// s.sampleVestingSchedules2 = []types.VestingSchedule{
	// 	{
	// 		ReleaseTime: types.MustParseRFC3339("2022-01-01T22:00:00+00:00"),
	// 		Weight:      sdk.MustNewDecFromStr("0.25"),
	// 	},
	// 	{
	// 		ReleaseTime: types.MustParseRFC3339("2022-05-01T22:00:00+00:00"),
	// 		Weight:      sdk.MustNewDecFromStr("0.25"),
	// 	},
	// 	{
	// 		ReleaseTime: types.MustParseRFC3339("2022-09-01T22:00:00+00:00"),
	// 		Weight:      sdk.MustNewDecFromStr("0.25"),
	// 	},
	// 	{
	// 		ReleaseTime: types.MustParseRFC3339("2022-12-01T22:00:00+00:00"),
	// 		Weight:      sdk.MustNewDecFromStr("0.25"),
	// 	},
	// }
	// s.sampleFixedPriceAuctions = []types.AuctionI{
	// 	types.NewFixedPriceAuction(
	// 		&types.BaseAuction{
	// 			Id:                    1,
	// 			Type:                  types.AuctionTypeFixedPrice,
	// 			Auctioneer:            s.addrs[4].String(),
	// 			SellingReserveAddress: types.SellingReserveAcc(1).String(),
	// 			PayingReserveAddress:  types.PayingReserveAcc(1).String(),
	// 			StartPrice:            sdk.OneDec(), // start price corresponds to the ratio of the paying coin
	// 			SellingCoin:           sdk.NewInt64Coin(denom1, 1_000_000_000_000),
	// 			PayingCoinDenom:       denom2,
	// 			VestingReserveAddress: types.VestingReserveAcc(1).String(),
	// 			VestingSchedules:      s.sampleVestingSchedules1,
	// 			WinningPrice:          sdk.ZeroDec(),
	// 			RemainingCoin:         sdk.NewInt64Coin(denom1, 1_000_000_000_000),
	// 			StartTime:             types.MustParseRFC3339("2023-01-01T00:00:00Z"),
	// 			EndTimes:              []time.Time{types.MustParseRFC3339("2023-01-10T00:00:00Z")},
	// 			Status:                types.AuctionStatusStandBy,
	// 		},
	// 	),
	// 	types.NewFixedPriceAuction(
	// 		&types.BaseAuction{
	// 			Id:                    2,
	// 			Type:                  types.AuctionTypeFixedPrice,
	// 			Auctioneer:            s.addrs[5].String(),
	// 			SellingReserveAddress: types.SellingReserveAcc(2).String(),
	// 			PayingReserveAddress:  types.PayingReserveAcc(2).String(),
	// 			StartPrice:            sdk.MustNewDecFromStr("0.5"),
	// 			SellingCoin:           sdk.NewInt64Coin(denom3, 1_000_000_000_000),
	// 			PayingCoinDenom:       denom4,
	// 			VestingReserveAddress: types.VestingReserveAcc(2).String(),
	// 			VestingSchedules:      s.sampleVestingSchedules2,
	// 			WinningPrice:          sdk.ZeroDec(),
	// 			RemainingCoin:         sdk.NewInt64Coin(denom3, 1_000_000_000_000),
	// 			StartTime:             types.MustParseRFC3339("2021-12-10T00:00:00Z"),
	// 			EndTimes:              []time.Time{types.MustParseRFC3339("2021-12-24T00:00:00Z")},
	// 			Status:                types.AuctionStatusStarted,
	// 		},
	// 	),
	// }
	// s.sampleFixedPriceBids = []types.Bid{
	// 	{
	// 		AuctionId: 2,
	// 		Sequence:  1,
	// 		Bidder:    s.addrs[0].String(),
	// 		Price:     sdk.MustNewDecFromStr("0.5"),
	// 		Coin:      sdk.NewInt64Coin(denom4, 20_000_000),
	// 		Height:    uint64(s.ctx.BlockHeight()),
	// 		Eligible:  true,
	// 	},
	// 	{
	// 		AuctionId: 2,
	// 		Sequence:  2,
	// 		Bidder:    s.addrs[0].String(),
	// 		Price:     sdk.MustNewDecFromStr("0.5"),
	// 		Coin:      sdk.NewInt64Coin(denom4, 30_000_000),
	// 		Height:    uint64(s.ctx.BlockHeight()),
	// 		Eligible:  false,
	// 	},
	// 	{
	// 		AuctionId: 2,
	// 		Sequence:  3,
	// 		Bidder:    s.addrs[1].String(),
	// 		Price:     sdk.MustNewDecFromStr("0.5"),
	// 		Coin:      sdk.NewInt64Coin(denom4, 50_000_000),
	// 		Height:    uint64(s.ctx.BlockHeight()),
	// 		Eligible:  true,
	// 	},
	// 	{
	// 		AuctionId: 2,
	// 		Sequence:  4,
	// 		Bidder:    s.addrs[1].String(),
	// 		Price:     sdk.MustNewDecFromStr("0.5"),
	// 		Coin:      sdk.NewInt64Coin(denom4, 50_000_000),
	// 		Height:    uint64(s.ctx.BlockHeight()),
	// 		Eligible:  false,
	// 	},
	// }
}

//
// Below are just shortcuts to frequently-used functions.
//

func (s *KeeperTestSuite) getBalances(addr sdk.AccAddress) sdk.Coins {
	return s.app.BankKeeper.GetAllBalances(s.ctx, addr)
}

func (s *KeeperTestSuite) getBalance(addr sdk.AccAddress, denom string) sdk.Coin {
	return s.app.BankKeeper.GetBalance(s.ctx, addr, denom)
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
	bid, err := s.keeper.PlaceBid(s.ctx, &types.MsgPlaceBid{
		AuctionId: auctionId,
		Bidder:    bidder.String(),
		Price:     price,
		Coin:      coin,
	})
	s.Require().NoError(err)

	return bid
}

func (s *KeeperTestSuite) cancelAuction(auctionId uint64, auctioneer sdk.AccAddress) {
	_, err := s.keeper.CancelAuction(s.ctx, &types.MsgCancelAuction{
		Auctioneer: auctioneer.String(),
		AuctionId:  auctionId,
	})
	s.Require().NoError(err)
}

//
// Below are useful helpers to write test code easily.
//

func (s *KeeperTestSuite) addr(addrNum int) sdk.AccAddress {
	addr := make(sdk.AccAddress, 20)
	binary.PutVarint(addr, int64(addrNum))
	return addr
}

func (s *KeeperTestSuite) fundAddr(addr sdk.AccAddress, amt sdk.Coins) {
	err := s.app.BankKeeper.MintCoins(s.ctx, types.ModuleName, amt)
	s.Require().NoError(err)
	err = s.app.BankKeeper.SendCoinsFromModuleToAccount(s.ctx, types.ModuleName, addr, amt)
	s.Require().NoError(err)
}

// // SetAuction is a convenient method to set an auction and reserve selling coin to the selling reserve account.
// func (s *KeeperTestSuite) SetAuction(auction types.AuctionI) {
// 	s.keeper.SetAuction(s.ctx, auction)
// 	err := s.keeper.ReserveSellingCoin(
// 		s.ctx,
// 		auction.GetId(),
// 		auction.GetAuctioneer(),
// 		auction.GetSellingCoin(),
// 	)
// 	s.Require().NoError(err)
// }

// // PlaceBid is a convenient method to bid and reserve paying coin to the paying reserve account.
// func (s *KeeperTestSuite) PlaceBid(bid types.Bid) {
// 	bidderAcc, err := sdk.AccAddressFromBech32(bid.Bidder)
// 	s.Require().NoError(err)

// 	nextSeq := s.keeper.GetNextSequenceWithUpdate(s.ctx, bid.AuctionId)
// 	s.keeper.SetBid(s.ctx, bid.AuctionId, nextSeq, bidderAcc, bid)

// 	err = s.keeper.ReservePayingCoin(
// 		s.ctx,
// 		bid.GetAuctionId(),
// 		bidderAcc,
// 		bid.Coin,
// 	)
// 	s.Require().NoError(err)
// }

// // PlaceBidWithCustom is a convenient method to bid with custom fields and
// // reserve paying coin to the paying reserve account.
// func (s *KeeperTestSuite) PlaceBidWithCustom(
// 	auctionId uint64,
// 	sequence uint64,
// 	bidder string,
// 	price sdk.Dec,
// 	coin sdk.Coin,
// ) {
// 	bidderAcc, err := sdk.AccAddressFromBech32(bidder)
// 	s.Require().NoError(err)

// 	s.keeper.SetBid(s.ctx, auctionId, sequence, bidderAcc, types.Bid{
// 		AuctionId: auctionId,
// 		Sequence:  sequence,
// 		Bidder:    bidderAcc.String(),
// 		Price:     price,
// 		Coin:      coin,
// 	})

// 	err = s.keeper.ReservePayingCoin(
// 		s.ctx,
// 		auctionId,
// 		bidderAcc,
// 		coin,
// 	)
// 	s.Require().NoError(err)
// }

// // CancelAuction is a convenient method to cancel the auction.
// func (s *KeeperTestSuite) CancelAuction(auction types.AuctionI) {
// 	err := s.keeper.ReleaseSellingCoin(s.ctx, auction)
// 	s.Require().NoError(err)

// 	_ = auction.SetRemainingCoin(sdk.NewCoin(auction.GetSellingCoin().Denom, sdk.ZeroInt()))
// 	_ = auction.SetStatus(types.AuctionStatusCancelled)

// 	s.keeper.SetAuction(s.ctx, auction)
// }

func parseCoins(s string) sdk.Coins {
	coins, err := sdk.ParseCoinsNormalized(s)
	if err != nil {
		panic(err)
	}
	return coins
}

// coinEq is a convenient method to test expected and got values of sdk.Coin.
func coinEq(exp, got sdk.Coin) (bool, string, string, string) {
	return exp.IsEqual(got), "expected:\t%v\ngot:\t\t%v", exp.String(), got.String()
}
