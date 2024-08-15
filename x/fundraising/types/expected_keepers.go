package types

import (
	"context"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
	GetModuleAddress(name string) sdk.AccAddress
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SendCoins(ctx context.Context, from, to sdk.AccAddress, amt sdk.Coins) error
	SpendableCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	InputOutputCoins(ctx context.Context, input banktypes.Input, outputs []banktypes.Output) error
	MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
}

// DistrKeeper defines the contract needed to be fulfilled for distribution keeper.
type DistrKeeper interface {
	FundCommunityPool(ctx context.Context, amount sdk.Coins, sender sdk.AccAddress) error
}

// ParamSubspace defines the expected Subspace interface for parameters.
type ParamSubspace interface {
	Get(context.Context, []byte, interface{})
	Set(context.Context, []byte, interface{})
}

// Event Hooks
// These can be utilized to communicate between a fundraising keeper and other keepers.
// The other keepers must implement this interface, which then the fundraising keeper can call.

// FundraisingHooks event hooks for fundraising auction and bid objects (noalias)
type FundraisingHooks interface {
	BeforeFixedPriceAuctionCreated(
		ctx context.Context,
		auctioneer string,
		startPrice math.LegacyDec,
		sellingCoin sdk.Coin,
		payingCoinDenom string,
		vestingSchedules []VestingSchedule,
		startTime time.Time,
		endTime time.Time,
	) error

	AfterFixedPriceAuctionCreated(
		ctx context.Context,
		auctionId uint64,
		auctioneer string,
		startPrice math.LegacyDec,
		sellingCoin sdk.Coin,
		payingCoinDenom string,
		vestingSchedules []VestingSchedule,
		startTime time.Time,
		endTime time.Time,
	) error

	BeforeBatchAuctionCreated(
		ctx context.Context,
		auctioneer string,
		startPrice math.LegacyDec,
		minBidPrice math.LegacyDec,
		sellingCoin sdk.Coin,
		payingCoinDenom string,
		vestingSchedules []VestingSchedule,
		maxExtendedRound uint32,
		extendedRoundRate math.LegacyDec,
		startTime time.Time,
		endTime time.Time,
	) error

	AfterBatchAuctionCreated(
		ctx context.Context,
		auctionId uint64,
		auctioneer string,
		startPrice math.LegacyDec,
		minBidPrice math.LegacyDec,
		sellingCoin sdk.Coin,
		payingCoinDenom string,
		vestingSchedules []VestingSchedule,
		maxExtendedRound uint32,
		extendedRoundRate math.LegacyDec,
		startTime time.Time,
		endTime time.Time,
	) error

	BeforeAuctionCanceled(
		ctx context.Context,
		auctionId uint64,
		auctioneer string,
	) error

	BeforeBidPlaced(
		ctx context.Context,
		auctionId uint64,
		bidId uint64,
		bidder string,
		bidType BidType,
		price math.LegacyDec,
		coin sdk.Coin,
	) error

	BeforeBidModified(
		ctx context.Context,
		auctionId uint64,
		bidId uint64,
		bidder string,
		bidType BidType,
		price math.LegacyDec,
		coin sdk.Coin,
	) error

	BeforeAllowedBiddersAdded(
		ctx context.Context,
		allowedBidders []AllowedBidder,
	) error

	BeforeAllowedBidderUpdated(
		ctx context.Context,
		auctionId uint64,
		bidder sdk.AccAddress,
		maxBidAmount math.Int,
	) error

	BeforeSellingCoinsAllocated(
		ctx context.Context,
		auctionId uint64,
		allocationMap map[string]math.Int,
		refundMap map[string]math.Int,
	) error
}

type FundraisingHooksWrapper struct{ FundraisingHooks }

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (FundraisingHooksWrapper) IsOnePerModuleType() {}
