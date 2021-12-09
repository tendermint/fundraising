<!-- order: 1 -->

# Concepts

## Fundraising Module

The `x/fundraising` Cosmos SDK module provides an oppotunity for new projects to onboard the ecosystem. It does not only allow projects to raise funds, but also increase their brand awareness before launching their projects.

## Auction Types

There are two different types of auctions in the fundraising module.

### Fixed Price Auction

A fixed price auction type is the most basic type of an auction that is first come first served way to raise funds. When an autioneer creates this fixed price auction, they must determine the fixed start price that is proportional to the paying coin denom. Once it is created, bidders can only bid with the start price of the auction. For example, an auctioneer sets 0.5 as start price for X coin and the paying coin is Y coin. How many X coins does a bidder receives if they bids 10Y coin? It is calculated as Y coin over start price (10/0.5), which results to 20X coin.  

### English Auction

An english auction type is an ascending dynamic auction that an auctioneer decides the starting price of the selling amount of coin and bidders bid to purchase the amounts of the coin. Competition for the price, not how fast bidders can bid.

## Extended Auction Round(s)
