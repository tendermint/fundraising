package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/fundraising/x/fundraising/types"

	_ "github.com/stretchr/testify/suite"
)

var _ types.FundraisingHooks = &MockFundraisingHooksReceiver{}

// MockFundraisingHooksReceiver event hooks for governance proposal object (noalias)
type MockFundraisingHooksReceiver struct {
	BeforeFixedPriceAuctionCreatedValid bool
	BeforeBatchAuctionCreatedValid      bool
	BeforeAuctionCanceledValid          bool
	BeforeBidPlacedValid                bool
	BeforeBidModifiedValid              bool
	BeforeAllowedBiddersAddedValid      bool
	BeforeAllowedBidderUpdatedValid     bool
}

func (h *MockFundraisingHooksReceiver) BeforeFixedPriceAuctionCreated(
	ctx sdk.Context,
	auctioneer string,
	startPrice sdk.Dec,
	sellingCoin sdk.Coin,
	payingCoinDenom string,
	vestingSchedules []types.VestingSchedule,
	startTime time.Time,
	endTime time.Time,
) {
	h.BeforeFixedPriceAuctionCreatedValid = true
}

func (h *MockFundraisingHooksReceiver) BeforeBatchAuctionCreated(
	ctx sdk.Context,
	auctioneer string,
	startPrice sdk.Dec,
	minBidPrice sdk.Dec,
	sellingCoin sdk.Coin,
	payingCoinDenom string,
	vestingSchedules []types.VestingSchedule,
	maxExtendedRound uint32,
	extendedRoundRate sdk.Dec,
	startTime time.Time,
	endTime time.Time,
) {
	h.BeforeBatchAuctionCreatedValid = true
}

func (h *MockFundraisingHooksReceiver) BeforeAuctionCanceled(
	ctx sdk.Context,
	auctionId uint64,
	auctioneer string,
) {
	h.BeforeAuctionCanceledValid = true
}

func (h *MockFundraisingHooksReceiver) BeforeBidPlaced(
	ctx sdk.Context,
	auctionId uint64,
	bidder string,
	bidType types.BidType,
	price sdk.Dec,
	coin sdk.Coin,
) {
	h.BeforeBidPlacedValid = true
}

func (h *MockFundraisingHooksReceiver) BeforeBidModified(
	ctx sdk.Context,
	auctionId uint64,
	bidder string,
	bidType types.BidType,
	price sdk.Dec,
	coin sdk.Coin,
) {
	h.BeforeBidModifiedValid = true
}

func (h *MockFundraisingHooksReceiver) BeforeAllowedBiddersAdded(
	ctx sdk.Context,
	auctionId uint64,
	allowedBidders []types.AllowedBidder,
) {
	h.BeforeAllowedBiddersAddedValid = true
}

func (h *MockFundraisingHooksReceiver) BeforeAllowedBidderUpdated(
	ctx sdk.Context,
	auctionId uint64,
	bidder sdk.AccAddress,
	maxBidAmount sdk.Int,
) {
	h.BeforeAllowedBidderUpdatedValid = true
}

func (s *KeeperTestSuite) TestHooks() {
	fundraisingHooksReceiver := MockFundraisingHooksReceiver{}

	// Set hooks
	s.keeper.SetHooks(types.NewMultiFundraisingHooks(&fundraisingHooksReceiver))

	s.Require().False(fundraisingHooksReceiver.BeforeFixedPriceAuctionCreatedValid)
	s.Require().False(fundraisingHooksReceiver.BeforeBatchAuctionCreatedValid)
	s.Require().False(fundraisingHooksReceiver.BeforeAuctionCanceledValid)
	s.Require().False(fundraisingHooksReceiver.BeforeBidPlacedValid)
	s.Require().False(fundraisingHooksReceiver.BeforeBidModifiedValid)
	s.Require().False(fundraisingHooksReceiver.BeforeAllowedBiddersAddedValid)
	s.Require().False(fundraisingHooksReceiver.BeforeAllowedBidderUpdatedValid)

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

	// Create a batch auction
	batchAuction := s.createBatchAuction(
		s.addr(1),
		parseDec("0.5"),
		parseDec("0.1"),
		parseCoin("1_000_000_000_000denom3"),
		"denom4",
		[]types.VestingSchedule{},
		1,
		sdk.MustNewDecFromStr("0.2"),
		time.Now().AddDate(0, 0, -1),
		time.Now().AddDate(0, 0, -1).AddDate(0, 2, 0),
		true,
	)
	s.Require().True(fundraisingHooksReceiver.BeforeBatchAuctionCreatedValid)

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
	s.Require().True(fundraisingHooksReceiver.BeforeAuctionCanceledValid)

	// Get already started batch auction
	auction, found := s.keeper.GetAuction(s.ctx, batchAuction.Id)
	s.Require().True(found)

	// Add allowed bidder
	s.addAllowedBidder(auction.GetId(), s.addr(3), parseInt("100_000_000_000"))
	s.Require().True(fundraisingHooksReceiver.BeforeAllowedBiddersAddedValid)

	// Update the allowed bidder
	s.keeper.UpdateAllowedBidder(s.ctx, auction.GetId(), s.addr(3), parseInt("110_000_000_000"))
	s.Require().True(fundraisingHooksReceiver.BeforeAllowedBidderUpdatedValid)

	// Place a bid
	bid := s.placeBidBatchWorth(auction.GetId(), s.addr(3), parseDec("0.55"), parseCoin("5_000_000denom4"), sdk.NewInt(10_000_000), true)
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
}
