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
		parseDec("0.1"),
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
		parseDec("0.1"),
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

func (s *KeeperTestSuite) TestAllocateSellingCoin_FixedPriceAuction() {
	auction := s.createFixedPriceAuction(
		s.addr(0),
		parseDec("1"),
		parseDec("0.1"),
		parseCoin("1000_000_000_000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)

	_, found := s.keeper.GetAuction(s.ctx, auction.Id)
	s.Require().True(found)

	// Place bids
	s.placeBidFixedPrice(auction.Id, s.addr(1), parseDec("1"), parseCoin("100_000_000denom2"), true)
	s.placeBidFixedPrice(auction.Id, s.addr(2), parseDec("1"), parseCoin("200_000_000denom2"), true)
	s.placeBidFixedPrice(auction.Id, s.addr(3), parseDec("1"), parseCoin("200_000_000denom2"), true)

	// Calculate allocation
	mInfo := s.keeper.CalculateFixedPriceAllocation(s.ctx, auction)

	// Distribute selling coin
	err := s.keeper.AllocateSellingCoin(s.ctx, auction, mInfo)
	s.Require().NoError(err)

	err = s.keeper.ReleaseRemainingSellingCoin(s.ctx, auction)
	s.Require().NoError(err)

	// The selling reserve account balance must be zero
	s.Require().True(s.getBalance(auction.GetSellingReserveAddress(), auction.SellingCoin.Denom).IsZero())

	// The bidders must have the selling coin
	s.Require().False(s.getBalance(s.addr(1), auction.GetSellingCoin().Denom).IsZero())
	s.Require().False(s.getBalance(s.addr(2), auction.GetSellingCoin().Denom).IsZero())
	s.Require().False(s.getBalance(s.addr(3), auction.GetSellingCoin().Denom).IsZero())
}

func (s *KeeperTestSuite) TestAllocateVestingPayingCoin() {
	auction := s.createFixedPriceAuction(
		s.addr(0),
		parseDec("1"),
		parseDec("0.1"),
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

func (s *KeeperTestSuite) TestCancelAuction() {
	standByAuction := s.createFixedPriceAuction(
		s.addr(0),
		parseDec("1"),
		parseDec("0.1"),
		parseCoin("500_000_000_000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		time.Now().AddDate(0, 1, 0),
		time.Now().AddDate(0, 1, 0).AddDate(0, 1, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStandBy, standByAuction.GetStatus())

	// Cancel the auction
	auction := s.cancelAuction(standByAuction.GetId(), s.addr(0))
	s.Require().Equal(types.AuctionStatusCancelled, auction.GetStatus())

	// The selling reserve balance must be zero
	sellingReserveAddr := auction.GetSellingReserveAddress()
	sellingCoinDenom := auction.GetSellingCoin().Denom
	s.Require().True(s.getBalance(sellingReserveAddr, sellingCoinDenom).IsZero())
}

func (s *KeeperTestSuite) TestAddAllowedBidders() {
	startedAuction := s.createFixedPriceAuction(
		s.addr(0),
		parseDec("0.5"),
		parseDec("0.1"),
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
		parseDec("0.1"),
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
		parseDec("0.1"),
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

// Example of "JH_ex0" in Sheet
func (s *KeeperTestSuite) TestCalculateAllocation_Many() {
	auction := s.createBatchAuction(
		s.addr(0),
		parseDec("1"),
		parseDec("0.1"),
		parseCoin("1000_000_000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		1,
		sdk.MustNewDecFromStr("0.2"),
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	s.placeBidBatchMany(auction.Id, s.addr(1), parseDec("1"), parseCoin("500_000_000denom1"), sdk.NewInt(1000_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(2), parseDec("0.9"), parseCoin("500_000_000denom1"), sdk.NewInt(1000_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(3), parseDec("0.8"), parseCoin("500_000_000denom1"), sdk.NewInt(1000_000_000), true)

	a, found := s.keeper.GetAuction(s.ctx, auction.Id)
	s.Require().True(found)

	mInfo := s.keeper.CalculateBatchAllocation(s.ctx, a)

	// Checking
	s.Require().Equal(mInfo.MatchedLen, int64(2))
	s.Require().Equal(mInfo.MatchedPrice, parseDec("0.9"))
	s.Require().Equal(mInfo.TotalMatchedAmount, sdk.NewInt(1000_000_000))
	s.Require().Equal(mInfo.AllocationMap[s.addr(1).String()], sdk.NewInt(500_000_000))
	s.Require().Equal(mInfo.AllocationMap[s.addr(2).String()], sdk.NewInt(500_000_000))
	s.Require().Equal(mInfo.AllocationMap[s.addr(3).String()], sdk.NewInt(0))
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(1).String()], sdk.NewInt(450_000_000))
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(2).String()], sdk.NewInt(450_000_000))
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(3).String()], sdk.NewInt(0))
	s.Require().Equal(mInfo.RefundMap[s.addr(1).String()], sdk.NewInt(50_000_000))
	s.Require().Equal(mInfo.RefundMap[s.addr(2).String()].Abs(), sdk.NewInt(0).Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(3).String()], sdk.NewInt(400_000_000))

	// Distribute selling coin
	err := s.keeper.AllocateSellingCoin(s.ctx, auction, mInfo)
	s.Require().NoError(err)

	s.Require().Equal(s.getBalance(auction.GetSellingReserveAddress(), auction.SellingCoin.Denom).Amount.Abs(), auction.SellingCoin.Amount.Sub(mInfo.TotalMatchedAmount).Abs())

	err = s.keeper.ReleaseRemainingSellingCoin(s.ctx, auction)
	s.Require().NoError(err)

	// The selling reserve account balance must be zero
	s.Require().True(s.getBalance(auction.GetSellingReserveAddress(), auction.SellingCoin.Denom).IsZero())

	// The auctioneer must have sellingCoin.Amount - TotalMatchedAmount
	s.Require().Equal(s.getBalance(s.addr(0), auction.GetSellingCoin().Denom).Amount, auction.SellingCoin.Amount.Sub(mInfo.TotalMatchedAmount).Abs())

	// The bidders must have the matched selling coin
	s.Require().Equal(s.getBalance(s.addr(1), auction.GetSellingCoin().Denom).Amount, sdk.NewInt(500_000_000))
	s.Require().Equal(s.getBalance(s.addr(2), auction.GetSellingCoin().Denom).Amount, sdk.NewInt(500_000_000))
	s.Require().True(s.getBalance(s.addr(3), auction.GetSellingCoin().Denom).IsZero())

	// Refund payingCoin
	err = s.keeper.RefundPayingCoin(s.ctx, auction, mInfo)
	s.Require().NoError(err)
}

// Example of "JH_ex0.1" in Sheet
func (s *KeeperTestSuite) TestCalculateAllocation_Worth() {
	auction := s.createBatchAuction(
		s.addr(0),
		parseDec("1"),
		parseDec("0.1"),
		parseCoin("1500_000_000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		1,
		sdk.MustNewDecFromStr("0.2"),
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	s.placeBidBatchWorth(auction.Id, s.addr(1), parseDec("1"), parseCoin("500_000_000denom2"), sdk.NewInt(1500_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(2), parseDec("0.9"), parseCoin("500_000_000denom2"), sdk.NewInt(1500_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(3), parseDec("0.8"), parseCoin("500_000_000denom2"), sdk.NewInt(1500_000_000), true)

	a, found := s.keeper.GetAuction(s.ctx, auction.Id)
	s.Require().True(found)

	mInfo := s.keeper.CalculateBatchAllocation(s.ctx, a)

	// Checking
	s.Require().Equal(mInfo.MatchedLen, int64(2))
	s.Require().Equal(mInfo.MatchedPrice, parseDec("0.9"))
	matchingPrice := parseDec("0.9")
	MatchedAmt := sdk.NewInt(500_000_000).ToDec().QuoTruncate(matchingPrice).TruncateInt()

	s.Require().Equal(mInfo.TotalMatchedAmount, MatchedAmt.Add(MatchedAmt))
	s.Require().Equal(mInfo.AllocationMap[s.addr(1).String()], MatchedAmt)
	s.Require().Equal(mInfo.AllocationMap[s.addr(2).String()], MatchedAmt)
	s.Require().Equal(mInfo.AllocationMap[s.addr(3).String()], sdk.NewInt(0))
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(1).String()], sdk.NewInt(500_000_000))
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(2).String()], sdk.NewInt(500_000_000))
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(3).String()], sdk.NewInt(0))
	s.Require().Equal(mInfo.RefundMap[s.addr(1).String()].Abs(), sdk.NewInt(0).Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(2).String()].Abs(), sdk.NewInt(0).Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(3).String()], sdk.NewInt(500_000_000))

	// Distribute selling coin
	err := s.keeper.AllocateSellingCoin(s.ctx, auction, mInfo)
	s.Require().NoError(err)

	err = s.keeper.ReleaseRemainingSellingCoin(s.ctx, auction)
	s.Require().NoError(err)

	// The selling reserve account balance must be zero
	s.Require().True(s.getBalance(auction.GetSellingReserveAddress(), auction.SellingCoin.Denom).IsZero())

	// The auctioneer must have sellingCoin.Amount - TotalMatchedAmount
	s.Require().Equal(s.getBalance(s.addr(0), auction.GetSellingCoin().Denom).Amount, auction.SellingCoin.Amount.Sub(mInfo.TotalMatchedAmount).Abs())

	// The bidders must have the matched selling coin
	s.Require().Equal(s.getBalance(s.addr(1), auction.GetSellingCoin().Denom).Amount, MatchedAmt)
	s.Require().Equal(s.getBalance(s.addr(2), auction.GetSellingCoin().Denom).Amount, MatchedAmt)
	s.Require().True(s.getBalance(s.addr(3), auction.GetSellingCoin().Denom).IsZero())

	// Refund payingCoin
	err = s.keeper.RefundPayingCoin(s.ctx, auction, mInfo)
	s.Require().NoError(err)
}

// Example of "JH_ex0.2" in Sheet
func (s *KeeperTestSuite) TestCalculateAllocation_Mixed() {
	auction := s.createBatchAuction(
		s.addr(0),
		parseDec("1"),
		parseDec("0.1"),
		parseCoin("1700_000_000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		1,
		sdk.MustNewDecFromStr("0.2"),
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	s.placeBidBatchMany(auction.Id, s.addr(1), parseDec("1"), parseCoin("500_000_000denom1"), sdk.NewInt(1500_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(2), parseDec("0.9"), parseCoin("500_000_000denom2"), sdk.NewInt(1500_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(3), parseDec("0.8"), parseCoin("500_000_000denom2"), sdk.NewInt(1500_000_000), true)

	a, found := s.keeper.GetAuction(s.ctx, auction.Id)
	s.Require().True(found)

	mInfo := s.keeper.CalculateBatchAllocation(s.ctx, a)

	// Checking
	s.Require().Equal(mInfo.MatchedLen, int64(2))
	s.Require().Equal(mInfo.MatchedPrice, parseDec("0.9"))
	matchingPrice := parseDec("0.9")
	MatchedAmt1 := sdk.NewInt(500_000_000)
	MatchedAmt2 := sdk.NewInt(500_000_000).ToDec().QuoTruncate(matchingPrice).TruncateInt()

	s.Require().Equal(mInfo.TotalMatchedAmount, sdk.NewInt(500_000_000).Add(MatchedAmt2))
	s.Require().Equal(mInfo.AllocationMap[s.addr(1).String()], MatchedAmt1)
	s.Require().Equal(mInfo.AllocationMap[s.addr(2).String()], MatchedAmt2)
	s.Require().Equal(mInfo.AllocationMap[s.addr(3).String()], sdk.NewInt(0))
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(1).String()], sdk.NewInt(450_000_000))
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(2).String()], sdk.NewInt(500_000_000))
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(3).String()], sdk.NewInt(0))
	s.Require().Equal(mInfo.RefundMap[s.addr(1).String()], sdk.NewInt(50_000_000))
	s.Require().Equal(mInfo.RefundMap[s.addr(2).String()].Abs(), sdk.NewInt(0).Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(3).String()], sdk.NewInt(500_000_000))

	// Distribute selling coin
	err := s.keeper.AllocateSellingCoin(s.ctx, auction, mInfo)
	s.Require().NoError(err)

	err = s.keeper.ReleaseRemainingSellingCoin(s.ctx, auction)
	s.Require().NoError(err)

	// The selling reserve account balance must be zero
	s.Require().True(s.getBalance(auction.GetSellingReserveAddress(), auction.SellingCoin.Denom).IsZero())

	// The auctioneer must have sellingCoin.Amount - TotalMatchedAmount
	s.Require().Equal(s.getBalance(s.addr(0), auction.GetSellingCoin().Denom).Amount, auction.SellingCoin.Amount.Sub(mInfo.TotalMatchedAmount).Abs())

	// The bidders must have the matched selling coin
	s.Require().Equal(s.getBalance(s.addr(1), auction.GetSellingCoin().Denom).Amount, MatchedAmt1)
	s.Require().Equal(s.getBalance(s.addr(2), auction.GetSellingCoin().Denom).Amount, MatchedAmt2)
	s.Require().True(s.getBalance(s.addr(3), auction.GetSellingCoin().Denom).IsZero())

	// Refund payingCoin
	err = s.keeper.RefundPayingCoin(s.ctx, auction, mInfo)
	s.Require().NoError(err)
}

