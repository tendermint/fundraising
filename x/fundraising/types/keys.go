package types

import (
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
	AuctionIdKey = []byte{0x11} // the key to retrieve the latest auction id

	AuctionKeyPrefix = []byte{0x21} // the prefix to retrieve the auction from an auction id

	BidKeyPrefix    = []byte{0x31} // the prefix to retrieve the bid from the  auction id
	BidderKeyPrefix = []byte{0x32} // the prefix to retrieve the bid from the bidder address
)

// GetAuctionKey returns the store key to retrieve the auction from the index field.
func GetAuctionKey(auctionID uint64) []byte {
	return append(AuctionKeyPrefix, sdk.Uint64ToBigEndian(auctionID)...)
}

// GetBidKey returns the store key to retrieve the bid from the index fields.
func GetBidKey(auctionID uint64) []byte {
	return append(BidKeyPrefix, sdk.Uint64ToBigEndian(auctionID)...)
}

// GetBidderKey returns the store key to retrieve the sequence number from the bidder address.
func GetBidderKey(bidderAcc sdk.AccAddress) []byte {
	return append(BidderKeyPrefix, address.MustLengthPrefix(bidderAcc)...)
}
