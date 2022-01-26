package fundraising_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/fundraising/x/fundraising"
	"github.com/tendermint/fundraising/x/fundraising/types"

	_ "github.com/stretchr/testify/suite"
)

func (s *ModuleTestSuite) TestEndBlockerStandByStatus() {
	standByAuction := s.createFixedPriceAuction(
		s.addr(0),
		sdk.MustNewDecFromStr("0.5"),
		sdk.NewInt64Coin("denom1", 500_000_000_000),
		"denom2",
		[]types.VestingSchedule{},
		types.MustParseRFC3339("2023-01-01T00:00:00Z"),
		types.MustParseRFC3339("2023-05-01T00:00:00Z"),
		true,
	)
	s.Require().Equal(types.AuctionStatusStandBy, standByAuction.GetStatus())

	// Modify current time
	s.ctx = s.ctx.WithBlockTime(standByAuction.StartTime.AddDate(0, 0, 1))
	fundraising.EndBlocker(s.ctx, s.keeper)

	auction, found := s.keeper.GetAuction(s.ctx, standByAuction.GetId())
	s.Require().True(found)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())
}

func (s *ModuleTestSuite) TestEndBlockerStartedStatus() {
	auctioneer := s.addr(0)

	startedAuction := s.createFixedPriceAuction(
		auctioneer,
		sdk.OneDec(),
		sdk.NewInt64Coin("denom1", 500_000_000_000),
		"denom2",
		[]types.VestingSchedule{
			{
				ReleaseTime: types.MustParseRFC3339("2024-01-01T00:00:00Z"),
				Weight:      sdk.MustNewDecFromStr("0.5"),
			},
			{
				ReleaseTime: types.MustParseRFC3339("2024-06-01T00:00:00Z"),
				Weight:      sdk.MustNewDecFromStr("0.5"),
			},
		},
		types.MustParseRFC3339("2022-01-01T00:00:00Z"),
		types.MustParseRFC3339("2023-05-01T00:00:00Z"),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, startedAuction.GetStatus())

	auctionId := startedAuction.GetId()
	payingCoinDenom := startedAuction.GetPayingCoinDenom()
	sellingCoin := startedAuction.GetSellingCoin()

	bid1 := s.placeBid(auctionId, s.addr(1), sdk.OneDec(), sdk.NewInt64Coin(payingCoinDenom, 20_000_000), true)
	bid2 := s.placeBid(auctionId, s.addr(2), sdk.OneDec(), sdk.NewInt64Coin(payingCoinDenom, 20_000_000), true)
	bid3 := s.placeBid(auctionId, s.addr(3), sdk.OneDec(), sdk.NewInt64Coin(payingCoinDenom, 20_000_000), true)

	totalBidCoin := bid1.Coin.Add(bid2.Coin).Add(bid3.Coin)
	receiveAmt := totalBidCoin.Amount.ToDec().QuoTruncate(startedAuction.GetStartPrice()).TruncateInt()
	receiveCoin := sdk.NewCoin(sellingCoin.Denom, receiveAmt)

	payingReserve := s.getBalance(startedAuction.GetPayingReserveAddress(), payingCoinDenom)
	s.Require().True(coinEq(totalBidCoin, payingReserve))

	// Modify the current block time a day after the end time
	s.ctx = s.ctx.WithBlockTime(startedAuction.GetEndTimes()[0].AddDate(0, 0, 1))
	fundraising.EndBlocker(s.ctx, s.keeper)

	// The remaining selling coin must be returned to the auctioneer
	auctioneerBalance := s.getBalance(auctioneer, sellingCoin.Denom)
	s.Require().Equal(startedAuction.GetSellingCoin(), auctioneerBalance.Add(receiveCoin))
}

func (s *ModuleTestSuite) TestEndBlockerVestingStatus() {
	// s.SetAuction(s.sampleFixedPriceAuctions[1])

	// auction, found := s.keeper.GetAuction(s.ctx, 2)
	// s.Require().True(found)
	// s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	// totalBidCoin := sdk.NewInt64Coin(s.sampleFixedPriceAuctions[1].GetPayingCoinDenom(), 0)
	// for _, bid := range s.sampleFixedPriceBids {
	// 	s.PlaceBid(bid)

	// 	totalBidCoin = totalBidCoin.Add(bid.Coin)
	// }

	// // set the current block time a day after so that it gets finished
	// s.ctx = s.ctx.WithBlockTime(auction.GetEndTimes()[0].AddDate(0, 0, 1))
	// fundraising.EndBlocker(s.ctx, s.keeper)

	// vestingReserve := s.app.BankKeeper.GetBalance(
	// 	s.ctx,
	// 	auction.GetVestingReserveAddress(),
	// 	auction.GetPayingCoinDenom(),
	// )
	// s.Require().Equal(totalBidCoin, vestingReserve)

	// s.ctx = s.ctx.WithBlockTime(types.MustParseRFC3339("2022-04-02T00:00:00Z"))
	// fundraising.EndBlocker(s.ctx, s.keeper)

	// queues := s.keeper.GetVestingQueuesByAuctionId(s.ctx, auction.GetId())
	// s.Require().Len(queues, 4)
	// s.Require().True(queues[0].Released)
	// s.Require().True(queues[1].Released)
	// s.Require().False(queues[2].Released)
	// s.Require().False(queues[3].Released)

	// // auctioneer should have received two released amounts
	// auctioneerBalance := s.app.BankKeeper.GetBalance(
	// 	s.ctx,
	// 	s.addrs[5],
	// 	auction.GetPayingCoinDenom(),
	// )
	// s.Require().Equal(
	// 	totalBidCoin.Amount.Quo(sdk.NewInt(2)),
	// 	auctioneerBalance.Amount.Sub(initialBalances.AmountOf(denom4)),
	// )
}
