package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

func TestParams(t *testing.T) {
	defaultParams := types.DefaultParams()

	paramsStr := `auction_creation_fee:<denom:"stake" amount:"100000000" > extended_period:1 `
	require.Equal(t, paramsStr, defaultParams.String())
}

func TestParamsValidate(t *testing.T) {
	require.NoError(t, types.DefaultParams().Validate())

	testCases := []struct {
		name        string
		configure   func(*types.Params)
		expectedErr string
	}{
		{
			"EmptyAuctionCreationFee",
			func(params *types.Params) {
				params.AuctionCreationFee = sdk.NewCoins()
			},
			"",
		},
		{
			"EmptyPlaceBidFee",
			func(params *types.Params) {
				params.PlaceBidFee = sdk.NewCoins()
			},
			"",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			params := types.DefaultParams()
			tc.configure(&params)
			err := params.Validate()

			if tc.expectedErr != "" {
				require.EqualError(t, err, tc.expectedErr)
			} else {
				require.Nil(t, err)
			}
		})
	}
}