// Example of "JH_ex0.01" in Sheet for MaxBidAmountLimit
func (s *KeeperTestSuite) TestCalculateAllocation_Many_Limited() {
	auction := s.createBatchAuction(
		s.addr(0),
		parseDec("1"),
		parseDec("0.1"),
		parseCoin("1000_000_000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		1,
		sdk.MustNewDecFromStr("0.2"),
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	s.placeBidBatchMany(auction.Id, s.addr(1), parseDec("1"), parseCoin("400_000_000denom1"), sdk.NewInt(400_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(2), parseDec("0.9"), parseCoin("400_000_000denom1"), sdk.NewInt(400_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(3), parseDec("0.8"), parseCoin("400_000_000denom1"), sdk.NewInt(400_000_000), true)

	a, found := s.keeper.GetAuction(s.ctx, auction.Id)
	s.Require().True(found)

	mInfo := s.keeper.CalculateBatchAllocation(s.ctx, a)

	// Checking
	s.Require().Equal(mInfo.MatchedLen, int64(2))
	s.Require().Equal(mInfo.MatchedPrice, parseDec("0.9"))
	s.Require().Equal(mInfo.TotalMatchedAmount, sdk.NewInt(800_000_000))
	s.Require().Equal(mInfo.AllocationMap[s.addr(1).String()], sdk.NewInt(400_000_000))
	s.Require().Equal(mInfo.AllocationMap[s.addr(2).String()], sdk.NewInt(400_000_000))
	s.Require().Equal(mInfo.AllocationMap[s.addr(3).String()], sdk.NewInt(0))
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(1).String()], sdk.NewInt(360_000_000))
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(2).String()], sdk.NewInt(360_000_000))
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(3).String()], sdk.NewInt(0))
	s.Require().Equal(mInfo.RefundMap[s.addr(1).String()], sdk.NewInt(40_000_000))
	s.Require().Equal(mInfo.RefundMap[s.addr(2).String()].Abs(), sdk.NewInt(0).Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(3).String()], sdk.NewInt(320_000_000))

	// Distribute selling coin
	err := s.keeper.AllocateSellingCoin(s.ctx, auction, mInfo)
	s.Require().NoError(err)

	err = s.keeper.ReleaseRemainingSellingCoin(s.ctx, auction)
	s.Require().NoError(err)

	// The selling reserve account balance must be zero
	s.Require().True(s.getBalance(auction.GetSellingReserveAddress(), auction.SellingCoin.Denom).IsZero())

	// The auctioneer must have sellingCoin.Amount - TotalMatchedAmount
	s.Require().Equal(s.getBalance(s.addr(0), auction.GetSellingCoin().Denom).Amount, sdk.NewInt(200_000_000))

	// The bidders must have the matched selling coin
	s.Require().Equal(s.getBalance(s.addr(1), auction.GetSellingCoin().Denom).Amount, sdk.NewInt(400_000_000))
	s.Require().Equal(s.getBalance(s.addr(2), auction.GetSellingCoin().Denom).Amount, sdk.NewInt(400_000_000))
	s.Require().True(s.getBalance(s.addr(3), auction.GetSellingCoin().Denom).IsZero())

	// Refund payingCoin
	err = s.keeper.RefundPayingCoin(s.ctx, auction, mInfo)
	s.Require().NoError(err)
}

// Example of "JH_ex0.11" in Sheet for MaxBidAmountLimit
func (s *KeeperTestSuite) TestCalculateAllocation_Worth_Limited() {
	auction := s.createBatchAuction(
		s.addr(0),
		parseDec("1"),
		parseDec("0.1"),
		parseCoin("1500_000_000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		1,
		sdk.MustNewDecFromStr("0.2"),
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	s.placeBidBatchWorth(auction.Id, s.addr(1), parseDec("1"), parseCoin("400_000_000denom2"), sdk.NewInt(400_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(2), parseDec("0.9"), parseCoin("360_000_000denom2"), sdk.NewInt(400_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(3), parseDec("0.8"), parseCoin("320_000_000denom2"), sdk.NewInt(400_000_000), true)

	a, found := s.keeper.GetAuction(s.ctx, auction.Id)
	s.Require().True(found)

	mInfo := s.keeper.CalculateBatchAllocation(s.ctx, a)

	// Checking
	s.Require().Equal(mInfo.MatchedLen, int64(3))
	s.Require().Equal(mInfo.MatchedPrice, parseDec("0.8"))
	s.Require().Equal(mInfo.TotalMatchedAmount, sdk.NewInt(1200_000_000))
	s.Require().Equal(mInfo.AllocationMap[s.addr(1).String()], sdk.NewInt(400_000_000))
	s.Require().Equal(mInfo.AllocationMap[s.addr(2).String()], sdk.NewInt(400_000_000))
	s.Require().Equal(mInfo.AllocationMap[s.addr(3).String()], sdk.NewInt(400_000_000))
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(1).String()], sdk.NewInt(320_000_000))
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(2).String()], sdk.NewInt(320_000_000))
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(3).String()], sdk.NewInt(320_000_000))
	s.Require().Equal(mInfo.RefundMap[s.addr(1).String()], sdk.NewInt(80_000_000))
	s.Require().Equal(mInfo.RefundMap[s.addr(2).String()], sdk.NewInt(40_000_000))
	s.Require().Equal(mInfo.RefundMap[s.addr(3).String()].Abs(), sdk.NewInt(0_000_000).Abs())

	// Distribute selling coin
	err := s.keeper.AllocateSellingCoin(s.ctx, auction, mInfo)
	s.Require().NoError(err)

	err = s.keeper.ReleaseRemainingSellingCoin(s.ctx, auction)
	s.Require().NoError(err)

	// The selling reserve account balance must be zero
	s.Require().True(s.getBalance(auction.GetSellingReserveAddress(), auction.SellingCoin.Denom).IsZero())

	// The auctioneer must have sellingCoin.Amount - TotalMatchedAmount
	s.Require().Equal(s.getBalance(s.addr(0), auction.GetSellingCoin().Denom).Amount, sdk.NewInt(300_000_000))

	// The bidders must have the matched selling coin
	s.Require().Equal(s.getBalance(s.addr(1), auction.GetSellingCoin().Denom).Amount, sdk.NewInt(400_000_000))
	s.Require().Equal(s.getBalance(s.addr(2), auction.GetSellingCoin().Denom).Amount, sdk.NewInt(400_000_000))
	s.Require().Equal(s.getBalance(s.addr(3), auction.GetSellingCoin().Denom).Amount, sdk.NewInt(400_000_000))

	// Refund payingCoin
	err = s.keeper.RefundPayingCoin(s.ctx, auction, mInfo)
	s.Require().NoError(err)
}

