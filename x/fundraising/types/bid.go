package types

import (
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (b Bid) GetBidder() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(b.Bidder)
	if err != nil {
		panic(err)
	}
	return addr
}

// SanitizeReverseBids sorts bids in descending order.
func SanitizeReverseBids(bids []Bid) []Bid {
	sort.SliceStable(bids, func(i, j int) bool {
		return bids[i].Price.GT(bids[j].Price)
	})
	return bids
}

func (b *Bid) SetWinner(status bool) {
	b.IsWinner = status
}

func (b Bid) GetExchangedSellingAmount() sdk.Int {
	return b.Coin.Amount.ToDec().QuoTruncate(b.Price).TruncateInt()
}
