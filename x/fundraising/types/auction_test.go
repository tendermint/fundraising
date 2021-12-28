package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

func TestIsAuctionStarted(t *testing.T) {
	now := types.ParseTime("2021-12-01T00:00:00Z")

	for _, tc := range []struct {
		startTimeStr string
		active       bool
	}{
		{"2021-11-01T00:00:00Z", true},
		{"2021-11-15T23:59:59Z", true},
		{"2021-11-20T00:00:00Z", true},
		{"2021-12-01T00:00:00Z", true},
		{"2021-12-01T00:00:01Z", false},
		{"2021-12-10T00:00:00Z", false},
	} {
		require.Equal(t, tc.active, types.IsAuctionStarted(types.ParseTime(tc.startTimeStr), now))
	}
}

func TestIsAuctionFinished(t *testing.T) {
	now := types.ParseTime("2021-12-10T00:00:00Z")

	for _, tc := range []struct {
		endTimeStr string
		active     bool
	}{
		{"2021-11-01T00:00:00Z", true},
		{"2021-11-15T23:59:59Z", true},
		{"2021-11-20T00:00:00Z", true},
		{"2021-12-10T00:00:00Z", true},
		{"2021-12-11T00:00:00Z", false},
		{"2021-12-30T00:00:00Z", false},
	} {
		require.Equal(t, tc.active, types.IsAuctionFinished(types.ParseTime(tc.endTimeStr), now))
	}
}
