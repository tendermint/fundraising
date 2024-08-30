package keeper_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	tmrand "github.com/cometbft/cometbft/libs/rand"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/tendermint/fundraising/app"
	testkeeper "github.com/tendermint/fundraising/testutil/keeper"
	"github.com/tendermint/fundraising/testutil/sample"
	"github.com/tendermint/fundraising/testutil/testutil/simapp"
	"github.com/tendermint/fundraising/x/fundraising/keeper"
	"github.com/tendermint/fundraising/x/fundraising/types"
)

const maxAddress = 100

type KeeperTestSuite struct {
	suite.Suite

	app       *app.App
	ctx       sdk.Context
	keeper    keeper.Keeper
	msgServer types.MsgServer
	addresses []sdk.AccAddress
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) SetupTest() {
	for i := 0; i < maxAddress; i++ {
		s.addresses = append(s.addresses, sample.AccAddress())
	}
	chainID := "chain-" + tmrand.NewRand().Str(6)

	var err error
	s.app, err = simapp.New(chainID)
	s.Require().NoError(err)

	s.ctx = s.app.BaseApp.NewContext(false)
	s.ctx = s.ctx.WithBlockTime(time.Now()) // set to current time
	s.keeper = s.app.FundraisingKeeper
	s.msgServer = keeper.NewMsgServerImpl(s.keeper)
}

//
// Below are just shortcuts to frequently-used functions.
//

func (s *KeeperTestSuite) createFixedPriceAuction(
	auctioneer sdk.AccAddress,
	startPrice math.LegacyDec,
	sellingCoin sdk.Coin,
	payingCoinDenom string,
	vestingSchedules []types.VestingSchedule,
	startTime time.Time,
	endTime time.Time,
	fund bool,
) *types.FixedPriceAuction {
	params, err := s.keeper.Params.Get(s.ctx)
	s.Require().NoError(err)

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
	startPrice math.LegacyDec,
	minBidPrice math.LegacyDec,
	sellingCoin sdk.Coin,
	payingCoinDenom string,
	vestingSchedules []types.VestingSchedule,
	maxExtendedRound uint32,
	extendedRoundRate math.LegacyDec,
	startTime time.Time,
	endTime time.Time,
	fund bool,
) *types.BatchAuction {
	params, err := s.keeper.Params.Get(s.ctx)
	s.Require().NoError(err)

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

func (s *KeeperTestSuite) addAllowedBidder(auctionId uint64, bidder sdk.AccAddress, maxBidAmt math.Int) error {
	allowedBidder, err := s.keeper.AllowedBidder.Get(s.ctx, collections.Join(auctionId, bidder))
	if err == nil {
		maxBidAmt = maxBidAmt.Add(allowedBidder.MaxBidAmount)
	}
	if err != nil && !errors.Is(err, collections.ErrNotFound) {
		return err
	}

	return s.keeper.AllowedBidder.Set(s.ctx, collections.Join(auctionId, bidder), types.NewAllowedBidder(auctionId, bidder, maxBidAmt))
}

func (s *KeeperTestSuite) placeBidFixedPrice(
	auctionId uint64,
	bidder sdk.AccAddress,
	price math.LegacyDec,
	coin sdk.Coin,
	fund bool,
) types.Bid {
	auction, err := s.keeper.Auction.Get(s.ctx, auctionId)
	s.Require().NoError(err)

	var fundAmt math.Int
	var fundCoin sdk.Coin
	var maxBidAmt math.Int

	if coin.Denom == auction.GetPayingCoinDenom() {
		fundCoin = coin
		maxBidAmt = math.LegacyNewDecFromInt(coin.Amount).QuoTruncate(price).TruncateInt()
	} else {
		fundAmt = math.LegacyNewDecFromInt(coin.Amount).Mul(price).Ceil().TruncateInt()
		fundCoin = sdk.NewCoin(auction.GetPayingCoinDenom(), fundAmt)
		maxBidAmt = coin.Amount
	}

	if fund {
		s.fundAddr(bidder, sdk.NewCoins(fundCoin))
	}

	err = s.addAllowedBidder(auctionId, bidder, maxBidAmt)
	s.Require().NoError(err)

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
	price math.LegacyDec,
	coin sdk.Coin,
	maxBidAmt math.Int,
	fund bool,
) types.Bid {
	if fund {
		s.fundAddr(bidder, sdk.NewCoins(coin))
	}

	err := s.addAllowedBidder(auctionId, bidder, maxBidAmt)
	s.Require().NoError(err)

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
	price math.LegacyDec,
	coin sdk.Coin,
	maxBidAmt math.Int,
	fund bool,
) types.Bid {
	auction, err := s.keeper.Auction.Get(s.ctx, auctionId)
	s.Require().NoError(err)

	if fund {
		fundAmt := math.LegacyNewDecFromInt(coin.Amount).Mul(price).Ceil().TruncateInt()
		fundCoin := sdk.NewCoin(auction.GetPayingCoinDenom(), fundAmt)

		s.fundAddr(bidder, sdk.NewCoins(fundCoin))
	}

	err = s.addAllowedBidder(auctionId, bidder, maxBidAmt)
	s.Require().NoError(err)

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

// Below are useful helpers to write test code easily.
func (s *KeeperTestSuite) addr(index int) sdk.AccAddress {
	if index >= maxAddress {
		panic(fmt.Sprintf("invalid address index %d", index))
	}
	return s.addresses[index]
}

func (s *KeeperTestSuite) fundAddr(addr sdk.AccAddress, coins sdk.Coins) {
	err := testkeeper.FundAccount(s.app.BankKeeper, s.ctx, addr, coins)
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
	auction, err := s.keeper.Auction.Get(s.ctx, auctionId)
	s.Require().NoError(err)

	payingCoinDenom := auction.GetPayingCoinDenom()
	bids, err := s.keeper.GetBidsByAuctionId(s.ctx, auctionId)
	s.Require().NoError(err)

	bids = types.SortBids(bids)

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
func bidSellingAmount(price math.LegacyDec, coin sdk.Coin) math.Int {
	return math.LegacyNewDecFromInt(coin.Amount).QuoTruncate(price).TruncateInt()
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

// parseInt parses string and returns math.Int.
func parseInt(s string) math.Int {
	s = strings.ReplaceAll(s, "_", "")
	amt, ok := math.NewIntFromString(s)
	if !ok {
		panic("failed to convert string to math.Int")
	}
	return amt
}

// parseDec parses string and returns math.LegacyDec.
func parseDec(s string) math.LegacyDec {
	return math.LegacyMustNewDecFromStr(s)
}

// coinEq is a convenient method to test expected and got values of sdk.Coin.
func coinEq(exp, got sdk.Coin) (bool, string, string, string) {
	return exp.Equal(got), "expected:\t%v\ngot:\t\t%v", exp.String(), got.String()
}
