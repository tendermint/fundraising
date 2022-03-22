<!-- order: 7 -->

# Events

The `fundraising` module emits the following events:

## Handlers

### MsgCreateFixedPriceAuction

| Type                       | Attribute Key         | Attribute Value            |
| -------------------------- | --------------------- | -------------------------- |
| create_fixed_price_auction | auction_id            | {auctionId}                |
| create_fixed_price_auction | auctioneer_address    | {auctioneerAddress}        |
| create_fixed_price_auction | start_price           | {startPrice}               |
| create_fixed_price_auction | selling_pool_address  | {SellingReserveAddress}    |
| create_fixed_price_auction | paying_pool_address   | {PayingReserveAddress}     |
| create_fixed_price_auction | vesting_pool_address  | {VestingReserveAddress}    |
| create_fixed_price_auction | selling_coin          | {sellingCoin}              |
| create_fixed_price_auction | paying_coin_denom     | {payingCoinDenom}          |
| create_fixed_price_auction | start_time            | {startTime}                |
| create_fixed_price_auction | end_time              | {endTime}                  |
| create_fixed_price_auction | auction_status        | {auctionStatus}            |
| message                    | module                | fundraising                |
| message                    | action                | create_fixed_price_auction |
| message                    | auctioneer            | {auctioneerAddress}        |

### MsgCreateBatchAuction

| Type                      | Attribute Key        | Attribute Value            |  
| ------------------------- | -------------------- | -------------------------- |
| create_batch_auction      | auction_id           | {auctionId}                |
| create_batch_auction      | auctioneer_address   | {auctioneerAddress}        |
| create_batch_auction      | start_price          | {startPrice}               |
| create_batch_auction      | selling_pool_address | {SellingReserveAddress}    |
| create_batch_auction      | paying_pool_address  | {PayingReserveAddress}     |
| create_batch_auction      | vesting_pool_address | {VestingReserveAddress}    |
| create_batch_auction      | selling_coin         | {sellingCoin}              |
| create_batch_auction      | paying_coin_denom    | {payingCoinDenom}          |
| create_batch_auction      | start_time           | {startTime}                |
| create_batch_auction      | end_time             | {endTime}                  |
| create_batch_auction      | auction_status       | {auctionStatus}            |
| create_batch_auction      | min_bid_price        | {minBidPrice}              |
| create_batch_auction      | matched_price        | {matchedPrice}             |
| create_batch_auction      | max_extended_round   | {maxExtendedRound}         |
| create_batch_auction      | extended_round_rate  | {extendedRoundRate}        |
| message                   | module               | fundraising                |
| message                   | action               | create_batch_auction       |
| message                   | auctioneer           | {auctioneerAddress}        | 



### MsgCancelAuction

| Type           | Attribute Key | Attribute Value     |
| -------------- | ------------- | ------------------- |
| cancel_auction | auction_id    | {auctionId}         |
| message        | module        | fundraising         |
| message        | action        | cancel_auction      |
| message        | auctioneer    | {auctioneerAddress} | 

### MsgPlaceBid

| Type      | Attribute Key  | Attribute Value |
| --------- | -------------- | --------------- |
| place_bid | bidder_address | {bidderAddress} |
| place_bid | bid_price      | {bidPrice}      |
| place_bid | bid_coin       | {bidCoin}       |
| place_bid | bid_amount     | {bidAmount}     |
| message   | module         | fundraising     |
| message   | action         | place_bid       |
| message   | bidder         | {bidderAddress} | 
