package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/fundraising/x/fundraising"
	"github.com/tendermint/fundraising/x/fundraising/types"

	_ "github.com/stretchr/testify/suite"
)

func (s *KeeperTestSuite) TestGRPCParams() {
	resp, err := s.querier.Params(sdk.WrapSDKContext(s.ctx), &types.QueryParamsRequest{})
	s.Require().NoError(err)
	s.Require().Equal(s.keeper.GetParams(s.ctx), resp.Params)
}

func (s *KeeperTestSuite) TestGRPCAuctions() {
	s.createFixedPriceAuction(
		s.addr(0),
		sdk.MustNewDecFromStr("1.0"),
		sdk.NewInt64Coin("denom1", 500_000_000_000),
		"denom2",
		[]types.VestingSchedule{},
		types.MustParseRFC3339("2022-01-01T00:00:00Z"),
		types.MustParseRFC3339("2022-06-10T00:00:00Z"),
		true,
	)

	s.createFixedPriceAuction(
		s.addr(1),
		sdk.MustNewDecFromStr("0.5"),
		sdk.NewInt64Coin("denom3", 1_000_000_000_000),
		"denom4",
		[]types.VestingSchedule{},
		types.MustParseRFC3339("2025-01-01T00:00:00Z"),
		types.MustParseRFC3339("2025-06-10T00:00:00Z"),
		true,
	)

	for _, tc := range []struct {
		name      string
		req       *types.QueryAuctionsRequest
		expectErr bool
		postRun   func(*types.QueryAuctionsResponse)
	}{
		{
			"nil request",
			nil,
			true,
			nil,
		},
		{
			"invalid type",
			&types.QueryAuctionsRequest{
				Type: "invalid",
			},
			true,
			nil,
		},
		{
			"query all auctions",
			&types.QueryAuctionsRequest{},
			false,
			func(resp *types.QueryAuctionsResponse) {
				s.Require().Len(resp.Auctions, 2)
			},
		},
		{
			"query all auctions by auction status",
			&types.QueryAuctionsRequest{
				Status: types.AuctionStatusStandBy.String(),
			},
			false,
			func(resp *types.QueryAuctionsResponse) {
				s.Require().Len(resp.Auctions, 1)
			},
		},
		{
			"query all auction by auction type",
			&types.QueryAuctionsRequest{
				Type: types.AuctionTypeFixedPrice.String(),
			},
			false,
			func(resp *types.QueryAuctionsResponse) {
				auctions, err := types.UnpackAuctions(resp.Auctions)
				s.Require().NoError(err)
				s.Require().Len(auctions, 2)

				for _, auction := range auctions {
					s.Require().Equal(types.AuctionTypeFixedPrice, auction.GetType())
				}
			},
		},
	} {
		s.Run(tc.name, func() {
			resp, err := s.querier.Auctions(sdk.WrapSDKContext(s.ctx), tc.req)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				tc.postRun(resp)
			}
		})
	}
}

func (s *KeeperTestSuite) TestGRPCAuction() {
	auction := s.createFixedPriceAuction(
		s.addr(0),
		sdk.MustNewDecFromStr("0.5"),
		sdk.NewInt64Coin("denom1", 500_000_000_000),
		"denom2",
		[]types.VestingSchedule{},
		types.MustParseRFC3339("2022-01-01T00:00:00Z"),
		types.MustParseRFC3339("2022-03-01T00:00:00Z"),
		true,
	)

	for _, tc := range []struct {
		name      string
		req       *types.QueryAuctionRequest
		expectErr bool
		postRun   func(*types.QueryAuctionResponse)
	}{
		{
			"nil request",
			nil,
			true,
			nil,
		},
		{
			"id not found",
			&types.QueryAuctionRequest{
				AuctionId: 5,
			},
			true,
			nil,
		},
		{
			"query by id",
			&types.QueryAuctionRequest{
				AuctionId: 1,
			},
			false,
			func(resp *types.QueryAuctionResponse) {
				a, err := types.UnpackAuction(resp.Auction)
				s.Require().NoError(err)

				s.Require().Equal(auction.GetId(), a.GetId())
				s.Require().Equal(auction.GetAuctioneer(), a.GetAuctioneer())
				s.Require().Equal(auction.GetType(), a.GetType())
				s.Require().Equal(auction.GetStartPrice(), a.GetStartPrice())
				s.Require().Equal(auction.GetStartTime(), a.GetStartTime())
				s.Require().Equal(auction.GetEndTimes(), a.GetEndTimes())
			},
		},
	} {
		s.Run(tc.name, func() {
			resp, err := s.querier.Auction(sdk.WrapSDKContext(s.ctx), tc.req)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				tc.postRun(resp)
			}
		})
	}
}

