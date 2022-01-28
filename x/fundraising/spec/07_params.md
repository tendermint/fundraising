<!-- order: 8 -->

# Parameters

The `fundraising` module contains the following parameters:

| Key                        | Type      | Example                                                           |
| -------------------------- | --------- | ----------------------------------------------------------------- |
| AuctionCreationFee         | sdk.Coins | [{"denom":"stake","amount":"100000000"}]                          |
| ExtendedPeriod             | uint32    | 3600 * 24                                                         |
| FeeCollectorAddress        | string    | cosmos1kxyag8zx2j9m8063m92qazaxqg63xv5h7z5jxz8yr27tuk67ne8q0lzjm9 |

## AuctionCreationFee

`AuctionCreationFee` is the fee required to pay to create an auction. This fee prevents from spamming attack.

## ExtendedPeriod

`ExtendedPeriod` is the extended period that determines how long the extended auction round is.

## FeeCollectorAddress

`FeeCollectorAddress` is the fee collector account address that collect fees in the fundraising module, such as auction creation fees.