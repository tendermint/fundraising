package types

import (
	sdkerrors "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

// NewAllowedBidder returns a new AllowedBidder.
func NewAllowedBidder(auctionId uint64, bidderAddr sdk.AccAddress, maxBidAmount math.Int) AllowedBidder {
	return AllowedBidder{
		AuctionId:    auctionId,
		Bidder:       bidderAddr.String(),
		MaxBidAmount: maxBidAmount,
	}
}

// GetBidder returns the bidder account address.
func (ab AllowedBidder) GetBidder() (sdk.AccAddress, error) {
	return sdk.AccAddressFromBech32(ab.Bidder)
}

// Validate validates allowed bidder object.
func (ab AllowedBidder) Validate() error {
	if _, err := ab.GetBidder(); err != nil {
		return sdkerrors.Wrap(errors.ErrInvalidAddress, err.Error())
	}
	if ab.MaxBidAmount.IsNil() {
		return ErrInvalidMaxBidAmount
	}
	if !ab.MaxBidAmount.IsPositive() {
		return ErrInvalidMaxBidAmount
	}
	return nil
}
