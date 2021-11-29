package types

import (
	time "time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewVestingSchedule creates a new VestingSchedule object.
func NewVestingSchedule(time time.Time, weight sdk.Dec) VestingSchedule {
	return VestingSchedule{
		DistributedTime: time,
		Weight:          weight,
	}
}

// ValidateVestingSchedules validates the vesting schedules.
// Each weight of the vesting schedule must be positive and total weight must be equal to 1.
// If a number of schedule equals to zero, the auctioneer doesn't want any vesting schedule.
func ValidateVestingSchedules(schedules []VestingSchedule) error {
	if len(schedules) == 0 {
		return nil
	}

	totalWeight := sdk.ZeroDec()
	for _, s := range schedules {
		if !s.Weight.IsPositive() {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "vesting weight must be positive")
		}
		if s.Weight.GT(sdk.OneDec()) {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "total vesting weight must not greater than 1")
		}
		totalWeight = totalWeight.Add(s.Weight)
	}

	if !totalWeight.Equal(sdk.OneDec()) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "total vesting weight must be equal to 1")
	}
	return nil
}
