package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewAllowedBidder returns a new AllowedBidder.
func NewAllowedBidder(bidderAddr sdk.AccAddress, maxBidAmount sdk.Int) AllowedBidder {
	return AllowedBidder{
		Bidder:       bidderAddr.String(),
		MaxBidAmount: maxBidAmount,
	}
}

// GetBidder returns the bidder account address.
func (ab AllowedBidder) GetBidder() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(ab.Bidder)
	if err != nil {
		panic(err)
	}
	return addr
}

// Validate validates allowed bidder object.
func (ab AllowedBidder) Validate() error {
	if _, err := sdk.AccAddressFromBech32(ab.Bidder); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}
	if ab.MaxBidAmount.IsNil() {
		return ErrInvalidMaxBidAmount
	}
	if !ab.MaxBidAmount.IsPositive() {
		return ErrInvalidMaxBidAmount
	}
	return nil
}
