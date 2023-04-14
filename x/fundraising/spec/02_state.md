<!-- order: 2 -->

# State

The `fundraising` module keeps track of the auction and bid states.

## Auction Interface

The auction interface exposes methods to read and write standard auction information.

Note that all of these methods operate on a auction struct that confirms to the interface. In order to write the auction to the store, the auction keeper is required.

```go
// AuctionI is an interface that inherits the BaseAuction and exposes common functions 
// to get and set standard auction data.
type AuctionI interface {
	GetId() uint64
	SetId(uint64) error

	GetType() AuctionType
	SetType(AuctionType) error

	GetAuctioneer() string
	SetAuctioneer(string) error

	GetSellingReserveAddress() string
	SetSellingReserveAddress(string) error

	GetPayingReserveAddress() string
	SetPayingReserveAddress(string) error

	GetStartPrice() sdk.Dec
	SetStartPrice(sdk.Dec) error

	GetSellingCoin() sdk.Coin
	SetSellingCoin(sdk.Coin) error

	GetPayingCoinDenom() string
	SetPayingCoinDenom(string) error

	GetVestingReserveAddress() string
	SetVestingReserveAddress(string) error

	GetVestingSchedules() []VestingSchedule
	SetVestingSchedules([]VestingSchedule) error
	
	GetStartTime() time.Time
	SetStartTime(time.Time) error

	GetEndTimes() []time.Time
	SetEndTimes([]time.Time) error

	GetStatus() AuctionStatus
	SetStatus(AuctionStatus) error

	ShouldAuctionStarted(t time.Time) bool
	ShouldAuctionClosed(t time.Time) bool

	Validate() error
}
```

## BaseAuction

A base auction stores all requisite fields directly in a struct.

```go
// BaseAuction defines a base auction type. It contains all the necessary fields
// for basic auction functionality. Any custom auction type should extend this
// type for additional functionality (e.g. english auction, fixed price auction).
type BaseAuction struct {
	Id                    uint64            // id of the auction
	Type                  AuctionType       // the auction type; currently FixedPrice and English are supported
	Auctioneer            string            // the owner of the auction
	SellingReserveAddress string            // the reserve account to collect selling coins from the auctioneer
	PayingReserveAddress  string            // the reserve account to collect paying coins from the bidders
	StartPrice            sdk.Dec           // the starting price
	SellingCoin           sdk.Coin          // the selling amount of coin
	PayingCoinDenom       string            // the denom that the auctioneer receives to raise funds
	VestingReserveAddress string            // the reserve account that releases the accumulated paying coins based on the schedules
	VestingSchedules      []VestingSchedule // the vesting schedules for the auction
	StartTime             time.Time         // the start time of the auction
	EndTimes              []time.Time       // the end times of the auction; it is an array since extended round(s) can occur
	Status                AuctionStatus     // the auction status
}
```

```go
// AllowedBidder defines a bidder who is allowed to bid with max number of bids.
type AllowedBidder struct {
	Bidder          string  // a bidder who is allowed to bid
	MaxBidAmount    math.Int // a maximum amount of bids per bidder
}
```

## Vesting
```go
// VestingSchedule defines the vesting schedule for the owner of an auction.
type VestingSchedule struct {
	ReleaseTime time.Time // the release time for vesting coin distribution
	Weight      sdk.Dec   // the vesting weight for the schedule
}

// VestingQueue defines the vesting queue.
type VestingQueue struct {
	AuctionId       uint64    // id of the auction
	Auctioneer      string    // the owner of the auction
	PayingCoin      sdk.Coin  // the paying amount of coin for the vesting
	ReleaseTime     time.Time // the release time of the vesting 
	Released        bool      // the distribution status 
}
```

## Auction Type

```go
// AuctionType is the type of an auction.
type AuctionType uint32

const (
	// AUCTION_TYPE_UNSPECIFIED defines an invalid auction type
	AuctionTypeNil AuctionType = 0
	// AUCTION_TYPE_FIXED_PRICE defines the fixed price auction type
	AuctionTypeFixedPrice AuctionType = 1
	// AUCTION_TYPE_BATCH defines the batch auction type
	AuctionTypeBatch AuctionType = 2
)

// FixedPriceAuction defines the fixed price auction type
type FixedPriceAuction struct {
	*BaseAuction
	RemainingSellingCoin sdk.Coin // the remaining amount of coin to sell
}

// BatchAuction defines the batch auction type 
type BatchAuction struct {
    *BaseAuction
	MinBidPrice			sdk.Dec	// the minimum bid price that bidders must provide
	MatchedPrice		sdk.Dec	// the matched price of the auction (a.k.a., winning price)
    MaxExtendedRound    uint32  // the maximum number of extended rounds
    ExtendedRate        sdk.Dec // the rate that determines if the auction needs another round; compared to the number of the matched bids at the previous end time.
}
```

