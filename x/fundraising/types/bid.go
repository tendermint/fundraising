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

func (b *Bid) SetMatched(status bool) {
	b.IsMatched = status
}

func (b Bid) GetBidSellingAmount(payingDenom string) sdk.Int {
	var bidSellingAmt sdk.Int

	if b.Coin.Denom == payingDenom {
		bidSellingAmt = b.Coin.Amount.ToDec().QuoTruncate(b.Price).TruncateInt()
	} else {
		bidSellingAmt = b.Coin.Amount
	}

	return bidSellingAmt
}

func (b Bid) GetBidPayingAmount(payingDenom string) sdk.Int {
	var bidPayingAmt sdk.Int

	if b.Coin.Denom == payingDenom {
		bidPayingAmt = b.Coin.Amount
	} else {
		bidPayingAmt = b.Coin.Amount.ToDec().Mul(b.Price).Ceil().TruncateInt()
	}

	return bidPayingAmt
}
