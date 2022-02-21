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
	LastAuctionIdKey = []byte{0x11} // key to retrieve the latest auction id
	LastBidIdKey     = []byte{0x12} // key to retrieve the latest bid id number from the auction id

	AuctionKeyPrefix = []byte{0x21} // prefix to retrieve the auction from an auction id

	BidKeyPrefix      = []byte{0x31} // prefix to retrieve the bid from the auction id and bid id number
	BidIndexKeyPrefix = []byte{0x32} // prefix to retrieve the auction id and bid id by iterating the bidder address

	VestingQueueKeyPrefix = []byte{0x41} // prefix to retrieve the vesting queues from the auction id and vesting release time
)

// GetBidIdKey returns the store key to retrieve the latest Bid Id from the index fields.
func GetBidIdKey(auctionId uint64) []byte {
	return append(LastBidIdKey, sdk.Uint64ToBigEndian(auctionId)...)
}

// GetAuctionKey returns the store key to retrieve the auction from the index field.
func GetAuctionKey(auctionId uint64) []byte {
	return append(AuctionKeyPrefix, sdk.Uint64ToBigEndian(auctionId)...)
}

// GetBidKey returns the store key to retrieve the bid from the index fields.
func GetBidKey(auctionId uint64, bidId uint64) []byte {
	return append(append(BidKeyPrefix, sdk.Uint64ToBigEndian(auctionId)...), sdk.Uint64ToBigEndian(bidId)...)
}

// GetBidAuctionIDKey returns the store key to retrieve the bid from the auction id.
func GetBidAuctionIDKey(auctionId uint64) []byte {
	return append(BidKeyPrefix, sdk.Uint64ToBigEndian(auctionId)...)
}

// GetBidIndexKey returns the store key to retrieve the bid id from the index fields.
func GetBidIndexKey(bidderAddr sdk.AccAddress, auctionId uint64, bidId uint64) []byte {
	return append(append(append(BidIndexKeyPrefix, address.MustLengthPrefix(bidderAddr)...), sdk.Uint64ToBigEndian(auctionId)...), sdk.Uint64ToBigEndian(bidId)...)
}

// GetBidByBidderPrefix returns a key prefix used to iterate
// bids by a bidder.
func GetBidIndexByBidderPrefix(bidderAddr sdk.AccAddress) []byte {
	return append(BidIndexKeyPrefix, address.MustLengthPrefix(bidderAddr)...)
}

// GetVestingQueueKey returns the store key to retrieve the vesting queue from the index fields.
func GetVestingQueueKey(auctionId uint64, timestamp time.Time) []byte {
	return append(append(VestingQueueKeyPrefix, sdk.Uint64ToBigEndian(auctionId)...), sdk.FormatTimeBytes(timestamp)...)
}

// GetVestingQueueByAuctionIdPrefix returns a key prefix used to iterate
// vesting queues by an auction id.
func GetVestingQueueByAuctionIdPrefix(auctionId uint64) []byte {
	return append(VestingQueueKeyPrefix, sdk.Uint64ToBigEndian(auctionId)...)
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
