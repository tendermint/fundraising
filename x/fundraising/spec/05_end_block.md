<!-- order: 5 -->

At the end of each block, the winners are logically calculated for all the running auctions with `StatusStarted`.

Winner calculating logic will be performed for auctions with ongoing status

- For auction with `EnglishAuction` type.
    - When the case of `Extended` is zero
        - Get all valid bid messages on the auction
        - Calculate winner price and determine which bid messages are win
        - Mark winning bid message's `isWinner` as true based on the above calculation
        - Update `Endtime`, i.e. add last item of `Endtime` with `ExtendedAuctionPeriod`, then append it to `Endtime`
    - When the case of `Extended` is not zero
        - Get all valid bid messages on the auction
        - Calculate winner price and determine which bid messages are win
        - Make winning bid list and compare it with previous winner list
            - If difference between winning bid list and previous winner list exceeds `ExtendedRate`, update current winning bid message's `isWinner` to true and the rest to false. After that, update `Endtime`
            - If not, mark the `isWinner` of all winning bid message as true and the rest to false. Then, update `AuctionStatus` to "Selling coin distributing ready"

- For auction with `FixedPriceAuction`
    - Get all valid bid messages on the auction
    - For case of sum of `Coin`/`Price` of all bid messages is less than or equal to amount of `SellingCoin`
        - Mark all bid message's `isWinner` as true
    - For case of sum of `Coin`/`Price` of all bid messages is greater than amount of `SellingCoin`
        - Mark bid message's `isWinner` as true until cumulated `Coin`/`Price` is greater than or equal to amount of `SellingCoin` in the order of ascending `Sequence`
        - The last winner will receive the amount equal to the `SellingCoin` minus the cumulative amount up to the person in front of him/her.
        - Update `AuctionStatus` to "Selling coin distributing ready"

Distribution of Selling coin
- Distribute selling coin for auction with "Selling coin distributing ready" status after winner calculation
- Auction winner gets `determinedAmount` at the `winningPrice`
- Update `AuctionStatus` to "Selling coin distributed"

Distribution of Paying coin
- Distribute paying coin for auction with "Selling coin distributed" status
- Paying coin will be distributed according to `VestingSchedule`
- If `Time` in `VestingSchedule` has passed based on the block time, paying coin will move `PayingPoolAddress` to Auctioneer's account.
- The amounts of Paying coin to move are `TotalPayingcoin` * `Weight` of this vesting time / sum of `Weight` of vesting time later than or equal to this vesting time
- If all paying coin are distributed, update `AuctionStatus` to "Closed"





