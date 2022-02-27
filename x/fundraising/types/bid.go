package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (b Bid) GetBidder() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(b.Bidder)
	if err != nil {
		panic(err)
	}
	return addr
}

func (b *Bid) SetWinner(status bool) {
	b.IsWinner = status
}

func (b Bid) GetExchangedSellingAmount() sdk.Int {
	return b.Coin.Amount.ToDec().QuoTruncate(b.Price).TruncateInt()
}