// Example of "JH_ex0.2" in Sheet for MaxBidAmountLimit
func (s *KeeperTestSuite) TestCalculateAllocation_Mixed_Limited() {
	auction := s.createBatchAuction(
		s.addr(0),
		parseDec("1"),
		parseDec("0.1"),
		parseCoin("1700_000_000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		1,
		sdk.MustNewDecFromStr("0.2"),
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	s.placeBidBatchMany(auction.Id, s.addr(1), parseDec("1"), parseCoin("500_000_000denom1"), sdk.NewInt(600_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(2), parseDec("0.9"), parseCoin("500_000_000denom2"), sdk.NewInt(600_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(3), parseDec("0.8"), parseCoin("450_000_000denom2"), sdk.NewInt(600_000_000), true)

	a, found := s.keeper.GetAuction(s.ctx, auction.Id)
	s.Require().True(found)

	mInfo := s.keeper.CalculateBatchAllocation(s.ctx, a)

	// Checking
	s.Require().Equal(mInfo.MatchedLen, int64(3))
	s.Require().Equal(mInfo.MatchedPrice, parseDec("0.8"))
	s.Require().Equal(mInfo.TotalMatchedAmount, sdk.NewInt(1662_500_000))
	s.Require().Equal(mInfo.AllocationMap[s.addr(1).String()], sdk.NewInt(500_000_000))
	s.Require().Equal(mInfo.AllocationMap[s.addr(2).String()], sdk.NewInt(600_000_000))
	s.Require().Equal(mInfo.AllocationMap[s.addr(3).String()], sdk.NewInt(562_500_000))
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(1).String()], sdk.NewInt(400_000_000))
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(2).String()], sdk.NewInt(480_000_000))
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(3).String()], sdk.NewInt(450_000_000))
	s.Require().Equal(mInfo.RefundMap[s.addr(1).String()], sdk.NewInt(100_000_000))
	s.Require().Equal(mInfo.RefundMap[s.addr(2).String()], sdk.NewInt(20_000_000))
	s.Require().Equal(mInfo.RefundMap[s.addr(3).String()].Abs(), sdk.NewInt(0).Abs())

	// Distribute selling coin
	err := s.keeper.AllocateSellingCoin(s.ctx, auction, mInfo)
	s.Require().NoError(err)

	err = s.keeper.ReleaseRemainingSellingCoin(s.ctx, auction)
	s.Require().NoError(err)

	// The selling reserve account balance must be zero
	s.Require().True(s.getBalance(auction.GetSellingReserveAddress(), auction.SellingCoin.Denom).IsZero())

	// The auctioneer must have sellingCoin.Amount - TotalMatchedAmount
	s.Require().Equal(s.getBalance(s.addr(0), auction.GetSellingCoin().Denom).Amount.Abs(), sdk.NewInt(37_500_000).Abs())

	// The bidders must have the matched selling coin
	s.Require().Equal(s.getBalance(s.addr(1), auction.GetSellingCoin().Denom).Amount, sdk.NewInt(500_000_000))
	s.Require().Equal(s.getBalance(s.addr(2), auction.GetSellingCoin().Denom).Amount, sdk.NewInt(600_000_000))
	s.Require().Equal(s.getBalance(s.addr(3), auction.GetSellingCoin().Denom).Amount, sdk.NewInt(562_500_000))

	// Refund payingCoin
	err = s.keeper.RefundPayingCoin(s.ctx, auction, mInfo)
	s.Require().NoError(err)
}

// Example of "JH_ex1" in Sheet
func (s *KeeperTestSuite) TestCalculateAllocation_Mixed2() {
	auction := s.createBatchAuction(
		s.addr(0),
		parseDec("1"),
		parseDec("0.1"),
		parseCoin("5000_000_000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		1,
		sdk.MustNewDecFromStr("0.2"),
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	s.placeBidBatchMany(auction.Id, s.addr(1), parseDec("1"), parseCoin("200_000_000denom1"), sdk.NewInt(5000_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(2), parseDec("0.8"), parseCoin("500_000_000denom2"), sdk.NewInt(5000_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(3), parseDec("0.9"), parseCoin("500_000_000denom1"), sdk.NewInt(5000_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(1), parseDec("1.1"), parseCoin("300_000_000denom2"), sdk.NewInt(0), true)
	s.placeBidBatchMany(auction.Id, s.addr(5), parseDec("1.2"), parseCoin("300_000_000denom1"), sdk.NewInt(5000_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(4), parseDec("0.8"), parseCoin("100_000_000denom1"), sdk.NewInt(5000_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(2), parseDec("0.7"), parseCoin("100_000_000denom1"), sdk.NewInt(0), true)
	s.placeBidBatchMany(auction.Id, s.addr(6), parseDec("0.5"), parseCoin("100_000_000denom1"), sdk.NewInt(5000_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(3), parseDec("0.8"), parseCoin("300_000_000denom2"), sdk.NewInt(0), true)
	s.placeBidBatchWorth(auction.Id, s.addr(7), parseDec("0.6"), parseCoin("500_000_000denom2"), sdk.NewInt(5000_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(8), parseDec("0.8"), parseCoin("500_000_000denom1"), sdk.NewInt(5000_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(9), parseDec("0.6"), parseCoin("600_000_000denom1"), sdk.NewInt(5000_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(6), parseDec("0.5"), parseCoin("500_000_000denom2"), sdk.NewInt(0), true)
	s.placeBidBatchMany(auction.Id, s.addr(10), parseDec("0.6"), parseCoin("100_000_000denom1"), sdk.NewInt(5000_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(3), parseDec("0.7"), parseCoin("800_000_000denom2"), sdk.NewInt(0), true)

	a, found := s.keeper.GetAuction(s.ctx, auction.Id)
	s.Require().True(found)

	mInfo := s.keeper.CalculateBatchAllocation(s.ctx, a)

	// Checking
	s.Require().Equal(mInfo.MatchedLen, int64(10))
	matchingPrice := parseDec("0.7")
	s.Require().Equal(mInfo.MatchedPrice, matchingPrice)

	MatchedAmt1 := sdk.NewInt(300_000_000).ToDec().QuoTruncate(matchingPrice).TruncateInt().Add(sdk.NewInt(200_000_000))
	MatchedAmt2 := sdk.NewInt(500_000_000).ToDec().QuoTruncate(matchingPrice).TruncateInt().Add(sdk.NewInt(100_000_000))
	tMatchedAmt3 := sdk.NewInt(300_000_000).ToDec().QuoTruncate(matchingPrice).TruncateInt().Add(sdk.NewInt(500_000_000))
	MatchedAmt3 := tMatchedAmt3.Add(sdk.NewInt(800_000_000).ToDec().QuoTruncate(matchingPrice).TruncateInt())
	MatchedAmt4 := sdk.NewInt(100_000_000)
	MatchedAmt5 := sdk.NewInt(300_000_000)
	MatchedAmt8 := sdk.NewInt(500_000_000)
	MatchedAmt_Zero := sdk.NewInt(0)
	TotalMatchedAmt := MatchedAmt1.Add(MatchedAmt2).Add(MatchedAmt3).Add(MatchedAmt4).Add(MatchedAmt5).Add(MatchedAmt8)

	s.Require().Equal(mInfo.TotalMatchedAmount, TotalMatchedAmt)
	s.Require().Equal(mInfo.AllocationMap[s.addr(1).String()], MatchedAmt1)
	s.Require().Equal(mInfo.AllocationMap[s.addr(2).String()], MatchedAmt2)
	s.Require().Equal(mInfo.AllocationMap[s.addr(3).String()], MatchedAmt3)
	s.Require().Equal(mInfo.AllocationMap[s.addr(4).String()], MatchedAmt4)
	s.Require().Equal(mInfo.AllocationMap[s.addr(5).String()], MatchedAmt5)
	s.Require().Equal(mInfo.AllocationMap[s.addr(6).String()], MatchedAmt_Zero)
	s.Require().Equal(mInfo.AllocationMap[s.addr(7).String()], MatchedAmt_Zero)
	s.Require().Equal(mInfo.AllocationMap[s.addr(8).String()], MatchedAmt8)
	s.Require().Equal(mInfo.AllocationMap[s.addr(9).String()], MatchedAmt_Zero)
	s.Require().Equal(mInfo.AllocationMap[s.addr(10).String()], MatchedAmt_Zero)

	ReservedMatchedAmt1 := sdk.NewInt(200_000_000).ToDec().Mul(matchingPrice).Ceil().TruncateInt().Add(sdk.NewInt(300_000_000))
	ReservedMatchedAmt2 := sdk.NewInt(100_000_000).ToDec().Mul(matchingPrice).Ceil().TruncateInt().Add(sdk.NewInt(500_000_000))
	ReservedMatchedAmt3 := sdk.NewInt(500_000_000).ToDec().Mul(matchingPrice).Ceil().TruncateInt().Add(sdk.NewInt(1100_000_000))
	ReservedMatchedAmt4 := sdk.NewInt(100_000_000).ToDec().Mul(matchingPrice).Ceil().TruncateInt()
	ReservedMatchedAmt5 := sdk.NewInt(300_000_000).ToDec().Mul(matchingPrice).Ceil().TruncateInt()
	ReservedMatchedAmt8 := sdk.NewInt(500_000_000).ToDec().Mul(matchingPrice).Ceil().TruncateInt()
	ReservedMatchedAmt_Zero := sdk.NewInt(0)

	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(1).String()], ReservedMatchedAmt1)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(2).String()], ReservedMatchedAmt2)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(3).String()], ReservedMatchedAmt3)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(4).String()], ReservedMatchedAmt4)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(5).String()], ReservedMatchedAmt5)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(6).String()], ReservedMatchedAmt_Zero)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(7).String()], ReservedMatchedAmt_Zero)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(8).String()], ReservedMatchedAmt8)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(9).String()], ReservedMatchedAmt_Zero)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(10).String()], ReservedMatchedAmt_Zero)

	RefundAmt1 := sdk.NewInt(60_000_000)
	RefundAmt2 := sdk.NewInt(0)
	RefundAmt3 := sdk.NewInt(100_000_000)
	RefundAmt4 := sdk.NewInt(10_000_000)
	RefundAmt5 := sdk.NewInt(150_000_000)
	RefundAmt6 := sdk.NewInt(550_000_000)
	RefundAmt7 := sdk.NewInt(500_000_000)
	RefundAmt8 := sdk.NewInt(50_000_000)
	RefundAmt9 := sdk.NewInt(360_000_000)
	RefundAmt10 := sdk.NewInt(60_000_000)

	s.Require().Equal(mInfo.RefundMap[s.addr(1).String()].Abs(), RefundAmt1.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(2).String()].Abs(), RefundAmt2.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(3).String()].Abs(), RefundAmt3.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(4).String()].Abs(), RefundAmt4.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(5).String()].Abs(), RefundAmt5.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(6).String()].Abs(), RefundAmt6.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(7).String()].Abs(), RefundAmt7.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(8).String()].Abs(), RefundAmt8.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(9).String()].Abs(), RefundAmt9.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(10).String()].Abs(), RefundAmt10.Abs())

	// Distribute selling coin
	err := s.keeper.AllocateSellingCoin(s.ctx, auction, mInfo)
	s.Require().NoError(err)

	err = s.keeper.ReleaseRemainingSellingCoin(s.ctx, auction)
	s.Require().NoError(err)

	// The selling reserve account balance must be zero
	s.Require().True(s.getBalance(auction.GetSellingReserveAddress(), auction.SellingCoin.Denom).IsZero())

	// The auctioneer must have sellingCoin.Amount - TotalMatchedAmount
	s.Require().Equal(s.getBalance(s.addr(0), auction.GetSellingCoin().Denom).Amount, auction.SellingCoin.Amount.Sub(mInfo.TotalMatchedAmount))

	// The bidders must have the matched selling coin
	s.Require().Equal(s.getBalance(s.addr(1), auction.GetSellingCoin().Denom).Amount, MatchedAmt1)
	s.Require().Equal(s.getBalance(s.addr(2), auction.GetSellingCoin().Denom).Amount, MatchedAmt2)
	s.Require().Equal(s.getBalance(s.addr(3), auction.GetSellingCoin().Denom).Amount, MatchedAmt3)
	s.Require().Equal(s.getBalance(s.addr(4), auction.GetSellingCoin().Denom).Amount, MatchedAmt4)
	s.Require().Equal(s.getBalance(s.addr(5), auction.GetSellingCoin().Denom).Amount, MatchedAmt5)
	s.Require().Equal(s.getBalance(s.addr(6), auction.GetSellingCoin().Denom).Amount.Abs(), MatchedAmt_Zero.Abs())
	s.Require().Equal(s.getBalance(s.addr(7), auction.GetSellingCoin().Denom).Amount.Abs(), MatchedAmt_Zero.Abs())
	s.Require().Equal(s.getBalance(s.addr(8), auction.GetSellingCoin().Denom).Amount, MatchedAmt8)
	s.Require().Equal(s.getBalance(s.addr(9), auction.GetSellingCoin().Denom).Amount.Abs(), MatchedAmt_Zero.Abs())
	s.Require().Equal(s.getBalance(s.addr(10), auction.GetSellingCoin().Denom).Amount.Abs(), MatchedAmt_Zero.Abs())

	// Refund payingCoin
	err = s.keeper.RefundPayingCoin(s.ctx, auction, mInfo)
	s.Require().NoError(err)
}

