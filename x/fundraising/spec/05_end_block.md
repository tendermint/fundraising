<!-- order: 5 -->

At the end of each end block, the `fundraising` module operates the following executions based on auction type.

## FixedPriceAuctionType

The module first gets all auctions registered in the store and proceed operations depending on auction status.

If the auction's status is `AuctionStatusVesting`, it gets a list of vesting queues in the store and look up the release time of each vesting queue to see if the module needs to distribute the paying coin to the auctioneer.

If the auction's status is `AuctionStatusStandBy`, it compares the current time and the start time of the auction and updte the status if necessary.

If the auction's status `AuctionStatusStarted`, it distribute the allocated paying coin for bidders for the auction and set vesting schedules if they are defined when the auction is created.

## EnglishAuctionType

Work in progress.


