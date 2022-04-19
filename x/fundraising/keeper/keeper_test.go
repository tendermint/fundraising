package keeper_test

import (
	"encoding/binary"
	"fmt"
	"strings"
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

	return auction.(*types.FixedPriceAuction)
}

func (s *KeeperTestSuite) createBatchAuction(
	auctioneer sdk.AccAddress,
	startPrice sdk.Dec,
	minBidPrice sdk.Dec,
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
		MinBidPrice:       minBidPrice,
		SellingCoin:       sellingCoin,
		PayingCoinDenom:   payingCoinDenom,
		VestingSchedules:  vestingSchedules,
		MaxExtendedRound:  maxExtendedRound,
		ExtendedRoundRate: extendedRoundRate,
		StartTime:         startTime,
		EndTime:           endTime,
	})
	s.Require().NoError(err)

	return auction.(*types.BatchAuction)
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
	auction, found := s.keeper.GetAuction(s.ctx, auctionId)
	s.Require().True(found)

	var fundAmt sdk.Int
	var fundCoin sdk.Coin
	var maxBidAmt sdk.Int

	if coin.Denom == auction.GetPayingCoinDenom() {
		fundCoin = coin
		maxBidAmt = coin.Amount.ToDec().QuoTruncate(price).TruncateInt()
	} else {
		fundAmt = coin.Amount.ToDec().Mul(price).Ceil().TruncateInt()
		fundCoin = sdk.NewCoin(auction.GetPayingCoinDenom(), fundAmt)
		maxBidAmt = coin.Amount
	}

	if fund {
		s.fundAddr(bidder, sdk.NewCoins(fundCoin))
	}

	s.addAllowedBidder(auctionId, bidder, maxBidAmt)

	b, err := s.keeper.PlaceBid(s.ctx, &types.MsgPlaceBid{
		AuctionId: auctionId,
		Bidder:    bidder.String(),
		BidType:   types.BidTypeFixedPrice,
		Price:     price,
		Coin:      coin,
	})
	s.Require().NoError(err)

	return b
}

func (s *KeeperTestSuite) placeBidBatchWorth(
	auctionId uint64,
	bidder sdk.AccAddress,
	price sdk.Dec,
	coin sdk.Coin,
	maxBidAmt sdk.Int,
	fund bool,
) types.Bid {
	if fund {
		s.fundAddr(bidder, sdk.NewCoins(coin))
	}

	s.addAllowedBidder(auctionId, bidder, maxBidAmt)

	b, err := s.keeper.PlaceBid(s.ctx, &types.MsgPlaceBid{
		AuctionId: auctionId,
		Bidder:    bidder.String(),
		BidType:   types.BidTypeBatchWorth,
		Price:     price,
		Coin:      coin,
	})
	s.Require().NoError(err)

	return b
}

