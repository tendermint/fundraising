<!-- order: 8 -->

# Parameters

The `fundraising` module contains the following parameters:

| Key                        | Type      | Example                                                           |
| -------------------------- | --------- | ----------------------------------------------------------------- |
| AuctionCreationFee         | sdk.Coins | [{"denom":"stake","amount":"100000000"}]                          |
| ExtendedPeriod             | uint32    | 3600 * 24                                                         |
| AuctionFeeCollector        | string    | cosmos1t2gp44cx86rt8gxv64lpt0dggveg98y4ma2wlnfqts7d4m4z70vqrzud4t |

## AuctionCreationFee

`AuctionCreationFee` is the fee required to pay to create an auction. This fee prevents from spamming attack.

## ExtendedPeriod

`ExtendedPeriod` is the extended period that determines how long the extended auction round is.