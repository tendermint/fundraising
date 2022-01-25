package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/fundraising/x/fundraising"
	"github.com/tendermint/fundraising/x/fundraising/keeper"
	"github.com/tendermint/fundraising/x/fundraising/types"

	_ "github.com/stretchr/testify/suite"
)

func (s *KeeperTestSuite) TestSellingPoolReserveAmountInvariant() {
	k, ctx, auction := s.keeper, s.ctx, s.sampleFixedPriceAuctions[1]

	k.SetAuction(s.ctx, auction)

	_, broken := keeper.SellingPoolReserveAmountInvariant(k)(ctx)
	s.Require().True(broken)

	err := k.ReserveSellingCoin(
		ctx,
		auction.GetId(),
		auction.GetAuctioneer(),
		auction.GetSellingCoin(),
	)
	s.Require().NoError(err)

	_, broken = keeper.SellingPoolReserveAmountInvariant(k)(ctx)
	s.Require().False(broken)

	// in reality, although it is not possible for an exploiter to have the same token denom
	// but it is safe to test the case anyway
	exploiterAcc := s.addrs[2]
	sendingCoins := sdk.NewCoins(
		sdk.NewInt64Coin(denom1, 500_000_000),
		sdk.NewInt64Coin(denom2, 500_000_000),
		sdk.NewInt64Coin(denom3, 500_000_000),
		sdk.NewInt64Coin(denom4, 500_000_000),
	)
	err = s.app.BankKeeper.SendCoins(ctx, exploiterAcc, auction.GetSellingReserveAddress(), sendingCoins)
	s.Require().NoError(err)

	_, broken = keeper.SellingPoolReserveAmountInvariant(k)(ctx)
	s.Require().False(broken)
}

func (s *KeeperTestSuite) TestPayingPoolReserveAmountInvariant() {
	k, ctx, auction := s.keeper, s.ctx, s.sampleFixedPriceAuctions[1]

	s.SetAuction(auction)

	for _, bid := range s.sampleFixedPriceBids {
		bidderAcc, err := sdk.AccAddressFromBech32(bid.Bidder)
		s.Require().NoError(err)
		k.SetBid(ctx, bid.AuctionId, bid.Sequence, bidderAcc, bid)

		err = k.ReservePayingCoin(ctx, bid.GetAuctionId(), bidderAcc, bid.Coin)
		s.Require().NoError(err)
	}

	_, broken := keeper.PayingPoolReserveAmountInvariant(k)(ctx)
	s.Require().False(broken)

	// in reality, although it is not possible for an exploiter to have the same token denom
	// but it is safe to test the case anyway
	exploiterAcc := s.addrs[2]
	sendingCoins := sdk.NewCoins(
		sdk.NewInt64Coin(denom1, 500_000_000),
		sdk.NewInt64Coin(denom2, 500_000_000),
		sdk.NewInt64Coin(denom3, 500_000_000),
		sdk.NewInt64Coin(denom4, 500_000_000),
	)
	err := s.app.BankKeeper.SendCoins(ctx, exploiterAcc, auction.GetPayingReserveAddress(), sendingCoins)
	s.Require().NoError(err)

	_, broken = keeper.PayingPoolReserveAmountInvariant(k)(ctx)
	s.Require().False(broken)
}

func (s *KeeperTestSuite) TestVestingPoolReserveAmountInvariant() {
	k, ctx, auction := s.keeper, s.ctx, s.sampleFixedPriceAuctions[1]

	s.SetAuction(auction)

	for _, bid := range s.sampleFixedPriceBids {
		bidderAcc, err := sdk.AccAddressFromBech32(bid.Bidder)
		s.Require().NoError(err)
		k.SetBid(ctx, bid.AuctionId, bid.Sequence, bidderAcc, bid)

		err = k.ReservePayingCoin(ctx, bid.GetAuctionId(), bidderAcc, bid.Coin)
		s.Require().NoError(err)
	}

	// set the current block time a day before second auction so that it gets finished
	ctx = ctx.WithBlockTime(s.sampleFixedPriceAuctions[1].GetEndTimes()[0].AddDate(0, 0, -1))
	fundraising.EndBlocker(ctx, k)

	// make first and second vesting queues over
	ctx = ctx.WithBlockTime(types.MustParseRFC3339("2022-04-02T00:00:00Z"))
	fundraising.EndBlocker(ctx, k)

	_, broken := keeper.VestingPoolReserveAmountInvariant(k)(ctx)
	s.Require().False(broken)

	// in reality, although it is not possible for an exploiter to have the same token denom
	// but it is safe to test the case anyway
	exploiterAcc := s.addrs[2]
	sendingCoins := sdk.NewCoins(
		sdk.NewInt64Coin(denom1, 500_000_000),
		sdk.NewInt64Coin(denom2, 500_000_000),
		sdk.NewInt64Coin(denom3, 500_000_000),
		sdk.NewInt64Coin(denom4, 500_000_000),
	)
	err := s.app.BankKeeper.SendCoins(ctx, exploiterAcc, auction.GetVestingReserveAddress(), sendingCoins)
	s.Require().NoError(err)

	_, broken = keeper.VestingPoolReserveAmountInvariant(k)(ctx)
	s.Require().False(broken)
}

func (s *KeeperTestSuite) TestAuctionStatusStatesInvariant() {
	k, ctx := s.keeper, s.ctx

	s.SetAuction(s.sampleFixedPriceAuctions[0])

	_, broken := keeper.AuctionStatusStatesInvariant(k)(ctx)
	s.Require().False(broken)

	s.SetAuction(s.sampleFixedPriceAuctions[1])

	_, broken = keeper.AuctionStatusStatesInvariant(k)(ctx)
	s.Require().False(broken)

	// set the current block time a day after so that it gets finished
	ctx = ctx.WithBlockTime(s.sampleFixedPriceAuctions[1].GetEndTimes()[0].AddDate(0, 0, 1))
	fundraising.EndBlocker(ctx, k)

	_, broken = keeper.AuctionStatusStatesInvariant(k)(ctx)
	s.Require().False(broken)

	// set the current block time a day after so that all vesting queues get released
	ctx = ctx.WithBlockTime(s.sampleFixedPriceAuctions[1].GetVestingSchedules()[3].GetReleaseTime().AddDate(0, 0, 1))
	fundraising.EndBlocker(ctx, k)

	_, broken = keeper.AuctionStatusStatesInvariant(k)(ctx)
	s.Require().False(broken)
}