// Example of "JH_ex1.01" in Sheet for the same MaxBidAmountLimit value
func (s *KeeperTestSuite) TestCalculateAllocation_Mixed2_LimitedSame() {
	auction := s.createBatchAuction(
		s.addr(0),
		parseDec("1"),
		parseDec("0.1"),
		parseCoin("5000_000_000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		1,
		sdk.MustNewDecFromStr("0.2"),
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	s.placeBidBatchMany(auction.Id, s.addr(1), parseDec("1"), parseCoin("200_000_000denom1"), sdk.NewInt(700_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(2), parseDec("0.8"), parseCoin("500_000_000denom2"), sdk.NewInt(700_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(3), parseDec("0.9"), parseCoin("500_000_000denom1"), sdk.NewInt(700_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(1), parseDec("1.1"), parseCoin("300_000_000denom2"), sdk.NewInt(0), true)
	s.placeBidBatchMany(auction.Id, s.addr(5), parseDec("1.2"), parseCoin("300_000_000denom1"), sdk.NewInt(700_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(4), parseDec("0.8"), parseCoin("100_000_000denom1"), sdk.NewInt(700_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(2), parseDec("0.7"), parseCoin("100_000_000denom1"), sdk.NewInt(0), true)
	s.placeBidBatchMany(auction.Id, s.addr(6), parseDec("0.5"), parseCoin("100_000_000denom1"), sdk.NewInt(700_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(3), parseDec("0.8"), parseCoin("300_000_000denom2"), sdk.NewInt(0), true)
	s.placeBidBatchWorth(auction.Id, s.addr(7), parseDec("0.6"), parseCoin("400_000_000denom2"), sdk.NewInt(700_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(8), parseDec("0.8"), parseCoin("500_000_000denom1"), sdk.NewInt(700_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(9), parseDec("0.6"), parseCoin("600_000_000denom1"), sdk.NewInt(700_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(6), parseDec("0.5"), parseCoin("350_000_000denom2"), sdk.NewInt(0), true)
	s.placeBidBatchMany(auction.Id, s.addr(10), parseDec("0.6"), parseCoin("100_000_000denom1"), sdk.NewInt(700_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(3), parseDec("0.7"), parseCoin("490_000_000denom2"), sdk.NewInt(0), true)

	a, found := s.keeper.GetAuction(s.ctx, auction.Id)
	s.Require().True(found)

	mInfo := s.keeper.CalculateBatchAllocation(s.ctx, a)

	// Checking
	s.Require().Equal(mInfo.MatchedLen, int64(13))
	matchingPrice := parseDec("0.6")
	s.Require().Equal(mInfo.MatchedPrice, matchingPrice)

	MatchedAmt1 := sdk.NewInt(700_000_000)
	MatchedAmt2 := sdk.NewInt(700_000_000)
	MatchedAmt3 := sdk.NewInt(700_000_000)
	MatchedAmt4 := sdk.NewInt(100_000_000)
	MatchedAmt5 := sdk.NewInt(300_000_000)
	MatchedAmt6 := sdk.NewInt(0)
	MatchedAmt7 := sdk.NewInt(400_000_000).ToDec().QuoTruncate(matchingPrice).TruncateInt()
	MatchedAmt8 := sdk.NewInt(500_000_000)
	MatchedAmt9 := sdk.NewInt(600_000_000)
	MatchedAmt10 := sdk.NewInt(100_000_000)

	TotalMatchedAmt := sdk.NewInt(3700_000_000).Add(MatchedAmt7)

	s.Require().Equal(mInfo.TotalMatchedAmount, TotalMatchedAmt)
	s.Require().Equal(mInfo.AllocationMap[s.addr(1).String()], MatchedAmt1)
	s.Require().Equal(mInfo.AllocationMap[s.addr(2).String()], MatchedAmt2)
	s.Require().Equal(mInfo.AllocationMap[s.addr(3).String()], MatchedAmt3)
	s.Require().Equal(mInfo.AllocationMap[s.addr(4).String()], MatchedAmt4)
	s.Require().Equal(mInfo.AllocationMap[s.addr(5).String()], MatchedAmt5)
	s.Require().Equal(mInfo.AllocationMap[s.addr(6).String()], MatchedAmt6)
	s.Require().Equal(mInfo.AllocationMap[s.addr(7).String()], MatchedAmt7)
	s.Require().Equal(mInfo.AllocationMap[s.addr(8).String()], MatchedAmt8)
	s.Require().Equal(mInfo.AllocationMap[s.addr(9).String()], MatchedAmt9)
	s.Require().Equal(mInfo.AllocationMap[s.addr(10).String()], MatchedAmt10)

	ReservedMatchedAmt1 := sdk.NewInt(420_000_000)
	ReservedMatchedAmt2 := sdk.NewInt(420_000_000)
	ReservedMatchedAmt3 := sdk.NewInt(420_000_000)
	ReservedMatchedAmt4 := sdk.NewInt(60_000_000)
	ReservedMatchedAmt5 := sdk.NewInt(180_000_000)
	ReservedMatchedAmt6 := sdk.NewInt(0)
	ReservedMatchedAmt7 := sdk.NewInt(400_000_000)
	ReservedMatchedAmt8 := sdk.NewInt(300_000_000)
	ReservedMatchedAmt9 := sdk.NewInt(360_000_000)
	ReservedMatchedAmt10 := sdk.NewInt(60_000_000)

	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(1).String()], ReservedMatchedAmt1)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(2).String()], ReservedMatchedAmt2)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(3).String()], ReservedMatchedAmt3)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(4).String()], ReservedMatchedAmt4)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(5).String()], ReservedMatchedAmt5)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(6).String()], ReservedMatchedAmt6)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(7).String()], ReservedMatchedAmt7)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(8).String()], ReservedMatchedAmt8)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(9).String()], ReservedMatchedAmt9)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(10).String()], ReservedMatchedAmt10)

	RefundAmt1 := sdk.NewInt(80_000_000)
	RefundAmt2 := sdk.NewInt(150_000_000)
	RefundAmt3 := sdk.NewInt(820_000_000)
	RefundAmt4 := sdk.NewInt(20_000_000)
	RefundAmt5 := sdk.NewInt(180_000_000)
	RefundAmt6 := sdk.NewInt(400_000_000)
	RefundAmt7 := sdk.NewInt(0)
	RefundAmt8 := sdk.NewInt(100_000_000)
	RefundAmt9 := sdk.NewInt(0)
	RefundAmt10 := sdk.NewInt(0)

	s.Require().Equal(mInfo.RefundMap[s.addr(1).String()].Abs(), RefundAmt1.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(2).String()].Abs(), RefundAmt2.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(3).String()].Abs(), RefundAmt3.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(4).String()].Abs(), RefundAmt4.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(5).String()].Abs(), RefundAmt5.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(6).String()].Abs(), RefundAmt6.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(7).String()].Abs(), RefundAmt7.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(8).String()].Abs(), RefundAmt8.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(9).String()].Abs(), RefundAmt9.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(10).String()].Abs(), RefundAmt10.Abs())

	// Distribute selling coin
	err := s.keeper.AllocateSellingCoin(s.ctx, auction, mInfo)
	s.Require().NoError(err)

	err = s.keeper.ReleaseRemainingSellingCoin(s.ctx, auction)
	s.Require().NoError(err)

	// The selling reserve account balance must be zero
	s.Require().True(s.getBalance(auction.GetSellingReserveAddress(), auction.SellingCoin.Denom).IsZero())

	// The auctioneer must have sellingCoin.Amount - TotalMatchedAmount
	s.Require().Equal(s.getBalance(s.addr(0), auction.GetSellingCoin().Denom).Amount, auction.SellingCoin.Amount.Sub(mInfo.TotalMatchedAmount))

	// The bidders must have the matched selling coin
	s.Require().Equal(s.getBalance(s.addr(1), auction.GetSellingCoin().Denom).Amount, MatchedAmt1)
	s.Require().Equal(s.getBalance(s.addr(2), auction.GetSellingCoin().Denom).Amount, MatchedAmt2)
	s.Require().Equal(s.getBalance(s.addr(3), auction.GetSellingCoin().Denom).Amount, MatchedAmt3)
	s.Require().Equal(s.getBalance(s.addr(4), auction.GetSellingCoin().Denom).Amount, MatchedAmt4)
	s.Require().Equal(s.getBalance(s.addr(5), auction.GetSellingCoin().Denom).Amount, MatchedAmt5)
	s.Require().Equal(s.getBalance(s.addr(6), auction.GetSellingCoin().Denom).Amount.Abs(), MatchedAmt6.Abs())
	s.Require().Equal(s.getBalance(s.addr(7), auction.GetSellingCoin().Denom).Amount, MatchedAmt7)
	s.Require().Equal(s.getBalance(s.addr(8), auction.GetSellingCoin().Denom).Amount, MatchedAmt8)
	s.Require().Equal(s.getBalance(s.addr(9), auction.GetSellingCoin().Denom).Amount, MatchedAmt9)
	s.Require().Equal(s.getBalance(s.addr(10), auction.GetSellingCoin().Denom).Amount, MatchedAmt10)

	// Refund payingCoin
	err = s.keeper.RefundPayingCoin(s.ctx, auction, mInfo)
	s.Require().NoError(err)
}

