<!-- order: 4 -->

# Messages

Messages (Msg) are objects that trigger state transitions. Msgs are wrapped in transactions (Txs) that clients submit to the network. The Cosmos SDK wraps and unwraps `fundraising` module messages from transactions.

## MsgCreateFixedPriceAuction

```go
// MsgCreateFixedPriceAuction defines an SDK message for creating a fixed price type auction
type MsgCreateFixedPriceAuction struct {	
	Auctioneer          string            // the owner of the auction
	StartPrice          sdk.Dec           // the starting price for the auction
	SellingCoin         sdk.Coin          // the selling coin for the auction
	PayingCoinDenom     string            // the denom that the auctioneer receives to raise funds
	VestingSchedules    []VestingSchedule // the vesting schedules for the auction
	StartTime           time.Time         // the start time of the auction
	EndTime             time.Time         // the end time of the auction
}
```
## MsgCreateBatchAuction

```go
// MsgCreateBatchAuction defines an SDK message for creating a batch type auction
type MsgCreateBatchAuction struct {
	Auctioneer       string            // the owner of the auction
	StartPrice       sdk.Dec           // the starting price for the auction
	MinBidPrice      sdk.Dec           // the minimum bid price that bidders must provide
	SellingCoin      sdk.Coin          // the selling coin for the auction
	PayingCoinDenom  string            // the denom that the auctioneer receives to raise funds
	VestingSchedules []VestingSchedule // the vesting schedules for the auction
	MaxExtendedRound uint32            // a maximum number of extended rounds
	ExtendedRate     sdk.Dec           // rate that determines if the auction needs another round, compared to the number of the matched bids at the previous end time.
	StartTime        time.Time         // the start time of the auction
	EndTime          time.Time         // the end times of the auction
}
```

## MsgCancelAuction

```go
// MsgCancelAuction defines an SDK message for cancelling an auction
type MsgCancelAuction struct {
	Auctioneer      string // the owner of the auction
	AuctionId       uint64 // id of the auction
}
```

## MsgPlaceBid
```go
// MsgPlaceBid defines an SDK message for placing a bid for the auction
// Bid price must be the start price for FixedPriceAuction whereas it can only be increased for EnglishAuction
type MsgPlaceBid struct {
	AuctionId       uint64   // id of the auction
	Bidder          string   // account that places a bid for the auction
	Type            BidType  // bid type; currently How-Much-Worth-To-Buy and How-Many-Coins-To-Buy are supported.
	Price           sdk.Dec  // bid price to bid for the auction
	Coin            sdk.Coin // targeted amount of coin that the bidder bids; the denom must be either the denom or SellingCoin or PayingCoinDenom
}
```

## MsgModifyBid
```go
// MsgModifyBid defines an SDK message for modifying a bid for the auction by replacing the existing bid by a new one.
// MsgModifyBid only applies for BatchAuction.
// Price cannot be lower than that of the original bid. If Price is the same as that of the original bid, the amount of BiddingCoin should be larger than that of the original bid. 
// The amount of BiddingCoin cannot be smaller than that of the original bid. If the amount of BiddingCoin is the same that of the original bid, Price should be higher than that of the original bid.  
// If the amount of BiddingCoin is the same that of the original bid, Price should be higher than that of the original bid. 
type MsgModifyBid struct {
	AuctionId       uint64   // id of the auction
	Bidder          string   // account that places a bid for the auction 
	BidId           uint64   // id of the bid of the bidder
	Price           sdk.Dec  // bid price to bid for the auction
	Coin            sdk.Coin // targeted amount of coin that the bidder bids; the denom must be either the denom or SellingCoin or PayingCoinDenom
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
	AuctionId       uint64        // id of the auction
	AllowedBidder   AllowedBidder // the bidder and their maximum bid amount
}
```