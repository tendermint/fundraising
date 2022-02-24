package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/fundraising/x/fundraising"
	"github.com/tendermint/fundraising/x/fundraising/keeper"
	"github.com/tendermint/fundraising/x/fundraising/types"

	_ "github.com/stretchr/testify/suite"
)

func (s *KeeperTestSuite) TestSellingPoolReserveAmountInvariant() {
	// Create a fixed price auction that has started status
	auction := s.createFixedPriceAuction(
		s.addr(0),
		sdk.MustNewDecFromStr("0.5"),
		sdk.NewInt64Coin("denom1", 500_000_000_000),
		"denom2",
		[]types.VestingSchedule{},
		types.MustParseRFC3339("2022-01-01T00:00:00Z"),
		types.MustParseRFC3339("2022-06-10T00:00:00Z"),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	_, broken := keeper.SellingPoolReserveAmountInvariant(s.keeper)(s.ctx)
	s.Require().False(broken)

	// Although it is not possible for an exploiter to have the same token denom in reality,
	// it is safe to test the case anyway
	exploiterAddr := s.addr(1)
	sellingReserveAddr := auction.GetSellingReserveAddress()
	s.sendCoins(exploiterAddr, sellingReserveAddr, sdk.NewCoins(
		sdk.NewInt64Coin("denom1", 500_000_000),
		sdk.NewInt64Coin("denom2", 500_000_000),
		sdk.NewInt64Coin("denom3", 500_000_000),
		sdk.NewInt64Coin("denom4", 500_000_000),
	), true)

	_, broken = keeper.SellingPoolReserveAmountInvariant(s.keeper)(s.ctx)
	s.Require().False(broken)
}

func (s *KeeperTestSuite) TestPayingPoolReserveAmountInvariant() {
	k, ctx := s.keeper, s.ctx

	auction := s.createFixedPriceAuction(
		s.addr(0),
		sdk.OneDec(),
		sdk.NewInt64Coin("denom3", 500_000_000_000),
		"denom4",
		[]types.VestingSchedule{},
		types.MustParseRFC3339("2022-01-01T00:00:00Z"),
		types.MustParseRFC3339("2022-03-10T00:00:00Z"),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	s.addAllowedBidder(auction.Id, s.addr(1), exchangedSellingAmount(parseDec("1"), parseCoin("20000000denom4")))
	s.addAllowedBidder(auction.Id, s.addr(2), exchangedSellingAmount(parseDec("1"), parseCoin("35000000denom4")))
	s.addAllowedBidder(auction.Id, s.addr(3), exchangedSellingAmount(parseDec("1"), parseCoin("35000000denom4")))

	s.placeBid(auction.GetId(), s.addr(1), types.BidTypeFixedPrice, sdk.OneDec(), parseCoin("20000000denom4"), true)
	s.placeBid(auction.GetId(), s.addr(2), types.BidTypeFixedPrice, sdk.OneDec(), parseCoin("20000000denom4"), true)
	s.placeBid(auction.GetId(), s.addr(2), types.BidTypeFixedPrice, sdk.OneDec(), parseCoin("15000000denom4"), true)
	s.placeBid(auction.GetId(), s.addr(3), types.BidTypeFixedPrice, sdk.OneDec(), parseCoin("35000000denom4"), true)

	_, broken := keeper.PayingPoolReserveAmountInvariant(k)(ctx)
	s.Require().False(broken)

	// Although it is not possible for an exploiter to have the same token denom in reality,
	// it is safe to test the case anyway
	exploiterAddr := s.addr(1)
	payingReserveAddr := auction.GetPayingReserveAddress()
	s.sendCoins(exploiterAddr, payingReserveAddr, sdk.NewCoins(
		sdk.NewInt64Coin("denom1", 500_000_000),
		sdk.NewInt64Coin("denom2", 500_000_000),
		sdk.NewInt64Coin("denom3", 500_000_000),
		sdk.NewInt64Coin("denom4", 500_000_000),
	), true)

	_, broken = keeper.PayingPoolReserveAmountInvariant(k)(ctx)
	s.Require().False(broken)
}

func (s *KeeperTestSuite) TestVestingPoolReserveAmountInvariant() {
	k, ctx := s.keeper, s.ctx

	auction := s.createFixedPriceAuction(
		s.addr(0),
		sdk.OneDec(),
		sdk.NewInt64Coin("denom3", 500_000_000_000),
		"denom4",
		[]types.VestingSchedule{
			{
				ReleaseTime: types.MustParseRFC3339("2023-01-01T00:00:00Z"),
				Weight:      sdk.MustNewDecFromStr("0.25"),
			},
			{
				ReleaseTime: types.MustParseRFC3339("2023-05-01T00:00:00Z"),
				Weight:      sdk.MustNewDecFromStr("0.25"),
			},
			{
				ReleaseTime: types.MustParseRFC3339("2023-09-01T00:00:00Z"),
				Weight:      sdk.MustNewDecFromStr("0.25"),
			},
			{
				ReleaseTime: types.MustParseRFC3339("2023-12-01T00:00:00Z"),
				Weight:      sdk.MustNewDecFromStr("0.25"),
			},
		},
		types.MustParseRFC3339("2022-01-01T00:00:00Z"),
		types.MustParseRFC3339("2022-03-10T00:00:00Z"),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	s.addAllowedBidder(auction.Id, s.addr(1), exchangedSellingAmount(parseDec("1"), parseCoin("20000000denom4")))
	s.addAllowedBidder(auction.Id, s.addr(2), exchangedSellingAmount(parseDec("1"), parseCoin("35000000denom4")))
	s.addAllowedBidder(auction.Id, s.addr(3), exchangedSellingAmount(parseDec("1"), parseCoin("35000000denom4")))

	s.placeBid(auction.GetId(), s.addr(1), types.BidTypeFixedPrice, sdk.OneDec(), parseCoin("20000000denom4"), true)
	s.placeBid(auction.GetId(), s.addr(2), types.BidTypeFixedPrice, sdk.OneDec(), parseCoin("20000000denom4"), true)
	s.placeBid(auction.GetId(), s.addr(2), types.BidTypeFixedPrice, sdk.OneDec(), parseCoin("15000000denom4"), true)
	s.placeBid(auction.GetId(), s.addr(3), types.BidTypeFixedPrice, sdk.OneDec(), parseCoin("35000000denom4"), true)

	// Make the auction ended
	ctx = ctx.WithBlockTime(auction.GetEndTimes()[0].AddDate(0, 0, 1))
	fundraising.EndBlocker(ctx, k)

	// Make first and second vesting queues over
	ctx = ctx.WithBlockTime(auction.GetVestingSchedules()[0].GetReleaseTime().AddDate(0, 0, 1))
	fundraising.EndBlocker(ctx, k)

	_, broken := keeper.VestingPoolReserveAmountInvariant(k)(ctx)
	s.Require().False(broken)

	// Although it is not possible for an exploiter to have the same token denom in reality,
	// it is safe to test the case anyway
	exploiterAddr := s.addr(1)
	vestingReserveAddr := auction.GetVestingReserveAddress()
	s.sendCoins(exploiterAddr, vestingReserveAddr, sdk.NewCoins(
		sdk.NewInt64Coin("denom1", 500_000_000),
		sdk.NewInt64Coin("denom2", 500_000_000),
		sdk.NewInt64Coin("denom3", 500_000_000),
		sdk.NewInt64Coin("denom4", 500_000_000),
	), true)

	_, broken = keeper.VestingPoolReserveAmountInvariant(k)(ctx)
	s.Require().False(broken)
}

func (s *KeeperTestSuite) TestAuctionStatusStatesInvariant() {
	k, ctx := s.keeper, s.ctx

	standByAuction := s.createFixedPriceAuction(
		s.addr(0),
		sdk.MustNewDecFromStr("0.35"),
		sdk.NewInt64Coin("denom1", 500_000_000_000),
		"denom2",
		[]types.VestingSchedule{},
		types.MustParseRFC3339("2023-01-01T00:00:00Z"),
		types.MustParseRFC3339("2023-03-01T00:00:00Z"),
		true,
	)
	s.Require().Equal(types.AuctionStatusStandBy, standByAuction.GetStatus())

	_, broken := keeper.AuctionStatusStatesInvariant(k)(ctx)
	s.Require().False(broken)

	startedAuction := s.createFixedPriceAuction(
		s.addr(1),
		sdk.MustNewDecFromStr("0.5"),
		sdk.NewInt64Coin("denom3", 500_000_000_000),
		"denom4",
		[]types.VestingSchedule{
			{
				ReleaseTime: types.MustParseRFC3339("2023-01-01T00:00:00Z"),
				Weight:      sdk.MustNewDecFromStr("0.25"),
			},
			{
				ReleaseTime: types.MustParseRFC3339("2023-05-01T00:00:00Z"),
				Weight:      sdk.MustNewDecFromStr("0.25"),
			},
			{
				ReleaseTime: types.MustParseRFC3339("2023-09-01T00:00:00Z"),
				Weight:      sdk.MustNewDecFromStr("0.25"),
			},
			{
				ReleaseTime: types.MustParseRFC3339("2023-12-01T00:00:00Z"),
				Weight:      sdk.MustNewDecFromStr("0.25"),
			},
		},
		types.MustParseRFC3339("2022-01-01T00:00:00Z"),
		types.MustParseRFC3339("2022-03-01T00:00:00Z"),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, startedAuction.GetStatus())

	_, broken = keeper.AuctionStatusStatesInvariant(k)(ctx)
	s.Require().False(broken)

	// set the current block time a day after so that it gets finished
	ctx = ctx.WithBlockTime(startedAuction.GetEndTimes()[0].AddDate(0, 0, 1))
	fundraising.EndBlocker(ctx, k)

	_, broken = keeper.AuctionStatusStatesInvariant(k)(ctx)
	s.Require().False(broken)

	// set the current block time a day after so that all vesting queues get released
	ctx = ctx.WithBlockTime(startedAuction.GetVestingSchedules()[3].GetReleaseTime().AddDate(0, 0, 1))
	fundraising.EndBlocker(ctx, k)

	_, broken = keeper.AuctionStatusStatesInvariant(k)(ctx)
	s.Require().False(broken)
}
