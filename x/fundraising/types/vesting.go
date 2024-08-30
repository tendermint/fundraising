package types

import (
	"time"

	sdkerrors "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ValidateVestingSchedules validates the vesting schedules.
// Each weight of the vesting schedule must be positive and total weight must be equal to 1.
// If a number of schedule equals to zero, the auctioneer doesn't want any vesting schedule.
// The release times must be chronological for vesting schedules. Otherwise it returns an error.
func ValidateVestingSchedules(schedules []VestingSchedule, endTime time.Time) error {
	if len(schedules) == 0 {
		return nil
	}

	totalWeight := math.LegacyZeroDec()
	prevReleaseTime := MustParseRFC3339("0001-01-01T00:00:00Z")

	for _, s := range schedules {
		if !s.Weight.IsPositive() {
			return sdkerrors.Wrapf(ErrInvalidVestingSchedules, "vesting weight must be positive")
		}

		if !s.ReleaseTime.After(endTime) { // ReleaseTime <= EndTime
			return sdkerrors.Wrapf(ErrInvalidVestingSchedules, "release time must be set after the end time")
		}

		if !s.ReleaseTime.After(prevReleaseTime) {
			return sdkerrors.Wrapf(ErrInvalidVestingSchedules, "release time must be chronological")
		}

		if s.Weight.GT(math.LegacyOneDec()) {
			return sdkerrors.Wrapf(ErrInvalidVestingSchedules, "vesting weight must not be greater than 1")
		}

		totalWeight = totalWeight.Add(s.Weight)
		prevReleaseTime = s.ReleaseTime
	}

	if !totalWeight.Equal(math.LegacyOneDec()) {
		return sdkerrors.Wrapf(ErrInvalidVestingSchedules, "total vesting weight must be equal to 1")
	}

	return nil
}

// NewVestingQueue returns a new VestingQueue.
func NewVestingQueue(auctionId uint64, auctioneer sdk.AccAddress, payingCoin sdk.Coin, releaseTime time.Time, released bool) VestingQueue {
	return VestingQueue{
		AuctionId:   auctionId,
		Auctioneer:  auctioneer.String(),
		PayingCoin:  payingCoin,
		ReleaseTime: releaseTime,
		Released:    released,
	}
}

// ShouldRelease returns true when the vesting queue is ready to release the paying coin.
// It checks if the release time is equal or before the given time t and released value is false.
func (vq VestingQueue) ShouldRelease(t time.Time) bool {
	return !vq.GetReleaseTime().After(t) && !vq.Released
}

// SetReleased sets released status of the vesting queue.
func (vq *VestingQueue) SetReleased(status bool) {
	vq.Released = status
}
