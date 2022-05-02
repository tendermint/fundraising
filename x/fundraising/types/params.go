package types

import (
	"fmt"

	"gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
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

// Parameter store keys.
var (
	KeyAuctionCreationFee = []byte("AuctionCreationFee")
	KeyPlaceBidFee        = []byte("PlaceBidFee")
	KeyExtendedPeriod     = []byte("ExtendedPeriod")

	DefaultAuctionCreationFee = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100_000_000)))
	DefaultPlaceBidFee        = sdk.Coins{}
	DefaultExtendedPeriod     = uint32(1)
)

var _ paramstypes.ParamSet = (*Params)(nil)

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramstypes.KeyTable {
	return paramstypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns the default fundraising module parameters.
func DefaultParams() Params {
	return Params{
		AuctionCreationFee: DefaultAuctionCreationFee,
		PlaceBidFee:        DefaultPlaceBidFee,
		ExtendedPeriod:     DefaultExtendedPeriod,
	}
}

// ParamSetPairs implements paramstypes.ParamSet.
func (p *Params) ParamSetPairs() paramstypes.ParamSetPairs {
	return paramstypes.ParamSetPairs{
		paramstypes.NewParamSetPair(KeyAuctionCreationFee, &p.AuctionCreationFee, validateAuctionCreationFee),
		paramstypes.NewParamSetPair(KeyPlaceBidFee, &p.PlaceBidFee, validatePlaceBidFee),
		paramstypes.NewParamSetPair(KeyExtendedPeriod, &p.ExtendedPeriod, validateExtendedPeriod),
	}
}

// String returns a human readable string representation of the parameters.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// Validate validates parameters.
func (p Params) Validate() error {
	for _, v := range []struct {
		value     interface{}
		validator func(interface{}) error
	}{
		{p.AuctionCreationFee, validateAuctionCreationFee},
		{p.ExtendedPeriod, validateExtendedPeriod},
	} {
		if err := v.validator(v.value); err != nil {
			return err
		}
	}
	return nil
}

func validateAuctionCreationFee(i interface{}) error {
	v, ok := i.(sdk.Coins)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if err := v.Validate(); err != nil {
		return err
	}

	return nil
}

func validatePlaceBidFee(i interface{}) error {
	v, ok := i.(sdk.Coins)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if err := v.Validate(); err != nil {
		return err
	}

	return nil
}

func validateExtendedPeriod(i interface{}) error {
	_, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}
