---
Title: REST APIs
Description: A high-level overview of gRPC-gateway REST routes in fundraising module.
---

# API Reference 

This document provides a high-level overview of what gRPC-gateway REST routes are supported in the fundraising module. 

## Swagger Documentation

## gRPC-gateway REST Routes

To test out the following command line interface, you must set up a local network. By simply running `make localnet` under the root project directory, you can start the local network. It requires Ignite CLI, but if you don't have Ignite CLI set up in your local machine, see [this guide](https://docs.ignite.com/welcome/install) to install it.

* [Params](#Params)
* [Auctions](#Auctions)
* [Auction](#Auction)
* [AllowedBidders](#AllowedBidders)
* [AllowedBidder](#AllowedBidder)
* [Bids](#Bids)
* [Vestings](#Vestings)

## REST Routes

+++ https://github.com/tendermint/fundraising/blob/main/proto/fundraising/query.proto#L14-L42

### Params

Query the values set as fundraising parameters

Example endpoint: 

<!-- markdown-link-check-disable-next-line -->
http://localhost:1317/cosmos/fundraising/v1beta1/params

Result:

```json
{
  "auction_creation_fee": [
    {
      "denom": "stake",
      "amount": "100000000"
    }
  ],
  "extended_period": 1,
  "auction_fee_collector": "cosmos1t2gp44cx86rt8gxv64lpt0dggveg98y4ma2wlnfqts7d4m4z70vqrzud4t"
}
```

### Auctions

Query for all auuctions on a network

Example endpoint: 

<!-- markdown-link-check-disable-next-line -->
http://localhost:1317/cosmos/fundraising/v1beta1/auctions

| **Query String**  |  **Description**   | **Example** |
| :---------------- | :----------------- | :---------- |
| status            | The auction status | {endpoint}/cosmos/fundraising/v1beta1/auctions?type=AUCTION_TYPE_FIXED_PRICE |
| type              | The auction type   | {endpoint}/cosmos/fundraising/v1beta1/auctions?status=AUCTION_STATUS_STANDBY |

Result:

```json
{
  "auctions": [
    {
      "@type": "/tendermint.fundraising.FixedPriceAuction",
      "base_auction": {
        "id": "1",
        "type": "AUCTION_TYPE_FIXED_PRICE",
        "auctioneer": "cosmos1mzgucqnfr2l8cj5apvdpllhzt4zeuh2cshz5xu",
        "selling_reserve_address": "cosmos1wl90665mfk3pgg095qhmlgha934exjvv437acgq42zw0sg94flestth4zu",
        "paying_reserve_address": "cosmos17gk7a5ys8pxuexl7tvyk3pc9tdmqjjek03zjemez4eqvqdxlu92qdhphm2",
        "start_price": "2.000000000000000000",
        "selling_coin": {
          "denom": "denom1",
          "amount": "1000000000000"
        },
        "paying_coin_denom": "denom2",
        "vesting_reserve_address": "cosmos1q4x4k4qsr4jwrrugnplhlj52mfd9f8jn5ck7r4ykdpv9wczvz4dqe8vrvt",
        "vesting_schedules": [
          {
            "release_time": "2023-01-01T00:00:00Z",
            "weight": "0.500000000000000000"
          },
          {
            "release_time": "2023-06-01T00:00:00Z",
            "weight": "0.500000000000000000"
          }
        ],
        "start_time": "2022-05-01T00:00:00Z",
        "end_times": [
          "2022-06-01T00:00:00Z"
        ],
        "status": "AUCTION_STATUS_STARTED"
      },
      "remaining_selling_coin": {
        "denom": "denom1",
        "amount": "999997500000"
      }
    },
    {
      "@type": "/tendermint.fundraising.BatchAuction",
      "base_auction": {
        "id": "2",
        "type": "AUCTION_TYPE_BATCH",
        "auctioneer": "cosmos1mzgucqnfr2l8cj5apvdpllhzt4zeuh2cshz5xu",
        "selling_reserve_address": "cosmos197ewwasd96k2fh3nx5m76zvqxpzjcxuyq65rwgw0aa2edmwafgfqfa5qqz",
        "paying_reserve_address": "cosmos1s3cspws3lsqfvtjcz9jvpx7kjm93npmwjq8p4xfu3fcjj5jz9pks20uja6",
        "start_price": "0.300000000000000000",
        "selling_coin": {
          "denom": "denom1",
          "amount": "1000000000000"
        },
        "paying_coin_denom": "denom2",
        "vesting_reserve_address": "cosmos1pye9kv5f8s9n8uxnr0uznsn3klq57vqz8h2ya6u0v4w5666lqdfqjrw0qu",
        "vesting_schedules": [
          {
            "release_time": "2023-01-01T00:00:00Z",
            "weight": "0.500000000000000000"
          },
          {
            "release_time": "2023-06-01T00:00:00Z",
            "weight": "0.500000000000000000"
          }
        ],
        "start_time": "2022-05-01T00:00:00Z",
        "end_times": [
          "2022-06-01T00:00:00Z"
        ],
        "status": "AUCTION_STATUS_STARTED"
      },
      "min_bid_price": "0.100000000000000000",
      "matched_price": "0.000000000000000000",
      "max_extended_round": 2,
      "extended_round_rate": "0.150000000000000000"
    }
  ],
  "pagination": {
    "next_key": null,
    "total": "2"
  }
}
```

### Auction

Query for the specific auction with the auction id

Example endpoint: 

<!-- markdown-link-check-disable-next-line -->
http://localhost:1317/cosmos/fundraising/v1beta1/auctions/1

Result:

```json
{
  "auction": {
    "@type": "/tendermint.fundraising.FixedPriceAuction",
    "base_auction": {
      "id": "1",
      "type": "AUCTION_TYPE_FIXED_PRICE",
      "auctioneer": "cosmos1mzgucqnfr2l8cj5apvdpllhzt4zeuh2cshz5xu",
      "selling_reserve_address": "cosmos1wl90665mfk3pgg095qhmlgha934exjvv437acgq42zw0sg94flestth4zu",
      "paying_reserve_address": "cosmos17gk7a5ys8pxuexl7tvyk3pc9tdmqjjek03zjemez4eqvqdxlu92qdhphm2",
      "start_price": "2.000000000000000000",
      "selling_coin": {
        "denom": "denom1",
        "amount": "1000000000000"
      },
      "paying_coin_denom": "denom2",
      "vesting_reserve_address": "cosmos1q4x4k4qsr4jwrrugnplhlj52mfd9f8jn5ck7r4ykdpv9wczvz4dqe8vrvt",
      "vesting_schedules": [
        {
          "release_time": "2023-01-01T00:00:00Z",
          "weight": "0.500000000000000000"
        },
        {
          "release_time": "2023-06-01T00:00:00Z",
          "weight": "0.500000000000000000"
        }
      ],
      "start_time": "2022-05-01T00:00:00Z",
      "end_times": [
        "2022-06-01T00:00:00Z"
      ],
      "status": "AUCTION_STATUS_STARTED"
    },
    "remaining_selling_coin": {
      "denom": "denom1",
      "amount": "999997500000"
    }
  }
}
```

### Bids

Query for all bids of the auction with the given auction id

Example endpoint: 

<!-- markdown-link-check-disable-next-line -->
http://localhost:1317/cosmos/fundraising/v1beta1/auctions/1/bids

| **Query String**  |  **Description**   | **Example** |
| :---------------- | :----------------- | :---------- |
| bidder            | The bidder address | {endpoint}/cosmos/fundraising/v1beta1/auctions/1/bids?bidder=cosmos1mc60p3ul372mepchm9shd9r456kur958t4v8ld |
| is_matched        | The matched status | {endpoint}/cosmos/fundraising/v1beta1/auctions/1/bids?is_matched=true |

Result:

```json
{
  "bids": [
    {
      "auction_id": "1",
      "bidder": "cosmos1mzgucqnfr2l8cj5apvdpllhzt4zeuh2cshz5xu",
      "id": "1",
      "type": "BID_TYPE_FIXED_PRICE",
      "price": "2.000000000000000000",
      "coin": {
        "denom": "denom2",
        "amount": "5000000"
      },
      "is_matched": true
    }
  ],
  "pagination": {
    "next_key": null,
    "total": "1"
  }
}
```

### AllowedBidders

Query for all allowed bidders list for the auction

Example endpoint:

<!-- markdown-link-check-disable-next-line -->
http://localhost:1317/cosmos/fundraising/v1beta1/auctions/1/allowed_bidders

Result:

```json
{
  "allowed_bidders": [
    {
      "bidder": "cosmos185fflsvwrz0cx46w6qada7mdy92m6kx4gqx0ny",
      "max_bid_amount": "5000000000"
    },
    {
      "bidder": "cosmos1mzgucqnfr2l8cj5apvdpllhzt4zeuh2cshz5xu",
      "max_bid_amount": "1000000000"
    }
  ],
  "pagination": {
    "next_key": null,
    "total": "2"
  }
}
```


### AllowedBidders

Query for a specific allowed bidder for the auction

Example endpoint:

<!-- markdown-link-check-disable-next-line -->
http://localhost:1317/cosmos/fundraising/v1beta1/auctions/1/allowed_bidders/cosmos185fflsvwrz0cx46w6qada7mdy92m6kx4gqx0ny

Result:

```json
{
  "allowed_bidder": {
    "bidder": "cosmos185fflsvwrz0cx46w6qada7mdy92m6kx4gqx0ny",
    "max_bid_amount": "5000000000"
  }
}
```

### Vestings

Query for all vesting queues 

Example endpoint: 

<!-- markdown-link-check-disable-next-line -->
http://localhost:1317/cosmos/fundraising/v1beta1/auctions/1/vestings

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