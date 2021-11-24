package types

import (
	"gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter store keys
var (
	KeyFundraisingCreationFee = []byte("FundraisingCreationFee")
	KeyExtendedPeriod         = []byte("ExtendedPeriod")

	DefaultFundraisingCreationFee = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100_000_000)))
	DefaultExtendedPeriod         = uint32(1)
)

var _ paramstypes.ParamSet = (*Params)(nil)

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramstypes.KeyTable {
	return paramstypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns the default farming module parameters.
func DefaultParams() Params {
	return Params{
		FundraisingCreationFee: DefaultFundraisingCreationFee,
		ExtendedPeriod:         DefaultExtendedPeriod,
	}
}

// ParamSetPairs implements paramstypes.ParamSet.
func (p *Params) ParamSetPairs() paramstypes.ParamSetPairs {
	return paramstypes.ParamSetPairs{
		paramstypes.NewParamSetPair(KeyFundraisingCreationFee, &p.FundraisingCreationFee, validateFundraisingCreationFee),
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
		{p.FundraisingCreationFee, validateFundraisingCreationFee},
		{p.ExtendedPeriod, validateExtendedPeriod},
	} {
		if err := v.validator(v.value); err != nil {
			return err
		}
	}
	return nil
}

func validateFundraisingCreationFee(i interface{}) error {
	// TODO: not implmented yet
	return nil
}

func validateExtendedPeriod(i interface{}) error {
	// TODO: not implmented yet
	return nil
}
