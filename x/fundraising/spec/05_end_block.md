<!-- order: 5 -->

At the end of each end block, the `fundraising` module operates the following executions based on auction type.

## Auction Status Transition

The module first gets all auctions registered in the store and proceed operations depending on auction status.

If the auction status is `AuctionStatusStandBy` and if the start time of the auction is passed, the auction status is updated to `AuctionStatusStarted`. 

For a batch auction, if the auction status is `AuctionStatusStarted` and if an end time of `EndTimes` of the auction is arrived yet, `MatchedPrice` is calculated and the matched bids that have the bid price higher than or equal to `MatchedPrice` are counted. According to `MaxExtendedRound` and `ExtendedRate`, whether the auction ends or the auction is extended with another extended round is determined. 


If the auction status is `AuctionStatusStarted` and if the last end time of the auction is arrived, or,
if `RemainingSellingCoin` is equal to zero for a fixed price auction, 
- the auction status is updated to `AuctionStatusVesting`,
- a list of `VestingQueue` is generated according to `VestingSchedules`,
- `MatchedPrice` is calculated and updated for the auction,
- the amount of `SellingCoin` is released from `SellingReserveAddress` to each matched bidders,
- if `SellingCoin` is not sold out, the remaining selling coin is sent from `SellingReserveAddress` to `Auctioneer`,
- the amount of `PayingCoin` corresponding to the amount of the sold `SellingCoin` is reserved in `VestingReserveAddress` from `PayingReserveAddress`, and 
- the remaining amount of `PayingCoin` in `PayingReserveAddress` is refunded from `PayingReserveAddress` to the bidders.


If the auction status is `AuctionStatusVesting` and if the last release time of the vesting schedule is arrived, the auction status is updated to `AuctionStatusFinished`.



## Calculation the Matched Price, Distribution of Selling Coins, and Refund of Paying Coins

This section provides the exact mathematical calculation of the matched price. For an illustrative example, see [here](../../../docs/Tutorials/demo/README.md).

The parameters to be calculated are defined as

- `X` : the matched price of a selling coin; This is `MatchedPrice` of the auction.
- `S_n` : total amount of selling coins to be distributed to the `n`-th bidder; This is how many selling coins will be distributed to the `n`-th bidder
- `R_n` : remaining amount of reserved paying coins of the `n`-th bidder to be refunded; This is how many paying coins will be refunded to the `n`-th bidder


For a fixed price auction, the above parameters for matched bids are calculated as
  - `X` = `MatchedPrice` = `BidPrice` = `StartPrice`,
  - `S_n` = `PayingCoin`/`MatchedPrice`,
  - `R_n` = 0.
  
For a batch auction, how to calculate the above parameters are described below.

To calculate `X`, `S_n(X)`, and `R_n(X)`, the given parameters from the auction and the bids are:

- `S` : total amount of `SellingCoin` for the auction
- `N` : total number of bidders for the auction
- `n` : index of a bidder
- `B_n` : total number of bids that the `n`-th bidder places for the auction
- `S^max_n` : max bid amount on how many `SellingCoin` the `n`-th bidder can get from this auction
- `b` : index of a bid of the `n`-th bidder
- `I_{n,b}(X)` : matched bid checker to check if `BidPrice` of the `b`-th bid of the `n`-th bidder is greater than or equal to `X`
    - If `BidPrice` of the `b`-th bid of the `n`-th bidder is greater than or equal to `X`, then `I_{n,b}(X)` =1.
    - If `BidPrice` of the `b`-th bid of the `n`-th bidder is less than `X`, then `I_{n,b}(X)` =0.
- `T_{n,1}` : set of `BidTypeFixedPrice` bids of the `n`-th bidder; This is only for a fixed price auction.
- `T_{n,2}` : set of `BidTypeBatchWorth` bids of the `n`-th bidder; This is only for a batch auction.
- `T_{n,3}` : set of `BidTypeBatchMany` bids of the `n`-th bidder; This is only for a batch auction.
- `P_{n,b}` : amount of `PayingCoin` reserved for the `b`-th bid of the `n`-th bidder;
- `S_{n,b}` : amount of selling coins to be released to the `n`-th bidder for the `b`-th bid
    - For a `BidTypeBatchWorth` bid, this varies according to the price and `P_{n,b}`=`S_{n,b}`  	&times; `MatchedPrice`.
    - For a `BidTypeBatchMany` bid, this is what the bidder places and `P_{n,b}`=`S_{n,b}`  	&times; `BidPrice` &le; `S_{n,b}`  	&times; `MatchedPrice` 


## Matched Price for Batch Auction

The matched price, `X`, (which is `MatchedPrice`&ge; `BidPrice`) is determined as the minimum price among the ordered prices that satisfy the following inequality. 

<p align="center"><img src="https://render.githubusercontent.com/render/math?math=\displaystyle \sum_n \max \left( \sum_{b \in T_{n,2}} \frac{P_{n, b}}{X}\cdot I_{n,b}(X) %2B \sum_{b \in T_{n,3}} S_{n, b}\cdot I_{n,b}(X) , {S^{\text{max}}_n} \right) \leq S "></p>

<!--- Plus sign should be replaced by %2B in math here. -->

## Distribution of Selling Coins

The amount `S_n` of selling coins to be distributed to the `n`-th bidder is calculated as

<p align="center"><img src="https://render.githubusercontent.com/render/math?math=\displaystyle S_n=\max \left(  \sum_{b \in T_{n,2}} \frac{P_{n, b}}{X}\cdot I_{n,b}(X) %2B \sum_{b \in T_{n,3}} S_{n, b}\cdot I_{n,b}(X), {S^{\text{max}}_n} \right). "></p>

## Refund of Paying Coins

The remaining amount `R_n` of reserved paying coins to be refundedl to the `n`-th bidder is calculated as

<p align="center"><img src="https://render.githubusercontent.com/render/math?math=\displaystyle R_n=\sum_{\text{all }\,  b} P_{n, b} - S_n \cdot X . "></p>




