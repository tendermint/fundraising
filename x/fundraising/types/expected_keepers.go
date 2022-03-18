package types

import (
	time "time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// AccountKeeper defines the expected account keeper.
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	GetModuleAddress(name string) sdk.AccAddress
}

// BankKeeper defines the expected bank send keeper.
type BankKeeper interface {
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	InputOutputCoins(ctx sdk.Context, inputs []banktypes.Input, outputs []banktypes.Output) error
}

// DistrKeeper is the keeper of the distribution store
type DistrKeeper interface {
	FundCommunityPool(ctx sdk.Context, amount sdk.Coins, sender sdk.AccAddress) error
}

// Event Hooks
// These can be utilized to communicate between a fundraising keeper and other keepers.
// The other keepers must implement this interface, which then the fundraising keeper can call.

// FundraisingHooks event hooks for fundraising auction and bid objects (noalias)
type FundraisingHooks interface {
	BeforeFixedPriceAuctionCreated(
		ctx sdk.Context,
		auctioneer string,
		startPrice sdk.Dec,
		minBidPrice sdk.Dec,
		sellingCoin sdk.Coin,
		payingCoinDenom string,
		vestingSchedules []VestingSchedule,
		startTime time.Time,
		endTime time.Time,
	)

	// BeforeBatchAuctionCreated(

	// )

	BeforeAuctionCanceled(
		ctx sdk.Context,
		auctionId uint64,
		auctioneer string,
	)

	BeforeBidPlaced(
		ctx sdk.Context,
		auctionId uint64,
		bidder string,
		bidType BidType,
		price sdk.Dec,
		coin sdk.Coin,
	)

	BeforeBidModified(
		ctx sdk.Context,
		auctionId uint64,
		bidder string,
		bidType BidType,
		price sdk.Dec,
		coin sdk.Coin,
	)

	BeforeAllowedBidderAdded(
		ctx sdk.Context,
		auctionId uint64,
		allowedBidder AllowedBidder,
	)
}
