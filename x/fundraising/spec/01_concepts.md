<!-- order: 1 -->

# Concepts

## Fundraising Module

The `x/fundraising` module is a Cosmos SDK module that provides a functionality to raise funds for a new project to onboard the ecosystem. It helps them to increase their brand awareness before launching a project. 

## Important Design Decision

The module is fundamentally designed to delegate authorization to an external module to add allowed bidder list `AllowedBidders` for an auction. When an auction is created, it is always closed state. It means that there is no single bidder who is authorized to place a bid for the selling coin unless they are added in the auction's allowed bidders list.

## Auction Type

The module allows the creation of two different types of an auction. 

* `FixedPriceAuction` 
* `BatchAuction`

## Fixed Price Auction

A fixed price auction is to sell a given amount of coins on a first-come and first-served basis. An external module creates a fixed price auction by setting parameters, such as start price for each coin, how many coins they are selling, what coin they are accepting in return, vesting schedules for them to receive paying coin, when to start and end, and so forth. When the auction is created successfully, the external module needs to add allowed bidders by using the implemented `AddAllowedBidders` function. During this step, they can limit a bidder’s maximum bid amount `MaxBidAmount`.

### What an auctioneer does:

A fixed price auction must determine the following parameters:

- `AllowedBidders`: the list of the bidders to be allowed to participate in the auction,
- `StartPrice`: fixed amount of the paying coins to get a selling coins (i.e., amount of paying coins per selling coin),
- `SellingCoin`: the denom and total amount of selling coin to be auctioned,
- `PayingCoinDenom`: the denom of coin to be used for payment,
- `StartTime`: when the auction starts,
- `EndTime`: when the auction ends,
- `VestingSchedules`: the vesting schedules to allocate the sold amounts of paying coins to the auctioneer.

The auctioneer can cancel the auction before `StartTime`.

In `AllowedBidders`, each bidder can be set with `MaxBidAmount`, which is the maximum number of selling coins that the bidder can get.

### What a bidder can/cannot do:

A bidder only listed in `AllowedBidders` can place a new bid with a fixed amount of either paying coins or selling coins. 
A bidder cannot modify or cancel the existing bid it previously placed.

### When the auction ends:

The auction will end either when `EndTime` is arrived or when the entire `SellingCoin` is sold out.



## Batch Auction

This batch auction allows each bidder to participate in the auction by placing limit orders with a bid price chosen freely at any time within the auction period. An order book is created to record the bids with various bid prices.

### What an auctioneer does:

When an auctioneer creates this batch auction, it must determine the following parameters.

- `AllowedBidders`: the list of the bidders to be allowed to participate in the auction,
- `SellingCoin`: the denom and total amount of selling coins to be auctioned,
- `PayingCoinDenom`: the denom of coins to be used for payment,
- `StartTime`: when the auction starts,
- `EndTimes`: when the auction ends including the possible extended rounds,
- `VestingSchedules`: the vesting schedules to allocate the sold amounts of paying coins to the auctioneer,
- `MinBidPrice`: the minimum bid price that the bidders must place a bid with,
- `MaxExtendedRound`: the maximum number of additional round for bidding,
- `ExtendedRoundRate`: the condition in a reduction rate of the number of the matched bids.

The auctioneer can cancel the auction before `StartTime`.
In `AllowedBidders`, each bidder can be set with `MaxBidAmount`, which is the maximum number of selling coins that the bidder can get.
Note that the extended round is to prevent the auction sniping, which is, e.g., to bid large amount of selling coins with a bid price slightly higher than the matched price, where this kind of last moment bid as auction sniping results in a sudden reduction of the matched bids. 
In order to provide more opportunity to bidders in case of auction sniping, the extended round is given if the reduction of the matched bids are more than `ExtendedRoundRate` compared to the number of matched bids at the previous end time.

### What a bidder can/cannot do:

A bidder only listed in `AllowedBidders` can do the following behaviors during the auction period between `StartTime` and `EndTimes`.
1. Place a new bid
    - This auction provides two options for bidder to choose: 1) How-Much-Worth-To-Buy and 2) How-Many-Coins-To-Buy
        - (`BidType` of `BidTypeBatchWorth`) How-Much-Worth-To-Buy (fixed `PayingCoin`/varying `SellingCoin`): A bidder places a bid with a fixed amount of the paying coins and, if it wins, the bidder gets the selling coins, where the amount of the selling coins varies depending on the matched price determined after the auction period ends.
        - (`BidType` of `BidTypeBatchMany`) How-Many-Coins-To-Buy (varying `PayingCoin`/fixed `SellingCoin`): A bidder places a bid for a fixed amount of the selling coin that the bidder wants to get if it wins. After the auction period ends, the remaining paying coins will be refunded depending on the matched price.
2. Modify the existing bid by replacing with a new one only with higher price and/or larger quantity
    - The bidder can replace its existing bid, which is previously placed, by a new one with the same `BidType` between `BidTypeBatchWorth` and `BidTypeBatchMany`.

A bidder cannot do the following behaviors during the auction period.

1. Cancel the existing bid 
2. Modify the existing bid by replacing with a new one with lower price or smaller quantity.

### When the auction ends:

The auction will end when the last time of `EndTimes` is arrived.

### How `MatchedPrice` is determined:

Once the auction period ends, the bids are ordered in descending order of the bid prices to determine `MatchedPrice`. `MatchedPrice` is determined by finding the lowest price among the bid prices satisfying that the total amount of selling coins placed at more than or equal to the price is less the entire offering `SellingCoin`.
The bidders who placed at the higher price than the matched price become the matched bidders and get the selling coins at the same price, which is `MatchedPrice`. 
