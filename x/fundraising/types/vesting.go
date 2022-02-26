package types

import (
	time "time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// ValidateVestingSchedules validates the vesting schedules.
// Each weight of the vesting schedule must be positive and total weight must be equal to 1.
// If a number of schedule equals to zero, the auctioneer doesn't want any vesting schedule.
// The release times must be chronological for vesting schedules. Otherwise it returns an error.
func ValidateVestingSchedules(schedules []VestingSchedule, endTime time.Time) error {
	if len(schedules) == 0 {
		return nil
	}

	// initialize timestamp with max time and total weight with zero
	ts := MustParseRFC3339("0001-01-01T00:00:00Z")
	totalWeight := sdk.ZeroDec()

	for _, s := range schedules {
		if !s.Weight.IsPositive() {
			return sdkerrors.Wrapf(ErrInvalidVestingSchedules, "vesting weight must be positive")
		}

		if s.Weight.GT(sdk.OneDec()) {
			return sdkerrors.Wrapf(ErrInvalidVestingSchedules, "each vesting weight must not be greater than 1")
		}
		totalWeight = totalWeight.Add(s.Weight)

		if !s.ReleaseTime.After(endTime) {
			return sdkerrors.Wrapf(ErrInvalidVestingSchedules, "release time must be after the end time")
		}

		if !s.ReleaseTime.After(ts) {
			return sdkerrors.Wrapf(ErrInvalidVestingSchedules, "release time must be chronological")
		}
		ts = s.ReleaseTime
	}

	if !totalWeight.Equal(sdk.OneDec()) {
		return sdkerrors.Wrapf(ErrInvalidVestingSchedules, "total vesting weight must be equal to 1")
	}

	return nil
}

// ShouldRelease returns true when the vesting queue is ready to release the paying coin.
// It checks if the release time is equal or before the given time t and released value is false.
func (vq VestingQueue) ShouldRelease(t time.Time) bool {
	return !vq.GetReleaseTime().After(t) && !vq.Released
}

func (vq *VestingQueue) SetReleased(status bool) {
	vq.Released = status
}
