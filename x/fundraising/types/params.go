package types

import (
	"fmt"

	"gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter store keys.
var (
	KeyAuctionCreationFee  = []byte("AuctionCreationFee")
	KeyExtendedPeriod      = []byte("ExtendedPeriod")
	KeyAuctionFeeCollector = []byte("AuctionFeeCollector")

	DefaultAuctionCreationFee  = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100_000_000)))
	DefaultExtendedPeriod      = uint32(1)
	DefaultAuctionFeeCollector = sdk.AccAddress(address.Module(ModuleName, []byte("AuctionFeeCollector"))).String()
)

var _ paramstypes.ParamSet = (*Params)(nil)

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramstypes.KeyTable {
	return paramstypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns the default fundraising module parameters.
func DefaultParams() Params {
	return Params{
		AuctionCreationFee:  DefaultAuctionCreationFee,
		ExtendedPeriod:      DefaultExtendedPeriod,
		AuctionFeeCollector: DefaultAuctionFeeCollector,
	}
}

// ParamSetPairs implements paramstypes.ParamSet.
func (p *Params) ParamSetPairs() paramstypes.ParamSetPairs {
	return paramstypes.ParamSetPairs{
		paramstypes.NewParamSetPair(KeyAuctionCreationFee, &p.AuctionCreationFee, validateAuctionCreationFee),
		paramstypes.NewParamSetPair(KeyExtendedPeriod, &p.ExtendedPeriod, validateExtendedPeriod),
		paramstypes.NewParamSetPair(KeyAuctionFeeCollector, &p.AuctionFeeCollector, validateAuctionFeeCollector),
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
		{p.AuctionFeeCollector, validateAuctionFeeCollector},
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
	_, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateAuctionFeeCollector(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == "" {
		return fmt.Errorf("auction fee collector address must not be empty")
	}

	_, err := sdk.AccAddressFromBech32(v)
	if err != nil {
		return fmt.Errorf("invalid account address: %v", v)
	}

	return nil
}
