---
Title: REST APIs
Description: A high-level overview of gRPC-gateway REST routes in fundraising module.
---

# API Reference 

This document provides a high-level overview of what gRPC-gateway REST routes are supported in the fundraising module. 

## Swagger Documentation

- Swagger Cosmos SDK Fundraising Module [REST and gRPC Gateway docs](https://app.swaggerhub.com/apis-docs/gravity-devs/fundraising/v0.1.0)

## gRPC-gateway REST Routes

To test out the following command line interface, you must set up a local network. By simply running `make localnet` under the root project directory, you can start the local network. It requires Ignite CLI, but if you don't have Starport set up in your local machine, see this [install Starport guide](https://docs.ignite.com/#install-starport) to install it.  

* [Params](#Params)
* [Auctions](#Auctions)
* [Auction](#Auction)
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
|                   |                    |

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

Query for all bids of the auction with the given auction id

Example endpoint: 

<!-- markdown-link-check-disable-next-line -->
http://localhost:1317/cosmos/fundraising/v1beta1/auctions/1/bids

| **Query String**  |  **Description**   | **Example** |
| :---------------- | :----------------- | :---------- |
| bidder            | The bidder address | {endpoint}/cosmos/fundraising/v1beta1/auctions/1/bids?bidder=cosmos1mc60p3ul372mepchm9shd9r456kur958t4v8ld |
| eligible          | The eligible       | {endpoint}/cosmos/fundraising/v1beta1/auctions/1/bids?eligible=false |
|                   |                    |

Result:

```json
{
  "bids": [
    {
      "auction_id": "1",
      "sequence": "1",
      "bidder": "cosmos1mc60p3ul372mepchm9shd9r456kur958t4v8ld",
      "price": "1.000000000000000000",
      "coin": {
        "denom": "denom2",
        "amount": "5000000"
      },
      "eligible": false
    }
  ],
  "pagination": {
    "next_key": null,
    "total": "1"
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