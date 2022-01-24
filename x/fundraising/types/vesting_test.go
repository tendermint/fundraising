package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

func TestIsVestingReleasable(t *testing.T) {
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
			require.Equal(t, tc.expResult, tc.vq.IsVestingReleasable(now))
		})
	}
}
