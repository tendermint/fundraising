package types

import (
	"bytes"
	time "time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	// ModuleName defines the module name
	ModuleName = "fundraising"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for the fundraising module
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_fundraising"
)

var (
	LastAuctionIdKey   = []byte{0x11} // key to retrieve the latest auction id
	LastBidIdKeyPrefix = []byte{0x12}

	AuctionKeyPrefix       = []byte{0x21}
	AllowedBidderKeyPrefix = []byte{0x22}

	BidKeyPrefix         = []byte{0x31}
	BidIndexKeyPrefix    = []byte{0x32}
	MatchedBidsLenPrefix = []byte{0x33}

	VestingQueueKeyPrefix = []byte{0x41}
)

// GetLastBidIdKey returns the store key to retrieve the latest bid id.
func GetLastBidIdKey(auctionId uint64) []byte {
	return append(LastBidIdKeyPrefix, sdk.Uint64ToBigEndian(auctionId)...)
}

// GetAuctionKey returns the store key to retrieve the auction object.
func GetAuctionKey(auctionId uint64) []byte {
	return append(AuctionKeyPrefix, sdk.Uint64ToBigEndian(auctionId)...)
}

// GetAllowedBidderKey returns the store key to retrieve the auction's allowed bidder object.
func GetAllowedBidderKey(auctionId uint64, bidder sdk.AccAddress) []byte {
	return append(append(AllowedBidderKeyPrefix, sdk.Uint64ToBigEndian(auctionId)...), address.MustLengthPrefix(bidder)...)
}

// GetAllowedBiddersByAuctionKeyPrefix returns the store key prefix to iterate allowed bidders by auction.
func GetAllowedBiddersByAuctionKeyPrefix(auctionId uint64) []byte {
	return append(AllowedBidderKeyPrefix, sdk.Uint64ToBigEndian(auctionId)...)
}

// GetBidKey returns the store key to retrieve the bid object.
func GetBidKey(auctionId uint64, bidId uint64) []byte {
	return append(append(BidKeyPrefix, sdk.Uint64ToBigEndian(auctionId)...), sdk.Uint64ToBigEndian(bidId)...)
}

// GetBidByAuctionIdPrefix returns the prefix to iterate all bids by the auction id.
func GetBidByAuctionIdPrefix(auctionId uint64) []byte {
	return append(BidKeyPrefix, sdk.Uint64ToBigEndian(auctionId)...)
}

// GetBidIndexKey returns the index key to retrieve the auction id and bid id to get the bid object.
func GetBidIndexKey(bidder sdk.AccAddress, auctionId uint64, bidId uint64) []byte {
	return append(append(append(BidIndexKeyPrefix, address.MustLengthPrefix(bidder)...), sdk.Uint64ToBigEndian(auctionId)...), sdk.Uint64ToBigEndian(bidId)...)
}

// GetBidByBidderPrefix returns the prefix to iterate all bids by a bidder.
func GetBidIndexByBidderPrefix(bidder sdk.AccAddress) []byte {
	return append(BidIndexKeyPrefix, address.MustLengthPrefix(bidder)...)
}

// GetVestingQueueKey returns the store key to retrieve the vesting queue from the index fields.
func GetVestingQueueKey(auctionId uint64, releaseTime time.Time) []byte {
	return append(append(VestingQueueKeyPrefix, sdk.Uint64ToBigEndian(auctionId)...), sdk.FormatTimeBytes(releaseTime)...)
}

// GetVestingQueueByAuctionIdPrefix returns a key prefix used to iterate vesting queues by an auction id.
func GetVestingQueueByAuctionIdPrefix(auctionId uint64) []byte {
	return append(VestingQueueKeyPrefix, sdk.Uint64ToBigEndian(auctionId)...)
}

func GetLastMatchedBidsLenKey(auctionId uint64) []byte {
	return append(MatchedBidsLenPrefix, sdk.Uint64ToBigEndian(auctionId)...)
}

// ParseBidIndexKey parses bid index key.
func ParseBidIndexKey(key []byte) (auctionId, bidId uint64) {
	if !bytes.HasPrefix(key, BidIndexKeyPrefix) {
		panic("key does not have proper prefix")
	}

	addrLen := key[1]
	bytesLen := 8
	auctionId = sdk.BigEndianToUint64(key[2+addrLen:])
	bidId = sdk.BigEndianToUint64(key[2+addrLen+byte(bytesLen):])
	return
}

// SplitAuctionIdBidIdKey splits the auction id and bid id.
func SplitAuctionIdBidIdKey(key []byte) (auctionId, bidId uint64) {
	bytesLen := 8
	auctionId = sdk.BigEndianToUint64(key)
	bidId = sdk.BigEndianToUint64(key[byte(bytesLen):])
	return
}
