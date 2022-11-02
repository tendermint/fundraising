package cli

// DONTCOVER

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagAuctionStatus = "status"
	FlagAuctionType   = "type"
	FlagBidderAddr    = "bidder-addr"
	FlagIsMatched     = "is-matched"
)

// flagSetAuctions returns a set of defined flags to query the auctions.
func flagSetAuctions() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.String(FlagAuctionStatus, "", "The auction status; AUCTION_STATUS_STANDBY, AUCTION_STATUS_STARTED, and etc.")
	fs.String(FlagAuctionType, "", "The auction type; AUCTION_TYPE_FIXED_PRICE or AUCTION_TYPE_ENGLISH")

	return fs
}

// flagSetBids returns a set of defined flags to query the bids.
func flagSetBids() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.String(FlagBidderAddr, "", "The bech32 address of the bidder account")
	fs.String(FlagIsMatched, "", "The bid that is successfully matched (a.k.a., winner)")

	return fs
}
