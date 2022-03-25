package types

// Event types for the farming module.
const (
	EventTypeCreateFixedPriceAuction = "create_fixed_price_auction"
	EventTypeCreateBatchAuction      = "create_batch_auction"
	EventTypeCancelAuction           = "cancel_auction"
	EventTypePlaceBid                = "place_bid"

	AttributeKeyAuctionId             = "auction_id" //nolint:golint
	AttributeKeyAuctioneerAddress     = "auctioneer_address"
	AttributeKeySellingReserveAddress = "selling_pool_address"
	AttributeKeyPayingReserveAddress  = "paying_pool_address"
	AttributeKeyVestingReserveAddress = "vesting_pool_address"
	AttributeKeyStartPrice            = "start_price"
	AttributeKeySellingCoin           = "selling_coin"
	AttributeKeyRemainingSellingCoin  = "remaining_selling_coin"
	AttributeKeyVestingSchedules      = "vesting_schedules"
	AttributeKeyPayingCoinDenom       = "paying_coin_denom"
	AttributeKeyAuctionStatus         = "auction_status"
	AttributeKeyStartTime             = "start_time"
	AttributeKeyEndTime               = "end_time"
	AttributeKeyBidderAddress         = "bidder_address"
	AttributeKeyBidPrice              = "bid_price"
	AttributeKeyBidCoin               = "bid_coin"
	AttributeKeyBidAmount             = "bid_amount"
	AttributeKeyMinBidPrice           = "min_bid_price"
	AttributeKeyMaxExtendedRound      = "maximum_extended_round"
	AttributeKeyExtendedRoundRate     = "extended_round_rate"
)
