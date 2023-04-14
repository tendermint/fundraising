<!-- order: 8 -->

# Hooks

Other modules may register operations to execute when a certain event has
occurred within fundraising. These events can be registered to execute right `Before` the fundraising event (as per the hook name). 
The following hooks can registered with fundraising:

```go
BeforeFixedPriceAuctionCreated(
    ctx sdk.Context,
    auctioneer string,
    startPrice sdk.Dec,
    sellingCoin sdk.Coin,
    payingCoinDenom string,
    vestingSchedules []VestingSchedule,
    startTime time.Time,
    endTime time.Time,
)

AfterFixedPriceAuctionCreated(
    ctx sdk.Context,
    auctionId uint64,
    auctioneer string,
    startPrice sdk.Dec,
    sellingCoin sdk.Coin,
    payingCoinDenom string,
    vestingSchedules []VestingSchedule,
    startTime time.Time,
    endTime time.Time,
)

BeforeBatchAuctionCreated(
    ctx sdk.Context,
    auctioneer string,
    startPrice sdk.Dec,
    minBidPrice sdk.Dec,
    sellingCoin sdk.Coin,
    payingCoinDenom string,
    vestingSchedules []VestingSchedule,
    maxExtendedRound uint32,
    extendedRoundRate sdk.Dec,
    startTime time.Time,
    endTime time.Time,
)

AfterBatchAuctionCreated(
    ctx sdk.Context,
    auctionId uint64,
    auctioneer string,
    startPrice sdk.Dec,
    minBidPrice sdk.Dec,
    sellingCoin sdk.Coin,
    payingCoinDenom string,
    vestingSchedules []VestingSchedule,
    maxExtendedRound uint32,
    extendedRoundRate sdk.Dec,
    startTime time.Time,
    endTime time.Time,
)

BeforeAuctionCanceled(
    ctx sdk.Context,
    auctionId uint64,
    auctioneer string,
)

BeforeBidPlaced(
    ctx sdk.Context,
    auctionId uint64,
    bidId uint64,
    bidder string,
    bidType BidType,
    price sdk.Dec,
    coin sdk.Coin,
)

BeforeBidModified(
    ctx sdk.Context,
    auctionId uint64,
    bidId uint64,
    bidder string,
    bidType BidType,
    price sdk.Dec,
    coin sdk.Coin,
)

BeforeAllowedBiddersAdded(
    ctx sdk.Context,
    allowedBidders []AllowedBidder,
)

BeforeAllowedBidderUpdated(
    ctx sdk.Context,
    auctionId uint64,
    bidder sdk.AccAddress,
    maxBidAmount math.Int,
)
```