func (s *KeeperTestSuite) TestGRPCBids() {
	auction := s.createFixedPriceAuction(
		s.addr(0),
		sdk.MustNewDecFromStr("1.0"),
		sdk.NewInt64Coin("denom1", 500_000_000_000),
		"denom2",
		[]types.VestingSchedule{},
		types.MustParseRFC3339("2022-01-01T00:00:00Z"),
		types.MustParseRFC3339("2022-03-01T00:00:00Z"),
		true,
	)

	bidder1 := s.addr(1)
	bidder2 := s.addr(2)
	bidder3 := s.addr(3)

	bid1 := s.placeBid(auction.GetId(), bidder1, sdk.OneDec(), sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 20_000_000), true)
	bid2 := s.placeBid(auction.GetId(), bidder1, sdk.OneDec(), sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 20_000_000), true)
	bid3 := s.placeBid(auction.GetId(), bidder2, sdk.OneDec(), sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 15_000_000), true)
	bid4 := s.placeBid(auction.GetId(), bidder3, sdk.OneDec(), sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 35_000_000), true)

	// Make bid4 not eligible
	bid4.Eligible = false
	s.keeper.SetBid(s.ctx, auction.GetId(), bid4.Sequence, bidder3, bid4)

	for _, tc := range []struct {
		name      string
		req       *types.QueryBidsRequest
		expectErr bool
		postRun   func(*types.QueryBidsResponse)
	}{
		{
			"nil request",
			nil,
			true,
			nil,
		},
		{
			"query by id",
			&types.QueryBidsRequest{
				AuctionId: 1,
			},
			false,
			func(resp *types.QueryBidsResponse) {
				s.Require().Len(resp.Bids, 4)
				s.Require().True(coinEq(bid1.GetCoin(), resp.Bids[0].Coin))
				s.Require().True(coinEq(bid2.GetCoin(), resp.Bids[1].Coin))
				s.Require().True(coinEq(bid3.GetCoin(), resp.Bids[2].Coin))
				s.Require().True(coinEq(bid4.GetCoin(), resp.Bids[3].Coin))
			},
		},
		{
			"query by bidder address",
			&types.QueryBidsRequest{
				AuctionId: 1,
				Bidder:    bidder1.String(),
			},
			false,
			func(resp *types.QueryBidsResponse) {
				s.Require().Len(resp.Bids, 2)
			},
		},
		{
			"query by eligible",
			&types.QueryBidsRequest{
				AuctionId: 1,
				Eligible:  "true",
			},
			false,
			func(resp *types.QueryBidsResponse) {
				s.Require().Len(resp.Bids, 3)
			},
		},
		{
			"query by eligible",
			&types.QueryBidsRequest{
				AuctionId: 1,
				Eligible:  "false",
			},
			false,
			func(resp *types.QueryBidsResponse) {
				s.Require().Len(resp.Bids, 1)
			},
		},
		{
			"query by both bidder address and eligible #1",
			&types.QueryBidsRequest{
				AuctionId: 1,
				Bidder:    bidder3.String(),
				Eligible:  "false",
			},
			false,
			func(resp *types.QueryBidsResponse) {
				s.Require().Len(resp.Bids, 1)
			},
		},
		{
			"query by both bidder address and eligible #2",
			&types.QueryBidsRequest{
				AuctionId: 1,
				Bidder:    bidder3.String(),
				Eligible:  "true",
			},
			false,
			func(resp *types.QueryBidsResponse) {
				s.Require().Len(resp.Bids, 0)
			},
		},
	} {
		s.Run(tc.name, func() {
			resp, err := s.querier.Bids(sdk.WrapSDKContext(s.ctx), tc.req)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				tc.postRun(resp)
			}
		})
	}
}

