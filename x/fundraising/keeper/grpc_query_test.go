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
	for _, auction := range s.sampleFixedPriceAuctions {
		s.SetAuction(auction)
	}

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
			"query all",
			&types.QueryAuctionsRequest{},
			false,
			func(resp *types.QueryAuctionsResponse) {
				s.Require().Len(resp.Auctions, 2)
			},
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
			"query by type",
			&types.QueryAuctionsRequest{Type: types.AuctionTypeFixedPrice.String()},
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
	for _, auction := range s.sampleFixedPriceAuctions {
		s.SetAuction(auction)
	}

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
			"query by id",
			&types.QueryAuctionRequest{AuctionId: 1},
			false,
			func(resp *types.QueryAuctionResponse) {
				plan, err := types.UnpackAuction(resp.Auction)
				s.Require().NoError(err)
				s.Require().Equal(plan.GetId(), uint64(1))
			},
		},
		{
			"id not found",
			&types.QueryAuctionRequest{AuctionId: 5},
			true,
			nil,
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
	for _, auction := range s.sampleFixedPriceAuctions {
		s.SetAuction(auction)
	}

	for _, bid := range s.sampleFixedPriceBids {
		s.PlaceBid(bid)
	}

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
			&types.QueryBidsRequest{AuctionId: 2},
			false,
			func(resp *types.QueryBidsResponse) {
				s.Require().Len(resp.Bids, 4)
				s.Require().True(coinEq(sdk.NewInt64Coin(denom4, 20_000_000), resp.Bids[0].Coin))
				s.Require().True(coinEq(sdk.NewInt64Coin(denom4, 30_000_000), resp.Bids[1].Coin))
				s.Require().True(coinEq(sdk.NewInt64Coin(denom4, 50_000_000), resp.Bids[2].Coin))
				s.Require().True(coinEq(sdk.NewInt64Coin(denom4, 50_000_000), resp.Bids[3].Coin))
			},
		},
		{
			"query by bidder address",
			&types.QueryBidsRequest{AuctionId: 2, Bidder: s.addrs[0].String()},
			false,
			func(resp *types.QueryBidsResponse) {
				s.Require().Len(resp.Bids, 2)
				for _, bid := range resp.GetBids() {
					s.Require().Equal(s.addrs[0].String(), bid.Bidder)
				}
			},
		},
		{
			"query by eligible",
			&types.QueryBidsRequest{AuctionId: 2, Eligible: "true"},
			false,
			func(resp *types.QueryBidsResponse) {
				s.Require().Len(resp.Bids, 2)
			},
		},
		{
			"query by eligible",
			&types.QueryBidsRequest{AuctionId: 2, Eligible: "false"},
			false,
			func(resp *types.QueryBidsResponse) {
				s.Require().Len(resp.Bids, 2)
			},
		},
		{
			"query by both bidder address and eligible #1",
			&types.QueryBidsRequest{AuctionId: 2, Bidder: s.addrs[0].String(), Eligible: "false"},
			false,
			func(resp *types.QueryBidsResponse) {
				s.Require().Len(resp.Bids, 1)
			},
		},
		{
			"query by both bidder address and eligible #2",
			&types.QueryBidsRequest{AuctionId: 2, Bidder: s.addrs[1].String(), Eligible: "true"},
			false,
			func(resp *types.QueryBidsResponse) {
				s.Require().Len(resp.Bids, 1)
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
	s.SetAuction(s.sampleFixedPriceAuctions[1])

	for _, bid := range s.sampleFixedPriceBids {
		s.PlaceBid(bid)
	}

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
			"query by id and sequence",
			&types.QueryBidRequest{
				AuctionId: 2,
				Sequence:  1,
			},
			false,
			func(resp *types.QueryBidResponse) {
				s.Require().Equal(uint64(2), resp.Bid.GetAuctionId())
				s.Require().Equal(uint64(1), resp.Bid.GetSequence())
			},
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
	s.SetAuction(s.sampleFixedPriceAuctions[1])

	auction, found := s.keeper.GetAuction(s.ctx, 2)
	s.Require().True(found)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	for _, bid := range s.sampleFixedPriceBids {
		s.PlaceBid(bid)
	}

	// set the current block time a day after so that it gets finished
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
				AuctionId: 2,
			},
			false,
			func(resp *types.QueryVestingsResponse) {
				s.Require().Len(resp.Vestings, 4)
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
