<!-- order: 4 -->

# Messages

Messages (Msg) are objects that trigger state transitions. Msgs are wrapped in transactions (Txs) that clients submit to the network. The Cosmos SDK wraps and unwraps `fundraising` module messages from transactions.

## MsgCreateEnglishAuction

```go
// MsgCreateEnglishAuction defines a SDK message for creating a English type auction
type MsgCreateEnglishAuction struct {
	StartPrice       sdk.Dec           // starting price of the auction
	SellingCoin      sdk.Coin          // selling coin for the auction
	PayingCoinDenom string             // the paying coin denom that a bidder needs to bid for
	VestingAddress   string            // the vesting account that releases the paying amount of coins based on the schedules
	VestingSchedules []VestingSchedule // vesting schedules for the auction
	MaximumBidPrice  sdk.Dec           // the maximum bid price for the auction
	ExtendRate       sdk.Dec           // rate that decides if the auction needs another round
	StartTime        time.Time         // start time of the auction
	EndTime          time.Time         // end time of the auction
}
```

## MsgCreateFixedPriceAuction

```go
// MsgCreateFixedPriceAuction defines a SDK message for creating a fixed price type auction
type MsgCreateFixedPriceAuction struct {
	StartPrice       sdk.Dec           // starting price of the auction
	SellingCoin      sdk.Coin          // selling coin for the auction
	PayingCoinDenom  string            // the paying denom that participants need to bid for
	VestingAddress   string            // the vesting account that releases the paying amount of coins based on the schedules
	VestingSchedules []VestingSchedule // vesting schedules for the auction
	StartTime        time.Time         // start time of the auction
	EndTime          time.Time         // end time of the auction
}
```


## MsgCancelFundraising

```go
// MsgCancelFundraising defines a SDK message for cancelling an auction
type MsgCancelFundraising struct {
	Id uint64 // id of the auction
}
```

## MsgPlaceBid

```go
// MsgPlaceBid defines a SDK message for placing a bid for the auction
type MsgPlaceBid struct {
	AuctionId uint64   // id of the auction
	Price     sdk.Dec  // increasing bid price is only possible
	Coin      sdk.Coin // paying amount of coin that the bidder bids
}
```