package types_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/tendermint/fundraising/x/fundraising/types"
)

func TestConvertToSellingAmount(t *testing.T) {
	payingCoinDenom := "denom2" // auction paying coin denom

	testCases := []struct {
		bid         types.Bid
		expectedAmt math.Int
	}{
		{
			bid: types.Bid{
				Price: math.LegacyMustNewDecFromStr("0.5"),
				Coin:  sdk.NewCoin("denom1", math.NewInt(100_000)),
			},
			expectedAmt: math.NewInt(100_000),
		},
		{
			bid: types.Bid{
				Price: math.LegacyMustNewDecFromStr("0.5"),
				Coin:  sdk.NewCoin("denom2", math.NewInt(100_000)),
			},
			expectedAmt: math.NewInt(200_000),
		},
		{
			bid: types.Bid{
				Price: math.LegacyMustNewDecFromStr("0.1"),
				Coin:  sdk.NewCoin("denom1", math.NewInt(100_000)),
			},
			expectedAmt: math.NewInt(100_000),
		},
		{
			bid: types.Bid{
				Price: math.LegacyMustNewDecFromStr("0.1"),
				Coin:  sdk.NewCoin("denom2", math.NewInt(100_000)),
			},
			expectedAmt: math.NewInt(1_000_000),
		},
		{
			bid: types.Bid{
				Price: math.LegacyMustNewDecFromStr("3"),
				Coin:  sdk.NewCoin("denom2", math.NewInt(4)),
			},
			expectedAmt: math.NewInt(1),
		},
	}

	for _, tc := range testCases {
		sellingAmt := tc.bid.ConvertToSellingAmount(payingCoinDenom)
		require.Equal(t, tc.expectedAmt, sellingAmt)
	}
}

func TestConvertToPayingAmount(t *testing.T) {
	payingCoinDenom := "denom2" // auction paying coin denom

	testCases := []struct {
		bid         types.Bid
		expectedAmt math.Int
	}{
		{
			bid: types.Bid{
				Price: math.LegacyMustNewDecFromStr("0.5"),
				Coin:  sdk.NewCoin("denom1", math.NewInt(100_000)),
			},
			expectedAmt: math.NewInt(50_000),
		},
		{
			bid: types.Bid{
				Price: math.LegacyMustNewDecFromStr("0.5"),
				Coin:  sdk.NewCoin("denom2", math.NewInt(100_000)),
			},
			expectedAmt: math.NewInt(100_000),
		},
		{
			bid: types.Bid{
				Price: math.LegacyMustNewDecFromStr("0.1"),
				Coin:  sdk.NewCoin("denom1", math.NewInt(100_000)),
			},
			expectedAmt: math.NewInt(10_000),
		},
		{
			bid: types.Bid{
				Price: math.LegacyMustNewDecFromStr("0.1"),
				Coin:  sdk.NewCoin("denom2", math.NewInt(100_000)),
			},
			expectedAmt: math.NewInt(100_000),
		},
		{
			bid: types.Bid{
				Price: math.LegacyMustNewDecFromStr("0.33"),
				Coin:  sdk.NewCoin("denom1", math.NewInt(100_000)),
			},
			expectedAmt: math.NewInt(33000),
		},
	}

	for _, tc := range testCases {
		payingAmt := tc.bid.ConvertToPayingAmount(payingCoinDenom)
		require.Equal(t, tc.expectedAmt, payingAmt)
	}
}

func TestSetMatched(t *testing.T) {
	bidder := sdk.AccAddress(crypto.AddressHash([]byte("Bidder")))

	bid := types.NewBid(
		1,
		bidder,
		1,
		types.BidTypeFixedPrice,
		math.LegacyMustNewDecFromStr("0.5"),
		sdk.NewCoin("denom1", math.NewInt(100_000)),
		false,
	)
	require.False(t, bid.IsMatched)
	require.Equal(t, bidder, bid.GetBidder())

	bid.SetMatched(true)
	require.True(t, bid.IsMatched)
}
