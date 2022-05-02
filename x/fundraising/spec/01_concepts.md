<!-- order: 1 -->

# Concepts

## Fundraising Module

The `x/fundraising` module is a Cosmos SDK module that provides a functionality to raise funds for a new project to onboard the ecosystem. It helps them to increase their brand awareness before launching their project. 

## Design Decision

The module is fundamentally designed to delegate authorization to an external module to add allowed bidder list for an auction. When an auction is created, it is closed state. It means that there is no bidder who is authorized to place a bid. The bidder must be added by an external module. 

## Auction Type

The module allows the creation of the following auction types:

* `FixedPriceAuction` 
* `BatchAuction`

## Fixed Price Auction

A FixedPriceAuction is to sell a given amount of coins in a fixed price. It is first-come and first-served basis. The module expects an external module or a project to create a fixed price auction by setting parameters needed for an auction, such as how many coins to sell, what type of coin denomination is payable by a bidder in exchange for the selling coin, fixed price for each selling coin, start and end time for the auction, and etc.  When an auction is created successfully by paying a creation fee, the module expects an external module (being as an auctioneer) to add allowed bidders list for the auction. In this process, the external module has a control over who can place a bid and the bidder’s maximum amount to place a bid for the auction. When an auction is started, allowed bidders can start to place their bids until the auction ends; however, as it is first-come and first-served basis, the selling amount of coin can be sold at any time. The distribution of selling coin will occur when the auction is ended.

### What an auctioneer does:

A fixed price auction must determine the following parameters:

- `StartPrice`: fixed amount of the paying coins to get a selling coins (i.e., amount of paying coins per selling coin),
- `SellingCoin`: the denom and total amount of selling coin to be auctioned,
- `PayingCoinDenom`: the denom of coin to be used for payment,
- `StartTime`: when the auction starts,
- `EndTime`: when the auction ends,
- `VestingSchedules`: the vesting schedules to allocate the sold amounts of paying coins to the auctioneer.

Note that the auctioneer can cancel the auction as long as an auction has not started.

### What a bidder can/cannot do:

As explained in `Design Decision`, bidders are not allowed to place their bids unless they are listed in `AllowedBidders` for an auction. Allowed bidders can place their bids either with paying coin denom (willing to pay in exchange of the selling coin) or selling coin denom (how many selling coins that a bidder is willing to buy). The module takes care of it. Once bids are placed, they can't be canceled. Bids can only be modified with higher bidding price or increasing bidding amount.

## Batch Auction

A `BatchAuction` is different from `FixedPriceAuction`. It provides an interesting and dynamic way for bidders to participate in an auction. The module expects an external module (being as an auctioneer) to create a batch auction by setting parameters needed for an auction. The creation process is the same as a fixed price auction creation. For a batch auction. there is no fixed price. A matched price (final price) gets determined by a number of bids with their bidding prices and amounts. When an auction is started, allowed bidders start to place their bids with the bidding price that they think each selling coin is worth. Bidders’ bidding amount is reserved once they place their bids until the end of an auction. They can’t cancel their bids, but they can modify them with higher bidding price or increasing amount. Since there is no guarantee that a bid gets matched to win the auction, it is up to allowed bidders to carefully monitor the demand and place or modify their bids accordingly. At the end of an auction, the module gets all recorded bids and calculate a matched price (final price) with a number of bids with bidding prices and amounts. 

### What an auctioneer does:

When an auctioneer creates this batch auction, it must determine the following parameters.

- `SellingCoin`: the denom and total amount of selling coins to be auctioned,
- `PayingCoinDenom`: the denom of coins to be used for payment,
- `StartTime`: when the auction starts,
- `EndTimes`: when the auction ends including the possible extended rounds,
- `VestingSchedules`: the vesting schedules to allocate the sold amounts of paying coins to the auctioneer,
- `MinBidPrice`: the minimum bid price that the bidders must place a bid with,
- `MaxExtendedRound`: the maximum number of additional round for bidding,
- `ExtendedRoundRate`: the condition in a reduction rate of the number of the matched bids.

Note that the auctioneer can cancel the auction as long as an auction has not started. Also, the extended round is to prevent the auction sniping technique, which is, e.g., to bid large amount of selling coins with a bid price slightly higher than the matched price, where this kind of last moment bid as auction sniping results in a sudden reduction of the matched bids. 

In order to provide more opportunity to bidders in case of auction sniping, the extended round is triggered if the reduction of the matched bids are more than `ExtendedRoundRate` compared to the number of matched bids at the previous end time.

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
