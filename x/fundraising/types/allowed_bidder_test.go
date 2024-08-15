package types_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

func TestValidate_AllowedBidder(t *testing.T) {
	testBidderAddr := sdk.AccAddress(crypto.AddressHash([]byte("TestBidder")))

	testCases := []struct {
		allowedBidder types.AllowedBidder
		expectedErr   bool
	}{
		{
			types.NewAllowedBidder(1, testBidderAddr, math.NewInt(100_000_000)),
			false,
		},
		{
			types.NewAllowedBidder(1, sdk.AccAddress{}, math.NewInt(100_000_000)),
			true,
		},
		{
			types.NewAllowedBidder(1, testBidderAddr, math.NewInt(0)),
			true,
		},
		{
			types.NewAllowedBidder(1, testBidderAddr, math.ZeroInt()),
			true,
		},
	}

	for _, tc := range testCases {
		err := tc.allowedBidder.Validate()
		if tc.expectedErr {
			require.Error(t, err)
		} else {
			require.Equal(t, testBidderAddr, tc.allowedBidder.GetBidder())
			require.NoError(t, err)
		}
	}
}
