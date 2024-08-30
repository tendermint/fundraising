package types

import "cosmossdk.io/collections"

const (
	// ModuleName defines the module name
	ModuleName = "fundraising"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_fundraising"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

var (
	// ParamsKey is the prefix to retrieve all Params
	ParamsKey = collections.NewPrefix("p_fundraising")

	// BidKey is the prefix to retrieve all Bid
	BidKey = collections.NewPrefix("bid/value/")
	// BidCountKey is the prefix to retrieve all Bid cound
	BidCountKey = collections.NewPrefix("bid/count/")

	// AuctionKey is the prefix to retrieve all Auction
	AuctionKey = collections.NewPrefix("auction/value/")
	// AuctionCountKey is the prefix to retrieve all Auction count
	AuctionCountKey = collections.NewPrefix("auction/count/")

	// AllowedBidderKey is the prefix to retrieve all AllowedBidder
	AllowedBidderKey = collections.NewPrefix("AllowedBidder/value/")

	// VestingQueueKey is the prefix to retrieve all VestingQueue
	VestingQueueKey = collections.NewPrefix("VestingQueue/value/")

	// MatchedBidsLenKey is the prefix to retrieve all MatchedBidsLen
	MatchedBidsLenKey = collections.NewPrefix("MatchedBidsLen/value/")
)
