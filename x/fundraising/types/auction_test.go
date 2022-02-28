package types_test

import (
	"testing"
	time "time"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

func TestShouldAuctionStarted(t *testing.T) {
	auction := types.BaseAuction{
		Id:                    1,
		Type:                  types.AuctionTypeFixedPrice,
		AllowedBidders:        nil,
		Auctioneer:            sdk.AccAddress(crypto.AddressHash([]byte("Auctioneer"))).String(),
		SellingReserveAddress: types.SellingReserveAddress(1).String(),
		PayingReserveAddress:  types.PayingReserveAddress(1).String(),
		StartPrice:            sdk.MustNewDecFromStr("0.5"),
		MinBidPrice:           sdk.MustNewDecFromStr("0.1"),
		SellingCoin:           sdk.NewInt64Coin("denom3", 1_000_000_000_000),
		PayingCoinDenom:       "denom4",
		VestingReserveAddress: types.VestingReserveAddress(1).String(),
		VestingSchedules:      []types.VestingSchedule{},
		MatchedPrice:          sdk.ZeroDec(),
		RemainingSellingCoin:  sdk.NewInt64Coin("denom3", 1_000_000_000_000),
		StartTime:             types.MustParseRFC3339("2021-12-01T00:00:00Z"),
		EndTimes:              []time.Time{types.MustParseRFC3339("2021-12-15T00:00:00Z")},
		Status:                types.AuctionStatusStandBy,
	}

	for _, tc := range []struct {
		currentTime string
		expected    bool
	}{
		{"2021-11-01T00:00:00Z", false},
		{"2021-11-15T23:59:59Z", false},
		{"2021-11-20T00:00:00Z", false},
		{"2021-12-01T00:00:00Z", true},
		{"2021-12-01T00:00:01Z", true},
		{"2021-12-10T00:00:00Z", true},
		{"2022-01-01T00:00:00Z", true},
	} {
		require.Equal(t, tc.expected, auction.ShouldAuctionStarted(types.MustParseRFC3339(tc.currentTime)))
	}
}

func TestShouldAuctionFinished(t *testing.T) {
	auction := types.BaseAuction{
		Id:                    1,
		Type:                  types.AuctionTypeFixedPrice,
		AllowedBidders:        nil,
		Auctioneer:            sdk.AccAddress(crypto.AddressHash([]byte("Auctioneer"))).String(),
		SellingReserveAddress: types.SellingReserveAddress(1).String(),
		PayingReserveAddress:  types.PayingReserveAddress(1).String(),
		StartPrice:            sdk.MustNewDecFromStr("0.5"),
		MinBidPrice:           sdk.MustNewDecFromStr("0.1"),
		SellingCoin:           sdk.NewInt64Coin("denom3", 1_000_000_000_000),
		PayingCoinDenom:       "denom4",
		VestingReserveAddress: types.VestingReserveAddress(1).String(),
		VestingSchedules:      []types.VestingSchedule{},
		MatchedPrice:          sdk.ZeroDec(),
		RemainingSellingCoin:  sdk.NewInt64Coin("denom3", 1_000_000_000_000),
		StartTime:             types.MustParseRFC3339("2021-12-01T00:00:00Z"),
		EndTimes:              []time.Time{types.MustParseRFC3339("2021-12-15T00:00:00Z")},
		Status:                types.AuctionStatusStandBy,
	}

	for _, tc := range []struct {
		currentTime string
		expected    bool
	}{
		{"2021-11-01T00:00:00Z", false},
		{"2021-11-15T23:59:59Z", false},
		{"2021-11-20T00:00:00Z", false},
		{"2021-12-15T00:00:00Z", true},
		{"2021-12-15T00:00:01Z", true},
		{"2021-12-30T00:00:00Z", true},
		{"2022-01-01T00:00:00Z", true},
	} {
		require.Equal(t, tc.expected, auction.ShouldAuctionFinished(types.MustParseRFC3339(tc.currentTime)))
	}
}
