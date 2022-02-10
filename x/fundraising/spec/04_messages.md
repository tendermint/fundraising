<!-- order: 4 -->

# Messages

Messages (Msg) are objects that trigger state transitions. Msgs are wrapped in transactions (Txs) that clients submit to the network. The Cosmos SDK wraps and unwraps `fundraising` module messages from transactions.

## MsgCreateFixedPriceAuction

```go
// MsgCreateFixedPriceAuction defines a SDK message for creating a fixed price type auction
type MsgCreateFixedPriceAuction struct {
	Auctioneer       string            // the owner of the auction
	StartPrice       sdk.Dec           // the starting price for the auction; it is proportional to the price of paying coin denom
	SellingCoin      sdk.Coin          // the selling coin for the auction
	PayingCoinDenom  string            // the denom that the auctioneer receives to raise funds
	VestingSchedules []VestingSchedule // the vesting schedules for the auction
	StartTime        time.Time         // the start time of the auction
	EndTime          time.Time         // the end time of the auction
}
```
## MsgCreateEnglishAuction

```go
// MsgCreateEnglishAuction defines a SDK message for creating a English type auction
type MsgCreateEnglishAuction struct {
	Auctioneer       string            // the owner of the auction
	StartPrice       sdk.Dec           // the starting price for the auction
	SellingCoin      sdk.Coin          // the selling coin for the auction
	PayingCoinDenom  string            // the denom that the auctioneer receives to raise funds
	VestingSchedules []VestingSchedule // the vesting schedules for the auction
	MaximumBidPrice  sdk.Dec           // the maximum bid price that bidders can bid for the auction
	ExtendRate       sdk.Dec           // the rate that determines if the auction needs an another round
	StartTime        time.Time         // the start time of the auction
	EndTime          time.Time         // the end time of the auction
}
```

## MsgCancelAuction

```go
// MsgCancelAuction defines a SDK message for cancelling an auction
type MsgCancelAuction struct {
	Auctioneer string // the owner of the auction
	AuctionId  uint64 // id of the auction
}
```

## MsgPlaceBid
```go
// MsgPlaceBid defines a SDK message for placing a bid for the auction
// Bid price must be the start price for FixedPriceAuction whereas it can only be increased for EnglishAuction
type MsgPlaceBid struct {
	AuctionId uint64   // id of the auction
	Bidder    string   // account that places a bid for the auction
	Price     sdk.Dec  // bid price to bid for the auction
	Coin      sdk.Coin // paying amount of coin that the bidder bids
}
```


## MsgAddAllowedBidder

This message is a custom message that is created for testing purpose only. It adds an allowed bidder to `AllowedBidders` for the auction. 
It is accessible when you build `fundraisingd` binary by the following command:

```bash
make install-testing
```

```go
// MsgAddAllowedBidder defines a SDK message to add an allowed bidder
type MsgAddAllowedBidder struct {
	AuctionId     uint64        // the id of the auction
	AllowedBidder AllowedBidder // the bidder and their maximum bid amount
}
```