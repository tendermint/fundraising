<!-- order: 3 -->

# State Transitions

This document describes the state transaction operations in the fundraising module.

## Coin Reservation for Fundraising Module Message

Transaction confirmation causes state transition on the Cosmos SDK [x/bank] module. Some messages on the fundraising module require coin reservation before confirmation.

The coin reserve processes for each message type are:

### MsgCreateFixedPriceAuction

When `MsgCreateFixedPriceAuction` is confirmed, 
- `SellingCoin` of the auctioneer is reserved in `SellingReserveAddress`, and 
- the created auction status is set to `AuctionStatusStandBy`.

### MsgCreateBatchAuction

When `MsgCreateBatchAuction` is confirmed,
- `SellingCoin` of the auctioneer is reserved in `SellingReserveAddress`, and
- the created auction status is set to `AuctionStatusStandBy`.

### MsgCancelAuction

When `MsgCancelAuction` is confirmed for the auction in `AuctionStatusStandBy`,
- `SellingCoin` reserved in `SellingReserveAddress` is refunded to the auctioneer, and
- the auction status is changed from `AuctionStatusStandBy` to `AuctionStatusCancelled`.

### MsgPlaceBid

When `MsgPlaceBid` is confirmed, `PayingCoin` of the bidder is reserved in `PayingReserveAddress`.

For a fixed price auction, when `MsgPlaceBid` is confirmed, `RemainingSellingCoin` is updated based on the message.

### MsgModifyBid

When `MsgModifyBid` is confirmed for an existing bid, the difference of `PayingCoin`  of the modifying bid and `PayingCoin`  of the existing bid is reserved in `PayingReserveAddress`.