// Example of "JH_ex1.1" in Sheet for different MaxBidAmountLimit values
func (s *KeeperTestSuite) TestCalculateAllocation_Mixed2_LimitedDifferent() {
	auction := s.createBatchAuction(
		s.addr(0),
		parseDec("1"),
		parseDec("0.1"),
		parseCoin("5000_000_000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		1,
		sdk.MustNewDecFromStr("0.2"),
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	s.placeBidBatchMany(auction.Id, s.addr(1), parseDec("1"), parseCoin("200_000_000denom1"), sdk.NewInt(1000_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(2), parseDec("0.8"), parseCoin("500_000_000denom2"), sdk.NewInt(1000_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(3), parseDec("0.9"), parseCoin("500_000_000denom1"), sdk.NewInt(800_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(1), parseDec("1.1"), parseCoin("300_000_000denom2"), sdk.NewInt(0), true)
	s.placeBidBatchMany(auction.Id, s.addr(5), parseDec("1.2"), parseCoin("300_000_000denom1"), sdk.NewInt(600_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(4), parseDec("0.8"), parseCoin("100_000_000denom1"), sdk.NewInt(800_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(2), parseDec("0.7"), parseCoin("100_000_000denom1"), sdk.NewInt(0), true)
	s.placeBidBatchMany(auction.Id, s.addr(6), parseDec("0.5"), parseCoin("100_000_000denom1"), sdk.NewInt(600_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(3), parseDec("0.8"), parseCoin("300_000_000denom2"), sdk.NewInt(0), true)
	s.placeBidBatchWorth(auction.Id, s.addr(7), parseDec("0.6"), parseCoin("200_000_000denom2"), sdk.NewInt(400_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(8), parseDec("0.8"), parseCoin("400_000_000denom1"), sdk.NewInt(400_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(9), parseDec("0.6"), parseCoin("200_000_000denom1"), sdk.NewInt(200_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(6), parseDec("0.5"), parseCoin("300_000_000denom2"), sdk.NewInt(0), true)
	s.placeBidBatchMany(auction.Id, s.addr(10), parseDec("0.6"), parseCoin("100_000_000denom1"), sdk.NewInt(200_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(3), parseDec("0.7"), parseCoin("560_000_000denom2"), sdk.NewInt(0), true)

	a, found := s.keeper.GetAuction(s.ctx, auction.Id)
	s.Require().True(found)

	mInfo := s.keeper.CalculateBatchAllocation(s.ctx, a)

	// Checking
	s.Require().Equal(mInfo.MatchedLen, int64(15))
	matchingPrice := parseDec("0.5")
	s.Require().Equal(mInfo.MatchedPrice, matchingPrice)

	MatchedAmt1 := sdk.NewInt(800_000_000)
	MatchedAmt2 := sdk.NewInt(1000_000_000)
	MatchedAmt3 := sdk.NewInt(800_000_000)
	MatchedAmt4 := sdk.NewInt(100_000_000)
	MatchedAmt5 := sdk.NewInt(300_000_000)
	MatchedAmt6 := sdk.NewInt(600_000_000)
	MatchedAmt7 := sdk.NewInt(400_000_000)
	MatchedAmt8 := sdk.NewInt(400_000_000)
	MatchedAmt9 := sdk.NewInt(200_000_000)
	MatchedAmt10 := sdk.NewInt(100_000_000)

	TotalMatchedAmt := sdk.NewInt(4700_000_000)

	s.Require().Equal(mInfo.TotalMatchedAmount, TotalMatchedAmt)
	s.Require().Equal(mInfo.AllocationMap[s.addr(1).String()], MatchedAmt1)
	s.Require().Equal(mInfo.AllocationMap[s.addr(2).String()], MatchedAmt2)
	s.Require().Equal(mInfo.AllocationMap[s.addr(3).String()], MatchedAmt3)
	s.Require().Equal(mInfo.AllocationMap[s.addr(4).String()], MatchedAmt4)
	s.Require().Equal(mInfo.AllocationMap[s.addr(5).String()], MatchedAmt5)
	s.Require().Equal(mInfo.AllocationMap[s.addr(6).String()], MatchedAmt6)
	s.Require().Equal(mInfo.AllocationMap[s.addr(7).String()], MatchedAmt7)
	s.Require().Equal(mInfo.AllocationMap[s.addr(8).String()], MatchedAmt8)
	s.Require().Equal(mInfo.AllocationMap[s.addr(9).String()], MatchedAmt9)
	s.Require().Equal(mInfo.AllocationMap[s.addr(10).String()], MatchedAmt10)

	ReservedMatchedAmt1 := sdk.NewInt(400_000_000)
	ReservedMatchedAmt2 := sdk.NewInt(500_000_000)
	ReservedMatchedAmt3 := sdk.NewInt(400_000_000)
	ReservedMatchedAmt4 := sdk.NewInt(50_000_000)
	ReservedMatchedAmt5 := sdk.NewInt(150_000_000)
	ReservedMatchedAmt6 := sdk.NewInt(300_000_000)
	ReservedMatchedAmt7 := sdk.NewInt(200_000_000)
	ReservedMatchedAmt8 := sdk.NewInt(200_000_000)
	ReservedMatchedAmt9 := sdk.NewInt(100_000_000)
	ReservedMatchedAmt10 := sdk.NewInt(50_000_000)

	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(1).String()], ReservedMatchedAmt1)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(2).String()], ReservedMatchedAmt2)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(3).String()], ReservedMatchedAmt3)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(4).String()], ReservedMatchedAmt4)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(5).String()], ReservedMatchedAmt5)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(6).String()], ReservedMatchedAmt6)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(7).String()], ReservedMatchedAmt7)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(8).String()], ReservedMatchedAmt8)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(9).String()], ReservedMatchedAmt9)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(10).String()], ReservedMatchedAmt10)

	RefundAmt1 := sdk.NewInt(100_000_000)
	RefundAmt2 := sdk.NewInt(70_000_000)
	RefundAmt3 := sdk.NewInt(910_000_000)
	RefundAmt4 := sdk.NewInt(30_000_000)
	RefundAmt5 := sdk.NewInt(210_000_000)
	RefundAmt6 := sdk.NewInt(50_000_000)
	RefundAmt7 := sdk.NewInt(0)
	RefundAmt8 := sdk.NewInt(120_000_000)
	RefundAmt9 := sdk.NewInt(20_000_000)
	RefundAmt10 := sdk.NewInt(10_000_000)

	s.Require().Equal(mInfo.RefundMap[s.addr(1).String()].Abs(), RefundAmt1.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(2).String()].Abs(), RefundAmt2.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(3).String()].Abs(), RefundAmt3.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(4).String()].Abs(), RefundAmt4.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(5).String()].Abs(), RefundAmt5.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(6).String()].Abs(), RefundAmt6.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(7).String()].Abs(), RefundAmt7.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(8).String()].Abs(), RefundAmt8.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(9).String()].Abs(), RefundAmt9.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(10).String()].Abs(), RefundAmt10.Abs())

	// Distribute selling coin
	err := s.keeper.AllocateSellingCoin(s.ctx, auction, mInfo)
	s.Require().NoError(err)

	err = s.keeper.ReleaseRemainingSellingCoin(s.ctx, auction)
	s.Require().NoError(err)

	// The selling reserve account balance must be zero
	s.Require().True(s.getBalance(auction.GetSellingReserveAddress(), auction.SellingCoin.Denom).IsZero())

	// The auctioneer must have sellingCoin.Amount - TotalMatchedAmount
	s.Require().Equal(s.getBalance(s.addr(0), auction.GetSellingCoin().Denom).Amount, sdk.NewInt(300_000_000))

	// The bidders must have the matched selling coin
	s.Require().Equal(s.getBalance(s.addr(1), auction.GetSellingCoin().Denom).Amount, MatchedAmt1)
	s.Require().Equal(s.getBalance(s.addr(2), auction.GetSellingCoin().Denom).Amount, MatchedAmt2)
	s.Require().Equal(s.getBalance(s.addr(3), auction.GetSellingCoin().Denom).Amount, MatchedAmt3)
	s.Require().Equal(s.getBalance(s.addr(4), auction.GetSellingCoin().Denom).Amount, MatchedAmt4)
	s.Require().Equal(s.getBalance(s.addr(5), auction.GetSellingCoin().Denom).Amount, MatchedAmt5)
	s.Require().Equal(s.getBalance(s.addr(6), auction.GetSellingCoin().Denom).Amount, MatchedAmt6)
	s.Require().Equal(s.getBalance(s.addr(7), auction.GetSellingCoin().Denom).Amount, MatchedAmt7)
	s.Require().Equal(s.getBalance(s.addr(8), auction.GetSellingCoin().Denom).Amount, MatchedAmt8)
	s.Require().Equal(s.getBalance(s.addr(9), auction.GetSellingCoin().Denom).Amount, MatchedAmt9)
	s.Require().Equal(s.getBalance(s.addr(10), auction.GetSellingCoin().Denom).Amount, MatchedAmt10)

	// Refund payingCoin
	err = s.keeper.RefundPayingCoin(s.ctx, auction, mInfo)
	s.Require().NoError(err)
}

