---
Title: Fundraisingd
Description: A high-level overview of how the command-line interface (CLI) works for the fundraising module.
---
# CLI Reference

## Synopsis

This document provides a high-level overview of how the command line (CLI) interface works for the `fundraising` module. To set up a local testing environment, it requires the latest Ignite CLI. If you don't have Ignite CLI set up in your local machine, see [this guide](https://docs.ignite.com/welcome/install) to install it. Run this command under the project root directory `$ ignite chain serve -c config-test.yml` or simply `$ make localnet`.

Note that [jq](https://stedolan.github.io/jq/) is recommended to be installed as it is used to process JSON throughout the document.

## Command Line Interface

- [Transaction](#Transaction)
  - [CreateFixedPriceAuction](#CreateFixedPriceAuction)
  - [CreateBatchAuction](#CreateBatchAuction)
  - [CancelAuction](#CancelAuction)
  - [AddAllowedBidder](#AddAllowedBidder)
  - [PlaceBid](#PlaceBid)
  - [ModifyBid](#ModifyBid)
- [Query](#Query)
  - [Params](#Params)
  - [Auctions](#Auctions)
  - [Auction](#Auction)
  - [AllowedBidder](#AllowedBidder)
  - [AllowedBidders](#AllowedBidders)
  - [Bids](#Bids)
  - [Vestings](#Vestings)

# Transaction

+++ https://github.com/tendermint/fundraising/blob/main/proto/fundraising/tx.proto#L12-L35

## CreateFixedPriceAuction

An auctioneer can create a fixed price auction by setting the following parameters. In a fixed price auction, `start_price` is the matched price and bidders can buy the selling coins on a first-come, first-served basis. See the [spec](https://github.com/tendermint/fundraising/blob/main/x/fundraising/spec/01_concepts.md#auction-types) for a detailed and technical information about a fixed priced auction type.

Usage

```bash
create-fixed-price-auction [file]
```

| **Argument** |  **Description**                                                   |
| :----------- | :----------------------------------------------------------------- |
| file         | file that contains required fields to create a fixed price auction |

Field description of the input file

| **Field**         |  **Description**                                                                    |
| :---------------- | :---------------------------------------------------------------------------------- |
| start_price       | The starting price of the selling coin; it is proportional to the paying coin denom | 
| selling_coin      | The selling amount of coin for the auction                                          | 
| paying_coin_denom | The paying coin denom that bidders use to bid with                                  | 
| vesting_schedules | The vesting schedules that release the paying coins to the autioneer                | 
| start_time        | The start time of the auction                                                       | 
| end_time          | The end time of the auction                                                         | 

Example of input as JSON:

An auctioneer creates a fixed price auction for `1000000000000denom1` selling coin where the start price is 2.0. It means that the price of 1 of denom1 is 2 of denom2. The auctioneer sets their vesting schedules for themselves to receive the accumulated paying coin amount when the auction ends. This is a gesture for taking responsibility for their auction participants.

```json
{
  "start_price": "2.000000000000000000",
  "selling_coin": {
    "denom": "denom1",
    "amount": "1000000000000"
  },
  "paying_coin_denom": "denom2",
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
  "end_time": "2022-06-01T00:00:00Z"
}
```

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

#
# Tips
#
# Query all auctions
fundraisingd q fundraising auctions -o json | jq

# Query to see if the selling coin is safely reserved
fundraisingd q bank balances <selling_reserve_address> -o json | jq
```

## CreateBatchAuction

An auctioneer can create a batch auction by setting the following parameters. Differently from a fixed price auction, `start_price` does not affect the determination of the matched price, but is provided by the auctioneer as a reference price to bidders. See the [spec](https://github.com/tendermint/fundraising/blob/main/x/fundraising/spec/01_concepts.md#auction-types) for a detailed and technical information about a batch auction type.

Usage

```bash
create-batch-auction [file]
```

| **Argument** |  **Description**                                       |
| :----------- | :----------------------------------------------------- |
| file         | file that contains required fields for a batch auction |

Field description of the input file

| **Field**           |  **Description**                                                                    |
| :------------------ | :---------------------------------------------------------------------------------- |
| start_price         | The starting price of the selling coin; it is proportional to the paying coin denom. This might not be the matched price and it can be used as a reference price from the auctioneer. | 
| min_bid_price       | The minimum bid price that bidders must place with                                  |
| selling_coin        | The selling amount of coin for the auction                                          | 
| paying_coin_denom   | The paying coin denom that bidders use to bid with                                  | 
| vesting_schedules   | The vesting schedules that release the paying coins to the autioneer                | 
| max_extended_round  | The maximum number of extended rounds that provides additional opportunity for the bidders to place bids when more than a certain ratio of the number of the matched bids are reduced compared to the previous end time  |
| extended_round_rate | The threshold reduction of the number of the matched bids are reduced compared to the previous end time to decide the necessity of another extended round | 
| start_time          | The start time of the auction                                                       | 
| end_time            | The end time of the auction                                                         | 

Example of input as JSON:

An auctioneer creates a batch price auction for `1000000000000denom1` selling coin where the start price is `0.3` and the minimum bid price is `0.1`. Unlike a fixed price auction, the start price might not be the matched price. It is only used as a reference of what the auctioneer considers their selling coin price is. The auctioneer sets their vesting schedules for themselves to receive the accumulated paying coin amount. This is a gesture for taking responsibility for auction participants.

```json
{
  "start_price": "0.300000000000000000",
  "min_bid_price": "0.100000000000000000",
  "selling_coin": {
    "denom": "denom1",
    "amount": "1000000000000"
  },
  "paying_coin_denom": "denom2",
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
  "max_extended_round": 2,
  "extended_round_rate": "0.150000000000000000",
  "start_time": "2022-05-01T00:00:00Z",
  "end_time": "2022-06-01T00:00:00Z"
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

#
# Tips
#
# Query all auctions
fundraisingd q fundraising auctions -o json | jq

# Query to see if the selling coin is safely reserved
fundraisingd q bank balances <selling_reserve_address> -o json | jq
```

## CancelAuction

This command is useful for an auctioneer when the auctioneer made mistake(s) on some values of the auction. The module doesn't support update functionality. Instead, the module allows them to cancel an auction and recreate it with correct values. Note that it can only be cancelled when the auction has not started yet.

Usage

```bash
cancel [auction-id]
```

| **Argument** |  **Description** |
| :----------- | :--------------- |
| auction-id   | auction id       |

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

## AddAllowedBidder

**Important Note**: the module is fundamentally designed to delegate authorization to an external module to add allowed bidder list for an auction. When an auction is created, it is closed state; meaning that no bidders are allowed to place a bid unless they are authorized. 

`AddAllowedBidder` CLI command is a special command that is built for **testing purpose**. It adds an allowed bidder for the auction and this command is only available when you build binary `fundraisingd` with `config-test.yml` file which passes `enableAddAllowedBidder` ldflags true under the hood.

Usage

```bash
add-allowed-bidder [auction-id] [bidder] [max-bid-amount]
```

| **Argument**   |  **Description**                                     |
| :------------- | :--------------------------------------------------- |
| auction-id     | auction id                                           |
| bidder         | bidder address                                       |
| max-bid-amount | maximum bid amount that the bidder is allowed to bid |

Example command:

```bash
#
# Once again, this CLI command is not available in mainnet environment
#

# Bob adds himself in the auction's allowed bidders list
fundraisingd tx fundraising add-allowed-bidder 1 cosmos1mzgucqnfr2l8cj5apvdpllhzt4zeuh2cshz5xu 2000000000 \
--chain-id fundraising \
--from bob \
--keyring-backend test \
--broadcast-mode block \
--yes \
--output json | jq

# Steve adds himself in the auction's allowed bidders list
fundraisingd tx fundraising add-allowed-bidder 2 cosmos185fflsvwrz0cx46w6qada7mdy92m6kx4gqx0ny 5000000000 \
--chain-id fundraising \
--from steve \
--keyring-backend test \
--broadcast-mode block \
--yes \
--output json | jq

#
# Tips
#
# Query the allowed bidders list for the auction
fundraisingd q fundraising allowed-bidders 1 -o json | jq
fundraisingd q fundraising allowed-bidder 1 cosmos1mzgucqnfr2l8cj5apvdpllhzt4zeuh2cshz5xu -o json | jq
fundraisingd q fundraising allowed-bidder 1 cosmos185fflsvwrz0cx46w6qada7mdy92m6kx4gqx0ny -o json | jq
```

## PlaceBid

This command is used for a bidder to place a bid to the auction where the bidder must exist in the list of the allowed bidders for the auction. 

Usage

```bash
bid [auction-id] [bid-type] [price] [coin]
```

| **Argument** |  **Description**                     |
| :----------- | :----------------------------------- |
| auction-id   | auction id | 
| bid-type     | 1) fixed-price (fp or f), 2) batch-worth (bw or w), and 3) batch-many (bm or m) where 1 is only for `FixedPriceAuction` and 2 and 3 are for `BatchAuction` |
| price        | bid price of a selling coin as the unit of a paying coin. For `FixedPriceAuction`, it must be the start price of the auction. For `BatchAuction`, the price must be higher than or equal to the minimum bid price of the auction | 
| coin         | how many coins to bid. The denom can be both selling and paying coin denom for fixed-price. However, batch-worth must be paying coin denom and batch-many must be selling coin denom |

Example command:

```bash
# Place a fixed price bid type for the fixed price auction
fundraisingd tx fundraising bid 1 fixed-price 2.0 5000000denom2 \
--chain-id fundraising \
--from bob \
--keyring-backend test \
--broadcast-mode block \
--yes \
--output json | jq

# Place a batch-worth bid type for the batch auction
fundraisingd tx fundraising bid 2 batch-worth 0.35 10000000denom2 \
--chain-id fundraising \
--from steve \
--keyring-backend test \
--broadcast-mode block \
--yes \
--output json | jq

# Place a batch-many bid type for the batch auction
fundraisingd tx fundraising bid 2 batch-many 0.4 10000000denom1 \
--chain-id fundraising \
--from steve \
--keyring-backend test \
--broadcast-mode block \
--yes \
--output json | jq

#
# Tips
#
# Query all bids that belong to the auction
fundraisingd q fundraising bids 1 -o json | jq
fundraisingd q fundraising bids 2 -o json | jq
```

## ModifyBid

This command is used for modifying the bid. It is only supported for `BatchAuction`. The bidder is allowed to modify the bid only with the same bid type and they must provide either higher bid price or larger bid amount. Lowering bid price or lesser bid amount is restricted.

Usage

```bash
modify-bid [auction-id] [bid-id] [price] [coin]
```

| **Argument**|  **Description**              |
| :---------- | :---------------------------- |
| auction-id  | auction id                    | 
| bid-id      | bid id that the bidder placed |
| price       | bid price of a selling coin as the unit of a paying coin. For `BatchAuction`, the price cannot be lower than the original bid price of the auction | 
| coin        | how many coins to bid. The denom must be the same as the modifying bid. The amount cannot be smaller than that of the original coin amount. |


Example command:

```bash
# Modify the bid price
fundraisingd tx fundraising modify-bid 2 1 0.38 10000000denom2 \
--chain-id fundraising \
--from steve \
--keyring-backend test \
--broadcast-mode block \
--yes \
--output json | jq

# Modify the bid amount
fundraisingd tx fundraising modify-bid 2 2 0.4 15000000denom1 \
--chain-id fundraising \
--from steve \
--keyring-backend test \
--broadcast-mode block \
--yes \
--output json | jq

#
# Tips
#
# Query all bids that belong to the auction
fundraisingd q fundraising bids 1 -o json | jq
```

# Query

+++ https://github.com/tendermint/fundraising/blob/main/proto/fundraising/query.proto#L15-L63

## Params 

This command is used to query the current fundraising parameters information.

Usage

```bash
params
```

Example commands:

```bash
# Query the values set as fundraising parameters
fundraisingd q fundraising params \
--output json | jq
```

## Auctions

This command is used to query the information of all auctions.

Usage

```bash
auctions
```

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

# Query for all auctions with the given auction type
# Ref: https://github.com/tendermint/fundraising/blob/main/x/fundraising/spec/02_state.md#auction-type
fundraisingd q fundraising auctions \
--type AUCTION_TYPE_FIXED_PRICE \
-o json | jq
```

## Auction

This command is used by an auctioneer to query the information of a specific auction.

Usage

```bash
auction [auction-id]
```

Example command:

```bash
# Query for the specific auction with the auction id
fundraisingd q fundraising auction 1 \
-o json | jq
```

## AllowedBidder

This command is used to query the specific allowed bidder information.

Usage

```bash
allowed-bidder [auction-id] [bidder]
```

Example command:

```bash
# Query for a specific allowed bidders for the auction
fundraisingd q fundraising allowed-bidder 1 cosmos1mzgucqnfr2l8cj5apvdpllhzt4zeuh2cshz5xu \
-o json | jq
```

## AllowedBidders

This command is used to query all allowed bidders list for the auction.

Usage

```bash
allowed-bidders [auction-id]
```

Example command:

```bash
# Query for a specific allowed bidders for the auction
fundraisingd q fundraising allowed-bidders 1 \
-o json | jq
```

## Bids

This command is used by an auctioneer to query the information of all the bids of a specific auction.

```bash
bids [auction-id]
```

Example command:

```bash
# Query for all bids of the auction with the given auction id
fundraisingd q fundraising bids 1 \
-o json | jq
```

## Vestings

This command is used by an auctioneer to query vesting information. It only returns results when the auction is in vesting status.

```bash
vestings [auction-id]
```

Example command:

```bash
# Query for all vesting queues 
fundraisingd q fundraising vestings 1 \
-o json | jq
```
