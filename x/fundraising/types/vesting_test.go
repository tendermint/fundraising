package types_test

import (
	"testing"
	time "time"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

func TestShouldRelease(t *testing.T) {
	now := types.MustParseRFC3339("2021-12-10T00:00:00Z")

	testCases := []struct {
		name      string
		vq        types.VestingQueue
		expResult bool
	}{
		{
			"the release time is already passed the current block time",
			types.VestingQueue{
				AuctionId:   1,
				Auctioneer:  sdk.AccAddress(crypto.AddressHash([]byte("Auctioneer"))).String(),
				PayingCoin:  sdk.NewInt64Coin("denom1", 10000000),
				ReleaseTime: types.MustParseRFC3339("2021-11-01T00:00:00Z"),
				Released:    false,
			},
			true,
		},
		{
			"the release time is exactly the same time as the current block time",
			types.VestingQueue{
				AuctionId:   1,
				Auctioneer:  sdk.AccAddress(crypto.AddressHash([]byte("Auctioneer"))).String(),
				PayingCoin:  sdk.NewInt64Coin("denom1", 10000000),
				ReleaseTime: now,
				Released:    false,
			},
			true,
		},
		{
			"the release time has not passed the current block time",
			types.VestingQueue{
				AuctionId:   1,
				Auctioneer:  sdk.AccAddress(crypto.AddressHash([]byte("Auctioneer"))).String(),
				PayingCoin:  sdk.NewInt64Coin("denom1", 10000000),
				ReleaseTime: types.MustParseRFC3339("2022-01-30T00:00:00Z"),
				Released:    false,
			},
			false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expResult, tc.vq.ShouldRelease(now))
		})
	}
}

func TestValidateVestingSchedules(t *testing.T) {
	for _, tc := range []struct {
		name        string
		schedules   []types.VestingSchedule
		endTime     time.Time
		expectedErr string
	}{
		{
			"happy case",
			[]types.VestingSchedule{
				{ReleaseTime: types.MustParseRFC3339("9999-01-01T00:00:00Z"), Weight: sdk.MustNewDecFromStr("1.0")},
			},
			types.MustParseRFC3339("2022-03-01T00:00:00Z"),
			"",
		},
		{
			"invalid case #1",
			[]types.VestingSchedule{
				{ReleaseTime: types.MustParseRFC3339("9999-01-01T00:00:00Z"), Weight: sdk.MustNewDecFromStr("-1.0")},
			},
			types.MustParseRFC3339("2022-03-01T00:00:00Z"),
			"vesting weight must be positive: invalid vesting schedules",
		},
		{
			"invalid case #2",
			[]types.VestingSchedule{
				{ReleaseTime: types.MustParseRFC3339("2022-01-01T00:00:00Z"), Weight: sdk.MustNewDecFromStr("1.0")},
			},
			types.MustParseRFC3339("2022-03-01T00:00:00Z"),
			"release time must be set after the end time: invalid vesting schedules",
		},
		{
			"invalid case #3",
			[]types.VestingSchedule{
				{ReleaseTime: types.MustParseRFC3339("9999-01-01T00:00:00Z"), Weight: sdk.MustNewDecFromStr("2.0")},
			},
			types.MustParseRFC3339("2022-03-01T00:00:00Z"),
			"vesting weight must not be greater than 1: invalid vesting schedules",
		},
		{
			"invalid case #4",
			[]types.VestingSchedule{
				{ReleaseTime: types.MustParseRFC3339("2022-06-01T00:00:00Z"), Weight: sdk.MustNewDecFromStr("0.25")},
				{ReleaseTime: types.MustParseRFC3339("2022-04-01T00:00:00Z"), Weight: sdk.MustNewDecFromStr("0.25")},
				{ReleaseTime: types.MustParseRFC3339("2022-09-01T00:00:00Z"), Weight: sdk.MustNewDecFromStr("0.25")},
				{ReleaseTime: types.MustParseRFC3339("2022-12-01T00:00:00Z"), Weight: sdk.MustNewDecFromStr("0.25")},
			},
			types.MustParseRFC3339("2022-03-01T00:00:00Z"),
			"release time must be chronological: invalid vesting schedules",
		},
		{
			"invalid case #5",
			[]types.VestingSchedule{
				{ReleaseTime: types.MustParseRFC3339("2022-05-01T00:00:00Z"), Weight: sdk.MustNewDecFromStr("0.5")},
				{ReleaseTime: types.MustParseRFC3339("2022-06-01T00:00:00Z"), Weight: sdk.MustNewDecFromStr("0.5")},
				{ReleaseTime: types.MustParseRFC3339("2022-07-01T00:00:00Z"), Weight: sdk.MustNewDecFromStr("0.5")},
				{ReleaseTime: types.MustParseRFC3339("2022-08-01T00:00:00Z"), Weight: sdk.MustNewDecFromStr("0.5")},
			},
			types.MustParseRFC3339("2022-03-01T00:00:00Z"),
			"total vesting weight must be equal to 1: invalid vesting schedules",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := types.ValidateVestingSchedules(tc.schedules, tc.endTime)
			if tc.expectedErr == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedErr)
			}
		})
	}
}

func TestSetReleased(t *testing.T) {
	vestingQueue := types.NewVestingQueue(
		1,
		sdk.AccAddress(crypto.AddressHash([]byte("Auctioneer"))),
		sdk.NewInt64Coin("denom1", 10000000),
		types.MustParseRFC3339("2021-11-01T00:00:00Z"),
		false,
	)
	require.False(t, vestingQueue.Released)

	vestingQueue.SetReleased(true)
	require.True(t, vestingQueue.Released)
}
