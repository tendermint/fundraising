package types_test

import (
	"testing"
	time "time"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

func TestIsAuctionStarted(t *testing.T) {
	auction := types.NewFixedPriceAuction(
		&types.BaseAuction{
			Id:                    1,
			Type:                  types.AuctionTypeFixedPrice,
			Auctioneer:            sdk.AccAddress(crypto.AddressHash([]byte("Auctioneer"))).String(),
			SellingReserveAddress: types.SellingReserveAcc(1).String(),
			PayingReserveAddress:  types.PayingReserveAcc(1).String(),
			StartPrice:            sdk.MustNewDecFromStr("0.5"),
			SellingCoin:           sdk.NewInt64Coin("denom3", 1_000_000_000_000),
			PayingCoinDenom:       "denom4",
			VestingReserveAddress: types.VestingReserveAcc(1).String(),
			VestingSchedules:      []types.VestingSchedule{},
			WinningPrice:          sdk.ZeroDec(),
			RemainingCoin:         sdk.NewInt64Coin("denom3", 1_000_000_000_000),
			StartTime:             types.MustParseRFC3339("2021-12-01T00:00:00Z"),
			EndTimes:              []time.Time{types.MustParseRFC3339("2021-12-15T00:00:00Z")},
			Status:                types.AuctionStatusStandBy,
		},
	)

	for _, tc := range []struct {
		currentTime string
		expResult   bool
	}{
		{"2021-11-01T00:00:00Z", false},
		{"2021-11-15T23:59:59Z", false},
		{"2021-11-20T00:00:00Z", false},
		{"2021-12-01T00:00:00Z", true},
		{"2021-12-01T00:00:01Z", true},
		{"2021-12-10T00:00:00Z", true},
		{"2022-01-01T00:00:00Z", true},
	} {
		require.Equal(t, tc.expResult, auction.IsAuctionStarted(types.MustParseRFC3339(tc.currentTime)))
	}
}

func TestIsAuctionFinished(t *testing.T) {
	auction := types.NewFixedPriceAuction(
		&types.BaseAuction{
			Id:                    1,
			Type:                  types.AuctionTypeFixedPrice,
			Auctioneer:            sdk.AccAddress(crypto.AddressHash([]byte("Auctioneer"))).String(),
			SellingReserveAddress: types.SellingReserveAcc(1).String(),
			PayingReserveAddress:  types.PayingReserveAcc(1).String(),
			StartPrice:            sdk.MustNewDecFromStr("0.5"),
			SellingCoin:           sdk.NewInt64Coin("denom3", 1_000_000_000_000),
			PayingCoinDenom:       "denom4",
			VestingReserveAddress: types.VestingReserveAcc(1).String(),
			VestingSchedules:      []types.VestingSchedule{},
			WinningPrice:          sdk.ZeroDec(),
			RemainingCoin:         sdk.NewInt64Coin("denom3", 1_000_000_000_000),
			StartTime:             types.MustParseRFC3339("2021-12-01T00:00:00Z"),
			EndTimes:              []time.Time{types.MustParseRFC3339("2021-12-15T00:00:00Z")},
			Status:                types.AuctionStatusStandBy,
		},
	)

	for _, tc := range []struct {
		currentTime string
		expResult   bool
	}{
		{"2021-11-01T00:00:00Z", false},
		{"2021-11-15T23:59:59Z", false},
		{"2021-11-20T00:00:00Z", false},
		{"2021-12-15T00:00:00Z", true},
		{"2021-12-15T00:00:01Z", true},
		{"2021-12-30T00:00:00Z", true},
		{"2022-01-01T00:00:00Z", true},
	} {
		require.Equal(t, tc.expResult, auction.IsAuctionFinished(types.MustParseRFC3339(tc.currentTime)))
	}
}
