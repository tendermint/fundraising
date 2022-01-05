<!-- order: 1 -->

# Concepts

## Fundraising Module

The `x/fundraising` Cosmos SDK module provides an oppotunity for new projects to onboard the ecosystem. It allows the projects to raise funds, but also increase their brand awareness before launching their projects.

## Auction Types

There are two different types of auctions in the fundraising module.

### Fixed Price Auction

A fixed price auction type is the most basic type of an auction that is first come first served way to raise funds. When an autioneer creates this fixed price auction, they must determine the fixed start price that is proportional to the paying coin denom. Once it is created, bidders can only bid with the start price of the auction. For example, an auctioneer sets 0.5 as start price for X coin and the paying coin is Y coin. How many X coins does a bidder receives if they bids 10Y coin? It is calculated as Y coin over start price (10/0.5), which results to 20X coin.  

### English Auction

An english auction type is an ascending dynamic auction where the bidding starts with the starting price which is set by the auctioneer and increases with the continuous bidding from the different bidders until the end time. During the bidding time, bidders can bid to purchase the amount of coins and they can only increase the bid price. This auction type is not about how fast bidders bid for the coin but it is about price competition.

## Extended Auction Round(s)

The concept of extended auction round(s) is formed out of auction sniping. It is the technique in a timed online auction where a sniper waits until the last second to bid slightly above the current highest bid to purchase the majority of the selling coin. This gives other bidders no time to outbid the sniper and provides hard feelings among other bidders. Therfore, the extended auction round(s) is there to provent from the auction sniping. There is going to be one or more auction rounds if the extend rate is exceeded the previous winning price. 