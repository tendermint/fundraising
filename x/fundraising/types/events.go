package types

// Event types for the farming module.
const (
	EventTypeCreateFixedPriceAuction = "create_fixed_price_auction"
	EventTypeCreateEnglishAuction    = "create_english_auction"
	EventTypeCancelAuction           = "cancel_auction"
	EventTypePlaceBid                = "place_bid"

	AttributeKeyAuctionId          = "auction_id" //nolint:golint
	AttributeKeyAuctioneerAddress  = "auctioneer_address"
	AttributeKeySellingPoolAddress = "selling_pool_address"
	AttributeKeyPayingPoolAddress  = "paying_pool_address"
	AttributeKeyVestingPoolAddress = "vesting_pool_address"
	AttributeKeyStartPrice         = "start_price"
	AttributeKeySellingCoin        = "selling_coin"
	AttributeKeyVestingSchedules   = "vesting_schedules"
	AttributeKeyPayingCoinDenom    = "paying_coin_denom"
	AttributeKeyAuctionStatus      = "auction_status"
	AttributeKeyStartTime          = "start_time"
	AttributeKeyEndTime            = "end_time"
	AttributeKeyBidderAddress      = "bidder_address"
	AttributeKeyBidPrice           = "bid_price"
	AttributeKeyBidCoin            = "bid_coin"
)
