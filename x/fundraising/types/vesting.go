package types

import (
	time "time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewVestingSchedule creates a new VestingSchedule object.
func NewVestingSchedule(time time.Time, weight sdk.Dec) VestingSchedule {
	return VestingSchedule{
		ReleaseTime: time,
		Weight:      weight,
	}
}

// NewVestingQueue creates a new VestingQueue object.
func NewVestingQueue(auctionId uint64, auctioneer string, payingCoin sdk.Coin, releaseTime time.Time) VestingQueue {
	return VestingQueue{
		AuctionId:   auctionId,
		Auctioneer:  auctioneer,
		PayingCoin:  payingCoin,
		ReleaseTime: releaseTime,
	}
}

// ValidateVestingSchedules validates the vesting schedules.
// Each weight of the vesting schedule must be positive and total weight must be equal to 1.
// If a number of schedule equals to zero, the auctioneer doesn't want any vesting schedule.
// The release times must be chronological for vesting schedules. Otherwise it returns an error.
func ValidateVestingSchedules(schedules []VestingSchedule) error {
	if len(schedules) == 0 {
		return nil
	}

	// initialize timestamp with max time and total weight with zero
	ts := ParseTime("0001-01-01T00:00:00Z")
	totalWeight := sdk.ZeroDec()

	for _, s := range schedules {
		if !s.Weight.IsPositive() {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "vesting weight must be positive")
		}
		if s.Weight.GT(sdk.OneDec()) {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "each vesting weight must not be greater than 1")
		}
		totalWeight = totalWeight.Add(s.Weight)

		if !s.ReleaseTime.After(ts) {
			return sdkerrors.Wrapf(ErrInvalidVestingSchedules, "release time must be chronological")
		}
		ts = s.ReleaseTime
	}

	if !totalWeight.Equal(sdk.OneDec()) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "total vesting weight must be equal to 1")
	}

	return nil
}
