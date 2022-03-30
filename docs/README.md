# Documentation

How to use the fundraising module documentation.


- [Documentation](#documentation)
  - [Overview](#overview)
  - [More Documentations](#more-documentations)

## Overview

The main purpose of the `fundraising` module is to sell a certain amount of selling coins with a proper price. 
The characteristics of how to determine the matched price, the matched bids, and the matched selling coins differentiate the auction types. 
The outline of the flow of progressing an auction is following.

1. Create an auction
    - An auctioneer creates an auction.
2. Add bidder(s) to the list of the allowed bidders
    - The auctioneer adds bidder(s) as the allowed bidders that enables to participates in the auction.
3. Place/modify a bid
    - The bidders added in the list of the allowed bidders place bids and modify the bids according to the auction types.
4. Calculation of the matched price, the matched bids, and the matched selling coins
    - The matched price, the matched bids, and the matched selling coins are calculated based on the placed bids and the auction type.
5. Allocation and refund coins
    - The matched selling coins are distributed to the matched bidders.
    - The remaining selling coins are refunded to the auctioneer.
    - The matched paying coins are reserved in the vesting address.
    - The remaining paying coins are refunded to the bidders.
6. Vesting the matched paying coins
    - According to the vesting schedule, the pre-configured amount of the matched paying coins are vested to the auctioneer.



## More Documentations
The following documentations further provide the explanations on `fundraising` module.

* [How-Tos](./How-To/README.md)
   - How to use API and CLI
* [Tutorials](./Tutorials/README.md)
  - How to proceed with the auction and how to calculate the matched price 