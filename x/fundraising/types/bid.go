package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func (b Bid) GetBidder() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(b.Bidder)
	if err != nil {
		panic(err)
	}
	return addr
}

// func (b Bid) CalculateReceivingSellingCoin(sellingCoinDenom string) sdk.Coin {
// 	receiveCoin := sdk.Coin{}
// 	switch b.Type {
// 	case BidTypeFixedPrice:
// 		receiveAmt := b.Coin.Amount.ToDec().QuoTruncate(b.Price).TruncateInt()
// 		receiveCoin = sdk.NewCoin(sellingCoinDenom, receiveAmt)
// 	case BidTypeBatchWorth:
// 		// TODO: not implemented yet
// 	case BidTypeBatchMany:
// 		// TODO: not implemented yet
// 	}

// 	return receiveCoin
// }
