# Demo

This is to provide some examples of how the auctions are going through by using the `fundraising` module.

## Fixed Price Auction 
A fixed price auction is to sell a given amount of coins on a first-come, first-served basis.
In other words, once a bidder places a bid with the pre-determined bid price first (i.e., `StartPrice`), then the bidder can get the selling coins. 
When all the selling coins are sold out, then the auction becomes ended. 

For example, suppose that the amount of `SellingCoin` is 200 and its denom is CoinS. 
From Bidder 1 to Bidder 6, the bidders place bids with 10, 40, 30, 20, 70, and 30 CoinS, respectively.
Then, when Bidder 6's bid is confirmed, all the `SellingCoin` are sold out and, therefore, the auction becomes ended.
This auction process is illustrated in the following figure.

![alt text][fixedPriceAuction_example]


## Batch Auction
A batch auction allows each bidder to participate in the auction by placing limit orders with a bid price chosen freely at any time within the auction period. 
An order book is created to record the bids with various bid prices.

The following figure illustrates how the bids in a batch auction are placed. 
![alt text][batchAuction_bid_example]


When the end time of the auction are arrived, the calculation of the matched price are performed.
According to the matched price, the matched bids (and also the matched bidders) are determined. 
Regardless of the bid price, all the matched bidders can get the selling coins with the same price, which is the matched price.

The following illustrates how the matched price is calculated.
![alt text][batchAuction_cal_example]

For the details on the calculation, please see [here](../../../x/fundraising/spec/05_end_block.md).


[fixedPriceAuction_example]: ./figures/FixedPriceAuction_bid_example.gif "Example of Fixed Price Auction"
[batchAuction_bid_example]: ./figures/BatchAuction_bid_example.gif "Example of Bidding to Batch Auction"
[batchAuction_cal_example]: ./figures/BatchAuction_cal_example.gif "Example of Matched Price Calculation in Batch Auction"
