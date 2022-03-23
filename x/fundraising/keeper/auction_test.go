package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/tendermint/fundraising/x/fundraising"
	"github.com/tendermint/fundraising/x/fundraising/types"

	_ "github.com/stretchr/testify/suite"
)

func (s *KeeperTestSuite) TestFixedPriceAuction_AuctionStatus() {
	standByAuction := s.createFixedPriceAuction(
		s.addr(0),
		parseDec("0.5"),
		parseCoin("5000_000_000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		time.Now().AddDate(0, 6, 0),
		time.Now().AddDate(0, 6, 0).AddDate(0, 1, 0),
		true,
	)

	auction, found := s.keeper.GetAuction(s.ctx, standByAuction.GetId())
	s.Require().True(found)
	s.Require().Equal(types.AuctionStatusStandBy, auction.GetStatus())

	feePool := s.app.DistrKeeper.GetFeePool(s.ctx)
	auctionCreationFee := s.keeper.GetParams(s.ctx).AuctionCreationFee
	s.Require().True(feePool.CommunityPool.IsEqual(sdk.NewDecCoinsFromCoins(auctionCreationFee...)))

	startedAuction := s.createFixedPriceAuction(
		s.addr(1),
		parseDec("0.5"),
		parseCoin("1000_000_000_000denom3"),
		"denom4",
		[]types.VestingSchedule{},
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)

	auction, found = s.keeper.GetAuction(s.ctx, startedAuction.GetId())
	s.Require().True(found)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())
}