// Example of "JH_ex2" in Sheet without MaxBidAmountLimit value
func (s *KeeperTestSuite) TestCalculateAllocation_Mixed3() {
	auction := s.createBatchAuction(
		s.addr(0),
		parseDec("10"),
		parseDec("0.1"),
		parseCoin("2500_000_000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		1,
		sdk.MustNewDecFromStr("0.2"),
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	s.placeBidBatchMany(auction.Id, s.addr(1), parseDec("10"), parseCoin("200_000_000denom1"), sdk.NewInt(2500_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(2), parseDec("11"), parseCoin("2000_000_000denom2"), sdk.NewInt(2500_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(3), parseDec("10.5"), parseCoin("500_000_000denom1"), sdk.NewInt(2500_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(4), parseDec("10.2"), parseCoin("1500_000_000denom2"), sdk.NewInt(2500_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(5), parseDec("10.8"), parseCoin("300_000_000denom1"), sdk.NewInt(2500_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(6), parseDec("11.4"), parseCoin("2500_000_000denom2"), sdk.NewInt(2500_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(7), parseDec("11.3"), parseCoin("100_000_000denom1"), sdk.NewInt(2500_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(8), parseDec("9.9"), parseCoin("2500_000_000denom2"), sdk.NewInt(2500_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(9), parseDec("10.1"), parseCoin("300_000_000denom1"), sdk.NewInt(2500_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(10), parseDec("10.45"), parseCoin("2000_000_000denom2"), sdk.NewInt(2500_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(11), parseDec("10.75"), parseCoin("150_000_000denom1"), sdk.NewInt(2500_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(12), parseDec("10.99"), parseCoin("1500_000_000denom2"), sdk.NewInt(2500_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(13), parseDec("10.2"), parseCoin("200_000_000denom1"), sdk.NewInt(2500_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(14), parseDec("9.87"), parseCoin("2000_000_000denom2"), sdk.NewInt(2500_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(15), parseDec("10.25"), parseCoin("200_000_000denom1"), sdk.NewInt(2500_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(16), parseDec("10.48"), parseCoin("2500_000_000denom2"), sdk.NewInt(2500_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(17), parseDec("10.52"), parseCoin("180_000_000denom1"), sdk.NewInt(2500_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(4), parseDec("10.8"), parseCoin("220_000_000denom1"), sdk.NewInt(0), true)
	s.placeBidBatchWorth(auction.Id, s.addr(5), parseDec("10.5"), parseCoin("1500_000_000denom2"), sdk.NewInt(0), true)
	s.placeBidBatchMany(auction.Id, s.addr(6), parseDec("9.7"), parseCoin("250_000_000denom1"), sdk.NewInt(0), true)

	a, found := s.keeper.GetAuction(s.ctx, auction.Id)
	s.Require().True(found)

	mInfo := s.keeper.CalculateBatchAllocation(s.ctx, a)

	// Checking
	s.Require().Equal(mInfo.MatchedLen, int64(11))
	matchingPrice := parseDec("10.48")
	s.Require().Equal(mInfo.MatchedPrice, matchingPrice)

	MatchedAmt1 := sdk.NewInt(0)
	MatchedAmt2 := sdk.NewInt(2000_000_000).ToDec().QuoTruncate(matchingPrice).TruncateInt()
	MatchedAmt3 := sdk.NewInt(500_000_000)
	MatchedAmt4 := sdk.NewInt(220_000_000)
	MatchedAmt5 := sdk.NewInt(1500_000_000).ToDec().QuoTruncate(matchingPrice).TruncateInt().Add(sdk.NewInt(300_000_000))
	MatchedAmt6 := sdk.NewInt(2500_000_000).ToDec().QuoTruncate(matchingPrice).TruncateInt()
	MatchedAmt7 := sdk.NewInt(100_000_000)
	MatchedAmt8 := sdk.NewInt(0)
	MatchedAmt9 := sdk.NewInt(0)
	MatchedAmt10 := sdk.NewInt(0)
	MatchedAmt11 := sdk.NewInt(150_000_000)
	MatchedAmt12 := sdk.NewInt(1500_000_000).ToDec().QuoTruncate(matchingPrice).TruncateInt()
	MatchedAmt13 := sdk.NewInt(0)
	MatchedAmt14 := sdk.NewInt(0)
	MatchedAmt15 := sdk.NewInt(0)
	MatchedAmt16 := sdk.NewInt(2500_000_000).ToDec().QuoTruncate(matchingPrice).TruncateInt()
	MatchedAmt17 := sdk.NewInt(180_000_000)

	TotalMatchedAmt := MatchedAmt2.Add(MatchedAmt3).Add(MatchedAmt4).Add(MatchedAmt5).Add(MatchedAmt6).Add(MatchedAmt7).Add(MatchedAmt11).Add(MatchedAmt12).Add(MatchedAmt16).Add(MatchedAmt17)

	s.Require().Equal(mInfo.TotalMatchedAmount, TotalMatchedAmt)
	s.Require().Equal(mInfo.AllocationMap[s.addr(1).String()], MatchedAmt1)
	s.Require().Equal(mInfo.AllocationMap[s.addr(2).String()], MatchedAmt2)
	s.Require().Equal(mInfo.AllocationMap[s.addr(3).String()], MatchedAmt3)
	s.Require().Equal(mInfo.AllocationMap[s.addr(4).String()], MatchedAmt4)
	s.Require().Equal(mInfo.AllocationMap[s.addr(5).String()], MatchedAmt5)
	s.Require().Equal(mInfo.AllocationMap[s.addr(6).String()], MatchedAmt6)
	s.Require().Equal(mInfo.AllocationMap[s.addr(7).String()], MatchedAmt7)
	s.Require().Equal(mInfo.AllocationMap[s.addr(8).String()], MatchedAmt8)
	s.Require().Equal(mInfo.AllocationMap[s.addr(9).String()], MatchedAmt9)
	s.Require().Equal(mInfo.AllocationMap[s.addr(10).String()], MatchedAmt10)
	s.Require().Equal(mInfo.AllocationMap[s.addr(11).String()], MatchedAmt11)
	s.Require().Equal(mInfo.AllocationMap[s.addr(12).String()], MatchedAmt12)
	s.Require().Equal(mInfo.AllocationMap[s.addr(13).String()], MatchedAmt13)
	s.Require().Equal(mInfo.AllocationMap[s.addr(14).String()], MatchedAmt14)
	s.Require().Equal(mInfo.AllocationMap[s.addr(15).String()], MatchedAmt15)
	s.Require().Equal(mInfo.AllocationMap[s.addr(16).String()], MatchedAmt16)
	s.Require().Equal(mInfo.AllocationMap[s.addr(17).String()], MatchedAmt17)

	ReservedMatchedAmt1 := sdk.NewInt(0)
	ReservedMatchedAmt2 := sdk.NewInt(2000_000_000)
	ReservedMatchedAmt3 := MatchedAmt3.ToDec().Mul(matchingPrice).Ceil().TruncateInt()
	ReservedMatchedAmt4 := MatchedAmt4.ToDec().Mul(matchingPrice).Ceil().TruncateInt()
	ReservedMatchedAmt5 := sdk.NewInt(300_000_000).ToDec().Mul(matchingPrice).Ceil().TruncateInt().Add(sdk.NewInt(1500_000_000))
	ReservedMatchedAmt6 := sdk.NewInt(2500_000_000)
	ReservedMatchedAmt7 := MatchedAmt7.ToDec().Mul(matchingPrice).Ceil().TruncateInt()
	ReservedMatchedAmt8 := sdk.NewInt(0)
	ReservedMatchedAmt9 := sdk.NewInt(0)
	ReservedMatchedAmt10 := sdk.NewInt(0)
	ReservedMatchedAmt11 := MatchedAmt11.ToDec().Mul(matchingPrice).Ceil().TruncateInt()
	ReservedMatchedAmt12 := sdk.NewInt(1500_000_000)
	ReservedMatchedAmt13 := sdk.NewInt(0)
	ReservedMatchedAmt14 := sdk.NewInt(0)
	ReservedMatchedAmt15 := sdk.NewInt(0)
	ReservedMatchedAmt16 := sdk.NewInt(2500_000_000)
	ReservedMatchedAmt17 := MatchedAmt17.ToDec().Mul(matchingPrice).Ceil().TruncateInt()

	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(1).String()], ReservedMatchedAmt1)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(2).String()], ReservedMatchedAmt2)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(3).String()], ReservedMatchedAmt3)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(4).String()], ReservedMatchedAmt4)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(5).String()], ReservedMatchedAmt5)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(6).String()], ReservedMatchedAmt6)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(7).String()], ReservedMatchedAmt7)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(8).String()], ReservedMatchedAmt8)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(9).String()], ReservedMatchedAmt9)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(10).String()], ReservedMatchedAmt10)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(11).String()], ReservedMatchedAmt11)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(12).String()], ReservedMatchedAmt12)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(13).String()], ReservedMatchedAmt13)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(14).String()], ReservedMatchedAmt14)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(15).String()], ReservedMatchedAmt15)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(16).String()], ReservedMatchedAmt16)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(17).String()], ReservedMatchedAmt17)

	RefundAmt1 := sdk.NewInt(200_000_000).ToDec().Mul(parseDec("10")).Ceil().TruncateInt()
	RefundAmt2 := sdk.NewInt(2000_000_000).Sub(ReservedMatchedAmt2)
	RefundAmt3 := sdk.NewInt(500_000_000).ToDec().Mul(parseDec("10.5")).Ceil().TruncateInt().Sub(ReservedMatchedAmt3)
	RefundAmt4 := sdk.NewInt(220_000_000).ToDec().Mul(parseDec("10.8")).Ceil().TruncateInt().Add(sdk.NewInt(1500_000_000)).Sub(ReservedMatchedAmt4)
	RefundAmt5 := sdk.NewInt(300_000_000).ToDec().Mul(parseDec("10.8")).Ceil().TruncateInt().Add(sdk.NewInt(1500_000_000)).Sub(ReservedMatchedAmt5)
	RefundAmt6 := sdk.NewInt(250_000_000).ToDec().Mul(parseDec("9.7")).Ceil().TruncateInt().Add(sdk.NewInt(2500_000_000)).Sub(ReservedMatchedAmt6)
	RefundAmt7 := sdk.NewInt(100_000_000).ToDec().Mul(parseDec("11.3")).Ceil().TruncateInt().Sub(ReservedMatchedAmt7)
	RefundAmt8 := sdk.NewInt(2500_000_000)
	RefundAmt9 := sdk.NewInt(300_000_000).ToDec().Mul(parseDec("10.1")).Ceil().TruncateInt()
	RefundAmt10 := sdk.NewInt(2000_000_000)
	RefundAmt11 := sdk.NewInt(150_000_000).ToDec().Mul(parseDec("10.75")).Ceil().TruncateInt().Sub(ReservedMatchedAmt11)
	RefundAmt12 := sdk.NewInt(1500_000_000).Sub(ReservedMatchedAmt12)
	RefundAmt13 := sdk.NewInt(200_000_000).ToDec().Mul(parseDec("10.2")).Ceil().TruncateInt()
	RefundAmt14 := sdk.NewInt(2000_000_000)
	RefundAmt15 := sdk.NewInt(200_000_000).ToDec().Mul(parseDec("10.25")).Ceil().TruncateInt()
	RefundAmt16 := sdk.NewInt(2500_000_000).Sub(ReservedMatchedAmt16)
	RefundAmt17 := sdk.NewInt(180_000_000).ToDec().Mul(parseDec("10.52")).Ceil().TruncateInt().Sub(ReservedMatchedAmt17)

	s.Require().Equal(mInfo.RefundMap[s.addr(1).String()].Abs(), RefundAmt1.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(2).String()].Abs(), RefundAmt2.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(3).String()].Abs(), RefundAmt3.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(4).String()].Abs(), RefundAmt4.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(5).String()].Abs(), RefundAmt5.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(6).String()].Abs(), RefundAmt6.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(7).String()].Abs(), RefundAmt7.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(8).String()].Abs(), RefundAmt8.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(9).String()].Abs(), RefundAmt9.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(10).String()].Abs(), RefundAmt10.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(11).String()].Abs(), RefundAmt11.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(12).String()].Abs(), RefundAmt12.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(13).String()].Abs(), RefundAmt13.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(14).String()].Abs(), RefundAmt14.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(15).String()].Abs(), RefundAmt15.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(16).String()].Abs(), RefundAmt16.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(17).String()].Abs(), RefundAmt17.Abs())

	// Distribute selling coin
	err := s.keeper.AllocateSellingCoin(s.ctx, auction, mInfo)
	s.Require().NoError(err)

	err = s.keeper.ReleaseRemainingSellingCoin(s.ctx, auction)
	s.Require().NoError(err)

	// The selling reserve account balance must be zero
	s.Require().True(s.getBalance(auction.GetSellingReserveAddress(), auction.SellingCoin.Denom).IsZero())

	// The auctioneer must have sellingCoin.Amount - TotalMatchedAmount
	s.Require().Equal(s.getBalance(s.addr(0), auction.GetSellingCoin().Denom).Amount, auction.SellingCoin.Amount.Sub(mInfo.TotalMatchedAmount))

	// The bidders must have the matched selling coin
	s.Require().Equal(s.getBalance(s.addr(1), auction.GetSellingCoin().Denom).Amount.Abs(), MatchedAmt1.Abs())
	s.Require().Equal(s.getBalance(s.addr(2), auction.GetSellingCoin().Denom).Amount.Abs(), MatchedAmt2.Abs())
	s.Require().Equal(s.getBalance(s.addr(3), auction.GetSellingCoin().Denom).Amount.Abs(), MatchedAmt3.Abs())
	s.Require().Equal(s.getBalance(s.addr(4), auction.GetSellingCoin().Denom).Amount.Abs(), MatchedAmt4.Abs())
	s.Require().Equal(s.getBalance(s.addr(5), auction.GetSellingCoin().Denom).Amount.Abs(), MatchedAmt5.Abs())
	s.Require().Equal(s.getBalance(s.addr(6), auction.GetSellingCoin().Denom).Amount.Abs(), MatchedAmt6.Abs())
	s.Require().Equal(s.getBalance(s.addr(7), auction.GetSellingCoin().Denom).Amount.Abs(), MatchedAmt7.Abs())
	s.Require().Equal(s.getBalance(s.addr(8), auction.GetSellingCoin().Denom).Amount.Abs(), MatchedAmt8.Abs())
	s.Require().Equal(s.getBalance(s.addr(9), auction.GetSellingCoin().Denom).Amount.Abs(), MatchedAmt9.Abs())
	s.Require().Equal(s.getBalance(s.addr(10), auction.GetSellingCoin().Denom).Amount.Abs(), MatchedAmt10.Abs())
	s.Require().Equal(s.getBalance(s.addr(11), auction.GetSellingCoin().Denom).Amount.Abs(), MatchedAmt11.Abs())
	s.Require().Equal(s.getBalance(s.addr(12), auction.GetSellingCoin().Denom).Amount.Abs(), MatchedAmt12.Abs())
	s.Require().Equal(s.getBalance(s.addr(13), auction.GetSellingCoin().Denom).Amount.Abs(), MatchedAmt13.Abs())
	s.Require().Equal(s.getBalance(s.addr(14), auction.GetSellingCoin().Denom).Amount.Abs(), MatchedAmt14.Abs())
	s.Require().Equal(s.getBalance(s.addr(15), auction.GetSellingCoin().Denom).Amount.Abs(), MatchedAmt15.Abs())
	s.Require().Equal(s.getBalance(s.addr(16), auction.GetSellingCoin().Denom).Amount.Abs(), MatchedAmt16.Abs())
	s.Require().Equal(s.getBalance(s.addr(17), auction.GetSellingCoin().Denom).Amount.Abs(), MatchedAmt17.Abs())

	// Refund payingCoin
	err = s.keeper.RefundPayingCoin(s.ctx, auction, mInfo)
	s.Require().NoError(err)
}