func (s *KeeperTestSuite) placeBidBatchMany(
	auctionId uint64,
	bidder sdk.AccAddress,
	price sdk.Dec,
	coin sdk.Coin,
	maxBidAmt sdk.Int,
	fund bool,
) types.Bid {
	auction, found := s.keeper.GetAuction(s.ctx, auctionId)
	s.Require().True(found)

	if fund {
		fundAmt := coin.Amount.ToDec().Mul(price).Ceil().TruncateInt()
		fundCoin := sdk.NewCoin(auction.GetPayingCoinDenom(), fundAmt)

		s.fundAddr(bidder, sdk.NewCoins(fundCoin))
	}

	s.addAllowedBidder(auctionId, bidder, maxBidAmt)

	b, err := s.keeper.PlaceBid(s.ctx, &types.MsgPlaceBid{
		AuctionId: auctionId,
		Bidder:    bidder.String(),
		BidType:   types.BidTypeBatchMany,
		Price:     price,
		Coin:      coin,
	})
	s.Require().NoError(err)

	return b
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

// fullString is a helper function that returns a full output of the matching result.
// it includes all bids sorted in descending order, allocation, refund, and matching info.
// it is useful for debugging.
func (s *KeeperTestSuite) fullString(auctionId uint64, mInfo keeper.MatchingInfo) string {
	auction, found := s.keeper.GetAuction(s.ctx, auctionId)
	s.Require().True(found)

	payingCoinDenom := auction.GetPayingCoinDenom()
	bids := s.keeper.GetBidsByAuctionId(s.ctx, auctionId)
	bids = types.SortByBidPrice(bids)

	var b strings.Builder

	// Bids
	b.WriteString("[Bids]\n")
	b.WriteString("+--------------------bidder---------------------+-id-+---------price---------+---------type---------+-----reserve-amount-----+-------bid-amount-------+\n")
	for _, bid := range bids {
		reserveAmt := bid.ConvertToPayingAmount(payingCoinDenom)
		bidAmt := bid.ConvertToSellingAmount(payingCoinDenom)

		_, _ = fmt.Fprintf(&b, "| %28s | %2d | %21s | %20s | %22s | %22s |\n", bid.Bidder, bid.Id, bid.Price.String(), bid.Type, reserveAmt, bidAmt)
	}
	b.WriteString("+-----------------------------------------------+----+-----------------------+----------------------+------------------------+------------------------+\n\n")

	// Allocation
	b.WriteString("[Allocation]\n")
	b.WriteString("+--------------------bidder---------------------+------allocated-amount------+\n")
	for bidder, allocatedAmt := range mInfo.AllocationMap {
		_, _ = fmt.Fprintf(&b, "| %28s | %26s |\n", bidder, allocatedAmt)
	}
	b.WriteString("+-----------------------------------------------+----------------------------+\n\n")

	// Refund
	if mInfo.RefundMap != nil {
		b.WriteString("[Refund]\n")
		b.WriteString("+--------------------bidder---------------------+------refund-amount------+\n")
		for bidder, refundAmt := range mInfo.RefundMap {
			_, _ = fmt.Fprintf(&b, "| %30s | %23s |\n", bidder, refundAmt)
		}
		b.WriteString("+-----------------------------------------------+-------------------------+\n\n")
	}

	b.WriteString("[MatchingInfo]\n")
	b.WriteString("+-matched-len-+------matched-price------+------total-matched-amount------+\n")
	_, _ = fmt.Fprintf(&b, "| %11d | %23s | %30s |\n", mInfo.MatchedLen, mInfo.MatchedPrice.String(), mInfo.TotalMatchedAmount)
	b.WriteString("+-------------+-------------------------+--------------------------------+")

	return b.String()
}

// bodSellingAmount exchanges to selling coin amount (PayingCoinAmount/Price).
func bidSellingAmount(price sdk.Dec, coin sdk.Coin) sdk.Int {
	return coin.Amount.ToDec().QuoTruncate(price).TruncateInt()
}

// parseCoin parses string and returns sdk.Coin.
func parseCoin(s string) sdk.Coin {
	s = strings.ReplaceAll(s, "_", "")
	coin, err := sdk.ParseCoinNormalized(s)
	if err != nil {
		panic(err)
	}
	return coin
}

// parseCoins parses string and returns sdk.Coins.
func parseCoins(s string) sdk.Coins {
	s = strings.ReplaceAll(s, "_", "")
	coins, err := sdk.ParseCoinsNormalized(s)
	if err != nil {
		panic(err)
	}
	return coins
}

// parseInt parses string and returns sdk.Int.
func parseInt(s string) sdk.Int {
	s = strings.ReplaceAll(s, "_", "")
	amt, ok := sdk.NewIntFromString(s)
	if !ok {
		panic("failed to convert string to sdk.Int")
	}
	return amt
}

// parseDec parses string and returns sdk.Dec.
func parseDec(s string) sdk.Dec {
	return sdk.MustNewDecFromStr(s)
}

// coinEq is a convenient method to test expected and got values of sdk.Coin.
func coinEq(exp, got sdk.Coin) (bool, string, string, string) {
	return exp.IsEqual(got), "expected:\t%v\ngot:\t\t%v", exp.String(), got.String()
}
