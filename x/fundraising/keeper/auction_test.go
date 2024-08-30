package keeper_test

import (
	"time"

	"cosmossdk.io/collections"
	sdkerrors "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	_ "github.com/stretchr/testify/suite"

	"github.com/tendermint/fundraising/x/fundraising/types"
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

	auction, err := s.keeper.Auction.Get(s.ctx, standByAuction.GetId())
	s.Require().NoError(err)
	s.Require().Equal(types.AuctionStatusStandBy, auction.GetStatus())

	feePool, err := s.app.DistrKeeper.FeePool.Get(s.ctx)
	s.Require().NoError(err)

	params, err := s.keeper.Params.Get(s.ctx)
	s.Require().NoError(err)

	auctionCreationFee := params.AuctionCreationFee
	s.Require().True(feePool.CommunityPool.Equal(sdk.NewDecCoinsFromCoins(auctionCreationFee...)))

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

	auction, err = s.keeper.Auction.Get(s.ctx, startedAuction.GetId())
	s.Require().NoError(err)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())
}

func (s *KeeperTestSuite) TestFixedPriceAuction_BidWithPayingCoinDenom() {
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

	auction, err := s.keeper.Auction.Get(s.ctx, startedAuction.GetId())
	s.Require().NoError(err)

	s.placeBidFixedPrice(auction.GetId(), s.addr(1), auction.GetStartPrice(), parseCoin("1_000_000denom2"), true)
	s.placeBidFixedPrice(auction.GetId(), s.addr(1), auction.GetStartPrice(), parseCoin("1_000_000denom2"), true)
	s.placeBidFixedPrice(auction.GetId(), s.addr(1), auction.GetStartPrice(), parseCoin("1_000_000denom2"), true)

	// Make sure allocate amount is equal to the total bid amount made by the same bidder
	mInfo, err := s.keeper.CalculateFixedPriceAllocation(s.ctx, auction)
	s.Require().NoError(err)

	allocateAmt := mInfo.AllocationMap[s.addr(1).String()]
	s.Require().Equal(allocateAmt, parseCoin("6_000_000denom1").Amount)
}

func (s *KeeperTestSuite) TestFixedPriceAuction_BidWithSellingCoinDenom() {
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

	auction, err := s.keeper.Auction.Get(s.ctx, startedAuction.GetId())
	s.Require().NoError(err)

	s.placeBidFixedPrice(auction.GetId(), s.addr(1), auction.GetStartPrice(), parseCoin("1_000_000denom1"), true)
	s.placeBidFixedPrice(auction.GetId(), s.addr(1), auction.GetStartPrice(), parseCoin("1_000_000denom1"), true)
	s.placeBidFixedPrice(auction.GetId(), s.addr(1), auction.GetStartPrice(), parseCoin("1_000_000denom1"), true)

	// Make sure allocate amount is equal to the total bid amount made by the same bidder
	mInfo, err := s.keeper.CalculateFixedPriceAllocation(s.ctx, auction)
	s.Require().NoError(err)

	allocateAmt := mInfo.AllocationMap[s.addr(1).String()]
	s.Require().Equal(allocateAmt, parseCoin("3_000_000denom1").Amount)
}

