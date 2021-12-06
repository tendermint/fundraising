package types

import (
	"bytes"
	fmt "fmt"
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
	AuctionIdKey = []byte{0x11} // key to retrieve the latest auction id
	SequenceKey  = []byte{0x12} // key to retrieve the latest sequence number from the auction id

	AuctionKeyPrefix = []byte{0x21} // prefix to retrieve the auction from an auction id

	BidKeyPrefix      = []byte{0x31} // prefix to retrieve the bid from the auction id and sequence number
	BidIndexKeyPrefix = []byte{0x32} // prefix to retrieve the auction id and sequence by iterating the bidder address

	VestingQueueKeyPrefix = []byte{0x41} // prefix to retrieve the vesting queues from the auction id and vesting release time
)

// GetAuctionKey returns the store key to retrieve the auction from the index field.
func GetAuctionKey(auctionID uint64) []byte {
	return append(AuctionKeyPrefix, sdk.Uint64ToBigEndian(auctionID)...)
}

// GetBidKey returns the store key to retrieve the bid from the index fields.
func GetBidKey(auctionID uint64, sequence uint64) []byte {
	return append(append(BidKeyPrefix, sdk.Uint64ToBigEndian(auctionID)...), sdk.Uint64ToBigEndian(sequence)...)
}

// GetBidIndexKey returns the store key to retrieve the sequence number from the index fields.
func GetBidIndexKey(bidderAcc sdk.AccAddress, auctionID uint64, sequence uint64) []byte {
	return append(append(append(BidIndexKeyPrefix, address.MustLengthPrefix(bidderAcc)...), sdk.Uint64ToBigEndian(auctionID)...), sdk.Uint64ToBigEndian(sequence)...)
}

// GetBidByBidderPrefix returns a key prefix used to iterate
// bids by a bidder.
func GetBidIndexByBidderPrefix(bidderAcc sdk.AccAddress) []byte {
	return append(BidIndexKeyPrefix, address.MustLengthPrefix(bidderAcc)...)
}

// GetVestingQueueKey returns the store key to retrieve the vesting queue from the index fields.
func GetVestingQueueKey(timestamp time.Time, auctionID uint64) []byte {
	idBz := sdk.Uint64ToBigEndian(auctionID)
	timeBz := sdk.FormatTimeBytes(timestamp)
	timeBzL := len(timeBz)
	prefixL := len(VestingQueueKeyPrefix)

	bz := make([]byte, prefixL+8+timeBzL+8)

	// copy the prefix
	copy(bz[:prefixL], VestingQueueKeyPrefix)

	// copy the encoded time bytes length
	copy(bz[prefixL:prefixL+8], sdk.Uint64ToBigEndian(uint64(timeBzL)))

	// copy the encoded time bytes
	copy(bz[prefixL+8:prefixL+8+timeBzL], timeBz)

	// copy the encoded auction id
	copy(bz[prefixL+8+timeBzL:], idBz)

	return bz
}

// ParseBidKey returnes the auction id and sequence from a key created
// from GetBidKey.
func ParseBidKey(key []byte) (auctionID uint64, sequence uint64) {
	if !bytes.HasPrefix(key, BidKeyPrefix) {
		panic("key does not have proper prefix")
	}
	bytesLen := 8
	auctionID = sdk.BigEndianToUint64(key[1:])
	sequence = sdk.BigEndianToUint64(key[1+bytesLen:])
	return
}

func ParseBidIndexKey(key []byte) (auctionID, sequence uint64) {
	if !bytes.HasPrefix(key, BidIndexKeyPrefix) {
		panic("key does not have proper prefix")
	}

	addrLen := key[1]
	bytesLen := 8
	auctionID = sdk.BigEndianToUint64(key[2+addrLen:])
	sequence = sdk.BigEndianToUint64(key[2+addrLen+byte(bytesLen):])
	return
}

// ParseVestingQueueKey returns the encoded time and auction id from a key created
// from GetVestingQueueKey.
func ParseVestingQueueKey(bz []byte) (time.Time, uint64, error) {
	prefixL := len(VestingQueueKeyPrefix)
	if prefix := bz[:prefixL]; !bytes.Equal(prefix, VestingQueueKeyPrefix) {
		return time.Time{}, 0, fmt.Errorf("invalid prefix; expected: %X, got: %X", VestingQueueKeyPrefix, prefix)
	}

	timeBzL := sdk.BigEndianToUint64(bz[prefixL : prefixL+8])
	ts, err := sdk.ParseTimeBytes(bz[prefixL+8 : prefixL+8+int(timeBzL)])
	if err != nil {
		return time.Time{}, 0, err
	}

	auctionID := sdk.BigEndianToUint64(bz[prefixL+8+int(timeBzL):])

	return ts, auctionID, nil
}
