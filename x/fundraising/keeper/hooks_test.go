package keeper_test

import (
	"context"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/stretchr/testify/suite"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

var _ types.FundraisingHooks = &MockFundraisingHooksReceiver{}

// MockFundraisingHooksReceiver event hooks for governance proposal object (noalias)
type MockFundraisingHooksReceiver struct {
	BeforeFixedPriceAuctionCreatedValid bool
	AfterFixedPriceAuctionCreatedValid  bool
	BeforeBatchAuctionCreatedValid      bool
	AfterBatchAuctionCreatedValid       bool
	BeforeAuctionCanceledValid          bool
	BeforeBidPlacedValid                bool
	BeforeBidModifiedValid              bool
	BeforeAllowedBiddersAddedValid      bool
	BeforeAllowedBidderUpdatedValid     bool
	BeforeSellingCoinsAllocatedValid    bool
}

func (h *MockFundraisingHooksReceiver) BeforeFixedPriceAuctionCreated(
	ctx context.Context,
	auctioneer string,
	startPrice math.LegacyDec,
	sellingCoin sdk.Coin,
	payingCoinDenom string,
	vestingSchedules []types.VestingSchedule,
	startTime time.Time,
	endTime time.Time,
) error {
	h.BeforeFixedPriceAuctionCreatedValid = true
	return nil
}

func (h *MockFundraisingHooksReceiver) AfterFixedPriceAuctionCreated(
	ctx context.Context,
	auctionId uint64,
	auctioneer string,
	startPrice math.LegacyDec,
	sellingCoin sdk.Coin,
	payingCoinDenom string,
	vestingSchedules []types.VestingSchedule,
	startTime time.Time,
	endTime time.Time,
) error {
	h.AfterFixedPriceAuctionCreatedValid = true
	return nil
}

func (h *MockFundraisingHooksReceiver) BeforeBatchAuctionCreated(
	ctx context.Context,
	auctioneer string,
	startPrice math.LegacyDec,
	minBidPrice math.LegacyDec,
	sellingCoin sdk.Coin,
	payingCoinDenom string,
	vestingSchedules []types.VestingSchedule,
	maxExtendedRound uint32,
	extendedRoundRate math.LegacyDec,
	startTime time.Time,
	endTime time.Time,
) error {
	h.BeforeBatchAuctionCreatedValid = true
	return nil
}

func (h *MockFundraisingHooksReceiver) AfterBatchAuctionCreated(
	ctx context.Context,
	auctionId uint64,
	auctioneer string,
	startPrice math.LegacyDec,
	minBidPrice math.LegacyDec,
	sellingCoin sdk.Coin,
	payingCoinDenom string,
	vestingSchedules []types.VestingSchedule,
	maxExtendedRound uint32,
	extendedRoundRate math.LegacyDec,
	startTime time.Time,
	endTime time.Time,
) error {
	h.AfterBatchAuctionCreatedValid = true
	return nil
}

func (h *MockFundraisingHooksReceiver) BeforeAuctionCanceled(
	ctx context.Context,
	auctionId uint64,
	auctioneer string,
) error {
	h.BeforeAuctionCanceledValid = true
	return nil
}

func (h *MockFundraisingHooksReceiver) BeforeBidPlaced(
	ctx context.Context,
	auctionId uint64,
	bidId uint64,
	bidder string,
	bidType types.BidType,
	price math.LegacyDec,
	coin sdk.Coin,
) error {
	h.BeforeBidPlacedValid = true
	return nil
}

func (h *MockFundraisingHooksReceiver) BeforeBidModified(
	ctx context.Context,
	auctionId uint64,
	bidId uint64,
	bidder string,
	bidType types.BidType,
	price math.LegacyDec,
	coin sdk.Coin,
) error {
	h.BeforeBidModifiedValid = true
	return nil
}

func (h *MockFundraisingHooksReceiver) BeforeAllowedBiddersAdded(
	ctx context.Context,
	allowedBidders []types.AllowedBidder,
) error {
	h.BeforeAllowedBiddersAddedValid = true
	return nil
}

func (h *MockFundraisingHooksReceiver) BeforeAllowedBidderUpdated(
	ctx context.Context,
	auctionId uint64,
	bidder sdk.AccAddress,
	maxBidAmount math.Int,
) error {
	h.BeforeAllowedBidderUpdatedValid = true
	return nil
}

func (h *MockFundraisingHooksReceiver) BeforeSellingCoinsAllocated(
	ctx context.Context,
	auctionId uint64,
	allocationMap map[string]math.Int,
	refundMap map[string]math.Int,
) error {
	h.BeforeSellingCoinsAllocatedValid = true
	return nil
}

func (s *KeeperTestSuite) TestHooks() {
	fundraisingHooksReceiver := MockFundraisingHooksReceiver{}

	// Set hooks
	s.keeper.SetHooks(types.NewMultiFundraisingHooks(&fundraisingHooksReceiver))

	s.Require().False(fundraisingHooksReceiver.BeforeFixedPriceAuctionCreatedValid)
	s.Require().False(fundraisingHooksReceiver.AfterFixedPriceAuctionCreatedValid)
	s.Require().False(fundraisingHooksReceiver.BeforeBatchAuctionCreatedValid)
	s.Require().False(fundraisingHooksReceiver.AfterBatchAuctionCreatedValid)
	s.Require().False(fundraisingHooksReceiver.BeforeAuctionCanceledValid)
	s.Require().False(fundraisingHooksReceiver.BeforeBidPlacedValid)
	s.Require().False(fundraisingHooksReceiver.BeforeBidModifiedValid)
	s.Require().False(fundraisingHooksReceiver.BeforeAllowedBiddersAddedValid)
	s.Require().False(fundraisingHooksReceiver.BeforeAllowedBidderUpdatedValid)
	s.Require().False(fundraisingHooksReceiver.BeforeSellingCoinsAllocatedValid)

	// Create a fixed price auction
	s.createFixedPriceAuction(
		s.addr(0),
		parseDec("2.0"),
		parseCoin("1_000_000_000_000denom1"),
		"denom2",
		[]types.VestingSchedule{},
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)
	s.Require().True(fundraisingHooksReceiver.BeforeFixedPriceAuctionCreatedValid)
	s.Require().True(fundraisingHooksReceiver.AfterFixedPriceAuctionCreatedValid)

	// Create a batch auction
	batchAuction := s.createBatchAuction(
		s.addr(1),
		parseDec("0.5"),
		parseDec("0.1"),
		parseCoin("1_000_000_000_000denom3"),
		"denom4",
		[]types.VestingSchedule{},
		1,
		math.LegacyMustNewDecFromStr("0.2"),
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)
	s.Require().True(fundraisingHooksReceiver.BeforeBatchAuctionCreatedValid)
	s.Require().True(fundraisingHooksReceiver.AfterBatchAuctionCreatedValid)

	// Create auction that is stand by status
	standByAuction := s.createFixedPriceAuction(
		s.addr(2),
		parseDec("2.0"),
		parseCoin("1_000_000_000_000denom5"),
		"denom6",
		[]types.VestingSchedule{},
		time.Now().AddDate(0, 1, 0),
		time.Now().AddDate(0, 3, 0),
		true,
	)

	// Cancel the auction
	err := s.keeper.CancelAuction(s.ctx, &types.MsgCancelAuction{
		Auctioneer: standByAuction.Auctioneer,
		AuctionId:  standByAuction.Id,
	})
	s.Require().NoError(err)
	s.Require().True(fundraisingHooksReceiver.BeforeAuctionCanceledValid)

	// Get already started batch auction
	auction, err := s.keeper.Auction.Get(s.ctx, batchAuction.Id)
	s.Require().NoError(err)

	// Add allowed bidder
	allowedBidders := []types.AllowedBidder{types.NewAllowedBidder(auction.GetId(), s.addr(3), parseInt("100_000_000_000"))}
	s.Require().NoError(s.keeper.AddAllowedBidders(s.ctx, auction.GetId(), allowedBidders))
	s.Require().True(fundraisingHooksReceiver.BeforeAllowedBiddersAddedValid)

	// Update the allowed bidder
	err = s.keeper.UpdateAllowedBidder(s.ctx, auction.GetId(), s.addr(3), parseInt("110_000_000_000"))
	s.Require().NoError(err)
	s.Require().True(fundraisingHooksReceiver.BeforeAllowedBidderUpdatedValid)

	// Place a bid
	bid := s.placeBidBatchWorth(auction.GetId(), s.addr(3), parseDec("0.55"), parseCoin("5_000_000denom4"), math.NewInt(10_000_000), true)
	s.Require().True(fundraisingHooksReceiver.BeforeBidPlacedValid)

	// Modify the bid
	s.fundAddr(bid.GetBidder(), sdk.NewCoins(parseCoin("1_000_000denom4")))
	err = s.keeper.ModifyBid(s.ctx, &types.MsgModifyBid{
		AuctionId: bid.AuctionId,
		BidId:     bid.Id,
		Bidder:    bid.Bidder,
		Price:     bid.Price,
		Coin:      parseCoin("6_000_000denom4"),
	})
	s.Require().NoError(err)
	s.Require().True(fundraisingHooksReceiver.BeforeBidModifiedValid)

	// Calculate fixed price allocation
	mInfo, err := s.keeper.CalculateFixedPriceAllocation(s.ctx, auction)
	s.Require().NoError(err)

	// Allocate the selling coin
	err = s.keeper.AllocateSellingCoin(s.ctx, auction, mInfo)
	s.Require().NoError(err)
	s.Require().True(fundraisingHooksReceiver.BeforeSellingCoinsAllocatedValid)
}
