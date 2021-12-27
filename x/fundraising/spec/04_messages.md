<!-- order: 4 -->

# Messages

Messages (Msg) are objects that trigger state transitions. Msgs are wrapped in transactions (Txs) that clients submit to the network. The Cosmos SDK wraps and unwraps `fundraising` module messages from transactions.

## MsgCreateFixedPriceAuction

```go
// MsgCreateFixedPriceAuction defines a SDK message for creating a fixed price type auction
type MsgCreateFixedPriceAuction struct {
	Auctioneer       string            // account that creates the auction
	StartPrice       sdk.Dec           // starting price of the selling coin proportional to the paying coin
	SellingCoin      sdk.Coin          // selling amount of coin for the auction
	PayingCoinDenom  string            // paying coin denom that bidders need to bid with
	VestingSchedules []VestingSchedule // vesting schedules that release the paying amount of coins to the autioneer
	StartTime        time.Time         // start time of the auction
	EndTime          time.Time         // end time of the auction
}
```
## MsgCreateEnglishAuction

```go
// MsgCreateEnglishAuction defines a SDK message for creating a English type auction
type MsgCreateEnglishAuction struct {
	Auctioneer       string            // account that creates the auction
	StartPrice       sdk.Dec           // starting price of the selling coin; it is proportional to the price of paying coin denom
	SellingCoin      sdk.Coin          // selling amount of coin for the auction
	PayingCoinDenom  string            // paying coin denom that bidders need to bid with
	VestingSchedules []VestingSchedule // vesting schedules that release the paying amount of coins to the autioneer
	MaximumBidPrice  sdk.Dec           // maximum bid price that bidders can bid for the auction
	ExtendRate       sdk.Dec           // rate that determines if the auction needs an another round
	StartTime        time.Time         // start time of the auction
	EndTime          time.Time         // end time of the auction
}
```

## MsgCancelAuction

```go
// MsgCancelAuction defines a SDK message for cancelling an auction
type MsgCancelAuction struct {
	Auctioneer string // account that creates the auction
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