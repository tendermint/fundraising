package keeper_test

import (
	"time"

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
		parseDec("1"),
		parseCoin("5000000000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		time.Now().AddDate(0, 6, 0),
		time.Now().AddDate(0, 6, 0).AddDate(0, 1, 0),
		true,
	)

	s.createFixedPriceAuction(
		s.addr(1),
		sdk.MustNewDecFromStr("0.5"),
		parseCoin("1000000000000denom3"),
		"denom4",
		[]types.VestingSchedule{},
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
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
		parseCoin("500000000000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		time.Now().AddDate(0, 6, 0),
		time.Now().AddDate(0, 6, 0).AddDate(0, 1, 0),
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
				s.Require().Equal(auction.GetStartTime().UTC(), a.GetStartTime())
				s.Require().Equal(auction.GetEndTimes()[0].UTC(), a.GetEndTimes()[0])
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
		parseDec("1"),
		parseCoin("500000000000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 1, 0),
		true,
	)

	bidder1 := s.addr(1)
	bidder2 := s.addr(2)
	bidder3 := s.addr(3)

	s.addAllowedBidder(auction.Id, bidder1, exchangeToSellingAmount(parseDec("1"), parseCoin("400000000denom2")))
	s.addAllowedBidder(auction.Id, bidder2, exchangeToSellingAmount(parseDec("1"), parseCoin("150000000denom2")))
	s.addAllowedBidder(auction.Id, bidder3, exchangeToSellingAmount(parseDec("1"), parseCoin("350000000denom2")))

	bid1 := s.placeBid(auction.GetId(), bidder1, types.BidTypeFixedPrice, parseDec("1"), parseCoin("20000000denom2"), true)
	bid2 := s.placeBid(auction.GetId(), bidder1, types.BidTypeFixedPrice, parseDec("1"), parseCoin("20000000denom2"), true)
	bid3 := s.placeBid(auction.GetId(), bidder2, types.BidTypeFixedPrice, parseDec("1"), parseCoin("15000000denom2"), true)
	bid4 := s.placeBid(auction.GetId(), bidder3, types.BidTypeFixedPrice, parseDec("1"), parseCoin("35000000denom2"), true)

	// Make bid4 not eligible
	bid4.IsWinner = false
	s.keeper.SetBid(s.ctx, bid4)

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
				s.Require().True(coinEq(bid1.Coin, resp.Bids[0].Coin))
				s.Require().True(coinEq(bid2.Coin, resp.Bids[1].Coin))
				s.Require().True(coinEq(bid3.Coin, resp.Bids[2].Coin))
				s.Require().True(coinEq(bid4.Coin, resp.Bids[3].Coin))
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
				IsWinner:  "true",
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
				IsWinner:  "false",
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
				IsWinner:  "false",
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
				IsWinner:  "true",
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
		parseDec("1"),
		parseCoin("500000000000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 1, 0),
		true,
	)

	s.addAllowedBidder(auction.Id, s.addr(1), exchangeToSellingAmount(parseDec("1"), parseCoin("20000000denom2")))
	bid := s.placeBid(auction.GetId(), s.addr(1), types.BidTypeFixedPrice, parseDec("1"), parseCoin("20000000denom2"), true)

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
			"bid id not found",
			&types.QueryBidRequest{
				AuctionId: 2,
				BidId:     5,
			},
			true,
			nil,
		},
		{
			"query by id and bid id",
			&types.QueryBidRequest{
				AuctionId: 1,
				BidId:     1,
			},
			false,
			func(resp *types.QueryBidResponse) {
				s.Require().Equal(bid.AuctionId, resp.Bid.AuctionId)
				s.Require().Equal(bid.GetBidder(), resp.Bid.GetBidder())
				s.Require().Equal(bid.Id, resp.Bid.Id)
				s.Require().Equal(bid.Coin, resp.Bid.Coin)
				s.Require().Equal(bid.IsWinner, resp.Bid.IsWinner)
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
		parseDec("1"),
		parseCoin("1000000000000denom1"),
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
