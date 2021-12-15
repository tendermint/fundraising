package types

import (
	"fmt"

	"gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter store keys.
var (
	KeyAuctionCreationFee = []byte("AuctionCreationFee")
	KeyExtendedPeriod     = []byte("ExtendedPeriod")

	DefaultAuctionCreationFee = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100_000_000)))
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
		ExtendedPeriod:     DefaultExtendedPeriod,
	}
}

// ParamSetPairs implements paramstypes.ParamSet.
func (p *Params) ParamSetPairs() paramstypes.ParamSetPairs {
	return paramstypes.ParamSetPairs{
		paramstypes.NewParamSetPair(KeyAuctionCreationFee, &p.AuctionCreationFee, validateAuctionCreationFee),
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

func validateExtendedPeriod(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("extended period must be positive: %d", v)
	}

	return nil
}
