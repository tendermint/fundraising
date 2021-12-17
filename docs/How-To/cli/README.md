---
Title: Fundraisingd
Description: A high-level overview of how the command-line interface (CLI) works for the fundraising module.
---

# CLI Reference

This document provides a high-level overview of how the command line (CLI) interface works for the fundraising module. The executable name is called `fundraisingd`.

## Command Line Interface

To test out the following command line interface, you must set up a local network. By simply running `make localnet` under the root project directory, you can start the local network. It requires [Starport](https://starport.com/), but if you don't have Starport set up in your local machine, see this [install Starport guide](https://docs.starport.network/) to install it.  

- [Transaction](#Transaction)
    * [CreateFixedPriceAuction](#CreateFixedPriceAuction)
    * [CreateEnglishAuction](#CreateEnglishAuction)
    * [CancelAuction](#CancelAuction)
    * [PlaceBid](#PlaceBid)
- [Query](#Query)
    * [Params](#Params)
    * [Auctions](#Auctions)
    * [Auction](#Auction)
    * [Bids](#Bids)
    * [Vestings](#Vestings)

## Transaction

+++ https://github.com/tendermint/fundraising/blob/main/proto/fundraising/tx.proto#L14-L29

### CreateFixedPriceAuction

This command is one way for an auctioneer to create an auction to raise funds for their project. It is the most basic and simple type of an auction that has first come first served characteristic.
When an auctioneer creates a fixed price auction type, they must determine the fixed starting price. It is proportional to the paying coin denom that they want to get funded for. Once it is created, bidders start to bid with the same starting price and amount of coin. See the [spec](https://github.com/tendermint/fundraising/blob/main/x/fundraising/spec/01_concepts.md#auction-types) for a detailed information about the fixed price auction type.

JSON example:

In this JSON example, an auctioneer plans to create a fixed price auction that sells `denom1` coin with an amount of `1000000000000`, and the starting price is `1.0` that is proportional to the paying coin denom `denom2`. It means that the fixed starting price of `denom1` is the same as `denom2` price. The auction starts at `2021-12-10T00:00:00Z` and ends at `2021-12-10T00:00:00Z`. As soon as the auction starts, bidders can bid for the auction with the fixed start price and amount of coin that they are willing to bid. When it ends, the paying amount of coin that is reserved for all bids is expected to be released based on the vesting schedules. 

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
      "release_time": "2022-01-01T00:00:00Z",
      "weight": "0.500000000000000000"
    },
    {
      "release_time": "2022-12-01T00:00:00Z",
      "weight": "0.500000000000000000"
    }
  ],
  "start_time": "2021-12-10T00:00:00Z",
  "end_time": "2021-12-10T00:00:00Z"
}
```

Reference the description of each field:

| **Field**         |  **Description**                                                              |
| :---------------- | :---------------------------------------------------------------------------- |
| start_price       | The starting price of the selling coin, proportional to the paying coin denom | 
| selling_coin      | The selling amount of coin for the auction                                    | 
| paying_coin_denom | The paying coin denom that bidders use to bid with                            | 
| vesting_schedules | The vesting schedules that release the paying coins to the autioneer          | 
| start_time        | The start time of the auction                                                 | 
| end_time          | The end time of the auction                                                   | 
|                   |                                                                               | 

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

Result:

```result
{
  "@type": "/cosmos.tx.v1beta1.Tx",
  "body": {
    "messages": [
      {
        "@type": "/tendermint.fundraising.fundraising.MsgCreateFixedPriceAuction",
        "auctioneer": "cosmos1m4ys0e222x45657hrg9y2gadfxtcqja270rdkg",
        "start_price": "1.000000000000000000",
        "selling_coin": {
          "denom": "denom1",
          "amount": "1000000000000"
        },
        "paying_coin_denom": "denom2",
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
        "start_time": "2021-12-01T00:00:00Z",
        "end_time": "2021-12-30T00:00:00Z"
      }
    ],
    "memo": "",
    "timeout_height": "0",
    "extension_options": [],
    "non_critical_extension_options": []
  },
  "auth_info": {
    "signer_infos": [
      {
        "public_key": {
          "@type": "/cosmos.crypto.secp256k1.PubKey",
          "key": "A8IlstomF7Z1qDMYBL1rhpWwM47IgJSHkq+e4zzeg2Xw"
        },
        "mode_info": {
          "single": {
            "mode": "SIGN_MODE_DIRECT"
          }
        },
        "sequence": "6"
      }
    ],
    "fee": {
      "amount": [],
      "gas_limit": "200000",
      "payer": "",
      "granter": ""
    }
  },
  "signatures": [
    "2Vjdw6VhsQ7Laxli8Wm9ESmBqChJMqeerX2HEUmUVC8d/467gDYC4TQSHsRJRMFXm65quWxekwkQgTUoY7+HPA=="
  ]
}
```

### CreateEnglishAuction

Example command:

```bash
TODO: IT IS BEING DEVELOPED
```

Result:

```result

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

Result:

```result
{
  "@type": "/cosmos.tx.v1beta1.Tx",
  "body": {
    "messages": [
      {
        "@type": "/tendermint.fundraising.fundraising.MsgCancelAuction",
        "auctioneer": "cosmos1xg6ngnzf9kz9606kx45z2g3eeskre7cm4effpq",
        "auction_id": "1"
      }
    ],
    "memo": "",
    "timeout_height": "0",
    "extension_options": [],
    "non_critical_extension_options": []
  },
  "auth_info": {
    "signer_infos": [
      {
        "public_key": {
          "@type": "/cosmos.crypto.secp256k1.PubKey",
          "key": "Aq7NW7m/FazN7NVy0bQqm3U7RD/ySZ34DDrw0RJ9rGsI"
        },
        "mode_info": {
          "single": {
            "mode": "SIGN_MODE_DIRECT"
          }
        },
        "sequence": "1"
      }
    ],
    "fee": {
      "amount": [],
      "gas_limit": "200000",
      "payer": "",
      "granter": ""
    }
  },
  "signatures": [
    "mvfN/nIzivLX4pRGpC2nTsHUNfucbf5oA605MLpg5ksO5kegjQ7brB5QlGM9qpRczXYxvguY1pjOivaWUCtUdw=="
  ]
}
```

### PlaceBid

Example command:

```bash
fundraisingd tx fundraising bid 1 1.0 5000000denom2 \
--chain-id fundraising \
--from alice \
--keyring-backend test \
--broadcast-mode block \
--yes \
--output json | jq
```

Result:

```result
{
  "@type": "/cosmos.tx.v1beta1.Tx",
  "body": {
    "messages": [
      {
        "@type": "/tendermint.fundraising.fundraising.MsgPlaceBid",
        "auction_id": "1",
        "bidder": "cosmos15ghqvkhllee5uvy400pw2fuh4d45ayykuzm2ts",
        "price": "1.000000000000000000",
        "coin": {
          "denom": "denom2",
          "amount": "5000000"
        }
      }
    ],
    "memo": "",
    "timeout_height": "0",
    "extension_options": [],
    "non_critical_extension_options": []
  },
  "auth_info": {
    "signer_infos": [
      {
        "public_key": {
          "@type": "/cosmos.crypto.secp256k1.PubKey",
          "key": "Asinw4aW3b+D81bkY2ALm5K2RlFxQTRdfPOrs36CnkLF"
        },
        "mode_info": {
          "single": {
            "mode": "SIGN_MODE_DIRECT"
          }
        },
        "sequence": "6"
      }
    ],
    "fee": {
      "amount": [],
      "gas_limit": "200000",
      "payer": "",
      "granter": ""
    }
  },
  "signatures": [
    "YqPjBdxyeP5ufw+XGw37dZzEjI/UPqiBz7JMMzxUJ9EA9oRfVsTKFl4BZ9vsaRztJw4a2O6SyE00ACamjYqECg=="
  ]
}
```

## Query

+++ https://github.com/tendermint/fundraising/blob/main/proto/fundraising/query.proto#L14-L42


### Params 

```bash
# Query the values set as fundraising parameters
fundraisingd q fundraising params --output json | jq
```

```result
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

```result
{
  "auctions": [
    {
      "@type": "/tendermint.fundraising.FixedPriceAuction",
      "base_auction": {
        "id": "1",
        "type": "AUCTION_TYPE_FIXED_PRICE",
        "auctioneer": "cosmos1m4ys0e222x45657hrg9y2gadfxtcqja270rdkg",
        "selling_pool_address": "cosmos18xzvtd72y9j8xyf8a36z5jjhth7qgtcwhh8lz7yee3tvxqn6ll5quh78zq",
        "paying_pool_address": "cosmos18permjyqvk5flft8ey9egr7hd4ry8tauqt4f9mg9knn4vvtkry9sujucrl",
        "start_price": "1.000000000000000000",
        "selling_coin": {
          "denom": "denom1",
          "amount": "1000000000000"
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
          "amount": "999995000000"
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

Example command:

```bash
# Query for the specific auction with the auction id
fundraisingd q fundraising auction 1 -o json | jq
```

Result:

```result
{
  "auction": {
    "@type": "/tendermint.fundraising.FixedPriceAuction",
    "base_auction": {
      "id": "1",
      "type": "AUCTION_TYPE_FIXED_PRICE",
      "auctioneer": "cosmos1m4ys0e222x45657hrg9y2gadfxtcqja270rdkg",
      "selling_pool_address": "cosmos18xzvtd72y9j8xyf8a36z5jjhth7qgtcwhh8lz7yee3tvxqn6ll5quh78zq",
      "paying_pool_address": "cosmos18permjyqvk5flft8ey9egr7hd4ry8tauqt4f9mg9knn4vvtkry9sujucrl",
      "start_price": "1.000000000000000000",
      "selling_coin": {
        "denom": "denom1",
        "amount": "1000000000000"
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
        "amount": "999995000000"
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

Example command:

```bash
# Query for all bids of the auction with the given auction id
fundraisingd q fundraising bids 1 \
-o json | jq
```

Result:

```result
{
  "bids": [
    {
      "auction_id": "1",
      "sequence": "1",
      "bidder": "cosmos15ghqvkhllee5uvy400pw2fuh4d45ayykuzm2ts",
      "price": "1.000000000000000000",
      "coin": {
        "denom": "denom2",
        "amount": "5000000"
      },
      "height": "1457",
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

This command is used by an auctioneer to query vesting information. It only returns results when the auction is in vesting status

Example command:

```bash
# Query for all vesting queues 
fundraisingd q fundraising vestings 1 \
-o json | jq
```

Result:

```result
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