// Example of "JH_ex2" in Sheet for different MaxBidAmountLimit values
func (s *KeeperTestSuite) TestCalculateAllocation_Mixed3_LimitedDifferent() {
	auction := s.createBatchAuction(
		s.addr(0),
		parseDec("10"),
		parseDec("0.1"),
		parseCoin("2500_000_000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		1,
		sdk.MustNewDecFromStr("0.2"),
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	s.placeBidBatchMany(auction.Id, s.addr(1), parseDec("10"), parseCoin("200_000_000denom1"), sdk.NewInt(500_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(2), parseDec("11"), parseCoin("2000_000_000denom2"), sdk.NewInt(500_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(3), parseDec("10.5"), parseCoin("500_000_000denom1"), sdk.NewInt(500_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(4), parseDec("10.2"), parseCoin("1500_000_000denom2"), sdk.NewInt(200_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(5), parseDec("10.8"), parseCoin("200_000_000denom1"), sdk.NewInt(200_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(6), parseDec("11.4"), parseCoin("2200_000_000denom2"), sdk.NewInt(200_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(7), parseDec("11.3"), parseCoin("100_000_000denom1"), sdk.NewInt(200_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(8), parseDec("9.9"), parseCoin("1900_000_000denom2"), sdk.NewInt(200_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(9), parseDec("10.1"), parseCoin("200_000_000denom1"), sdk.NewInt(200_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(10), parseDec("10.45"), parseCoin("2000_000_000denom2"), sdk.NewInt(200_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(11), parseDec("10.75"), parseCoin("100_000_000denom1"), sdk.NewInt(100_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(12), parseDec("10.99"), parseCoin("1050_000_000denom2"), sdk.NewInt(100_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(13), parseDec("10.2"), parseCoin("100_000_000denom1"), sdk.NewInt(100_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(14), parseDec("9.87"), parseCoin("980_000_000denom2"), sdk.NewInt(100_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(15), parseDec("10.25"), parseCoin("100_000_000denom1"), sdk.NewInt(100_000_000), true)
	s.placeBidBatchWorth(auction.Id, s.addr(16), parseDec("10.48"), parseCoin("1000_000_000denom2"), sdk.NewInt(100_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(17), parseDec("10.52"), parseCoin("100_000_000denom1"), sdk.NewInt(100_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(4), parseDec("10.8"), parseCoin("200_000_000denom1"), sdk.NewInt(0), true)
	s.placeBidBatchWorth(auction.Id, s.addr(5), parseDec("10.5"), parseCoin("1500_000_000denom2"), sdk.NewInt(0), true)
	s.placeBidBatchMany(auction.Id, s.addr(6), parseDec("9.7"), parseCoin("200_000_000denom1"), sdk.NewInt(0), true)

	a, found := s.keeper.GetAuction(s.ctx, auction.Id)
	s.Require().True(found)

	mInfo := s.keeper.CalculateBatchAllocation(s.ctx, a)

	// Checking
	s.Require().Equal(mInfo.MatchedLen, int64(16))
	matchingPrice := parseDec("10.1")
	s.Require().Equal(mInfo.MatchedPrice, matchingPrice)

	MatchedAmt1 := sdk.NewInt(0)
	MatchedAmt2 := sdk.NewInt(2000_000_000).ToDec().QuoTruncate(matchingPrice).TruncateInt()
	MatchedAmt3 := sdk.NewInt(500_000_000)
	MatchedAmt4 := sdk.NewInt(200_000_000)
	MatchedAmt5 := sdk.NewInt(200_000_000)
	MatchedAmt6 := sdk.NewInt(200_000_000)
	MatchedAmt7 := sdk.NewInt(100_000_000)
	MatchedAmt8 := sdk.NewInt(0)
	MatchedAmt9 := sdk.NewInt(200_000_000)
	MatchedAmt10 := sdk.NewInt(2000_000_000).ToDec().QuoTruncate(matchingPrice).TruncateInt()
	MatchedAmt11 := sdk.NewInt(100_000_000)
	MatchedAmt12 := sdk.NewInt(100_000_000)
	MatchedAmt13 := sdk.NewInt(100_000_000)
	MatchedAmt14 := sdk.NewInt(0)
	MatchedAmt15 := sdk.NewInt(100_000_000)
	MatchedAmt16 := sdk.NewInt(1000_000_000).ToDec().QuoTruncate(matchingPrice).TruncateInt()
	MatchedAmt17 := sdk.NewInt(100_000_000)

	TotalMatchedAmt := MatchedAmt2.Add(MatchedAmt3).Add(MatchedAmt4).Add(MatchedAmt5).Add(MatchedAmt6).Add(MatchedAmt7).Add(MatchedAmt9).Add(MatchedAmt10).Add(MatchedAmt11).Add(MatchedAmt12).Add(MatchedAmt13).Add(MatchedAmt15).Add(MatchedAmt16).Add(MatchedAmt17)

	s.Require().Equal(mInfo.TotalMatchedAmount, TotalMatchedAmt)
	s.Require().Equal(mInfo.AllocationMap[s.addr(1).String()], MatchedAmt1)
	s.Require().Equal(mInfo.AllocationMap[s.addr(2).String()], MatchedAmt2)
	s.Require().Equal(mInfo.AllocationMap[s.addr(3).String()], MatchedAmt3)
	s.Require().Equal(mInfo.AllocationMap[s.addr(4).String()], MatchedAmt4)
	s.Require().Equal(mInfo.AllocationMap[s.addr(5).String()], MatchedAmt5)
	s.Require().Equal(mInfo.AllocationMap[s.addr(6).String()], MatchedAmt6)
	s.Require().Equal(mInfo.AllocationMap[s.addr(7).String()], MatchedAmt7)
	s.Require().Equal(mInfo.AllocationMap[s.addr(8).String()], MatchedAmt8)
	s.Require().Equal(mInfo.AllocationMap[s.addr(9).String()], MatchedAmt9)
	s.Require().Equal(mInfo.AllocationMap[s.addr(10).String()], MatchedAmt10)
	s.Require().Equal(mInfo.AllocationMap[s.addr(11).String()], MatchedAmt11)
	s.Require().Equal(mInfo.AllocationMap[s.addr(12).String()], MatchedAmt12)
	s.Require().Equal(mInfo.AllocationMap[s.addr(13).String()], MatchedAmt13)
	s.Require().Equal(mInfo.AllocationMap[s.addr(14).String()], MatchedAmt14)
	s.Require().Equal(mInfo.AllocationMap[s.addr(15).String()], MatchedAmt15)
	s.Require().Equal(mInfo.AllocationMap[s.addr(16).String()], MatchedAmt16)
	s.Require().Equal(mInfo.AllocationMap[s.addr(17).String()], MatchedAmt17)

	ReservedMatchedAmt1 := sdk.NewInt(0)
	ReservedMatchedAmt2 := sdk.NewInt(2000_000_000)
	ReservedMatchedAmt3 := sdk.NewInt(5050_000_000)
	ReservedMatchedAmt4 := sdk.NewInt(2020_000_000)
	ReservedMatchedAmt5 := sdk.NewInt(2020_000_000)
	ReservedMatchedAmt6 := sdk.NewInt(2020_000_000)
	ReservedMatchedAmt7 := sdk.NewInt(1010_000_000)
	ReservedMatchedAmt8 := sdk.NewInt(0)
	ReservedMatchedAmt9 := sdk.NewInt(2020_000_000)
	ReservedMatchedAmt10 := sdk.NewInt(2000_000_000)
	ReservedMatchedAmt11 := sdk.NewInt(1010_000_000)
	ReservedMatchedAmt12 := sdk.NewInt(1010_000_000)
	ReservedMatchedAmt13 := sdk.NewInt(1010_000_000)
	ReservedMatchedAmt14 := sdk.NewInt(0)
	ReservedMatchedAmt15 := sdk.NewInt(1010_000_000)
	ReservedMatchedAmt16 := sdk.NewInt(1000_000_000)
	ReservedMatchedAmt17 := sdk.NewInt(1010_000_000)

	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(1).String()], ReservedMatchedAmt1)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(2).String()], ReservedMatchedAmt2)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(3).String()], ReservedMatchedAmt3)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(4).String()], ReservedMatchedAmt4)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(5).String()], ReservedMatchedAmt5)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(6).String()], ReservedMatchedAmt6)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(7).String()], ReservedMatchedAmt7)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(8).String()], ReservedMatchedAmt8)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(9).String()], ReservedMatchedAmt9)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(10).String()], ReservedMatchedAmt10)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(11).String()], ReservedMatchedAmt11)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(12).String()], ReservedMatchedAmt12)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(13).String()], ReservedMatchedAmt13)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(14).String()], ReservedMatchedAmt14)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(15).String()], ReservedMatchedAmt15)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(16).String()], ReservedMatchedAmt16)
	s.Require().Equal(mInfo.ReservedMatchedMap[s.addr(17).String()], ReservedMatchedAmt17)

	RefundAmt1 := sdk.NewInt(2000_000_000)
	RefundAmt2 := sdk.NewInt(0)
	RefundAmt3 := sdk.NewInt(200_000_000)
	RefundAmt4 := sdk.NewInt(1640_000_000)
	RefundAmt5 := sdk.NewInt(1640_000_000)
	RefundAmt6 := sdk.NewInt(2120_000_000)
	RefundAmt7 := sdk.NewInt(120_000_000)
	RefundAmt8 := sdk.NewInt(1900_000_000)
	RefundAmt9 := sdk.NewInt(0)
	RefundAmt10 := sdk.NewInt(0)
	RefundAmt11 := sdk.NewInt(65_000_000)
	RefundAmt12 := sdk.NewInt(40_000_000)
	RefundAmt13 := sdk.NewInt(10_000_000)
	RefundAmt14 := sdk.NewInt(980_000_000)
	RefundAmt15 := sdk.NewInt(15_000_000)
	RefundAmt16 := sdk.NewInt(0)
	RefundAmt17 := sdk.NewInt(42_000_000)

	s.Require().Equal(mInfo.RefundMap[s.addr(1).String()].Abs(), RefundAmt1.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(2).String()].Abs(), RefundAmt2.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(3).String()].Abs(), RefundAmt3.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(4).String()].Abs(), RefundAmt4.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(5).String()].Abs(), RefundAmt5.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(6).String()].Abs(), RefundAmt6.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(7).String()].Abs(), RefundAmt7.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(8).String()].Abs(), RefundAmt8.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(9).String()].Abs(), RefundAmt9.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(10).String()].Abs(), RefundAmt10.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(11).String()].Abs(), RefundAmt11.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(12).String()].Abs(), RefundAmt12.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(13).String()].Abs(), RefundAmt13.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(14).String()].Abs(), RefundAmt14.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(15).String()].Abs(), RefundAmt15.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(16).String()].Abs(), RefundAmt16.Abs())
	s.Require().Equal(mInfo.RefundMap[s.addr(17).String()].Abs(), RefundAmt17.Abs())

	// Distribute selling coin
	err := s.keeper.AllocateSellingCoin(s.ctx, auction, mInfo)
	s.Require().NoError(err)

	err = s.keeper.ReleaseRemainingSellingCoin(s.ctx, auction)
	s.Require().NoError(err)

	// The selling reserve account balance must be zero
	s.Require().True(s.getBalance(auction.GetSellingReserveAddress(), auction.SellingCoin.Denom).IsZero())

	// The auctioneer must have sellingCoin.Amount - TotalMatchedAmount
	s.Require().Equal(s.getBalance(s.addr(0), auction.GetSellingCoin().Denom).Amount, auction.SellingCoin.Amount.Sub(mInfo.TotalMatchedAmount))

	// The bidders must have the matched selling coin
	s.Require().Equal(s.getBalance(s.addr(1), auction.GetSellingCoin().Denom).Amount.Abs(), MatchedAmt1.Abs())
	s.Require().Equal(s.getBalance(s.addr(2), auction.GetSellingCoin().Denom).Amount.Abs(), MatchedAmt2.Abs())
	s.Require().Equal(s.getBalance(s.addr(3), auction.GetSellingCoin().Denom).Amount.Abs(), MatchedAmt3.Abs())
	s.Require().Equal(s.getBalance(s.addr(4), auction.GetSellingCoin().Denom).Amount.Abs(), MatchedAmt4.Abs())
	s.Require().Equal(s.getBalance(s.addr(5), auction.GetSellingCoin().Denom).Amount.Abs(), MatchedAmt5.Abs())
	s.Require().Equal(s.getBalance(s.addr(6), auction.GetSellingCoin().Denom).Amount.Abs(), MatchedAmt6.Abs())
	s.Require().Equal(s.getBalance(s.addr(7), auction.GetSellingCoin().Denom).Amount.Abs(), MatchedAmt7.Abs())
	s.Require().Equal(s.getBalance(s.addr(8), auction.GetSellingCoin().Denom).Amount.Abs(), MatchedAmt8.Abs())
	s.Require().Equal(s.getBalance(s.addr(9), auction.GetSellingCoin().Denom).Amount.Abs(), MatchedAmt9.Abs())
	s.Require().Equal(s.getBalance(s.addr(10), auction.GetSellingCoin().Denom).Amount.Abs(), MatchedAmt10.Abs())
	s.Require().Equal(s.getBalance(s.addr(11), auction.GetSellingCoin().Denom).Amount.Abs(), MatchedAmt11.Abs())
	s.Require().Equal(s.getBalance(s.addr(12), auction.GetSellingCoin().Denom).Amount.Abs(), MatchedAmt12.Abs())
	s.Require().Equal(s.getBalance(s.addr(13), auction.GetSellingCoin().Denom).Amount.Abs(), MatchedAmt13.Abs())
	s.Require().Equal(s.getBalance(s.addr(14), auction.GetSellingCoin().Denom).Amount.Abs(), MatchedAmt14.Abs())
	s.Require().Equal(s.getBalance(s.addr(15), auction.GetSellingCoin().Denom).Amount.Abs(), MatchedAmt15.Abs())
	s.Require().Equal(s.getBalance(s.addr(16), auction.GetSellingCoin().Denom).Amount.Abs(), MatchedAmt16.Abs())
	s.Require().Equal(s.getBalance(s.addr(17), auction.GetSellingCoin().Denom).Amount.Abs(), MatchedAmt17.Abs())

	// Refund payingCoin
	err = s.keeper.RefundPayingCoin(s.ctx, auction, mInfo)
	s.Require().NoError(err)
}
