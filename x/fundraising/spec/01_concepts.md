<!-- order: 1 -->

# Concepts

## Fundraising Module

The `x/fundraising` Cosmos SDK module is a module to raise funds as an auction of coins. This fundraising module provides an opportunity for a new project to onboard the ecosystem. It does not only allow the project to raise funds, but also increase its brand awareness before launching the project.

## Auction Types

This fundraising module provides two different types of auctions: 1) Fixed Price Auction and 2) Batch Auction.

## Fixed Price Auction

This fixed price auction is to sell a given amount of coins on a first-come, first-served basis.

### What an auctioneer does:
When an auctioneer creates a fixed price auction, it must determine the following parameters.

- `StartPrice`: fixed amount of the paying coins to get a selling coins (i.e., amount of paying coins per selling coin),
- `SellingCoin`: the denom and total amount of selling coin to be auctioned,
- `PayingCoinDenom`: the denom of coin to be used for payment,
- `StartTime`: when the auction starts,
- `EndTime`: when the auction ends,
- `VestingSchedules`: the vesting schedules to receive the sold amounts.

The auctioneer can cancel the auction before `StartTime`.

### What a bidder can/cannot do:

A bidder can place a new bid with a fixed amount of paying coins. 
A bidder cannot modify or cancel the existing bid it previously placed.

### When the auction ends:

The auction will end either when `EndTime` is arrived or when the entire `SellingCoin` is sold out.



## Batch Auction

This batch auction allows each bidder to participate in the auction by placing limit orders with a bid price chosen freely at any time within the auction period. An order book is created to record the bids with various bid prices.

### What an auctioneer does:

When an auctioneer creates this batch auction, it must determine the following parameters.

- `SellingCoin`: the denom and total amount of selling coins to be auctioned,
- `PayingCoinDenom`: the denom of coins to be used for payment,
- `StartTime`: when the auction starts,
- `EndTime`: when the auction ends,
- `VestingSchedules`: the vesting schedules to receive the sold amounts.

The auctioneer can cancel the auction before `StartTime`.

### What a bidder can/cannot do:

A bidder can do the following behaviors during the auction period.
1. Place a new bid
    - This auction provides two options for bidder to choose: 1) How-Much-Worth-To-Buy and 2) How-Many-Coins-To-Buy
        - (**Option A**) How-Much-Worth-To-Buy (fixed `PayingCoin`/varying `SellingCoin`): A bidder offers with a fixed amount of the paying coins and, if it wins, the bidder gets the selling coins, where the amount of the selling coins varies depending on the winning price determined after the auction period ends.
        - (**Option B**) How-Many-Coins-To-Buy (varying `PayingCoin`/fixed `SellingCoin`): A bidder offers for a fixed amount of the selling coin that the bidder wants to get if it wins. After the auction period ends, the remaining paying coins will be refunded depending on the winning price.
2. Replace the existing bid by a new one only with higher price and/or more quantities
    - The bidder can replace its existing bid, which is previously placed, by a new one with the same option between Option A and Option B.

A bidder cannot do the following behaviors during the auction period.

1. Cancel the existing bid
2. Replace the existing bid by a new one with lower price or fewer quantities.

### When the auction ends:

The auction will end when Auction End Time is arrived.

### How `WinningPrice` is determined:

Once the auction period ends, the bids are ordered in descending order of the bid prices to determine `WinningPrice`. `WinningPrice` is determined by finding the lowest price among the bid prices satisfying that the total amount of selling coins placed at more than or equal to the price is less the entire offering `SellingCoin`.
The bidders who placed at the higher price than the winning price become the winning bidders and get the selling coins at the same price, which is `WinningPrice`. 