package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/fundraising/x/fundraising/types"

	_ "github.com/stretchr/testify/suite"
)

func (s *KeeperTestSuite) TestBidIterators() {
	// create a fixed price auction with already started status
	s.SetAuction(s.sampleFixedPriceAuctions[1])

	auction, found := s.keeper.GetAuction(s.ctx, s.sampleFixedPriceAuctions[1].GetId())
	s.Require().True(found)

	for _, bid := range s.sampleFixedPriceBids {
		s.PlaceBid(bid)
	}

	bids := s.keeper.GetBids(s.ctx)
	s.Require().Len(bids, 4)

	bidsById := s.keeper.GetBidsByAuctionId(s.ctx, auction.GetId())
	s.Require().Len(bidsById, 4)

	bidsByBidder := s.keeper.GetBidsByBidder(s.ctx, s.addrs[0])
	s.Require().Len(bidsByBidder, 2)
}

func (s *KeeperTestSuite) TestBidSequence() {
	s.SetAuction(s.sampleFixedPriceAuctions[1])

	for _, bid := range s.sampleFixedPriceBids {
		s.PlaceBid(bid)
	}

	auction, found := s.keeper.GetAuction(s.ctx, 2)
	s.Require().True(found)

	bidsById := s.keeper.GetBidsByAuctionId(s.ctx, auction.GetId())
	s.Require().Len(bidsById, 4)
	s.Require().Equal(uint64(5), s.keeper.GetNextSequenceWithUpdate(s.ctx, auction.GetId()))

	// create a new auction with auction
	s.SetAuction(types.NewFixedPriceAuction(
		&types.BaseAuction{
			Id:                    3,
			Type:                  types.AuctionTypeFixedPrice,
			Auctioneer:            s.addrs[4].String(),
			SellingReserveAddress: types.SellingReserveAcc(3).String(),
			PayingReserveAddress:  types.PayingReserveAcc(3).String(),
			StartPrice:            sdk.MustNewDecFromStr("0.5"),
			SellingCoin:           sdk.NewInt64Coin(denom3, 1_000_000_000_000),
			PayingCoinDenom:       denom4,
			VestingReserveAddress: types.VestingReserveAcc(3).String(),
			VestingSchedules:      []types.VestingSchedule{},
			WinningPrice:          sdk.ZeroDec(),
			RemainingCoin:         sdk.NewInt64Coin(denom3, 1_000_000_000_000),
			StartTime:             types.MustParseRFC3339("2021-12-10T00:00:00Z"),
			EndTimes:              []time.Time{types.MustParseRFC3339("2022-12-20T00:00:00Z")},
			Status:                types.AuctionStatusStarted,
		},
	))

	auction, found = s.keeper.GetAuction(s.ctx, 3)
	s.Require().True(found)

	// sequence must start with 1
	bidsById = s.keeper.GetBidsByAuctionId(s.ctx, auction.GetId())
	s.Require().Len(bidsById, 0)
	s.Require().Equal(uint64(1), s.keeper.GetNextSequenceWithUpdate(s.ctx, auction.GetId()))
}
