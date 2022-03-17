package types_test

import (
	"testing"
	time "time"

	"github.com/stretchr/testify/require"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

func TestUnpackAuction(t *testing.T) {
	auction := []types.AuctionI{
		types.NewFixedPriceAuction(
			types.NewBaseAuction(
				1,
				types.AuctionTypeFixedPrice,
				nil,
				sdk.AccAddress(crypto.AddressHash([]byte("Auctioneer"))).String(),
				types.SellingReserveAddress(1).String(),
				types.PayingReserveAddress(1).String(),
				sdk.MustNewDecFromStr("0.5"),
				sdk.NewInt64Coin("denom3", 1_000_000_000_000),
				"denom4",
				types.VestingReserveAddress(1).String(),
				[]types.VestingSchedule{},
				sdk.NewInt64Coin("denom3", 1_000_000_000_000),
				time.Now().AddDate(0, 0, -1),
				[]time.Time{time.Now().AddDate(0, 1, -1)},
				types.AuctionStatusStarted,
			),
		),
	}

	any, err := types.PackAuction(auction[0])
	require.NoError(t, err)

	marshaled, err := any.Marshal()
	require.NoError(t, err)

	var any2 codectypes.Any
	err = any2.Unmarshal(marshaled)
	require.NoError(t, err)

	reMarshal, err := any2.Marshal()
	require.NoError(t, err)
	require.Equal(t, marshaled, reMarshal)

	_, err = types.UnpackAuction(&any2)
	require.NoError(t, err)
}

func TestShouldAuctionStarted(t *testing.T) {
	auction := types.BaseAuction{
		Id:                    1,
		Type:                  types.AuctionTypeFixedPrice,
		AllowedBidders:        nil,
		Auctioneer:            sdk.AccAddress(crypto.AddressHash([]byte("Auctioneer"))).String(),
		SellingReserveAddress: types.SellingReserveAddress(1).String(),
		PayingReserveAddress:  types.PayingReserveAddress(1).String(),
		StartPrice:            sdk.MustNewDecFromStr("0.5"),
		SellingCoin:           sdk.NewInt64Coin("denom3", 1_000_000_000_000),
		PayingCoinDenom:       "denom4",
		VestingReserveAddress: types.VestingReserveAddress(1).String(),
		VestingSchedules:      []types.VestingSchedule{},
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
		SellingCoin:           sdk.NewInt64Coin("denom3", 1_000_000_000_000),
		PayingCoinDenom:       "denom4",
		VestingReserveAddress: types.VestingReserveAddress(1).String(),
		VestingSchedules:      []types.VestingSchedule{},
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