func (s *KeeperTestSuite) TestFixedPriceAuction_CalculateAllocation_BidwithPayingCoinDenom() {
	startedAuction := s.createFixedPriceAuction(
		s.addr(0),
		parseDec("0.5"),
		parseCoin("1_000_000_000_000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)

	a, found := s.keeper.GetAuction(s.ctx, startedAuction.GetId())
	s.Require().True(found)

	bidder := s.addr(1)
	s.addAllowedBidder(a.GetId(), bidder, parseInt("1_000_000_000"))

	s.placeBidFixedPrice(a.GetId(), bidder, a.GetStartPrice(), parseCoin("1_000_000denom2"), true)
	s.placeBidFixedPrice(a.GetId(), bidder, a.GetStartPrice(), parseCoin("1_000_000denom2"), true)
	s.placeBidFixedPrice(a.GetId(), bidder, a.GetStartPrice(), parseCoin("1_000_000denom2"), true)

	// Make sure allocate amount is equal to the total bid amount made by the same bidder
	mInfo := s.keeper.CalculateFixedPriceAllocation(s.ctx, a)
	allocateAmt := mInfo.AllocationMap[bidder.String()]
	s.Require().Equal(allocateAmt, parseCoin("6_000_000denom1").Amount)
}

func (s *KeeperTestSuite) TestFixedPriceAuction_CalculateAllocation_BidwithSellingCoinDenom() {
	startedAuction := s.createFixedPriceAuction(
		s.addr(0),
		parseDec("0.5"),
		parseCoin("1_000_000_000_000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)

	a, found := s.keeper.GetAuction(s.ctx, startedAuction.GetId())
	s.Require().True(found)

	bidder := s.addr(1)
	s.addAllowedBidder(a.GetId(), bidder, parseInt("1_000_000_000"))

	s.placeBidFixedPrice(a.GetId(), bidder, a.GetStartPrice(), parseCoin("2_000_000denom1"), true)
	s.placeBidFixedPrice(a.GetId(), bidder, a.GetStartPrice(), parseCoin("2_000_000denom1"), true)
	s.placeBidFixedPrice(a.GetId(), bidder, a.GetStartPrice(), parseCoin("2_000_000denom1"), true)

	// Make sure allocate amount is equal to the total bid amount made by the same bidder
	mInfo := s.keeper.CalculateFixedPriceAllocation(s.ctx, a)
	allocateAmt := mInfo.AllocationMap[bidder.String()]
	s.Require().Equal(allocateAmt, parseCoin("6_000_000denom1").Amount)
}

func (s *KeeperTestSuite) TestFixedPriceAuction_CalculateAllocation_mixed() {
	startedAuction := s.createFixedPriceAuction(
		s.addr(0),
		parseDec("0.5"),
		parseCoin("1_000_000_000_000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)

	a, found := s.keeper.GetAuction(s.ctx, startedAuction.GetId())
	s.Require().True(found)

	bidder := s.addr(1)
	s.addAllowedBidder(a.GetId(), bidder, parseInt("1_000_000_000"))

	s.placeBidFixedPrice(a.GetId(), bidder, a.GetStartPrice(), parseCoin("2_000_000denom1"), true)
	s.placeBidFixedPrice(a.GetId(), bidder, a.GetStartPrice(), parseCoin("2_000_000denom1"), true)
	s.placeBidFixedPrice(a.GetId(), bidder, a.GetStartPrice(), parseCoin("2_000_000denom2"), true)

	// Make sure allocate amount is equal to the total bid amount made by the same bidder
	mInfo := s.keeper.CalculateFixedPriceAllocation(s.ctx, a)
	allocateAmt := mInfo.AllocationMap[bidder.String()]
	s.Require().Equal(allocateAmt, parseCoin("8_000_000denom2").Amount)
}

func (s *KeeperTestSuite) TestFixedPriceAuction_AllocateSellingCoin() {
	auction := s.createFixedPriceAuction(
		s.addr(0),
		parseDec("0.5"),
		parseCoin("1000_000_000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)

	_, found := s.keeper.GetAuction(s.ctx, auction.Id)
	s.Require().True(found)

	// Place bids
	s.placeBidFixedPrice(auction.Id, s.addr(1), parseDec("0.5"), parseCoin("100_000_000denom2"), true)
	s.placeBidFixedPrice(auction.Id, s.addr(2), parseDec("0.5"), parseCoin("100_000_000denom1"), true)
	s.placeBidFixedPrice(auction.Id, s.addr(3), parseDec("0.5"), parseCoin("200_000_000denom1"), true)

	// Calculate allocation
	mInfo := s.keeper.CalculateFixedPriceAllocation(s.ctx, auction)

	// Distribute selling coin
	err := s.keeper.AllocateSellingCoin(s.ctx, auction, mInfo)
	s.Require().NoError(err)

	err = s.keeper.ReleaseRemainingSellingCoin(s.ctx, auction)
	s.Require().NoError(err)

	// The selling reserve account balance must be zero
	s.Require().True(s.getBalance(auction.GetSellingReserveAddress(), auction.SellingCoin.Denom).IsZero())

	// The auctioneer must have sellingCoin.Amount - TotalMatchedAmount
	s.Require().Equal(s.getBalance(s.addr(0), auction.GetSellingCoin().Denom), parseCoin("500_000_000denom1"))
	s.Require().Equal(s.getBalance(auction.GetPayingReserveAddress(), auction.GetPayingCoinDenom()), parseCoin("250_000_000denom2"))

	// The bidders must have the matched selling coin
	s.Require().Equal(s.getBalance(s.addr(1), auction.GetSellingCoin().Denom), parseCoin("200_000_000denom1"))
	s.Require().Equal(s.getBalance(s.addr(2), auction.GetSellingCoin().Denom), parseCoin("100_000_000denom1"))
	s.Require().Equal(s.getBalance(s.addr(3), auction.GetSellingCoin().Denom), parseCoin("200_000_000denom1"))
}

func (s *KeeperTestSuite) TestFixedPriceAuction_AllocateVestingPayingCoin() {
	auction := s.createFixedPriceAuction(
		s.addr(0),
		parseDec("1"),
		parseCoin("1000_000_000_000denom1"),
		"denom2",
		[]types.VestingSchedule{
			{
				ReleaseTime: time.Now().AddDate(0, 6, 0),
				Weight:      sdk.MustNewDecFromStr("0.25"),
			},
			{
				ReleaseTime: time.Now().AddDate(0, 9, 0),
				Weight:      sdk.MustNewDecFromStr("0.25"),
			},
			{
				ReleaseTime: time.Now().AddDate(1, 0, 0),
				Weight:      sdk.MustNewDecFromStr("0.25"),
			},
			{
				ReleaseTime: time.Now().AddDate(1, 3, 0),
				Weight:      sdk.MustNewDecFromStr("0.25"),
			},
		},
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	// Place bids
	s.placeBidFixedPrice(auction.GetId(), s.addr(1), parseDec("1"), parseCoin("100_000_000denom2"), true)
	s.placeBidFixedPrice(auction.GetId(), s.addr(1), parseDec("1"), parseCoin("200_000_000denom2"), true)
	s.placeBidFixedPrice(auction.GetId(), s.addr(1), parseDec("1"), parseCoin("200_000_000denom2"), true)

	// Calculate allocation
	mInfo := s.keeper.CalculateFixedPriceAllocation(s.ctx, auction)

	// Distribute selling coin
	err := s.keeper.AllocateSellingCoin(s.ctx, auction, mInfo)
	s.Require().NoError(err)

	// Apply vesting schedules
	err = s.keeper.ApplyVestingSchedules(s.ctx, auction)
	s.Require().NoError(err)

	// All of the vesting queues must not be released yet
	vqs := s.keeper.GetVestingQueuesByAuctionId(s.ctx, auction.GetId())
	s.Require().Equal(4, len(vqs))
	for _, vq := range vqs {
		s.Require().False(vq.Released)
	}

	// Change the block time to release two vesting schedules
	s.ctx = s.ctx.WithBlockTime(vqs[0].GetReleaseTime().AddDate(0, 4, 1))
	fundraising.EndBlocker(s.ctx, s.keeper)

	// Distribute paying coin
	err = s.keeper.AllocateVestingPayingCoin(s.ctx, auction)
	s.Require().NoError(err)

	// First two vesting queues must be released
	for i, vq := range s.keeper.GetVestingQueuesByAuctionId(s.ctx, auction.GetId()) {
		if i == 0 || i == 1 {
			s.Require().True(vq.Released)
		} else {
			s.Require().False(vq.Released)
		}
	}

	// Change the block time
	s.ctx = s.ctx.WithBlockTime(vqs[3].GetReleaseTime().AddDate(0, 0, 1))
	fundraising.EndBlocker(s.ctx, s.keeper)
	s.Require().NoError(s.keeper.AllocateVestingPayingCoin(s.ctx, auction))

	// All of the vesting queues must be released
	for _, vq := range s.keeper.GetVestingQueuesByAuctionId(s.ctx, auction.GetId()) {
		s.Require().True(vq.Released)
	}

	finishedAuction, found := s.keeper.GetAuction(s.ctx, auction.GetId())
	s.Require().True(found)
	s.Require().Equal(types.AuctionStatusFinished, finishedAuction.GetStatus())
}

func (s *KeeperTestSuite) TestFixedPriceAuction_CancelAuction() {
	auction := s.createFixedPriceAuction(
		s.addr(0),
		parseDec("1"),
		parseCoin("500_000_000_000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 1, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	// Not found auction
	err := s.keeper.CancelAuction(s.ctx, &types.MsgCancelAuction{
		Auctioneer: auction.Auctioneer,
		AuctionId:  10,
	})
	s.Require().Error(err, sdkerrors.ErrNotFound)

	// Unauthorized
	err = s.keeper.CancelAuction(s.ctx, &types.MsgCancelAuction{
		Auctioneer: s.addr(10).String(),
		AuctionId:  auction.Id,
	})
	s.Require().Error(err, sdkerrors.ErrUnauthorized)

	// Invalid auction status
	err = s.keeper.CancelAuction(s.ctx, &types.MsgCancelAuction{
		Auctioneer: auction.Auctioneer,
		AuctionId:  auction.Id,
	})
	s.Require().Error(err, types.ErrInvalidAuctionStatus)

	// Forcefully update auction status
	err = auction.SetStatus(types.AuctionStatusStandBy)
	s.Require().NoError(err)
	s.keeper.SetAuction(s.ctx, auction)

	// Cancel the auction
	err = s.keeper.CancelAuction(s.ctx, &types.MsgCancelAuction{
		Auctioneer: auction.Auctioneer,
		AuctionId:  auction.Id,
	})
	s.Require().NoError(err)

	// Verify the status
	a, found := s.keeper.GetAuction(s.ctx, auction.GetId())
	s.Require().True(found)
	s.Require().Equal(types.AuctionStatusCancelled, a.GetStatus())

	// The selling reserve balance must be zero
	sellingReserveAddr := a.GetSellingReserveAddress()
	sellingCoinDenom := a.GetSellingCoin().Denom
	s.Require().True(s.getBalance(sellingReserveAddr, sellingCoinDenom).IsZero())
}

func (s *KeeperTestSuite) TestBatchAuction_AuctionStatus() {
	standByAuction := s.createBatchAuction(
		s.addr(0),
		parseDec("1"),
		parseDec("0.1"),
		parseCoin("5000_000_000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		1,
		sdk.MustNewDecFromStr("0.2"),
		time.Now().AddDate(0, 6, 0),
		time.Now().AddDate(0, 6, 0).AddDate(0, 1, 0),
		true,
	)

	auction, found := s.keeper.GetAuction(s.ctx, standByAuction.GetId())
	s.Require().True(found)
	s.Require().Equal(types.AuctionStatusStandBy, auction.GetStatus())

	feePool := s.app.DistrKeeper.GetFeePool(s.ctx)
	auctionCreationFee := s.keeper.GetParams(s.ctx).AuctionCreationFee
	s.Require().True(feePool.CommunityPool.IsEqual(sdk.NewDecCoinsFromCoins(auctionCreationFee...)))

	startedAuction := s.createBatchAuction(
		s.addr(1),
		parseDec("0.5"),
		parseDec("0.1"),
		parseCoin("5000_000_000denom3"),
		"denom4",
		[]types.VestingSchedule{},
		1,
		sdk.MustNewDecFromStr("0.2"),
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)

	auction, found = s.keeper.GetAuction(s.ctx, startedAuction.GetId())
	s.Require().True(found)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())
}

func (s *KeeperTestSuite) TestAddAllowedBidders() {
	startedAuction := s.createFixedPriceAuction(
		s.addr(0),
		parseDec("0.5"),
		parseCoin("500_000_000_000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)

	auction, found := s.keeper.GetAuction(s.ctx, startedAuction.GetId())
	s.Require().True(found)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())
	s.Require().Len(auction.GetAllowedBidders(), 0)

	for _, tc := range []struct {
		name        string
		bidders     []types.AllowedBidder
		expectedErr error
	}{
		{
			"single bidder",
			[]types.AllowedBidder{
				{Bidder: s.addr(1).String(), MaxBidAmount: sdk.NewInt(100_000_000)},
			},
			nil,
		},
		{
			"multiple bidders",
			[]types.AllowedBidder{
				{Bidder: s.addr(1).String(), MaxBidAmount: sdk.NewInt(100_000_000)},
				{Bidder: s.addr(2).String(), MaxBidAmount: sdk.NewInt(500_000_000)},
				{Bidder: s.addr(3).String(), MaxBidAmount: sdk.NewInt(800_000_000)},
			},
			nil,
		},
		{

			"empty bidders",
			[]types.AllowedBidder{},
			types.ErrEmptyAllowedBidders,
		},
		{
			"zero maximum bid amount",
			[]types.AllowedBidder{
				{Bidder: s.addr(1).String(), MaxBidAmount: sdk.NewInt(0)},
			},
			types.ErrInvalidMaxBidAmount,
		},
		{
			"negative maximum bid amount",
			[]types.AllowedBidder{
				{Bidder: s.addr(1).String(), MaxBidAmount: sdk.NewInt(-1)},
			},
			types.ErrInvalidMaxBidAmount,
		},
		{
			"exceed the total selling amount",
			[]types.AllowedBidder{
				{Bidder: s.addr(1).String(), MaxBidAmount: sdk.NewInt(500_000_000_001)},
			},
			types.ErrInsufficientRemainingAmount,
		},
	} {
		s.Run(tc.name, func() {
			err := s.keeper.AddAllowedBidders(s.ctx, auction.GetId(), tc.bidders)
			if tc.expectedErr != nil {
				s.Require().ErrorIs(err, tc.expectedErr)
				return
			}
			s.Require().NoError(err)
		})
	}
}

func (s *KeeperTestSuite) TestAddAllowedBidders_Length() {
	startedAuction := s.createFixedPriceAuction(
		s.addr(0),
		parseDec("0.5"),
		parseCoin("500_000_000_000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)

	auction, found := s.keeper.GetAuction(s.ctx, startedAuction.GetId())
	s.Require().True(found)
	s.Require().Len(auction.GetAllowedBidders(), 0)

	// Add some bidders
	s.Require().NoError(s.keeper.AddAllowedBidders(s.ctx, auction.GetId(), []types.AllowedBidder{
		{Bidder: s.addr(1).String(), MaxBidAmount: sdk.NewInt(100_000_000)},
		{Bidder: s.addr(2).String(), MaxBidAmount: sdk.NewInt(500_000_000)},
	}))

	auction, found = s.keeper.GetAuction(s.ctx, auction.GetId())
	s.Require().True(found)
	s.Require().Len(auction.GetAllowedBidders(), 2)

	// Add more bidders
	s.Require().NoError(s.keeper.AddAllowedBidders(s.ctx, auction.GetId(), []types.AllowedBidder{
		{Bidder: s.addr(3).String(), MaxBidAmount: sdk.NewInt(100_000_000)},
		{Bidder: s.addr(4).String(), MaxBidAmount: sdk.NewInt(100_000_000)},
		{Bidder: s.addr(5).String(), MaxBidAmount: sdk.NewInt(100_000_000)},
	}))

	auction, found = s.keeper.GetAuction(s.ctx, auction.GetId())
	s.Require().True(found)
	s.Require().Len(auction.GetAllowedBidders(), 5)
}

func (s *KeeperTestSuite) TestUpdateAllowedBidder() {
	startedAuction := s.createFixedPriceAuction(
		s.addr(0),
		parseDec("0.5"),
		parseCoin("500_000_000_000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)

	auction, found := s.keeper.GetAuction(s.ctx, startedAuction.GetId())
	s.Require().True(found)
	s.Require().Len(auction.GetAllowedBidders(), 0)

	// Add 5 bidders with different maximum bid amount
	s.Require().NoError(s.keeper.AddAllowedBidders(s.ctx, auction.GetId(), []types.AllowedBidder{
		{Bidder: s.addr(1).String(), MaxBidAmount: sdk.NewInt(100_000_000)},
		{Bidder: s.addr(2).String(), MaxBidAmount: sdk.NewInt(200_000_000)},
		{Bidder: s.addr(3).String(), MaxBidAmount: sdk.NewInt(300_000_000)},
		{Bidder: s.addr(4).String(), MaxBidAmount: sdk.NewInt(400_000_000)},
		{Bidder: s.addr(5).String(), MaxBidAmount: sdk.NewInt(500_000_000)},
	}))

	auction, found = s.keeper.GetAuction(s.ctx, startedAuction.GetId())
	s.Require().True(found)
	s.Require().Len(auction.GetAllowedBidders(), 5)

	for _, tc := range []struct {
		name         string
		bidder       sdk.AccAddress
		maxBidAmount sdk.Int
		expectedErr  error
	}{
		{
			"update bidder's maximum bid amount",
			s.addr(1),
			sdk.NewInt(555_000_000_000),
			nil,
		},
		{
			"bidder not found",
			s.addr(10),
			sdk.NewInt(300_000_000),
			sdkerrors.Wrapf(sdkerrors.ErrNotFound, "bidder %s is not found", s.addr(10).String()),
		},
		{
			"zero maximum bid amount value",
			s.addr(1),
			sdk.NewInt(0),
			types.ErrInvalidMaxBidAmount,
		},
		{
			"negative maximum bid amount value",
			s.addr(1),
			sdk.NewInt(-1),
			types.ErrInvalidMaxBidAmount,
		},
	} {
		s.Run(tc.name, func() {
			err := s.keeper.UpdateAllowedBidder(s.ctx, auction.GetId(), tc.bidder, tc.maxBidAmount)
			if tc.expectedErr != nil {
				s.Require().ErrorIs(err, tc.expectedErr)
				return
			}
			s.Require().NoError(err)

			auction, found = s.keeper.GetAuction(s.ctx, auction.GetId())
			s.Require().True(found)
			s.Require().Len(auction.GetAllowedBidders(), 5)

			// Check if it is sucessfully updated
			allowedBiddersMap := auction.GetAllowedBiddersMap()
			maxBidAmt := allowedBiddersMap[tc.bidder.String()]
			s.Require().Equal(tc.maxBidAmount, maxBidAmt)
		})
	}
}