func (s *KeeperTestSuite) TestFixedPriceAuction_BidWithBoth() {
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

	auction, err := s.keeper.Auction.Get(s.ctx, startedAuction.GetId())
	s.Require().NoError(err)

	s.placeBidFixedPrice(auction.GetId(), s.addr(1), auction.GetStartPrice(), parseCoin("2_000_000denom1"), true)
	s.placeBidFixedPrice(auction.GetId(), s.addr(1), auction.GetStartPrice(), parseCoin("2_000_000denom1"), true)
	s.placeBidFixedPrice(auction.GetId(), s.addr(1), auction.GetStartPrice(), parseCoin("2_000_000denom2"), true)
	s.placeBidFixedPrice(auction.GetId(), s.addr(2), auction.GetStartPrice(), parseCoin("1_000_000denom1"), true)
	s.placeBidFixedPrice(auction.GetId(), s.addr(2), auction.GetStartPrice(), parseCoin("1_000_000denom2"), true)

	// Make sure allocate amount is equal to the total bid amount made by the same bidder
	mInfo, err := s.keeper.CalculateFixedPriceAllocation(s.ctx, auction)
	s.Require().NoError(err)

	s.Require().Equal(mInfo.AllocationMap[s.addr(1).String()], parseCoin("8_000_000denom2").Amount)
	s.Require().Equal(mInfo.AllocationMap[s.addr(2).String()], parseCoin("3_000_000denom2").Amount)
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

	_, err := s.keeper.Auction.Get(s.ctx, auction.GetId())
	s.Require().NoError(err)

	// Place bids
	s.placeBidFixedPrice(auction.Id, s.addr(1), parseDec("0.5"), parseCoin("100_000_000denom2"), true)
	s.placeBidFixedPrice(auction.Id, s.addr(2), parseDec("0.5"), parseCoin("100_000_000denom1"), true)
	s.placeBidFixedPrice(auction.Id, s.addr(3), parseDec("0.5"), parseCoin("200_000_000denom1"), true)

	// Calculate allocation
	mInfo, err := s.keeper.CalculateFixedPriceAllocation(s.ctx, auction)
	s.Require().NoError(err)

	// Distribute selling coin
	err = s.keeper.AllocateSellingCoin(s.ctx, auction, mInfo)
	s.Require().NoError(err)

	err = s.keeper.RefundRemainingSellingCoin(s.ctx, auction)
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

func (s *KeeperTestSuite) TestFixedPriceAuction_ReleaseVestingPayingCoin() {
	auction := s.createFixedPriceAuction(
		s.addr(0),
		parseDec("1"),
		parseCoin("1000_000_000_000denom1"),
		"denom2",
		[]types.VestingSchedule{
			{
				ReleaseTime: time.Now().AddDate(0, 6, 0),
				Weight:      math.LegacyMustNewDecFromStr("0.25"),
			},
			{
				ReleaseTime: time.Now().AddDate(0, 9, 0),
				Weight:      math.LegacyMustNewDecFromStr("0.25"),
			},
			{
				ReleaseTime: time.Now().AddDate(1, 0, 0),
				Weight:      math.LegacyMustNewDecFromStr("0.25"),
			},
			{
				ReleaseTime: time.Now().AddDate(1, 3, 0),
				Weight:      math.LegacyMustNewDecFromStr("0.25"),
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
	mInfo, err := s.keeper.CalculateFixedPriceAllocation(s.ctx, auction)
	s.Require().NoError(err)

	// Distribute selling coin
	err = s.keeper.AllocateSellingCoin(s.ctx, auction, mInfo)
	s.Require().NoError(err)

	// Apply vesting schedules
	err = s.keeper.ApplyVestingSchedules(s.ctx, auction)
	s.Require().NoError(err)

	// All of the vesting queues must not be released yet
	vqs, err := s.keeper.GetVestingQueuesByAuctionId(s.ctx, auction.GetId())
	s.Require().NoError(err)

	s.Require().Equal(4, len(vqs))
	for _, vq := range vqs {
		s.Require().False(vq.Released)
	}

	// Change the block time to release two vesting schedules
	s.ctx = s.ctx.WithBlockTime(vqs[0].GetReleaseTime().AddDate(0, 4, 1))
	err = s.keeper.BeginBlocker(s.ctx)
	s.Require().NoError(err)

	// Distribute paying coin
	err = s.keeper.ReleaseVestingPayingCoin(s.ctx, auction)
	s.Require().NoError(err)

	// First two vesting queues must be released
	vqsa, err := s.keeper.GetVestingQueuesByAuctionId(s.ctx, auction.GetId())
	s.Require().NoError(err)

	for i, vq := range vqsa {
		if i == 0 || i == 1 {
			s.Require().True(vq.Released)
		} else {
			s.Require().False(vq.Released)
		}
	}

	// Change the block time
	s.ctx = s.ctx.WithBlockTime(vqs[3].GetReleaseTime().AddDate(0, 0, 1))
	err = s.keeper.BeginBlocker(s.ctx)
	s.Require().NoError(err)
	s.Require().NoError(s.keeper.ReleaseVestingPayingCoin(s.ctx, auction))

	// All of the vesting queues must be released
	vqByAuction, err := s.keeper.GetVestingQueuesByAuctionId(s.ctx, auction.GetId())
	s.Require().NoError(err)
	for _, vq := range vqByAuction {
		s.Require().True(vq.Released)
	}

	finishedAuction, err := s.keeper.Auction.Get(s.ctx, auction.GetId())
	s.Require().NoError(err)
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
	s.Require().Error(err, errors.ErrNotFound)

	// Unauthorized
	err = s.keeper.CancelAuction(s.ctx, &types.MsgCancelAuction{
		Auctioneer: s.addr(10).String(),
		AuctionId:  auction.Id,
	})
	s.Require().Error(err, errors.ErrUnauthorized)

	// Invalid auction status
	err = s.keeper.CancelAuction(s.ctx, &types.MsgCancelAuction{
		Auctioneer: auction.Auctioneer,
		AuctionId:  auction.Id,
	})
	s.Require().Error(err, types.ErrInvalidAuctionStatus)

	// Forcefully update auction status
	err = auction.SetStatus(types.AuctionStatusStandBy)
	s.Require().NoError(err)
	err = s.keeper.Auction.Set(s.ctx, auction.GetId(), auction)
	s.Require().NoError(err)

	// Cancel the auction
	err = s.keeper.CancelAuction(s.ctx, &types.MsgCancelAuction{
		Auctioneer: auction.Auctioneer,
		AuctionId:  auction.Id,
	})
	s.Require().NoError(err)

	// Verify the status
	a, err := s.keeper.Auction.Get(s.ctx, auction.GetId())
	s.Require().NoError(err)
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
		math.LegacyMustNewDecFromStr("0.2"),
		time.Now().AddDate(0, 6, 0),
		time.Now().AddDate(0, 6, 0).AddDate(0, 1, 0),
		true,
	)

	auction, err := s.keeper.Auction.Get(s.ctx, standByAuction.GetId())
	s.Require().NoError(err)
	s.Require().Equal(types.AuctionStatusStandBy, auction.GetStatus())

	feePool, err := s.app.DistrKeeper.FeePool.Get(s.ctx)
	s.Require().NoError(err)

	params, err := s.keeper.Params.Get(s.ctx)
	s.Require().NoError(err)
	auctionCreationFee := params.AuctionCreationFee
	s.Require().True(feePool.CommunityPool.Equal(sdk.NewDecCoinsFromCoins(auctionCreationFee...)))

	startedAuction := s.createBatchAuction(
		s.addr(1),
		parseDec("0.5"),
		parseDec("0.1"),
		parseCoin("5000_000_000denom3"),
		"denom4",
		[]types.VestingSchedule{},
		1,
		math.LegacyMustNewDecFromStr("0.2"),
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)

	auction, err = s.keeper.Auction.Get(s.ctx, startedAuction.GetId())
	s.Require().NoError(err)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())
}

func (s *KeeperTestSuite) TestBatchAuction_MaxNumVestingSchedules() {
	batchAuction := types.NewMsgCreateBatchAuction(
		s.addr(0).String(),
		parseDec("1"),
		parseDec("0.1"),
		parseCoin("5000_000_000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		1,
		math.LegacyMustNewDecFromStr("0.2"),
		time.Now().AddDate(0, 6, 0),
		time.Now().AddDate(0, 6, 0).AddDate(0, 1, 0),
	)

	params, err := s.keeper.Params.Get(s.ctx)
	s.Require().NoError(err)

	s.fundAddr(s.addr(0), params.AuctionCreationFee.Add(batchAuction.SellingCoin))

	// Invalid max extended round
	batchAuction.MaxExtendedRound = types.MaxExtendedRound + 1

	_, err = s.keeper.CreateBatchAuction(s.ctx, batchAuction)
	s.Require().EqualError(err, sdkerrors.Wrap(errors.ErrInvalidRequest, "exceed maximum extended round").Error())

	batchAuction.MaxExtendedRound = 1

	// Invalid number of vesting schedules
	numSchedules := types.MaxNumVestingSchedules + 1
	schedules := make([]types.VestingSchedule, numSchedules)
	totalWeight := math.LegacyZeroDec()
	for i := range schedules {
		var schedule types.VestingSchedule
		if i == numSchedules-1 {
			schedule.Weight = math.LegacyOneDec().Sub(totalWeight)
		} else {
			schedule.Weight = math.LegacyOneDec().Quo(math.LegacyNewDec(int64(numSchedules)))
		}
		schedule.ReleaseTime = time.Now().AddDate(0, 0, i)

		totalWeight = totalWeight.Add(schedule.Weight)
		schedules[i] = schedule
	}
	batchAuction.VestingSchedules = schedules

	_, err = s.keeper.CreateBatchAuction(s.ctx, batchAuction)
	s.Require().EqualError(err, sdkerrors.Wrap(errors.ErrInvalidRequest, "exceed maximum number of vesting schedules").Error())
}

func (s *KeeperTestSuite) TestFixedPriceAuction_MaxNumVestingSchedules() {
	fixedPriceAuction := types.NewMsgCreateFixedPriceAuction(
		s.addr(0).String(),
		parseDec("0.5"),
		parseCoin("500_000_000_000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
	)

	params, err := s.keeper.Params.Get(s.ctx)
	s.Require().NoError(err)
	s.fundAddr(s.addr(0), params.AuctionCreationFee.Add(fixedPriceAuction.SellingCoin))

	// Invalid number of vesting schedules
	numSchedules := types.MaxNumVestingSchedules + 1
	schedules := make([]types.VestingSchedule, numSchedules)
	totalWeight := math.LegacyZeroDec()
	for i := range schedules {
		var schedule types.VestingSchedule
		if i == numSchedules-1 {
			schedule.Weight = math.LegacyOneDec().Sub(totalWeight)
		} else {
			schedule.Weight = math.LegacyOneDec().Quo(math.LegacyNewDec(int64(numSchedules)))
		}
		schedule.ReleaseTime = time.Now().AddDate(0, 0, i)

		totalWeight = totalWeight.Add(schedule.Weight)
		schedules[i] = schedule
	}
	fixedPriceAuction.VestingSchedules = schedules

	_, err = s.keeper.CreateFixedPriceAuction(s.ctx, fixedPriceAuction)
	s.Require().EqualError(err, sdkerrors.Wrap(errors.ErrInvalidRequest, "exceed maximum number of vesting schedules").Error())
}

func (s *KeeperTestSuite) TestInvalidEndTime() {
	params, err := s.keeper.Params.Get(s.ctx)
	s.Require().NoError(err)

	fixedPriceAuction := types.NewMsgCreateFixedPriceAuction(
		s.addr(0).String(),
		parseDec("0.5"),
		parseCoin("500_000_000_000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		types.MustParseRFC3339("2022-03-01T00:00:00Z"),
		types.MustParseRFC3339("2022-01-01T00:00:00Z"),
	)
	s.fundAddr(s.addr(0), params.AuctionCreationFee.Add(fixedPriceAuction.SellingCoin))

	_, err = s.keeper.CreateFixedPriceAuction(s.ctx, fixedPriceAuction)
	s.Require().EqualError(err, sdkerrors.Wrap(errors.ErrInvalidRequest, "end time must be set after the current time").Error())

	batchAuction := types.NewMsgCreateBatchAuction(
		s.addr(1).String(),
		parseDec("1"),
		parseDec("0.1"),
		parseCoin("5000_000_000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		1,
		math.LegacyMustNewDecFromStr("0.2"),
		types.MustParseRFC3339("2022-03-01T00:00:00Z"),
		types.MustParseRFC3339("2022-01-01T00:00:00Z"),
	)
	s.fundAddr(s.addr(1), params.AuctionCreationFee.Add(batchAuction.SellingCoin))

	_, err = s.keeper.CreateBatchAuction(s.ctx, batchAuction)
	s.Require().EqualError(err, sdkerrors.Wrap(errors.ErrInvalidRequest, "end time must be set after the current time").Error())
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

	auction, err := s.keeper.Auction.Get(s.ctx, startedAuction.GetId())
	s.Require().NoError(err)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	bidderByAuction, err := s.keeper.GetAllowedBiddersByAuction(s.ctx, startedAuction.Id)
	s.Require().Len(bidderByAuction, 0)

	// Invalid auction id
	err = s.keeper.AddAllowedBidders(s.ctx, 10, []types.AllowedBidder{
		{AuctionId: 1, Bidder: s.addr(1).String(), MaxBidAmount: math.NewInt(100_000_000)},
	})
	s.Require().Error(err)

	for _, tc := range []struct {
		name        string
		bidders     []types.AllowedBidder
		expectedErr error
	}{
		{
			"single bidder",
			[]types.AllowedBidder{
				{AuctionId: 1, Bidder: s.addr(1).String(), MaxBidAmount: math.NewInt(100_000_000)},
			},
			nil,
		},
		{
			"multiple bidders",
			[]types.AllowedBidder{
				{AuctionId: 1, Bidder: s.addr(1).String(), MaxBidAmount: math.NewInt(100_000_000)},
				{AuctionId: 1, Bidder: s.addr(2).String(), MaxBidAmount: math.NewInt(500_000_000)},
				{AuctionId: 1, Bidder: s.addr(3).String(), MaxBidAmount: math.NewInt(800_000_000)},
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
				{AuctionId: 1, Bidder: s.addr(1).String(), MaxBidAmount: math.NewInt(0)},
			},
			types.ErrInvalidMaxBidAmount,
		},
		{
			"negative maximum bid amount",
			[]types.AllowedBidder{
				{AuctionId: 1, Bidder: s.addr(1).String(), MaxBidAmount: math.NewInt(-1)},
			},
			types.ErrInvalidMaxBidAmount,
		},
		{
			"exceed the total selling amount",
			[]types.AllowedBidder{
				{AuctionId: 1, Bidder: s.addr(1).String(), MaxBidAmount: math.NewInt(500_000_000_001)},
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

	auction, err := s.keeper.Auction.Get(s.ctx, startedAuction.GetId())
	s.Require().NoError(err)
	bidderByAuction, err := s.keeper.GetAllowedBiddersByAuction(s.ctx, auction.GetId())
	s.Require().NoError(err)
	s.Require().Len(bidderByAuction, 0)

	// Add some bidders
	s.Require().NoError(s.keeper.AddAllowedBidders(s.ctx, auction.GetId(), []types.AllowedBidder{
		{AuctionId: 1, Bidder: s.addr(1).String(), MaxBidAmount: math.NewInt(100_000_000)},
		{AuctionId: 1, Bidder: s.addr(2).String(), MaxBidAmount: math.NewInt(500_000_000)},
	}))

	auction, err = s.keeper.Auction.Get(s.ctx, auction.GetId())
	s.Require().NoError(err)

	bidderByAuction, err = s.keeper.GetAllowedBiddersByAuction(s.ctx, auction.GetId())
	s.Require().NoError(err)
	s.Require().Len(bidderByAuction, 2)

	// Add more bidders
	s.Require().NoError(s.keeper.AddAllowedBidders(s.ctx, auction.GetId(), []types.AllowedBidder{
		{AuctionId: 1, Bidder: s.addr(3).String(), MaxBidAmount: math.NewInt(100_000_000)},
		{AuctionId: 1, Bidder: s.addr(4).String(), MaxBidAmount: math.NewInt(100_000_000)},
		{AuctionId: 1, Bidder: s.addr(5).String(), MaxBidAmount: math.NewInt(100_000_000)},
	}))

	auction, err = s.keeper.Auction.Get(s.ctx, auction.GetId())
	s.Require().NoError(err)
	bidderByAuction, err = s.keeper.GetAllowedBiddersByAuction(s.ctx, auction.GetId())
	s.Require().NoError(err)
	s.Require().Len(bidderByAuction, 5)
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

	auction, err := s.keeper.Auction.Get(s.ctx, startedAuction.GetId())
	s.Require().NoError(err)
	bidderByAuction, err := s.keeper.GetAllowedBiddersByAuction(s.ctx, startedAuction.GetId())
	s.Require().NoError(err)
	s.Require().Len(bidderByAuction, 0)

	// Invalid auction id
	err = s.keeper.UpdateAllowedBidder(s.ctx, 10, s.addr(1), math.NewInt(100_000_000))
	s.Require().Error(err)

	// Add 5 bidders with different maximum bid amount
	s.Require().NoError(s.keeper.AddAllowedBidders(s.ctx, auction.GetId(), []types.AllowedBidder{
		{AuctionId: 1, Bidder: s.addr(1).String(), MaxBidAmount: math.NewInt(100_000_000)},
		{AuctionId: 1, Bidder: s.addr(2).String(), MaxBidAmount: math.NewInt(200_000_000)},
		{AuctionId: 1, Bidder: s.addr(3).String(), MaxBidAmount: math.NewInt(300_000_000)},
		{AuctionId: 1, Bidder: s.addr(4).String(), MaxBidAmount: math.NewInt(400_000_000)},
		{AuctionId: 1, Bidder: s.addr(5).String(), MaxBidAmount: math.NewInt(500_000_000)},
	}))
	bidderByAuction, err = s.keeper.GetAllowedBiddersByAuction(s.ctx, auction.GetId())
	s.Require().NoError(err)
	s.Require().Len(bidderByAuction, 5)

	for _, tc := range []struct {
		name         string
		bidder       sdk.AccAddress
		maxBidAmount math.Int
		expectedErr  error
	}{
		{
			"update bidder's maximum bid amount",
			s.addr(1),
			math.NewInt(555_000_000_000),
			nil,
		},
		{
			"bidder not found",
			s.addr(10),
			math.NewInt(300_000_000),
			sdkerrors.Wrapf(errors.ErrNotFound, "bidder %s is not found", s.addr(10).String()),
		},
		{
			"zero maximum bid amount value",
			s.addr(1),
			math.NewInt(0),
			types.ErrInvalidMaxBidAmount,
		},
		{
			"negative maximum bid amount value",
			s.addr(1),
			math.NewInt(-1),
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

			auction, err = s.keeper.Auction.Get(s.ctx, auction.GetId())
			s.Require().NoError(err)

			allowedBidders, err := s.keeper.GetAllowedBiddersByAuction(s.ctx, startedAuction.Id)
			s.Require().NoError(err)
			s.Require().Len(allowedBidders, 5)

			// Check if it is successfully updated
			allowedBidder, err := s.keeper.AllowedBidder.Get(s.ctx, collections.Join(auction.GetId(), tc.bidder))
			s.Require().NoError(err)
			s.Require().Equal(tc.maxBidAmount, allowedBidder.MaxBidAmount)
		})
	}
}

func (s *KeeperTestSuite) TestRefundPayingCoin() {
	auction := s.createBatchAuction(
		s.addr(0),
		parseDec("1.0"),
		parseDec("0.5"),
		parseCoin("100_000_000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		1,
		math.LegacyMustNewDecFromStr("0.2"),
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	s.placeBidBatchMany(auction.Id, s.addr(1), parseDec("1"), parseCoin("50_000_000denom1"), math.NewInt(1_000_000_000), true)
	refundBid := s.placeBidBatchMany(auction.Id, s.addr(2), parseDec("0.9"), parseCoin("60_000_000denom1"), math.NewInt(1_000_000_000), true)

	a, err := s.keeper.Auction.Get(s.ctx, auction.Id)
	s.Require().NoError(err)

	mInfo, err := s.keeper.CalculateBatchAllocation(s.ctx, a)
	s.Require().NoError(err)

	err = s.keeper.RefundPayingCoin(s.ctx, a, mInfo)
	s.Require().NoError(err)

	expectedAmt := refundBid.ConvertToPayingAmount(auction.GetPayingCoinDenom())
	bidderBalance := s.getBalance(s.addr(2), auction.GetPayingCoinDenom()).Amount
	s.Require().Equal(expectedAmt, bidderBalance)
}

func (s *KeeperTestSuite) TestCloseFixedPriceAuction() {
	auction := s.createFixedPriceAuction(
		s.addr(0),
		parseDec("1"),
		parseCoin("1_000_000_000_000denom1"),
		"denom2",
		[]types.VestingSchedule{
			{ReleaseTime: types.MustParseRFC3339("2023-01-01T00:00:00Z"), Weight: math.LegacyOneDec()},
		},
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	s.placeBidFixedPrice(auction.Id, s.addr(1), parseDec("1"), parseCoin("250_000_000denom1"), true)
	s.placeBidFixedPrice(auction.Id, s.addr(2), parseDec("1"), parseCoin("250_000_000denom1"), true)
	s.placeBidFixedPrice(auction.Id, s.addr(3), parseDec("1"), parseCoin("250_000_000denom1"), true)
	s.placeBidFixedPrice(auction.Id, s.addr(4), parseDec("1"), parseCoin("250_000_000denom1"), true)

	a, err := s.keeper.Auction.Get(s.ctx, auction.Id)
	s.Require().NoError(err)

	err = s.keeper.CloseFixedPriceAuction(s.ctx, a)
	s.Require().NoError(err)

	s.Require().Equal(parseCoin("999000000000denom1"), s.getBalance(s.addr(0), a.GetSellingCoin().Denom))
	s.Require().Equal(parseCoin("0denom2"), s.getBalance(s.addr(0), a.GetPayingCoinDenom()))
	s.Require().Equal(parseCoin("250_000_000denom1"), s.getBalance(s.addr(1), a.GetSellingCoin().Denom))
	s.Require().Equal(parseCoin("250_000_000denom1"), s.getBalance(s.addr(2), a.GetSellingCoin().Denom))
	s.Require().Equal(parseCoin("250_000_000denom1"), s.getBalance(s.addr(3), a.GetSellingCoin().Denom))
	s.Require().Equal(parseCoin("250_000_000denom1"), s.getBalance(s.addr(4), a.GetSellingCoin().Denom))

	vestingQueues, err := s.keeper.VestingQueues(s.ctx)
	s.Require().NoError(err)
	s.Require().Len(vestingQueues, len(a.GetVestingSchedules()))
}

func (s *KeeperTestSuite) TestCloseBatchAuction() {
	// Close a batch auction right away by setting MaxExtendedRound to 0 value
	maxExtendedRound := uint32(0)

	auction := s.createBatchAuction(
		s.addr(0),
		parseDec("0.5"),
		parseDec("0.1"),
		parseCoin("10_000_000_000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		maxExtendedRound,
		math.LegacyMustNewDecFromStr("0.2"),
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	s.placeBidBatchMany(auction.Id, s.addr(1), parseDec("0.9"), parseCoin("200_000_000denom1"), math.NewInt(1_000_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(2), parseDec("0.8"), parseCoin("200_000_000denom1"), math.NewInt(1_000_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(3), parseDec("0.7"), parseCoin("100_000_000denom1"), math.NewInt(1_000_000_000), true)

	a, err := s.keeper.Auction.Get(s.ctx, auction.Id)
	s.Require().NoError(err)

	err = s.keeper.CloseBatchAuction(s.ctx, a)
	s.Require().NoError(err)

	s.Require().Equal(parseCoin("350000000denom2"), s.getBalance(s.addr(0), a.GetPayingCoinDenom()))
	s.Require().Equal(parseCoin("9500000000denom1"), s.getBalance(s.addr(0), a.GetSellingCoin().Denom))
	s.Require().Equal(parseCoin("200_000_000denom1"), s.getBalance(s.addr(1), a.GetSellingCoin().Denom))
	s.Require().Equal(parseCoin("200_000_000denom1"), s.getBalance(s.addr(2), a.GetSellingCoin().Denom))
	s.Require().Equal(parseCoin("100_000_000denom1"), s.getBalance(s.addr(3), a.GetSellingCoin().Denom))

	vestingQueues, err := s.keeper.VestingQueues(s.ctx)
	s.Require().NoError(err)
	s.Require().Len(vestingQueues, len(a.GetVestingSchedules()))
}

func (s *KeeperTestSuite) TestCloseBatchAuction_ExtendRound() {
	// Extend round for a batch auction by setting MaxExtendedRound to non zero value
	maxExtendedRound := uint32(5)
	extendedRoundRate := parseDec("0.2")

	auction := s.createBatchAuction(
		s.addr(0),
		parseDec("0.5"),
		parseDec("0.1"),
		parseCoin("10_000_000_000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		maxExtendedRound,
		extendedRoundRate,
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	s.placeBidBatchMany(auction.Id, s.addr(1), parseDec("0.9"), parseCoin("200_000_000denom1"), math.NewInt(1_000_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(2), parseDec("0.8"), parseCoin("200_000_000denom1"), math.NewInt(1_000_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(3), parseDec("0.7"), parseCoin("100_000_000denom1"), math.NewInt(1_000_000_000), true)

	a, err := s.keeper.Auction.Get(s.ctx, auction.Id)
	s.Require().NoError(err)
	s.Require().Len(a.GetEndTimes(), 1)

	s.keeper.CloseBatchAuction(s.ctx, auction)

	// Extended round must be triggered
	a, err = s.keeper.Auction.Get(s.ctx, auction.Id)
	s.Require().NoError(err)
	s.Require().Len(a.GetEndTimes(), 2)

	// Auction sniping occurs
	s.placeBidBatchMany(auction.Id, s.addr(4), parseDec("0.85"), parseCoin("9_800_000_000denom1"), math.NewInt(100_000_000_000), true)

	s.keeper.CloseBatchAuction(s.ctx, a)

	a, err = s.keeper.Auction.Get(s.ctx, auction.Id)
	s.Require().NoError(err)
	s.Require().Len(a.GetEndTimes(), 3)
}

func (s *KeeperTestSuite) TestCloseBatchAuction_Valid() {
	// Extend round for a batch auction by setting MaxExtendedRound to non zero value
	maxExtendedRound := uint32(5)
	extendedRoundRate := parseDec("0.2")

	auction := s.createBatchAuction(
		s.addr(0),
		parseDec("0.5"),
		parseDec("0.1"),
		parseCoin("10_000_000_000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		maxExtendedRound,
		extendedRoundRate,
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)
	s.Require().Equal(types.AuctionStatusStarted, auction.GetStatus())

	s.placeBidBatchMany(auction.Id, s.addr(1), parseDec("0.9"), parseCoin("200_000_000denom1"), math.NewInt(1_000_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(2), parseDec("0.8"), parseCoin("200_000_000denom1"), math.NewInt(1_000_000_000), true)
	s.placeBidBatchMany(auction.Id, s.addr(3), parseDec("0.7"), parseCoin("100_000_000denom1"), math.NewInt(1_000_000_000), true)

	a, err := s.keeper.Auction.Get(s.ctx, auction.Id)
	s.Require().NoError(err)
	s.Require().Len(a.GetEndTimes(), 1)

	err = s.keeper.CloseBatchAuction(s.ctx, auction)
	s.Require().NoError(err)

	// Extended round must be triggered
	a, err = s.keeper.Auction.Get(s.ctx, auction.Id)
	s.Require().NoError(err)
	s.Require().Len(a.GetEndTimes(), 2)

	// Auction sniping occurs
	s.placeBidBatchMany(auction.Id, s.addr(4), parseDec("0.85"), parseCoin("9_500_000_000denom1"), math.NewInt(100_000_000_000), true)

	err = s.keeper.CloseBatchAuction(s.ctx, a)
	s.Require().NoError(err)

	a, err = s.keeper.Auction.Get(s.ctx, auction.Id)
	s.Require().NoError(err)
	s.Require().Len(a.GetEndTimes(), 2)
}
