---
Title: Fundraisingd
Description: A high-level overview of how the command-line interface (CLI) works for the fundraising module.
---

# CLI Reference

This document provides a high-level overview of how the command line (CLI) interface works for the fundraising module. The executable file name is called `fundraisingd`.

## Command Line Interface

To test out the following commands, you must set up a local network. By simply running `$ make localnet` under the root project directory, you can start the local network. It requires the latest [Starport](https://starport.com/). If you don't have `Starport` set up in your local machine, see this [Starport guide](https://docs.starport.network/) to install it.  

- [Transaction](#Transaction)
    * [CreateFixedPriceAuction](#CreateFixedPriceAuction)
    * [CreateBatchAuction](#CreateBatchAuction)
    * [CancelAuction](#CancelAuction)
    * [AddAllowedBidder](#AddAllowedBidder)
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

This command provides an opportunity for an auctioneer to create an auction to raise funds for their project. It is the most basic and simple type of an auction that has first come first served basis. When an auctioneer creates a fixed price auction, they must determine the fixed starting price. It is proportional to the paying coin denom. To give you an example, an auctioneer sells X coin and plans to receive Y coin for the auction. The price of X coin is determined by the proportion of Y coin. Let's assume that the price of Y coin is currently $30 and the auctioneer wants to sell their X coin for $15, then they must set 0.5 as the fixed starting price. Once the auction is successfully created, bidders can now start to bid. The bidders must provide the same start price when they bid. See the [spec](https://github.com/tendermint/fundraising/blob/main/x/fundraising/spec/01_concepts.md#auction-types) for a detailed and technical information about the fixed price auction type.

JSON example:

In this JSON example, an auctioneer plans to create a fixed price auction that plans to sell `1000000000000denom1` coin, and the starting price is `1.0` which means that the price of `denom1` is the same as `denom2`. The auction starts at `2022-01-21T00:00:00Z` and ends at `2022-02-21T00:00:00Z`. As soon as the auction starts, bidders can now bid for the auction with any amount of coin they are willing bid with the fixed start price. When it ends, the paying amount of coin that is reserved for all bids is expected to be released based on the vesting schedules and if the selling coin is not entirely sold out, it transfers it back to the auctioneer.

```json
{
  "start_price": "1.000000000000000000",
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
  "end_time": "2022-03-01T00:00:00Z"
}
```

Reference the description of each field:

| **Field**         |  **Description**                                                                    |
| :---------------- | :---------------------------------------------------------------------------------- |
| start_price       | The starting price of the selling coin; it is proportional to the paying coin denom | 
| min_bid_price     | The minimum bid price that bidders must provide                                     |
| selling_coin      | The selling amount of coin for the auction                                          | 
| paying_coin_denom | The paying coin denom that bidders use to bid with                                  | 
| vesting_schedules | The vesting schedules that release the paying coins to the autioneer                | 
| start_time        | The start time of the auction                                                       | 
| end_time          | The end time of the auction                                                         | 
|                   |                                                                                     |

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

```json
{
  "@type": "/cosmos.tx.v1beta1.Tx",
  "body": {
    "messages": [
      {
        "@type": "/tendermint.fundraising.MsgCreateFixedPriceAuction",
        "auctioneer": "cosmos1dncsflcfknkmlmt3t6836tkd3mu742e2wh4r70",
        "start_price": "1.000000000000000000",
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
        "end_time": "2022-03-01T00:00:00Z"
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
          "key": "A3mbh7d1pTgT3xSDyXHjdpcaxm58t0azRCXeGP0EsKsQ"
        },
        "mode_info": {
          "single": {
            "mode": "SIGN_MODE_DIRECT"
          }
        },
        "sequence": "0"
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
    "UahsRZ27hATh0xu7M/IWFvNvaFESpQ+W0RmQhQql3ERsnYdTDrFP81/MxyYxuX4WNBUv4+3FyhOwEQ7hqlU+MQ=="
  ]
}
```

### CreateBatchAuction

This command is another type of an auction for an auctioneer to raise funds for their project. See the [spec](https://github.com/tendermint/fundraising/blob/main/x/fundraising/spec/01_concepts.md#auction-types) for a detailed and technical information about a batch auction type.

JSON example:

```json
{
  "start_price": "0.500000000000000000",
  "min_bid_price": "0.100000000000000000",
  "selling_coin": {
    "denom": "denom1",
    "amount": "1000000000000"
  },
  "paying_coin_denom": "denom2",
  "vesting_schedules": [
    {
      "release_time": "2023-06-01T00:00:00Z",
      "weight": "0.500000000000000000"
    },
    {
      "release_time": "2023-12-01T00:00:00Z",
      "weight": "0.500000000000000000"
    }
  ],
  "max_extended_round": 2,
  "extended_round_rate": "0.150000000000000000",
  "start_time": "2022-02-01T00:00:00Z",
  "end_time": "2022-06-20T00:00:00Z"
}
```

Reference the description of each field:

| **Field**           |  **Description**                                                                    |
| :------------------ | :---------------------------------------------------------------------------------- |
| start_price         | The starting price of the selling coin; it is proportional to the paying coin denom | 
| min_bid_price       | The minimum bid price that bidders must provide                                     |
| selling_coin        | The selling amount of coin for the auction                                          | 
| paying_coin_denom   | The paying coin denom that bidders use to bid with                                  | 
| vesting_schedules   | The vesting schedules that release the paying coins to the autioneer                | 
| max_extended_round  | The number of extended rounds                                                       | 
| extended_round_rate | The rate that determines if the auction needs to run another round                  | 
| start_time          | The start time of the auction                                                       | 
| end_time            | The end time of the auction                                                         | 
|                     |                                                                                     |


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

Result:

```json
{
  "@type": "/cosmos.tx.v1beta1.Tx",
  "body": {
    "messages": [
      {
        "@type": "/tendermint.fundraising.MsgCreateBatchAuction",
        "auctioneer": "cosmos1ygsq4lnaernkz02un4fyksdzhzm6aazqpktj9p",
        "start_price": "0.100000000000000000",
        "min_bid_price": "0.100000000000000000",
        "selling_coin": {
          "denom": "denom1",
          "amount": "1000000000000"
        },
        "paying_coin_denom": "denom2",
        "vesting_schedules": [
          {
            "release_time": "2023-06-01T00:00:00Z",
            "weight": "0.500000000000000000"
          },
          {
            "release_time": "2023-12-01T00:00:00Z",
            "weight": "0.500000000000000000"
          }
        ],
        "max_extended_round": 2,
        "extended_round_rate": "0.150000000000000000",
        "start_time": "2022-02-01T00:00:00Z",
        "end_time": "2022-06-20T00:00:00Z"
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
          "key": "A1NJdw96iIRrlhyrPmkWVKNcbrd8mhCRXb4InQqjU1Vm"
        },
        "mode_info": {
          "single": {
            "mode": "SIGN_MODE_DIRECT"
          }
        },
        "sequence": "0"
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
    "UJmyeX5azpoTCmJIUhzr7UqmipUadPHlLuYSfuYZonRKiunRj6JkJQ4xWzzvE05ehsoWXODBALtp4Brmnr87WA=="
  ]
}
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

```json
{
  "@type": "/cosmos.tx.v1beta1.Tx",
  "body": {
    "messages": [
      {
        "@type": "/tendermint.fundraising.MsgCancelAuction",
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

Result:

```json
{
  "@type": "/cosmos.tx.v1beta1.Tx",
  "body": {
    "messages": [
      {
        "@type": "/tendermint.fundraising.MsgAddAllowedBidder",
        "auction_id": "1",
        "allowed_bidder": {
          "bidder": "cosmos1tfzynkllgxdpmrcknx2j5d0hj9zd82tceyfa5n",
          "max_bid_amount": "1000000000"
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
          "key": "A8VLxM/RDIlFEtOe7rfzA2Am55/Zam2n+oq1+I/Ovkbv"
        },
        "mode_info": {
          "single": {
            "mode": "SIGN_MODE_DIRECT"
          }
        },
        "sequence": "0"
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
    "D49R49OD1YIBzdVVy5g1yIc8AvKbII6f8n3NpJDDHbY3O2vX/dwsoC2TX5eWRSXGgJ92+PfIZIek5PrsZWyfxQ=="
  ]
}
```

### PlaceBid

Example command:

```bash
fundraisingd tx fundraising bid 1 1.0 5000000denom2 \
--chain-id fundraising \
--from steve \
--keyring-backend test \
--broadcast-mode block \
--yes \
--output json | jq
```

Result:

```json
{
  "@type": "/cosmos.tx.v1beta1.Tx",
  "body": {
    "messages": [
      {
        "@type": "/tendermint.fundraising.MsgPlaceBid",
        "auction_id": "1",
        "bidder": "cosmos1tfzynkllgxdpmrcknx2j5d0hj9zd82tceyfa5n",
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
          "key": "A8VLxM/RDIlFEtOe7rfzA2Am55/Zam2n+oq1+I/Ovkbv"
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
    "Ahrvo4CXneHxTd0Hgyt+HdZXmrhKhm1ijo5Tf7/K7OcK4P5590UlDpoqJ7ofLB738AGt+3rJ+cHy+K09KqBFaA=="
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

```json
{
  "auction_creation_fee": [
    {
      "denom": "stake",
      "amount": "100000000"
    }
  ],
  "extended_period": 1,
  "fee_collector_address": "cosmos1kxyag8zx2j9m8063m92qazaxqg63xv5h7z5jxz8yr27tuk67ne8q0lzjm9"
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
