<!-- order: 7 -->

# Events

The `fundraising` module emits the following events:

## EndBlocker

| Type                   | Attribute Key        | Attribute Value        |
| ---------------------- | -------------------- | ---------------------- |
| TBD  | TBD               | {TBD}                 |

## Handlers

### MsgCreateFixedPriceAuction

| Type                       | Attribute Key         | Attribute Value            |
| -------------------------- | --------------------- | -------------------------- |
| create_fixed_price_auction | auction_id            | {auctionId}                |
| create_fixed_price_auction | auctioneer_address    | {auctioneerAddress}        |
| create_fixed_price_auction | selling_pool_address  | {sellingPoolAddress}       |
| create_fixed_price_auction | paying_pool_address   | {payingPoolAddress}        |
| create_fixed_price_auction | vesting_pool_address  | {vestingPoolAddress}       |
| create_fixed_price_auction | start_price           | {startPrice}               |
| create_fixed_price_auction | selling_coin          | {sellingCoin}              |
| create_fixed_price_auction | vesting_schedules     | {vestingSchedules}         |
| create_fixed_price_auction | paying_coin_denom     | {payingCoinDenom}          |
| create_fixed_price_auction | auction_status        | {auctionStatus}            |
| create_fixed_price_auction | start_time            | {startTime}                |
| create_fixed_price_auction | end_time              | {endTime}                  |
| message                    | module                | fundraising                |
| message                    | action                | create_fixed_price_auction |
| message                    | auctioneer            | {auctioneerAddress}        | 

### MsgCreateEnglishAuction

| Type                      | Attribute Key        | Attribute Value            |  
| ------------------------- | -------------------- | -------------------------- |
| TBD                       | TBD                  | {TBD}                      |
| message                   | module               | fundraising                |
| message                   | action               | create_english_auction     |
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
| place_bid | price          | {price}         |
| place_bid | coin           | {coin}          |
| message   | module         | fundraising     |
| message   | action         | place_bid       |
| message   | bidder         | {bidderAddress} | 
