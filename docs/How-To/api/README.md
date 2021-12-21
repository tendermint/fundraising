---
Title: REST APIs
Description: A high-level overview of gRPC-gateway REST routes in fundraising module.
---

# API Reference 

This document provides a high-level overview of what gRPC-gateway REST routes are supported in the fundraising module. 

## Swagger Documentation

## gRPC-gateway REST Routes

To test out the following command line interface, you must set up a local network. By simply running `make localnet` under the root project directory, you can start the local network. It requires [Starport](https://starport.com/), but if you don't have Starport set up in your local machine, see this [install Starport guide](https://docs.starport.network/) to install it.  

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
        "auctioneer": "cosmos1k64etk30a60sv9hqla4qctwur65acq9zjdd8jt",
        "selling_pool_address": "cosmos18xzvtd72y9j8xyf8a36z5jjhth7qgtcwhh8lz7yee3tvxqn6ll5quh78zq",
        "paying_pool_address": "cosmos18permjyqvk5flft8ey9egr7hd4ry8tauqt4f9mg9knn4vvtkry9sujucrl",
        "start_price": "1.000000000000000000",
        "selling_coin": {
          "denom": "denom1",
          "amount": "10000000000"
        },
        "paying_coin_denom": "denom2",
        "vesting_pool_address": "cosmos1gukaqt783nhz79uhcqklsty7lc7jfyy8scn5ke4x7v0m3rkpt4dst7y4l3",
        "vesting_schedules": [
          {
            "release_time": "2022-01-01T00:00:00Z",
            "weight": "0.500000000000000000"
          },
          {
            "release_time": "2022-12-01T00:00:00Z",
            "weight": "0.500000000000000000"
          }
        ],
        "winning_price": "0.000000000000000000",
        "remaining_coin": {
          "denom": "denom1",
          "amount": "9995000000"
        },
        "start_time": "2021-12-01T00:00:00Z",
        "end_times": [
          "2021-12-30T00:00:00Z"
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
      "auctioneer": "cosmos1k64etk30a60sv9hqla4qctwur65acq9zjdd8jt",
      "selling_pool_address": "cosmos18xzvtd72y9j8xyf8a36z5jjhth7qgtcwhh8lz7yee3tvxqn6ll5quh78zq",
      "paying_pool_address": "cosmos18permjyqvk5flft8ey9egr7hd4ry8tauqt4f9mg9knn4vvtkry9sujucrl",
      "start_price": "1.000000000000000000",
      "selling_coin": {
        "denom": "denom1",
        "amount": "10000000000"
      },
      "paying_coin_denom": "denom2",
      "vesting_pool_address": "cosmos1gukaqt783nhz79uhcqklsty7lc7jfyy8scn5ke4x7v0m3rkpt4dst7y4l3",
      "vesting_schedules": [
        {
          "release_time": "2022-01-01T00:00:00Z",
          "weight": "0.500000000000000000"
        },
        {
          "release_time": "2022-12-01T00:00:00Z",
          "weight": "0.500000000000000000"
        }
      ],
      "winning_price": "0.000000000000000000",
      "remaining_coin": {
        "denom": "denom1",
        "amount": "9995000000"
      },
      "start_time": "2021-12-01T00:00:00Z",
      "end_times": [
        "2021-12-30T00:00:00Z"
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
| eligible          | The eligible   | {endpoint}/cosmos/fundraising/v1beta1/auctions/1/bids?eligible=false |
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
      "height": "230",
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
      "vested": false
    },
    {
      "auction_id": 1,
      "auctioneer": "cosmos1m4ys0e222x45657hrg9y2gadfxtcqja270rdkg",
      "paying_coin": "denom2",
      "release_time": "2022-12-01T00:00:00Z",
      "vested": false
    }
  ]
}
```