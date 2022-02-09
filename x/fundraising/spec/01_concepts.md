<!-- order: 1 -->

# Concepts

## Fundraising Module

The `x/fundraising` Cosmos SDK module is a module to raise funds as an auction of tokens. This fundraising module provides an opportunity for new projects to onboard the ecosystem. It does not only allow projects to raise funds, but also increase their brand awareness before launching their projects.

## Auction Types

There are two different types of auctions in the fundraising module.

### Fixed Price Auction

This fixed price auction is to sell a given amount of tokens in first-come, first-served basis.

#### What an auctioneer does
When an auctioneer creates this fixed price auction, it must determine the following parameters.

- **ReceivingToken**: the denom of tokens to be auctioned,
- **PayingToken**: the denom of tokens to be used for payment,
- **Price**: fixed amount of the paying tokens to get a receiving token,
- **AuctionStartTime**: when the auction starts,
- **AuctionEndTime**: when the auction ends,
- **OfferingQuantity**: total amount of receiving tokens to be auctioned.

#### What bidders can/cannot do

A bidder can place a new bid with a fixed amount of paying tokens. 
A bidder cannot modify or cancel the existing bid it previously placed.

#### When the auction ends

The auction will end either when the auction end time is arrived or when the entire offering quantity is sold out.



### Order Book Auction (Batch Auction)

This order book auction allows each bidder to participate in the auction by placing limit orders with the price chosen freely and at any time within the auction period. An order book is created to record the bids with various bid prices.

#### What an auctioneer does

When an auctioneer creates this order book auction, it must determine the following parameters.

- **ReceivingToken**: the denom of tokens to be auctioned,
- **PayingToken**: the denom of tokens to be used for payment,
- **AuctionStartTime**: when the auction starts,
- **AuctionEndTime**: when the auction ends,
- **OfferingQuantity**: total amount of receiving tokens to be auctioned.

#### What bidders can/cannot do

A bidder can do the following behaviors during the auction period:
1. Place a new bid
    - This auction provides two options the bidders for bidding: 1) How-Much-Worth-To-Buy and 2) How-Many-Tokens-To-Buy
        - (**Option A**) How-Much-Worth-To-Buy (fixed paying tokens/varying receiving tokens): A bidder offers with a fixed amount of the paying tokens and, if win, the bidder gets the receiving tokens, where the amount of the receiving tokens varies depending on the price of the receiving token.
        - (**Option B**) How-Many-Tokens-To-Buy (varying paying tokens/fixed receiving tokens): A bidder offers for a fixed amount of the receiving token that the bidder wants to get if win. The residual paying tokens the bidder placed can be refunded depending on the last price.
    - Each bidder can choose one of the above two options. The two options mean 1) how much worth in paying tokens of the receiving tokens the user wants to buy, and 2) how many receiving tokens the user wants to buy, respectively.
2. Replace the existing bid by a new one only with higher price and/or more quantities
    - The bidder can replace its existing bid, which is previously placed,  by a new one with the same option between Option A and Option B.

A bidder cannot do the following behaviors during the auction period:

1. Cancel the existing bid
2. Replace the existing bid by a new one with lower price or fewer quantities.

#### When the auction ends

The auction will end when the auction end time is arrived.

#### How the offering price is determined

Once the auction period ends, the bids are ordered in descending order of the bid prices to determine the offering price. The offering price is determined by finding the lowest price among the bid prices satisfying that the total amount of receiving tokens placed at more than or equal to the price is less the entire offering quantity.
The bidders who placed at the higher price than the offering price become the winning bidders and get the receiving tokens at the same price, which is the offering price. 