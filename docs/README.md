# Documentation

How to use the fundraising module documentation.


- [Documentation](#documentation)
  - [Overview](#overview)
  - [More Documentations](#more-documentations)

## Overview

The main purpose of the `fundraising` module is to sell a certain amount of selling coins with a proper price. 
The characteristics of how to determine the matched price, the matched bids, and the matched selling coins differentiate the auction types. 
The outline of the flow of progressing an auction is following.

1. Create an auction
    - An auctioneer creates an auction.
2. Add bidder(s) to the list of the allowed bidders
    - The auctioneer adds bidder(s) as the allowed bidders that enables to participates in the auction.
3. Place/modify a bid
    - The bidders added in the list of the allowed bidders place bids and modify the bids according to the auction types.
4. Calculation of the matched price, the matched bids, and the matched selling coins
    - The matched price, the matched bids, and the matched selling coins are calculated based on the placed bids and the auction type.
5. Allocation and refund coins
    - The matched selling coins are distributed to the matched bidders.
    - The remaining selling coins are refunded to the auctioneer.
    - The matched paying coins are reserved in the vesting address.
    - The remaining paying coins are refunded to the bidders.
6. Vesting the matched paying coins
    - According to the vesting schedule, the pre-configured amount of the matched paying coins are vested to the auctioneer.



## More Documentations
The following documentations further provide the explanations on `fundraising` module.

* [How-Tos](./How-To/README.md)
   - How to use API and CLI
* [Tutorials](./Tutorials/README.md)
  - How to proceed with the auction and how to calculate the matched price 


<!-- we use the  *Grand Unified Theory of Documentation* (David Laing) as described by [Divio](https://documentation.divio.com/) as a basis for our documentation strategy.

This approach outlines four specific use cases for documentation:

* [Explanation](./Explanation/README.md)
* [How-Tos](./How-To/README.md)
* [Tutorials](./Tutorials/README.md)

For further background please see [the ADR relating to the documentation structure](./Explanation/ADR/adr-002-docs-structure.md). 

## Contributing

* Write all documentation following [Google Documentation Best Practice](https://google.github.io/styleguide/docguide/best_practices.html)
* Generate as much documentation as possible from the code.
* Raise a PR for all documentation changes
* Follow our [Code of Conduct](../CONTRIBUTING.md)

## Reference

- [Google Style Guide for Markdown](https://github.com/google/styleguide/blob/gh-pages/docguide/style.md)
- [Write the Docs global community](https://www.writethedocs.org/)
- [Write the Docs Code of Conduct](https://www.writethedocs.org/code-of-conduct/#the-principles)
- [The good docs project](https://github.com/thegooddocsproject)
- [Readme editor](https://readme.so/editor) -->