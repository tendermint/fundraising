package types

import (
	"fmt"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// DefaultGenesisState returns the default fundraising genesis state
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params:               DefaultParams(),
		Auctions:             []*codectypes.Any{},
		AllowedBidderRecords: []AllowedBidderRecord{},
		Bids:                 []Bid{},
		VestingQueues:        []VestingQueue{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}

	for _, a := range gs.Auctions {
		auction, err := UnpackAuction(a)
		if err != nil {
			return err
		}

		if err := auction.Validate(); err != nil {
			return err
		}
	}

	for _, r := range gs.AllowedBidderRecords {
		if err := r.Validate(); err != nil {
			return err
		}
	}

	for _, b := range gs.Bids {
		if err := b.Validate(); err != nil {
			return err
		}
	}

	for _, q := range gs.VestingQueues {
		if err := q.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// Validate validates AllowedBidder.
func (r AllowedBidderRecord) Validate() error {
	if r.AuctionId == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "auction id cannot be 0")
	}
	return r.AllowedBidder.Validate()
}

// Validate validates Bid.
func (b Bid) Validate() error {
	if _, err := sdk.AccAddressFromBech32(b.Bidder); err != nil {
		return err
	}
	if !b.Price.IsPositive() {
		return fmt.Errorf("bid price must be positive value: %s", b.Price.String())
	}
	if err := b.Coin.Validate(); err != nil {
		return err
	}
	if !b.Coin.Amount.IsPositive() {
		return fmt.Errorf("coin amount must be positive: %s", b.Coin.Amount.String())
	}
	return nil
}

// Validate validates VestingQueue.
func (q VestingQueue) Validate() error {
	if _, err := sdk.AccAddressFromBech32(q.Auctioneer); err != nil {
		return err
	}
	if err := q.PayingCoin.Validate(); err != nil {
		return fmt.Errorf("paying coin is invalid: %v", err)
	}
	return nil
}