func (s *KeeperTestSuite) TestGRPCBid() {
	auction := s.createFixedPriceAuction(
		s.addr(0),
		sdk.MustNewDecFromStr("1.0"),
		sdk.NewInt64Coin("denom1", 500_000_000_000),
		"denom2",
		[]types.VestingSchedule{},
		types.MustParseRFC3339("2022-01-01T00:00:00Z"),
		types.MustParseRFC3339("2022-03-01T00:00:00Z"),
		true,
	)

	bid := s.placeBid(auction.GetId(), s.addr(1), sdk.OneDec(), sdk.NewInt64Coin(auction.GetPayingCoinDenom(), 20_000_000), true)

	for _, tc := range []struct {
		name      string
		req       *types.QueryBidRequest
		expectErr bool
		postRun   func(*types.QueryBidResponse)
	}{
		{
			"nil request",
			nil,
			true,
			nil,
		},
		{
			"id not found",
			&types.QueryBidRequest{
				AuctionId: 5,
			},
			true,
			nil,
		},
		{
			"sequence not found",
			&types.QueryBidRequest{
				AuctionId: 2,
				Sequence:  5,
			},
			true,
			nil,
		},
		{
			"query by id and sequence",
			&types.QueryBidRequest{
				AuctionId: 1,
				Sequence:  1,
			},
			false,
			func(resp *types.QueryBidResponse) {
				s.Require().Equal(bid.GetAuctionId(), resp.Bid.GetAuctionId())
				s.Require().Equal(bid.GetBidder(), resp.Bid.GetBidder())
				s.Require().Equal(bid.GetSequence(), resp.Bid.GetSequence())
				s.Require().Equal(bid.GetCoin(), resp.Bid.GetCoin())
				s.Require().Equal(bid.GetEligible(), resp.Bid.GetEligible())
			},
		},
	} {
		s.Run(tc.name, func() {
			resp, err := s.querier.Bid(sdk.WrapSDKContext(s.ctx), tc.req)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				tc.postRun(resp)
			}
		})
	}
}

func (s *KeeperTestSuite) TestGRPCVestings() {
	auction := s.createFixedPriceAuction(
		s.addr(0),
		sdk.OneDec(),
		sdk.NewInt64Coin("denom1", 1_000_000_000_000),
		"denom2",
		[]types.VestingSchedule{
			{
				ReleaseTime: types.MustParseRFC3339("2022-06-01T00:00:00Z"),
				Weight:      sdk.MustNewDecFromStr("0.5"),
			},
			{
				ReleaseTime: types.MustParseRFC3339("2022-12-01T00:00:00Z"),
				Weight:      sdk.MustNewDecFromStr("0.5"),
			},
		},
		types.MustParseRFC3339("2022-01-01T00:00:00Z"),
		types.MustParseRFC3339("2022-03-01T00:00:00Z"),
		true,
	)

	// Set the current block time a day after so that it gets finished
	s.ctx = s.ctx.WithBlockTime(auction.GetEndTimes()[0].AddDate(0, 0, 1))
	fundraising.EndBlocker(s.ctx, s.keeper)

	for _, tc := range []struct {
		name      string
		req       *types.QueryVestingsRequest
		expectErr bool
		postRun   func(*types.QueryVestingsResponse)
	}{
		{
			"nil request",
			nil,
			true,
			nil,
		},
		{
			"query by id",
			&types.QueryVestingsRequest{
				AuctionId: 1,
			},
			false,
			func(resp *types.QueryVestingsResponse) {
				s.Require().Len(resp.Vestings, 2)
			},
		},
	} {
		s.Run(tc.name, func() {
			resp, err := s.querier.Vestings(sdk.WrapSDKContext(s.ctx), tc.req)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				tc.postRun(resp)
			}
		})
	}
}
