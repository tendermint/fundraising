package types

import (
	"sort"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	"github.com/tendermint/tendermint/crypto"
)

// MustParseRFC3339 parses string time to time in RFC3339 format.
// This is used only for internal testing purpose.
func MustParseRFC3339(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

// DeriveAddress derives an address with the given address length type, module name, and
// address derivation name. It is used to derive reserve account addresses for selling, paying, and vesting.
func DeriveAddress(addressType AddressType, moduleName, name string) sdk.AccAddress {
	switch addressType {
	case AddressType32Bytes:
		return sdk.AccAddress(address.Module(moduleName, []byte(name)))
	case AddressType20Bytes:
		return sdk.AccAddress(crypto.AddressHash([]byte(moduleName + name)))
	default:
		return sdk.AccAddress{}
	}
}

// SortByBidPrice sorts bid array by bid price in descending order.
func SortByBidPrice(bids []Bid) []Bid {
	sort.Slice(bids, func(i, j int) bool {
		if bids[i].Price.GT(bids[j].Price) {
			return true
		}
		return bids[i].Id < bids[j].Id
	})
	return bids
}
