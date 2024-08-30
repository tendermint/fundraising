package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// MaxNumVestingSchedules is the maximum number of vesting schedules in an auction
	// It prevents from a malicious auctioneer to set an infinite number of vesting schedules
	// when they create an auction
	MaxNumVestingSchedules = 100

	// MaxExtendedRound is the maximum extend rounds for a batch auction to have
	// It prevents from a batch auction to extend its rounds forever
	MaxExtendedRound = 30
)

var (
	DefaultAuctionCreationFee = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(100_000_000)))
	DefaultPlaceBidFee        = sdk.Coins{}
	DefaultExtendedPeriod     = uint32(1)
)

// NewParams creates a new Params instance.
func NewParams(
	auctionCreationFee,
	placeBidFee sdk.Coins,
	extendedPeriod uint32,
) Params {
	return Params{AuctionCreationFee: auctionCreationFee, PlaceBidFee: placeBidFee, ExtendedPeriod: extendedPeriod}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(DefaultAuctionCreationFee, DefaultPlaceBidFee, DefaultExtendedPeriod)
}

// Validate validates the set of params.
func (p Params) Validate() error {
	if err := validateAuctionCreationFee(p.AuctionCreationFee); err != nil {
		return err
	}
	if err := validatePlaceBidFee(p.PlaceBidFee); err != nil {
		return err
	}
	if err := validateExtendedPeriod(p.ExtendedPeriod); err != nil {
		return err
	}

	return nil
}

func validateAuctionCreationFee(v sdk.Coins) error {
	return v.Validate()
}

func validatePlaceBidFee(v sdk.Coins) error {
	return v.Validate()
}

func validateExtendedPeriod(uint32) error {
	return nil
}