## Auction Status

```go
// AuctionStatus is the status of an auction
type AuctionStatus uint32

const (
	// AUCTION_STATUS_UNSPECIFIED defines an invalid auction status
	StatusNil AuctionStatus = 0
	// AUCTION_STATUS_STANDY_BY defines an auction status before it opens
	StatusStandBy AuctionStatus = 1
	// AUCTION_STATUS_STARTED defines an auction status that is started
	StatusStarted AuctionStatus = 2
	// AUCTION_STATUS_VESTING defines an auction status that is in distribution based on the vesting schedules
	StatusVesting AuctionStatus = 3
	// AUCTION_STATUS_FINISHED defines an auction status that is finished 
	StatusFinished AuctionStatus = 4
	// AUCTION_STATUS_CANCELLED defines an auction sttus that is cancelled
	StatusCancelled AuctionStatus = 5
)
```

## Bid

```go
// Bid defines a standard bid for an auction.
type Bid struct {
	AuctionId 	uint64   // id of the auction
	Bidder    	string   // the account that bids for the auction
	Id        	uint64   // id of the bid of the bidder
	Type      	BidType  // the bid type; currently Fixed-Price, How-Much-Worth-To-Buy and How-Many-Coins-To-Buy are supported.
	Price     	sdk.Dec  // the price for the bid
	Coin      	sdk.Coin // targeted amount of coin that the bidder bids; the denom must be either the denom or SellingCoin or PayingCoinDenom
	IsMatched	bool     // the bid that is determined to be matched (a.k.a., winner) when an auction ends; default value is false
}
```

## Bid Type

```go
// BidType is the type of a bid.
type BidType uint32

const (
	// BID_TYPE_UNSPECIFIED defines an invalid bid type
	BidTypeNil          BidType = 0
 	// BID_TYPE_FIXED_PRICE defines a bid type for a fixed price auction type
	BidTypeFixedPrice   BidType = 1
	// Bid_TYPE_BATCH_WORTH defines a bid type for How-Much-Worth-to-Buy of a batch auction
	BidTypeBatchWorth   BidType = 2
	// Bid_TYPE_BATCH_MANY defines a bid type for How-Many-Coins-to-Buy of a batch auction
	BidTypeBatchMany    BidType = 3
)
```

For `FixedPriceAuction`,
- `BidType` must be set to `BidTypeFixedPrice`,
- `BidPrice` must be set as `StartPrice` in `BaseAuction`, and
- the denom of `BidCoin` can be set to either `PayingCoinDenom` or the denom of `SellingCoin`.

For `BatchAuction`and `BidTypeBatchWorth`,
- the denom of `BidCoin` must be set as `PayingCoinDenom`.

For `BatchAuction`and `BidTypeBatchMany`,
- the denom of `BidCoin` must be set as the denom of `SellingCoin`.

## Parameters

- ModuleName: `fundraising`
- RouterKey: `fundraising`
- StoreKey: `fundraising`
- QuerierRoute: `fundraising`

## Stores

Stores are KVStores in the multi-store. The key to find the store is the first parameter in the list.

### The key for the latest auction id

- `LastAuctionIdKey: 0x11 -> Uint64Value(lastAuctionId)`

### The key for the latest bid id

- `LastBidIdKey: 0x12 | AuctionId -> Uint64Value(lastBidId)`

### The key to retrieve the auction object from the auction id

- `AuctionKey: 0x21 | AuctionId -> ProtocolBuffer(Auction)`

### The key to retrieve the allowed bidder object for the auction

- `AllowedBidderKey: 0x22 | AuctionId | BidderAddrLen (1 byte) | BidderAddr -> ProtocolBuffer(AllowedBidder)`

### The key to retrieve the bid object from the auction id and bid id

- `BidKey: 0x31 | AuctionId | BidId -> ProtocolBuffer(Bid)`

### The index key to retrieve the bid object from the bidder address

- `BidIndexKey: 0x32 | BidderAddrLen (1 byte) | BidderAddr | AuctionId | BidId -> nil`

### The key to retrieve the last matched bids length 

- `LastMatchedBidsLenKey: 0x33 | AuctionId -> Uint64Value(lastMatchedBidsLen)`

### The key to retrieve the vesting queue object from the  auction id and 

- `VestingQueueKey: 0x41 | AuctionId | sdk.FormatTimeBytes(releaseTime) -> ProtocolBuffer(VestingQueue)`