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
	s.app = simapp.New(app.DefaultNodeHome)
	s.ctx = s.app.BaseApp.NewContext(false, tmproto.Header{})
	s.ctx = s.ctx.WithBlockTime(time.Now()) // set to current time
	s.keeper = s.app.FundraisingKeeper
	s.querier = keeper.Querier{Keeper: s.keeper}
	s.msgServer = keeper.NewMsgServerImpl(s.keeper)
}

//
// Below are just shortcuts to frequently-used functions.
//

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

func (s *KeeperTestSuite) createBatchAuction(
	auctioneer sdk.AccAddress,
	startPrice sdk.Dec,
	sellingCoin sdk.Coin,
	payingCoinDenom string,
	vestingSchedules []types.VestingSchedule,
	maxExtendedRound uint32,
	extendedRoundRate sdk.Dec,
	startTime time.Time,
	endTime time.Time,
	fund bool,
) *types.BatchAuction {
	params := s.keeper.GetParams(s.ctx)
	if fund {
		s.fundAddr(auctioneer, params.AuctionCreationFee.Add(sellingCoin))
	}

	auction, err := s.keeper.CreateBatchAuction(s.ctx, &types.MsgCreateBatchAuction{
		Auctioneer:        auctioneer.String(),
		StartPrice:        startPrice,
		SellingCoin:       sellingCoin,
		PayingCoinDenom:   payingCoinDenom,
		VestingSchedules:  vestingSchedules,
		MaxExtendedRound:  maxExtendedRound,
		ExtendedRoundRate: extendedRoundRate,
		StartTime:         startTime,
		EndTime:           endTime,
	})
	s.Require().NoError(err)

	return auction
}

func (s *KeeperTestSuite) addAllowedBidder(auctionId uint64, bidder sdk.AccAddress, maxBidAmt sdk.Int) {
	auction, found := s.keeper.GetAuction(s.ctx, auctionId)
	s.Require().True(found)

	prevMaxBidAmt, found := auction.GetAllowedBiddersMap()[bidder.String()]
	if found {
		maxBidAmt = maxBidAmt.Add(prevMaxBidAmt)
	}

	err := s.keeper.AddAllowedBidders(s.ctx, auctionId, []types.AllowedBidder{
		{Bidder: bidder.String(), MaxBidAmount: maxBidAmt},
	})
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) placeBidFixedPrice(
	auctionId uint64,
	bidder sdk.AccAddress,
	price sdk.Dec,
	coin sdk.Coin,
	fund bool,
) types.Bid {
	fundAmt := coin.Amount
	fundCoin := coin

	if fund {
		s.fundAddr(bidder, sdk.NewCoins(fundCoin))
	}

	s.addAllowedBidder(auctionId, bidder, fundAmt)

	bid, err := s.keeper.PlaceBid(s.ctx, &types.MsgPlaceBid{
		AuctionId: auctionId,
		Bidder:    bidder.String(),
		BidType:   types.BidTypeFixedPrice,
		Price:     price,
		Coin:      coin,
	})
	s.Require().NoError(err)

	return bid
}

func (s *KeeperTestSuite) placeBidBatchWorth(
	auctionId uint64,
	bidder sdk.AccAddress,
	price sdk.Dec,
	coin sdk.Coin,
	fund bool,
) types.Bid {
	fundAmt := coin.Amount
	fundCoin := coin

	if fund {
		s.fundAddr(bidder, sdk.NewCoins(fundCoin))
	}

	s.addAllowedBidder(auctionId, bidder, fundAmt)

	bid, err := s.keeper.PlaceBid(s.ctx, &types.MsgPlaceBid{
		AuctionId: auctionId,
		Bidder:    bidder.String(),
		BidType:   types.BidTypeBatchWorth,
		Price:     price,
		Coin:      coin,
	})
	s.Require().NoError(err)

	return bid
}

func (s *KeeperTestSuite) placeBidBatchMany(
	auctionId uint64,
	bidder sdk.AccAddress,
	price sdk.Dec,
	coin sdk.Coin,
	fund bool,
) types.Bid {
	auction, found := s.keeper.GetAuction(s.ctx, auctionId)
	s.Require().True(found)

	fundAmt := coin.Amount.ToDec().Mul(price).Ceil().TruncateInt()
	fundCoin := sdk.NewCoin(auction.GetPayingCoinDenom(), fundAmt)

	if fund {
		s.fundAddr(bidder, sdk.NewCoins(fundCoin))
	}

	s.addAllowedBidder(auctionId, bidder, fundAmt)

	bid, err := s.keeper.PlaceBid(s.ctx, &types.MsgPlaceBid{
		AuctionId: auctionId,
		Bidder:    bidder.String(),
		BidType:   types.BidTypeBatchMany,
		Price:     price,
		Coin:      coin,
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

// exchangedSellingAmount exchanges to selling coin amount (PayingCoinAmount/Price).
func exchangedSellingAmount(price sdk.Dec, coin sdk.Coin) sdk.Int {
	return coin.Amount.ToDec().QuoTruncate(price).TruncateInt()
}

// parseCoin parses and returns sdk.Coin.
func parseCoin(s string) sdk.Coin {
	coin, err := sdk.ParseCoinNormalized(s)
	if err != nil {
		panic(err)
	}
	return coin
}

// parseCoins parses and returns sdk.Coins.
func parseCoins(s string) sdk.Coins {
	coins, err := sdk.ParseCoinsNormalized(s)
	if err != nil {
		panic(err)
	}
	return coins
}

// parseDec is a shortcut for sdk.MustNewDecFromStr.
func parseDec(s string) sdk.Dec {
	return sdk.MustNewDecFromStr(s)
}

// coinEq is a convenient method to test expected and got values of sdk.Coin.
func coinEq(exp, got sdk.Coin) (bool, string, string, string) {
	return exp.IsEqual(got), "expected:\t%v\ngot:\t\t%v", exp.String(), got.String()
}
