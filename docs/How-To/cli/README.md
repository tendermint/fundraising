---
Title: Fundraisingd
Description: A high-level overview of how the command-line interface (CLI) works for the fundraising module.
---

# CLI Reference

This document provides a high-level overview of how the command line (CLI) interface works for the fundraising module. The executable file name is called `fundraisingd`.

## Command Line Interface

To test out the following commands, you must set up a local network. By simply running `$ make localnet` under the root project directory, you can start the local network. It requires the latest [Starport](https://starport.com/). If you don't have `Starport` set up in your local machine, see this [Starport guide](https://docs.starport.network/) to install it.  

- [CLI Reference](#cli-reference)
  - [Command Line Interface](#command-line-interface)
  - [Transaction](#transaction)
    - [CreateFixedPriceAuction](#createfixedpriceauction)
    - [CreateBatchAuction](#createbatchauction)
    - [CancelAuction](#cancelauction)
    - [AddAllowedBidder](#addallowedbidder)
    - [PlaceBid](#placebid)
  - [Query](#query)
    - [Params](#params)
    - [Auctions](#auctions)
    - [Auction](#auction)
    - [Bids](#bids)
    - [Vestings](#vestings)

## Transaction

+++ https://github.com/tendermint/fundraising/blob/main/proto/fundraising/tx.proto#L14-L29

### CreateFixedPriceAuction

An auctioneer can create a fixed price auction by setting the following parameters. In a fixed price auction, `start_price` is the matched price and bidders can buy the selling coins on a first-come, first-served basis. See the [spec](https://github.com/tendermint/fundraising/blob/main/x/fundraising/spec/01_concepts.md#auction-types) for a detailed and technical information about a fixed priced auction type.

| **Field**         |  **Description**                                                                    |
| :---------------- | :---------------------------------------------------------------------------------- |
| allowed_bidders | The list of allowed bidders that can participate in the auction, with a maximum possible bid amount for each bidder. It is empty when an auction is created. The module is designed to delegate permission to an external module to add its allowed bidders to the auction. |
| start_price       | The starting price of the selling coin; it is proportional to the paying coin denom. This is the matched price. | 
| selling_coin      | The selling amount of coin for the auction                                      | 
| paying_coin_denom | The paying coin denom that bidders use to bid with                                  | 
| vesting_schedules | The vesting schedules that release the paying coins to the autioneer                | 
| start_time        | The start time of the auction                                                       | 
| end_time          | The end time of the auction                                                         | 
|                   |                                                                                     |

Example of input as JSON:

```json
{
  "start_price": "1.000000000000000000",
  "selling_coin": {
    "denom": "denom1",
    "amount": "1000000000000"
  },
  "paying_coin_denom": "denom2",
  "vesting_schedules": [
    {
      "release_time": "2022-06-21T00:00:00Z",
      "weight": "0.500000000000000000"
    },
    {
      "release_time": "2022-12-21T00:00:00Z",
      "weight": "0.500000000000000000"
    }
  ],
  "start_time": "2022-02-01T00:00:00Z",
  "end_time": "2022-03-01T00:00:00Z"
}
```
the auctioneer can create a fixed price auction to `1000,000,000,000denom1` coin, where the starting price is `2.0` which means that the price of `1denom1` is `2denom2`. The auction starts at `2022-02-01T00:00:00Z` and ends at `2022-03-01T00:00:00Z`.

Example command:

```bash
# Create a fixed price auction
fundraisingd tx fundraising create-fixed-price-auction auction.json \
--chain-id fundraising \
--from bob \
--keyring-backend test \
--broadcast-mode block \
--yes \
--output json | jq
```


### CreateBatchAuction

An auctioneer can create a batch auction by setting the following parameters. Differently from a fixed price auction,  start_price does not affect the determination of the matched price, but is provided by the auctioneer as a reference price to bidders. See the [spec](https://github.com/tendermint/fundraising/blob/main/x/fundraising/spec/01_concepts.md#auction-types) for a detailed and technical information about a batch auction type.


| **Field**         |  **Description**                                                                    |
| :---------------- | :---------------------------------------------------------------------------------- |
| allowed_bidders | The list of allowed bidders that can participate in the auction, with a maximum possible bid amount for each bidder. It is empty when an auction is created. The module is designed to delegate permission to an external module to add its allowed bidders to the auction. |
| start_price       | The starting price of the selling coin; it is proportional to the paying coin denom. This is the matched price. | 
|min_bid_price | The minimum bid price that bidders must place with.|
| selling_coin      | The selling amount of coin for the auction                                      | 
| paying_coin_denom | The paying coin denom that bidders use to bid with                                  | 
| vesting_schedules | The vesting schedules that release the paying coins to the autioneer                | 
| start_time        | The start time of the auction                                                       | 
| end_times          | The list of the end times of the auction in consideration of the extended rounds                                                         | 
| max_extended_round   | The maximum number of extended rounds that provides additional opportunity for the bidders to place bids when more than a certain ratio of the number of the matched bids are reduced compared to the previous end time  |
| extended_round_rate | The threshold reduction of the number of the matched bids are reduced compared to the previous end time to decide the necessity of another extended round |              

Example of input as JSON:

```json
{
  "allowed_bidders": [],
	"start_price": "2.000000000000000000",
  "min_bid_price": "0.100000000000000000",
  "selling_coin": {
    "denom": "denom1",
    "amount": "1000000000000"
  },
  "paying_coin_denom": "denom2",
  "vesting_schedules": [
    {
      "release_time": "2022-06-21T00:00:00Z",
      "weight": "0.500000000000000000"
    },
    {
      "release_time": "2022-12-21T00:00:00Z",
      "weight": "0.500000000000000000"
    }
  ],
  "start_time": "2022-02-01T00:00:00Z",
  "end_times": [
		"2022-03-01T00:00:00Z",
		"2022-03-02T00:00:00Z", 
		"2022-03-03T00:00:00Z", 
		"2022-03-04T00:00:00Z"
	],
	"max_extended_round": "3",
	"extended_round_rate": "0.05"
}
```


Example command:

```bash
# Create a batch auction
fundraisingd tx fundraising create-batch-auction auction-batch.json \
--chain-id fundraising \
--from bob \
--keyring-backend test \
--broadcast-mode block \
--yes \
--output json | jq
```


### CancelAuction

This command is useful for an auctioneer when the auctioneer made mistake(s) on some values of the auction. The module doesn't support update functionality. Instead, the module allows them to cancel an auction and recreate it with correct values. Note that it can only be cancelled when the auction has not started yet.

Example command:

```bash
# Cancel the auction
fundraisingd tx fundraising cancel 1 \
--chain-id fundraising \
--from bob \
--keyring-backend test \
--broadcast-mode block \
--yes \
--output json | jq
```


### AddAllowedBidder

**Important Note**: the `fundraising` module is designed in a way that all auctions are closed when they are created. It means that no one can place a bid unless they are allowed. The module expects an external module (a module that imports and uses the `fundraising` module) to control a list of allowed bidder for an auction. There are functions, such as `AddAllowedBidders()` and `UpdateAllowedBidder()` implemented for the external module to use. 

For testing purpose, there is a custom message called `MsgAddAllowedBidder`. It adds a single allowed bidder for the auction and this message is only available when you build `fundraisingd` with `config-test.yml` file. Running `make localnet` is automatically using `config-test.yml`. Under the hood, a custom `enableAddAllowedBidder` ldflags is passed to build configuration in `config-test.yml` file.

Example command:

```bash
# Add steve's address to allowed bidder list
fundraisingd tx fundraising add-allowed-bidder 1 1000000000 \
--chain-id fundraising \
--from steve \
--keyring-backend test \
--broadcast-mode block \
--yes \
--output json | jq
```



### PlaceBid
This is for a bidder to place a new bid to the auction, where the bidder should be in the list of the allowed bidders. 

Usage
```bash
bid [auction-id] [bid-type] [price] [coin]
```

| **Argument**      |  **Description**                     |
| :---------------- | :----------------------------------- |
| auction-id        | auction ID that the bid corresponds to. | 
| bid-type  | bid type among 1) fixed-price (fp or f), 2) batch-worth (bw or w), and 3) batch-many  (bm or m), where 1) is only for `FixedPriceAuction` and 2)&3) are only for `BatchAuction`.|
| price     | bid price (dec type) of a selling coin as the unit of a paying coin. For fixed-price type, this price must be the same as `StartPrice` of the auction. For batch-worth and batch-many, this price must be higher than or equal to `MinBidPrice` of the auction. | 
| coin      | how many coins to bid, where the denom should be of the paying coin for the bid types of fixed-price and batch-worth, and of the selling coin for the bid type of batch-many.|

Example command:

```bash
fundraisingd tx fundraising bid 1 fixed-price 1.0 5000000denom2 \
--chain-id fundraising \
--from steve \
--keyring-backend test \
--broadcast-mode block \
--yes \
--output json | jq
```


## Query

+++ https://github.com/tendermint/fundraising/blob/main/proto/fundraising/query.proto#L14-L42


### Params 

```bash
# Query the values set as fundraising parameters
fundraisingd q fundraising params --output json | jq
```



### Auctions

Example commands:

```bash
# Query for all auuctions on a network
fundraisingd q fundraising auctions \
-o json | jq

# Query for all auctions with the given auction status
# Ref: https://github.com/tendermint/fundraising/blob/main/x/fundraising/spec/02_state.md#auction-status
fundraisingd q fundraising auctions \
--status AUCTION_STATUS_STANDBY \
-o json | jq

# Query for all auctions with the given auctino type
# Ref: https://github.com/tendermint/fundraising/blob/main/x/fundraising/spec/02_state.md#auction-type
fundraisingd q fundraising auctions \
--type AUCTION_TYPE_FIXED_PRICE \
-o json | jq
```

Result:

```json
{
  "auctions": [
    {
      "@type": "/tendermint.fundraising.FixedPriceAuction",
      "base_auction": {
        "id": "1",
        "type": "AUCTION_TYPE_FIXED_PRICE",
        "allowed_bidders": [
          {
            "bidder": "cosmos1tfzynkllgxdpmrcknx2j5d0hj9zd82tceyfa5n",
            "max_bid_amount": "1000000000"
          }
        ],
        "auctioneer": "cosmos1dncsflcfknkmlmt3t6836tkd3mu742e2wh4r70",
        "selling_reserve_address": "cosmos1wl90665mfk3pgg095qhmlgha934exjvv437acgq42zw0sg94flestth4zu",
        "paying_reserve_address": "cosmos17gk7a5ys8pxuexl7tvyk3pc9tdmqjjek03zjemez4eqvqdxlu92qdhphm2",
        "start_price": "1.000000000000000000",
        "selling_coin": {
          "denom": "denom1",
          "amount": "1000000000000"
        },
        "paying_coin_denom": "denom2",
        "vesting_reserve_address": "cosmos1q4x4k4qsr4jwrrugnplhlj52mfd9f8jn5ck7r4ykdpv9wczvz4dqe8vrvt",
        "vesting_schedules": [
          {
            "release_time": "2022-06-21T00:00:00Z",
            "weight": "0.500000000000000000"
          },
          {
            "release_time": "2022-12-21T00:00:00Z",
            "weight": "0.500000000000000000"
          }
        ],
        "winning_price": "0.000000000000000000",
        "remaining_coin": {
          "denom": "denom1",
          "amount": "999995000000"
        },
        "start_time": "2022-02-01T00:00:00Z",
        "end_times": [
          "2022-03-01T00:00:00Z"
        ],
        "status": "AUCTION_STATUS_STARTED"
      }
    }
  ],
  "pagination": {
    "next_key": null,
    "total": "1"
  }
}
```

### Auction

Example command:

```bash
# Query for the specific auction with the auction id
fundraisingd q fundraising auction 1 -o json | jq
```

Result:

```json
{
  "auction": {
    "@type": "/tendermint.fundraising.FixedPriceAuction",
    "base_auction": {
      "id": "1",
      "type": "AUCTION_TYPE_FIXED_PRICE",
      "allowed_bidders": [
        {
          "bidder": "cosmos1tfzynkllgxdpmrcknx2j5d0hj9zd82tceyfa5n",
          "max_bid_amount": "1000000000"
        }
      ],
      "auctioneer": "cosmos1dncsflcfknkmlmt3t6836tkd3mu742e2wh4r70",
      "selling_reserve_address": "cosmos1wl90665mfk3pgg095qhmlgha934exjvv437acgq42zw0sg94flestth4zu",
      "paying_reserve_address": "cosmos17gk7a5ys8pxuexl7tvyk3pc9tdmqjjek03zjemez4eqvqdxlu92qdhphm2",
      "start_price": "1.000000000000000000",
      "selling_coin": {
        "denom": "denom1",
        "amount": "1000000000000"
      },
      "paying_coin_denom": "denom2",
      "vesting_reserve_address": "cosmos1q4x4k4qsr4jwrrugnplhlj52mfd9f8jn5ck7r4ykdpv9wczvz4dqe8vrvt",
      "vesting_schedules": [
        {
          "release_time": "2022-06-21T00:00:00Z",
          "weight": "0.500000000000000000"
        },
        {
          "release_time": "2022-12-21T00:00:00Z",
          "weight": "0.500000000000000000"
        }
      ],
      "winning_price": "0.000000000000000000",
      "remaining_coin": {
        "denom": "denom1",
        "amount": "999995000000"
      },
      "start_time": "2022-02-01T00:00:00Z",
      "end_times": [
        "2022-03-01T00:00:00Z"
      ],
      "status": "AUCTION_STATUS_STARTED"
    }
  }
}
```

### Bids

Example command:

```bash
# Query for all bids of the auction with the given auction id
fundraisingd q fundraising bids 1 \
-o json | jq
```

Result:

```json
{
  "bids": [
    {
      "auction_id": "1",
      "sequence": "1",
      "bidder": "cosmos1tfzynkllgxdpmrcknx2j5d0hj9zd82tceyfa5n",
      "price": "1.000000000000000000",
      "coin": {
        "denom": "denom2",
        "amount": "5000000"
      },
      "height": "1407",
      "eligible": true
    }
  ],
  "pagination": {
    "next_key": null,
    "total": "1"
  }
}
```

### Vestings

This command is used by an auctioneer to query vesting information. It only returns results when the auction is in vesting status

Example command:

```bash
# Query for all vesting queues 
fundraisingd q fundraising vestings 1 \
-o json | jq
```

Result:

```json
{
  "vesting_queues": [
    {
      "auction_id": 1,
      "auctioneer": "cosmos1m4ys0e222x45657hrg9y2gadfxtcqja270rdkg",
      "paying_coin": "denom2",
      "release_time": "2022-01-01T00:00:00Z",
      "released": false
    },
    {
      "auction_id": 1,
      "auctioneer": "cosmos1m4ys0e222x45657hrg9y2gadfxtcqja270rdkg",
      "paying_coin": "denom2",
      "release_time": "2022-12-01T00:00:00Z",
      "released": false
    }
  ]
}
```
