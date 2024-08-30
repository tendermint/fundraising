package types

import (
	"sort"
	"time"

	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
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

// SortBids sorts bid array by bid price in descending order.
func SortBids(bids []Bid) []Bid {
	sort.Slice(bids, func(i, j int) bool {
		if bids[i].Price.GT(bids[j].Price) {
			return true
		}
		return bids[i].Id < bids[j].Id
	})
	return bids
}

func BidsByPrice(bids []Bid) (prices []math.LegacyDec, bidsByPrice map[string][]Bid) {
	bids = SortBids(bids)

	bidsByPrice = map[string][]Bid{} // price => []Bid

	for _, bid := range bids {
		priceStr := bid.Price.String()
		bidsByPrice[priceStr] = append(bidsByPrice[priceStr], bid)
	}

	// Sort prices in descending order.
	prices = make([]math.LegacyDec, len(bidsByPrice))
	i := 0
	for priceStr := range bidsByPrice {
		prices[i] = math.LegacyMustNewDecFromStr(priceStr)
		i++
	}
	sort.Slice(prices, func(i, j int) bool {
		return prices[i].GT(prices[j])
	})
	return
}
