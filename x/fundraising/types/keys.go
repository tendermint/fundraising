package types

import (
	"bytes"

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
	SequenceKey  = []byte{0x12} // key to retrieve the latest sequence number

	AuctionKeyPrefix = []byte{0x21} // the prefix to retrieve the auction from an auction id

	SequenceIndexKeyPrefix = []byte{0x31} // the prefix to retrieve the bid from the combination of auction id and sequence number
	BidKeyPrefix           = []byte{0x32} // the prefix to retrieve the bid from the combination of auction id and auctioneer address
	BidderKeyPrefix        = []byte{0x33} // the prefix to retrieve the sequence number from the bidder address
)

// GetAuctionKey returns the store key to retrieve the auction from the index field.
func GetAuctionKey(auctionID uint64) []byte {
	return append(AuctionKeyPrefix, sdk.Uint64ToBigEndian(auctionID)...)
}

// GetSequenceIndexKey returns the store key to retrieve the bid from the index fields.
func GetSequenceIndexKey(auctionID uint64, sequence uint64) []byte {
	return append(append(SequenceIndexKeyPrefix, sdk.Uint64ToBigEndian(auctionID)...), sdk.Uint64ToBigEndian(sequence)...)
}

// GetBidKey returns the store key to retrieve the bid from the index fields.
func GetBidKey(auctionID uint64, auctioneerAcc sdk.AccAddress) []byte {
	return append(append(BidKeyPrefix, sdk.Uint64ToBigEndian(auctionID)...), address.MustLengthPrefix(auctioneerAcc)...)
}

// GetBidderKey returns the store key to retrieve the sequence number from the bidder address.
func GetBidderKey(bidderAcc sdk.AccAddress) []byte {
	return append(BidderKeyPrefix, address.MustLengthPrefix(bidderAcc)...)
}

// ParseSequenceIndexKey parses the store key to retrieve the index fields.
func ParseSequenceIndexKey(key []byte) (auctionID uint64, sequence uint64) {
	if !bytes.HasPrefix(key, SequenceIndexKeyPrefix) {
		panic("key does not have proper prefix")
	}

	bytesLen := 8
	auctionID = sdk.BigEndianToUint64(key[1:])
	sequence = sdk.BigEndianToUint64(key[1+bytesLen:])

	return
}

// ParseBidKey parses the store key to retrieve the index fields.
func ParseBidKey(key []byte) (auctionID uint64, auctioneerAcc sdk.AccAddress) {
	if !bytes.HasPrefix(key, BidKeyPrefix) {
		panic("key does not have proper prefix")
	}
	// TODO: not implemented yet
	return
}
